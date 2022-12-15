package ws_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestWs(t *testing.T) {
	var (
		conns = make([]*websocket.Conn, 4)
		pids  = make([]uuid.UUID, 4)
	)

	// Streamerを起動
	h := ws.NewHub()
	s := ws.NewStreamer(h, echo.New().Logger) // TODO: loggerのためにechoを使っているのを直す

	for i := 0; i < 4; i++ {
		pids[i] = uuid.New()

		// Websocketクライアントを接続
		server := httptest.NewServer(&httpHandler{t, s, pids[i]})
		server.URL = "ws" + strings.TrimPrefix(server.URL, "http")
		c, _, err := websocket.DefaultDialer.Dial(server.URL, nil)
		require.NoError(t, err)

		conns[i] = c
		defer c.Close()

		// 接続を確認
		res := readWsResponse(t, c)
		resbody, err := res.Body.AsWsResponseBodyConnected()
		require.NoError(t, err)
		require.Equal(t, oapi.WsResponseTypeConnected, res.Type)
		require.Equal(t, oapi.WsResponseBodyConnected{PlayerId: pids[i]}, resbody)
	}
}
