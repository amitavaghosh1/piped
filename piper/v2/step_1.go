// Simple producer consumer
package v2

import (
	"context"
	"fmt"
	"piped/piper"

	"github.com/sirupsen/logrus"
)

type Result[E any] struct {
	Data E
	Err  error
}

type User *piper.User

type Step1[E User] struct {
	subscriber *piper.UserProducer
}

func NewStep1[E User](subscriber *piper.UserProducer) *Step1[E] {
	return &Step1[E]{subscriber: subscriber}
}

func (s *Step1[E]) producer(ctx context.Context, opts piper.Opts) chan Result[E] {
	ch := make(chan Result[E], 1)

	go func() {
		defer close(ch)

		for res := range s.subscriber.Next(ctx, opts) {
			ch <- Result[E]{Data: res.Data, Err: res.Err}

			if res.Err != nil {
				return
			}
		}
	}()

	return ch
}

func (s *Step1[E]) consumer(ctx context.Context, in chan Result[E]) error {
	for result := range in {
		if result.Err != nil {
			logrus.Error(result.Err)
			return result.Err
		}

		var user *piper.User = result.Data
		fmt.Println(user.Name, "|", user.GetRank())
	}

	return nil
}

func (s *Step1[E]) supervisor(ctx context.Context, opts piper.Opts) {
	for {
		in := s.producer(ctx, opts)
		if err := s.consumer(ctx, in); err != nil {
			return
		}
	}
}
