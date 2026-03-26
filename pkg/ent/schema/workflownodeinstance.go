package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/schema/mixin"
)

type WorkflowNodeInstance struct {
	ent.Schema
}

func (WorkflowNodeInstance) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "workflow_node_instances"},
	}
}

func (WorkflowNodeInstance) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (WorkflowNodeInstance) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("workflow_instance_id"),
		field.Int64("workflow_id"),
		field.String("node_code").MaxLen(64),
		field.Int64("job_id"),
		field.Int64("job_instance_id").Optional().Nillable(),
		field.String("status").MaxLen(32).Default("pending"),
		field.String("error_message").MaxLen(512).Optional().Nillable(),
		field.Time("start_time").Optional().Nillable(),
		field.Time("finish_time").Optional().Nillable(),
	}
}

func (WorkflowNodeInstance) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("workflow_instance_id", "node_code").Unique(),
		index.Fields("workflow_instance_id"),
		index.Fields("workflow_id"),
		index.Fields("status"),
	}
}
