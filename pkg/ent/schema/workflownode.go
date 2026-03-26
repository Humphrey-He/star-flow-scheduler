package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent/schema/mixin"
)

type WorkflowNode struct {
	ent.Schema
}

func (WorkflowNode) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "workflow_nodes"},
	}
}

func (WorkflowNode) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

func (WorkflowNode) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("workflow_id"),
		field.String("node_code").MaxLen(64),
		field.String("node_name").MaxLen(128),
		field.String("job_code").MaxLen(64),
		field.String("upstream_codes").Optional().Nillable(),
		field.String("trigger_condition").MaxLen(32).Default("all_success"),
		field.String("fail_strategy").MaxLen(32).Default("stop"),
		field.Int("timeout_ms").Default(60000),
		field.Int("sort_order").Default(0),
	}
}

func (WorkflowNode) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("workflow_id", "node_code").Unique(),
		index.Fields("workflow_id"),
	}
}
