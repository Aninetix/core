package ancore

import (
	"context"
	"fmt"
	"os"

	"github.com/Aninetix/core/aninterface"
	"github.com/Aninetix/core/anware"
	"github.com/Aninetix/core/internal/anconfig"
	"github.com/Aninetix/core/internal/anflags"
	"github.com/Aninetix/core/internal/anlocal"
	"github.com/Aninetix/core/internal/anlogger"
	"github.com/Aninetix/core/internal/helpers"
)

type AnCore struct {
	Flags  any
	Config any
	Logger aninterface.AnLogger
	AnWare *anware.AnWare
	Data   aninterface.StaticData
}

type InitOptions struct {
	LogPath    string
	ConfigPath string
	Debug      *bool
}

type Option func(*InitOptions)

func WithLogPath(p string) Option {
	return func(o *InitOptions) { o.LogPath = p }
}

func WithConfigPath(p string) Option {
	return func(o *InitOptions) { o.ConfigPath = p }
}

func WithDebug(b bool) Option {
	return func(o *InitOptions) { o.Debug = &b }
}

func ptrBool(b bool) *bool {
	return &b
}

func InitCore[F any, C any](opts ...Option) (*F, *C, aninterface.AnLogger) {
	// default options, e.g. from flags
	var flg F
	var cfg C

	if err := anflags.ParseFlags(&flg); err != nil {
		panic(err)
	}

	o := InitOptions{
		LogPath:    helpers.GetFieldString(&flg, "LogPath"),
		ConfigPath: helpers.GetFieldString(&flg, "ConfigPath"),
		Debug:      ptrBool(helpers.GetFieldBool(&flg, "Debug")),
	}

	// override with provided optional params
	for _, opt := range opts {
		opt(&o)
	}

	// --- LOGGER ---
	logger := anlogger.NewLogger(o.LogPath, *o.Debug)

	// --- CONFIG ---
	if err := anconfig.LoadConfig(o.ConfigPath, &cfg); err != nil {
		logger.Error(fmt.Sprintf("Erreur chargement config: %v", err))
		os.Exit(1)
	}

	return &flg, &cfg, logger
}

func BootCore(flg any, cfg any, logger aninterface.AnLogger, ctx context.Context, cancel context.CancelFunc) AnCore {
	anStaticData := anlocal.LoadStaticData()

	return AnCore{
		Data:   anStaticData,
		Logger: logger,
		AnWare: anware.NewAnWare(ctx, cancel, logger),
		Flags:  flg,
		Config: cfg,
	}
}

func (core *AnCore) Run() {
	core.Logger.Info("[ANCORE] Booting AnCore...")
	core.AnWare.AutoLoadModules(core.Data, core.Config, core.Logger)
	core.AnWare.Run()
	core.Logger.Info("[ANCORE] AnCore is running.")
}
