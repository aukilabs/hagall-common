package crypt

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestSignature(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	testMessage := "secret message"
	signature, err := Sign(privateKey, testMessage)
	require.NoError(t, err)

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	walletAddr, err := ValidateSignedMessage(testMessage, signature)
	require.NoError(t, err)
	require.Equal(t, address, walletAddr)
}
