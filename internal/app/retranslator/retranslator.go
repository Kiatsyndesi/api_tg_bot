package retranslator

import (
	"github.com/Kiatsyndesi/api_tg_bot/internal/app/consumer"
	"github.com/Kiatsyndesi/api_tg_bot/internal/app/producer"
	"github.com/Kiatsyndesi/api_tg_bot/internal/app/repo"
	"github.com/Kiatsyndesi/api_tg_bot/internal/app/sender"
	"github.com/Kiatsyndesi/api_tg_bot/internal/model"
	"github.com/gammazero/workerpool"
	"time"
)

type IRetranslator interface {
	Start()
	Close()
}

type Config struct {
	ChannelSize uint64

	ConsumerCount   uint64
	ConsumerTimeout time.Duration
	ConsumerSize    uint64

	ProducerCount uint64
	WorkerCount   int

	Repo   repo.EventRepo
	Sender sender.EventSender
}

type Retranslator struct {
	events     chan model.PhonesEvent
	consumer   consumer.Consumer
	producer   producer.Producer
	workerPool *workerpool.WorkerPool
}

func NewRetranslator(cfg Config) *Retranslator {
	events := make(chan model.PhonesEvent, cfg.ChannelSize)
	workerPool := workerpool.New(cfg.WorkerCount)

	consumerForRetrans := consumer.NewConsumer(cfg.ConsumerCount, events, cfg.Repo, cfg.ConsumerSize, cfg.ConsumerTimeout)
	producerForRetrans := producer.NewKafkaProducer(cfg.ProducerCount, events, workerPool, cfg.Sender)

	return &Retranslator{
		events:     events,
		consumer:   consumerForRetrans,
		producer:   producerForRetrans,
		workerPool: workerPool,
	}
}

func (r *Retranslator) Start() {
	r.producer.Start()
	r.consumer.Start()
}

func (r *Retranslator) Close() {
	r.consumer.Close()
	r.producer.Close()
	r.workerPool.StopWait()
}
