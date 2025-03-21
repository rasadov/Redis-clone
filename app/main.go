package main

import (
	"fmt"
	"github.com/rasadov/redis-clone/app/eventloop"
	"github.com/rasadov/redis-clone/app/models"
	"github.com/rasadov/redis-clone/app/rdb"
	"log"
	"net"
	"os"
)

var (
	Storage     *models.InMemoryStorage
	eventLoop   *eventloop.EventLoop
	rdbFilepath string
)

func init() {
	var err error

	InitConfig(os.Args)
	Storage = models.NewInMemoryStorage(0)
	rdbFilepath = Config["dir"] + "/" + Config["dbfilename"]
	if rdbFilepath != "" {
		err = rdb.LoadRDB(rdbFilepath, Storage)
	}
	if rdbFilepath == "" || err != nil {
		fmt.Println("rdb file not found. Initializing empty database")
	}

	eventLoop = &eventloop.EventLoop{
		Tasks: make(chan eventloop.Task, 100),
		Stop:  make(chan bool),
	}
}

func main() {
	wg := eventloop.InitEventLoop(eventLoop, 5)

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatal("Failed to bind to port 6379:", err)
	}
	defer listener.Close()
	defer func() {
		if rdbFilepath != "" {
			err := rdb.SaveRDB(rdbFilepath, Storage)
			if err != nil {
				fmt.Println("Error while saving rdb file:", err)
			}
		}
	}()

	fmt.Println("Listening on port 6379...")

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}

			// Add connection handling to the event loop
			eventloop.AddToEventLoop(eventLoop, eventloop.Task{
				MainTask:   handleConnection,
				IsBlocking: true,
			}, conn)
		}
	}()

	wg.Wait()
}
