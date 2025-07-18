package goodsdto

type SKUProductDetail struct {
	SKUID        uint64            `json:"sku_id"`
	ProductID    uint64            `json:"product_id"`
	ProductTitle string            `json:"title"`
	Description  string            `json:"description"`
	Specs        map[string]string `json:"specs"`
	Price        float64           `json:"price"`
	SeckillPrice float64           `json:"seckill_price,omitempty"`
	Stock        int               `json:"stock"`
	IsSeckill    bool              `json:"is_seckill"`
	IsActive     bool              `json:"is_active"`
}

type GoodsModuleResponseCode int

const (
	SKUNotFound      GoodsModuleResponseCode = 20001
	GoodsNotFound    GoodsModuleResponseCode = 20002
	SKUInternalError GoodsModuleResponseCode = 2003
)

func (r GoodsModuleResponseCode) Int() int {
	return int(r)
}
