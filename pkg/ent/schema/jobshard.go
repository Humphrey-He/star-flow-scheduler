package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/schema/mixin"
)

type JobShard struct {
	ent.Schema
}

func (JobShard) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "job_shards"},
	}
}

func (JobShard) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (JobShard) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("instance_id"),
		field.String("shard_no").MaxLen(64).Unique(),
		field.Int("shard_index"),
		field.Int("shard_total"),
		field.String("route_key").MaxLen(128).Optional().Nillable(),
		field.Int64("executor_id").Optional().Nillable(),
		field.String("status").MaxLen(32),
		field.String("payload").Optional().Nillable(),
		field.Int("retry_count").Default(0),
		field.Time("start_time").Optional().Nillable(),
		field.Time("finish_time").Optional().Nillable(),
		field.String("error_message").MaxLen(1024).Optional().Nillable(),
	}
}

func (JobShard) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("instance_id", "status"),
		index.Fields("instance_id", "shard_index").Unique(),
	}
}
