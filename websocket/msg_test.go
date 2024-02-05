package websocket

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aukilabs/hagall-common/messages/hagallpb"
	"github.com/aukilabs/hagall-common/messages/vikjapb"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSendReceive(t *testing.T) {
	req := hagallpb.ParticipantJoinRequest{
		Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_REQUEST,
		Timestamp: timestamppb.Now(),
		RequestId: 42,
		SessionId: "tedisinthehotspring",
	}

	s := httptest.NewServer(websocket.Server{
		Handler: func(ws *websocket.Conn) {
			msg, n, err := Receive(ws)
			require.NoError(t, err)
			require.NotZero(t, n)

			var pjr hagallpb.ParticipantJoinRequest
			err = msg.DataTo(&pjr)
			require.NoError(t, err)
			require.Equal(t, req.Type, pjr.Type)
			require.True(t, req.Timestamp.AsTime().Equal(pjr.Timestamp.AsTime()))
			require.Equal(t, req.RequestId, pjr.RequestId)
			require.Equal(t, req.SessionId, pjr.SessionId)

			msg, err = MsgFromProto(&hagallpb.ParticipantJoinResponse{
				Type:          hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_RESPONSE,
				Timestamp:     timestamppb.Now(),
				RequestId:     pjr.RequestId,
				SessionId:     pjr.SessionId,
				ParticipantId: 1,
			})
			require.NoError(t, err)

			n, err = Send(ws, msg)
			require.NoError(t, err)
			require.NotZero(t, n)
		},
	})
	defer s.Close()

	endpoint := strings.ReplaceAll(s.URL, "http://", "ws://")
	ws, err := websocket.Dial(endpoint, "", "http://localhost")
	require.NoError(t, err)
	defer ws.Close()

	msg, err := MsgFromProto(&req)
	require.NoError(t, err)

	n, err := Send(ws, msg)
	require.NoError(t, err)
	require.NotZero(t, n)

	msg, n, err = Receive(ws)
	require.NoError(t, err)
	require.NotZero(t, n)

	var res hagallpb.ParticipantJoinResponse
	err = msg.DataTo(&res)
	require.NoError(t, err)
	require.Equal(t, hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_RESPONSE, res.Type)
	require.NotNil(t, res.Timestamp)
	require.Equal(t, req.RequestId, res.RequestId)
	require.Equal(t, req.SessionId, res.SessionId)
	require.Equal(t, uint32(1), res.ParticipantId)
}

func TestProtoMsgType(t *testing.T) {
	utests := []struct {
		scenario string
		msg      ProtoMsg
		expected string
	}{
		{
			scenario: "default msg type",
			msg:      &hagallpb.Msg{},
			expected: hagallpb.MsgType_MSG_TYPE_ERROR_RESPONSE.String(),
		},
		{
			scenario: "defined msg type",
			msg:      &hagallpb.Msg{Type: hagallpb.MsgType_MSG_TYPE_ENTITY_UPDATE_POSE},
			expected: hagallpb.MsgType_MSG_TYPE_ENTITY_UPDATE_POSE.String(),
		},
		{
			scenario: "undefined msg type",
			msg:      &hagallpb.Msg{Type: hagallpb.MsgType(666)},
			expected: "666",
		},
	}

	for _, u := range utests {
		s := ProtoMsgType(u.msg)
		require.Equal(t, u.expected, s)
	}
}

func BenchmarkProtoMsgType(b *testing.B) {
	msg := &hagallpb.Msg{Type: hagallpb.MsgType_MSG_TYPE_ENTITY_UPDATE_POSE}

	for i := 0; i < b.N; i++ {
		ProtoMsgType(msg)
	}
}

func BenchmarkMsgTypeString(b *testing.B) {
	msg := Msg{Type: hagallpb.MsgType_MSG_TYPE_ENTITY_UPDATE_POSE}

	for i := 0; i < b.N; i++ {
		msg.TypeString()
	}
}

