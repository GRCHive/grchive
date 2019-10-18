package core

import (
	"container/list"
	"sync"
)

// Handles passing messages to relevant listeners.
// The users of this file must be able to treat the functionality
// in this file as being able to pass/receive messages to all
// running instances of the webserver.

type MessageType uint
type MessageSubtype string
type MessagePayload interface{}

const (
	UpdateDisplaySettingsForProcessFlowNode MessageType = iota
)

type ListenerMap map[MessageType]map[MessageSubtype]*list.List

var registeredListeners = make(ListenerMap)
var registerMutex sync.RWMutex

func SendMessage(typ MessageType, subtyp MessageSubtype, payload MessagePayload) {
	// TODO: Make this work for multiple instances (e.g. Cloud PubSub?)
	ReceiveMessage(typ, subtyp, payload)
}

func ReceiveMessage(typ MessageType, subtyp MessageSubtype, payload MessagePayload) {
	registerMutex.RLock()
	defer registerMutex.RUnlock()

	subtypeListeners, ok := registeredListeners[typ]
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

func RegisterListener(typ MessageType, subtyp MessageSubtype, c chan MessagePayload) *list.Element {
	registerMutex.Lock()
	defer registerMutex.Unlock()

	subtypeListeners, ok := registeredListeners[typ]
	if !ok {
		registeredListeners[typ] = make(map[MessageSubtype]*list.List)
		subtypeListeners = registeredListeners[typ]
	}

	listeners, ok := subtypeListeners[subtyp]
	if !ok {
		subtypeListeners[subtyp] = list.New()
		listeners = subtypeListeners[subtyp]
	}

	return listeners.PushBack(c)
}

func UnregisterListener(typ MessageType, subtyp MessageSubtype, e *list.Element) {
	registerMutex.Lock()
	defer registerMutex.Unlock()

	subtypeListeners, ok := registeredListeners[typ]
	if !ok {
		return
	}

	listeners, ok := subtypeListeners[subtyp]
	if !ok {
		return
	}

	listeners.Remove(e)
}
