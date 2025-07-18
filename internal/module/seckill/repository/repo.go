package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	seckilldto "github.com/algorithm9/flash-deal/internal/module/seckill/dto"
	"github.com/algorithm9/flash-deal/internal/shared/ent/gen"
	"github.com/algorithm9/flash-deal/internal/shared/ent/gen/seckillactivity"
	"github.com/algorithm9/flash-deal/internal/shared/redisclient"
	"github.com/algorithm9/flash-deal/pkg/errorx"
)

type SeckillRepository interface {
	FindSeckillActivities(ctx context.Context) ([]*gen.SeckillActivity, error)
	FindSeckillActivityDetail(ctx context.Context, skuID uint64) (*seckilldto.SeckillSkuDetail, error)
}

type seckillRepositoryImp struct {
	db *gen.Client
}

func NewSeckillRepository(db *gen.Client) SeckillRepository {
	return &seckillRepositoryImp{
		db: db,
	}
}

func (r *seckillRepositoryImp) FindSeckillActivities(ctx context.Context) ([]*gen.SeckillActivity, error) {
	now := time.Now().UTC()
	all, err := r.db.SeckillActivity.Query().
		Where(seckillactivity.StartTimeLTE(now),
			seckillactivity.EndTimeGTE(now)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (r *seckillRepositoryImp) FindSeckillActivityDetail(ctx context.Context, skuID uint64) (*seckilldto.SeckillSkuDetail, error) {
	skuDetail, err := r.db.QueryContext(ctx, skuDetailSql(), skuID)
	if err != nil {
		return nil, errorx.Wrap(seckilldto.ActivityInternalError.Int(), http.StatusInternalServerError, "failed to query sku", err)
	}
	defer skuDetail.Close()

	var s seckilldto.SeckillSkuDetail
	var specsRaw []byte
	for skuDetail.Next() {
		if err = skuDetail.Scan(&s.ID, &s.Title,
			&s.Description, &specsRaw); err != nil {
			return nil, errorx.Wrap(seckilldto.ActivityInternalError.Int(), http.StatusInternalServerError, "failed to query sku", err)
		}
	}
	if err = skuDetail.Err(); err != nil {
		return nil, errorx.Wrap(seckilldto.ActivityInternalError.Int(), http.StatusInternalServerError, "failed to query sku", err)
	}

	if s.ID == 0 {
		return nil, nil
	}
	var specs map[string]string
	if err := json.Unmarshal(specsRaw, &specs); err != nil {
		return nil, err
	}

	s.Specs = specs

	return &s, nil
}

func skuDetailSql() string {
	return `
SELECT 
    s.id,
  	p.title,
  	p.description,
  	s.specs
FROM skus s
JOIN products p ON s.product_id = p.id
WHERE s.id = ?
`
}

type LuaRepo interface {
	DecrStock(ctx context.Context, activityID, skuID uint64) (bool, error)
	Exists(ctx context.Context, key string) (bool, error)
}

type QueueProducer interface {
	Send(ctx context.Context, msg seckilldto.SeckillRequest) error
}

type luaRepoImpl struct {
	rdb *redisclient.Client
	lua *redis.Script
}

func NewLuaRepo(rdb *redisclient.Client) LuaRepo {
	return &luaRepoImpl{
		rdb: rdb,
		lua: redis.NewScript(`
			local stock = tonumber(redis.call("GET", KEYS[1]))
			if not stock or stock <= 0 then
				return 0
			end
			redis.call("DECR", KEYS[1])
			return 1
		`),
	}
}

func (r *luaRepoImpl) DecrStock(ctx context.Context, activityID, skuID uint64) (bool, error) {
	key := fmt.Sprintf("sku:stock:%d:%d", activityID, skuID)
	res, err := r.lua.Run(ctx, r.rdb.Client, []string{key}).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

func (r *luaRepoImpl) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.rdb.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
