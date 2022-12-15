package wshandler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
)

type WsHandler interface {
	HandleEvent(req *oapi.WsRequest) error
}

type WsResponseSender interface {
	Bloadcast(roomID uuid.UUID, res *oapi.WsResponse)
}

type wsHandler struct {
	playerID uuid.UUID
	roomRepo repository.RoomRepository
	sender   WsResponseSender
}

func NewWsHandler(playerID uuid.UUID, roomRepo repository.RoomRepository, sender WsResponseSender) WsHandler {
	return &wsHandler{
		playerID: playerID,
		roomRepo: roomRepo,
		sender:   sender,
	}
}

func (h *wsHandler) HandleEvent(req *oapi.WsRequest) error {
	switch req.Type {
	case oapi.WsRequestTypeGameStartEvent:
		return h.handleGameStartEvent(req.Body)

	case oapi.WsRequestTypeCardEvent:
		return h.handleCardEvent(req.Body)

	case oapi.WsRequestTypeLifeEvent:
		return h.handleLifeEvent(req.Body)

	case oapi.WsRequestTypeRailMergeEvent:
		return h.handleRailMergeEvent(req.Body)

	default:
		return errors.New("invalid request type")
	}
}
