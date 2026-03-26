package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/schema/mixin"
)

type WorkflowDefinition struct {
	ent.Schema
}

func (WorkflowDefinition) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "workflow_definitions"},
	}
}

func (WorkflowDefinition) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (WorkflowDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.String("workflow_code").MaxLen(64).Unique(),
		field.String("workflow_name").MaxLen(128),
		field.String("description").MaxLen(512).Optional().Nillable(),
		field.String("status").MaxLen(32).Default("enabled"),
	}
}

func (WorkflowDefinition) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
	}
}
