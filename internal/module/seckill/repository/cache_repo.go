package repository

import (
	"bytes"
	"context"
	"encoding/gob"
	"reflect"
	"time"

	seckilldto "github.com/algorithm9/flash-deal/internal/module/seckill/dto"
	"github.com/algorithm9/flash-deal/internal/shared/keys"
	"github.com/algorithm9/flash-deal/internal/shared/redisclient"
	"github.com/algorithm9/flash-deal/pkg/cache"
	"github.com/algorithm9/flash-deal/pkg/errorx"
	"github.com/algorithm9/flash-deal/pkg/logger"
)

type Cache interface {
	SetSkuDetail(ctx context.Context, skuID uint64, detail *seckilldto.SeckillSkuDetail) error
	GetSkuDetail(ctx context.Context, skuID uint64) (*seckilldto.SeckillSkuDetail, error)
	SetResult(ctx context.Context, userID, activityID, skuID uint64, result string) error
	GetResult(ctx context.Context, userID, activityID, skuID uint64) (string, error)
}

func NewCache(c cache.Cache, rdb *redisclient.Client) Cache {
	return &cacheImp{cache: c, rdb: rdb}
}

type cacheImp struct {
	cache cache.Cache
	rdb   *redisclient.Client
}

func (c *cacheImp) SetSkuDetail(ctx context.Context, skuID uint64, detail *seckilldto.SeckillSkuDetail) error {
	key := keys.SKUDetailKey(skuID)
	data, err := c.structToBytes(detail)
	if err != nil {
		return errorx.WrapErr(err, "serialize failed")
	}

	return c.SetData(ctx, key, data)
}

func (c *cacheImp) GetSkuDetail(ctx context.Context, skuID uint64) (*seckilldto.SeckillSkuDetail, error) {
	key := keys.SKUDetailKey(skuID)
	ids := map[string]uint64{"skuID": skuID}

	decode := func(data []byte) (interface{}, error) {
		detail := new(seckilldto.SeckillSkuDetail)
		if err := c.bytesToStruct(data, detail); err != nil {
			return nil, err
		}
		return detail, nil
	}

	result, err := c.getWithFallback(ctx, key, decode, ids)
	if err != nil {
		return nil, err
	}
	return result.(*seckilldto.SeckillSkuDetail), nil
}

func (c *cacheImp) structToBytes(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(obj)
	return buf.Bytes(), err
}

func (c *cacheImp) bytesToStruct(data []byte, obj interface{}) error {
	if len(data) == 0 {
		return errorx.NewErrMsg("empty data")
	}
	if reflect.ValueOf(obj).Kind() != reflect.Ptr || reflect.ValueOf(obj).IsNil() {
		return errorx.NewErrMsg("invalid object pointer")
	}
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(obj)
}

func (c *cacheImp) SetResult(ctx context.Context, userID, activityID, skuID uint64, result string) error {
	key := keys.SeckillOrderedKey(userID, activityID, skuID)
	data := []byte(result)

	return c.SetData(ctx, key, data)
}

func (c *cacheImp) GetResult(ctx context.Context, userID, activityID, skuID uint64) (string, error) {
	key := keys.SeckillOrderedKey(userID, activityID, skuID)
	ids := map[string]uint64{"userID": userID, "activityID": activityID, "skuID": skuID}

	decode := func(data []byte) (interface{}, error) {
		return string(data), nil
	}

	result, err := c.getWithFallback(ctx, key, decode, ids)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (c *cacheImp) SetData(ctx context.Context, key string, data []byte) error {
	if err := c.rdb.Client.Set(ctx, key, data, time.Hour).Err(); err != nil {
		return errorx.WrapErr(err, "redis set failed")
	}

	if err := c.cache.Set(key, data, time.Hour); err != nil {
		logger.L().Err(err).Str("key", key).Msg("failed to set local cache")
	}
	return nil
}

func (c *cacheImp) getWithFallback(
	ctx context.Context,
	key string,
	decodeFunc func([]byte) (interface{}, error),
	ids map[string]uint64,
) (interface{}, error) {

	if data, err := c.cache.Get(key); err == nil {
		result, err := decodeFunc(data)
		if err == nil {
			return result, nil
		}
		logger.L().Err(err).Fields(logFields(ids)).Msg("failed to decode local cache")
	}

	data, err := c.rdb.Client.Get(ctx, key).Bytes()
	if err != nil {
		if redisclient.IsNil(err) {
			return nil, errorx.NewErrMsg("data not found")
		}
		return nil, errorx.WrapErrMsgF(err, "key:%s,failed to get redis cache", key)
	}

	result, err := decodeFunc(data)
	if err != nil {
		return nil, errorx.WrapErrMsgF(err, "key:%s,failed to decode data", key)
	}

	go func() {
		if err := c.cache.Set(key, data, time.Hour); err != nil {
			logger.L().Err(err).Fields(logFields(ids)).Msg("failed to backfill local cache")
		}
	}()

	return result, nil
}

// 辅助函数：生成日志字段
func logFields(ids map[string]uint64) map[string]interface{} {
	fields := make(map[string]interface{})
	for k, v := range ids {
		fields[k] = v
	}
	return fields
}
