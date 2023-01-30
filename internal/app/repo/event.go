package repo

import "github.com/Kiatsyndesi/api_tg_bot/internal/model"

type EventRepo interface {
	Lock(n uint64) ([]model.PhonesEvent, error)
	Unlock(eventIDs []uint64) error
	Add(event []model.PhonesEvent) error
	Remove(eventIDs []uint64) error
}
