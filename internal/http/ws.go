package http

import (
	"bytes"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// Wraps gorilla [websocket.Upgrader] to offer new methods.
type wsUpgrader struct {
	websocket.Upgrader
}

// Wraps gorilla [websocket.Conn] to offer new methods.
type wsConn struct {
	*websocket.Conn
}

var upgrader = wsUpgrader{}

func (u *wsUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*wsConn, error) {
	conn, err := u.Upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	return &wsConn{
		Conn: conn,
	}, nil
}

// Checks if a client wants to upgrade to websocket connection via the "Upgrade" header.
func (u wsUpgrader) WantsToUpgrade(req http.Request) bool {
	upgrade := false
	for _, header := range req.Header["Upgrade"] {
		if header == "websocket" {
			upgrade = true
			break
		}
	}

	return upgrade
}

// Renders a template and writes it in the websocket connection.
func (ws wsConn) WriteTemplate(ectx echo.Context, tpl string, v any) error {
	var buf bytes.Buffer

	err := ectx.Echo().Renderer.Render(&buf, tpl, v, ectx)
	if err != nil {
		return err
	}

	return ws.WriteMessage(websocket.TextMessage, buf.Bytes())
}
