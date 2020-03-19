package core

import (
	"container/list"
	"sync"
)

// Handles passing messages to relevant listeners.
// This functionality is only to pass messages between different parts of
// a single instance of the webserver. If it needs to be passed to others
// then a message should be sent out via RabbitMQ.

type MessageType uint
type MessageSubtype string
type MessagePayload interface{}

const (
	UpdateDisplaySettingsForProcessFlowNode MessageType = iota
	MHUserNotification
)

type SubtypeListenerMap map[MessageSubtype]*list.List
type ListenerMap map[MessageType]SubtypeListenerMap

type MessageHub interface {
	SendMessage(typ MessageType, subtyp MessageSubtype, payload MessagePayload)
	ReceiveMessage(typ MessageType, subtyp MessageSubtype, payload MessagePayload)
	RegisterListener(typ MessageType, subtyp MessageSubtype, c chan MessagePayload) *list.Element
	UnregisterListener(typ MessageType, subtyp MessageSubtype, e *list.Element)
}

type RealMessageHub struct {
	registeredListeners ListenerMap
	registerMutex       sync.RWMutex
}

func (m RealMessageHub) NumTotalListeners() int {
	total := 0
	for k, _ := range m.registeredListeners {
		total = total + m.NumMessageTypeListeners(k)
	}
	return total
}

func (m RealMessageHub) NumMessageTypeListeners(typ MessageType) int {
	total := 0
	for k, _ := range m.registeredListeners[typ] {
		total = total + m.NumSubTypeListeners(typ, k)
	}
	return total
}

func (m RealMessageHub) NumSubTypeListeners(typ MessageType, subtyp MessageSubtype) int {
	return m.registeredListeners[typ][subtyp].Len()
}

func (m *RealMessageHub) SendMessage(typ MessageType, subtyp MessageSubtype, payload MessagePayload) {
	m.registerMutex.RLock()
	defer m.registerMutex.RUnlock()

	// TODO: Make this work for multiple instances (e.g. Cloud PubSub?)
	m.ReceiveMessage(typ, subtyp, payload)
}

func (m *RealMessageHub) ReceiveMessage(typ MessageType, subtyp MessageSubtype, payload MessagePayload) {
	m.registerMutex.RLock()
	defer m.registerMutex.RUnlock()

	subtypeListeners, ok := m.registeredListeners[typ]
	if !ok {
		return
	}

	listeners, ok := subtypeListeners[subtyp]
	if !ok {
		return
	}

	for e := listeners.Front(); e != nil; e = e.Next() {
		l := e.Value
		l.(chan MessagePayload) <- payload
	}
}

func (m *RealMessageHub) RegisterListener(typ MessageType, subtyp MessageSubtype, c chan MessagePayload) *list.Element {
	m.registerMutex.Lock()
	defer m.registerMutex.Unlock()

	subtypeListeners, ok := m.registeredListeners[typ]
	if !ok {
		m.registeredListeners[typ] = make(map[MessageSubtype]*list.List)
		subtypeListeners = m.registeredListeners[typ]
	}

	listeners, ok := subtypeListeners[subtyp]
	if !ok {
		subtypeListeners[subtyp] = list.New()
		listeners = subtypeListeners[subtyp]
	}

	return listeners.PushBack(c)
}

func (m *RealMessageHub) UnregisterListener(typ MessageType, subtyp MessageSubtype, e *list.Element) {
	m.registerMutex.Lock()
	defer m.registerMutex.Unlock()

	subtypeListeners, ok := m.registeredListeners[typ]
	if !ok {
		return
	}

	listeners, ok := subtypeListeners[subtyp]
	if !ok {
		return
	}

	listeners.Remove(e)
}

func CreateRealMessageHub(m ListenerMap) RealMessageHub {
	return RealMessageHub{
		registeredListeners: m,
	}
}

var DefaultMessageHub = CreateRealMessageHub(make(ListenerMap))
