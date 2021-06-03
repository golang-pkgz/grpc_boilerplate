package grpc_boilerplate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseConnectionString(t *testing.T) {
	hostport, token, err := parseConnectionString("h2c://localhost:50002")
	assert.Equal(t, hostport, "localhost:50002")
	assert.Equal(t, token, "")
	assert.Nil(t, err)

	hostport, token, err = parseConnectionString("h2c://localhost:50002")
	assert.Equal(t, hostport, "localhost:50002")
	assert.Equal(t, token, "")
	assert.Nil(t, err)

	hostport, token, err = parseConnectionString("h2c://secret@localhost:50002")
	assert.Equal(t, hostport, "localhost:50002")
	assert.Equal(t, token, "secret")
	assert.Nil(t, err)

	hostport, token, err = parseConnectionString("h2c://secret@localhost")
	assert.Equal(t, hostport, "")
	assert.Equal(t, token, "")
	assert.Equal(t, err.Error(), "host:port does contain port: 'localhost'")
}
