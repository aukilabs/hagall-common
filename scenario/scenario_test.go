package scenario

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aukilabs/go-tooling/pkg/errors"
	"github.com/aukilabs/hagall-common/messages/hagallpb"
	hwebsocket "github.com/aukilabs/hagall-common/websocket"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFilterByType(t *testing.T) {
	s := httptest.NewServer(websocket.Server{
		Handler: func(conn *websocket.Conn) {
			_, _, err := hwebsocket.Receive(conn)
			require.NoError(t, err)

			msg, err := hwebsocket.MsgFromProto(&hagallpb.Msg{
				Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_RESPONSE,
				Timestamp: timestamppb.Now(),
			})
			require.NoError(t, err)

			_, err = hwebsocket.Send(conn, msg)
			require.NoError(t, err)

			msg, err = hwebsocket.MsgFromProto(&hagallpb.Msg{
				Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_BROADCAST,
				Timestamp: timestamppb.Now(),
			})
			require.NoError(t, err)

			_, err = hwebsocket.Send(conn, msg)
			require.NoError(t, err)
		},
	})
	defer s.Close()

	endpoint := strings.ReplaceAll(s.URL, "http://", "ws://")
	conn, err := websocket.Dial(endpoint, "", "http://localhost")
	require.NoError(t, err)
	defer conn.Close()

	err = NewScenario(conn).
		Send(func() hwebsocket.ProtoMsg {
			return &hagallpb.Msg{
				Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_REQUEST,
				Timestamp: timestamppb.Now(),
			}
		}).
		Receive(FilterByType(hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_BROADCAST), func(msg hwebsocket.Msg) error {
			require.Equal(t, hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_BROADCAST, msg.Type)
			return nil
		}).
		Run(context.Background())
	require.NoError(t, err)
}

func TestFilterByRequestID(t *testing.T) {
	s := httptest.NewServer(websocket.Server{
		Handler: func(conn *websocket.Conn) {
			_, _, err := hwebsocket.Receive(conn)
			require.NoError(t, err)

			msg, err := hwebsocket.MsgFromProto(&hagallpb.Response{
				Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_RESPONSE,
				Timestamp: timestamppb.Now(),
				RequestId: 1,
			})
			require.NoError(t, err)

			_, err = hwebsocket.Send(conn, msg)
			require.NoError(t, err)

			msg, err = hwebsocket.MsgFromProto(&hagallpb.Response{
				Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_BROADCAST,
				Timestamp: timestamppb.Now(),
				RequestId: 2,
			})
			require.NoError(t, err)

			_, err = hwebsocket.Send(conn, msg)
			require.NoError(t, err)
		},
	})
	defer s.Close()

	endpoint := strings.ReplaceAll(s.URL, "http://", "ws://")
	conn, err := websocket.Dial(endpoint, "", "http://localhost")
	require.NoError(t, err)
	defer conn.Close()

	err = NewScenario(conn).
		Send(func() hwebsocket.ProtoMsg {
			return &hagallpb.Msg{
				Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_REQUEST,
				Timestamp: timestamppb.Now(),
			}
		}).
		Receive(FilterByRequestID(2), func(msg hwebsocket.Msg) error {
			var res hagallpb.Response
			err := msg.DataTo(&res)
			require.NoError(t, err)

			require.Equal(t, uint32(2), res.RequestId)
			return nil
		}).
		Run(context.Background())
	require.NoError(t, err)
}

func TestCustomCheck(t *testing.T) {
	s := httptest.NewServer(websocket.Server{
		Handler: func(conn *websocket.Conn) {
			_, _, err := hwebsocket.Receive(conn)
			require.NoError(t, err)

			msg, err := hwebsocket.MsgFromProto(&hagallpb.Msg{
				Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_RESPONSE,
				Timestamp: timestamppb.Now(),
			})
			require.NoError(t, err)

			_, err = hwebsocket.Send(conn, msg)
			require.NoError(t, err)

			msg, err = hwebsocket.MsgFromProto(&hagallpb.Msg{
				Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_BROADCAST,
				Timestamp: timestamppb.Now(),
			})

			_, err = hwebsocket.Send(conn, msg)
			require.NoError(t, err)
		},
	})
	defer s.Close()

	endpoint := strings.ReplaceAll(s.URL, "http://", "ws://")
	conn, err := websocket.Dial(endpoint, "", "http://localhost")
	require.NoError(t, err)
	defer conn.Close()

	t.Run("custom check succeed", func(t *testing.T) {
		err = NewScenario(conn).
			Send(func() hwebsocket.ProtoMsg {
				return &hagallpb.Msg{
					Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_REQUEST,
					Timestamp: timestamppb.Now(),
				}
			}).
			Receive(func(msg hwebsocket.Msg) error {
				if msg.Type != hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_RESPONSE {
					return errors.New("bad msg type")
				}
				return nil
			}).
			Run(context.Background())
		require.NoError(t, err)
	})

	t.Run("custom check fails", func(t *testing.T) {
		err = NewScenario(conn).
			Receive(func(msg hwebsocket.Msg) error {
				if msg.Type != hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_RESPONSE {
					return errors.New("bad msg type")
				}
				return nil
			}).
			Run(context.Background())
		require.Error(t, err)
	})
}
