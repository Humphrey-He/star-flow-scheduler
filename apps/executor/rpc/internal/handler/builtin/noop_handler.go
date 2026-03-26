package builtin

import "context"

type NoopHandler struct{}

func (h *NoopHandler) Name() string {
	return "noop"
}

func (h *NoopHandler) Execute(ctx context.Context, payload []byte) error {
	return nil
}
