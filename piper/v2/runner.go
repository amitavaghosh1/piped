package v2

import (
	"context"
	"piped/piper"
)

func Run(ctx context.Context) {
	opts := piper.Opts{"demand": 10}

	// s1 := NewStep1(piper.NewUserProducer())
	// s1.supervisor(ctx, opts)

	// s2 := NewStep2(piper.NewUserProducer())
	// s2.supervisor(ctx, opts)

	s3 := &Step3{}
	s3.Supervisors(ctx, opts)
}
