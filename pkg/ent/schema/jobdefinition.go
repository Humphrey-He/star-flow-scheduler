package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/schema/mixin"
)

type JobDefinition struct {
	ent.Schema
}

func (JobDefinition) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "job_definitions"},
	}
}

func (JobDefinition) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (JobDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.String("job_code").MaxLen(64).Unique().Comment("Task unique code"),
		field.String("job_name").MaxLen(128).Comment("Task name"),
		field.String("job_type").MaxLen(32).Comment("cron/delay/once/dag"),
		field.String("schedule_expr").MaxLen(128).Optional().Nillable().Comment("Cron expression"),
		field.Int64("delay_ms").Optional().Nillable().Comment("Delay in milliseconds"),
		field.String("execute_mode").MaxLen(32).Default("standalone").Comment("standalone/sharding/dag"),
		field.String("handler_name").MaxLen(128).Default("").Comment("Executor handler name"),
		field.String("handler_payload").Optional().Nillable().Comment("Default task payload"),
		field.Int("timeout_ms").Default(60000),
		field.Int("retry_limit").Default(3),
		field.String("retry_backoff").Default("1s,3s,5s").Comment("Retry backoff settings"),
		field.Int("priority").Default(5),
		field.Int("shard_total").Default(1),
		field.String("route_strategy").MaxLen(32).Default("least_load").Comment("least_load/hash/tag"),
		field.String("executor_tag").MaxLen(64).Optional().Nillable().Comment("Executor tag"),
		field.String("idempotent_key_expr").MaxLen(128).Optional().Nillable().Comment("Idempotent key expression"),
		field.String("status").MaxLen(32).Default("enabled").Comment("enabled/disabled"),
		field.String("created_by").MaxLen(64).Optional().Nillable(),
		field.String("updated_by").MaxLen(64).Optional().Nillable(),
	}
}

func (JobDefinition) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("job_type", "status"),
		index.Fields("handler_name"),
	}
}
