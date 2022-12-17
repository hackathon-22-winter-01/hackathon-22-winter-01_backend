package ws_test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository/repoimpl"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/consts"
	"github.com/shiguredo/websocket"
	"github.com/stretchr/testify/require"
)

func TestWs(t *testing.T) {
	var (
		conns = make([]*websocket.Conn, consts.PlayerLimit)
		pids  = make([]uuid.UUID, consts.PlayerLimit)
		wg    = new(sync.WaitGroup)
	)

	// プレイヤーIDのセットアップ
	for i := 0; i < consts.PlayerLimit; i++ {
		pids[i] = uuid.New()
	}

	// 部屋を作成
	roomRepo := repoimpl.NewRoomRepository()
	room, err := roomRepo.CreateRoom(domain.NewPlayer(pids[0], "player0"))
	require.NoError(t, err)

	// プレイヤー1~3を部屋に参加
	for i := 1; i < consts.PlayerLimit; i++ {
		err := roomRepo.JoinRoom(room.ID, domain.NewPlayer(pids[i], fmt.Sprintf("player%d", i)))
		require.NoError(t, err)
	}

	// Streamerを起動
	h := ws.NewHub(roomRepo)
	s := ws.NewStreamer(h)

	// n個のクライアントをWebsocketに接続
	for i := 0; i < consts.PlayerLimit; i++ {
		// Websocketクライアントを接続
		server := httptest.NewServer(&httpHandler{t, s, pids[i]})
		server.URL = "ws" + strings.TrimPrefix(server.URL, "http")
		c, _, err := websocket.DefaultDialer.Dial(server.URL, nil)
		require.NoError(t, err)

		conns[i] = c
		defer c.Close()

		// 接続を確認
		readWsResponse[oapi.WsResponseBodyConnected](t, c).
			Equal(
				oapi.WsResponseTypeConnected,
				oapi.WsResponseBodyConnected{PlayerId: pids[i]},
			)
	}

	// オーナーがゲーム開始リクエストを送信
	b := oapi.WsRequest_Body{}
	require.NoError(t, b.FromWsRequestBodyGameStartEvent(
		oapi.WsRequestBodyGameStartEvent{},
	))
	mustWriteWsRequest(t, conns[0], oapi.WsRequestTypeGameStartEvent, b)

	// 各プレイヤーはゲーム開始通知を受信
	forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
		readWsResponse[oapi.WsResponseBodyGameStarted](t, c).
			Equal(
				oapi.WsResponseTypeGameStarted,
				oapi.WsResponseBodyGameStarted{
					Players: []oapi.Player{
						{Id: pids[0], Life: consts.MaxLife, MainRail: oapi.Rail{Index: 3}, Rails: []oapi.Rail{{}, {}, {}, {Index: 3}, {}, {}, {}}},
						{Id: pids[1], Life: consts.MaxLife, MainRail: oapi.Rail{Index: 3}, Rails: []oapi.Rail{{}, {}, {}, {Index: 3}, {}, {}, {}}},
						{Id: pids[2], Life: consts.MaxLife, MainRail: oapi.Rail{Index: 3}, Rails: []oapi.Rail{{}, {}, {}, {Index: 3}, {}, {}, {}}},
						{Id: pids[3], Life: consts.MaxLife, MainRail: oapi.Rail{Index: 3}, Rails: []oapi.Rail{{}, {}, {}, {Index: 3}, {}, {}, {}}},
					},
				},
				cmpopts.IgnoreFields(oapi.Rail{}, "Id"),
			)
	})

	t.Run("プレイヤー1がプレイヤー0に対してカードを出してレールを生成する", func(t *testing.T) {
		// NOTE:
		// 今のところはt.Parallel()を付けずに実行する
		// t.Parallel()をつける場合はc.Close()を実行しない必要がある

		// プレイヤー1がプレイヤー0に対してカードを出す
		b := oapi.WsRequest_Body{}
		require.NoError(t, b.FromWsRequestBodyCardEvent(
			oapi.WsRequestBodyCardEvent{
				Id:       uuid.New(),
				TargetId: pids[0],
				Type:     oapi.CardTypePullShark,
			},
		))
		mustWriteWsRequest(t, conns[1], oapi.WsRequestTypeCardEvent, b)

		// 各プレイヤーは結果を受信する
		forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
			readWsResponse[oapi.WsResponseBodyRailCreated](t, c).
				Equal(
					oapi.WsResponseTypeRailCreated,
					oapi.WsResponseBodyRailCreated{
						AttackerId: pids[1],
						CardType:   oapi.CardTypePullShark,
						NewRail:    oapi.Rail{Index: randint(6)},
						ParentRail: oapi.Rail{Index: 3},
						TargetId:   pids[0],
					},
					cmpopts.IgnoreFields(oapi.Rail{}, "Id"), // TODO: Idがなくなったら消す
				)
		})
	})

	t.Run("プレイヤー1がプレイヤー0に対してカードを出して障害物を生成する", func(t *testing.T) {
		// プレイヤー1がプレイヤー0に対してカードを出す
		b := oapi.WsRequest_Body{}
		require.NoError(t, b.FromWsRequestBodyCardEvent(
			oapi.WsRequestBodyCardEvent{
				Id:       uuid.New(),
				TargetId: pids[0],
				Type:     oapi.CardTypePairExtraordinaire,
			},
		))
		mustWriteWsRequest(t, conns[1], oapi.WsRequestTypeCardEvent, b)

		// 各プレイヤーは結果を受信する
		forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
			readWsResponse[oapi.WsResponseBodyBlockCreated](t, c).
				Equal(
					oapi.WsResponseTypeBlockCreated,
					oapi.WsResponseBodyBlockCreated{
						Attack:     30,
						AttackerId: pids[1],
						CardType:   oapi.CardTypePairExtraordinaire,
						Delay:      2,
						TargetId:   pids[0],
					},
				)
		})
	})

	t.Run("プレイヤー0が障害物に当たってライフが1減少する", func(t *testing.T) {
		// プレイヤー0がライフ減少のリクエストを出す
		b := oapi.WsRequest_Body{}
		require.NoError(t, b.FromWsRequestBodyLifeEvent(
			oapi.WsRequestBodyLifeEvent{
				Type: oapi.LifeEventTypeDamaged,
				Diff: 1,
			},
		))
		mustWriteWsRequest(t, conns[0], oapi.WsRequestTypeLifeEvent, b)

		// 各プレイヤーは結果を受信する
		forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
			readWsResponse[oapi.WsResponseBodyLifeChanged](t, c).
				Equal(
					oapi.WsResponseTypeLifeChanged,
					oapi.WsResponseBodyLifeChanged{
						CardType: oapi.CardTypeNone,
						PlayerId: pids[0],
						NewLife:  99,
					},
				)
		})
	})
}
