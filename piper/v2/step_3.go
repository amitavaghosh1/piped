// Reversing the relationship between consumer and producer
// Step 1 to 2, the producer produces and consumer consumers from the channel
// But, this can lead to  back pressure

package v2

import (
	"context"
	"math/rand"
	"piped/piper"
	"sync"

	"github.com/sirupsen/logrus"
)

// type Producer[E User] func(context.Context, piper.Opts) chan Result[E]
// type ProducerConsumer[E User] func(context.Context, Producer[E], piper.Opts) chan Result[E]

type Producer[E any] interface {
	Next(context.Context, piper.Opts) chan Result[E]
}

type UserProducer[E User] struct {
	subscriber *piper.UserProducer
	once       *sync.Once
	updates    chan Result[E]
}

func NewUserProducer[E User](subscriber *piper.UserProducer, opts piper.Opts) *UserProducer[E] {
	return &UserProducer[E]{
		subscriber: subscriber,
		once:       &sync.Once{},
		updates:    make(chan Result[E], getDemandOrDefault(opts)),
	}
}

func (up *UserProducer[E]) Next(ctx context.Context, opts piper.Opts) chan Result[E] {
	up.once.Do(func() { go up.loop(ctx, opts) })
	return up.updates
}

func (up *UserProducer[E]) loop(ctx context.Context, opts piper.Opts) {
	// defer func() {
	// 	logrus.Println("up closing")
	// }()
	defer close(up.updates)

	for res := range up.subscriber.Next(ctx, opts) {
		up.updates <- Result[E]{Data: res.Data, Err: res.Err}

		if res.Err != nil {
			return
		}
	}
}

// Behaves much like a producer
// type ProducerConsumer[E User] interface {
// 	Next(context.Context, piper.Opts) chan Result[E]
// }

type KillMonger[E User] struct {
	Producer Producer[E]
	updates  chan Result[E]
}

func NewKillerProducerConsumer[E User](producer Producer[E], opts piper.Opts) *KillMonger[E] {
	return &KillMonger[E]{
		Producer: producer,
		updates:  make(chan Result[E], getDemandOrDefault(opts)),
	}
}

func (pc *KillMonger[E]) Next(ctx context.Context, opts piper.Opts) chan Result[E] {
	// go pc.loop(ctx, opts, pc.updates)
	// return pc.updates

	ch := make(chan Result[E], getDemandOrDefault(opts))
	go pc.loop(ctx, opts, ch)

	return ch
}

func (pc *KillMonger[E]) loop(ctx context.Context, opts piper.Opts, ch chan Result[E]) {
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

		// var user *piper.User = result.Data
		// fmt.Println(user.Name, "|", user.Kills, "|", user.GetRank())
	}

	return nil
}

type Step3 struct {
}

func (s *Step3) Supervisor(ctx context.Context, opts piper.Opts) {
	userProducer := NewUserProducer(piper.NewUserProducer(opts), opts)
	killMonger := NewKillerProducerConsumer[User](userProducer, opts)

	consumer := &UserConsumer[User]{LinksTo: killMonger}

	for {
		if err := consumer.Handle(ctx, opts); err != nil {
			// logrus.Error(err)
			return
		}
	}
}

func (s *Step3) Supervisors(ctx context.Context, opts piper.Opts) {
	userProducer := NewUserProducer(piper.NewUserProducer(opts), opts)
	killMonger := NewKillerProducerConsumer[User](userProducer, opts)

	consumers := []*UserConsumer[User]{
		{LinksTo: killMonger},
		{LinksTo: killMonger},
	}

	var wg sync.WaitGroup
	wg.Add(len(consumers))

	for i, consumer := range consumers {
		go func(w *sync.WaitGroup, c *UserConsumer[User], idx int) {
			defer w.Done()

			for {
				// logrus.Println("waiting ", idx)
				if err := c.Handle(ctx, opts); err != nil {
					logrus.Error(err)
					return
				}
			}
		}(&wg, consumer, i)
	}

	wg.Wait()
}

func getDemandOrDefault(opts piper.Opts) int {
	demand, ok := opts["demand"].(int)
	if !ok {
		return 1
	}

	return demand
}
