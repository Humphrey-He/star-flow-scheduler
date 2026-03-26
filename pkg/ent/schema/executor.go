package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/schema/mixin"
)

type Executor struct {
	ent.Schema
}

func (Executor) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "executors"},
	}
}

func (Executor) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (Executor) Fields() []ent.Field {
	return []ent.Field{
		field.String("executor_code").MaxLen(64).Unique(),
		field.String("host").MaxLen(128),
		field.String("ip").MaxLen(64),
		field.String("grpc_addr").MaxLen(128),
		field.String("http_addr").MaxLen(128).Optional().Nillable(),
		field.String("tags").MaxLen(256).Optional().Nillable().Comment("Comma separated tags"),
		field.Int("capacity").Default(100),
		field.Int("current_load").Default(0),
		field.String("version").MaxLen(32).Optional().Nillable(),
		field.String("status").MaxLen(32).Default("online"),
		field.Time("last_heartbeat_at"),
		field.JSON("metadata", map[string]interface{}{}).Optional().Nillable(),
	}
}

func (Executor) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status", "last_heartbeat_at"),
	}
}
