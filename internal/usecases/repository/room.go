package repository

import (
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
)

var CommonRoomID = uuid.New()

type RoomRepository interface {
	JoinRoom(jr *JoinRoomArgs) (*domain.Room, uuid.UUID, error)
	CreateRoom(cr *CreateRoomArgs) (*domain.Room, error)
	GetRoom(rid uuid.UUID) (*domain.Room, error)
	DeleteRoom(rid uuid.UUID) error
	FindRoom(roomID uuid.UUID) (*domain.Room, error)
}

type CreateRoomArgs struct {
	PlayerName string
}

type JoinRoomArgs struct {
	RoomId     uuid.UUID
	PlayerName string
}
