package sender

import "github.com/Kiatsyndesi/api_tg_bot/internal/model"

type EventSender interface {
	Send(phoneEvent *model.PhonesEvent) error
}
