package repoimpl

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/config"
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

func (r *roomRepository) FindRoomFromPlayerID(playerID uuid.UUID) (*domain.Room, error) {
	var room *domain.Room

	r.roomMap.Range(func(_ uuid.UUID, value *domain.Room) bool {
		for _, player := range value.Players.Clone() {
			if player.ID == playerID {
				room = value
				return false
			}
		}
		return true
	})

	if room == nil {
		return nil, errors.New("指定したプレイヤーが属する部屋は存在しません")
	}

	return room, nil
}

func (r *roomRepository) JoinRoom(roomID uuid.UUID, player *domain.Player) error {
	room, ok := r.roomMap.Load(roomID)
	if !ok {
		return errors.New("部屋が存在しません")
	}

	if len(room.Players.Clone()) >= config.PlayerLimit {
		return errors.New("部屋が満員です")
	}

	room.Players.Append(player)

	return nil
}

func (r *roomRepository) CreateRoom(player *domain.Player) (*domain.Room, error) {
	roomID := uuid.New()

	room := domain.NewRoom(roomID)

	room.Players.Append(player)

	_, ok := r.roomMap.LoadOrStore(roomID, room)

	if ok {
		return nil, errors.New("部屋が既に存在します")
	}

	return room, nil
}

func (r *roomRepository) DeleteRoom(roomID uuid.UUID) error {
	_, ok := r.roomMap.Load(roomID)
	if !ok {
		return errors.New("部屋が存在しません")
	}

	r.roomMap.Delete(roomID)

	return nil
}
