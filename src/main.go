package main

import "github.com/FrMnJ/postago/src/queue"

func main() {
	mailQueue, err := queue.NewMailQueue()
	if err != nil {
		panic(err)
	}
	mailQueue.MainLoop()
}
