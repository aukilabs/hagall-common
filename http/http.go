package http

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/aukilabs/hagall-common/errors"
	"github.com/aukilabs/hagall-common/logs"
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

// Normalize IP address in different format
func NormalizeIP(in string) (net.IP, error) {
	ipStr := strings.TrimSpace(in)

	// get first address if ipStr contains comma
	if strings.Index(in, ",") > 0 {
		ipStr = strings.TrimSpace(strings.Split(in, ",")[0])
	}

	// split ipStr if it contains IP & port information
	ip, _, err := net.SplitHostPort(ipStr)
	if err == nil {
		ipStr = ip
	}

	netIP := net.ParseIP(ipStr)
	if netIP == nil {
		return net.IP{}, errors.New("error parsing ip")
	}

	// ParseIP stores ipv4 in 4in6 format (64 bits) , pgconn serialize it into ipv6 format i.e '1.0.0.1' -> '::ffff:1.0.0.1'
	// code below convert netIP to into pure ipv4 if it matches ipv4 length (32 bits)
	ip4 := netIP.To4()
	if len(ip4) == net.IPv4len {
		netIP = ip4
	}

	return netIP, nil
}
