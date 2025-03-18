package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var (
	Storage    *InMemoryStorage
	ConfigFile = "conf.json"
)

func main() {
	InitConfig(os.Args)

	Storage = NewInMemoryStorage()
	//err := ReadConfig(ConfigFile, Storage)

	//if err != nil {
	//	log.Fatal("Issue with reading the configuration file")
	//}

	eventLoop := &EventLoop{
		mainTasks: make(chan Task, 100),
		stop:      make(chan bool),
	}

	wg := InitEventLoop(eventLoop, 5)

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatal("Failed to bind to port 6379:", err)
	}
	defer listener.Close()

	fmt.Println("Listening on port 6379...")

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}

			// Add connection handling to the event loop
			AddToEventLoop(eventLoop, Task{
				MainTask:   handleConnection,
				IsBlocking: true,
			}, conn)
		}
	}()

	//err = WriteConfig(ConfigFile, Storage.data)
	//if err != nil {
	//	fmt.Println("Error writing to configuration file")
	//}
	wg.Wait()
}
