package crypt

import (
	"crypto/sha256"

	"github.com/aukilabs/go-tooling/pkg/errors"
)

type hdsClient interface {
	Secret() string
}

type cacher func(func() ([]byte, error)) ([]byte, error)

// A secret provider using Hagall token from HDS client.
type HagallSecretProvider struct {
	client hdsClient
	cacher
}

// Returns a new HagallSecretProvider.
func NewHagallSecretProvider(client hdsClient) HagallSecretProvider {
	return HagallSecretProvider{
		client: client,
		cacher: keyCacher(client),
	}
}

// GetKey generates a 256-bit key using sha256 hash of Hagall secret with cache.
func (h HagallSecretProvider) GetKey() ([]byte, error) {
	return h.cacher(h.getKey)
}

// getKey generates a 256-bit key using sha256 hash of Hagall secret.
func (h HagallSecretProvider) getKey() ([]byte, error) {
	return sha256hash([]byte(h.client.Secret()))
}

func sha256hash(buf []byte) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(buf))
	if err != nil {
		return nil, errors.New("failed to hash secret").Wrap(err)
	}
	return hash.Sum(nil), nil
}

// keyCacher returns cacher function that caches key generated from fn.
// Cache is invalidated when client secret changes.
func keyCacher(client hdsClient) cacher {
	var cachedSecret string
	var cachedKey []byte
	return func(fn func() ([]byte, error)) ([]byte, error) {
		var err error
		if len(cachedKey) == 0 ||
			cachedSecret != client.Secret() {
			cachedKey, err = fn()
			if err != nil {
				return nil, err
			}
			cachedSecret = client.Secret()
		}
		return cachedKey, nil
	}
}
