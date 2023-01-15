package gompscqueue

import "sync"

func (o *MPSCQueue) WaitGroup() *sync.WaitGroup {
	return o.producers
}

func (o *MPSCQueue) MessagesCh() chan []interface{} {
	return o.msgCh
}

func (o *MPSCQueue) ErrCh() chan error {
	return o.errCh
}

func (o *MPSCQueue) DoneCh() chan bool {
	return o.consumerDoneCh
}
