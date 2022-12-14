package ws_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/handler"
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

	// 環境変数の設定
	t.Setenv("APP_PORT", ":8081")

	// APIサーバーを起動
	e := echo.New()
	hub := ws.NewHub()
	st := ws.NewStreamer(hub, e.Logger)
	h := handler.New(st)
	oapi.RegisterHandlersWithBaseURL(e, h, "/api/v1")

	go e.Logger.Fatal(e.Start(":8081"))

	time.Sleep(3 * time.Second)

	for i := 0; i < 4; i++ {
		// Websocketクライアントを接続
		c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081/api/v1/ws", nil)
		require.NoError(t, err)

		defer c.Close()

		conns[i] = c

		// 接続を確認 & プレイヤーIDを取得
		res := readWsResponse(t, c)
		b, err := new(oapi.WsResponse_Body).AsWsResponseBodyConnected()
		require.NoError(t, err)

		pids[i] = b.PlayerId

		expectedResBody := new(oapi.WsResponse_Body)
		require.NoError(t, expectedResBody.FromWsResponseBodyConnected(oapi.WsResponseBodyConnected{PlayerId: pids[i]}))
		require.Equal(t, oapi.WsResponseTypeConnected, res.Type)
		require.Equal(t, expectedResBody, &res.Body)
	}
}
