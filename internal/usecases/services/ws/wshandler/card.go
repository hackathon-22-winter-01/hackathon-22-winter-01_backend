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
	b, err := body.AsWsRequestBodyCardForAllEvent()
	if err != nil {
		return err
	}

	fmap := map[oapi.CardType]func(reqbody oapi.WsRequestBodyCardForAllEvent, now time.Time) ([]*oapi.WsResponse, error){
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

	usedRails := domain.CalcUsedRails(targetPlayer.BranchEvents)
	parentRailIndex := domain.GetParentRailIndex(targetRail.Index, usedRails)

	targetPlayer.BranchEvents = append(targetPlayer.BranchEvents, domain.NewBranchEvent(
		uuid.New(),
		cardType,
		now,
		domain.BranchEventMerged,
		h.playerID,
		targetRail.Index,
		parentRailIndex,
	))

	res, err := oapi.NewWsResponseRailMerged(
		now,
		targetRail.Index,
		parentRailIndex,
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

	oapiCardType := oapi.CardTypeOpenSourcerer

	res, err := oapi.NewWsResponseLifeChanged(
		now,
		h.playerID,
		&oapiCardType,
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
		targetRail.Index,
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
		targetRail.Index,
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
		targetRail.Index,
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

	usedRails := domain.CalcUsedRails(targetPlayer.BranchEvents)
	blockedRails := domain.CalcBlockedRails(targetPlayer.BlockEvents)

	// usedRails[i] == false && blockedRails[i] == false のとき、レールiは空いている
	emptyRails := make([]int, 0, consts.RailLimit)

	for i := 0; i < consts.RailLimit; i++ {
		if !usedRails[i] && !blockedRails[i] {
			emptyRails = append(emptyRails, i)
		}
	}

	if len(emptyRails) == 0 {
		targetPlayer.JustCardEvents = append(targetPlayer.JustCardEvents, domain.NewJustCardEvent(
			uuid.New(),
			cardType,
			now,
		))

		return oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now), nil
	}

	// 空いているレールをランダムで指定する
	newRailIndex := emptyRails[rand.Intn(len(emptyRails))]

	// newRailIndexからmain方向に向かって一番近い使用中のレールを親として分岐させる
	parentRailIndex := domain.GetParentRailIndex(newRailIndex, usedRails)

	targetPlayer.BranchEvents = append(targetPlayer.BranchEvents, domain.NewBranchEvent(
		uuid.New(),
		cardType,
		now,
		domain.BranchEventCreated,
		h.playerID,
		newRailIndex,
		parentRailIndex,
	))

	res, err := oapi.NewWsResponseRailCreated(
		now,
		newRailIndex,
		parentRailIndex,
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
		targetRail.Index,
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
		targetRail.Index,
		delay,
		attack,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handleOoops(reqbody oapi.WsRequestBodyCardForAllEvent, now time.Time) ([]*oapi.WsResponse, error) {
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

		r, err := oapi.NewWsResponseBlockCreated(
			now,
			h.playerID,
			targetPlayer.ID,
			oapi.CardTypeOoops,
			targetRail.Index,
			delay,
			attack,
		)
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
	usedRails := domain.CalcUsedRails(p.BranchEvents)
	blockedRails := domain.CalcBlockedRails(p.BlockEvents)

	// 存在しているかつブロックされていないレールを対象にする
	targetRails := make([]int, 0)

	for i := 0; i < consts.RailLimit; i++ {
		if !allowMain && i == consts.RailLimit/2 {
			continue
		}

		if usedRails[i] && !blockedRails[i] {
			targetRails = append(targetRails, i)
		}
	}

	if len(targetRails) == 0 {
		return nil, false
	}

	// ランダムにブロック対称のレールを取得する
	targetRailIndex := targetRails[rand.Intn(len(targetRails))]

	return &domain.Rail{Index: targetRailIndex}, true
}
