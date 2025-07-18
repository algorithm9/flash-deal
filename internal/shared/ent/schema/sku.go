// ent/schema/sku.go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	"github.com/algorithm9/flash-deal/internal/shared/ent/schema/mixin"
)

// SKU holds the SKU info.
type SKU struct {
	ent.Schema
}

func (SKU) Annotations() []schema.Annotation {
	return append(mixin.DefaultAnnotations, entsql.Annotation{Table: "skus"})
}

func (SKU) Mixin() []ent.Mixin {
	return mixin.DefaultMixin
}

func (SKU) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("product_id").Comment("所属商品ID"),
		field.JSON("specs", map[string]string{}).Comment("属性组合，如颜色和容量"),
		field.Float("price").Comment("原价"),
		field.Int("stock").Comment("库存"),
	}
}

func (SKU) Edges() []ent.Edge {
	return []ent.Edge{}
}
