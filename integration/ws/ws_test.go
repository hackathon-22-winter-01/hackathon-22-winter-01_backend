package ws_test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository/repoimpl"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/consts"
	"github.com/labstack/echo/v4"
	"github.com/shiguredo/websocket"
	"github.com/stretchr/testify/require"
)

func TestWs(t *testing.T) {
	var (
		conns = make([]*websocket.Conn, consts.PlayerLimit)
		pids  = make([]uuid.UUID, consts.PlayerLimit)
		cards = make([][]oapi.Card, consts.PlayerLimit)
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
	forEachClientAsync(t, wg, conns, func(i int, c *websocket.Conn) {
		players := make([]oapi.Player, consts.PlayerLimit)
		for j, pid := range pids {
			players[j] = oapi.Player{
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

		// カードを記録しておく
		cards[i] = resbody.Cards
	})

	t.Run("クライアント1がクライアント0に対してカードを出して障害物を生成する", func(t *testing.T) {
		// FIXME: t.Parallel()を付けて実行するとFAILする
		// Received unexpected error:
		// write tcp 127.0.0.1:38046->127.0.0.1:35689: use of closed network connection

		// クライアント1がクライアント0に対してカードを出す
		card := cards[1][1]
		b := oapi.WsRequest_Body{}
		require.NoError(t, b.FromWsRequestBodyCardEvent(
			oapi.WsRequestBodyCardEvent{
				Id:       card.Id,
				TargetId: pids[0],
				Type:     card.Type,
			},
		))
		mustWriteWsRequest(t, conns[1], oapi.WsRequestTypeCardEvent, b)

		// 各クライアントは結果を受信する
		forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
			res := readWsResponse(t, c)
			resbody, err := res.Body.AsWsResponseBodyBlockCreated()
			require.NoError(t, err)
			require.Equal(t, oapi.WsResponseTypeBlockCreated, res.Type)
			require.Equal(t, oapi.WsResponseBodyBlockCreated{
				AttackerId: pids[1],
				TargetId:   pids[0],
			}, resbody)
		})
	})
}
