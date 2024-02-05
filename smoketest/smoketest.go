package smoketest

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/aukilabs/hagall-common/errors"
	"github.com/aukilabs/hagall-common/messages/hagallpb"
	"github.com/aukilabs/hagall-common/scenario"
	hwebsocket "github.com/aukilabs/hagall-common/websocket"
	"golang.org/x/net/websocket"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RunSmokeTestOptions struct {
	FromEndpoint       string
	ToEndpoint         string
	ToEndpointToken    string
	Timeout            time.Duration
	MaxSessionIDLength int
}

func RunSmokeTest(ctx context.Context, opts RunSmokeTestOptions) (SmokeTestResults, error) {
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	stRes := SmokeTestResults{
		FromEndpoint: opts.FromEndpoint,
		ToEndpoint:   opts.ToEndpoint,
	}

	serverWSEndpoint := strings.ReplaceAll(opts.ToEndpoint, "https://", "wss://")
	serverWSEndpoint = strings.ReplaceAll(serverWSEndpoint, "http://", "ws://")

	cfg, err := websocket.NewConfig(serverWSEndpoint, opts.ToEndpoint)
	if err != nil {
		stRes.Status = StatusError
		return stRes, errors.New("creating websocket config failed").Wrap(err)
	}
	cfg.Header.Set("Authorization", "Bearer "+opts.ToEndpointToken)
	cfg.Header.Set("User-Agent", "SmokeTest (Go WebSocket Client golang.org/x/net/websocket)")

	conn, err := websocket.DialConfig(cfg)
	if err != nil {
		stRes.Status = StatusFailed
		return stRes, errors.New("dial websocket failed").Wrap(err)
	}
	defer conn.Close()

	entityPose := hagallpb.Pose{
		Px: float32(rand.Intn(100)),
		Py: float32(rand.Intn(100)),
		Pz: float32(rand.Intn(100)),
		Rx: float32(rand.Intn(100)),
		Ry: float32(rand.Intn(100)),
		Rz: float32(rand.Intn(100)),
		Rw: float32(rand.Intn(100)),
	}

	measure := latencies{}

	err = scenario.NewScenario(conn).
		Send(func() hwebsocket.ProtoMsg {
			measure.start()
			return &hagallpb.ParticipantJoinRequest{
				Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_REQUEST,
				Timestamp: timestamppb.Now(),
				RequestId: 1,
			}
		}).
		Receive(
			scenario.FilterByType(hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_RESPONSE),
			scenario.FilterByRequestID(1),
			func(msg hwebsocket.Msg) error {
				var res hagallpb.ParticipantJoinResponse
				if err := msg.DataTo(&res); err != nil {
					return errors.New("unmarshaling join response protobuf failed").Wrap(err)
				}

				if res.Timestamp.Seconds == 0 ||
					len(res.SessionId) == 0 ||
					len(res.SessionId) > opts.MaxSessionIDLength ||
					res.ParticipantId == 0 {
					return errors.New("invalid participant join response").
						WithTag("timestamp", res.Timestamp.AsTime()).
						WithTag("session_id", res.SessionId).
						WithTag("participant_id", res.ParticipantId)
				}

				measure.end()

				return nil
			}).
		Receive(
			scenario.FilterByType(hagallpb.MsgType_MSG_TYPE_SESSION_STATE),
			func(msg hwebsocket.Msg) error {
				var res hagallpb.SessionState
				if err := msg.DataTo(&res); err != nil {
					return errors.New("unmarshaling session state protobuf failed").Wrap(err)
				}

				if res.Timestamp.Seconds == 0 {
					return errors.New("invalid session state").
						WithTag("timestamp", res.Timestamp.AsTime())
				}
				return nil
			}).
		Send(func() hwebsocket.ProtoMsg {
			measure.start()
			return &hagallpb.EntityAddRequest{
				Type:      hagallpb.MsgType_MSG_TYPE_ENTITY_ADD_REQUEST,
				Timestamp: timestamppb.Now(),
				RequestId: 2,
				Pose:      &entityPose,
			}
		}).
		Receive(
			scenario.FilterByRequestID(2),
			scenario.FilterByType(hagallpb.MsgType_MSG_TYPE_ENTITY_ADD_RESPONSE),
			func(msg hwebsocket.Msg) error {
				var res hagallpb.EntityAddResponse
				if err := msg.DataTo(&res); err != nil {
					return errors.New("unmarshaling entity add response protobuf failed").Wrap(err)
				}

				if res.EntityId == 0 {
					return errors.New("empty entity id").
						WithTag("request_id", 2).
						WithTag("timestamp", msg.Time)
				}
				measure.end()
				return nil
			},
		).
		Run(ctx)
	if err != nil {
		stRes.Status = StatusFailed
		if errors.Is(err, context.Canceled) {
			stRes.Status = StatusTimeout
		}
		return stRes, errors.New("smoke test failed").Wrap(err)
	}

	return SmokeTestResults{
		FromEndpoint:    opts.FromEndpoint,
		ToEndpoint:      opts.ToEndpoint,
		LatencyMilliSec: measure.average(),
		Status:          StatusSuccess,
	}, nil
}

type latency struct {
	startTime time.Time
	endTime   time.Time
}

type latencies []latency

func (l *latencies) start() {
	*l = append(*l, latency{
		startTime: time.Now(),
	})
}

func (l *latencies) end() {
	(*l)[len((*l))-1].endTime = time.Now()
}

func (l *latencies) average() float64 {
	var total float64
	for _, lt := range *l {
		dur := lt.endTime.Sub(lt.startTime)
		total += float64(dur.Milliseconds())
	}
	return total / float64(len(*l))
}
