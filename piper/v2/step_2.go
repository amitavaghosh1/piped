// Improving the api for step_2

package v2

import (
	"context"
	"math/rand"
	"piped/piper"

	"github.com/sirupsen/logrus"
)

type Step2[E User] struct {
	subscriber *piper.UserProducer
}

func NewStep2[E User](subscriber *piper.UserProducer) *Step2[E] {
	return &Step2[E]{subscriber: subscriber}
}

func (s *Step2[E]) producer(ctx context.Context, opts piper.Opts) chan Result[E] {
	ch := make(chan Result[E], getDemandOrDefault(opts))

	go func() {
		defer close(ch)

		for {
			for res := range s.subscriber.Next(ctx, opts) {
				ch <- Result[E]{Data: res.Data, Err: res.Err}

				if res.Err != nil {
					return
				}
			}
		}
	}()

	return ch
}

func (s *Step2[E]) killGenerator(ctx context.Context, opts piper.Opts, in chan Result[E]) chan Result[E] {
	ch := make(chan Result[E], getDemandOrDefault(opts))

	go func() {
		defer close(ch)

		for res := range in {
			if res.Err != nil {
				ch <- Result[E]{Data: res.Data, Err: res.Err}
				return
			}

			var user *piper.User = res.Data
			user.Kills = user.Kills + (rand.Intn(8) + 4)

			ch <- Result[E]{Data: user, Err: res.Err}
		}
	}()

	return ch
}

func (s *Step2[E]) consumer(ctx context.Context, in chan Result[E]) error {
	for result := range in {
		if result.Err != nil {
			logrus.Error(result.Err)
			return result.Err
		}

		// var user *piper.User = result.Data
		// fmt.Println(user.Name, "|", user.Kills, "|", user.GetRank())
	}

	return nil
}

func (s *Step2[E]) supervisor(ctx context.Context, opts piper.Opts) {
	in := s.producer(ctx, opts)
	killgen := s.killGenerator(ctx, opts, in)
	if err := s.consumer(ctx, killgen); err != nil {
		return
	}
}
