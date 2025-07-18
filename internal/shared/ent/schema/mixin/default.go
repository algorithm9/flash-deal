package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
)

var DefaultMixin = []ent.Mixin{
	TimeMixin{},
	IDMixin{},
}

var DefaultAnnotations = []schema.Annotation{
	entsql.Annotation{
		Charset:   "utf8mb4",
		Collation: "utf8mb4_0900_ai_ci",
	},
}
