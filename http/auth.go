package http

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/aukilabs/hagall-common/errors"
	"github.com/aukilabs/hagall-common/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func MakeJWTSecret() string {
	return base64.RawURLEncoding.EncodeToString([]byte(uuid.NewString()))
}

func MakeAuthorizationHeader(token string) string {
	return "Bearer " + token
}

func SignIdentity(endpoint, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"endpoint": endpoint})
	return token.SignedString([]byte(secret))
}

func VerifyHagallUserAccessToken(token, secret string) error {
	var claims models.HagallUserClaim

	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		// Further validations like expiration are validated by the jwt package.
		return []byte(secret), nil
	})
	if err != nil {
		var validationError *jwt.ValidationError
		if errors.As(err, &validationError) {
			if validationError.Inner != nil {
				return errors.New("parse token error").
					WithTag("jwt_error_flags", validationError.Errors).
					Wrap(validationError.Inner)
			} else {
				return errors.New("parse token error").
					WithTag("jwt_error_flags", validationError.Errors).
					Wrap(err)
			}
		}
	}
	return err
}

// GenerateHagallUserAccessToken generate a Hagall user access token with the
// given secret.
func GenerateHagallUserAccessToken(appKey, secret string, ttl time.Duration) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.HagallUserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "HDS",
			Subject:   "",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			ID:        uuid.NewString(),
		},
		AppKey: appKey,
	})

	return token.SignedString([]byte(secret))
}

func GetAppKeyFromHTTPRequest(r *http.Request) string {
	switch auth := r.Header.Get("Authorization"); {
	case strings.HasPrefix(auth, "Basic"):
		key, _, _ := r.BasicAuth()
		return key

	default:
		return ""
	}
}

// Parses the Hagall user token and returns the app key.
func GetAppKeyFromHagallUserToken(token string) string {
	var claims models.HagallUserClaim
	jwt.ParseWithClaims(token, &claims, nil)
	return claims.AppKey
}

// Returns the Hagall user token from a http request.
func GetUserTokenFromHTTPRequest(r *http.Request) string {
	var token string
	for _, finder := range []tokenFinder{
		tokenFromHeader,
		tokenFromQuery,
		tokenFromCookie,
	} {
		if token = finder(r); token != "" {
			return token
		}
	}
	return ""
}

type tokenFinder func(*http.Request) string

func tokenFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if !strings.HasPrefix(bearer, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(bearer, "Bearer ")
}

func tokenFromQuery(r *http.Request) string {
	return r.URL.Query().Get("access_token")
}

func tokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}
