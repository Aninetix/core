package anware

import (
	"fmt"
	"reflect"

	"github.com/Aninetix/core/aninterface"
)

type ModuleDescriptor struct {
	Name string

	New func(
		local aninterface.StaticData,
		cfg any,
		logger aninterface.AnLogger,
	) AnModule

	ConfigType any
}

type ConfigValidator interface {
	Validate() error
}

var moduleRegistry = map[string]ModuleDescriptor{}

func RegisterModule(desc ModuleDescriptor) {
	if _, ok := moduleRegistry[desc.Name]; ok {
		panic("[ANWARE] FATAL: module already registered: " + desc.Name)
	}
	moduleRegistry[desc.Name] = desc
}

// var moduleRegistry = map[string]ModuleFactory{}

type ModuleFactory interface {
	Build(local aninterface.StaticData, config any, flags any, logger aninterface.AnLogger) AnModule
}

type GenericModuleConstructor[F any, C any] struct {
	fn func(local aninterface.StaticData, config C, flags F, logger aninterface.AnLogger) AnModule
}

func (g GenericModuleConstructor[F, C]) Build(local aninterface.StaticData, config any, flags any, logger aninterface.AnLogger) AnModule {
	return g.fn(
		local,
		config.(C),
		flags.(F),
		logger,
	)
}

func (m *AnWare) AutoLoadModules(
	staticData aninterface.StaticData,
	appConfig any,
	logger aninterface.AnLogger,
) {
	for name, desc := range moduleRegistry {

		cfg, err := extractSubConfig(appConfig, name, desc.ConfigType)
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		cfgVal := reflect.ValueOf(cfg)
		if cfgVal.Kind() == reflect.Ptr {
			cfgVal = cfgVal.Elem()
		}

		if cfgVal.IsZero() {
			logger.Info("[ANWARE] module disabled (empty config): " + name)
			continue
		}

		if v, ok := cfg.(ConfigValidator); ok {
			if err := v.Validate(); err != nil {
				logger.Error(
					fmt.Sprintf(
						"[ANWARE] module %s disabled: invalid config: %v",
						name,
						err,
					),
				)
				continue
			}
		}

		mod := desc.New(staticData, cfg, logger)

		m.routes[name] = make(chan AnWareEvent, 128)
		m.mods[name] = mod

		logger.Info("[ANWARE] auto-loaded module: " + name)
	}
}
