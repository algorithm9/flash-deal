package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	"github.com/algorithm9/flash-deal/internal/shared/ent/schema/mixin"
)

// Product holds the SPU info.
type Product struct {
	ent.Schema
}

func (Product) Annotations() []schema.Annotation {
	return mixin.DefaultAnnotations
}

func (Product) Mixin() []ent.Mixin {
	return mixin.DefaultMixin
}

func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").Comment("商品标题"),
		field.String("description").Optional().Nillable().Comment("商品描述"),
	}
}

func (Product) Edges() []ent.Edge {
	return []ent.Edge{}
}
