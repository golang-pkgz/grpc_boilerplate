package grpc_boilerplate

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const default_api_token_header_name = "api_token"

func ServerTokenAuth(api_token string, api_token_header_name string) grpc.UnaryServerInterceptor {
	if api_token_header_name == "" {
		api_token_header_name = default_api_token_header_name
	}

	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if api_token != "" {
			if info.FullMethod != "/grpc.health.v1.Health/Check" {
				err := func(ctx context.Context, tokenHeaderName string, correct_token string) error {
					md, ok := metadata.FromIncomingContext(ctx)
					if !ok {
						return status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
					}

					authHeader, ok := md[api_token_header_name]
					if !ok {
						fmt.Println(md)
						return status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
					}

					token := authHeader[0]
					if token != correct_token {
						return status.Errorf(codes.Unauthenticated, "Wrong token")
					}

					return nil
				}(ctx, api_token_header_name, api_token)
				if err != nil {
					return nil, err
				}
			}
		}

		h, err := handler(ctx, req)
		return h, err
	}
}

func ClientTokenAuth(api_token string, api_token_header_name string) grpc.UnaryClientInterceptor {
	if api_token_header_name == "" {
		api_token_header_name = default_api_token_header_name
	}

	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if api_token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, api_token_header_name, api_token)
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}
