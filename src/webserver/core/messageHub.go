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
type MessagePayload interface{}

const (
	UpdateDisplaySettingsForProcessFlowNode MessageType = iota
)

type ListenerMap map[MessageType]*list.List

var registeredListeners = make(ListenerMap)
var registerMutex sync.RWMutex

func SendMessage(typ MessageType, payload MessagePayload) {
	// TODO: Make this work for multiple instances (e.g. Cloud PubSub?)
	ReceiveMessage(typ, payload)
}

func ReceiveMessage(typ MessageType, payload MessagePayload) {
	registerMutex.RLock()
	defer registerMutex.RUnlock()

	listeners, ok := registeredListeners[typ]
	if !ok {
		return
	}

	for e := listeners.Front(); e != nil; e = e.Next() {
		l := e.Value
		l.(chan MessagePayload) <- payload
	}
}

func RegisterListener(typ MessageType, c chan MessagePayload) *list.Element {
	registerMutex.Lock()
	defer registerMutex.Unlock()

	listeners, ok := registeredListeners[typ]
	if !ok {
		registeredListeners[typ] = list.New()
		listeners = registeredListeners[typ]
	}

	return listeners.PushBack(c)
}

func UnregisterListener(typ MessageType, e *list.Element) {
	registerMutex.Lock()
	defer registerMutex.Unlock()

	listeners, ok := registeredListeners[typ]
	if !ok {
		return
	}

	listeners.Remove(e)
}
