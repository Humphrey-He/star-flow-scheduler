package builtin

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type DemoHandler struct{}

func (h *DemoHandler) Name() string {
	return "demo_print"
}

func (h *DemoHandler) Execute(ctx context.Context, payload []byte) error {
	logx.WithContext(ctx).Infof("demo handler payload=%s", string(payload))
	return nil
}
