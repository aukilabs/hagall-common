package scenario

import (
	"context"
	"time"

	"github.com/aukilabs/go-tooling/pkg/errors"
	"github.com/aukilabs/go-tooling/pkg/logs"
	"github.com/aukilabs/hagall-common/messages/hagallpb"
	"github.com/aukilabs/hagall-common/scenario"
	hwebsocket "github.com/aukilabs/hagall-common/websocket"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// integrationTest tests hagall handling of common business logic such as:
// joining session, adding entity, entity component type and entity component
type integrationTest struct{}

// Run implement scenario's Run method
func (s integrationTest) Run(ctx context.Context, opts Options) error {
	ctx, cancel := context.WithTimeout(ctx, opts.IntegrationTestTimeout)
	defer cancel()

	start := time.Now()

	logs.WithTag("hds", opts.HDS).
		WithTag("hagall_endpoint", opts.HagallPublicEndpoint).
		WithTag("hagall", opts.Hagall).
		WithTag("timeout", opts.IntegrationTestTimeout).
		Info("starting integration test")

	if err := runIntegrationTest(ctx, opts); err != nil {
		logs.WithTag("test_duration", time.Since(start)).
			Error(errors.Newf("integration test failed!").Wrap(err))
		return err
	}

	logs.WithTag("test_duration", time.Since(start)).
		Info("integration test succeeded")

	return nil
}

// runIntegrationTest runs integration test with 2 clients connecting to hagall & hds
func runIntegrationTest(ctx context.Context, opts Options) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	clientA, err := connectToHagall(ctx, opts)
	if err != nil {
		return errors.New("connecting client to hagall failed").
			WithTag("client", "A").
			WithTag("hagall", opts.Hagall).
			Wrap(err)
	}
	defer clientA.Close()

	clientB, err := connectToHagall(ctx, opts)
	if err != nil {
		return errors.New("connecting client to hagall failed").
			WithTag("client", "B").
			WithTag("hagall", opts.Hagall).
			Wrap(err)
	}
	defer clientB.Close()

	var sessionID string
	var entityID uint32
	var entityEntityComponentTypeId uint32
	entityComponentData := []byte("a cool component")

	if err := scenario.NewScenario(clientA).
		Send(func() hwebsocket.ProtoMsg {
			return &hagallpb.ParticipantJoinRequest{
				Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_REQUEST,
				Timestamp: timestamppb.Now(),
				RequestId: 1,
			}
		}).
		Receive(scenario.FilterByType(
			hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_RESPONSE),
			func(msg hwebsocket.Msg) error {
				var res hagallpb.ParticipantJoinResponse
				if err := msg.DataTo(&res); err != nil {
					return errors.New("failed to join session").
						Wrap(err)
				}
				sessionID = res.SessionId
				return nil
			},
		).
		Send(func() hwebsocket.ProtoMsg {
			return &hagallpb.EntityAddRequest{
				Type:      hagallpb.MsgType_MSG_TYPE_ENTITY_ADD_REQUEST,
				Timestamp: timestamppb.Now(),
				RequestId: 2,
			}
		}).
		Receive(
			scenario.FilterByRequestID(2),
			scenario.FilterByType(hagallpb.MsgType_MSG_TYPE_ENTITY_ADD_RESPONSE),
			func(msg hwebsocket.Msg) error {
				var res hagallpb.EntityAddResponse
				if err := msg.DataTo(&res); err != nil {
					return errors.New("failed to add entity").
						Wrap(err)
				}

				entityID = res.EntityId
				return nil
			},
		).
		Send(func() hwebsocket.ProtoMsg {
			return &hagallpb.EntityComponentTypeAddRequest{
				Type:                    hagallpb.MsgType_MSG_TYPE_ENTITY_COMPONENT_TYPE_ADD_REQUEST,
				Timestamp:               timestamppb.Now(),
				RequestId:               3,
				EntityComponentTypeName: "cool component type",
			}
		}).
		Receive(
			scenario.FilterByRequestID(3),
			scenario.FilterByType(hagallpb.MsgType_MSG_TYPE_ENTITY_COMPONENT_TYPE_ADD_RESPONSE),
			func(msg hwebsocket.Msg) error {
				var res hagallpb.EntityComponentTypeAddResponse
				if err := msg.DataTo(&res); err != nil {
					return errors.New("failed to add component type").
						WithTag("entity_id", entityID).
						Wrap(err)
				}

				entityEntityComponentTypeId = res.EntityComponentTypeId
				return nil
			},
		).
		Send(func() hwebsocket.ProtoMsg {
			return &hagallpb.EntityComponentAddRequest{
				Type:                  hagallpb.MsgType_MSG_TYPE_ENTITY_COMPONENT_ADD_REQUEST,
				Timestamp:             timestamppb.Now(),
				RequestId:             4,
				EntityId:              entityID,
				EntityComponentTypeId: entityEntityComponentTypeId,
				Data:                  entityComponentData,
			}
		}).
		Receive(
			scenario.FilterByRequestID(4),
			scenario.FilterByType(hagallpb.MsgType_MSG_TYPE_ENTITY_COMPONENT_ADD_RESPONSE),
			func(msg hwebsocket.Msg) error {
				var res hagallpb.EntityComponentAddResponse
				if err := msg.DataTo(&res); err != nil {
					return errors.New("failed to add component").
						WithTag("entity_id", entityID).
						WithTag("componet_type_id", entityEntityComponentTypeId).
						Wrap(err)
				}
				return nil
			},
		).
		Run(ctx); err != nil {
		return errors.New("failed to run scenario A").
			WithTag("client", "A").
			Wrap(err)
	}

	if err := scenario.NewScenario(clientB).
		Send(func() hwebsocket.ProtoMsg {
			return &hagallpb.ParticipantJoinRequest{
				Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_REQUEST,
				Timestamp: timestamppb.Now(),
				RequestId: 1,
				SessionId: sessionID,
			}
		}).
		Receive(scenario.FilterByType(
			hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_RESPONSE),
			func(msg hwebsocket.Msg) error {
				var res hagallpb.ParticipantJoinResponse
				if err := msg.DataTo(&res); err != nil {
					return errors.New("failed to join session").
						Wrap(err)
				}
				return nil
			},
		).
		Receive(
			scenario.FilterByType(hagallpb.MsgType_MSG_TYPE_SESSION_STATE),
			func(msg hwebsocket.Msg) error {
				var state hagallpb.SessionState
				if err := msg.DataTo(&state); err != nil {
					return errors.New("failed to receive session state").
						Wrap(err)
				}

				if len(state.EntityComponents) != 1 {
					return errors.New("failed to received initial entity component state").
						Wrap(err)
				}
				entityComponent := state.EntityComponents[0]
				if entityEntityComponentTypeId != entityComponent.EntityComponentTypeId ||
					entityID != entityComponent.EntityId ||
					string(entityComponentData) != string(entityComponent.Data) {
					return errors.New("received invalid entity component data").
						WithTag("entity_id", entityComponent.EntityId).
						WithTag("entity_component_type_id", entityComponent.EntityComponentTypeId).
						WithTag("entity_component_data", string(entityComponent.Data)).
						Wrap(err)
				}
				return err
			},
		).
		Run(ctx); err != nil {
		return errors.New("failed to run scenarios").
			WithTag("client", "B").
			Wrap(err)
	}

	return nil
}
