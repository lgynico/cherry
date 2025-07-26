package cherryActor

import (
	"sync"

	cfacade "github.com/cherry-game/cherry/facade"
)

type (
	eventMap struct {
		m sync.Map
	}
)

func (p *eventMap) Register(name string, actor *Actor) {
	value, _ := p.m.LoadOrStore(name, &sync.Map{})

	actorMap := value.(*sync.Map)
	actorMap.LoadOrStore(actor, struct{}{})
}

func (p *eventMap) Unregister(actor *Actor) {
	p.m.Range(func(key, value any) bool {
		actorMap := value.(*sync.Map)
		actorMap.Delete(actor)
		return true
	})
}

func (p *eventMap) PostEvent(e cfacade.IEventData) {
	if value, ok := p.m.Load(e.Name()); ok {
		actorMap := value.(*sync.Map)
		actorMap.Range(func(key, _ any) bool {
			actor := key.(*Actor)
			if actor.State() == WorkerState {
				actor.event.Push(e)
			}

			return true
		})
	}
}
