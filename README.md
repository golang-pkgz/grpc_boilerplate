# Grpc boilerplate
Grpc server and client staff for golang

## Install

```sh
go get github.com/golang-pkgz/grpc_boilerplate
```

## Staff
### Log duration interceptor for server
Log every request duration

Example server
```go
grpcServer := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        grpc_boilerplate.ServerLogDuration(nil), // nil for log.Default() as logger
    ),
)
```

### Token authentication for server and client
Provide token header for authentication

Example server
```golang
grpcServer := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        grpc_boilerplate.ServerTokenAuth(
            "secret_token",
            "<header name or empty for default>",
        ),
    ),
)
```

Example client
```golang
conn, err := grpc.Dial(
    opts.ServerAddr,
    grpc.WithBlock(),
    grpc.WithInsecure(),
    grpc.WithUnaryInterceptor(grpc_boilerplate.ClientTokenAuth(
        "secret_token",
        "<header name or empty for default>",
    )),
)
```

### Dial from connectionstring for client
```golang
conn, err := grpc_boilerplate.DialFromConnectionString("h2c://clientTokenSecret@localhost:50002", grpc_boilerplate.DIAL_OPTS_DEFAULT...)

if err != nil {
    log.Fatalf("fail to dial: %v", err)
}
defer conn.Close()
...

```