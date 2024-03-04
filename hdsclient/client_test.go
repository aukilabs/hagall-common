package hdsclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aukilabs/go-tooling/pkg/logs"
	httpcmn "github.com/aukilabs/hagall-common/http"
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

	t.Run("registration success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		registeredCh := make(chan struct{})
		var clientVerification func(w http.ResponseWriter, r *http.Request)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			vw := httptest.NewRecorder()
			vr := httptest.NewRequest(http.MethodPost, "/registration", nil)

			body, err := io.ReadAll(req.Body)
			require.NoError(t, err)

			var postBody PostServerIn
			err = json.Unmarshal(body, &postBody)
			require.NoError(t, err)

			vr.Header.Set(httpcmn.HeaderHagallIDKey, "0x1")
			vr.Header.Set(httpcmn.HeaderHagallJWTSecretHeaderKey, httpcmn.MakeJWTSecret())
			vr.Header.Set(httpcmn.HeaderHagallRegistrationStateKey, postBody.State)

			clientVerification(vw, vr)

			require.Equal(t, http.StatusOK, vw.Result().StatusCode)
			w.WriteHeader(http.StatusOK)
			registeredCh <- struct{}{}
		}))

		client := NewClient(WithHagallEndpoint("http://test"),
			WithHDSEndpoint(server.URL),
			WithEncoder(json.Marshal),
			WithDecoder(json.Unmarshal),
			WithTransport(http.DefaultTransport))

		clientVerification = client.HandleServerRegistration

		clientCh := make(chan struct{})
		go func() {
			err := client.Pair(ctx, PairIn{
				Endpoint:             server.URL,
				HealthCheckTTL:       1 * time.Minute,
				RegistrationInterval: 100 * time.Millisecond,
				RegistrationRetries:  2,
			})

			require.ErrorIs(t, err, context.Canceled)
			require.Equal(t, RegistrationStatusRegistered, client.GetRegistrationStatus())
			clientCh <- struct{}{}
		}()
		<-registeredCh

		time.Sleep(100 * time.Millisecond)
		cancel()

		<-clientCh
	})

	t.Run("registration success with retries", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		retried := 0
		registeredCh := make(chan struct{})
		var clientVerification func(w http.ResponseWriter, r *http.Request)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			retried++
			if retried == 1 {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				vw := httptest.NewRecorder()
				vr := httptest.NewRequest(http.MethodPost, "/registration", nil)

				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)

				var postBody PostServerIn
				err = json.Unmarshal(body, &postBody)
				require.NoError(t, err)

				vr.Header.Set(httpcmn.HeaderHagallIDKey, "0x1")
				vr.Header.Set(httpcmn.HeaderHagallJWTSecretHeaderKey, httpcmn.MakeJWTSecret())
				vr.Header.Set(httpcmn.HeaderHagallRegistrationStateKey, postBody.State)

				clientVerification(vw, vr)

				require.Equal(t, http.StatusOK, vw.Result().StatusCode)
				w.WriteHeader(http.StatusOK)
				registeredCh <- struct{}{}
			}
		}))

		client := NewClient(WithHagallEndpoint("http://test"),
			WithHDSEndpoint(server.URL),
			WithEncoder(json.Marshal),
			WithDecoder(json.Unmarshal),
			WithTransport(http.DefaultTransport))

		clientVerification = client.HandleServerRegistration

		clientCh := make(chan struct{})
		go func() {
			err := client.Pair(ctx, PairIn{
				Endpoint:             server.URL,
				HealthCheckTTL:       1 * time.Minute,
				RegistrationInterval: 100 * time.Millisecond,
				RegistrationRetries:  2,
			})

			require.ErrorIs(t, err, context.Canceled)
			require.Equal(t, RegistrationStatusRegistered, client.GetRegistrationStatus())
			clientCh <- struct{}{}
		}()
		<-registeredCh

		time.Sleep(100 * time.Millisecond)
		cancel()

		<-clientCh
		require.Equal(t, 2, retried)
	})

	t.Run("registration failed exceeded maximum retries", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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

	t.Run("registration failed - wrong state", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		registeredCh := make(chan struct{})
		var clientVerification func(w http.ResponseWriter, r *http.Request)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			vw := httptest.NewRecorder()
			vr := httptest.NewRequest(http.MethodPost, "/registration", nil)

			vr.Header.Set(httpcmn.HeaderHagallIDKey, "0x1")
			vr.Header.Set(httpcmn.HeaderHagallJWTSecretHeaderKey, httpcmn.MakeJWTSecret())
			vr.Header.Set(httpcmn.HeaderHagallRegistrationStateKey, "wrong-state")

			clientVerification(vw, vr)

			require.Equal(t, http.StatusForbidden, vw.Result().StatusCode)
			w.WriteHeader(http.StatusForbidden)
			registeredCh <- struct{}{}
		}))

		client := NewClient(WithHagallEndpoint("http://test"),
			WithHDSEndpoint(server.URL),
			WithEncoder(json.Marshal),
			WithDecoder(json.Unmarshal),
			WithTransport(http.DefaultTransport))

		clientVerification = client.HandleServerRegistration

		clientCh := make(chan struct{})
		go func() {
			err := client.Pair(ctx, PairIn{
				Endpoint:             server.URL,
				HealthCheckTTL:       1 * time.Minute,
				RegistrationInterval: 100 * time.Millisecond,
				RegistrationRetries:  1,
			})

			require.Error(t, err)
			require.Equal(t, RegistrationStatusFailed, client.GetRegistrationStatus())
			clientCh <- struct{}{}
		}()
		<-registeredCh
		<-clientCh
	})

	t.Run("registration failed - missing secret", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		registeredCh := make(chan struct{})
		var clientVerification func(w http.ResponseWriter, r *http.Request)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			vw := httptest.NewRecorder()
			vr := httptest.NewRequest(http.MethodPost, "/registration", nil)

			body, err := io.ReadAll(req.Body)
			require.NoError(t, err)

			var postBody PostServerIn
			err = json.Unmarshal(body, &postBody)
			require.NoError(t, err)

			vr.Header.Set(httpcmn.HeaderHagallIDKey, "0x1")
			vr.Header.Set(httpcmn.HeaderHagallRegistrationStateKey, postBody.State)

			clientVerification(vw, vr)

			require.Equal(t, http.StatusBadRequest, vw.Result().StatusCode)
			w.WriteHeader(http.StatusForbidden)
			registeredCh <- struct{}{}
		}))

		client := NewClient(WithHagallEndpoint("http://test"),
			WithHDSEndpoint(server.URL),
			WithEncoder(json.Marshal),
			WithDecoder(json.Unmarshal),
			WithTransport(http.DefaultTransport))

		clientVerification = client.HandleServerRegistration

		clientCh := make(chan struct{})
		go func() {
			err := client.Pair(ctx, PairIn{
				Endpoint:             server.URL,
				HealthCheckTTL:       1 * time.Minute,
				RegistrationInterval: 100 * time.Millisecond,
				RegistrationRetries:  1,
			})

			require.Error(t, err)
			require.Equal(t, RegistrationStatusFailed, client.GetRegistrationStatus())
			clientCh <- struct{}{}
		}()
		<-registeredCh
		<-clientCh
	})

	t.Run("registration failed - missing session id", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		registeredCh := make(chan struct{})
		var clientVerification func(w http.ResponseWriter, r *http.Request)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			vw := httptest.NewRecorder()
			vr := httptest.NewRequest(http.MethodPost, "/registration", nil)

			body, err := io.ReadAll(req.Body)
			require.NoError(t, err)

			var postBody PostServerIn
			err = json.Unmarshal(body, &postBody)
			require.NoError(t, err)

			vr.Header.Set(httpcmn.HeaderHagallJWTSecretHeaderKey, httpcmn.MakeJWTSecret())
			vr.Header.Set(httpcmn.HeaderHagallRegistrationStateKey, postBody.State)

			clientVerification(vw, vr)

			require.Equal(t, http.StatusBadRequest, vw.Result().StatusCode)
			w.WriteHeader(http.StatusForbidden)
			registeredCh <- struct{}{}
		}))

		client := NewClient(WithHagallEndpoint("http://test"),
			WithHDSEndpoint(server.URL),
			WithEncoder(json.Marshal),
			WithDecoder(json.Unmarshal),
			WithTransport(http.DefaultTransport))

		clientVerification = client.HandleServerRegistration

		clientCh := make(chan struct{})
		go func() {
			err := client.Pair(ctx, PairIn{
				Endpoint:             server.URL,
				HealthCheckTTL:       1 * time.Minute,
				RegistrationInterval: 100 * time.Millisecond,
				RegistrationRetries:  1,
			})

			require.Error(t, err)
			require.Equal(t, RegistrationStatusFailed, client.GetRegistrationStatus())
			clientCh <- struct{}{}
		}()
		<-registeredCh
		<-clientCh
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
