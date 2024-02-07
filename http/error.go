package http

import "github.com/aukilabs/go-tooling/pkg/errors"

var (
	ErrDuplicatedWalletAddress = "duplicated wallet address"
)

var (
	errorMessages = map[string]string{
		ErrDuplicatedWalletAddress: "Address already registered, please generate a new address",
	}
)

// GetErrorMessage returns custom error message from internal error.
func GetErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	if richErr, ok := err.(errors.Error); ok {
		if msg, ok := errorMessages[richErr.Message()]; ok {
			return msg
		}
	} else {
		if msg, ok := errorMessages[err.Error()]; ok {
			return msg
		}
	}

	err = errors.Unwrap(err)
	if err == nil {
		return ""
	}

	return GetErrorMessage(err)
}
