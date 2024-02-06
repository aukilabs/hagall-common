package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/aukilabs/go-tooling/pkg/errors"
)

const (
	// default nonce size for AES-GCM
	nonceSize = 12
)

// Encrypt encrypts buf using key with AES-GCM AEAD mode,
// nonce is prepend to encrypted text and result is hex encoded
func Encrypt(buf []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to initialize cipher").
			WithTag("key_length", len(key)).
			Wrap(err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("failed to initialize nonce").
			WithTag("nonce_size", nonceSize).
			Wrap(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to initialize encryptor").
			Wrap(err)
	}

	encrypted := aesGCM.Seal(nil, nonce, buf, nil)

	// prepend nonce into payload
	encrypted = append(nonce, encrypted...)
	return encrypted, nil
}

// Decrypt decrypts buf using key with AES-GCM AEAD mode
func Decrypt(buf []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to initialize cipher").Wrap(err)
	}

	if len(buf) < nonceSize {
		return nil, errors.New("decrypt failed: invalid payload")
	}

	// extract nonce from payload
	nonce := buf[0:nonceSize]
	encrypted := buf[nonceSize:]

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to initialize decryptor").Wrap(err)
	}

	decrypted, err := aesGCM.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt payload").Wrap(err)
	}

	return decrypted, nil
}
