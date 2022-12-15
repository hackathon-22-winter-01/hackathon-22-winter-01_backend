package oapi

import "github.com/google/uuid"

func WsResponseFromType(typ WsResponseType) *WsResponse {
	return &WsResponse{
		Type: typ,
	}
}

func NewWsResponseConnected(playerID uuid.UUID) (*WsResponse, error) {
	b := WsResponseBodyConnected{
		PlayerId: playerID,
	}

	res := WsResponseFromType(WsResponseTypeConnected)
	if err := res.Body.FromWsResponseBodyConnected(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseGameStarted(cards []Card, players []Player) (*WsResponse, error) {
	b := WsResponseBodyGameStarted{
		Cards:   cards,
		Players: players,
	}

	res := WsResponseFromType(WsResponseTypeGameStarted)
	if err := res.Body.FromWsResponseBodyGameStarted(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseLifeChanged() (*WsResponse, error) {
	b := WsResponseBodyLifeChanged{}

	res := WsResponseFromType(WsResponseTypeLifeChanged)
	if err := res.Body.FromWsResponseBodyLifeChanged(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseCardReset() (*WsResponse, error) {
	b := WsResponseBodyCardReset{}

	res := WsResponseFromType(WsResponseTypeCardReset)
	if err := res.Body.FromWsResponseBodyCardReset(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseRailCreated() (*WsResponse, error) {
	b := WsResponseBodyRailCreated{}

	res := WsResponseFromType(WsResponseTypeRailCreated)
	if err := res.Body.FromWsResponseBodyRailCreated(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseRailMerged() (*WsResponse, error) {
	b := WsResponseBodyRailMerged{}

	res := WsResponseFromType(WsResponseTypeRailMerged)
	if err := res.Body.FromWsResponseBodyRailMerged(b); err != nil {
		return nil, err
	}

	return res, nil
}

func NewWsResponseBlockCreated(attackerID uuid.UUID, targetID uuid.UUID) (*WsResponse, error) {
	b := WsResponseBodyBlockCreated{
		AttackerId: attackerID,
		TargetId:   targetID,
	}

	res := WsResponseFromType(WsResponseTypeBlockCreated)
	if err := res.Body.FromWsResponseBodyBlockCreated(b); err != nil {
		return nil, err
	}

	return res, nil
}
