package piper

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Opts map[string]interface{}

type User struct {
	ID    int
	Name  string
	Kills int
}

func (u *User) GetRank() string {
	kills := u.Kills

	if kills < 3 {
		return "Novice"
	}

	if kills < 7 {
		return "Pro"
	}

	return "Veteran"
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
			{ID: 1, Name: "avai", Kills: 0},
			{ID: 2, Name: "bvai", Kills: 1},
			{ID: 3, Name: "cvai", Kills: 2},
			{ID: 4, Name: "dvai", Kills: 3},
			{ID: 5, Name: "evai", Kills: 4},
			{ID: 6, Name: "fvai", Kills: 5},
			{ID: 7, Name: "gvai", Kills: 6},
			{ID: 8, Name: "hvai", Kills: 7},
			{ID: 9, Name: "ivai", Kills: 8},
			{ID: 10, Name: "jvai", Kills: 9},
			{ID: 11, Name: "kvai", Kills: 0},
			{ID: 12, Name: "lvai", Kills: 1},
			{ID: 13, Name: "mvai", Kills: 2},
			{ID: 14, Name: "nvai", Kills: 3},
			{ID: 15, Name: "ovai", Kills: 4},
			{ID: 16, Name: "pvai", Kills: 5},
			{ID: 17, Name: "qvai", Kills: 6},
			{ID: 19, Name: "rvai", Kills: 7},
			{ID: 20, Name: "svai", Kills: 8},
			{ID: 21, Name: "tvai", Kills: 9},
			{ID: 22, Name: "uvai", Kills: 20},
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
