package main

import (
	"context"
	// v1 "piped/piper/v1"
	v2 "piped/piper/v2"
)

func main() {
	ctx := context.Background()
	// v1.Run(ctx)

	v2.Run(ctx)
}
