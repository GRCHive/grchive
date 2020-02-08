package websocket

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/webcore"
	"net/http"
	"time"
)

const HeartbeatInterval int = 30

var upgrader = websocket.Upgrader{}

type WebsocketHandler = func(conn *websocket.Conn, r *http.Request)
type HTTPHandler = func(w http.ResponseWriter, r *http.Request)

func RegisterPaths(r *mux.Router) {
	s := r.PathPrefix(core.WebsocketPrefix).Subrouter()
	s.Use(webcore.CreateVerifyCSRFMiddleware(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	s.Use(webcore.ObtainUserSessionInContextMiddleware)

	s.HandleFunc(core.WebsocketProcessFlowNodeDisplaySettingsEndpoint,
		createWebsocketWrapper(processProcessFlowNodeDisplaySettings))
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

func createWebsocketWrapper(handler WebsocketHandler) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			core.Warning("Failed to upgrade to websocket: " + err.Error())
			return
		}

		websocketHeartbeat(c)

		defer c.Close()
		handler(c, r)
	}
}
