package token_auth

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Example client with token auth
func ExampleClientTokenAuth() {
	conn, err := grpc.Dial(
		"localhost:50002",
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(ClientTokenAuth(
			"secret_token",
			"<header name or empty for default>",
		)),
	)
	if err != nil {
		log.Fatalf("dial failed: %s", err.Error())
	}
	defer conn.Close()
	// ...
}

func checkAuthInvoker(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("")
	}

	authHeader, ok := md[DefaultApiTokenHeaderName]
	if !ok {
		return errors.New("")
	}

	token := authHeader[0]
	if token != "secret" {
		return errors.New("")
	}

	return nil
}

func TestClientTokenAuth(t *testing.T) {
	ta := ClientTokenAuth("secret", "")

	// No metadata
	err := ta(
		context.Background(),
		"",
		nil, nil,
		&grpc.ClientConn{},
		checkAuthInvoker,
	)
	assert.Error(t, err)

	// No token
	err = ta(
		metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{})),
		"",
		nil, nil,
		&grpc.ClientConn{},
		checkAuthInvoker,
	)
	assert.Error(t, err)

	// Wrong token
	err = ta(
		metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
			DefaultApiTokenHeaderName: "nosecret",
		})),
		"",
		nil, nil,
		&grpc.ClientConn{},
		checkAuthInvoker,
	)
	assert.Error(t, err)

	// Good token
	err = ta(
		metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
			DefaultApiTokenHeaderName: "secret",
		})),
		"",
		nil, nil,
		&grpc.ClientConn{},
		checkAuthInvoker,
	)
	assert.NoError(t, err)
}
