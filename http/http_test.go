package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aukilabs/go-tooling/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestHTTPError(t *testing.T) {
	// observe log level from stdout
	t.Run("log info on status code < 500", func(t *testing.T) {
		rec := httptest.NewRecorder()

		HTTPError(rec, http.StatusBadRequest, nil)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("log error on status code >= 500", func(t *testing.T) {
		rec := httptest.NewRecorder()

		HTTPError(rec, http.StatusInternalServerError, nil)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("return error message if defined", func(t *testing.T) {
		rec := httptest.NewRecorder()

		HTTPError(rec, http.StatusInternalServerError, errors.New("test error").Wrap(ErrDuplicatedWalletAddress))
		require.Equal(t, http.StatusInternalServerError, rec.Code)

		body, err := io.ReadAll(rec.Result().Body)
		require.NoError(t, err)
		require.Equal(t, "Wallet already registered for your endpoint or another endpoint", strings.TrimSpace(string(body)))
	})

	t.Run("return empty if undefined", func(t *testing.T) {
		rec := httptest.NewRecorder()

		HTTPError(rec, http.StatusInternalServerError, errors.New("test error"))
		require.Equal(t, http.StatusInternalServerError, rec.Code)

		body, err := io.ReadAll(rec.Result().Body)
		require.NoError(t, err)
		require.Empty(t, strings.TrimSpace(string(body)))
	})
}
