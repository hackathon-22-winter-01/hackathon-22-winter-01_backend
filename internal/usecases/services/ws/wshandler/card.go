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

	targetRailID, ok := getNonBlockingRailID(targetPlayer)
	if !ok {
		return oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now), nil
	}

	cardType := domain.CardTypeRefactoring

	delay, attack, err := cardType.DelayAndAttack()
	if err != nil {
		return nil, err
	}

	targetPlayer.BlockEvents = append(targetPlayer.BlockEvents, domain.NewBlockEvent(
		uuid.New(),
		cardType,
		now,
		h.playerID,
		targetPlayer.ID,
		targetRailID,
		delay,
		attack,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, delay, attack)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handlePairExtraordinaire(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	targetRailID, ok := getNonBlockingRailID(targetPlayer)
	if !ok {
		return oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now), nil
	}

	cardType := domain.CardTypePairExtraordinaire

	delay, attack, err := cardType.DelayAndAttack()
	if err != nil {
		return nil, err
	}

	targetPlayer.BlockEvents = append(targetPlayer.BlockEvents, domain.NewBlockEvent(
		uuid.New(),
		domain.CardTypePairExtraordinaire,
		now,
		h.playerID,
		targetPlayer.ID,
		targetRailID,
		delay,
		attack,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, delay, attack)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *wsHandler) handleLgtm(reqbody oapi.WsRequestBodyCardEvent, now time.Time, targetPlayer *domain.Player) (*oapi.WsResponse, error) {
	targetRailID, ok := getNonBlockingRailID(targetPlayer)
	if !ok {
		return oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now), nil
	}

	cardType := domain.CardTypeLgtm

	delay, attack, err := cardType.DelayAndAttack()
	if err != nil {
		return nil, err
	}

	targetPlayer.BlockEvents = append(targetPlayer.BlockEvents, domain.NewBlockEvent(
		uuid.New(),
		cardType,
		now,
		h.playerID,
		targetPlayer.ID,
		targetRailID,
		delay,
		attack,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, delay, attack)
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
	targetRailID, ok := getNonBlockingRailID(targetPlayer)
	if !ok {
		return oapi.WsResponseFromType(oapi.WsResponseTypeNoop, now), nil
	}

	cardType := domain.CardTypeStarstruck

	delay, attack, err := cardType.DelayAndAttack()
	if err != nil {
		return nil, err
	}

	targetPlayer.BlockEvents = append(targetPlayer.BlockEvents, domain.NewBlockEvent(
		uuid.New(),
		domain.CardTypeLgtm,
		now,
		h.playerID,
		targetPlayer.ID,
		targetRailID,
		delay,
		attack,
	))

	res, err := oapi.NewWsResponseBlockCreated(now, h.playerID, reqbody.TargetId, delay, attack)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// getNonBlockingRailID 現状のレールからランダムにブロック対象のレールを取得する
// 既にブロックされているブランチは取得できない
func getNonBlockingRailID(p *domain.Player) (uuid.UUID, bool) {
	// 既にブロックされているブランチのIDを取得
	blockBranchIDs := make(map[uuid.UUID]struct{})
	for _, e := range p.BlockEvents {
		// TODO: 解消時はmapから消す実装を書く
		blockBranchIDs[e.TargetRailID] = struct{}{}
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
			return id, true
		}
	}

	return uuid.Nil, false
}
