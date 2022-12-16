package wshandler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
)

func (h *wsHandler) handleRailMergeEvent(body oapi.WsRequest_Body) error {
	b, err := body.AsWsRequestBodyRailMergeEvent()
	if err != nil {
		return err
	}

	res, err := oapi.NewWsResponseRailMerged(jst.Now(), b.ChildId, b.ParentId, h.playerID)
	if err != nil {
		return err
	}

	h.sender.Broadcast(h.room.ID, res)

	target, ok := h.room.FindPlayer(h.playerID)
	if !ok {
		return errors.New("player not found")
	}

	l := len(target.Events)
	if l == 0 {
		return errors.New("メインのレールを削除することはできません")
	}

	beforeRails := target.Events[l-1].AfterRails
	if len(beforeRails) == 0 {
		return errors.New("メインのレールを削除することはできません")
	}

	afterRails := make([]*domain.Rail, 0)
	for _, r := range afterRails {
		if r.ID != b.ChildId {
			afterRails = append(afterRails, r)
		}
	}

	target.Events = append(target.Events, domain.NewRailEvent(
		uuid.New(),
		domain.RailEventMerged,
		uuid.Nil,
		h.playerID,
		beforeRails,
		afterRails,
	))

	return nil
}
