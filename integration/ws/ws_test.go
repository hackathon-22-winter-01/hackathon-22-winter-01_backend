package ws_test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository/repoimpl"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/consts"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestWs(t *testing.T) {
	var (
		conns = make([]*websocket.Conn, n)
		pids  = make([]uuid.UUID, n)
		wg    = new(sync.WaitGroup)
	)

	// Streamerを起動
	roomRepo := repoimpl.NewRoomRepository()
	h := ws.NewHub(roomRepo)
	s := ws.NewStreamer(h, echo.New().Logger) // TODO: loggerのためにechoを使っているのを直す

	// n個のクライアントをWebsocketに接続
	for i := 0; i < consts.PlayerLimit; i++ {
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

	// オーナーがゲーム開始リクエストを送信
	b := oapi.WsRequest_Body{}
	require.NoError(t, b.FromWsRequestBodyGameStartEvent(
		oapi.WsRequestBodyGameStartEvent{
			Name: fmt.Sprintf("player%d", 0),
		},
	))
	mustWriteWsRequest(t, conns[0], oapi.WsRequestTypeGameStartEvent, b)

	// 各クライアントはゲーム開始通知を受信
	forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
		players := make([]oapi.Player, consts.PlayerLimit)
		for i, pid := range pids {
			players[i] = oapi.Player{
				PlayerId: pid,
				Life:     3,
			}
		}

		res := readWsResponse(t, c)
		resbody, err := res.Body.AsWsResponseBodyGameStarted()
		require.NoError(t, err)
		require.Equal(t, oapi.WsResponseTypeGameStarted, res.Type)
		require.Equal(t, players, resbody.Players)
		require.Len(t, resbody.Cards, 5) // ここではCardsの中身は問わない
	})
}
