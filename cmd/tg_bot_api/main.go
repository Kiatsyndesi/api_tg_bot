package main

import (
	"github.com/Kiatsyndesi/api_tg_bot/internal/app/retranslator"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)

	cfg := retranslator.Config{
		ChannelSize:     512,
		ConsumerCount:   2,
		ConsumerTimeout: 1,
		ConsumerSize:    10,
		ProducerCount:   28,
		WorkerCount:     2,
	}

	retranslator := retranslator.NewRetranslator(cfg)
	retranslator.Start()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
