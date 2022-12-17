// Package oapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package oapi

import (
	"encoding/json"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

// Defines values for BlockEventType.
const (
	BlockEventTypeCanceled BlockEventType = "canceled"
	BlockEventTypeCrashed  BlockEventType = "crashed"
)

// Defines values for CardType.
const (
	CardTypeGalaxyBrain        CardType = "galaxyBrain"
	CardTypeLgtm               CardType = "lgtm"
	CardTypeOoops              CardType = "ooops"
	CardTypeOpenSourcerer      CardType = "openSourcerer"
	CardTypePairExtraordinaire CardType = "pairExtraordinaire"
	CardTypePullShark          CardType = "pullShark"
	CardTypeRefactoring        CardType = "refactoring"
	CardTypeStarstruck         CardType = "starstruck"
	CardTypeYolo               CardType = "yolo"
)

// Defines values for LifeEventType.
const (
	LifeEventTypeDamaged LifeEventType = "damaged"
	LifeEventTypeHealed  LifeEventType = "healed"
)

// Defines values for WsRequestType.
const (
	WsRequestTypeBlockEvent      WsRequestType = "blockEvent"
	WsRequestTypeCardEvent       WsRequestType = "cardEvent"
	WsRequestTypeCardEventForAll WsRequestType = "cardEventForAll"
	WsRequestTypeGameStartEvent  WsRequestType = "gameStartEvent"
	WsRequestTypeLifeEvent       WsRequestType = "lifeEvent"
)

// Defines values for WsResponseType.
const (
	WsResponseTypeBlockCanceled WsResponseType = "blockCanceled"
	WsResponseTypeBlockCrashed  WsResponseType = "blockCrashed"
	WsResponseTypeBlockCreated  WsResponseType = "blockCreated"
	WsResponseTypeCardReset     WsResponseType = "cardReset"
	WsResponseTypeConnected     WsResponseType = "connected"
	WsResponseTypeGameOverred   WsResponseType = "gameOverred"
	WsResponseTypeGameStarted   WsResponseType = "gameStarted"
	WsResponseTypeLifeChanged   WsResponseType = "lifeChanged"
	WsResponseTypeNoop          WsResponseType = "noop"
	WsResponseTypeRailCreated   WsResponseType = "railCreated"
	WsResponseTypeRailMerged    WsResponseType = "railMerged"
)

// BlockEventType ブロックに関するイベントの種類
type BlockEventType string

// Card カード情報
type Card struct {
	// Id カードUUID
	Id CardId `json:"id"`

	// Type カードの効果の種類
	Type CardType `json:"type"`
}

// CardId カードUUID
type CardId = openapi_types.UUID

// CardType カードの効果の種類
type CardType string

// CreateRoomRequest 部屋を作成する際に送信するリクエスト
type CreateRoomRequest struct {
	// PlayerName プレイヤーの名前
	PlayerName string `json:"playerName"`
}

// JoinRoomRequest 部屋に参加する際に送信するリクエスト
type JoinRoomRequest struct {
	// PlayerName プレイヤーの名前
	PlayerName string `json:"playerName"`

	// RoomId ルームUUID
	RoomId RoomId `json:"roomId"`
}

// LifeEventType ライフに関するイベントの種類
type LifeEventType string

// Player プレイヤー情報
type Player struct {
	// Id プレイヤーUUID
	Id PlayerId `json:"id"`

	// Life プレイヤーのライフ
	Life float32 `json:"life"`

	// MainRail レール情報
	MainRail Rail `json:"mainRail"`

	// Rails プレイヤーのレールのリスト
	Rails []Rail `json:"rails"`
}

// PlayerId プレイヤーUUID
type PlayerId = openapi_types.UUID

// Rail レール情報
type Rail struct {
	// Id レールUUID
	Id RailId `json:"id"`
}

// RailId レールUUID
type RailId = openapi_types.UUID

// Room 部屋情報
type Room struct {
	// Id ルームUUID
	Id RoomId `json:"id"`

	// Players プレイヤーのリスト
	Players []Player `json:"players"`

	// StartedAt ゲーム開始時刻
	StartedAt time.Time `json:"startedAt"`
}

// RoomId ルームUUID
type RoomId = openapi_types.UUID

// RoomResponse 部屋の情報
type RoomResponse struct {
	// PlayerId プレイヤーUUID
	PlayerId PlayerId `json:"playerId"`

	// Room 部屋情報
	Room Room `json:"room"`
}

// WsRequest Websocket接続中にサーバーに送信するオブジェクト
type WsRequest struct {
	// Body イベントの情報
	Body WsRequest_Body `json:"body"`

	// Type イベントの種類
	Type WsRequestType `json:"type"`
}

// WsRequest_Body イベントの情報
type WsRequest_Body struct {
	union json.RawMessage
}

// WsRequestBodyBlockEvent ブロックに関するイベントの情報
type WsRequestBodyBlockEvent struct {
	// RailId レールUUID
	RailId RailId `json:"railId"`

	// Type ブロックに関するイベントの種類
	Type BlockEventType `json:"type"`
}

// WsRequestBodyCardEvent カードに関するイベントの情報
type WsRequestBodyCardEvent struct {
	// Id カードUUID
	Id CardId `json:"id"`

	// TargetId プレイヤーUUID
	TargetId PlayerId `json:"targetId"`

	// Type カードの効果の種類
	Type CardType `json:"type"`
}

// WsRequestBodyCardEventForAll 全プレイヤーに影響を与えるカードに関するイベントの情報
type WsRequestBodyCardEventForAll struct {
	// Id カードUUID
	Id CardId `json:"id"`

	// Type カードの効果の種類
	Type CardType `json:"type"`
}

// WsRequestBodyGameStartEvent ゲーム開始時にサーバーに送信するオブジェクト
type WsRequestBodyGameStartEvent = map[string]interface{}

// WsRequestBodyLifeEvent ライフに関するイベントの情報
type WsRequestBodyLifeEvent struct {
	// Diff ライフの変化量
	Diff float32 `json:"diff"`

	// Type ライフに関するイベントの種類
	Type LifeEventType `json:"type"`
}

// WsRequestType イベントの種類
type WsRequestType string

// WsResponse Websocket接続中にサーバーから受信するオブジェクト
type WsResponse struct {
	// Body イベントの情報
	Body WsResponse_Body `json:"body"`

	// EventTime イベントの発生時刻
	EventTime time.Time `json:"eventTime"`

	// Type イベントの種類
	Type WsResponseType `json:"type"`
}

// WsResponse_Body イベントの情報
type WsResponse_Body struct {
	union json.RawMessage
}

// WsResponseBodyBlockCanceled 障害物の解消情報
type WsResponseBodyBlockCanceled struct {
	// RailId レールUUID
	RailId RailId `json:"railId"`
}

// WsResponseBodyBlockCrashed 障害物と衝突したときの情報
type WsResponseBodyBlockCrashed struct {
	// New 変動後のライフ
	New float32 `json:"new"`

	// PlayerId プレイヤーUUID
	PlayerId PlayerId `json:"playerId"`

	// RailId レールUUID
	RailId RailId `json:"railId"`
}

// WsResponseBodyBlockCreated 新規障害物の作成情報
type WsResponseBodyBlockCreated struct {
	// Attack 障害物と衝突したときに与えるダメージ
	Attack float32 `json:"attack"`

	// AttackerId プレイヤーUUID
	AttackerId PlayerId `json:"attackerId"`

	// Delay 障害物を解消するために必要な秒数
	Delay int `json:"delay"`

	// TargetId プレイヤーUUID
	TargetId PlayerId `json:"targetId"`
}

// WsResponseBodyCardReset 各プレイヤーのカードのリセット情報
type WsResponseBodyCardReset = []struct {
	// Cards リセットされたカードのリスト
	Cards []Card `json:"cards"`

	// PlayerId プレイヤーUUID
	PlayerId PlayerId `json:"playerId"`
}

// WsResponseBodyConnected 接続したプレイヤーのID
type WsResponseBodyConnected struct {
	// PlayerId プレイヤーUUID
	PlayerId PlayerId `json:"playerId"`
}

// WsResponseBodyGameOverred defines model for WsResponseBodyGameOverred.
type WsResponseBodyGameOverred struct {
	// PlayerId プレイヤーUUID
	PlayerId PlayerId `json:"playerId"`
}

// WsResponseBodyGameStarted ゲーム開始時の情報
type WsResponseBodyGameStarted struct {
	// Players 各プレイヤーの情報
	Players []Player `json:"players"`
}

// WsResponseBodyLifeChanged ライフの変動情報
type WsResponseBodyLifeChanged struct {
	// New 変動後のライフ
	New float32 `json:"new"`

	// PlayerId プレイヤーUUID
	PlayerId PlayerId `json:"playerId"`
}

// WsResponseBodyRailCreated 新規レールの作成情報
type WsResponseBodyRailCreated struct {
	// AttackerId プレイヤーUUID
	AttackerId PlayerId `json:"attackerId"`

	// Id レールUUID
	Id RailId `json:"id"`

	// ParentId レールUUID
	ParentId RailId `json:"parentId"`

	// TargetId プレイヤーUUID
	TargetId PlayerId `json:"targetId"`
}

// WsResponseBodyRailMerged レールのマージ情報
type WsResponseBodyRailMerged struct {
	// ChildId レールUUID
	ChildId RailId `json:"childId"`

	// ParentId レールUUID
	ParentId RailId `json:"parentId"`

	// PlayerId プレイヤーUUID
	PlayerId PlayerId `json:"playerId"`
}

// WsResponseType イベントの種類
type WsResponseType string

// ConnectToWsParams defines parameters for ConnectToWs.
type ConnectToWsParams struct {
	// PlayerId ユーザーUUID
	PlayerId PlayerId `form:"playerId" json:"playerId"`
}

// JoinRoomJSONRequestBody defines body for JoinRoom for application/json ContentType.
type JoinRoomJSONRequestBody = JoinRoomRequest

// CreateRoomJSONRequestBody defines body for CreateRoom for application/json ContentType.
type CreateRoomJSONRequestBody = CreateRoomRequest

// AsWsRequestBodyGameStartEvent returns the union data inside the WsRequest_Body as a WsRequestBodyGameStartEvent
func (t WsRequest_Body) AsWsRequestBodyGameStartEvent() (WsRequestBodyGameStartEvent, error) {
	var body WsRequestBodyGameStartEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsRequestBodyGameStartEvent overwrites any union data inside the WsRequest_Body as the provided WsRequestBodyGameStartEvent
func (t *WsRequest_Body) FromWsRequestBodyGameStartEvent(v WsRequestBodyGameStartEvent) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsRequestBodyGameStartEvent performs a merge with any union data inside the WsRequest_Body, using the provided WsRequestBodyGameStartEvent
func (t *WsRequest_Body) MergeWsRequestBodyGameStartEvent(v WsRequestBodyGameStartEvent) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsWsRequestBodyLifeEvent returns the union data inside the WsRequest_Body as a WsRequestBodyLifeEvent
func (t WsRequest_Body) AsWsRequestBodyLifeEvent() (WsRequestBodyLifeEvent, error) {
	var body WsRequestBodyLifeEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsRequestBodyLifeEvent overwrites any union data inside the WsRequest_Body as the provided WsRequestBodyLifeEvent
func (t *WsRequest_Body) FromWsRequestBodyLifeEvent(v WsRequestBodyLifeEvent) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsRequestBodyLifeEvent performs a merge with any union data inside the WsRequest_Body, using the provided WsRequestBodyLifeEvent
func (t *WsRequest_Body) MergeWsRequestBodyLifeEvent(v WsRequestBodyLifeEvent) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsWsRequestBodyCardEvent returns the union data inside the WsRequest_Body as a WsRequestBodyCardEvent
func (t WsRequest_Body) AsWsRequestBodyCardEvent() (WsRequestBodyCardEvent, error) {
	var body WsRequestBodyCardEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsRequestBodyCardEvent overwrites any union data inside the WsRequest_Body as the provided WsRequestBodyCardEvent
func (t *WsRequest_Body) FromWsRequestBodyCardEvent(v WsRequestBodyCardEvent) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsRequestBodyCardEvent performs a merge with any union data inside the WsRequest_Body, using the provided WsRequestBodyCardEvent
func (t *WsRequest_Body) MergeWsRequestBodyCardEvent(v WsRequestBodyCardEvent) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsWsRequestBodyBlockEvent returns the union data inside the WsRequest_Body as a WsRequestBodyBlockEvent
func (t WsRequest_Body) AsWsRequestBodyBlockEvent() (WsRequestBodyBlockEvent, error) {
	var body WsRequestBodyBlockEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsRequestBodyBlockEvent overwrites any union data inside the WsRequest_Body as the provided WsRequestBodyBlockEvent
func (t *WsRequest_Body) FromWsRequestBodyBlockEvent(v WsRequestBodyBlockEvent) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsRequestBodyBlockEvent performs a merge with any union data inside the WsRequest_Body, using the provided WsRequestBodyBlockEvent
func (t *WsRequest_Body) MergeWsRequestBodyBlockEvent(v WsRequestBodyBlockEvent) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsWsRequestBodyCardEventForAll returns the union data inside the WsRequest_Body as a WsRequestBodyCardEventForAll
func (t WsRequest_Body) AsWsRequestBodyCardEventForAll() (WsRequestBodyCardEventForAll, error) {
	var body WsRequestBodyCardEventForAll
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsRequestBodyCardEventForAll overwrites any union data inside the WsRequest_Body as the provided WsRequestBodyCardEventForAll
func (t *WsRequest_Body) FromWsRequestBodyCardEventForAll(v WsRequestBodyCardEventForAll) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsRequestBodyCardEventForAll performs a merge with any union data inside the WsRequest_Body, using the provided WsRequestBodyCardEventForAll
func (t *WsRequest_Body) MergeWsRequestBodyCardEventForAll(v WsRequestBodyCardEventForAll) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

func (t WsRequest_Body) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *WsRequest_Body) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsWsResponseBodyConnected returns the union data inside the WsResponse_Body as a WsResponseBodyConnected
func (t WsResponse_Body) AsWsResponseBodyConnected() (WsResponseBodyConnected, error) {
	var body WsResponseBodyConnected
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsResponseBodyConnected overwrites any union data inside the WsResponse_Body as the provided WsResponseBodyConnected
func (t *WsResponse_Body) FromWsResponseBodyConnected(v WsResponseBodyConnected) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsResponseBodyConnected performs a merge with any union data inside the WsResponse_Body, using the provided WsResponseBodyConnected
func (t *WsResponse_Body) MergeWsResponseBodyConnected(v WsResponseBodyConnected) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsWsResponseBodyGameStarted returns the union data inside the WsResponse_Body as a WsResponseBodyGameStarted
func (t WsResponse_Body) AsWsResponseBodyGameStarted() (WsResponseBodyGameStarted, error) {
	var body WsResponseBodyGameStarted
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsResponseBodyGameStarted overwrites any union data inside the WsResponse_Body as the provided WsResponseBodyGameStarted
func (t *WsResponse_Body) FromWsResponseBodyGameStarted(v WsResponseBodyGameStarted) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsResponseBodyGameStarted performs a merge with any union data inside the WsResponse_Body, using the provided WsResponseBodyGameStarted
func (t *WsResponse_Body) MergeWsResponseBodyGameStarted(v WsResponseBodyGameStarted) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsWsResponseBodyLifeChanged returns the union data inside the WsResponse_Body as a WsResponseBodyLifeChanged
func (t WsResponse_Body) AsWsResponseBodyLifeChanged() (WsResponseBodyLifeChanged, error) {
	var body WsResponseBodyLifeChanged
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsResponseBodyLifeChanged overwrites any union data inside the WsResponse_Body as the provided WsResponseBodyLifeChanged
func (t *WsResponse_Body) FromWsResponseBodyLifeChanged(v WsResponseBodyLifeChanged) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsResponseBodyLifeChanged performs a merge with any union data inside the WsResponse_Body, using the provided WsResponseBodyLifeChanged
func (t *WsResponse_Body) MergeWsResponseBodyLifeChanged(v WsResponseBodyLifeChanged) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsWsResponseBodyCardReset returns the union data inside the WsResponse_Body as a WsResponseBodyCardReset
func (t WsResponse_Body) AsWsResponseBodyCardReset() (WsResponseBodyCardReset, error) {
	var body WsResponseBodyCardReset
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsResponseBodyCardReset overwrites any union data inside the WsResponse_Body as the provided WsResponseBodyCardReset
func (t *WsResponse_Body) FromWsResponseBodyCardReset(v WsResponseBodyCardReset) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsResponseBodyCardReset performs a merge with any union data inside the WsResponse_Body, using the provided WsResponseBodyCardReset
func (t *WsResponse_Body) MergeWsResponseBodyCardReset(v WsResponseBodyCardReset) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsWsResponseBodyRailCreated returns the union data inside the WsResponse_Body as a WsResponseBodyRailCreated
func (t WsResponse_Body) AsWsResponseBodyRailCreated() (WsResponseBodyRailCreated, error) {
	var body WsResponseBodyRailCreated
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsResponseBodyRailCreated overwrites any union data inside the WsResponse_Body as the provided WsResponseBodyRailCreated
func (t *WsResponse_Body) FromWsResponseBodyRailCreated(v WsResponseBodyRailCreated) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsResponseBodyRailCreated performs a merge with any union data inside the WsResponse_Body, using the provided WsResponseBodyRailCreated
func (t *WsResponse_Body) MergeWsResponseBodyRailCreated(v WsResponseBodyRailCreated) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsWsResponseBodyRailMerged returns the union data inside the WsResponse_Body as a WsResponseBodyRailMerged
func (t WsResponse_Body) AsWsResponseBodyRailMerged() (WsResponseBodyRailMerged, error) {
	var body WsResponseBodyRailMerged
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsResponseBodyRailMerged overwrites any union data inside the WsResponse_Body as the provided WsResponseBodyRailMerged
func (t *WsResponse_Body) FromWsResponseBodyRailMerged(v WsResponseBodyRailMerged) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsResponseBodyRailMerged performs a merge with any union data inside the WsResponse_Body, using the provided WsResponseBodyRailMerged
func (t *WsResponse_Body) MergeWsResponseBodyRailMerged(v WsResponseBodyRailMerged) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsWsResponseBodyBlockCreated returns the union data inside the WsResponse_Body as a WsResponseBodyBlockCreated
func (t WsResponse_Body) AsWsResponseBodyBlockCreated() (WsResponseBodyBlockCreated, error) {
	var body WsResponseBodyBlockCreated
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsResponseBodyBlockCreated overwrites any union data inside the WsResponse_Body as the provided WsResponseBodyBlockCreated
func (t *WsResponse_Body) FromWsResponseBodyBlockCreated(v WsResponseBodyBlockCreated) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsResponseBodyBlockCreated performs a merge with any union data inside the WsResponse_Body, using the provided WsResponseBodyBlockCreated
func (t *WsResponse_Body) MergeWsResponseBodyBlockCreated(v WsResponseBodyBlockCreated) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsWsResponseBodyBlockCanceled returns the union data inside the WsResponse_Body as a WsResponseBodyBlockCanceled
func (t WsResponse_Body) AsWsResponseBodyBlockCanceled() (WsResponseBodyBlockCanceled, error) {
	var body WsResponseBodyBlockCanceled
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsResponseBodyBlockCanceled overwrites any union data inside the WsResponse_Body as the provided WsResponseBodyBlockCanceled
func (t *WsResponse_Body) FromWsResponseBodyBlockCanceled(v WsResponseBodyBlockCanceled) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsResponseBodyBlockCanceled performs a merge with any union data inside the WsResponse_Body, using the provided WsResponseBodyBlockCanceled
func (t *WsResponse_Body) MergeWsResponseBodyBlockCanceled(v WsResponseBodyBlockCanceled) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsWsResponseBodyBlockCrashed returns the union data inside the WsResponse_Body as a WsResponseBodyBlockCrashed
func (t WsResponse_Body) AsWsResponseBodyBlockCrashed() (WsResponseBodyBlockCrashed, error) {
	var body WsResponseBodyBlockCrashed
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsResponseBodyBlockCrashed overwrites any union data inside the WsResponse_Body as the provided WsResponseBodyBlockCrashed
func (t *WsResponse_Body) FromWsResponseBodyBlockCrashed(v WsResponseBodyBlockCrashed) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsResponseBodyBlockCrashed performs a merge with any union data inside the WsResponse_Body, using the provided WsResponseBodyBlockCrashed
func (t *WsResponse_Body) MergeWsResponseBodyBlockCrashed(v WsResponseBodyBlockCrashed) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

// AsWsResponseBodyGameOverred returns the union data inside the WsResponse_Body as a WsResponseBodyGameOverred
func (t WsResponse_Body) AsWsResponseBodyGameOverred() (WsResponseBodyGameOverred, error) {
	var body WsResponseBodyGameOverred
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsResponseBodyGameOverred overwrites any union data inside the WsResponse_Body as the provided WsResponseBodyGameOverred
func (t *WsResponse_Body) FromWsResponseBodyGameOverred(v WsResponseBodyGameOverred) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsResponseBodyGameOverred performs a merge with any union data inside the WsResponse_Body, using the provided WsResponseBodyGameOverred
func (t *WsResponse_Body) MergeWsResponseBodyGameOverred(v WsResponseBodyGameOverred) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(b, t.union)
	t.union = merged
	return err
}

func (t WsResponse_Body) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *WsResponse_Body) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
