package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/schema/mixin"
)

type JobInstance struct {
	ent.Schema
}

func (JobInstance) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "job_instances"},
	}
}

func (JobInstance) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (JobInstance) Fields() []ent.Field {
	return []ent.Field{
		field.String("instance_no").MaxLen(64).Unique().Comment("Instance number"),
		field.Int64("job_id"),
		field.Int64("workflow_id").Optional().Nillable(),
		field.String("trigger_type").MaxLen(32).Comment("cron/manual/delay/event/retry"),
		field.Time("trigger_time"),
		field.Time("scheduled_time"),
		field.Time("dispatch_time").Optional().Nillable(),
		field.Time("start_time").Optional().Nillable(),
		field.Time("finish_time").Optional().Nillable(),
		field.String("status").MaxLen(32).Comment("pending/dispatched/running/success/failed/retry_wait/dead/canceled"),
		field.Int("retry_count").Default(0),
		field.Int64("current_backoff_ms").Default(0),
		field.Int64("executor_id").Optional().Nillable(),
		field.Int("shard_total").Default(1),
		field.Int("success_shards").Default(0),
		field.Int("failed_shards").Default(0),
		field.String("payload").Optional().Nillable().Comment("Runtime payload"),
		field.String("result_summary").MaxLen(512).Optional().Nillable(),
		field.String("error_code").MaxLen(64).Optional().Nillable(),
		field.String("error_message").MaxLen(1024).Optional().Nillable(),
	}
}

func (JobInstance) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("job_id", "trigger_time"),
		index.Fields("status", "scheduled_time"),
		index.Fields("executor_id", "status"),
	}
}
