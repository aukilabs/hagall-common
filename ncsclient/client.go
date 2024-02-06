package ncsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aukilabs/go-tooling/pkg/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

type NCSClient struct {
	Endpoint  string
	Transport http.RoundTripper
}

func NewNCSClient(endpoint string, transport http.RoundTripper) NCSClient {
	if transport == nil {
		transport = http.DefaultTransport
	}

	return NCSClient{
		Endpoint:  NormalizeEndpoint(endpoint),
		Transport: transport,
	}
}

func NormalizeEndpoint(v string) string {
	v = strings.TrimSpace(v)
	return strings.TrimRight(v, "/")
}

type ReceiptPayload struct {
	Receipt   string `json:"receipt"`
	Hash      []byte `json:"hash"`
	Signature []byte `json:"signature"`
}

func NewReceiptPayload(receipt Receipt, privateKeyString string) (ReceiptPayload, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return ReceiptPayload{}, errors.New("hex to ecdsa error").Wrap(err)
	}

	receiptJSON, err := EncodeReceipt(receipt)
	if err != nil {
		return ReceiptPayload{}, errors.New("error encoding receipt").Wrap(err)
	}

	hash := crypto.Keccak256Hash([]byte(receiptJSON))
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return ReceiptPayload{}, errors.New("ecdsa sign error").Wrap(err)
	}

	return ReceiptPayload{
		Receipt:   receiptJSON,
		Hash:      hash.Bytes(),
		Signature: signature,
	}, nil
}

func NewReceiptPayloadFromParams(receipt string, hash []byte, signature []byte) ReceiptPayload {
	return ReceiptPayload{
		Receipt:   receipt,
		Hash:      hash,
		Signature: signature,
	}
}

func (c *NCSClient) PostReceipt(ctx context.Context, payload ReceiptPayload) error {
	path := "/receipt"

	body, err := json.Marshal(payload)
	if err != nil {
		return errors.New("marshalling receipt failed").Wrap(err)
	}

	url := c.Endpoint + path
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return errors.New("creating request failed").Wrap(err)
	}

	httpcli := http.Client{
		Transport: c.Transport,
	}

	_, err = httpcli.Do(req)
	if err != nil {
		return errors.New("posting receipt failed").Wrap(err)
	}

	return nil
}
