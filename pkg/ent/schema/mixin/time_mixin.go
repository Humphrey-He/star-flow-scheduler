package mixin

import (
    "time"

    "entgo.io/ent"
    entmixin "entgo.io/ent/schema/mixin"
    "entgo.io/ent/schema/field"
)

type TimeMixin struct {
    entmixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
    return []ent.Field{
        field.Time("created_at").Default(time.Now).Immutable(),
        field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
        field.Time("deleted_at").Optional().Nillable(),
    }
}
