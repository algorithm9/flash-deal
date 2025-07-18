package seckilldto

import (
	"time"

	"github.com/algorithm9/flash-deal/internal/shared/ent/gen"
)

type SeckillStatus string

const (
	Absence  SeckillStatus = "absence"
	Queueing SeckillStatus = "queueing"
	Success  SeckillStatus = "success"
	Fail     SeckillStatus = "fail"
)

func (s SeckillStatus) String() string {
	return string(s)
}

type SeckillRequest struct {
	UserID     uint64
	ActivityID uint64
	SKUID      uint64
}

type SeckillActivities struct {
	List []*SeckillActivity `json:"list"`
}

func ConvertToDTO(seckillActivities []*gen.SeckillActivity) *SeckillActivities {
	var result SeckillActivities
	for _, sa := range seckillActivities {
		result.List = append(result.List, &SeckillActivity{
			ID:           sa.ID,
			ActivityName: sa.Name,
			SKUID:        sa.SkuID,
			SeckillPrice: sa.SeckillPrice,
			Price:        sa.Price,
			Stock:        sa.Stock,
			StartTime:    sa.StartTime,
			EndTime:      sa.EndTime,
		})
	}
	return &result
}

type SeckillActivity struct {
	ID           uint64    `json:"id"`
	ActivityName string    `json:"name"` // activity name
	SKUID        uint64    `json:"sku_id"`
	SeckillPrice float64   `json:"seckill_price"`
	Price        float64   `json:"price"`
	Stock        int       `json:"stock"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}

type SeckillSkuDetail struct {
	ID          uint64            `json:"id"`
	Title       string            `json:"title"`       // 商品名称
	Description string            `json:"description"` // 商品描述
	Specs       map[string]string `json:"specs"`       // 商品属性
}

type SeckillModuleResponseCode int

const (
	SKUNotFound                SeckillModuleResponseCode = 30001
	GoodsNotFound              SeckillModuleResponseCode = 30002
	ActivityInternalError      SeckillModuleResponseCode = 30003
	InsufficientInventoryError SeckillModuleResponseCode = 30004
)

func (r SeckillModuleResponseCode) Int() int {
	return int(r)
}
