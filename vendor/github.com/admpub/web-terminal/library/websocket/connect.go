package websocket

import (
	"net/http"

	"github.com/admpub/websocket"
	"github.com/pkg/errors"
)

var DefaultUpgrader = websocket.Upgrader{
	// Cross/cors origin domain
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	// Resolve: Sec-WebSocket-Protocol Header
	//Subprotocols: []string{r.Header.Get("Sec-WebSocket-Protocol")},
	Subprotocols:    []string{"web-terminal"},
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
}

func Connect(
	w http.ResponseWriter, req *http.Request,
	onInit func(w Writer) error,
	onRecv func(w Writer, msgType int, data []byte) error,
	upgraders ...websocket.Upgrader,
) error {
	upgrader := DefaultUpgrader
	if len(upgraders) > 0 {
		upgrader = upgraders[0]
	}

	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	if err = onInit(ws); err != nil {
		return err
	}
	for {
		msgType, data, err := ws.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "websocket close or read message err")
		}

		if err = onRecv(ws, msgType, data); err != nil {
			return err
		}
	}
}
