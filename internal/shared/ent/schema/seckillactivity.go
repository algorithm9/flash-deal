package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/algorithm9/flash-deal/internal/shared/ent/schema/mixin"
)

// SeckillActivity holds the schema definition for the SeckillActivity entity.
type SeckillActivity struct {
	ent.Schema
}

func (SeckillActivity) Annotations() []schema.Annotation {
	return mixin.DefaultAnnotations
}

func (SeckillActivity) Mixin() []ent.Mixin {
	return mixin.DefaultMixin
}

// Fields of the SeckillActivity.
func (SeckillActivity) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("活动名称"),
		field.Uint64("sku_id").Comment("商品ID"),
		field.Float("price").Comment("原价"),
		field.Float("seckill_price").Comment("秒杀价格"),
		field.Int("stock").Comment("秒杀库存"),
		field.Time("start_time").Comment("开始时间"),
		field.Time("end_time").Comment("结束时间"),
	}
}

func (SeckillActivity) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sku_id"),
		index.Fields("start_time"),
		index.Fields("end_time"),
	}
}
