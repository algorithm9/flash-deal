package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/algorithm9/flash-deal/internal/shared/redisclient"
)

type RateLimiter struct {
	client *redisclient.Client
}

func NewRateLimiter(client *redisclient.Client) *RateLimiter {
	return &RateLimiter{client: client}
}

// UserLimit 用户级限流
func (l *RateLimiter) UserLimit(ctx context.Context, userID, activityID uint64, interval time.Duration) bool {
	key := fmt.Sprintf("limit:user:%d:%d", activityID, userID)
	return l.checkLimit(ctx, key, interval)
}

// IPLimit IP级限流
func (l *RateLimiter) IPLimit(ctx context.Context, ip string, activityID uint64, interval time.Duration) bool {
	key := fmt.Sprintf("limit:ip:%s:%d", ip, activityID)
	return l.checkLimit(ctx, key, interval)
}

// ActivityLimit 活动级限流
func (l *RateLimiter) ActivityLimit(ctx context.Context, activityID uint64, maxQPS int) bool {
	key := fmt.Sprintf("limit:activity:%d", activityID)
	now := time.Now().UnixMilli()
	minTime := now - 1000 // 1秒窗口

	// 使用ZSET实现滑动窗口限流
	l.client.Client.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", minTime))

	count, err := l.client.Client.ZCard(ctx, key).Result()
	if err != nil {
		return true
	}

	if count >= int64(maxQPS) {
		return true
	}

	l.client.Client.ZAdd(ctx, key, &redis.Z{
		Score:  float64(now),
		Member: now,
	})
	l.client.Client.Expire(ctx, key, 2*time.Second)

	return false
}

func (l *RateLimiter) checkLimit(ctx context.Context, key string, interval time.Duration) bool {
	val, err := l.client.Client.Get(ctx, key).Int64()
	if err != nil && !redisclient.IsNil(err) {
		return true
	}

	if val > 0 {
		return true
	}

	l.client.Client.Set(ctx, key, 1, interval)
	return false
}
