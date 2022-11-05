// Reversing the relationship between consumer and producer
// Step 1 to 2, the producer produces and consumer consumers from the channel
// But, this can lead to  back pressure

package v2

import (
	"context"
	"fmt"
	"math/rand"
	"piped/piper"

	"github.com/sirupsen/logrus"
)

// type Producer[E User] func(context.Context, piper.Opts) chan Result[E]
// type ProducerConsumer[E User] func(context.Context, Producer[E], piper.Opts) chan Result[E]

type Producer[E any] interface {
	Next(context.Context, piper.Opts) chan Result[E]
}

type UserProducer[E User] struct {
	subscriber *piper.UserProducer
}

func NewUserProducer[E User](subscriber *piper.UserProducer) *UserProducer[E] {
	return &UserProducer[E]{subscriber: subscriber}
}

func (up *UserProducer[E]) Next(ctx context.Context, opts piper.Opts) chan Result[E] {
	ch := make(chan Result[E], getDemandOrDefault(opts))

	go func() {
		defer close(ch)

		for res := range up.subscriber.Next(ctx, opts) {
			ch <- Result[E]{Data: res.Data, Err: res.Err}

			if res.Err != nil {
				return
			}
		}
	}()

	return ch
}

// Behaves much like a producer
// type ProducerConsumer[E User] interface {
// 	Next(context.Context, piper.Opts) chan Result[E]
// }

type KillMonger[E User] struct {
	Producer Producer[E]
}

func NewKillerProducerConsumer[E User](producer Producer[E]) *KillMonger[E] {
	return &KillMonger[E]{Producer: producer}
}

func (pc *KillMonger[E]) Next(ctx context.Context, opts piper.Opts) chan Result[E] {
	ch := make(chan Result[E], getDemandOrDefault(opts))

	go func(ctx context.Context, o piper.Opts) {
		defer close(ch)

		for res := range pc.Producer.Next(ctx, opts) {
			if res.Err != nil {
				ch <- res
				continue
			}

			var user *piper.User = res.Data
			user.Kills = user.Kills + (rand.Intn(8) + 4)

			ch <- Result[E]{Data: user, Err: res.Err}
		}
	}(ctx, opts)

	return ch
}

type Consumer[E any] interface {
	Handle(ctx context.Context, opts piper.Opts) error
}

type UserConsumer[E User] struct {
	LinksTo Producer[E]
}

func (uc *UserConsumer[E]) Handle(ctx context.Context, opts piper.Opts) error {
	for result := range uc.LinksTo.Next(ctx, opts) {
		if result.Err != nil {
			return result.Err
		}

		var user *piper.User = result.Data
		fmt.Println(user.Name, "|", user.Kills, "|", user.GetRank())
	}

	return nil
}

type Step3 struct {
}

func (s *Step3) Supervisor(ctx context.Context, opts piper.Opts) {
	userProducer := NewUserProducer(piper.NewUserProducer())
	killMonger := NewKillerProducerConsumer[User](userProducer)

	consumer := &UserConsumer[User]{LinksTo: killMonger}

	for {
		if err := consumer.Handle(ctx, opts); err != nil {
			logrus.Error(err)
			return
		}
	}
}

func getDemandOrDefault(opts piper.Opts) int {
	demand, ok := opts["demand"].(int)
	if !ok {
		return 1
	}

	// demand = 1 // for now
	return demand
}
