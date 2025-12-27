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
		panic("module already registered: " + desc.Name)
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

// func RegisterModule[F any, C any](
// 	name string,
// 	constructor func(local aninterface.StaticData, config C, flags F, logger aninterface.AnLogger) AnModule,
// ) {
// 	moduleRegistry[name] = GenericModuleConstructor[F, C]{fn: constructor}
// }

// func (m *AnWare) AutoLoadModules(staticData aninterface.StaticData, configData any, flags any, logger aninterface.AnLogger) {
// 	for name, constructor := range moduleRegistry {

// 		mod := constructor.Build(staticData, configData, flags, logger)

// 		m.routes[name] = make(chan AnWareEvent, 128)
// 		m.mods[name] = mod

// 		m.Logger.Info("[ANWARE] Auto-loaded module: " + name)
// 	}
// }

func (m *AnWare) AutoLoadModules(
	staticData aninterface.StaticData,
	appConfig any,
	logger aninterface.AnLogger,
) {
	for name, desc := range moduleRegistry {

		cfg := extractSubConfig(appConfig, name, desc.ConfigType)
		cfgVal := reflect.ValueOf(cfg)

		if cfgVal.Kind() == reflect.Ptr {
			cfgVal = cfgVal.Elem()
		}

		if cfgVal.IsZero() {
			fmt.Print("module disabled, config value not Set: " + name)
			// logger.Info("module disabled (empty config): " + name)
			continue
		}

		if v, ok := cfg.(ConfigValidator); ok {
			if err := v.Validate(); err != nil {
				logger.Error(
					fmt.Sprintf("[ANWARE] module %s disabled: invalid config: %v", name, err),
				)
				continue
			}
		}
		mod := desc.New(staticData, cfg, logger)

		m.routes[name] = make(chan AnWareEvent, 128)
		m.mods[name] = mod

		logger.Info("[ANWARE] Auto-loaded module: " + name)
	}
}
