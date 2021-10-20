package grpc_boilerplate

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/golang-pkgz/grpc_boilerplate/token_auth"
	"google.golang.org/grpc"
)

// DIAL_OPTS_DEFAULT useful with DialFromConnectionString
var DIAL_OPTS_DEFAULT []grpc.DialOption = []grpc.DialOption{
	grpc.WithBlock(),
	grpc.WithInsecure(),
}

func parseConnectionString(connectionString string) (string, string, error) {
	parsed, err := url.Parse(connectionString)
	if err != nil {
		return "", "", err
	}

	if parsed.Scheme != "h2c" && parsed.Scheme != "h2cs" {
		return "", "", fmt.Errorf("unknown scheme: '%s'", parsed.Scheme)
	}

	if parsed.Scheme == "h2cs" {
		return "", "", errors.New("h2cs scheme is not supported for now")
	}

	if !strings.Contains(parsed.Host, ":") {
		return "", "", fmt.Errorf("host:port does contain port: '%s'", parsed.Host)
	}

	return parsed.Host, parsed.User.Username(), nil
}

// DialFromConnectionString
// Connect from connectionString `h2c|h2cs://[<token>@]host:port`
//
// Usage:
// conn, err := grpc_boilerplate.DialFromConnectionString(cs, grpc_boilerplate.DIAL_OPTS_DEFAULT)
func DialFromConnectionString(connectionString string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	hostport, token, err := parseConnectionString(connectionString)

	if err != nil {
		return nil, err
	}

	if token != "" {
		opts = append(opts, grpc.WithUnaryInterceptor(token_auth.ClientTokenAuth(token, "")))
	}

	return grpc.Dial(hostport, opts...)
}
