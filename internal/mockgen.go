package internal

//go:generate mockgen -destination=./mocks/repo_mock.go -package=mocks github.com/Kiatsyndesi/api_tg_bot/internal/app/repo EventRepo
//go:generate mockgen -destination=./mocks/sender_mock.go -package=mocks github.com/Kiatsyndesi/api_tg_bot/internal/app/sender EventSender
