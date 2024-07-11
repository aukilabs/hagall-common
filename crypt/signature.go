package crypt

import (
	"crypto/ecdsa"
	"time"

	"github.com/aukilabs/go-tooling/pkg/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Sign signs a message with an ECDSA private key.
// Returns signature in 0x string format.
func Sign(privateKey *ecdsa.PrivateKey, message string) (string, error) {
	signature, err := crypto.Sign(crypto.Keccak256([]byte(message)), privateKey)
	if err != nil {
		return "", errors.New("signing message failed").
			Wrap(err)
	}

	return hexutil.Encode(signature), nil
}

// SignWithTimestamp signs a message with an ECDSA private key.
// Returns signature in 0x string format.
func SignWithTimestamp(privateKey *ecdsa.PrivateKey, message string) (string, string, error) {
	timestamp := time.Now().Format(time.RFC3339Nano)
	signature, err := Sign(privateKey, message+timestamp)
	if err != nil {
		return "", "", err
	}

	return signature, timestamp, nil
}

// ValidateSignedMessage validates a signed message and returns wallet address
// from signature. Expects signature in 0x string format (hex).
//
// Deprecated: ValidateSignedMessage is half-broken, as it does not check that
// the signature matches the message. You instead have to rely on checking the
// returned address matches the expected one. Use VerifySignedMessage instead.
func ValidateSignedMessage(message string, signature string) (common.Address, error) {
	if message == "" {
		return common.Address{}, errors.New("signed message is empty")
	}
	if signature == "" {
		return common.Address{}, errors.New("signature is empty")
	}

	hash := crypto.Keccak256([]byte(message))
	publicKeyECDSA, err := crypto.SigToPub(hash, common.FromHex(signature))
	if err != nil {
		return common.Address{}, err
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	signatureBytes := common.FromHex(signature)[:64]
	if !crypto.VerifySignature(publicKeyBytes, hash, signatureBytes) {
		return common.Address{}, errors.New("signature is incorrect")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}

// VerifySignedMessage verifies a signed message against a public key
// and returns wallet address recovered from signature. Expects signature
// in 0x string format (hex).
func VerifySignedMessage(publicKey *ecdsa.PublicKey, message, signature string) (common.Address, error) {
	if publicKey == nil {
		return common.Address{}, errors.New("public key is nil")
	}
	if message == "" {
		return common.Address{}, errors.New("message is empty")
	}
	if sigLen := len(signature); sigLen != 130 && sigLen != 132 {
		return common.Address{}, errors.New("invalid signature format")
	}

	hash := crypto.Keccak256([]byte(message))
	signatureBytes := common.FromHex(signature)[:64]
	publicKeyBytes := crypto.FromECDSAPub(publicKey)

	if !crypto.VerifySignature(publicKeyBytes, hash, signatureBytes) {
		return common.Address{}, errors.New("signature is incorrect")
	}

	publicKeyECDSA, err := crypto.SigToPub(hash, common.FromHex(signature))
	if err != nil {
		return common.Address{}, err
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}
