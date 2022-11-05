package v1

import (
	"context"
	"fmt"
	"piped/loadbalancers"
	"piped/piper"
	"sync"
)

type UserDispatcherV2 interface {
	Dispatch(ctx context.Context, result piper.UserResult) chan piper.UserResult
}

type UserBroadcastDispatcherV2 struct {
	Consumers []UserBroadcastConsumer
}

func (bd *UserBroadcastDispatcherV2) Dispatch(ctx context.Context, result piper.UserResult) chan piper.UserResult {
	res := make(chan piper.UserResult, 1)

	go func() {
		defer close(res)

		if result.Err != nil {
			fmt.Println("subscriber error ", result.Err)
			res <- result
			return
		}

		var wg sync.WaitGroup
		wg.Add(len(bd.Consumers))

		for i, consumer := range bd.Consumers {
			fmt.Println("invoking consumer ", i)
			go bd.work(ctx, &wg, consumer, result, res)
		}

		wg.Wait()
	}()

	return res
}

func (bd *UserBroadcastDispatcherV2) work(ctx context.Context, w *sync.WaitGroup, consumer UserBroadcastConsumer, result piper.UserResult, res chan piper.UserResult) {
	defer w.Done()

	res <- piper.UserResult{
		Data: result.Data,
		Err:  consumer.Handle(ctx, result).Err,
	}
}

type UserConsumerGroupDispatcher struct {
	Consumers []UserBroadcastConsumer
	Balancer  loadbalancers.Balancer
}

func (cd *UserConsumerGroupDispatcher) Dispatch(ctx context.Context, result piper.UserResult) chan piper.UserResult {
	res := make(chan piper.UserResult, 1)

	go func() {
		defer close(res)

		if result.Err != nil {
			fmt.Println("subscriber error ", result.Err)
			res <- result
			return
		}

		idx := cd.Balancer.Get()
		fmt.Println("using worker ", idx)

		consumer := cd.Consumers[idx]
		res <- piper.UserResult{Data: result.Data, Err: consumer.Handle(ctx, result).Err}
	}()

	return res

}

type UserBroadcastSupervisorV2 struct {
	Dispatcher UserDispatcherV2
	Subscribes piper.SimpleProducer
}

func (bs *UserBroadcastSupervisorV2) Run(ctx context.Context, opts piper.Opts) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			if err := bs.loop(ctx, opts); err != nil {
				return
			}
		}
	}()

	wg.Wait()
}

func (bs *UserBroadcastSupervisorV2) loop(ctx context.Context, opts piper.Opts) error {
	for result := range bs.Subscribes.Next(ctx, opts) {
		if result.Err != nil {
			fmt.Println("error ", result.Err)
			return result.Err
		}

		for res := range bs.Dispatcher.Dispatch(ctx, result) {
			if res.Err != nil {
				fmt.Println("consumer error ", res.Err)
				continue
			}
		}
	}

	return nil
}

type UserConsumerGroupSupervisor struct {
	Dispatcher UserDispatcherV2
	Subscribes piper.SimpleProducer
}

func (cs *UserConsumerGroupSupervisor) Run(ctx context.Context, opts piper.Opts) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			if err := cs.loop(ctx, opts); err != nil {
				return
			}
		}
	}()

	wg.Wait()
}

func (cs *UserConsumerGroupSupervisor) loop(ctx context.Context, opts piper.Opts) error {
	for result := range cs.Subscribes.Next(ctx, opts) {
		if result.Err != nil {
			fmt.Println("error ", result.Err)
			return result.Err
		}

		for res := range cs.Dispatcher.Dispatch(ctx, result) {
			if res.Err != nil {
				fmt.Println("consumer error ", res.Err)
				continue
			}
		}
	}

	return nil
}
