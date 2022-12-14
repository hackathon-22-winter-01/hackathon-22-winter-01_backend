package repository

import (
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
)

type RoomRepository interface {
	JoinRoom(jr *JoinRoomArgs) (*domain.Room, uuid.UUID, error)
	CreateRoom(cr *CreateRoomArgs) error
	GetRoom(rid uuid.UUID) (*domain.Room, error)
	DeleteRoom(rid uuid.UUID) error
}

type CreateRoomArgs struct {
	UserId uuid.UUID
}

type JoinRoomArgs struct {
	RoomId   uuid.UUID
	UserName string
}
