package anware

import (
	"github.com/Aninetix/core/aninterface"

	"context"
	"sync"
)

// --- Types ---

type AnWareEvent struct {
	Source string
	Target string
	Type   string
	Data   any
}

type AnModule interface {
	Name() string
	Param(ctx context.Context, in <-chan AnWareEvent, mw *AnWare)
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

// --- Constructeur ---

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
