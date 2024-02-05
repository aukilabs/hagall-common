package testing

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	hwebsocket "github.com/aukilabs/hagall-common/websocket"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
)

func MockHagall(t *testing.T, ctx context.Context, handler func(*websocket.Conn, hwebsocket.Msg)) *httptest.Server {
	server := httptest.NewServer(websocket.Server{
		Handshake: func(c *websocket.Config, r *http.Request) error {
			return nil
		},
		Handler: func(conn *websocket.Conn) {
			recvChan := make(chan hwebsocket.Msg)
			go func() {
				for {
					msg, n, err := hwebsocket.Receive(conn)
					if err != nil {
						if errors.Is(err, net.ErrClosed) ||
							errors.Is(err, io.EOF) {
							return
						}
					}
					require.NoError(t, err)
					require.NotZero(t, n)
					recvChan <- msg
				}
			}()

			for {
				select {
				case <-ctx.Done():
					conn.Close()
				case msg := <-recvChan:
					handler(conn, msg)
				}
			}
		},
	})
	return server
}

func SendProto(t *testing.T, conn *websocket.Conn, protoMsg hwebsocket.ProtoMsg) {
	msg, err := hwebsocket.MsgFromProto(protoMsg)
	require.NoError(t, err)
	n, err := hwebsocket.Send(conn, msg)
	require.NoError(t, err)
	require.NotZero(t, n)
}
