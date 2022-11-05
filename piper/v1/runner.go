package v1

import (
	"context"
	"piped/loadbalancers"
	"piped/piper"
)

func Run(ctx context.Context) {
	simpleProducer(ctx)
	broadCastSupervisorV2(ctx)
	consumerGroupSupervisor(ctx)
}

func simpleProducer(ctx context.Context) {
	userProducer := piper.NewUserProducer()
	simpleConsumer := NewUserDetailConsumer(userProducer)

	supervisor := &SimpleSupervisor{Consumer: simpleConsumer}
	supervisor.Run(ctx, piper.Opts{"demand": 10})

}

func broadCastSupervisor(ctx context.Context) {
	userProducer := piper.NewUserProducer()
	dispatcher := &UserBroadcastDispatcher{
		Consumers: []UserBroadcastConsumer{
			&UserIDBroadcastConsumer{},
			&UserNameBroadcastConsumer{},
		},
		Subscribes: userProducer,
	}

	supervisor := &UserBroadcastSupervisor{
		Dispatcher: dispatcher,
	}
	supervisor.Run(ctx, piper.Opts{"demand": 10})
}

func broadCastSupervisorV2(ctx context.Context) {
	userProducer := piper.NewUserProducer()
	dispatcher := &UserBroadcastDispatcherV2{
		Consumers: []UserBroadcastConsumer{
			&UserIDBroadcastConsumer{},
			&UserNameBroadcastConsumer{},
		},
	}

	supervisor := &UserBroadcastSupervisorV2{
		Subscribes: userProducer,
		Dispatcher: dispatcher,
	}

	supervisor.Run(ctx, piper.Opts{"demand": 10})
}

func consumerGroupSupervisor(ctx context.Context) {
	userProducer := piper.NewUserProducer()
	dispatcher := &UserConsumerGroupDispatcher{
		Consumers: []UserBroadcastConsumer{
			&UserIDBroadcastConsumer{},
			&UserIDBroadcastConsumer{},
		},
		Balancer: loadbalancers.NewRoundRobalancer(2),
	}

	supervisor := &UserConsumerGroupSupervisor{
		Subscribes: userProducer,
		Dispatcher: dispatcher,
	}

	supervisor.Run(ctx, piper.Opts{"demand": 10})
}
