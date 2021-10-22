package connectionstring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConnectionString(t *testing.T) {
	hostport, _, err := ParseConnectionString("h2c://localhost:50002")
	assert.Equal(t, hostport, "localhost:50002")
	assert.Nil(t, err)

	hostport, _, err = ParseConnectionString("h2c://localhost:50002")
	assert.Equal(t, hostport, "localhost:50002")
	assert.Nil(t, err)

	hostport, _, err = ParseConnectionString("h2c://secret@localhost:50002")
	assert.Equal(t, hostport, "localhost:50002")
	assert.Nil(t, err)

	_, _, err = ParseConnectionString("h2c://secret@localhost")
	assert.Equal(t, err.Error(), "host:port must contain port: 'localhost'")
}
