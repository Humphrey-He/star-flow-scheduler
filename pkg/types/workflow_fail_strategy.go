package types

type WorkflowFailStrategy string

const (
	WorkflowFailStrategyStop     WorkflowFailStrategy = "stop"
	WorkflowFailStrategyContinue WorkflowFailStrategy = "continue"
)
