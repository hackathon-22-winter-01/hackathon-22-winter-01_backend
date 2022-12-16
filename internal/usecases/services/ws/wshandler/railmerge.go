package wshandler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
)

// TODO: 呼ばれないので消す
func (h *wsHandler) handleRailMergeEvent(body oapi.WsRequest_Body) error {
	b, err := body.AsWsRequestBodyRailMergeEvent()
	if err != nil {
		return err
	}

	res, err := oapi.NewWsResponseRailMerged(jst.Now(), b.ChildId, b.ParentId, h.playerID)
	if err != nil {
		return err
	}

	if err := h.sender.Broadcast(h.room.ID, res); err != nil {
		return err
	}

	target, ok := h.room.FindPlayer(h.playerID)
	if !ok {
		return errors.New("player not found")
	}

	l := len(target.RailEvents)
	if l == 0 {
		return errors.New("メインのレールを削除することはできません")
	}

	beforeRails := target.RailEvents[l-1].AfterRails
	if len(beforeRails) == 0 {
		return errors.New("メインのレールを削除することはできません")
	}

	afterRails := make([]*domain.Rail, 0)
	for _, r := range afterRails {
		if r.ID != b.ChildId {
			afterRails = append(afterRails, r)
		}
	}

	target.RailEvents = append(target.RailEvents, domain.NewRailEvent(
		uuid.New(),
		domain.CardTypeNone,
		domain.RailEventMerged,
		uuid.Nil,
		h.playerID,
		afterRails,
	))

	return nil
}
