package ncsclient

import (
	"encoding/json"
	"time"
)

type Receipt struct {
	AppID            string    `json:"app_id"`
	ClientID         string    `json:"client_id"`
	SessionID        string    `json:"session_id"`
	HagallWalletAddr string    `json:"hagall_wallet_addr"`
	ParticipantID    int       `json:"participant_id"`
	CreatedAt        time.Time `json:"created_at"`
	SessionJoinedAt  time.Time `json:"session_joined_at"`
	BytesSent        int64     `json:"bytes_sent"`
	BytesReceived    int64     `json:"bytes_received"`
}

func DecodeReceipt(receiptEncoded string) (Receipt, error) {
	var r Receipt
	err := json.Unmarshal([]byte(receiptEncoded), &r)
	if err != nil {
		return Receipt{}, err
	}
	return r, nil
}

func EncodeReceipt(receipt Receipt) (string, error) {
	receiptEncoded, err := json.Marshal(receipt)
	if err != nil {
		return "", err
	}
	return string(receiptEncoded), nil
}
