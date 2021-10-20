package log_duration

import (
	"fmt"

	"google.golang.org/grpc"
)

func ExampleServerLogDuration() {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			ServerLogDuration(nil), // nil for log.Default() as logger
		),
	)
	fmt.Println(grpcServer)
	// ...
}
