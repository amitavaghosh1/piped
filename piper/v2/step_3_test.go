package v2

import (
	"context"
	"fmt"
	"piped/piper"
	"testing"
)

var table = []struct {
	demand int
}{{demand: 1}, {demand: 10}, {demand: 20}}

func BenchmarkStep3WithDemand(b *testing.B) {
	ctx := context.Background()

	for _, input := range table {
		b.Run(fmt.Sprintf("with demand: %d\n", input.demand), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s3 := &Step3{}
				s3.Supervisor(ctx, piper.Opts{"demand": input.demand})
			}
		})
	}
}

func BenchmarkStep3WithDemand_MultipleConsumer(b *testing.B) {
	ctx := context.Background()

	for _, input := range table {
		b.Run(fmt.Sprintf("with demand: %d\n", input.demand), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s3 := &Step3{}
				s3.Supervisors(ctx, piper.Opts{"demand": input.demand})
			}
		})
	}

}
