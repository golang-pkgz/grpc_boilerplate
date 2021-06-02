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
        grpc_boilerplate.ServerLogDuration(),
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
