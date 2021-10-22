package connectionstring

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/golang-pkgz/grpc_boilerplate/token_auth"
	"github.com/gorilla/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type optionsParser struct {
	ServerCrt string `schema:"ServerCrt"`
}

// Parse grpc connectionstring `h2c|h2cs://[<token>@]host:port[?ServerCrt=<path to server cert>]`
// Attempt to create generic connectionstring format for grpc connections
//
// schema
// * h2c: specifies insecure connection
// * h2cs: specifies secure connection
//         Optional loads server cert with `ServerCrt` query option
//
// token
// specifies token for token_auth interceptor
//
// host:port
// required server address and port
//
// ServerCrt
//		Path to server certificate (server.crt)
//		Works with TLS enabled (h2cs://...)
//		Has no effect if connection schema is insecure
func ParseConnectionString(connectionString string) (hostPort string, dialOptions []grpc.DialOption, err error) {
	parsed, err := url.Parse(connectionString)
	if err != nil {
		return
	}

	// Parse host:port
	if !strings.Contains(parsed.Host, ":") {
		err = fmt.Errorf("host:port must contain port: '%s'", parsed.Host)
		return
	}
	hostPort = parsed.Host

	// Validate schema
	if parsed.Scheme != "h2c" && parsed.Scheme != "h2cs" {
		err = fmt.Errorf("unknown scheme: '%s'", parsed.Scheme)
		return
	}

	// Parse query options
	dialOptions = make([]grpc.DialOption, 0)
	queryOptions := optionsParser{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(false)
	err = decoder.Decode(&queryOptions, parsed.Query())
	if err != nil {
		return
	}
	if parsed.Scheme == "h2cs" {
		// By default use empty tls config
		creds := credentials.NewTLS(&tls.Config{})

		// If server cert specified, load it
		if queryOptions.ServerCrt != "" {
			creds, err = credentials.NewClientTLSFromFile(queryOptions.ServerCrt, "")
			if err != nil {
				err = fmt.Errorf("could not load tls cert: %s", err)
				return
			}
		}

		dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))
	} else {
		if queryOptions.ServerCrt != "" {
			err = errors.New("ServerCrt cannot be set with insecure connection")
			return
		}
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	// Token auth support
	token := ""
	if parsed.User != nil {
		token = parsed.User.Username()
	}
	if token != "" {
		dialOptions = append(dialOptions, grpc.WithUnaryInterceptor(token_auth.ClientTokenAuth(token, "")))
	}

	return
}
