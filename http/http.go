package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aukilabs/go-tooling/pkg/errors"
	"github.com/aukilabs/go-tooling/pkg/logs"
)

func OK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func OKWithJSON(w http.ResponseWriter, out interface{}) {
	body, err := json.Marshal(out)
	if err != nil {
		InternalServerError(w, errors.New("encoding response failed").Wrap(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func NotModified(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotModified)
}

func MethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "", http.StatusMethodNotAllowed)
}

func BadRequest(w http.ResponseWriter, err error) {
	HTTPError(w, http.StatusBadRequest, err)
}

func Forbidden(w http.ResponseWriter, err error) {
	HTTPError(w, http.StatusForbidden, err)
}

func NotFound(w http.ResponseWriter) {
	HTTPError(w, http.StatusNotFound, nil)
}

func Unauthorized(w http.ResponseWriter, err error) {
	HTTPError(w, http.StatusUnauthorized, err)
}

func InternalServerError(w http.ResponseWriter, err error) {
	HTTPError(w, http.StatusInternalServerError, err)
}

func NotImplemented(w http.ResponseWriter) {
	HTTPError(w, http.StatusNotImplemented, nil)
}

func PaymentRequired(w http.ResponseWriter, err error) {
	HTTPError(w, http.StatusPaymentRequired, err)
}

func Conflict(w http.ResponseWriter, err error) {
	HTTPError(w, http.StatusConflict, err)
}

func HTTPError(w http.ResponseWriter, code int, err error) {
	logger := logs.WithTag("code", code)
	if err != nil {
		logger = logger.WithTag("error", err)
	}

	if code >= 500 {
		logger.Error(errors.New("http request returned an error"))
	} else {
		logger.Info("http request returned an error")
	}

	http.Error(w, GetErrorMessage(err), code)
}

func NormalizeEndpoint(v string) string {
	v = strings.TrimSpace(v)
	return strings.TrimRight(v, "/")
}
