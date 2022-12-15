package repository

import (
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
)

// TODO: 消す
var CommonRoomID = uuid.New()

type RoomRepository interface {
	FindRoom(roomID uuid.UUID) (*domain.Room, error)
	JoinRoom(roomID uuid.UUID, player *domain.Player) error
}