func TestSchedulerDispatch(t *testing.T) {
	t.Run("dispatch message", func(t *testing.T) {
		s := NewScheduler()

		s.Dispatch(context.Background(), Msg{
			Type: hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_REQUEST,
		})

		require.Len(t, s.queue, 1)
	})

	t.Run("dispatch update pose message", func(t *testing.T) {
		s := NewScheduler()

		up := hagallpb.EntityUpdatePose{
			Type:      hagallpb.MsgType_MSG_TYPE_ENTITY_UPDATE_POSE,
			Timestamp: timestamppb.Now(),
			EntityId:  42,
		}

		b, err := protobuf.Marshal(&up)
		require.NoError(t, err)

		msg := Msg{
			Type: up.Type,
			Time: up.Timestamp.AsTime(),
			body: b,
		}

		ctx := context.Background()
		for i := 0; i < 10; i++ {
			s.Dispatch(ctx, msg)
		}

		require.Empty(t, s.queue)
		require.Len(t, s.poseUpdates, 1)
		require.NotNil(t, s.poseUpdates[42])
	})

	t.Run("dispatch update entity component message", func(t *testing.T) {
		s := NewScheduler()

		ecu := hagallpb.EntityComponentUpdate{
			Type:                  hagallpb.MsgType_MSG_TYPE_ENTITY_COMPONENT_UPDATE,
			Timestamp:             timestamppb.Now(),
			EntityComponentTypeId: 21,
			EntityId:              42,
		}

		b, err := protobuf.Marshal(&ecu)
		require.NoError(t, err)

		msg := Msg{
			Type: ecu.Type,
			Time: ecu.Timestamp.AsTime(),
			body: b,
		}

		ctx := context.Background()
		for i := 0; i < 10; i++ {
			s.Dispatch(ctx, msg)
		}

		require.Empty(t, s.queue)
		require.Len(t, s.entityComponentUpdates, 1)
		require.Contains(t, s.entityComponentUpdates, fmt.Sprintf("%v:%v", 21, 42))
	})
}

func TestSchedulerHandleFrame(t *testing.T) {
	s := NewScheduler()

	eup := &hagallpb.EntityUpdatePose{
		Type:      hagallpb.MsgType_MSG_TYPE_ENTITY_UPDATE_POSE,
		Timestamp: timestamppb.Now(),
		EntityId:  42,
	}
	eupBytes, err := protobuf.Marshal(eup)
	require.NoError(t, err)
	s.Dispatch(context.Background(), Msg{
		Type: eup.Type,
		Time: eup.GetTimestamp().AsTime(),
		body: eupBytes,
	})

	ecu := &hagallpb.EntityComponentUpdate{
		Type:                  hagallpb.MsgType_MSG_TYPE_ENTITY_COMPONENT_UPDATE,
		Timestamp:             timestamppb.Now(),
		EntityComponentTypeId: 42,
		EntityId:              21,
	}
	ecuBytes, err := protobuf.Marshal(ecu)
	require.NoError(t, err)
	s.Dispatch(context.Background(), Msg{
		Type: ecu.Type,
		Time: ecu.GetTimestamp().AsTime(),
		body: ecuBytes,
	})

	require.Len(t, s.poseUpdates, 1)
	require.Len(t, s.entityComponentUpdates, 1)

	s.HandleFrame()
	require.Empty(t, s.poseUpdates)
	require.Empty(t, s.entityComponentUpdates)
	require.Len(t, s.queue, 2)
}

func TestSchedulerConsumer(t *testing.T) {
	s := NewScheduler()

	ctx := context.Background()

	s.Dispatch(ctx, Msg{
		Type: hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_REQUEST,
	})

	msg, err := s.Consume(ctx)
	require.NoError(t, err)
	require.Equal(t, hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_REQUEST, msg.Type)
	require.Empty(t, s.queue)
}

func TestMsgFromProto(t *testing.T) {
	utests := []struct {
		scenario string
		in       ProtoMsg
		out      Msg
	}{
		{
			scenario: "converting hagall proto message succeeds",
			in: &hagallpb.Msg{
				Type:      hagallpb.MsgType_MSG_TYPE_ENTITY_ADD_REQUEST,
				Timestamp: timestamppb.Now(),
			},
			out: Msg{
				Type: hagallpb.MsgType_MSG_TYPE_ENTITY_ADD_REQUEST,
			},
		},
		{
			scenario: "converting hagall module proto message succeeds",
			in: &vikjapb.EntityActionRequest{
				Type:      vikjapb.MsgType_MSG_TYPE_VIKJA_ENTITY_ACTION_REQUEST,
				Timestamp: timestamppb.Now(),
			},
			out: Msg{
				Type: vikjapb.MsgType_MSG_TYPE_VIKJA_ENTITY_ACTION_REQUEST,
			},
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			msg, err := MsgFromProto(u.in)
			require.NoError(t, err)
			require.Equal(t, u.out.Type, msg.Type)
			require.NotZero(t, msg.Time)
			require.NotEmpty(t, msg.body)
		})
	}
}
