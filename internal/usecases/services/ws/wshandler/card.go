package wshandler

import (
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/consts"
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

	targetRail, ok := getNonBlockingRail(targetPlayer, false)
	if !ok {
		targetPlayer.JustCardEvents = append(targetPlayer.JustCardEvents, domain.NewJustCardEvent(
			uuid.New(),
			cardType,
			now,
		))

		return oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now), nil
	}

	afterRails := domain.NewRails(targetPlayer.Main)
	if l := len(targetPlayer.BranchEvents); l > 0 {
		copy(afterRails[:], targetPlayer.BranchEvents[l-1].AfterRails[:])
	}

	var targetRailIndex int

	for i, r := range afterRails {
		if r.ID != targetRail.ID {
			afterRails[i] = nil
			targetRailIndex = i

			break
		}
	}

	targetPlayer.BranchEvents = append(targetPlayer.BranchEvents, domain.NewBranchEvent(
		uuid.New(),
		cardType,
		now,
		domain.BranchEventMerged,
		h.playerID,
		afterRails,
	))

	res, err := oapi.NewWsResponseRailMerged(
		now,
		oapi.NewRail(targetRail.ID, targetRailIndex),
		oapi.NewRail(targetPlayer.Main.ID, consts.RailLimit/2),
		h.playerID,
	)
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

	nowLife := domain.CalculateLife(targetPlayer.LifeEvents)
	targetPlayer.LifeEvents = append(targetPlayer.LifeEvents, domain.NewLifeEvent(
		uuid.New(),
		cardType,
		now,
		domain.LifeEventTypeHealed,
		consts.MaxLife-nowLife,
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

	targetRail, ok := getNonBlockingRail(targetPlayer, true)
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
		targetRail.ID,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, delay, attack)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handlePairExtraordinaire(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	cardType := domain.CardTypePairExtraordinaire

	targetRail, ok := getNonBlockingRail(targetPlayer, true)
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
		targetRail.ID,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, delay, attack)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handleLgtm(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	cardType := domain.CardTypeLgtm

	targetRail, ok := getNonBlockingRail(targetPlayer, true)
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
		targetRail.ID,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, delay, attack)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handlePullShark(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	cardType := domain.CardTypePullShark

	afterRails := domain.NewRails(targetPlayer.Main)
	if l := len(targetPlayer.BranchEvents); l > 0 {
		copy(afterRails[:], targetPlayer.BranchEvents[l-1].AfterRails[:])
	}

	emptys := []int{}

	for i, r := range afterRails {
		if r == nil {
			emptys = append(emptys, i)
		}
	}

	if len(emptys) == 0 {
		targetPlayer.JustCardEvents = append(targetPlayer.JustCardEvents, domain.NewJustCardEvent(
			uuid.New(),
			cardType,
			now,
		))

		return oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now), nil
	}

	railID := uuid.New()
	newRailIndex := emptys[rand.Intn(len(emptys))]
	afterRails[newRailIndex] = domain.NewRail(railID)

	targetPlayer.BranchEvents = append(targetPlayer.BranchEvents, domain.NewBranchEvent(
		uuid.New(),
		cardType,
		now,
		domain.BranchEventCreated,
		h.playerID,
		afterRails,
	))

	// newRailIndexからmainに向かってサーチし、始めに見つかったレールを親とする
	parent := oapi.NewRail(targetPlayer.Main.ID, consts.RailLimit/2)

	if newRailIndex < consts.RailLimit {
		for i := newRailIndex + 1; i < consts.RailLimit/2; i++ {
			if afterRails[i].ID != uuid.Nil {
				parent = oapi.NewRail(afterRails[i].ID, i)
				break
			}
		}
	} else {
		for i := newRailIndex - 1; i >= consts.RailLimit/2; i-- {
			if afterRails[i].ID != uuid.Nil {
				parent = oapi.NewRail(afterRails[i].ID, i)
				break
			}
		}
	}

	res, err := oapi.NewWsResponseRailCreated(
		now,
		oapi.NewRail(railID, newRailIndex),
		parent,
		h.playerID,
		targetPlayer.ID,
		oapi.CardTypePullShark,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handleStarstruck(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	cardType := domain.CardTypeStarstruck

	targetRail, ok := getNonBlockingRail(targetPlayer, true)
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
		targetRail.ID,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, delay, attack)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// getNonBlockingRail 現状のレールからランダムにブロック対象のレールを取得する
// 既にブロックされているブランチは取得できない
func getNonBlockingRail(p *domain.Player, allowMain bool) (*domain.Rail, bool) {
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

	shuffleRails := []*domain.Rail{p.Main}

	if l := len(p.BranchEvents); l > 0 {
		rails := p.BranchEvents[l-1].AfterRails
		shuffleRails = make([]*domain.Rail, len(rails))
		copy(shuffleRails, rails[:])

		rand.Shuffle(len(shuffleRails), func(i, j int) {
			shuffleRails[i], shuffleRails[j] = shuffleRails[j], shuffleRails[i]
		})
	}

	for _, r := range shuffleRails {
		if r == nil {
			continue
		}

		if _, ok := blockBranchIDs[r.ID]; !ok {
			if !allowMain && r.ID == p.Main.ID {
				continue
			}

			return r, true
		}
	}

	return nil, false
}
