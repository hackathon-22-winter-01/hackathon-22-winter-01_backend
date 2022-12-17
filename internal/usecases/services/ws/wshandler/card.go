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
		oapi.CardTypeZeroDay:            h.handleZeroDay,
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

func (h *wsHandler) handleCardForAllEvent(body oapi.WsRequest_Body) error {
	b, err := body.AsWsRequestBodycardForAllEvent()
	if err != nil {
		return err
	}

	fmap := map[oapi.CardType]func(reqbody oapi.WsRequestBodycardForAllEvent, now time.Time) ([]*oapi.WsResponse, error){
		oapi.CardTypeOoops: h.handleOoops,
	}

	f, ok := fmap[b.Type]
	if !ok {
		return errors.New("存在しないカードです")
	}

	res, err := f(b, jst.Now())
	if err != nil {
		return err
	}

	for _, r := range res {
		if err := h.sender.Broadcast(h.room.ID, r); err != nil {
			return err
		}
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
		if r.Index != targetRail.Index {
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
		targetRailIndex,
		consts.RailLimit/2,
		h.playerID,
		oapi.CardTypeYolo,
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

	res, err := oapi.NewWsResponseLifeChanged(
		now,
		h.playerID,
		oapi.CardTypeOpenSourcerer,
		domain.CalculateLife(targetPlayer.LifeEvents),
	)
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
		targetRail.Index,
	))

	res, err := oapi.NewWsResponseBlockCreated(
		now,
		h.playerID,
		reqbody.TargetId,
		oapi.CardTypeRefactoring,
		delay,
		attack,
	)
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
		targetRail.Index,
	))

	res, err := oapi.NewWsResponseBlockCreated(
		now,
		h.playerID,
		reqbody.TargetId,
		oapi.CardTypePairExtraordinaire,
		delay,
		attack,
	)
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
		targetRail.Index,
	))

	res, err := oapi.NewWsResponseBlockCreated(
		now,
		h.playerID,
		reqbody.TargetId,
		oapi.CardTypeLgtm,
		delay,
		attack,
	)
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

	// railID := uuid.New()
	newRailIndex := emptys[rand.Intn(len(emptys))]
	afterRails[newRailIndex] = domain.NewRail(newRailIndex)

	targetPlayer.BranchEvents = append(targetPlayer.BranchEvents, domain.NewBranchEvent(
		uuid.New(),
		cardType,
		now,
		domain.BranchEventCreated,
		h.playerID,
		afterRails,
	))

	// newRailIndexからmainに向かってサーチし、始めに見つかったレールを親とする
	parent := consts.RailLimit / 2

	if newRailIndex < consts.RailLimit {
		for i := newRailIndex + 1; i < consts.RailLimit/2; i++ {
			if afterRails[i] != nil {
				parent = i
				break
			}
		}
	} else {
		for i := newRailIndex - 1; i >= consts.RailLimit/2; i-- {
			if afterRails[i] != nil {
				parent = i
				break
			}
		}
	}

	res, err := oapi.NewWsResponseRailCreated(
		now,
		newRailIndex,
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
		targetRail.Index,
	))

	res, err := oapi.NewWsResponseBlockCreated(
		now,
		h.playerID,
		reqbody.TargetId,
		oapi.CardTypeStarstruck,
		delay,
		attack,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handleZeroDay(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	cardType := domain.CardTypeZeroDay

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
		targetRail.Index,
	))

	res, err := oapi.NewWsResponseBlockCreated(
		now,
		h.playerID,
		reqbody.TargetId,
		oapi.CardTypeZeroDay,
		delay,
		attack,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handleOoops(reqbody oapi.WsRequestBodycardForAllEvent, now time.Time) ([]*oapi.WsResponse, error) {
	cardType := domain.CardTypeOoops

	var res []*oapi.WsResponse

	for _, targetPlayer := range h.room.Players {
		targetRail, ok := getNonBlockingRail(targetPlayer, true)
		if !ok {
			targetPlayer.JustCardEvents = append(targetPlayer.JustCardEvents, domain.NewJustCardEvent(
				uuid.New(),
				cardType,
				now,
			))

			res = append(res, oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now))

			continue
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
			targetPlayer.ID,
			targetRail.Index,
		))

		r, err := oapi.NewWsResponseBlockCreated(now, h.playerID, targetPlayer.ID, oapi.CardTypeOoops, delay, attack)
		if err != nil {
			return nil, err
		}

		res = append(res, r)
	}

	return res, nil
}

// getNonBlockingRail 現状のレールからランダムにブロック対象のレールを取得する
// 既にブロックされているブランチは取得できない
func getNonBlockingRail(p *domain.Player, allowMain bool) (*domain.Rail, bool) {
	// 既にブロックされているブランチのIDを取得
	blockBranchIDs := make(map[int]struct{})

	for _, e := range p.BlockEvents {
		switch e.Type {
		case domain.BlockEventTypeCreated:
			blockBranchIDs[e.TargetRailIndex] = struct{}{}
		case domain.BlockEventTypeCanceled:
			delete(blockBranchIDs, e.TargetRailIndex)
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

		if _, ok := blockBranchIDs[r.Index]; !ok {
			if !allowMain && r.Index == p.Main.Index {
				continue
			}

			return r, true
		}
	}

	return nil, false
}
