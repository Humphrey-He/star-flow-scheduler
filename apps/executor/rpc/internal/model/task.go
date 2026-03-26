package model

type Task struct {
	InstanceNo    string
	ShardNo       string
	JobCode       string
	HandlerName   string
	Payload       []byte
	TimeoutMs     int32
	TraceID       string
	IdempotentKey string
	ShardIndex    int32
	ShardTotal    int32
}
