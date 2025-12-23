package anware

import (
	"fmt"
	"time"
)

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

func (m *AnWare) Broadcast(msg AnWareEvent) {
	msgCopy := msg
	for name := range m.mods {
		if name == msg.Source {
			continue
		}

		msgCopy.Target = name
		m.Send(msgCopy)
	}
}

func (m *AnWare) Send(msg AnWareEvent) {
	select {
	case m.bus <- msg:
	default:
		m.Logger.Info(fmt.Sprintf("[ANWARE] Bus full, event dropped: %+v", msg))
	}
}

func (m *AnWare) SendSync(
	source string,
	target string,
	msgType string,
	data any,
) (any, error) {

	replyCh := make(chan AnWareReply, 1)

	m.Send(AnWareEvent{
		Source:  source,
		Target:  target,
		Type:    msgType,
		Data:    data,
		ReplyTo: replyCh,
	})

	select {
	case reply := <-replyCh:
		return reply.Data, reply.Err

	case <-m.context.Done():
		return nil, fmt.Errorf("anware shutting down")

	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout waiting reply from %s", target)
	}
}

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

		if msg.ReplyTo != nil {
			msg.ReplyTo <- AnWareReply{
				Err: fmt.Errorf("target module not found: %s", msg.Target),
			}
		}
		return
	}

	select {
	case targetCh <- msg:
	default:
		m.Logger.Info(fmt.Sprintf("[ANWARE] Channel full for %s, event ignored: %+v", msg.Target, msg))

		if msg.ReplyTo != nil {
			msg.ReplyTo <- AnWareReply{
				Err: fmt.Errorf("target module busy: %s", msg.Target),
			}
		}
	}
}
