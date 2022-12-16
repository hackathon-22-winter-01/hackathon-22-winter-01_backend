package repository

import (
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
)

type CardRepository interface {
	DrawCards(roomID uuid.UUID, playerID uuid.UUID, num int) ([]*domain.Card, error)
}
