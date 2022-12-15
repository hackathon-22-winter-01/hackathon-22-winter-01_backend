package repoimpl

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/sync"
)

type roomRepository struct {
	roomMap sync.Map[uuid.UUID, *domain.Room]
}

func NewRoomRepository() repository.RoomRepository {
	return &roomRepository{
		roomMap: sync.Map[uuid.UUID, *domain.Room]{},
	}
}

func (r *roomRepository) FindRoom(roomID uuid.UUID) (*domain.Room, error) {
	room, ok := r.roomMap.Load(roomID)
	if !ok {
		return nil, errors.New("部屋が存在しません")
	}

	return room, nil
}

func (r *roomRepository) JoinRoom(roomID uuid.UUID, player *domain.Player) error {
	// TODO: ここでroomMapにroomIDが存在しない場合はエラーを返す
	room, ok := r.roomMap.LoadOrStore(roomID, domain.NewRoom(roomID))
	if !ok {
		log.Println("部屋が存在しないので作成しました。本来はエラーを返すべきです。")
	}

	room.Players = append(room.Players, player)

	return nil
}
