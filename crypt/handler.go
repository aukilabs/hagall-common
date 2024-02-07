package crypt

import (
	"bytes"
	"net/http"

	"github.com/aukilabs/go-tooling/pkg/errors"
	"github.com/aukilabs/go-tooling/pkg/logs"
)

const (
	acceptEncodingHeader = "Accept-Encoding"
)

type secretProvider interface {
	GetKey() ([]byte, error)
}

// HandleWithEncryption returns a http handler function that encrypts content
// returned by handler using key provided by secretProvider.
func HandleWithEncryption(provider secretProvider, handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseWriter := &responseEncrypter{
			secret: provider,
			writer: w,
		}

		removeCompression(r.Header)

		handler.ServeHTTP(responseWriter, r)

		if responseWriter.statusCode != http.StatusOK {
			w.WriteHeader(responseWriter.statusCode)
			w.Write(responseWriter.buf.Bytes())
			return
		}

		if responseWriter.buf.Len() > 0 {
			responseWriter.encryptResponse()
		}
	}
}

// responseEncrypter implements http.ResponseWriter interface to intercept
// plaintext response from upstream http.Handler. Outputs encrypted response.
type responseEncrypter struct {
	secret     secretProvider
	buf        bytes.Buffer
	writer     http.ResponseWriter
	statusCode int
}

// Returns underlying response writer header.
func (w responseEncrypter) Header() http.Header {
	return w.writer.Header()
}

// WriteHeader stores upstream statusCode.
func (w *responseEncrypter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// Write stores upstream content in buffer.
func (w *responseEncrypter) Write(buf []byte) (int, error) {
	return w.buf.Write(buf)
}

// encryptResponse encrypts buffer and outputs to underlying response writer.
func (w responseEncrypter) encryptResponse() {
	key, err := w.secret.GetKey()
	if err != nil {
		logs.Error(errors.New("failed getting key").Wrap(err))
		w.writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return status 403 when iv & key are empty (hagall unregistered)
	if len(key) == 0 {
		w.writer.WriteHeader(http.StatusForbidden)
		return
	}

	enc, err := Encrypt(w.buf.Bytes(), key)
	if err != nil {
		logs.Error(errors.New("failed encrypting content").Wrap(err))
		w.writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.writer.WriteHeader(w.statusCode)

	_, err = w.writer.Write(enc)
	if err != nil {
		logs.Error(errors.New("failed writing response").Wrap(err))
		return
	}
}

// removeCompression removes "gzip" from http Accept-Encoding request header.
func removeCompression(header http.Header) {
	acceptEncoding := header.Values(acceptEncodingHeader)

	header.Del(acceptEncodingHeader)

	for _, enc := range acceptEncoding {
		if enc == "gzip" {
			continue
		}
		header.Add(acceptEncodingHeader, enc)
	}
}
