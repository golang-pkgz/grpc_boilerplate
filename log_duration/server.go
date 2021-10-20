package log_duration

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

// ServerLogDuration logs every request method, duration and error
func ServerLogDuration(logger *log.Logger) grpc.UnaryServerInterceptor {
	if logger == nil {
		logger = log.Default()
	}

	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		h, err := handler(ctx, req)

		logger.Printf("Request - Method:%s\tDuration:%s\tError:%v\n", info.FullMethod, time.Since(start), err)
		return h, err
	}
}
