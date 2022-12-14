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

func NewWsResponseLifeChanged(eventTime time.Time, playerID uuid.UUID, cardType *CardType, newLife float32) (*WsResponse, error) {
	b := WsResponseBodyLifeChanged{
		PlayerId: playerID,
		CardType: cardType,
		NewLife:  newLife,
	}

	res := WsResponseFromType(WsResponseTypeLifeChanged, eventTime)

	if err := res.Body.FromWsResponseBodyLifeChanged(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseRailCreated(eventTime time.Time, newRail, parentRail RailIndex, attackerID, targetID uuid.UUID, cardType CardType) (*WsResponse, error) {
	b := WsResponseBodyRailCreated{
		NewRail:    newRail,
		ParentRail: parentRail,
		AttackerId: attackerID,
		TargetId:   targetID,
		CardType:   cardType,
	}

	res := WsResponseFromType(WsResponseTypeRailCreated, eventTime)

	if err := res.Body.FromWsResponseBodyRailCreated(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseRailMerged(eventTime time.Time, childRail, parentRail RailIndex, playerID uuid.UUID, cardType CardType) (*WsResponse, error) {
	b := WsResponseBodyRailMerged{
		ChildRail:  childRail,
		ParentRail: parentRail,
		PlayerId:   playerID,
		CardType:   cardType,
	}

	res := WsResponseFromType(WsResponseTypeRailMerged, eventTime)

	if err := res.Body.FromWsResponseBodyRailMerged(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseBlockCreated(eventTime time.Time, attackerID uuid.UUID, targetID uuid.UUID, cardType CardType, railIndex int, delay int, attack float32) (*WsResponse, error) {
	b := WsResponseBodyBlockCreated{
		AttackerId: attackerID,
		TargetId:   targetID,
		CardType:   cardType,
		RailIndex:  railIndex,
		Delay:      delay,
		Attack:     attack,
	}

	res := WsResponseFromType(WsResponseTypeBlockCreated, eventTime)

	if err := res.Body.FromWsResponseBodyBlockCreated(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseBlockCanceled(eventTime time.Time, targetID uuid.UUID, rail RailIndex, cardType *CardType) (*WsResponse, error) {
	b := WsResponseBodyBlockCanceled{
		TargetId:  targetID,
		RailIndex: rail,
		CardType:  cardType,
	}

	res := WsResponseFromType(WsResponseTypeBlockCanceled, eventTime)

	if err := res.Body.FromWsResponseBodyBlockCanceled(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseBlockCrashed(eventTime time.Time, targetID uuid.UUID, rail RailIndex, cardType *CardType) (*WsResponse, error) {
	b := WsResponseBodyBlockCrashed{
		TargetId:  targetID,
		RailIndex: rail,
		CardType:  cardType,
	}

	res := WsResponseFromType(WsResponseTypeBlockCrashed, eventTime)

	if err := res.Body.FromWsResponseBodyBlockCrashed(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseGameOverred(eventTime time.Time, playerID uuid.UUID) (*WsResponse, error) {
	b := WsResponseBodyGameOverred{
		PlayerId: playerID,
	}

	res := WsResponseFromType(WsResponseTypeGameOverred, eventTime)

	if err := res.Body.FromWsResponseBodyGameOverred(b); err != nil {
		return nil, err
	}

	return res, nil
}
