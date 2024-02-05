package metrics

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTP(t *testing.T) {
	s := httptest.NewServer(HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, err := io.ReadAll(r.Body)
		if err != nil || len(name) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Write([]byte("Hello, "))
		w.Write(name)
	})))
	defer s.Close()

	transport := HTTPTransport(http.DefaultTransport)

	t.Run("no payload is sent returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, s.URL, nil)
		res, err := transport.RoundTrip(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.NotNil(t, res.Body)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("payload sent returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, s.URL, bytes.NewBufferString("Ted"))
		res, err := transport.RoundTrip(req)
		require.NoError(t, err)
		defer res.Body.Close()

		require.NotNil(t, res.Body)
		require.Equal(t, http.StatusOK, res.StatusCode)

		greet, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, "Hello, Ted", string(greet))
	})
}

func TestResponseWriterHijack(t *testing.T) {
	t.Run("underlying writer is a hijacker", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			w := makeResponseWriter(res, func(statusCode, bytes int, err error) {})
			conn, rw, err := w.Hijack()
			require.NotNil(t, conn)
			require.NotNil(t, rw)
			require.NoError(t, err)
			conn.Close()
		}))

		_, err := s.Client().Get(s.URL)
		require.Error(t, err)

	})

	t.Run("underlying writer is not a hijacker", func(t *testing.T) {
		w := makeResponseWriter(httptest.NewRecorder(), func(statusCode, bytes int, err error) {})
		conn, rw, err := w.Hijack()
		require.Nil(t, conn)
		require.Nil(t, rw)
		require.Error(t, err)
	})
}

func TestDefaultPathFormater(t *testing.T) {
	utests := []struct {
		in  string
		out string
	}{
		{in: "/", out: "/"},
		{in: "/hello", out: "/hello"},
		{in: "/hello/", out: "/hello/"},
		{in: "/hello/world", out: "/hello/"},
		{in: "hello", out: "/hello"},
	}

	for _, u := range utests {
		require.Equal(t, u.out, DefaultPathFormater(0, u.in))
	}
}

func BenchmarkDefaultPathFormater(b *testing.B) {
	for n := 0; n < b.N; n++ {
		DefaultPathFormater(0, "/hello/world")
	}
}
