package crypt

import (
	"crypto/ecdsa"
	"time"

	"github.com/aukilabs/go-tooling/pkg/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Sign a message with a private key, return signature in 0x string format
func Sign(privateKey *ecdsa.PrivateKey, message string) (string, error) {
	signature, err := crypto.Sign(crypto.Keccak256Hash([]byte(message)).Bytes(), privateKey)
	if err != nil {
		return "", errors.New("signing message failed").
			Wrap(err)
	}

	return hexutil.Encode(signature), nil
}

// Sign a message with a private key, return signature in 0x string format
func SignWithTimestamp(privateKey *ecdsa.PrivateKey, message string) (string, string, error) {
	timestamp := time.Now().Format(time.RFC3339)
	signature, err := Sign(privateKey, message+timestamp)
	if err != nil {
		return "", "", err
	}

	return signature, timestamp, nil
}

// ValidateSignedMessage validates a signed message and returning wallet address from signature
func ValidateSignedMessage(message string, signature string) (common.Address, error) {
	if message == "" {
		return common.Address{}, errors.New("signed message is empty")
	}

	if signature == "" {
		return common.Address{}, errors.New("signature is empty")
	}

	publicKeyECDSA, err := crypto.SigToPub(crypto.Keccak256Hash([]byte(message)).Bytes(), common.FromHex(signature))
	if err != nil {
		return common.Address{}, err
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}
