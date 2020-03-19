package websocket

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"time"
)

// Send a fairly often heartbeat to make sure we close off websocket connections if they're dead.
const HeartbeatInterval int = 5

var upgrader = websocket.Upgrader{}

type WebsocketHandler = func(conn *websocket.Conn, r *http.Request, role *core.Role)
type HTTPHandler = func(w http.ResponseWriter, r *http.Request)

func RegisterPaths(r *mux.Router) {
	s := r.PathPrefix(core.WebsocketPrefix).Subrouter()
	s.Use(webcore.CreateVerifyCSRFMiddleware(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))

	s.HandleFunc(core.WebsocketProcessFlowNodeDisplaySettingsEndpoint,
		createWebsocketWrapper(processProcessFlowNodeDisplaySettings))

	s.HandleFunc(core.WebsocketUserNotificationsEndpoint,
		createWebsocketWrapper(processUserNotifications))
}

// This is mainly just for keep alive to ensure that we don't trigger
// NGINX or GKE's timeouts.
func websocketHeartbeat(conn *websocket.Conn) {
	go func() {
		for {
			err := conn.WriteControl(
				websocket.PingMessage,
				[]byte{},
				time.Now().UTC().Add(time.Second*time.Duration(5)))

			if err != nil {
				conn.Close()
				core.Warning("Failed to send ping: " + err.Error())
				break
			}

			time.Sleep(time.Second * time.Duration(HeartbeatInterval))
		}
	}()
}

func readWebsocketRole(conn *websocket.Conn) (*core.Role, error) {
	type Payload struct {
		ApiKey string
		OrgId  int32
		Error  error
	}
	ch := make(chan Payload)

	// The first message the user must sent is a message
	// with the user's API key.
	go func() {
		data := Payload{}
		err := conn.ReadJSON(&data)
		if err != nil {
			data.Error = err
		}
		ch <- data
	}()

	// Caveat here is that if the user's API key is expired then
	// there's no recourse. Oh well?
	select {
	case payload := <-ch:
		if payload.Error != nil {
			return nil, payload.Error
		}

		key, err := database.FindApiKey(core.RawApiKey(payload.ApiKey).Hash())
		if err != nil {
			return nil, err
		}

		role, err := webcore.GetRoleFromKey(key, payload.OrgId)
		if err != nil {
			// Not necessarily an error for certain paths.
			return nil, nil
		}

		return role, nil
	case <-time.After(10 * time.Second):
		return nil, errors.New("Timeout")
	}
}

func createWebsocketWrapper(handler WebsocketHandler) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			core.Warning("Failed to upgrade to websocket: " + err.Error())
			return
		}
		defer c.Close()

		role, err := readWebsocketRole(c)
		if err != nil {
			core.Warning("Failed to read role from websocket: " + err.Error())
			return
		}

		websocketHeartbeat(c)

		handler(c, r, role)
	}
}
