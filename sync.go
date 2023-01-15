package gompscqueue

import (
	"sync"
)

// MPSCQueue provides an option set to manage a minimal message system with multiple parallel producers and
// single consumer, leveraging Go sync.WaitGroup and channels to notify errors, results and completion
// of consuming the results from the single consumer worker.
type MPSCQueue struct {
	producers      *sync.WaitGroup
	consumerDoneCh chan bool
	msgCh          chan []interface{}
	errCh          chan error
}

func NewMPSCQueue(parallelism int) *MPSCQueue {
	wg := &sync.WaitGroup{}
	wg.Add(parallelism)

	msgCh := make(chan []interface{})

	errCh := make(chan error)

	doneCh := make(chan bool, 1)

	return &MPSCQueue{
		producers:      wg,
		consumerDoneCh: doneCh,
		msgCh:          msgCh,
		errCh:          errCh,
	}
}

// SendMessage sends a message as variadic parameter msg of type interface{} to the messages queue.
func (o *MPSCQueue) SendMessage(msg ...interface{}) {
	o.msgCh <- msg
}

// SendMessageAndComplete sends a message as variadic parameter msg of type interface{} to the messages queue,
// and eventually signals the completion of the current producer.
func (o *MPSCQueue) SendMessageAndComplete(msg ...interface{}) {
	defer o.SigProducerCompletion()
	o.msgCh <- msg
}

// SendError sends an error message of type error to the errors queue.
func (o *MPSCQueue) SendError(err error) {
	o.errCh <- err
}

// Consume listens for both messages and errors on queues and do something with them,
// as specified by msgHandler and errHandler functions.
func (o *MPSCQueue) Consume(msgHandler func(msg interface{}), errHandler func(err error)) {
	for o.errCh != nil || o.msgCh != nil {
		select {
		case p, ok := <-o.msgCh:

			// If the channel is still open.
			if ok {

				// Do something with the message.
				msgHandler(p)
				continue
			}
			o.msgCh = nil
		case e, ok := <-o.errCh:

			// If the channel is still open.
			if ok {

				// Do something with error.
				errHandler(e)
				continue
			}
			o.errCh = nil
		}
	}
	o.SigConsumerCompletion()
}

// SigProducerCompletion signals that a producer completed its work.
func (o *MPSCQueue) SigProducerCompletion() {
	o.producers.Done()
}

// SigConsumerCompletion signals that the consumer completed its work.
func (o *MPSCQueue) SigConsumerCompletion() {
	o.consumerDoneCh <- true
}

func (o *MPSCQueue) WaitAndClose() {

	// Wait for producers to complete.
	o.producers.Wait()
	close(o.msgCh)
	close(o.errCh)

	// Wait for consumers to complete.
	<-o.consumerDoneCh
}
