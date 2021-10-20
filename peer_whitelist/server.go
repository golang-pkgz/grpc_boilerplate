package peer_whitelist

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// ServerPeerWhitelist provides server client whitelist by CIDR(s)
func ServerPeerWhitelist(whitelist []*net.IPNet) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		requestPeer, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Internal, "Retrieving peer is failed")
		}

		// Convert peer net.Addr to net.Ip
		var peerIp net.IP
		switch addr := requestPeer.Addr.(type) {
		case *net.UDPAddr:
			peerIp = addr.IP
		case *net.TCPAddr:
			peerIp = addr.IP
		default:
			return nil, status.Errorf(codes.Internal, "Retrieving peer ip is failed")
		}

		for _, network := range whitelist {
			if network.Contains(peerIp) {
				return handler(ctx, req)
			}
		}

		return nil, status.Errorf(codes.PermissionDenied, "Denied")
	}
}

// ParseCIDRs shortcut for easy converting CIDR(s) string slice to net.IPNet slice
func ParseCIDRs(cidrs []string) ([]*net.IPNet, error) {
	nets := make([]*net.IPNet, 0)

	for _, cidr := range cidrs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return nets, err
		}

		nets = append(nets, network)
	}

	return nets, nil
}
