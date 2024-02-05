package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeIP(t *testing.T) {
	t.Run("ipv4", func(t *testing.T) {
		ip, err := NormalizeIP("1.1.1.1")
		require.NoError(t, err)
		require.Equal(t, "1.1.1.1", ip.String())
	})
	t.Run("ipv4 with port", func(t *testing.T) {
		ip, err := NormalizeIP("1.1.1.1:1234")
		require.NoError(t, err)
		require.Equal(t, "1.1.1.1", ip.String())
	})
	t.Run("multi ipv4 with comma separator", func(t *testing.T) {
		ip, err := NormalizeIP("1.1.1.1, 2.2.2.2")
		require.NoError(t, err)
		require.Equal(t, "1.1.1.1", ip.String())
	})
	t.Run("ipv6", func(t *testing.T) {
		ip, err := NormalizeIP("2001:d08::1")
		require.NoError(t, err)
		require.Equal(t, "2001:d08::1", ip.String())
	})
	t.Run("ipv6 localhost", func(t *testing.T) {
		ip, err := NormalizeIP("::1")
		require.NoError(t, err)
		require.Equal(t, "::1", ip.String())
	})
	t.Run("ipv6 with port", func(t *testing.T) {
		ip, err := NormalizeIP("[2001:d08::1]:1234")
		require.NoError(t, err)
		require.Equal(t, "2001:d08::1", ip.String())
	})
	t.Run("invalid ipv4", func(t *testing.T) {
		_, err := NormalizeIP("1.2:3456")
		require.Error(t, err)
	})
	t.Run("invalid ipv6", func(t *testing.T) {
		_, err := NormalizeIP("2001:d08::1::2")
		require.Error(t, err)
	})
}
