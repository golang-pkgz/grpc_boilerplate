package grpc_boilerplate

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func ServerLogDuration() grpc.UnaryServerInterceptor {

	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		h, err := handler(ctx, req)

		log.Printf("Request - Method:%s\tDuration:%s\tError:%v\n", info.FullMethod, time.Since(start), err)
		return h, err
	}
}
