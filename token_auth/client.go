package token_auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Provide token header for authentication
func ClientTokenAuth(apiToken string, apiTokenHeaderName string) grpc.UnaryClientInterceptor {
	if apiTokenHeaderName == "" {
		apiTokenHeaderName = DefaultApiTokenHeaderName
	}

	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if apiToken != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, apiTokenHeaderName, apiToken)
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}
