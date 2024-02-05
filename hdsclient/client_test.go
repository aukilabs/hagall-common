package hdsclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httpcmn "github.com/aukilabs/hagall-common/http"
	"github.com/aukilabs/hagall-common/logs"
	"github.com/stretchr/testify/require"
)

func TestClientHandleRegistration(t *testing.T) {
	setupTestLog(t)

	t.Run("handling non self initiated registration return a 403", func(t *testing.T) {
		c := NewClient()
		res := httptest.NewRecorder()
		c.HandleServerRegistration(res, httptest.NewRequest(http.MethodPost, "/", nil))
		require.Equal(t, http.StatusForbidden, res.Code)
	})

	t.Run("handling non initiated registration return a 403", func(t *testing.T) {
		c := NewClient()
		c.setRegistrationState("hello")

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(httpcmn.HeaderHagallRegistrationStateKey, "bye")

		res := httptest.NewRecorder()

		c.HandleServerRegistration(res, req)
		require.Equal(t, http.StatusForbidden, res.Code)

	})
}

func TestClientRegistration(t *testing.T) {
	setupTestLog(t)

	t.Run("registration success - pending verification", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		registeredCh := make(chan struct{})
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			registeredCh <- struct{}{}
		}))

		clientCh := make(chan struct{})
		go func() {
			client := NewClient(WithHagallEndpoint("http://test"),
				WithHDSEndpoint(server.URL),
				WithEncoder(json.Marshal),
				WithDecoder(json.Unmarshal),
				WithTransport(http.DefaultTransport))
			err := client.Pair(ctx, PairIn{
				Endpoint:             server.URL,
				HealthCheckTTL:       1 * time.Minute,
				RegistrationInterval: 100 * time.Millisecond,
				RegistrationRetries:  1,
			})
			require.ErrorIs(t, err, context.Canceled)
			require.Equal(t, RegistrationStatusPendingVerification, client.GetRegistrationStatus())
			clientCh <- struct{}{}
		}()
		<-registeredCh

		time.Sleep(100 * time.Millisecond)
		cancel()

		<-clientCh
	})

	t.Run("registration success with retries", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		retried := 0
		registeredCh := make(chan struct{})
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			retried++
			if retried == 1 {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				registeredCh <- struct{}{}
			}
		}))

		clientCh := make(chan struct{})
		go func() {
			client := NewClient(WithHagallEndpoint("http://test"),
				WithHDSEndpoint(server.URL),
				WithEncoder(json.Marshal),
				WithDecoder(json.Unmarshal),
				WithTransport(http.DefaultTransport))
			err := client.Pair(ctx, PairIn{
				Endpoint:             server.URL,
				HealthCheckTTL:       1 * time.Minute,
				RegistrationInterval: 100 * time.Millisecond,
				RegistrationRetries:  2,
			})

			require.ErrorIs(t, err, context.Canceled)
			require.Equal(t, RegistrationStatusPendingVerification, client.GetRegistrationStatus())
			clientCh <- struct{}{}
		}()
		<-registeredCh

		time.Sleep(100 * time.Millisecond)
		cancel()

		<-clientCh
		require.Equal(t, 2, retried)
	})

	t.Run("registration failed exceeded maximum retries", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		retried := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			retried++
			w.WriteHeader(http.StatusInternalServerError)
		}))

		clientCh := make(chan struct{})
		go func() {
			client := NewClient(WithHagallEndpoint("http://test"),
				WithHDSEndpoint(server.URL),
				WithEncoder(json.Marshal),
				WithDecoder(json.Unmarshal),
				WithTransport(http.DefaultTransport))
			err := client.Pair(ctx, PairIn{
				Endpoint:             server.URL,
				HealthCheckTTL:       1 * time.Minute,
				RegistrationInterval: 100 * time.Millisecond,
				RegistrationRetries:  3,
			})

			require.NotErrorIs(t, err, context.Canceled)
			require.Equal(t, RegistrationStatusFailed, client.GetRegistrationStatus())
			clientCh <- struct{}{}
		}()

		<-clientCh
		require.Equal(t, 3, retried)
	})

	t.Run("registration retries with incremental delay", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		retried := 0
		lastRetry := time.Now()
		registrationInterval := 100 * time.Millisecond
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			expected := lastRetry.Add(time.Duration(((retried) * int(registrationInterval))))
			require.GreaterOrEqual(t, time.Now(), expected)
			lastRetry = time.Now()
			retried++
			w.WriteHeader(http.StatusInternalServerError)
		}))

		clientCh := make(chan struct{})
		go func() {
			client := NewClient(WithHagallEndpoint("http://test"),
				WithHDSEndpoint(server.URL),
				WithEncoder(json.Marshal),
				WithDecoder(json.Unmarshal),
				WithTransport(http.DefaultTransport))
			err := client.Pair(ctx, PairIn{
				Endpoint:             server.URL,
				HealthCheckTTL:       1 * time.Minute,
				RegistrationInterval: registrationInterval,
				RegistrationRetries:  3,
			})

			require.NotErrorIs(t, err, context.Canceled)
			require.Equal(t, RegistrationStatusFailed, client.GetRegistrationStatus())
			clientCh <- struct{}{}
		}()

		<-clientCh
		require.Equal(t, 3, retried)
	})
}

func TestUnpair(t *testing.T) {
	t.Run("unpair", func(t *testing.T) {
		var unpaired bool
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method == http.MethodDelete {
				unpaired = true
			}
		}))

		client := NewClient(WithHagallEndpoint("http://test"),
			WithHDSEndpoint(server.URL),
			WithEncoder(json.Marshal),
			WithDecoder(json.Unmarshal),
			WithTransport(http.DefaultTransport))
		client.setRegistrationStatus(RegistrationStatusRegistered)

		err := client.Unpair()
		require.NoError(t, err)

		require.True(t, unpaired)
	})

	t.Run("unpair on panic", func(t *testing.T) {
		registeredCh := make(chan struct{})
		var unpaired bool
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method == http.MethodPost {
				w.WriteHeader(http.StatusOK)
				registeredCh <- struct{}{}
				return
			}

			if req.Method == http.MethodDelete {
				unpaired = true
			}
		}))

		client := NewClient(WithHagallEndpoint("http://test"),
			WithHDSEndpoint(server.URL),
			WithEncoder(json.Marshal),
			WithDecoder(json.Unmarshal),
			WithTransport(http.DefaultTransport))
		client.setRegistrationStatus(RegistrationStatusRegistered)

		defer func() {
			// unpair on panic
			if r := recover(); r != nil {
				err := client.Unpair()
				require.NoError(t, err)
				require.True(t, unpaired)
			}
		}()
		panic("test panic")
	})
}

func setupTestLog(tb testing.TB) {
	logs.SetLogger(func(e logs.Entry) { tb.Log(e) })
	logs.Encoder = func(v any) ([]byte, error) {
		return json.MarshalIndent(v, "", "  ")
	}
}
