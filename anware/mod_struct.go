package anware

import "context"

type BaseModule struct {
	name string

	ctx context.Context
	in  <-chan AnWareEvent
	mw  *AnWare
}

type AnModule interface {
	Name() string

	SetRuntime(
		name string,
		ctx context.Context,
		in <-chan AnWareEvent,
		mw *AnWare,
	)

	Start()
	Stop() error
}

func (b *BaseModule) SetRuntime(
	name string,
	ctx context.Context,
	in <-chan AnWareEvent,
	mw *AnWare,
) {
	b.name = name
	b.ctx = ctx
	b.in = in
	b.mw = mw
}

func (b *BaseModule) Name() string {
	return b.name
}

func (b *BaseModule) Ctx() context.Context {
	return b.ctx
}

func (b *BaseModule) In() <-chan AnWareEvent {
	return b.in
}

func (b *BaseModule) MW() *AnWare {
	return b.mw
}
