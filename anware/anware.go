package anware

import (
	"github.com/Aninetix/core/aninterface"

	"context"
	"sync"
)

type AnWareReply struct {
	Data any
	Err  error
}

type AnWareEvent struct {
	Source string
	Target string
	Type   string
	Data   any

	ReplyTo chan AnWareReply
}

type BaseModule struct {
	name string

	ctx context.Context
	in  <-chan AnWareEvent
	mw  *AnWare
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

// --- AnWare ---

type AnWare struct {
	routes map[string]chan AnWareEvent
	mods   map[string]AnModule
	bus    chan AnWareEvent
	wg     sync.WaitGroup

	context context.Context
	cancel  context.CancelFunc

	Logger aninterface.AnLogger
}

func NewAnWare(context context.Context, cancel context.CancelFunc, logger aninterface.AnLogger) *AnWare {
	return &AnWare{
		routes:  make(map[string]chan AnWareEvent),
		mods:    make(map[string]AnModule),
		bus:     make(chan AnWareEvent, 256),
		context: context,
		cancel:  cancel,
		Logger:  logger,
	}
}
