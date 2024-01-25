package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func initWebSocket(c *gin.Context) {
	socketHandler(c)
}

var websocketSvr func(*websocket.Conn, []byte) []byte

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}} //

func socketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if !ChkErr(err) {
		for {
			messageType, message, err := conn.ReadMessage()
			var _result []byte
			if ChkErrNormal(err) {
				break
			} else {
				_result = websocketSvr(conn, message)
			}
			err = conn.WriteMessage(messageType, _result)
			if ChkErrNormal(err) {
				break
			}
		}
	}
	conn.Close()
}
