package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"sync"
)

type ProcessFlowNodeDisplaySettingsPayload struct {
	NodeId   int64
	Settings map[string]interface{}
}

func processProcessFlowNodeDisplaySettings(conn *websocket.Conn) {
	// Channel to receive communications from the message hub
	// about relevant events to send to the user.
	var hubChannel chan core.MessagePayload = make(chan core.MessagePayload)
	ele := core.RegisterListener(core.UpdateDisplaySettingsForProcessFlowNode, hubChannel)
	defer core.UnregisterListener(core.UpdateDisplaySettingsForProcessFlowNode, ele)

	waitGroup := sync.WaitGroup{}

	// Spawn one goroutine to listen for user messages.
	// User messages should trigger a database query to update
	// the node display - we assume the user is authoritative.
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for {
			data := ProcessFlowNodeDisplaySettingsPayload{}
			err := conn.ReadJSON(&data)
			if err != nil {
				core.Warning("Failed to read websocket message: " + err.Error())
				break
			}

			// TODO: What's the best way to handle errors here?
			err = database.UpdateDisplaySettingsForProcessFlowNode(data.NodeId, data.Settings)
			if err != nil {
				core.Warning("Failed to update display settings: " + err.Error())
				break
			}
		}
	}()

	// Spawn another goroutine to listen for message hub messages
	// that will get passed once there is a database change to
	// a process flow node display.
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()

		for {
			payload := <-hubChannel
			// This might not work if we have to send these messages over pub/sub instead?
			// Or we can enforce the type in the message hub somehow?
			typedPayload := payload.(ProcessFlowNodeDisplaySettingsPayload)
			jsonMessage, err := json.Marshal(typedPayload)
			if err != nil {
				core.Warning("Failed to marshal paylod: " + err.Error())
				break
			}
			err = conn.WriteMessage(websocket.TextMessage, jsonMessage)
			if err != nil {
				core.Warning("Failed to write message to websocket: " + err.Error())
				break
			}
		}
	}()

	waitGroup.Wait()
}
