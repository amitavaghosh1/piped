package piper

import (
	"context"
	"fmt"
	"sync"
)

type Opts map[string]interface{}

type UserConsumer interface {
	Run(ctx context.Context, opts Opts) chan UserResult
}

type UserDetailConsumer struct {
	subscribes_to SimpleProducer
}

func NewUserDetailConsumer(subs SimpleProducer) *UserDetailConsumer {
	return &UserDetailConsumer{subscribes_to: subs}
}

func (c *UserDetailConsumer) Run(ctx context.Context, opts Opts) chan UserResult {
	res := make(chan UserResult, 1)

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

func (sv *SimpleSupervisor) Run(ctx context.Context, opts Opts) {
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
