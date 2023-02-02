package consumer

import (
	"github.com/Kiatsyndesi/api_tg_bot/internal/app/repo"
	"github.com/Kiatsyndesi/api_tg_bot/internal/model"
	"sync"
	"time"
)

/*Первоочередная цель консьюмера - получить данные из БД
У консьюмера есть 2 метода - открыть и закрыть соединение самого себя с БД

- Струкура консьюмера:
1. Кол-во консьюмеров
2. Канал с событиями по телефонам
3. Интерфейс для работы ивентами (repo)
4. Размер батча данных и таймаут с типом Duration(базовый тип пакета duration)
5. Канал с булевыми значениями для передачи done
6. вейтгруппа для синхронизации потоков

- Структура конфига:
1. Имеет все те же самые поля, что и консьюмер, только без инструментов для управления конкурентностью

- Функция инициализации консьюмера

- Метод для старта работы консьюмера:
1. запускаем цикл, кол-во итераций зависит от количества консьюмеров
2. в горутине создаем тикер
3. запускаем бесконечный цикл с селектом
4. если приходит из тикера - лочим ивенты, если ошибка - пропускаем, в другом случае пишем в канал с ивентами
5. В случае если из булева канала получаем что-либо, то делаем return из функции

- Метода закрытия:
1. Ждем выполнения всех горутин
2. закрываем канал с булевыми
*/

type ICustomer interface {
	Start()
	Close()
}

type Customer struct {
	n         uint64
	events    chan<- model.PhonesEvent
	repo      repo.EventRepo
	butchSize uint64
	timeout   time.Duration
	doneChan  chan bool
	wg        *sync.WaitGroup
}

type Config struct {
	n         uint64
	events    chan<- model.PhonesEvent
	repo      repo.EventRepo
	butchSize uint64
	timeout   time.Duration
}

func NewCustomer(
	numberOfCustomers uint64,
	eventsWriting chan<- model.PhonesEvent,
	repo repo.EventRepo,
	butchSize uint64,
	timeout time.Duration) Customer {

	wg := &sync.WaitGroup{}
	doneChan := make(chan bool)
	return Customer{
		n:         numberOfCustomers,
		events:    eventsWriting,
		repo:      repo,
		butchSize: butchSize,
		timeout:   timeout,
		doneChan:  doneChan,
		wg:        wg,
	}
}

func (c *Customer) Start() {
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)

			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(c.butchSize)

					if err != nil {
						continue
					}

					for _, event := range events {
						c.events <- event
					}
				case <-c.doneChan:
					return
				}
			}
		}()
	}
}

func (c *Customer) Close() {
	c.wg.Wait()
	close(c.doneChan)
}
