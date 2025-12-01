package anware

import "fmt"

// --- Lancement ---

func (m *AnWare) Run() {
	go m.dispatchLoop()

	for name, mod := range m.mods {
		m.wg.Add(1)

		ch := m.routes[name]

		mod.Param(m.context, ch, m)

		go func(mod AnModule, in <-chan AnWareEvent) {
			defer m.wg.Done()
			mod.Start()
		}(mod, ch)

		m.Logger.Info("[ANWARE] Module loaded: " + name)

	}
}

// --- Arrêt propre ---

func (m *AnWare) Shutdown() {
	m.Logger.Info("[ANWARE] Stopping AnWare...")

	for _, mod := range m.mods {
		if err := mod.Stop(); err != nil {
			m.Logger.Error(fmt.Sprintf("[ANWARE] Error stopping module %s: %v", mod.Name(), err))
		}
	}

	if m.cancel != nil {
		m.cancel()
	}

	close(m.bus)
	m.wg.Wait()
	m.Logger.Info("[ANWARE] All modules stopped.")
}

// --- Diffusion à TOUS les modules (pub/sub) ---

func (m *AnWare) Broadcast(msg AnWareEvent) {
	msgCopy := msg
	for name := range m.mods {
		if name == msg.Source {
			continue // évite d’envoyer à soi-même
		}

		msgCopy.Target = name
		m.Send(msgCopy)
	}
}

// --- Envoi d'événements ---

func (m *AnWare) Send(msg AnWareEvent) {
	select {
	case m.bus <- msg:
	default:
		m.Logger.Info(fmt.Sprintf("[ANWARE] Bus full, event dropped: %+v", msg))
	}
}

// --- Boucle de dispatch ---

func (m *AnWare) dispatchLoop() {
	for {
		select {
		case <-m.context.Done():
			return
		case msg, ok := <-m.bus:
			if !ok {
				return
			}
			m.LoopOfAnWare(msg)
			m.routeMessage(msg)
		}
	}
}

func (m *AnWare) LoopOfAnWare(msg AnWareEvent) {
	if msg.Target == "anWare" {
		if msg.Type == "exit" {
			m.Shutdown()
		}
	}
}

func (m *AnWare) routeMessage(msg AnWareEvent) {

	if msg.Target == "*" {
		m.Broadcast(msg)
		return
	}

	targetCh, found := m.routes[msg.Target]
	if !found {
		m.Logger.Info(fmt.Sprintf("[ANWARE] No module found for target: %s", msg.Target))
		return
	}

	select {
	case targetCh <- msg:
	default:
		m.Logger.Info(fmt.Sprintf("[ANWARE] Channel full for %s, event ignored: %+v", msg.Target, msg))
	}
}
