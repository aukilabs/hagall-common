package ncsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aukilabs/go-tooling/pkg/errors"
	httpcmn "github.com/aukilabs/hagall-common/http"
	"github.com/ethereum/go-ethereum/crypto"
)

// NCSClient is the Network Credit Service client.
type NCSClient struct {
	Endpoint  string
	Transport http.RoundTripper
}

// NewNCSClient returns a new NCS client with the NCS endpoint and
// a HTTP roundtripper.
func NewNCSClient(endpoint string, transport http.RoundTripper) NCSClient {
	if transport == nil {
		transport = http.DefaultTransport
	}

	return NCSClient{
		Endpoint:  httpcmn.NormalizeEndpoint(endpoint),
		Transport: transport,
	}
}

type ReceiptPayload struct {
	Receipt   string `json:"receipt"`
	Hash      []byte `json:"hash"`
	Signature []byte `json:"signature"`
}

// NewReceiptPayload returns a receipt payload signed with privateKeyString.
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

// PostReceipt sends receipts to Network Credit Service /receipt endpoint.
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
