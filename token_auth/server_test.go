package token_auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Example server with token auth
func ExampleServerTokenAuth() {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			ServerTokenAuth(
				"secret_token",
				"<header name or empty for default>",
			),
		),
	)
	fmt.Println(grpcServer)
}

func assertGrpcError(t *testing.T, err error, code codes.Code, msg string) {
	status, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, status.Code(), code)
	assert.Equal(t, status.Message(), msg)
}

func TestServerTokenAuth(t *testing.T) {
	mw := ServerTokenAuth("secret", "")

	// No metadata
	resp, err := mw(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{},
		func(ctx context.Context, req interface{}) (interface{}, error) { return true, nil },
	)
	assert.Nil(t, resp)
	assertGrpcError(t, err, codes.Internal, "Retrieving metadata is failed")

	// Healthchecks passed without token
	resp, err = mw(
		metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{})),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/grpc.health.v1.Health/Check"},
		func(ctx context.Context, req interface{}) (interface{}, error) { return true, nil },
	)
	assert.Equal(t, resp, true)
	assert.NoError(t, err)

	// No token
	resp, err = mw(
		metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{})),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/MyService/Ping"},
		func(ctx context.Context, req interface{}) (interface{}, error) { return true, nil },
	)
	assert.Nil(t, resp)
	assertGrpcError(t, err, codes.Unauthenticated, "Authorization token is not supplied")

	// Wrong token
	resp, err = mw(
		metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
			DefaultApiTokenHeaderName: "nosecret",
		})),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/HisService/Ping"},
		func(ctx context.Context, req interface{}) (interface{}, error) { return true, nil },
	)
	assert.Nil(t, resp)
	assertGrpcError(t, err, codes.Unauthenticated, "Wrong token")

	// Good token
	resp, err = mw(
		metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
			DefaultApiTokenHeaderName: "secret",
		})),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/ShkliService/Ping"},
		func(ctx context.Context, req interface{}) (interface{}, error) { return true, nil },
	)
	assert.Equal(t, resp, true)
	assert.Nil(t, err)
}
