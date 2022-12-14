package main

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/stretchr/testify/require"
)

func TestWs(t *testing.T) {
	// 環境変数の設定
	t.Setenv("APP_PORT", ":8081")

	// APIサーバーを起動
	go main()
	time.Sleep(1 * time.Second)

	// Websocketクライアントを接続
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081/api/v1/ws", nil)
	require.NoError(t, err)

	defer c.Close()

	var (
		playerID        uuid.UUID
		res             oapi.WsResponse
		expectedResBody = new(oapi.WsResponse_Body)
	)

	assertResponse := func(typ oapi.WsResponseType) {
		require.Equal(t, typ, res.Type)
		require.Equal(t, expectedResBody, res.Body)
	}

	// 接続を確認 & プレイヤーIDを取得
	require.NoError(t, c.ReadJSON(&res))
	b, err := res.Body.AsWsResponseBodyConnected()
	require.NoError(t, err)

	playerID = b.PlayerId

	require.NoError(t, expectedResBody.FromWsResponseBodyConnected(oapi.WsResponseBodyConnected{PlayerId: playerID}))

	assertResponse(oapi.WsResponseTypeConnected)
}
