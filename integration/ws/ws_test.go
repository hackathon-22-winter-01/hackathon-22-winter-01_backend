// nolint wsl
package ws_test

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository/repoimpl"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/config"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/optional"
	"github.com/shiguredo/websocket"
	"github.com/stretchr/testify/require"
)

// 短縮用の型定義

const (
	tGameStartEvent  = oapi.WsRequestTypeGameStartEvent
	tLifeEvent       = oapi.WsRequestTypeLifeEvent
	tCardEvent       = oapi.WsRequestTypeCardEvent
	tBlockEvent      = oapi.WsRequestTypeBlockEvent
	tCardForAllEvent = oapi.WsRequestTypeCardForAllEvent

	tConnected     = oapi.WsResponseTypeConnected
	tGameStarted   = oapi.WsResponseTypeGameStarted
	tLifeChanged   = oapi.WsResponseTypeLifeChanged
	tRailCreated   = oapi.WsResponseTypeRailCreated
	tRailMerged    = oapi.WsResponseTypeRailMerged
	tBlockCreated  = oapi.WsResponseTypeBlockCreated
	tBlockCanceled = oapi.WsResponseTypeBlockCanceled
	tBlockCrashed  = oapi.WsResponseTypeBlockCrashed
	tGameOverred   = oapi.WsResponseTypeGameOverred
	tNoop          = oapi.WsResponseTypeNoop
)

type (
	bGameStartEvent  = oapi.WsRequestBodyGameStartEvent
	bLifeEvent       = oapi.WsRequestBodyLifeEvent
	bCardEvent       = oapi.WsRequestBodyCardEvent
	bBlockEvent      = oapi.WsRequestBodyBlockEvent
	bCardForAllEvent = oapi.WsRequestBodyCardForAllEvent

	bConnected     = oapi.WsResponseBodyConnected
	bGameStarted   = oapi.WsResponseBodyGameStarted
	bLifeChanged   = oapi.WsResponseBodyLifeChanged
	bRailCreated   = oapi.WsResponseBodyRailCreated
	bRailMerged    = oapi.WsResponseBodyRailMerged
	bBlockCreated  = oapi.WsResponseBodyBlockCreated
	bBlockCanceled = oapi.WsResponseBodyBlockCanceled
	bBlockCrashed  = oapi.WsResponseBodyBlockCrashed
	bGameOverred   = oapi.WsResponseBodyGameOverred
)

var (
	roomRepo = repoimpl.NewRoomRepository()
	streamer = ws.NewStreamer(ws.NewHub(roomRepo))
)

