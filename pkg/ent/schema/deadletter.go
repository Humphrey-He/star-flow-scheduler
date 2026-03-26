package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/schema/mixin"
)

type DeadLetter struct {
	ent.Schema
}

func (DeadLetter) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "dead_letters"},
	}
}

func (DeadLetter) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (DeadLetter) Fields() []ent.Field {
	return []ent.Field{
		field.String("instance_no").MaxLen(64).Comment("Instance number"),
		field.Int64("job_id"),
		field.String("dead_reason").MaxLen(64),
		field.String("dead_message").MaxLen(1024).Optional().Nillable(),
		field.String("retry_snapshot").Optional().Nillable(),
		field.Int64("last_executor_id").Optional().Nillable(),
		field.String("status").MaxLen(32).Default("open"),
	}
}

func (DeadLetter) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status", "dead_reason"),
		index.Fields("instance_no"),
	}
}
