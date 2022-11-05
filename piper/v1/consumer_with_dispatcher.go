package v1

import (
	"context"
	"fmt"
	"piped/piper"
	"sync"
)

type UserBroadcastDispatcher struct {
	Consumers  []UserBroadcastConsumer
	Subscribes piper.SimpleProducer
}

func (bd *UserBroadcastDispatcher) Dispatch(ctx context.Context, opts piper.Opts) chan piper.UserResult {
	res := make(chan piper.UserResult, 1)

	go func() {
		for result := range bd.Subscribes.Next(ctx, opts) {
			if result.Err != nil {
				fmt.Println("subscriber error ", result.Err)
				res <- result
				close(res)
				return
			}

			for _, consumer := range bd.Consumers {
				go func(c UserBroadcastConsumer, result piper.UserResult) {
					res <- piper.UserResult{Data: result.Data, Err: c.Handle(ctx, result).Err}
				}(consumer, result)
			}
		}
	}()

	return res
}

type UserBroadcastSupervisor struct {
	Dispatcher *UserBroadcastDispatcher
}

func (bs *UserBroadcastSupervisor) Run(ctx context.Context, opts piper.Opts) {
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
	Handle(ctx context.Context, result piper.UserResult) piper.UserResult
}

type UserIDBroadcastConsumer struct {
	// In chan UserResult
}

func (ubc *UserIDBroadcastConsumer) Handle(ctx context.Context, result piper.UserResult) piper.UserResult {
	// var err = make(chan error, 1)

	// go func(res UserResult) {
	fmt.Println("id broadcaster ", result.Data.ID)
	// }(result)

	return piper.UserResult{Data: result.Data, Err: nil}
}

type UserNameBroadcastConsumer struct {
	// In chan UserResult
}

func (ubc *UserNameBroadcastConsumer) Handle(ctx context.Context, result piper.UserResult) piper.UserResult {
	// var err = make(chan error, 1)

	// go func(res UserResult) {
	fmt.Println("name broadcaster ", result.Data.Name)
	// }(result)

	return piper.UserResult{Data: result.Data, Err: nil}

}
