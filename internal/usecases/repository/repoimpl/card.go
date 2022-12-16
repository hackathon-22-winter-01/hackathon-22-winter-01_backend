package repoimpl

import "github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"

type cardRepository struct{}

func NewCardRepository() repository.CardRepository {
	return &cardRepository{}
}
