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
		conns     = make([]*websocket.Conn, consts.PlayerLimit)
		pids      = make([]uuid.UUID, consts.PlayerLimit)
		cards     = make([][]oapi.Card, consts.PlayerLimit)
		mainRails = make([]oapi.Rail, consts.PlayerLimit)
		rails     = make([][]oapi.Rail, consts.PlayerLimit)
		wg        = new(sync.WaitGroup)
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

	// 各プレイヤーはゲーム開始通知を受信
	forEachClientAsync(t, wg, conns, func(i int, c *websocket.Conn) {
		res := readWsResponse(t, c)
		resbody, err := res.Body.AsWsResponseBodyGameStarted()
		require.NoError(t, err)
		require.Equal(t, oapi.WsResponseTypeGameStarted, res.Type)
		require.Len(t, resbody.Cards, 5) // ここではCardsの中身は問わない
		require.Len(t, resbody.Players, consts.PlayerLimit)
		for j, p := range resbody.Players {
			require.Equal(t, pids[j], p.PlayerId)
			require.Equal(t, consts.MaxLife, p.Life)

			// レールを記録しておく
			if i == 0 {
				mainRails[j] = p.MainRail
				rails[j] = p.Rails
			}
		}

		// カードを記録しておく
		cards[i] = resbody.Cards
	})

	t.Run("プレイヤー1がプレイヤー0に対してカードを出してレールを生成する", func(t *testing.T) {
		// FIXME: t.Parallel()を付けて実行するとFAILする
		// Received unexpected error:
		// write tcp 127.0.0.1:38046->127.0.0.1:35689: use of closed network connection

		// プレイヤー1がプレイヤー0に対してカードを出す
		card := cards[1][0]
		b := oapi.WsRequest_Body{}
		require.NoError(t, b.FromWsRequestBodyCardEvent(
			oapi.WsRequestBodyCardEvent{
				Id:       card.Id,
				TargetId: pids[0],
				Type:     card.Type,
			},
		))
		mustWriteWsRequest(t, conns[1], oapi.WsRequestTypeCardEvent, b)

		// 各プレイヤーは結果を受信する
		forEachClientAsync(t, wg, conns, func(i int, c *websocket.Conn) {
			res := readWsResponse(t, c)
			resbody, err := res.Body.AsWsResponseBodyRailCreated()
			require.NoError(t, err)
			require.Equal(t, oapi.WsResponseTypeRailCreated, res.Type)
			// 作成されたレールのID以外の確認
			require.Equal(t, mainRails[0].Id, resbody.ParentId)
			require.Equal(t, pids[1], resbody.AttackerId)
			require.Equal(t, pids[0], resbody.TargetId)

			// レールの更新
			if i == 0 {
				rails[0] = append(rails[0], oapi.Rail{Id: resbody.Id})
			}
		})
	})

	t.Run("プレイヤー1がプレイヤー0に対してカードを出して障害物を生成する", func(t *testing.T) {
		// プレイヤー1がプレイヤー0に対してカードを出す
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

		// 各プレイヤーは結果を受信する
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

	t.Run("プレイヤー0が障害物に当たってライフが1減少する", func(t *testing.T) {
		// プレイヤー0がライフ減少のリクエストを出す
		b := oapi.WsRequest_Body{}
		require.NoError(t, b.FromWsRequestBodyLifeEvent(
			oapi.WsRequestBodyLifeEvent{
				Type: oapi.LifeEventTypeDecrement,
			},
		))
		mustWriteWsRequest(t, conns[0], oapi.WsRequestTypeLifeEvent, b)

		// 各プレイヤーは結果を受信する
		forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
			res := readWsResponse(t, c)
			resbody, err := res.Body.AsWsResponseBodyLifeChanged()
			require.NoError(t, err)
			require.Equal(t, oapi.WsResponseTypeLifeChanged, res.Type)
			require.Equal(t, oapi.WsResponseBodyLifeChanged{
				PlayerId: pids[0],
				New:      2,
			}, resbody)
		})
	})

	t.Run("プレイヤー0がレールをmainにマージする", func(t *testing.T) {
		// プレイヤー0がレールマージのリクエストを出す
		b := oapi.WsRequest_Body{}
		require.NoError(t, b.FromWsRequestBodyRailMergeEvent(
			oapi.WsRequestBodyRailMergeEvent{
				ChildId:  rails[0][1].Id, // rails[0] = [main, プレイヤー1のカードで生成されたレール]
				ParentId: mainRails[0].Id,
			},
		))
		mustWriteWsRequest(t, conns[0], oapi.WsRequestTypeRailMergeEvent, b)

		// 各プレイヤーは結果を受信する
		forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
			res := readWsResponse(t, c)
			resbody, err := res.Body.AsWsResponseBodyRailMerged()
			require.NoError(t, err)
			require.Equal(t, oapi.WsResponseTypeRailMerged, res.Type)
			require.Equal(t, oapi.WsResponseBodyRailMerged{
				ChildId:  rails[0][1].Id,
				ParentId: mainRails[0].Id,
				PlayerId: pids[0],
			}, resbody)
		})
	})
}
