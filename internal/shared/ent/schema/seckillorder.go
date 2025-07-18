package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/algorithm9/flash-deal/internal/shared/ent/schema/mixin"
)

// SeckillOrder holds the schema definition for the SeckillOrder entity.
type SeckillOrder struct {
	ent.Schema
}

func (SeckillOrder) Annotations() []schema.Annotation {
	return mixin.DefaultAnnotations
}

func (SeckillOrder) Mixin() []ent.Mixin {
	return mixin.DefaultMixin
}

func (SeckillOrder) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("user_id").Comment("用户ID"),
		field.Uint64("activity_id").Comment("活动ID"),
		field.Uint64("sku_id").Comment("sku ID"),
		field.Float("price").Comment("秒杀价"),
		field.Int("status").Default(0).Comment("状态: 0=待支付,1=已支付,2=超时"),
	}
}

func (SeckillOrder) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("activity_id"),
		index.Fields("sku_id"),
		index.Fields("user_id", "activity_id", "sku_id").Unique(),
	}
}
