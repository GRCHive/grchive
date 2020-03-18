package webcore

import (
	"encoding/json"
	"gitlab.com/grchive/grchive/core"
	"reflect"
)

func SendEventToRabbitMQ(event core.Event) {
	msg := EventMessage{
		Event: event,
	}

	if event.Object != nil {
		msg.ObjectType = core.GetBaseType(event.Object).String()
	}

	if event.IndirectObject != nil {
		msg.IndirectObjectType = core.GetBaseType(event.IndirectObject).String()
	}

	DefaultRabbitMQ.SendMessage(PublishMessage{
		Exchange: EVENT_EXCHANGE,
		Queue:    "",
		Body:     msg,
	})
}

func (m EventMessage) RecreateEvent() (*core.Event, error) {
	event := core.Event{
		Subject:   m.Event.Subject,
		Verb:      m.Event.Verb,
		Timestamp: m.Event.Timestamp,
	}

	if m.Event.Object != nil {
		rObject := reflect.New(core.TypeRegistry[m.ObjectType])
		object := rObject.Interface()

		raw, err := json.Marshal(m.Event.Object)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(raw, object)
		if err != nil {
			return nil, err
		}

		event.Object = rObject.Elem().Interface()
	}

	if m.Event.IndirectObject != nil {
		rObject := reflect.New(core.TypeRegistry[m.IndirectObjectType])
		object := rObject.Interface()

		raw, err := json.Marshal(m.Event.Object)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(raw, &object)
		if err != nil {
			return nil, err
		}

		event.IndirectObject = rObject.Elem().Interface()
	}

	return &event, nil
}
