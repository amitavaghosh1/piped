package piper

import (
	"context"
	"fmt"
	"sync"
)

type UserBroadcastDispatcher struct {
	Consumers  []UserBroadcastConsumer
	Subscribes SimpleProducer
}

func (bd *UserBroadcastDispatcher) Dispatch(ctx context.Context, opts Opts) chan UserResult {
	res := make(chan UserResult, 1)

	go func() {
		for result := range bd.Subscribes.Next(ctx, opts) {
			if result.Err != nil {
				fmt.Println("subscriber error ", result.Err)
				res <- result
				close(res)
				return
			}

			for _, consumer := range bd.Consumers {
				go func(c UserBroadcastConsumer, result UserResult) {
					res <- UserResult{Data: result.Data, Err: c.Handle(ctx, result)}
				}(consumer, result)
			}
		}
	}()

	return res
}

type UserBroadcastSupervisor struct {
	Dispatcher *UserBroadcastDispatcher
}

func (bs *UserBroadcastSupervisor) Run(ctx context.Context, opts Opts) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			for result := range bs.Dispatcher.Dispatch(ctx, opts) {
				if result.Err != nil {
					fmt.Println("error ", result.Err)
					return
				}
			}
		}
	}()

	wg.Wait()
}

type UserBroadcastConsumer interface {
	Handle(ctx context.Context, result UserResult) error
}

type UserIDBroadcastConsumer struct {
	// In chan UserResult
}

func (ubc *UserIDBroadcastConsumer) Handle(ctx context.Context, result UserResult) (err error) {
	// var err = make(chan error, 1)

	// go func(res UserResult) {
	fmt.Println("id broadcaster ", result.Data.ID)
	// }(result)

	return err
}

type UserNameBroadcastConsumer struct {
	// In chan UserResult
}

func (ubc *UserNameBroadcastConsumer) Handle(ctx context.Context, result UserResult) (err error) {
	// var err = make(chan error, 1)

	// go func(res UserResult) {
	fmt.Println("name broadcaster ", result.Data.Name)
	// }(result)

	return err

}
