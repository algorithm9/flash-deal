package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	"github.com/algorithm9/flash-deal/internal/shared/ent/schema/mixin"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return mixin.DefaultAnnotations
}

func (User) Mixin() []ent.Mixin {
	return mixin.DefaultMixin
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("phone").
			Unique(),
		field.String("password_hash").
			Sensitive(),
		field.Enum("status").
			Values("active", "locked", "deleted").
			Default("active"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{}
}
