package mixin

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type TimeMixin struct {
	ent.Mixin
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}
