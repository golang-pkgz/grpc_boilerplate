all: certs lint test

certs:
	mkdir certs
	openssl genrsa -out certs/server.key 4096  # server.key: a private RSA key to sign and authenticate the public key
	openssl req -new -x509 -sha256 -key certs/server.key -subj '/CN=server' -out certs/server.crt -days 3650  # server.pem/server.crt: self-signed X.509 public keys for distribution

lint:
	golangci-lint run

.PHONY: test
test: lint
	go test -coverprofile cover.out ./...
	go tool cover -func cover.out
	rm cover.out
