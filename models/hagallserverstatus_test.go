package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHagallServerStatusFromString(t *testing.T) {
	t.Run("valid statuses", func(t *testing.T) {
		statuses := []string{"offline", "online", "unhealthy"}
		for _, s := range statuses {
			_, err := HagallServerStatusFromString(s)
			require.NoError(t, err)
		}
	})
	t.Run("invalid statuses", func(t *testing.T) {
		statuses := []string{"-1", "a", ""}
		for _, s := range statuses {
			_, err := HagallServerStatusFromString(s)
			require.Error(t, err)
		}
	})
}
