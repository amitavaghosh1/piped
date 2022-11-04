package piper

import (
	"context"
	"fmt"
	"sync"
)

// SimpleProducer
// ProducerConsumer
// Consumer

type UserProducerConsumerSupervisor struct {
	Subscribes SimpleProducer
	Consumer   UserBroadcastConsumer
}

func (ucs *UserProducerConsumerSupervisor) Run(ctx context.Context, opts Opts) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for {
		}
	}()

	wg.Wait()
}

func (ucs *UserProducerConsumerSupervisor) loop(ctx context.Context, opts Opts) error {
	for result := range ucs.Subscribes.Next(ctx, opts) {
		if result.Err != nil {
			fmt.Println("error ", result.Err)
			return result.Err
		}

		if err  := ucs.Consumer.Handle(ctx, result); err != nil {
			
		}
	}
}
