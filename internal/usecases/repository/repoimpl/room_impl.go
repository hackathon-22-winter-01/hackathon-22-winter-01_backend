package repository

import (
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
)

const capacity = 4

type roomRepository struct {
	rooms map[uuid.UUID]*domain.Room
}

func NewRepository() repository.Repository {
	return &roomRepository{
		rooms: make(map[uuid.UUID]*domain.Room),
	}
}

func (r *roomRepository) JoinRoom(jr *repository.JoinRoomArgs) (*domain.Room, uuid.UUID, error) {
	room, ok := r.rooms[jr.RoomId]

	if !ok {
		return nil, uuid.Nil, repository.ErrNotFound
	}

	if capacity <= len(room.Players) {
		return nil, uuid.Nil, repository.ErrForbidden
	}

	uid := uuid.New()

	room.Players = append(room.Players, domain.Player{
		ID:   uid,
		Name: jr.UserName,
	})

	return room, uid, nil
}

func (r *roomRepository) CreateRoom(cr *repository.CreateRoomArgs) (*domain.Room, error) {

	rid := uuid.New()
	if _, ok := r.rooms[rid]; ok {
		return nil, repository.ErrAlreadyExists
	}

	uid := uuid.New()

	room := &domain.Room{
		ID: rid,
		Players: []domain.Player{
			{
				ID:   uid,
				Name: cr.UserName,
			},
		},
	}

	r.rooms[rid] = room

	return room, nil
}

func (r *roomRepository) GetRoom(rid uuid.UUID) (*domain.Room, error) {
	panic("implement me")
}

func (r *roomRepository) DeleteRoom(rid uuid.UUID) error {
	panic("implement me")
}
