package v1

import (
	"context"
	"fmt"
	"piped/piper"
	"sync"
)

type UserConsumer interface {
	Run(ctx context.Context, opts piper.Opts) chan piper.UserResult
}

type UserDetailConsumer struct {
	subscribes_to piper.SimpleProducer
}

func NewUserDetailConsumer(subs piper.SimpleProducer) *UserDetailConsumer {
	return &UserDetailConsumer{subscribes_to: subs}
}

func (c *UserDetailConsumer) Run(ctx context.Context, opts piper.Opts) chan piper.UserResult {
	res := make(chan piper.UserResult, 1)

	go func() {
		for result := range c.subscribes_to.Next(ctx, opts) {
			if result.Err == nil {
				fmt.Println("receiving ", result.Data.ID)
			}

			res <- result
		}

		close(res)
	}()

	return res
}

type SimpleSupervisor struct {
	Consumer UserConsumer
}

func (sv *SimpleSupervisor) Run(ctx context.Context, opts piper.Opts) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			for result := range sv.Consumer.Run(ctx, opts) {
				if result.Err != nil {
					fmt.Println("error ", result.Err)
					return
				}
			}
		}
	}()

	wg.Wait()
}
