package types

type DeadLetterStatus string

const (
	DeadLetterStatusOpen    DeadLetterStatus = "open"
	DeadLetterStatusRetried DeadLetterStatus = "retried"
	DeadLetterStatusClosed  DeadLetterStatus = "closed"
)
