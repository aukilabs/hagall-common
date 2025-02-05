package http

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTokenFinders(t *testing.T) {
	utests := []struct {
		scenario      string
		finder        tokenFinder
		request       *http.Request
		expectedToken string
	}{
		{
			scenario: "token is retrieved from bearer authorization header",
			finder:   tokenFromHeader,
			request: &http.Request{
				Header: http.Header{
					"Authorization": []string{"Bearer tedxyz"},
				},
			},
			expectedToken: "tedxyz",
		},
		{
			scenario: "token is not retrieved from bearer authorization header",
			finder:   tokenFromHeader,
			request: &http.Request{
				Header: http.Header{
					"Authorization": []string{"Basic tedxyz:qqq"},
				},
			},
			expectedToken: "",
		},
		{
			scenario: "token is retrieved from url parameters",
			finder:   tokenFromQuery,
			request: &http.Request{
				URL: &url.URL{
					Scheme:   "http",
					Host:     "wingchun.ted",
					RawQuery: "access_token=tedxyz",
				},
			},
			expectedToken: "tedxyz",
		},
		{
			scenario: "token is not retrieved from url parameters",
			finder:   tokenFromQuery,
			request: &http.Request{
				URL: &url.URL{
					Scheme: "http",
					Host:   "wingchun.ted",
				},
			},
			expectedToken: "",
		},
		{
			scenario: "token is retrieved from cookie",
			finder:   tokenFromCookie,
			request: &http.Request{
				Header: http.Header{
					"Cookie": []string{
						(&http.Cookie{Name: "access_token", Value: "tedxyz"}).String(),
					},
				},
			},
			expectedToken: "tedxyz",
		},
		{
			scenario:      "token is not retrieved from cookie",
			finder:        tokenFromCookie,
			request:       &http.Request{},
			expectedToken: "",
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			token := u.finder(u.request)
			require.Equal(t, u.expectedToken, token)
		})
	}
}

func TestVerifyHagallUserAccessToken(t *testing.T) {
	secret := MakeJWTSecret()

	t.Run("verification from user token", func(t *testing.T) {
		token, err := GenerateHagallUserAccessToken(
			"0x0",
			secret,
			time.Minute,
		)
		require.NoError(t, err)

		err = VerifyHagallUserAccessToken(token, secret)
		require.NoError(t, err)
	})

	t.Run("verification from legacy user token", func(t *testing.T) {
		token, err := generateLegacyHagallUserAccessToken(
			"0x0",
			secret,
			time.Minute,
		)
		require.NoError(t, err)

		err = VerifyHagallUserAccessToken(token, secret)
		require.NoError(t, err)
	})

	t.Run("allow tokens used <10 seconds before issuance", func(t *testing.T) {
		now := time.Now()
		tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, HagallUserClaim{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "HDS",
				Subject:   "",
				IssuedAt:  jwt.NewNumericDate(now.Add(9 * time.Second)),
				ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Minute)),
				ID:        uuid.NewString(),
			},
			AppKey: "0x0",
		})

		token, err := tokenWithClaims.SignedString([]byte(secret))
		require.NoError(t, err)

		err = VerifyHagallUserAccessToken(token, secret)
		require.NoError(t, err)
	})

	t.Run("disallow tokens used >10 seconds before issuance", func(t *testing.T) {
		now := time.Now()
		tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, HagallUserClaim{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "HDS",
				Subject:   "",
				IssuedAt:  jwt.NewNumericDate(now.Add(11 * time.Second)),
				ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Minute)),
				ID:        uuid.NewString(),
			},
			AppKey: "0x0",
		})

		token, err := tokenWithClaims.SignedString([]byte(secret))
		require.NoError(t, err)

		err = VerifyHagallUserAccessToken(token, secret)
		require.Error(t, err)
		require.Contains(t, err.Error(), "token used before issued")
	})

	t.Run("legacy verification from user token", func(t *testing.T) {
		token, err := GenerateHagallUserAccessToken(
			"0x0",
			secret,
			time.Minute,
		)
		require.NoError(t, err)

		err = verifyLegacyHagallUserAccessToken(token, secret)
		require.NoError(t, err)
	})
}

func TestGetAppKeyFromHagallUserToken(t *testing.T) {
	secret := MakeJWTSecret()

	token, err := GenerateHagallUserAccessToken(
		"0xTED",
		secret,
		time.Minute,
	)
	require.NoError(t, err)

	appKey := GetAppKeyFromHagallUserToken(token)
	require.Equal(t, "0xTED", appKey)
}

func TestGetAppKeyFromHTTPRequest(t *testing.T) {
	utests := []struct {
		scenario       string
		request        *http.Request
		expectedAppKey string
	}{
		{
			scenario: "no authorization",
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", nil)
				return r
			}(),
		},
		{
			scenario: "app key in basic authorization",
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", nil)
				r.SetBasicAuth("tedx", "42")
				return r
			}(),
			expectedAppKey: "tedx",
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			appKey := GetAppKeyFromHTTPRequest(u.request)
			require.Equal(t, u.expectedAppKey, appKey)
		})
	}
}

func generateLegacyHagallUserAccessToken(appKey, secret string, ttl time.Duration) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "HDS",
		Subject:   "",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		ID:        uuid.NewString(),
	})

	return token.SignedString([]byte(secret))
}

func verifyLegacyHagallUserAccessToken(token, secret string) error {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	return err
}
