package ws

import (
	"errors"

	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
)

func (h *Hub) handleEvent(req *oapi.WsRequest) error {
	switch req.Type {
	case oapi.WsRequestTypeGameStartEvent:
		_, err := req.Body.AsWsRequestBodyGameStartEvent()
		if err != nil {
			return err
		}

		res, err := h.sendGameStarted()
		if err != nil {
			return err
		}

		h.bloadcast(res)

		return nil

	case oapi.WsRequestTypeCardEvent:
		b, err := req.Body.AsWsRequestBodyCardEvent()
		if err != nil {
			return err
		}

		switch b.Type {
		case oapi.CardTypeCreateRail:
			res, err := h.sendRailCreated()
			if err != nil {
				return err
			}

			h.bloadcast(res)

		case oapi.CardTypeCreateBlock:
			res, err := h.sendBlockCreated()
			if err != nil {
				return err
			}

			h.bloadcast(res)

		default:
			return errors.New("invalid card type")
		}

		return nil

	case oapi.WsRequestTypeLifeEvent:
		b, err := req.Body.AsWsRequestBodyLifeEvent()
		if err != nil {
			return err
		}

		switch b.Type {
		case oapi.LifeEventTypeDecrement:
			res, err := h.sendLifeChanged()
			if err != nil {
				return err
			}

			h.bloadcast(res)

		default:
			return errors.New("invalid life type")
		}

		return nil

	case oapi.WsRequestTypeRailMergeEvent:
		res, err := h.sendRailMerged()
		if err != nil {
			return err
		}

		h.bloadcast(res)

		return nil

	default:
		return errors.New("invalid request type")
	}
}
