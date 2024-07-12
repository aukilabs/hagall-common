package crypt

import (
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSign(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	publicKey := crypto.FromECDSAPub(&privateKey.PublicKey)

	t.Run("succesful sign", func(t *testing.T) {
		msg := uuid.NewString()

		signature, err := Sign(privateKey, msg)
		require.NoError(t, err)
		require.True(t, strings.HasPrefix(signature, "0x"))
		signatureBytes := common.FromHex(signature)
		require.Equal(t, 65, len(signatureBytes))

		hash := crypto.Keccak256([]byte(msg))
		require.Equal(t, 32, len(hash))
		correct := crypto.VerifySignature(publicKey, hash, signatureBytes[:64])
		require.True(t, correct)
	})

	t.Run("different public key", func(t *testing.T) {
		msg := uuid.NewString()

		signature, err := Sign(privateKey, msg)
		require.NoError(t, err)
		signatureBytes := common.FromHex(signature)

		privateKey2, err := crypto.GenerateKey()
		require.NoError(t, err)
		publicKey2 := crypto.FromECDSAPub(&privateKey2.PublicKey)

		hash := crypto.Keccak256([]byte(msg))
		correct := crypto.VerifySignature(publicKey2, hash, signatureBytes[:64])
		require.False(t, correct)
	})
}

func TestSignWithTimestamp(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	publicKey := crypto.FromECDSAPub(&privateKey.PublicKey)

	t.Run("succesful sign", func(t *testing.T) {
		msg := uuid.NewString()

		signature, timestamp, err := SignWithTimestamp(privateKey, msg)
		require.NoError(t, err)
		require.True(t, strings.HasPrefix(signature, "0x"))
		signatureBytes := common.FromHex(signature)
		require.Equal(t, 65, len(signatureBytes))

		msgWithTimestamp := msg + timestamp
		hash := crypto.Keccak256([]byte(msgWithTimestamp))
		correct := crypto.VerifySignature(publicKey, hash, signatureBytes[:64])
		require.True(t, correct)
	})

	t.Run("different public key", func(t *testing.T) {
		msg := uuid.NewString()

		signature, timestamp, err := SignWithTimestamp(privateKey, msg)
		require.NoError(t, err)
		signatureBytes := common.FromHex(signature)
		require.Equal(t, 65, len(signatureBytes))

		privateKey2, err := crypto.GenerateKey()
		require.NoError(t, err)
		publicKey2 := crypto.FromECDSAPub(&privateKey2.PublicKey)

		msgWithTimestamp := msg + timestamp
		hash := crypto.Keccak256([]byte(msgWithTimestamp))
		correct := crypto.VerifySignature(publicKey2, hash, signatureBytes[:64])
		require.False(t, correct)
	})
}

func TestValidateSignedMessage(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	t.Run("empty message and signature", func(t *testing.T) {
		_, err := ValidateSignedMessage("", "")
		require.Error(t, err)
	})

	t.Run("empty signature", func(t *testing.T) {
		_, err := ValidateSignedMessage("message", "")
		require.Error(t, err)
	})

	t.Run("empty message", func(t *testing.T) {
		msg := uuid.NewString()
		signature, err := Sign(privateKey, msg)
		require.NoError(t, err)

		_, err = ValidateSignedMessage("", signature)
		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		msg := uuid.NewString()
		signature, err := Sign(privateKey, msg)
		require.NoError(t, err)

		recoveredAddr, err := ValidateSignedMessage(msg, signature)
		require.NoError(t, err)
		require.Equal(t, address, recoveredAddr)
	})

	t.Run("non-matching msg", func(t *testing.T) {
		msg := uuid.NewString()
		signature, err := Sign(privateKey, msg)
		require.NoError(t, err)

		msg2 := uuid.NewString()

		recoveredAddr, err := ValidateSignedMessage(msg2, signature)
		require.NoError(t, err) // validate doesn't actually verify that signature is correct for the message
		require.NotEqual(t, address, recoveredAddr)
	})

	t.Run("different timestamps", func(t *testing.T) {
		now := time.Now()
		later := now.Add(time.Hour)
		url := uuid.NewString()
		msg := url + now.Format(time.RFC3339Nano)
		msg2 := url + later.Format(time.RFC3339Nano)
		signature, err := Sign(privateKey, msg)
		require.NoError(t, err)

		recoveredAddr, err := ValidateSignedMessage(msg2, signature)
		require.NoError(t, err) // validate doesn't actually verify that signature is correct for the message
		require.NotEqual(t, address, recoveredAddr)
	})
}

func TestVerifySignedMessage(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	publicKey := &privateKey.PublicKey

	t.Run("empty signature", func(t *testing.T) {
		_, err := VerifySignedMessage(publicKey, "message", "")
		require.Error(t, err)
	})

	t.Run("empty message", func(t *testing.T) {
		msg := uuid.NewString()
		signature, err := Sign(privateKey, msg)
		require.NoError(t, err)

		_, err = VerifySignedMessage(publicKey, "", signature)
		require.Error(t, err)
	})

	t.Run("nil public key", func(t *testing.T) {
		msg := uuid.NewString()
		signature, err := Sign(privateKey, msg)
		require.NoError(t, err)

		_, err = VerifySignedMessage(nil, msg, signature)
		require.Error(t, err)
	})

	t.Run("invalid signature", func(t *testing.T) {
		_, err = VerifySignedMessage(publicKey, uuid.NewString(), "0x191675bf9f484d3ea93086cde7ed2ee4")
		require.Error(t, err)
	})

	t.Run("non-matching msg", func(t *testing.T) {
		msg := uuid.NewString()
		signature, err := Sign(privateKey, msg)
		require.NoError(t, err)

		msg2 := uuid.NewString()

		_, err = VerifySignedMessage(publicKey, msg2, signature)
		require.Error(t, err)
	})

	t.Run("non-matching public key", func(t *testing.T) {
		msg := uuid.NewString()
		signature, err := Sign(privateKey, msg)
		require.NoError(t, err)

		privateKey2, err := crypto.GenerateKey()
		require.NoError(t, err)
		publicKey2 := &privateKey2.PublicKey

		_, err = VerifySignedMessage(publicKey2, msg, signature)
		require.Error(t, err)
	})

	t.Run("different timestamps", func(t *testing.T) {
		now := time.Now()
		later := now.Add(time.Hour)
		url := uuid.NewString()
		msg := url + now.Format(time.RFC3339Nano)
		msg2 := url + later.Format(time.RFC3339Nano)
		signature, err := Sign(privateKey, msg)
		require.NoError(t, err)

		recoveredAddr, err := VerifySignedMessage(publicKey, msg2, signature)
		require.Error(t, err)
		require.Equal(t, common.Address{}, recoveredAddr)
	})

	t.Run("success", func(t *testing.T) {
		msg := uuid.NewString()
		signature, err := Sign(privateKey, msg)
		require.NoError(t, err)

		recoveredAddr, err := VerifySignedMessage(publicKey, msg, signature)
		require.NoError(t, err)
		require.Equal(t, address, recoveredAddr)
	})
}
