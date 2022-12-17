package ws_test

import (
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

var (
	roomRepo = repoimpl.NewRoomRepository()
	streamer = ws.NewStreamer(ws.NewHub(roomRepo))
)

func TestWs(t *testing.T) {
	var (
		conns = make([]*websocket.Conn, consts.PlayerLimit)
		ps    = []*domain.Player{
			domain.NewPlayer(uuid.MustParse("194a6a3c-4278-4be4-ba8e-2c3528d92b8f"), "player0"),
			domain.NewPlayer(uuid.MustParse("cf8b3659-31c4-439b-88b6-4d90dc7b6df9"), "player1"),
			domain.NewPlayer(uuid.MustParse("b46b92b4-ad86-43e9-93d6-216dd9efefd7"), "player2"),
			domain.NewPlayer(uuid.MustParse("07a86224-7419-47b7-965e-5ed4a6b05b22"), "player3"),
		}
		wg = new(sync.WaitGroup)
	)

	// 部屋を作成 & 全員が参加
	room, err := roomRepo.CreateRoom(ps[0])
	require.NoError(t, err)
	require.NoError(t, roomRepo.JoinRoom(room.ID, ps[1]))
	require.NoError(t, roomRepo.JoinRoom(room.ID, ps[2]))
	require.NoError(t, roomRepo.JoinRoom(room.ID, ps[3]))

	// n個のクライアントをWebsocketに接続
	for i := 0; i < consts.PlayerLimit; i++ {
		// Websocketクライアントを接続
		server := httptest.NewServer(&httpHandler{t, streamer, ps[i].ID})
		server.URL = "ws" + strings.TrimPrefix(server.URL, "http")
		c, _, err := websocket.DefaultDialer.Dial(server.URL, nil)
		require.NoError(t, err)

		conns[i] = c
		defer c.Close()

		// 接続を確認
		readWsResponse[oapi.WsResponseBodyConnected](t, c).
			Equal(
				oapi.WsResponseTypeConnected,
				oapi.WsResponseBodyConnected{PlayerId: ps[i].ID},
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
						{Id: ps[0].ID, Life: consts.MaxLife, MainRail: oapi.Rail{Index: 3}, Rails: []oapi.Rail{{}, {}, {}, {Index: 3}, {}, {}, {}}},
						{Id: ps[1].ID, Life: consts.MaxLife, MainRail: oapi.Rail{Index: 3}, Rails: []oapi.Rail{{}, {}, {}, {Index: 3}, {}, {}, {}}},
						{Id: ps[2].ID, Life: consts.MaxLife, MainRail: oapi.Rail{Index: 3}, Rails: []oapi.Rail{{}, {}, {}, {Index: 3}, {}, {}, {}}},
						{Id: ps[3].ID, Life: consts.MaxLife, MainRail: oapi.Rail{Index: 3}, Rails: []oapi.Rail{{}, {}, {}, {Index: 3}, {}, {}, {}}},
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
				TargetId: ps[0].ID,
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
						AttackerId: ps[1].ID,
						CardType:   oapi.CardTypePullShark,
						NewRail:    oapi.Rail{Index: randint(6)},
						ParentRail: oapi.Rail{Index: 3},
						TargetId:   ps[0].ID,
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
				TargetId: ps[0].ID,
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
						AttackerId: ps[1].ID,
						CardType:   oapi.CardTypePairExtraordinaire,
						Delay:      2,
						TargetId:   ps[0].ID,
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
						PlayerId: ps[0].ID,
						NewLife:  99,
					},
				)
		})
	})
}
