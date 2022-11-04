package piper

import (
	"context"
	"errors"
	"fmt"
	"time"
)

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
