package v1

import (
	"context"
	"math/rand"
	"piped/piper"
)

// SimpleProducer
// ProducerConsumer
// Consumer

func newProducerConsumer(ctx context.Context, producer piper.SimpleProducer, dispatcher UserBroadcastDispatcherV2) chan piper.UserResult {
	opts := piper.Opts{"demand": 10}

	resultChan := make(chan piper.UserResult, 1)

	go func() {

		for result := range producer.Next(ctx, opts) {
			go func(r piper.UserResult) {
				for res := range dispatcher.Dispatch(ctx, r) {
					resultChan <- res
				}
			}(result)
		}
	}()

	return resultChan
}

type GenerateKills struct {
}

func (g *GenerateKills) Handle(ctx context.Context, result piper.UserResult) piper.UserResult {
	if result.Err != nil {
		return result
	}

	r := result
	r.Data.Kills = r.Data.Kills + rand.Intn(8)

	return r
}

type KillDispatcher struct {
	Consumer []UserBroadcastConsumer
}

func (kd *KillDispatcher) Dispatch(ctx context.Context, result piper.UserResult) chan piper.UserResult {
	return make(chan piper.UserResult)
}

type UserProducerConsumer struct {
	Dispatcher UserBroadcastDispatcherV2
}

func (pc *UserProducerConsumer) Handle(ctx context.Context, result piper.UserResult, opts piper.Opts) chan piper.UserResult {
	resultChan := make(chan piper.UserResult, 1)

	go func() {
		for res := range pc.Dispatcher.Dispatch(ctx, result) {
			resultChan <- res
		}
	}()

	return resultChan
}
