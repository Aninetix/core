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

func InitCore[F any, C any]() (*F, *C, aninterface.AnLogger) {

	// --- FLAGS ---
	var flg F
	if err := anflags.ParseFlags(&flg); err != nil {
		panic(err)
	}

	// --- LOGGER ---
	logPath := helpers.GetFieldString(&flg, "LogPath")
	if logPath == "" {
		fmt.Println("Error Developer: LogPath is empty on flags Struct")
		os.Exit(1)
	}

	debugOn := helpers.GetFieldBool(&flg, "Debug")

	logger := anlogger.NewLogger(logPath, debugOn)

	// --- CONFIG ---
	var cfg C
	configPath := helpers.GetFieldString(&flg, "ConfigPath")
	if configPath == "" {
		fmt.Println("Error Developer: ConfigPath is empty on flags Struct")
		os.Exit(1)
	}

	if err := anconfig.LoadConfig(configPath, &cfg); err != nil {
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
	core.AnWare.AutoLoadModules(core.Data, core.Config, core.Flags, core.Logger)
	core.AnWare.Run()
	core.Logger.Info("[ANCORE] AnCore is running.")
}
