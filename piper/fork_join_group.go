package piper

import (
	"context"
	"fmt"
	"sync"
)

// Incomplete
type UserConsumerGroup struct {
	Consumer UserSyncConsumer

	MaxWorkers int

	taskChan chan UserResult
	wg       sync.WaitGroup
}

func NewConsumerGroup(consumer UserSyncConsumer, workerCount int) *UserConsumerGroup {
	return &UserConsumerGroup{
		Consumer:   consumer,
		MaxWorkers: workerCount,
		taskChan:   make(chan UserResult),
	}
}

func (cg *UserConsumerGroup) Run(ctx context.Context, result UserResult) chan UserResult {
	res := make(chan UserResult, 1)

	return res
}

func (cg *UserConsumerGroup) work(ctx context.Context) {
	for result := range cg.taskChan {
		res := result
		cg.Consumer.Run(ctx, res)
	}
}

type UserSyncConsumer struct {
}

func (uc *UserSyncConsumer) Run(ctx context.Context, result UserResult) error {
	fmt.Println("receiving user ", result.Data.ID)
	return nil
}
