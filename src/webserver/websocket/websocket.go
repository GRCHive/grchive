package websocket

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/webcore"
	"net/http"
)

var upgrader = websocket.Upgrader{}

type WebsocketHandler = func(conn *websocket.Conn)
type HTTPHandler = func(w http.ResponseWriter, r *http.Request)

func RegisterPaths(r *mux.Router) {
	s := r.PathPrefix(core.WebsocketPrefix).Subrouter()
	s.Use(webcore.CreateVerifyCSRFMiddleware(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))

	s.HandleFunc(core.WebsocketProcessFlowNodeDisplaySettingsEndpoint,
		createWebsocketWrapper(processProcessFlowNodeDisplaySettings))
}

func createWebsocketWrapper(handler WebsocketHandler) HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			core.Warning("Failed to upgrade to websocket: " + err.Error())
			return
		}

		defer c.Close()
		handler(c)
	}
}
