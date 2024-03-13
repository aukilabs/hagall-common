package http

import "github.com/aukilabs/go-tooling/pkg/errors"

var (
	ErrDuplicatedWalletAddress = errors.New("duplicated wallet address")
	ErrBadRequest              = errors.New("invalid request body")
)

// GetErrorMessage returns custom error message from internal error.
func GetErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	switch true {
	case errors.Is(err, ErrDuplicatedWalletAddress):
		return "Address already registered, please generate a new address"
	case errors.Is(err, ErrBadRequest):
		return "Invalid request body"
	}

	return ""
}
