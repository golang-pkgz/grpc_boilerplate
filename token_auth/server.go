package token_auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// DefaultApiTokenHeaderName used if apiTokenHeaderName is empty string
const DefaultApiTokenHeaderName = "api_token"

// ServerTokenAuth provides token header for authentication
func ServerTokenAuth(apiToken string, apiTokenHeaderName string) grpc.UnaryServerInterceptor {
	if apiTokenHeaderName == "" {
		apiTokenHeaderName = DefaultApiTokenHeaderName
	}

	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Allow all requests to healthcheck api
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Internal, "Retrieving metadata is failed")
		}

		authHeader, ok := md[apiTokenHeaderName]
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
		}

		token := authHeader[0]
		if token != apiToken {
			return nil, status.Errorf(codes.Unauthenticated, "Wrong token")
		}

		return handler(ctx, req)
	}
}
