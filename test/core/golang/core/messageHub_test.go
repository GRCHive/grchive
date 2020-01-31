package core_test

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"gitlab.com/grchive/grchive/core"
	"sync"
	"testing"
)

func createMessageHub() (core.RealMessageHub, core.ListenerMap) {
	listeners := make(core.ListenerMap)
	return core.CreateRealMessageHub(listeners), listeners
}

func TestCreateRealMessageHub(t *testing.T) {
	hub, listeners := createMessageHub()
	assert.Equal(t, 0, hub.NumTotalListeners())

	// Initialize map for testing
	listeners[0] = make(core.SubtypeListenerMap)
	l1 := &list.List{}
	listeners[0]["Test"] = l1

	l2 := &list.List{}
	listeners[0]["Teeth"] = l2

	listeners[5] = make(core.SubtypeListenerMap)
	l3 := &list.List{}
	listeners[5]["Blah"] = l3

	l1.PushBack(10)
	assert.Equal(t, 1, hub.NumTotalListeners())
	assert.Equal(t, 1, hub.NumMessageTypeListeners(0))
	assert.Equal(t, 1, hub.NumSubTypeListeners(0, "Test"))
	assert.Equal(t, 0, hub.NumSubTypeListeners(0, "Teeth"))
	assert.Equal(t, 0, hub.NumMessageTypeListeners(5))
	assert.Equal(t, 0, hub.NumSubTypeListeners(5, "Blah"))

	l2.PushBack(20)
	assert.Equal(t, 2, hub.NumTotalListeners())
	assert.Equal(t, 2, hub.NumMessageTypeListeners(0))
	assert.Equal(t, 1, hub.NumSubTypeListeners(0, "Test"))
	assert.Equal(t, 1, hub.NumSubTypeListeners(0, "Teeth"))
	assert.Equal(t, 0, hub.NumMessageTypeListeners(5))
	assert.Equal(t, 0, hub.NumSubTypeListeners(5, "Blah"))

	l3.PushBack(30)
	assert.Equal(t, 3, hub.NumTotalListeners())
	assert.Equal(t, 2, hub.NumMessageTypeListeners(0))
	assert.Equal(t, 1, hub.NumSubTypeListeners(0, "Test"))
	assert.Equal(t, 1, hub.NumSubTypeListeners(0, "Teeth"))
	assert.Equal(t, 1, hub.NumMessageTypeListeners(5))
	assert.Equal(t, 1, hub.NumSubTypeListeners(5, "Blah"))
}

func TestRegisterListener(t *testing.T) {
	hub, _ := createMessageHub()
	assert.Equal(t, 0, hub.NumTotalListeners())

	testChannel := make(chan core.MessagePayload)
	ele := hub.RegisterListener(0, "Test", testChannel)
	assert.Equal(t, 1, hub.NumTotalListeners())
	assert.Equal(t, 1, hub.NumMessageTypeListeners(0))
	assert.Equal(t, 1, hub.NumSubTypeListeners(0, "Test"))

	hub.UnregisterListener(0, "Test", ele)
	assert.Equal(t, 0, hub.NumTotalListeners())
	assert.Equal(t, 0, hub.NumMessageTypeListeners(0))
	assert.Equal(t, 0, hub.NumSubTypeListeners(0, "Test"))
}

func TestMultiThreadRegisterListener(t *testing.T) {
	hub, _ := createMessageHub()
	assert.Equal(t, 0, hub.NumTotalListeners())
	testChannel := make(chan core.MessagePayload)
	eles := make([]*list.Element, 26)

	wg := sync.WaitGroup{}

	register := func(total int, typ core.MessageType, subtype core.MessageSubtype, offset int) {
		defer wg.Done()

		for i := 0; i < total; i++ {
			eles[offset+i] = hub.RegisterListener(typ, subtype, testChannel)
		}
	}

	wg.Add(4)

	go register(5, 0, "Test", 0)
	go register(10, 0, "Test", 5)
	go register(3, 0, "Test2", 15)
	go register(8, 1, "Blah", 18)

	wg.Wait()

	assert.Equal(t, 26, hub.NumTotalListeners())
	assert.Equal(t, 18, hub.NumMessageTypeListeners(0))
	assert.Equal(t, 15, hub.NumSubTypeListeners(0, "Test"))
	assert.Equal(t, 3, hub.NumSubTypeListeners(0, "Test2"))
	assert.Equal(t, 8, hub.NumMessageTypeListeners(1))
	assert.Equal(t, 8, hub.NumSubTypeListeners(1, "Blah"))

	unregister := func(start int, num int, typ core.MessageType, subtype core.MessageSubtype) {
		defer wg.Done()

		for i := 0; i < num; i++ {
			hub.UnregisterListener(typ, subtype, eles[start+i])
		}
	}

	wg.Add(4)

	go unregister(0, 5, 0, "Test")
	go unregister(5, 10, 0, "Test")
	go unregister(15, 3, 0, "Test2")
	go unregister(18, 8, 1, "Blah")

	wg.Wait()

	assert.Equal(t, 0, hub.NumTotalListeners())
	assert.Equal(t, 0, hub.NumMessageTypeListeners(0))
	assert.Equal(t, 0, hub.NumSubTypeListeners(0, "Test"))
	assert.Equal(t, 0, hub.NumSubTypeListeners(0, "Test2"))
	assert.Equal(t, 0, hub.NumMessageTypeListeners(1))
	assert.Equal(t, 0, hub.NumSubTypeListeners(1, "Blah"))
}

func TestSendReceiveMessage(t *testing.T) {
	hub, _ := createMessageHub()
	c1 := make(chan core.MessagePayload)
	c2 := make(chan core.MessagePayload)

	hub.RegisterListener(0, "c1", c1)
	hub.RegisterListener(0, "c2", c2)

	m1 := "Test1"
	m2 := "Test2"
	m3 := "Test3"

	go func() {
		hub.SendMessage(0, "c1", m1)
	}()

	go func() {
		hub.SendMessage(0, "c2", m2)
	}()

	go func() {
		hub.SendMessage(0, "c3", m3)
	}()

	t1, t2 := <-c1, <-c2
	assert.Equal(t, m1, t1)
	assert.Equal(t, m2, t2)
}
