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
	http.Error(w, "", http.StatusNotFound)
}

func Unauthorized(w http.ResponseWriter, err error) {
	HTTPError(w, http.StatusUnauthorized, err)
}

func InternalServerError(w http.ResponseWriter, err error) {
	if err != nil {
		logs.WithTag("code", http.StatusInternalServerError).
			WithTag("error", err.Error()).
			Error(errors.New("internal error handling http request").Wrap(err))
	}

	http.Error(w, GetErrorMessage(err), http.StatusInternalServerError)
}

func NotImplemented(w http.ResponseWriter) {
	http.Error(w, "", http.StatusNotImplemented)
}

func PaymentRequired(w http.ResponseWriter, err error) {
	HTTPError(w, http.StatusPaymentRequired, err)
}

func Conflict(w http.ResponseWriter, err error) {
	HTTPError(w, http.StatusConflict, err)
}

func HTTPError(w http.ResponseWriter, code int, err error) {
	if err != nil {
		logs.WithTag("code", code).
			WithTag("error", err.Error()).
			Warn("http request returned an error")
	}

	http.Error(w, GetErrorMessage(err), code)
}

func NormalizeEndpoint(v string) string {
	v = strings.TrimSpace(v)
	return strings.TrimRight(v, "/")
}
