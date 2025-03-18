package main

type Task struct {
	MainTask   func()
	CallBack   func()
	IsBlocking bool
}
