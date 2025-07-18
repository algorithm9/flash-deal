package mixin

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(NowWithUTC).
			Annotations(
				entsql.Default("CURRENT_TIMESTAMP"),
			).
			Immutable(),
		field.Time("updated_at").
			Default(NowWithUTC).
			Annotations(
				entsql.Default("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"),
			).
			Immutable(),
	}
}

type IDMixin struct {
	mixin.Schema
}

var _ ent.Mixin = IDMixin{}

func (IDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").
			Immutable().
			Unique(),
	}
}

func NowWithUTC() time.Time {
	return time.Now().UTC()
}
