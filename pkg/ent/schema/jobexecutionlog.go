package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/schema/mixin"
)

type JobExecutionLog struct {
	ent.Schema
}

func (JobExecutionLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "job_execution_logs"},
	}
}

func (JobExecutionLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (JobExecutionLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("instance_no").MaxLen(64),
		field.String("shard_no").MaxLen(64).Optional().Nillable(),
		field.Int64("executor_id").Optional().Nillable(),
		field.String("log_level").MaxLen(16),
		field.String("phase").MaxLen(32),
		field.String("message").MaxLen(1024),
		field.String("trace_id").MaxLen(64).Optional().Nillable(),
	}
}

func (JobExecutionLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("instance_no"),
		index.Fields("shard_no"),
	}
}
