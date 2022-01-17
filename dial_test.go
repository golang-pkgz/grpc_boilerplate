package grpc_boilerplate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDialFromConnectionString(t *testing.T) {
	_, err := DialFromConnectionString("a\na", "h2c://127.0.0.1:50002")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "newline symbol not allowed in user agent string")
}
