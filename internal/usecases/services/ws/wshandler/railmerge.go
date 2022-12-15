package wshandler

import (
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
)

func (h *wsHandler) handleRailMergeEvent(_ oapi.WsRequest_Body) error {
	res, err := oapi.NewWsResponseRailMerged(jst.Now())
	if err != nil {
		return err
	}

	h.sender.Bloadcast(repository.CommonRoomID, res)

	return nil
}
