package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goodsdto "github.com/algorithm9/flash-deal/internal/module/product/dto"
	"github.com/algorithm9/flash-deal/internal/module/product/repository"
	"github.com/algorithm9/flash-deal/internal/shared/redisclient"
	"github.com/algorithm9/flash-deal/pkg/logger"
)

type GoodsService interface {
	GetSKUDetail(ctx context.Context, skuID uint64) (*goodsdto.SKUProductDetail, error)
}

type goodsServiceImp struct {
	repo repository.GoodsRepository
	rdb  *redisclient.Client
}

func NewGoodsService(repo repository.GoodsRepository, rdb *redisclient.Client) GoodsService {
	return &goodsServiceImp{repo: repo, rdb: rdb}
}

func (s *goodsServiceImp) GetSKUDetail(ctx context.Context, skuID uint64) (*goodsdto.SKUProductDetail, error) {
	cacheKey := fmt.Sprintf("goods:sku:%d", skuID)

	// 读取缓存
	cached, err := s.rdb.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中
		if cached == "null" {
			// 缓存了空值
			return nil, nil
		}
		var detail goodsdto.SKUProductDetail
		if err := json.Unmarshal([]byte(cached), &detail); err == nil {
			return &detail, nil
		}
		// 反序列化失败：可记录日志但不直接返回错误
		logger.L().Info().Err(err).Msgf("failed to unmarshal cached SKU detail")
	} else if !redisclient.IsNil(err) {
		// 其他 Redis 错误：记录日志
		logger.L().Info().Err(err).Msgf("redis get error")
	}

	// 缓存未命中或反序列化失败，查询数据库
	detail, err := s.repo.GetSKUWithProduct(ctx, skuID)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		// 缓存空值防止穿透
		if err := s.rdb.Client.Set(ctx, cacheKey, "null", time.Minute*1).Err(); err != nil {
			logger.L().Info().Err(err).Msgf("failed to cache null for sku %d", skuID)
		}
		return nil, nil
	}

	// 序列化并缓存
	if b, err := json.Marshal(detail); err == nil {
		if err := s.rdb.Client.Set(ctx, cacheKey, b, time.Minute*5).Err(); err != nil {
			logger.L().Info().Err(err).Msg("failed to cache sku detail")
		}
	} else {
		logger.L().Info().Err(err).Msg("failed to marshal sku detail")
	}

	return detail, nil
}
