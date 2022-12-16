package wshandler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
)

type WsHandler interface {
	HandleEvent(req *oapi.WsRequest) error
}

type WsResponseSender interface {
	Broadcast(roomID uuid.UUID, res *oapi.WsResponse) error
}

type wsHandler struct {
	playerID uuid.UUID
	room     *domain.Room // roomRepositoryの方がいいかも
	cardRepo repository.CardRepository
	sender   WsResponseSender
}

func NewWsHandler(playerID uuid.UUID, room *domain.Room, cardRepo repository.CardRepository, sender WsResponseSender) WsHandler {
	return &wsHandler{
		playerID: playerID,
		room:     room,
		cardRepo: cardRepo,
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
