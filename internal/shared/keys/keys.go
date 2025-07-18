package keys

import "fmt"

// SeckillStockKey 商品库存缓存 key（秒杀库存）
func SeckillStockKey(skuID uint64) string {
	return fmt.Sprintf("seckill:stock:%d", skuID)
}

// SKUDetailKey 商品详情缓存 key（SKU 信息）
func SKUDetailKey(skuID uint64) string {
	return fmt.Sprintf("sku:detail:%d", skuID)
}

// SeckillOrderedKey 秒杀是否已下单记录（防止重复下单）
func SeckillOrderedKey(userID, activityID, skuID uint64) string {
	return fmt.Sprintf("seckill:ordered:%d:%d:%d", userID, activityID, skuID)
}

// SeckillTokenKey 用户秒杀令牌 key（令牌桶或验证码）
func SeckillTokenKey(skuID, userID uint64) string {
	return fmt.Sprintf("seckill:token:%d:%d", skuID, userID)
}

// RateLimitUserKey 用户限流 key（单位时间内请求次数）
func RateLimitUserKey(userID uint64) string {
	return fmt.Sprintf("limit:user:%d", userID)
}

// RateLimitIPKey IP 限流 key（防止刷请求）
func RateLimitIPKey(ip string) string {
	return fmt.Sprintf("limit:ip:%s", ip)
}

// UserTokenKey 用户登录 token 缓存（token 验证）
func UserTokenKey(userID uint64) string {
	return fmt.Sprintf("user:token:%d", userID)
}

// SeckillUserLockKey 用户抢购锁 key（防止并发请求）
func SeckillUserLockKey(skuID, userID uint64) string {
	return fmt.Sprintf("seckill:lock:%d:%d", skuID, userID)
}
