package v2

// import (
// 	"context"
// 	"fmt"
// 	"piped/piper"
// 	"testing"
// )

// func BenchmarkStep2WithDemand(b *testing.B) {
// 	ctx := context.Background()

// 	for _, input := range table {
// 		b.Run(fmt.Sprintf("with demand: %d\n", input.demand), func(b *testing.B) {
// 			for i := 0; i < b.N; i++ {
// 				s2 := NewStep2(piper.NewUserProducer())
// 				s2.supervisor(ctx, piper.Opts{"demand": input.demand})

// 			}
// 		})
// 	}

// }
