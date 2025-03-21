package eventloop

import (
	"net"
	"sync"
)

type Task struct {
	MainTask   func(c net.Conn)
	IsBlocking bool
}

type EventLoop struct {
	Tasks chan Task
	Stop  chan bool
}

func AddToEventLoop(eventLoop *EventLoop, t Task, conn net.Conn) {
	originalHandler := t.MainTask

	t.MainTask = func(_ net.Conn) {
		originalHandler(conn)
	}

	eventLoop.Tasks <- t
}

func StopEventLoop(eventLoop *EventLoop) {
	eventLoop.Stop <- true
}

func InitEventLoop(eventLoop *EventLoop, workerPoolSize int) *sync.WaitGroup {
	var wg sync.WaitGroup

	// Add the event loop goroutine to the wait group to track it's execution
	wg.Add(1)

	// Worker pool for handling blocking tasks
	workerPool := make(chan struct{}, workerPoolSize)

	go func() {
		defer wg.Done()

		for {
			select {
			case task := <-eventLoop.Tasks:
				if task.IsBlocking {
					workerPool <- struct{}{}
					go func() {
						defer func() { <-workerPool }()
						task.MainTask(nil)
					}()
				} else {
					task.MainTask(nil)
				}

			case <-eventLoop.Stop:
				return
			}
		}
	}()

	return &wg
}
