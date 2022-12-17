package wshandler

import (
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
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

	res, err := f(b, jst.Now(), target)
	if err != nil {
		return err
	}

	if err := h.sender.Broadcast(h.room.ID, res); err != nil {
		return err
	}

	return nil
}

func (h *wsHandler) handleYolo(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	cardType := domain.CardTypeYolo

	targetRailID, ok := getNonBlockingRailID(targetPlayer, false)
	if !ok {
		targetPlayer.JustCardEvents = append(targetPlayer.JustCardEvents, domain.NewJustCardEvent(
			uuid.New(),
			cardType,
			now,
		))

		return oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now), nil
	}

	afterRails := []*domain.Rail{targetPlayer.Main}
	if l := len(targetPlayer.BranchEvents); l > 0 {
		afterRails = targetPlayer.BranchEvents[l-1].AfterRails
	}

	for _, r := range afterRails {
		if r.ID != targetRailID {
			afterRails = append(afterRails, r)
		}
	}

	targetPlayer.BranchEvents = append(targetPlayer.BranchEvents, domain.NewBranchEvent(
		uuid.New(),
		cardType,
		now,
		domain.BranchEventCreated,
		h.playerID,
		targetPlayer.ID,
		afterRails,
	))

	res, err := oapi.NewWsResponseRailMerged(now, targetRailID, targetPlayer.Main.ID, h.playerID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handleGalaxyBrain(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	cardType := domain.CardTypeGalaxyBrain

	targetPlayer.JustCardEvents = append(targetPlayer.JustCardEvents, domain.NewJustCardEvent(
		uuid.New(),
		cardType,
		now,
	))

	res := oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now)

	return res, nil
}

func (h *wsHandler) handleOpenSourcerer(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	cardType := domain.CardTypeOpenSourcerer

	targetPlayer.LifeEvents = append(targetPlayer.LifeEvents, domain.NewLifeEvent(
		uuid.New(),
		cardType,
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
	cardType := domain.CardTypeRefactoring

	if h.playerID != targetPlayer.ID {
		return nil, errors.New("targetID is different from playerID")
	}

	targetRailID, ok := getNonBlockingRailID(targetPlayer, true)
	if !ok {
		targetPlayer.JustCardEvents = append(targetPlayer.JustCardEvents, domain.NewJustCardEvent(
			uuid.New(),
			cardType,
			now,
		))

		return oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now), nil
	}

	delay, attack, err := cardType.DelayAndAttack()
	if err != nil {
		return nil, err
	}

	targetPlayer.BlockEvents = append(targetPlayer.BlockEvents, domain.NewBlockEvent(
		uuid.New(),
		cardType,
		now,
		domain.BlockEventTypeCreated,
		h.playerID,
		targetPlayer.ID,
		targetRailID,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, delay, attack)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handlePairExtraordinaire(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	cardType := domain.CardTypePairExtraordinaire

	targetRailID, ok := getNonBlockingRailID(targetPlayer, true)
	if !ok {
		targetPlayer.JustCardEvents = append(targetPlayer.JustCardEvents, domain.NewJustCardEvent(
			uuid.New(),
			cardType,
			now,
		))

		return oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now), nil
	}

	delay, attack, err := cardType.DelayAndAttack()
	if err != nil {
		return nil, err
	}

	targetPlayer.BlockEvents = append(targetPlayer.BlockEvents, domain.NewBlockEvent(
		uuid.New(),
		cardType,
		now,
		domain.BlockEventTypeCreated,
		h.playerID,
		targetPlayer.ID,
		targetRailID,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, delay, attack)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handleLgtm(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	cardType := domain.CardTypeLgtm

	targetRailID, ok := getNonBlockingRailID(targetPlayer, true)
	if !ok {
		targetPlayer.JustCardEvents = append(targetPlayer.JustCardEvents, domain.NewJustCardEvent(
			uuid.New(),
			cardType,
			now,
		))

		return oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now), nil
	}

	delay, attack, err := cardType.DelayAndAttack()
	if err != nil {
		return nil, err
	}

	targetPlayer.BlockEvents = append(targetPlayer.BlockEvents, domain.NewBlockEvent(
		uuid.New(),
		cardType,
		now,
		domain.BlockEventTypeCreated,
		h.playerID,
		targetPlayer.ID,
		targetRailID,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, delay, attack)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handlePullShark(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	cardType := domain.CardTypePullShark

	afterRails := []*domain.Rail{targetPlayer.Main}
	if l := len(targetPlayer.BranchEvents); l > 0 {
		afterRails = append(targetPlayer.BranchEvents[l-1].AfterRails, domain.NewRail())
	}

	targetPlayer.BranchEvents = append(targetPlayer.BranchEvents, domain.NewBranchEvent(
		uuid.New(),
		cardType,
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
	cardType := domain.CardTypeStarstruck

	targetRailID, ok := getNonBlockingRailID(targetPlayer, true)
	if !ok {
		targetPlayer.JustCardEvents = append(targetPlayer.JustCardEvents, domain.NewJustCardEvent(
			uuid.New(),
			cardType,
			now,
		))

		return oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now), nil
	}

	delay, attack, err := cardType.DelayAndAttack()
	if err != nil {
		return nil, err
	}

	targetPlayer.BlockEvents = append(targetPlayer.BlockEvents, domain.NewBlockEvent(
		uuid.New(),
		cardType,
		now,
		domain.BlockEventTypeCreated,
		h.playerID,
		targetPlayer.ID,
		targetRailID,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, delay, attack)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// getNonBlockingRailID 現状のレールからランダムにブロック対象のレールを取得する
// 既にブロックされているブランチは取得できない
func getNonBlockingRailID(p *domain.Player, allowMain bool) (uuid.UUID, bool) {
	// 既にブロックされているブランチのIDを取得
	blockBranchIDs := make(map[uuid.UUID]struct{})

	for _, e := range p.BlockEvents {
		switch e.Type {
		case domain.BlockEventTypeCreated:
			blockBranchIDs[e.TargetRailID] = struct{}{}
		case domain.BlockEventTypeCanceled:
			delete(blockBranchIDs, e.TargetRailID)
		}
	}

	railIDs := []uuid.UUID{p.Main.ID}

	if l := len(p.BranchEvents); l > 0 {
		rails := p.BranchEvents[l-1].AfterRails
		railIDs = make([]uuid.UUID, len(rails))

		for i, r := range rails {
			railIDs[i] = r.ID
		}

		rand.Shuffle(len(railIDs), func(i, j int) {
			railIDs[i], railIDs[j] = railIDs[j], railIDs[i]
		})
	}

	for _, id := range railIDs {
		if _, ok := blockBranchIDs[id]; !ok {
			if !allowMain && id == p.Main.ID {
				continue
			}

			return id, true
		}
	}

	return uuid.Nil, false
}
