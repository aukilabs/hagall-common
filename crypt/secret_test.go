package crypt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHagallSecretProvider(t *testing.T) {
	client := &mockHdsClient{secret: "my-secret"}
	p := NewHagallSecretProvider(client)

	key, err := p.GetKey()
	require.NoError(t, err)
	require.NotEmpty(t, key)

	cacheKey, err := p.GetKey()
	require.NoError(t, err)
	require.Equal(t, cacheKey, key)

	client.secret = "my-new-secret"
	newKey, err := p.GetKey()
	require.NoError(t, err)
	require.NotEqual(t, newKey, key)

}

type mockHdsClient struct {
	secret string
}

func (m mockHdsClient) Secret() string {
	return m.secret
}