func TestWs(t *testing.T) {
	var (
		conns = make([]*websocket.Conn, config.PlayerLimit)
		ps    = []*domain.Player{
			domain.NewPlayer(uuid.MustParse("194a6a3c-4278-4be4-ba8e-2c3528d92b8f"), "player0"),
			domain.NewPlayer(uuid.MustParse("cf8b3659-31c4-439b-88b6-4d90dc7b6df9"), "player1"),
			domain.NewPlayer(uuid.MustParse("b46b92b4-ad86-43e9-93d6-216dd9efefd7"), "player2"),
			domain.NewPlayer(uuid.MustParse("07a86224-7419-47b7-965e-5ed4a6b05b22"), "player3"),
		}
		wg = new(sync.WaitGroup)
	)

	// 全員のクライアントをWebsocketに接続&確認
	var roomID *uuid.UUID
	for i := 0; i < config.PlayerLimit; i++ {
		c := connectToWs(t, streamer, ws.ServeWsOpts{
			PlayerID:   ps[i].ID,
			PlayerName: ps[i].Name,
			RoomID:     optional.NewFromPtr(roomID),
		})
		conns[i] = c
		defer c.Close()

		// 既に参加しているメンバー全員にプレイヤーiの接続通知を送信
		forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
			if c == nil {
				return
			}

			readWsResponse[bConnected](t, c).
				Equal(tConnected, bConnected{
					PlayerId: ps[i].ID,
				})
		})

		if i == 0 {
			room, err := roomRepo.FindRoomFromPlayerID(ps[i].ID)
			require.NoError(t, err)
			roomID = &room.ID
		}
	}

	// オーナーがゲーム開始リクエストを送信
	oapi.WriteWsRequest(t, conns[0], tGameStartEvent, bGameStartEvent{})

	// 各プレイヤーはゲーム開始通知を受信
	forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
		readWsResponse[bGameStarted](t, c).
			Equal(tGameStarted, bGameStarted{
				Players: []oapi.Player{
					{Id: ps[0].ID, Name: "player0", Life: config.MaxLife},
					{Id: ps[1].ID, Name: "player1", Life: config.MaxLife},
					{Id: ps[2].ID, Name: "player2", Life: config.MaxLife},
					{Id: ps[3].ID, Name: "player3", Life: config.MaxLife},
				},
			})
	})

	// プレイヤー1がプレイヤー0に対して"Pull Shark"カードを出す
	// 敵のレールを増やす
	oapi.WriteWsRequest(t, conns[1], tCardEvent, bCardEvent{
		Id:       uuid.New(),
		TargetId: ps[0].ID,
		Type:     oapi.CardTypePullShark,
	})

	// 各プレイヤーは結果を受信する
	forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
		readWsResponse[bRailCreated](t, c).
			Equal(tRailCreated, bRailCreated{
				AttackerId: ps[1].ID,
				CardType:   oapi.CardTypePullShark,
				NewRail:    randint(6),
				ParentRail: 3,
				TargetId:   ps[0].ID,
			})
	})

	// プレイヤー1がプレイヤー0に対して"Pair Extraordinaire"カードを出す
	// レールに妨害を発生させる
	oapi.WriteWsRequest(t, conns[1], tCardEvent, bCardEvent{
		Id:       uuid.New(),
		TargetId: ps[0].ID,
		Type:     oapi.CardTypePairExtraordinaire,
	})

	// 各プレイヤーは結果を受信する
	forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
		readWsResponse[bBlockCreated](t, c).
			Equal(tBlockCreated, bBlockCreated{
				Attack:     30,
				AttackerId: ps[1].ID,
				CardType:   oapi.CardTypePairExtraordinaire,
				RailIndex:  randint(6),
				Delay:      2,
				TargetId:   ps[0].ID,
			})
	})

	// プレイヤー0が"Pair Extraordinaire"の妨害に衝突しライフ減少のリクエストを出す
	cardType := oapi.CardTypePairExtraordinaire
	oapi.WriteWsRequest(t, conns[0], tBlockEvent, bBlockEvent{
		CardType:  &cardType,
		RailIndex: randint(6),
		Type:      oapi.BlockEventTypeCrashed,
	})

	// 各プレイヤーは結果を受信する
	forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
		readWsResponse[bBlockCrashed](t, c).
			Equal(tBlockCrashed, bBlockCrashed{
				CardType:  &cardType,
				RailIndex: randint(6),
				TargetId:  ps[0].ID,
			})

		readWsResponse[bLifeChanged](t, c).
			Equal(tLifeChanged, bLifeChanged{
				CardType: &cardType,
				PlayerId: ps[0].ID,
				NewLife:  70, // = 100 - 30
			})
	})

	// プレイヤー0が自分に対して"YOLO"カードを出す
	// 自分の妨害がないレールをマージする
	oapi.WriteWsRequest(t, conns[0], tCardEvent, bCardEvent{
		Id:       uuid.New(),
		TargetId: ps[0].ID,
		Type:     oapi.CardTypeYolo,
	})

	// 各プレイヤーは結果を受信する
	forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
		readWsResponse[bRailMerged](t, c).
			Equal(tRailMerged, bRailMerged{
				CardType:   oapi.CardTypeYolo,
				ChildRail:  randint(6),
				ParentRail: 3,
				PlayerId:   ps[0].ID,
			})
	})

	// プレイヤー1が"Galaxy Brain"カードを出す
	// バックエンドでは何もしないので空のBodyを返す
	oapi.WriteWsRequest(t, conns[1], tCardEvent, bCardEvent{
		Id:       uuid.New(),
		TargetId: ps[0].ID,
		Type:     oapi.CardTypeGalaxyBrain,
	})

	// 各プレイヤーは結果を受信する
	forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
		readWsResponse[any](t, c).
			Equal(tNoop, nil)
	})

	// プレイヤー1が"Refactoring"カードを出す
	// 自身のレールに妨害を発生させる
	oapi.WriteWsRequest(t, conns[1], tCardEvent, bCardEvent{
		Id:       uuid.New(),
		TargetId: ps[1].ID,
		Type:     oapi.CardTypeRefactoring,
	})

	// 各プレイヤーは結果を受信する
	forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
		readWsResponse[bBlockCreated](t, c).
			Equal(tBlockCreated, bBlockCreated{
				Attack:     5,
				AttackerId: ps[1].ID,
				CardType:   oapi.CardTypeRefactoring,
				RailIndex:  3, // 1本(main)しかない
				Delay:      1,
				TargetId:   ps[1].ID,
			})
	})

	// プレイヤー2がプレイヤー1に対して"LGTM"カードを出す
	// 可能ならレールに妨害を発生させるが全てのレールに既に妨害があるので無効となる
	oapi.WriteWsRequest(t, conns[2], tCardEvent, bCardEvent{
		Id:       uuid.New(),
		TargetId: ps[1].ID,
		Type:     oapi.CardTypeLgtm,
	})

	// 各プレイヤーは結果を受信する
	forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
		readWsResponse[any](t, c).
			Equal(tNoop, nil)
	})

	// プレイヤー0が"Refactoring"の妨害に衝突しライフ減少のリクエストを出す
	cardType = oapi.CardTypeRefactoring
	oapi.WriteWsRequest(t, conns[0], tBlockEvent, bBlockEvent{
		CardType:  &cardType,
		RailIndex: 3,
		Type:      oapi.BlockEventTypeCrashed,
	})

	// 各プレイヤーは結果を受信する
	forEachClientAsync(t, wg, conns, func(_ int, c *websocket.Conn) {
		readWsResponse[bBlockCrashed](t, c).
			Equal(tBlockCrashed, bBlockCrashed{
				CardType:  &cardType,
				RailIndex: 3,
				TargetId:  ps[0].ID,
			})

		readWsResponse[bLifeChanged](t, c).
			Equal(tLifeChanged, bLifeChanged{
				CardType: &cardType,
				PlayerId: ps[0].ID,
				NewLife:  65, // = 70 - 5
			})
	})
}
