package service

import (
	"context"
	"net/http"

	seckilldto "github.com/algorithm9/flash-deal/internal/module/seckill/dto"
	"github.com/algorithm9/flash-deal/internal/module/seckill/repository"
	"github.com/algorithm9/flash-deal/pkg/errorx"
	"github.com/algorithm9/flash-deal/pkg/logger"
)

type Service interface {
	GetSeckillActivities(ctx context.Context) (*seckilldto.SeckillActivities, error)
	GetSeckillActivityDetail(ctx context.Context, id uint64) (*seckilldto.SeckillSkuDetail, error)
	Seckill(ctx context.Context, userID, activityID, skuID uint64) error
	GetResult(ctx context.Context, userID, activityID, skuID uint64) (string, error)
}

type serviceImpl struct {
	db   repository.SeckillRepository
	repo repository.LuaRepo
	q    repository.QueueProducer
	c    repository.Cache
}

func New(repo repository.LuaRepo,
	db repository.SeckillRepository,
	q repository.QueueProducer,
	c repository.Cache) Service {
	return &serviceImpl{repo: repo, db: db, q: q, c: c}
}

func (s *serviceImpl) GetSeckillActivities(ctx context.Context) (*seckilldto.SeckillActivities, error) {
	activities, err := s.db.FindSeckillActivities(ctx)
	if err != nil {
		return nil, errorx.Wrap(seckilldto.ActivityInternalError.Int(), http.StatusInternalServerError, "failed to query activity", err)
	}
	return seckilldto.ConvertToDTO(activities), nil
}

func (s *serviceImpl) GetSeckillActivityDetail(ctx context.Context, skuID uint64) (skuDetail *seckilldto.SeckillSkuDetail, err error) {

	skuDetail, err = s.c.GetSkuDetail(ctx, skuID)
	if err != nil {
		logger.L().Err(err).Uint64("skuID", skuID).Msg("failed to get sku detail from cache")
	}
	if skuDetail != nil {
		return
	}

	skuDetail, err = s.db.FindSeckillActivityDetail(ctx, skuID)
	if err != nil {
		return nil, errorx.Wrap(seckilldto.ActivityInternalError.Int(), http.StatusInternalServerError, "failed to query sku", err)
	}

	if skuDetail != nil {
		err := s.c.SetSkuDetail(ctx, skuID, skuDetail)
		if err != nil {
			logger.L().Err(err).Uint64("skuID", skuID).Msg("failed to set sku detail to cache")
		}
	}

	return skuDetail, nil
}

func (s *serviceImpl) Seckill(ctx context.Context, userID, activityID, skuID uint64) error {
	success, err := s.repo.DecrStock(ctx, activityID, skuID)
	if err != nil {
		return err
	}
	if !success {
		_ = s.c.SetResult(ctx, userID, activityID, skuID, seckilldto.Fail.String())
		return errorx.New(seckilldto.InsufficientInventoryError.Int(), http.StatusInternalServerError, "insufficient inventory")
	}
	err = s.q.Send(ctx, seckilldto.SeckillRequest{UserID: userID, ActivityID: activityID, SKUID: skuID})
	if err == nil {

	}
	return err
}

func (s *serviceImpl) GetResult(ctx context.Context, userID, activityID, skuID uint64) (string, error) {
	getResult, err := s.c.GetResult(ctx, userID, activityID, skuID)
	if err != nil {
		return "", errorx.WrapErrMsgF(err, "failed to get result from cache")
	}
	if getResult == "" {
		return seckilldto.Queueing.String(), nil
	}
	return getResult, nil
}
