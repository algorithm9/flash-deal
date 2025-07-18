package kafka

import (
	"context"
	"fmt"
	"strconv"
	"time"

	db "github.com/algorithm9/flash-deal/cmd/worker/workers/db/models"
	seckilldto "github.com/algorithm9/flash-deal/internal/module/seckill/dto"
	"github.com/algorithm9/flash-deal/pkg/errorx"
)

func (kc *Consumer) handleOrder(ctx context.Context, req seckilldto.SeckillRequest) error {
	if kc.hasOrdered(ctx, req.UserID, req.ActivityID, req.SKUID) {
		return errorx.NewErrMsgF("Block repeat order, user=%d, activity=%d, sku=%d", req.UserID, req.ActivityID, req.SKUID)
	}

	price, err := kc.getPrice(ctx, req.ActivityID, req.SKUID)
	if err != nil {
		return errorx.WrapErr(err, "Failed to get price")
	}

	if err := kc.insertOrderAndDecrDB(ctx, req.UserID, req.ActivityID, req.SKUID, price); err != nil {
		return errorx.WrapErr(err, "Failed to insert order")
	}

	kc.markOrdered(ctx, req.UserID, req.ActivityID, req.SKUID)
	return nil
}

func (kc *Consumer) getPrice(ctx context.Context, activityID, skuID uint64) (float64, error) {
	key := fmt.Sprintf("seckill:price:%d:%d", activityID, skuID)
	val, err := kc.rdb.Client.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(val, 64)
}

func (kc *Consumer) hasOrdered(ctx context.Context, userID, activityID, skuID uint64) bool {
	key := fmt.Sprintf("seckill:ordered:%d:%d:%d", userID, activityID, skuID)
	exists, _ := kc.rdb.Client.Exists(ctx, key).Result()
	return exists > 0
}

func (kc *Consumer) markOrdered(ctx context.Context, userID, activityID, skuID uint64) {
	key := fmt.Sprintf("seckill:ordered:%d:%d:%d", userID, activityID, skuID)
	kc.rdb.Client.SetNX(ctx, key, 1, 24*time.Hour)
}

func (kc *Consumer) insertOrderAndDecrDB(ctx context.Context, userID, activityID, skuID uint64, price float64) error {
	tx, err := kc.db.BeginTx(ctx, nil)
	if err != nil {
		return errorx.WrapErr(err, "failed to begin transaction")
	}

	// 1. decrease stock
	result, err := kc.queries.WithTx(tx).DecrStock(ctx, skuID)
	if err != nil {
		tx.Rollback()
		return errorx.WrapErr(err, "failed to decr stock")
	}

	// check stock
	if result == 0 {
		tx.Rollback()
		return errorx.NewErrMsg(fmt.Sprintf("activity id: %d, sku id: %d,insufficient inventory", activityID, skuID))
	}

	// 2. create order
	err = kc.queries.WithTx(tx).InsertOrder(ctx, db.InsertOrderParams{
		UserID:     userID,
		ActivityID: activityID,
		SkuID:      skuID,
		Price:      price,
	})
	if err != nil {
		tx.Rollback()
		return errorx.WrapErr(err, "failed to insert order")
	}

	if err := tx.Commit(); err != nil {
		return errorx.WrapErr(err, "failed to commit transaction")
	}
	return nil
}
