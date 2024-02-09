package crypt

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleWithEncryption(t *testing.T) {
	testBody := "super secret content"
	provider := mockProvider{}

	handle := HandleWithEncryption(provider, http.HandlerFunc((mockHandler(http.StatusOK, testBody))))
	rec := httptest.NewRecorder()

	handle(rec, &http.Request{})
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)

	key, err := provider.GetKey()
	require.NoError(t, err)

	decrypted, err := Decrypt(rec.Body.Bytes(), key)
	require.NoError(t, err)
	require.Equal(t, testBody, string(decrypted))
}

type mockProvider struct{}

func (m mockProvider) GetKey() ([]byte, error) {
	return sha256hash([]byte("super-secret"))
}

func mockHandler(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if len(body) > 0 {
			w.Write([]byte(body))
		}
	}
}
