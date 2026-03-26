package handler

import (
	"context"
	"testing"
)

type testHandler struct {
	name string
}

func (h *testHandler) Name() string {
	return h.name
}

func (h *testHandler) Execute(ctx context.Context, payload []byte) error {
	return nil
}

func TestRegistryRegisterAndGet(t *testing.T) {
	reg := NewRegistry()
	h := &testHandler{name: "t1"}

	if err := reg.Register(h); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := reg.Register(h); err == nil {
		t.Fatalf("expected duplicate error")
	}

	got, ok := reg.Get("t1")
	if !ok || got == nil {
		t.Fatalf("expected handler to be found")
	}
}
