package peer_whitelist

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func ExampleServerPeerWhitelist() {
	whitelistNetworks, err := ParseCIDRs([]string{
		"127.0.0.1/24",
	})

	if err != nil {
		log.Fatalln("")
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			ServerPeerWhitelist(whitelistNetworks), // nil for log.Default() as logger
		),
	)
	fmt.Println(grpcServer)
	// ...
}

func TestParseCIDRs(t *testing.T) {
	nets, err := ParseCIDRs([]string{
		"127.0..1/24",
	})
	assert.Error(t, err)
	assert.Empty(t, nets)

	nets, err = ParseCIDRs([]string{
		"127.0.0.1/24",
	})
	assert.NoError(t, err)
	assert.Len(t, nets, 1)
	assert.Equal(t, nets[0], &net.IPNet{IP: []byte{0x7f, 0x0, 0x0, 0x0}, Mask: []byte{0xff, 0xff, 0xff, 0x0}})
}

func assertGrpcError(t *testing.T, err error, code codes.Code, msg string) {
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, st.Code(), code)
	assert.Equal(t, st.Message(), msg)
}

func TestServerPeerWhitelist(t *testing.T) {
	whitelistNetworks, err := ParseCIDRs([]string{
		"127.0.0.1/24",
	})

	assert.NoError(t, err)

	mw := ServerPeerWhitelist(whitelistNetworks)

	// No peer
	resp, err := mw(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{},
		func(ctx context.Context, req interface{}) (interface{}, error) { return true, nil },
	)
	assert.Nil(t, resp)
	assertGrpcError(t, err, codes.Internal, "Retrieving peer is failed")

	// Bad peer
	resp, err = mw(
		peer.NewContext(context.Background(), &peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("1.2.3.4")}}),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/ShkleService/Pong"},
		func(ctx context.Context, req interface{}) (interface{}, error) { return true, nil },
	)
	assert.Nil(t, resp)
	assertGrpcError(t, err, codes.PermissionDenied, "Denied")

	// Healthchecks NOT passed with bad peer
	resp, err = mw(
		peer.NewContext(context.Background(), &peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("1.2.3.4")}}),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/grpc.health.v1.Health/Check"},
		func(ctx context.Context, req interface{}) (interface{}, error) { return true, nil },
	)
	assert.Nil(t, resp)
	assertGrpcError(t, err, codes.PermissionDenied, "Denied")

	// Good peer
	resp, err = mw(
		peer.NewContext(context.Background(), &peer.Peer{Addr: &net.TCPAddr{IP: net.ParseIP("127.0.0.2")}}),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/ShkleService/Pong"},
		func(ctx context.Context, req interface{}) (interface{}, error) { return true, nil },
	)
	assert.Equal(t, resp, true)
	assert.NoError(t, err)
}
