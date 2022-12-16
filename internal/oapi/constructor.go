package oapi

import (
	"time"

	"github.com/google/uuid"
)

func WsResponseFromType(typ WsResponseType, eventTime time.Time) *WsResponse {
	return &WsResponse{
		Type:      typ,
		EventTime: eventTime,
	}
}

func NewWsResponseConnected(eventTime time.Time, playerID uuid.UUID) (*WsResponse, error) {
	b := WsResponseBodyConnected{
		PlayerId: playerID,
	}

	res := WsResponseFromType(WsResponseTypeConnected, eventTime)

	if err := res.Body.FromWsResponseBodyConnected(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseGameStarted(eventTime time.Time, players []Player) (*WsResponse, error) {
	b := WsResponseBodyGameStarted{
		Players: players,
	}

	res := WsResponseFromType(WsResponseTypeGameStarted, eventTime)

	if err := res.Body.FromWsResponseBodyGameStarted(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseLifeChanged(eventTime time.Time, playerID uuid.UUID, newLife int) (*WsResponse, error) {
	b := WsResponseBodyLifeChanged{
		PlayerId: playerID,
		New:      newLife,
	}

	res := WsResponseFromType(WsResponseTypeLifeChanged, eventTime)

	if err := res.Body.FromWsResponseBodyLifeChanged(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseCardReset(eventTime time.Time) (*WsResponse, error) {
	b := WsResponseBodyCardReset{}

	res := WsResponseFromType(WsResponseTypeCardReset, eventTime)

	if err := res.Body.FromWsResponseBodyCardReset(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseRailCreated(eventTime time.Time, railID, parentRailID, attackerID, targetID uuid.UUID) (*WsResponse, error) {
	b := WsResponseBodyRailCreated{
		Id:         railID,
		ParentId:   parentRailID,
		AttackerId: attackerID,
		TargetId:   targetID,
	}

	res := WsResponseFromType(WsResponseTypeRailCreated, eventTime)

	if err := res.Body.FromWsResponseBodyRailCreated(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseRailMerged(eventTime time.Time, childID, parentID, playerID uuid.UUID) (*WsResponse, error) {
	b := WsResponseBodyRailMerged{
		ChildId:  childID,
		ParentId: parentID,
		PlayerId: playerID,
	}

	res := WsResponseFromType(WsResponseTypeRailMerged, eventTime)

	if err := res.Body.FromWsResponseBodyRailMerged(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseBlockCreated(eventTime time.Time, attackerID uuid.UUID, targetID uuid.UUID, delay int, attack int) (*WsResponse, error) {
	b := WsResponseBodyBlockCreated{
		AttackerId: attackerID,
		TargetId:   targetID,
		Delay:      delay,
		Attack:     attack,
	}

	res := WsResponseFromType(WsResponseTypeBlockCreated, eventTime)

	if err := res.Body.FromWsResponseBodyBlockCreated(b); err != nil {
		return nil, err
	}

	return res, nil
}
