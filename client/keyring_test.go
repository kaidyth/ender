package client

import (
	"testing"

	"github.com/99designs/keyring"
	"github.com/stretchr/testify/assert"
)

func TestKeyring(t *testing.T) {
	ring, err := GetKeyring("default")
	assert.Nil(t, err)

	err = ring.Set(keyring.Item{Key: "test", Data: []byte("com.kaidyth.ender")})
	assert.Nil(t, err)

	d, err := ring.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, "com.kaidyth.ender", string(d.Data))

	err = ring.Remove("test")
	assert.Nil(t, err)
	d, err = ring.Get("test")
	assert.NotNil(t, err)
	assert.Equal(t, keyring.Item{}, d)
}
