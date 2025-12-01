package anware

import "github.com/Aninetix/core/aninterface"

var moduleRegistry = map[string]ModuleFactory{}

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

func RegisterModule[F any, C any](
	name string,
	constructor func(local aninterface.StaticData, config C, flags F, logger aninterface.AnLogger) AnModule,
) {
	moduleRegistry[name] = GenericModuleConstructor[F, C]{fn: constructor}
}

func (m *AnWare) AutoLoadModules(staticData aninterface.StaticData, configData any, flags any, logger aninterface.AnLogger) {
	for name, constructor := range moduleRegistry {

		mod := constructor.Build(staticData, configData, flags, logger)

		m.routes[name] = make(chan AnWareEvent, 128)
		m.mods[name] = mod

		m.Logger.Info("[ANWARE] Auto-loaded module: " + name)
	}
}
