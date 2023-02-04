package producer

import (
	"github.com/Kiatsyndesi/api_tg_bot/internal/app/sender"
	"github.com/Kiatsyndesi/api_tg_bot/internal/model"
	"github.com/gammazero/workerpool"
	"sync"
	"time"
)

type IProducer interface {
	Start()
	Close()
}

type Producer struct {
	n             uint64
	eventsReading <-chan model.PhonesEvent
	sender        sender.EventSender
	workerPool    *workerpool.WorkerPool
	timeout       time.Duration
	doneChan      chan bool
	wg            *sync.WaitGroup
}

func NewKafkaProducer(
	n uint64,
	events <-chan model.PhonesEvent,
	workerPool *workerpool.WorkerPool,
	sender sender.EventSender,
) Producer {
	wg := &sync.WaitGroup{}
	doneChan := make(chan bool)

	return Producer{
		n:             n,
		eventsReading: events,
		sender:        sender,
		workerPool:    workerPool,
		doneChan:      doneChan,
		wg:            wg,
	}
}

func (p *Producer) Start() {
	for i := uint64(1); i < p.n; i++ {
		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			select {
			case event := <-p.eventsReading:
				if err := p.sender.Send(&event); err != nil {
					p.workerPool.Submit(func() {
						//делаем какую-нибудь логику если пришла ошибка
					})
				} else {
					p.workerPool.Submit(func() {
						//делаем какую-нибудь логику если все хорошо
					})
				}
			case <-p.doneChan:
				return
			}
		}()
	}
}

func (p *Producer) Close() {
	defer p.wg.Wait()
	close(p.doneChan)
}
