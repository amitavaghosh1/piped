package piper

import (
	"context"
	"fmt"
	"piped/loadbalancers"
	"sync"
)

type UserDispatcherV2 interface {
	Dispatch(ctx context.Context, result UserResult) chan UserResult
}

type UserBroadcastDispatcherV2 struct {
	Consumers []UserBroadcastConsumer
}

func (bd *UserBroadcastDispatcherV2) Dispatch(ctx context.Context, result UserResult) chan UserResult {
	res := make(chan UserResult, 1)

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

func (bd *UserBroadcastDispatcherV2) work(ctx context.Context, w *sync.WaitGroup, consumer UserBroadcastConsumer, result UserResult, res chan UserResult) {
	defer w.Done()

	res <- UserResult{
		Data: result.Data,
		Err:  consumer.Handle(ctx, result),
	}
}

type UserConsumerGroupDispatcher struct {
	Consumers []UserBroadcastConsumer
	Balancer  loadbalancers.Balancer
}

func (cd *UserConsumerGroupDispatcher) Dispatch(ctx context.Context, result UserResult) chan UserResult {
	res := make(chan UserResult, 1)

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
		res <- UserResult{Data: result.Data, Err: consumer.Handle(ctx, result)}
	}()

	return res

}

type UserBroadcastSupervisorV2 struct {
	Dispatcher UserDispatcherV2
	Subscribes SimpleProducer
}

func (bs *UserBroadcastSupervisorV2) Run(ctx context.Context, opts Opts) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			for result := range bs.Subscribes.Next(ctx, opts) {
				if result.Err != nil {
					fmt.Println("error ", result.Err)
					return
				}

				for res := range bs.Dispatcher.Dispatch(ctx, result) {
					if res.Err != nil {
						fmt.Println("consumer error ", res.Err)
						continue
					}
				}
			}
		}
	}()

	wg.Wait()
}

type UserConsumerGroupSupervisor struct {
	Dispatcher UserDispatcherV2
	Subscribes SimpleProducer
}

func (cs *UserConsumerGroupSupervisor) Run(ctx context.Context, opts Opts) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			for result := range cs.Subscribes.Next(ctx, opts) {
				if result.Err != nil {
					fmt.Println("error ", result.Err)
					return
				}

				for res := range cs.Dispatcher.Dispatch(ctx, result) {
					if res.Err != nil {
						fmt.Println("consumer error ", res.Err)
						continue
					}
				}
			}
		}
	}()

	wg.Wait()
}
