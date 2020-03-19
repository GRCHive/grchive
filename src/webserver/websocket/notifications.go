package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func processUserNotifications(conn *websocket.Conn, r *http.Request, role *core.Role) {
	userId, err := webcore.GetUserIdFromRequestUrl(r)
	if err != nil {
		core.Warning("Failed to get user id: " + err.Error())
		return
	}

	// Channel to receive communications from the message hub
	// about relevant events to send to the user.
	var hubChannel chan core.MessagePayload = make(chan core.MessagePayload)
	ele := core.DefaultMessageHub.RegisterListener(core.MHUserNotification,
		core.MessageSubtype(strconv.FormatInt(userId, 10)),
		hubChannel)
	defer core.DefaultMessageHub.UnregisterListener(core.MHUserNotification,
		core.MessageSubtype(strconv.FormatInt(userId, 10)),
		ele)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)

	// Wait for us to receive notifications about notifications and pass them on to the user if relevant.
	go func() {
		defer waitGroup.Done()

		for {
			// Try to ping client to detect if the connection is still alive.
			{
				err := conn.WriteControl(
					websocket.PingMessage,
					[]byte{},
					time.Now().UTC().Add(time.Second*time.Duration(5)))
				if err != nil {
					break
				}
			}

			select {
			case payload := <-hubChannel:
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
			case <-time.After(5 * time.Second):
				continue
			}
		}

	}()

	waitGroup.Wait()
}
