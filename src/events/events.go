package events

import "sync"

type Event struct {
	Name string
	Data string
}

type EventListener func(event Event)

type EventManager struct {
	listeners map[string][]EventListener
	mu        sync.Mutex
}

func NewEventManager() *EventManager {
	return &EventManager{
		listeners: make(map[string][]EventListener),
	}
}

func (em *EventManager) AddListener(eventName string, listener EventListener) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.listeners[eventName] = append(em.listeners[eventName], listener)
}

func (em *EventManager) Emit(event Event) {
	em.mu.Lock()
	defer em.mu.Unlock()
	if listeners, ok := em.listeners[event.Name]; ok {
		for _, listener := range listeners {
			go listener(event)
		}
	}
}
