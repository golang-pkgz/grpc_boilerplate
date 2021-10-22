package grpc_boilerplate

import (
	"github.com/golang-pkgz/grpc_boilerplate/connectionstring"
	"google.golang.org/grpc"
)

// Connect to grpc sever from connectionString `h2c|h2cs://[<token>@]host:port?<options>`
// See `connectionstring.ParseConnectionString` for `options` description
func DialFromConnectionString(userAgent string, connectionString string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	hostPort, parsed_opts, err := connectionstring.ParseConnectionString(connectionString)
	opts = append(opts, parsed_opts...)

	if err != nil {
		return nil, err
	}

	if userAgent != "" {
		opts = append(parsed_opts, grpc.WithUserAgent(userAgent))
	}

	return grpc.Dial(hostPort, opts...)
}
