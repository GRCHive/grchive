package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
	"strconv"
	"sync"
)

type ProcessFlowNodeDisplaySettingsPayload struct {
	NodeId   int64
	Settings map[string]interface{}
}

func processProcessFlowNodeDisplaySettings(conn *websocket.Conn, r *http.Request) {
	flowId, err := webcore.GetProcessFlowIdFromRequest(r)
	if err != nil {
		core.Warning("Failed to get flow id: " + err.Error())
		return
	}

	// Channel to receive communications from the message hub
	// about relevant events to send to the user.
	var hubChannel chan core.MessagePayload = make(chan core.MessagePayload)
	ele := core.RegisterListener(core.UpdateDisplaySettingsForProcessFlowNode,
		core.MessageSubtype(strconv.FormatInt(flowId, 10)),
		hubChannel)
	defer core.UnregisterListener(core.UpdateDisplaySettingsForProcessFlowNode,
		core.MessageSubtype(strconv.FormatInt(flowId, 10)),
		ele)

	// The user needs to know what the current settings are so do a query
	// for the display settings of every node in the process flow. I think this
	// needs to happen after the channel gets registered successfully so the
	// user is guaranteed to see all updates.
	nodeSettings, err := database.FindDisplaySettingsForProcessFlow(flowId)
	if err != nil {
		core.Warning("Failed to get initial node settings: " + err.Error())
		return
	}

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

		for nodeId, settings := range nodeSettings {
			jsonMessage, err := json.Marshal(ProcessFlowNodeDisplaySettingsPayload{
				NodeId:   nodeId,
				Settings: settings,
			})
			err = conn.WriteMessage(websocket.TextMessage, jsonMessage)
			if err != nil {
				core.Warning("Failed to write message to websocket: " + err.Error())
				break
			}
		}

		for {
			payload := <-hubChannel
			jsonMessage, err := json.Marshal(payload)
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