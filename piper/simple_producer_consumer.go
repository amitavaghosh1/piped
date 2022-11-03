package piper

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Opts map[string]interface{}

type User struct {
	ID   int
	Name string
}

type SimpleProducer interface {
	Next(ctx context.Context, opts Opts) chan UserResult
	Stop() error
}

type UserProducer struct {
	store   []User
	closing chan chan error
}

type UserResult struct {
	Data *User
	Err  error
}

// New, only used for initialization
func NewUserProducer() *UserProducer {
	return &UserProducer{
		closing: make(chan chan error, 1),
		store: []User{
			{ID: 1, Name: "avai"},
			{ID: 2, Name: "bvai"},
			{ID: 3, Name: "cvai"},
			{ID: 4, Name: "dvai"},
			{ID: 5, Name: "evai"},
			{ID: 6, Name: "fvai"},
			{ID: 7, Name: "gvai"},
			{ID: 8, Name: "hvai"},
			{ID: 9, Name: "ivai"},
			{ID: 10, Name: "jvai"},
			{ID: 11, Name: "kvai"},
			{ID: 12, Name: "lvai"},
			{ID: 13, Name: "mvai"},
			{ID: 14, Name: "nvai"},
			{ID: 15, Name: "ovai"},
			{ID: 16, Name: "pvai"},
			{ID: 17, Name: "qvai"},
			{ID: 19, Name: "rvai"},
			{ID: 20, Name: "svai"},
			{ID: 21, Name: "tvai"},
			{ID: 22, Name: "uvai"},
		}}
}

func (p *UserProducer) Stop() error {
	fmt.Println("stopping")

	errch := make(chan error)
	p.closing <- errch

	return <-errch
}

func (p *UserProducer) Next(ctx context.Context, opts Opts) chan UserResult {
	fmt.Println("invoked")

	res := make(chan UserResult, 1)

	demand, ok := opts["demand"].(int)
	if !ok {
		demand = 1
	}

	var err error

	go func() {
		sig := make(chan bool, 1)
		sig <- true

		defer close(sig)

		for {
			select {
			case errch := <-p.closing:
				errch <- err
				close(res)
				return
			case <-sig:
				var store []User

				if demand > len(p.store)-1 {
					store = p.store
				} else {
					store = p.store[:demand+1]
				}

				if demand > len(store) {
					res <- UserResult{Err: errors.New("no records")}
					close(res)
					return
				}

				p.store = p.store[len(store):]

				for _, item := range store {
					r := item
					res <- UserResult{Data: &r, Err: nil}
				}
			case <-time.After(2 * time.Second):
				sig <- true
			}
		}
	}()

	return res
}

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
