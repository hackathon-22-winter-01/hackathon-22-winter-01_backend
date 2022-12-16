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

// Defines values for CardType.
const (
	CardTypeGalaxyBrain        CardType = "galaxyBrain"
	CardTypeLgtm               CardType = "lgtm"
	CardTypeOpenSourcerer      CardType = "openSourcerer"
	CardTypePairExtraordinaire CardType = "pairExtraordinaire"
	CardTypePullShark          CardType = "pullShark"
	CardTypeRefactoring        CardType = "refactoring"
	CardTypeStarstruck         CardType = "starstruck"
	CardTypeYolo               CardType = "yolo"
)

// Defines values for LifeEventType.
const (
	LifeEventTypeDecrement LifeEventType = "decrement"
)

// Defines values for WsRequestType.
const (
	WsRequestTypeCardEvent      WsRequestType = "cardEvent"
	WsRequestTypeGameStartEvent WsRequestType = "gameStartEvent"
	WsRequestTypeLifeEvent      WsRequestType = "lifeEvent"
	WsRequestTypeRailMergeEvent WsRequestType = "railMergeEvent"
)

// Defines values for WsResponseType.
const (
	WsResponseTypeBlockCreated WsResponseType = "blockCreated"
	WsResponseTypeCardReset    WsResponseType = "cardReset"
	WsResponseTypeConnected    WsResponseType = "connected"
	WsResponseTypeGameStarted  WsResponseType = "gameStarted"
	WsResponseTypeLifeChanged  WsResponseType = "lifeChanged"
	WsResponseTypeRailCreated  WsResponseType = "railCreated"
	WsResponseTypeRailMerged   WsResponseType = "railMerged"
)

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
	Life int `json:"life"`

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

// Room defines model for Room.
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

// WsRequestBodyCardEvent カードに関するイベントの情報
type WsRequestBodyCardEvent struct {
	// Id カードUUID
	Id CardId `json:"id"`

	// TargetId プレイヤーUUID
	TargetId PlayerId `json:"targetId"`

	// Type カードの効果の種類
	Type CardType `json:"type"`
}

// WsRequestBodyGameStartEvent ゲーム開始時にサーバーに送信するオブジェクト
type WsRequestBodyGameStartEvent struct {
	// Name プレイヤーの名前
	Name string `json:"name"`
}

// WsRequestBodyLifeEvent ライフに関するイベントの情報
type WsRequestBodyLifeEvent struct {
	// Type ライフに関するイベントの種類
	Type LifeEventType `json:"type"`
}

// WsRequestBodyRailMergeEvent レールのマージに関するイベントの情報
type WsRequestBodyRailMergeEvent struct {
	// ChildId レールUUID
	ChildId RailId `json:"childId"`

	// ParentId レールUUID
	ParentId RailId `json:"parentId"`
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

// WsResponseBodyBlockCreated 新規障害物の作成情報
type WsResponseBodyBlockCreated struct {
	// Attack 障害物と衝突したときに与えるダメージ
	Attack int `json:"attack"`

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

// WsResponseBodyGameStarted ゲーム開始時の情報
type WsResponseBodyGameStarted struct {
	// Cards ゲーム開始時のカードのリスト
	Cards []Card `json:"cards"`

	// Players 各プレイヤーの情報
	Players []Player `json:"players"`
}

// WsResponseBodyLifeChanged ライフの変動情報
type WsResponseBodyLifeChanged struct {
	// New 変動後のライフ
	New int `json:"new"`

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

// AsWsRequestBodyRailMergeEvent returns the union data inside the WsRequest_Body as a WsRequestBodyRailMergeEvent
func (t WsRequest_Body) AsWsRequestBodyRailMergeEvent() (WsRequestBodyRailMergeEvent, error) {
	var body WsRequestBodyRailMergeEvent
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromWsRequestBodyRailMergeEvent overwrites any union data inside the WsRequest_Body as the provided WsRequestBodyRailMergeEvent
func (t *WsRequest_Body) FromWsRequestBodyRailMergeEvent(v WsRequestBodyRailMergeEvent) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeWsRequestBodyRailMergeEvent performs a merge with any union data inside the WsRequest_Body, using the provided WsRequestBodyRailMergeEvent
func (t *WsRequest_Body) MergeWsRequestBodyRailMergeEvent(v WsRequestBodyRailMergeEvent) error {
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

func (t WsResponse_Body) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *WsResponse_Body) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
