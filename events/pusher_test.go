package events

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aukilabs/hagall-common/logs"
	"github.com/stretchr/testify/require"
)

func TestPusher(t *testing.T) {
	initLogs(t)

	var count int
	resetCount := func() {
		count = 0
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var body eventPayload
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		count += len(body.Events)
	}))
	defer s.Close()

	t.Run("events are sent when batch is full", func(t *testing.T) {
		defer resetCount()

		l := Pusher{
			Endpoint:      s.URL,
			BatchSize:     2,
			FlushInterval: time.Minute,
		}
		l.Start()
		l.NewEvent("hello")
		l.NewEvent("world")

		time.Sleep(time.Millisecond * 10)
		l.Close()
		require.Equal(t, 2, count)
	})

	t.Run("events are sent when batch is full", func(t *testing.T) {
		defer resetCount()

		l := Pusher{
			Endpoint:      s.URL,
			FlushInterval: time.Millisecond * 10,
		}
		l.Start()
		l.NewEvent("hello")
		l.NewEvent("world")

		time.Sleep(time.Millisecond * 15)
		l.Close()
		require.Equal(t, 2, count)
	})

	t.Run("events are sent when logger is closed", func(t *testing.T) {
		defer resetCount()

		l := Pusher{
			Endpoint: s.URL,
		}
		l.Start()
		l.NewEvent("hello")
		l.NewEvent("world")

		time.Sleep(time.Millisecond * 10)
		l.Close()
		require.Equal(t, 2, count)
	})
}

func initLogs(t *testing.T) {
	logs.Encoder = func(v any) ([]byte, error) {
		return json.MarshalIndent(v, "", "  ")
	}

	logs.SetLogger(func(e logs.Entry) {
		t.Log(e)
	})
}
