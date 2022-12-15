package repository

import (
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
)

type PlayerRepository interface {
	FindUser(uid uuid.UUID) (*domain.Player, error)
	CreateUser(user *domain.Player) error
}
