# MPSC queue in Go

This is a simple implementation of a multiple producers, single consumer queue in Go.

It leverages `sync.WaitGroup` and `channel`s, and provides a queue for messages and a queue for errors.
Producers are synchronized through `sync.WaitGroup`.
Synchronization with the consumer happens through a dedicated `channel`.

## Usage

```go
package main

import (
	"fmt"
	
	gompscqueue "github.com/maxgio92/go-mpsc-queue"
)

func main() {
	parallelism := 10

	queue := gompscqueue.NewMPSCQueue(parallelism)

	for i := 0; i < parallelism; i++ {
		go func() {
			defer queue.SigProducerCompletion()

			// Here you do your work.

			queue.SendMessage("Hello world from producer!")
		}()
	}

	go queue.Consume(
		func(msg interface{}) {
			fmt.Printf("new message: %s\n", msg)
		},
		func(err error) {
			fmt.Errorf("error: %s\n", err.Error())
		},
	)

	queue.WaitAndClose()
}
```