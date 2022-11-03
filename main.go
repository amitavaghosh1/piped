package main

import (
	"context"
	"piped/piper"
)

func main() {
	ctx := context.Background()

	// simpleProducer(ctx)
	// broadCastSupervisor(ctx)
	broadCastSupervisorV2(ctx)
}

func simpleProducer(ctx context.Context) {
	userProducer := piper.NewUserProducer()
	simpleConsumer := piper.NewUserDetailConsumer(userProducer)

	supervisor := &piper.SimpleSupervisor{Consumer: simpleConsumer}
	supervisor.Run(ctx, piper.Opts{"demand": 10})

}

func broadCastSupervisor(ctx context.Context) {
	userProducer := piper.NewUserProducer()
	dispatcher := &piper.UserBroadcastDispatcher{
		Consumers: []piper.UserBroadcastConsumer{
			&piper.UserIDBroadcastConsumer{},
			&piper.UserNameBroadcastConsumer{},
		},
		Subscribes: userProducer,
	}

	supervisor := &piper.UserBroadcastSupervisor{
		Dispatcher: dispatcher,
	}
	supervisor.Run(ctx, piper.Opts{"demand": 10})
}

func broadCastSupervisorV2(ctx context.Context) {
	userProducer := piper.NewUserProducer()
	dispatcher := &piper.UserBroadcastDispatcherV2{
		Consumers: []piper.UserBroadcastConsumer{
			&piper.UserIDBroadcastConsumer{},
			&piper.UserNameBroadcastConsumer{},
		},
	}

	supervisor := &piper.UserBroadcastSupervisorV2{
		Subscribes: userProducer,
		Dispatcher: dispatcher,
	}

	supervisor.Run(ctx, piper.Opts{"demand": 10})

}
