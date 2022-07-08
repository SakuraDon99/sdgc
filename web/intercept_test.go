package web

import (
	"context"
	"testing"
)

func TestGinInterceptor(t *testing.T) {
	_ = Intercept(func(ctx context.Context) error { return nil })
}
