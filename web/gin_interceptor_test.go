package core

import (
	"context"
	"testing"
)

func TestGinInterceptor(t *testing.T) {
	_ = GinInterceptor(func(ctx context.Context) error { return nil })
}
