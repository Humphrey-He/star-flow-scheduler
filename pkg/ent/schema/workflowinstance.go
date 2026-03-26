package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/schema/mixin"
)

type WorkflowInstance struct {
	ent.Schema
}

func (WorkflowInstance) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "workflow_instances"},
	}
}

func (WorkflowInstance) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (WorkflowInstance) Fields() []ent.Field {
	return []ent.Field{
		field.String("workflow_instance_no").MaxLen(64).Unique(),
		field.Int64("workflow_id"),
		field.String("workflow_code").MaxLen(64),
		field.String("status").MaxLen(32),
		field.Time("start_time").Optional().Nillable(),
		field.Time("finish_time").Optional().Nillable(),
		field.String("error_message").MaxLen(1024).Optional().Nillable(),
	}
}

func (WorkflowInstance) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("workflow_code", "status"),
		index.Fields("workflow_id"),
	}
}
