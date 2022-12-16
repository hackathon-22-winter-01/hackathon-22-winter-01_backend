package wshandler

import (
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
)

func (h *wsHandler) handleCardEvent(body oapi.WsRequest_Body) error {
	b, err := body.AsWsRequestBodyCardEvent()
	if err != nil {
		return err
	}

	target, ok := h.room.FindPlayer(b.TargetId)
	if !ok {
		return errors.New("player not found")
	}

	fmap := map[oapi.CardType]func(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error){
		oapi.CardTypeYolo:               h.handleYolo,
		oapi.CardTypeGalaxyBrain:        h.handleGalaxyBrain,
		oapi.CardTypeOpenSourcerer:      h.handleOpenSourcerer,
		oapi.CardTypeRefactoring:        h.handleRefactoring,
		oapi.CardTypePairExtraordinaire: h.handlePairExtraordinaire,
		oapi.CardTypeLgtm:               h.handleLgtm,
		oapi.CardTypePullShark:          h.handlePullShark,
		oapi.CardTypeStarstruck:         h.handleStarstruck,
	}

	f, ok := fmap[b.Type]
	if !ok {
		return errors.New("存在しないカードです")
	}

	res, err := f(b, time.Now(), target)
	if err != nil {
		return err
	}

	if err := h.sender.Broadcast(h.room.ID, res); err != nil {
		return err
	}

	return nil
}

func (h *wsHandler) handleYolo(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	rails := []*domain.Rail{targetPlayer.Main}
	if l := len(targetPlayer.BranchEvents); l > 0 {
		rails = targetPlayer.BranchEvents[l-1].AfterRails
	}

	shuffledRails := []*domain.Rail{}
	copy(shuffledRails, rails)

	rand.Shuffle(len(shuffledRails), func(i, j int) {
		shuffledRails[i], shuffledRails[j] = shuffledRails[j], shuffledRails[i]
	})

	var childID uuid.UUID

	for _, rail := range shuffledRails {
		if rail.ID != targetPlayer.Main.ID && !rail.HasBlock {
			childID = rail.ID
			break
		}
	}

	if childID == uuid.Nil {
		return oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now), nil
	}

	afterRails := make([]*domain.Rail, 0, len(rails)-1)

	for _, rail := range rails {
		if rail.ID != childID {
			afterRails = append(afterRails, rail)
		}
	}

	targetPlayer.BranchEvents = append(targetPlayer.BranchEvents, domain.NewBranchEvent(
		uuid.New(),
		domain.CardTypeYolo,
		now,
		domain.BranchEventCreated,
		targetPlayer.Main.ID,
		childID,
		afterRails,
	))

	res, err := oapi.NewWsResponseRailMerged(now, childID, targetPlayer.Main.ID, h.playerID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handleGalaxyBrain(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	// TODO: 何かしらのイベントを追加したい
	res := oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now)

	return res, nil
}

func (h *wsHandler) handleOpenSourcerer(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	targetPlayer.LifeEvents = append(targetPlayer.LifeEvents, domain.NewLifeEvent(
		uuid.New(),
		domain.CardTypeOpenSourcerer,
		now,
		domain.LifeEventTypeHealed,
		30,
	))

	res, err := oapi.NewWsResponseLifeChanged(now, h.playerID, domain.CalculateLife(targetPlayer.LifeEvents))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handleRefactoring(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	if h.playerID != targetPlayer.ID {
		return nil, errors.New("targetID is different from playerID")
	}

	targetPlayer.BlockEvents = append(targetPlayer.BlockEvents, domain.NewBlockEvent(
		uuid.New(),
		domain.CardTypeRefactoring,
		now,
		h.playerID,
		targetPlayer.ID,
		reqbody.TargetId,
		1,
		5,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, 1, 5)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handlePairExtraordinaire(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	targetPlayer.BlockEvents = append(targetPlayer.BlockEvents, domain.NewBlockEvent(
		uuid.New(),
		domain.CardTypePairExtraordinaire,
		now,
		h.playerID,
		targetPlayer.ID,
		reqbody.TargetId,
		2,
		30,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, 2, 30)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handleLgtm(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	targetPlayer.BlockEvents = append(targetPlayer.BlockEvents, domain.NewBlockEvent(
		uuid.New(),
		domain.CardTypeLgtm,
		now,
		h.playerID,
		targetPlayer.ID,
		reqbody.TargetId,
		3,
		20,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, 3, 20)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handlePullShark(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	afterRails := []*domain.Rail{targetPlayer.Main}
	if l := len(targetPlayer.BranchEvents); l > 0 {
		afterRails = append(targetPlayer.BranchEvents[l-1].AfterRails, domain.NewRail())
	}

	targetPlayer.BranchEvents = append(targetPlayer.BranchEvents, domain.NewBranchEvent(
		uuid.New(),
		domain.CardTypePullShark,
		now,
		domain.BranchEventCreated,
		h.playerID,
		targetPlayer.ID,
		afterRails,
	))

	res, err := oapi.NewWsResponseRailCreated(now, uuid.New(), targetPlayer.Main.ID, h.playerID, targetPlayer.ID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handleStarstruck(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	targetPlayer.BlockEvents = append(targetPlayer.BlockEvents, domain.NewBlockEvent(
		uuid.New(),
		domain.CardTypeLgtm,
		now,
		h.playerID,
		targetPlayer.ID,
		reqbody.TargetId,
		5,
		50,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, 5, 50)
	if err != nil {
		return nil, err
	}

	return res, nil
}
