// Introducing dispatcher for consumers multiplexing
package v2

// import (
// 	"context"
// 	"fmt"
// 	"math/rand"
// 	"piped/piper"

// 	"github.com/sirupsen/logrus"
// )

// type Step4[E User] struct {
// 	subscriber *piper.UserProducer
// }

// func NewStep4[E User](subscriber *piper.UserProducer) *Step4[E] {
// 	return &Step4[E]{subscriber: subscriber}
// }
// func (s *Step4[E]) producer(ctx context.Context, opts piper.Opts) chan Result[E] {
// 	ch := make(chan Result[E], 1)

// 	go func() {
// 		defer close(ch)

// 		for {
// 			for res := range s.subscriber.Next(ctx, opts) {
// 				ch <- Result[E]{Data: res.Data, Err: res.Err}

// 				if res.Err != nil {
// 					return
// 				}
// 			}
// 		}
// 	}()

// 	return ch
// }

// func (s *Step4[E]) killGenerator(ctx context.Context, opts piper.Opts, in chan Result[E]) chan Result[E] {
// 	ch := make(chan Result[E], 1)

// 	go func() {
// 		defer close(ch)

// 		for res := range in {
// 			if res.Err != nil {
// 				ch <- Result[E]{Data: res.Data, Err: res.Err}
// 				return
// 			}

// 			var user *piper.User = res.Data
// 			user.Kills = user.Kills + (rand.Intn(8) + 4)

// 			ch <- Result[E]{Data: user, Err: res.Err}
// 		}
// 	}()

// 	return ch
// }

// func (s *Step4[E]) consumer(ctx context.Context, in chan Result[E]) error {
// 	for result := range in {
// 		if result.Err != nil {
// 			logrus.Error(result.Err)
// 			return result.Err
// 		}

// 		var user *piper.User = result.Data
// 		fmt.Println(user.Name, "|", user.Kills, "|", user.GetRank())
// 	}

// 	return nil
// }

// func (s *Step4[E]) supervisor(ctx context.Context, opts piper.Opts) {
// 	in := s.producer(ctx, opts)
// 	killgen := s.killGenerator(ctx, opts, in)
// 	if err := s.consumer(ctx, killgen); err != nil {
// 		return
// 	}
// }
