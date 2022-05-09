package utils

import (
	"context"
	"sync"
)

// Refurbisher handles serial tasks obsoleting each other.
type Refurbisher struct {
	mtx    sync.Mutex
	cancel func()
}

// NewRefurbisher returns a new Refurbisher.
func NewRefurbisher() *Refurbisher {
	return &Refurbisher{cancel: func() {}}
}

// Do runs f after cancelling the previous f (if still running).
func (r *Refurbisher) Do(ctx context.Context, f func(context.Context) error) error {
	var childCtx context.Context

	r.mtx.Lock()
	cancel := r.cancel
	childCtx, r.cancel = context.WithCancel(ctx)
	r.mtx.Unlock()

	cancel()

	return f(childCtx)
}
