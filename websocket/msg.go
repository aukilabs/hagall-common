package websocket

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/aukilabs/go-tooling/pkg/errors"
	"github.com/aukilabs/go-tooling/pkg/logs"
	"github.com/aukilabs/hagall-common/messages/hagallpb"
	"golang.org/x/net/websocket"
	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	schedulerQueueSize = 256
)

// ResponseSender
type ResponseSender interface {
	Send(ProtoMsg)
	SendMsg(Msg)
}

// Msg represents a Hagall WebSocket message to be handled.
type Msg struct {
	Type protoreflect.Enum
	Time time.Time

	body []byte
}

// DataTo stores the message in the given value. v should be a pointer.
func (m Msg) DataTo(v ProtoMsg) error {
	return protobuf.Unmarshal(m.body, v)
}

// MsgFromProto build Msg from ProtoMsg
func MsgFromProto(v ProtoMsg) (Msg, error) {
	protoMsgVal := reflect.Indirect(reflect.ValueOf(v))
	if protoMsgVal.Kind() != reflect.Struct {
		return Msg{}, errors.New("protobuf message is not a struct").WithTag("kind", protoMsgVal.Kind())
	}
	msgTypeVal := protoMsgVal.FieldByName("Type")
	msgType, ok := msgTypeVal.Interface().(protoreflect.Enum)
	if !ok {
		return Msg{}, errors.New("protobuf message type is not a protobuf enum").WithTag("type_field_type", msgTypeVal.Type())
	}

	b, err := protobuf.Marshal(v.(protoreflect.ProtoMessage))
	if err != nil {
		return Msg{}, errors.New("encoding protobuf message failed").Wrap(err)
	}

	return Msg{
		Type: msgType,
		body: b,
		Time: v.GetTimestamp().AsTime(),
	}, nil
}

func (m Msg) TypeString() string {
	if m.Type == nil {
		return ""
	}
	return protoTypes.Type(m.Type)
}

// ProtoMsg is the interface that describe a protobuf message.
type ProtoMsg interface {
	protoreflect.ProtoMessage

	GetTimestamp() *timestamppb.Timestamp
}

// Receiver represents a function that receives a message.
type Receiver func() (Msg, int, error)

// Sender represents a function that sends a message.
type Sender func(msg Msg) (int, error)

// Send sends the given msg through the web socket.
func Send(ws *websocket.Conn, msg Msg) (int, error) {
	var written int
	codec := websocket.Codec{
		Marshal: func(v interface{}) ([]byte, byte, error) {
			buf, _ := v.([]byte)
			written = len(buf)
			return buf, byte(websocket.BinaryFrame), nil
		},
	}

	if err := codec.Send(ws, msg.body); err != nil {
		return 0, errors.New("sending message failed").
			WithType(ErrTypeMsgSendfail).
			WithTag("msg_type", msg.TypeString).
			Wrap(err)
	}

	return written, nil
}

// Receive receives the incoming message from the web socket.
func Receive(ws *websocket.Conn) (Msg, int, error) {
	var body []byte

	codec := websocket.Codec{
		Unmarshal: func(data []byte, payloadType byte, v interface{}) (err error) {
			if payloadType != websocket.BinaryFrame {
				return errors.New("received invalid websocket payload type").
					WithTag("payload_type", payloadType)
			}
			body = data
			return protobuf.Unmarshal(body, v.(protoreflect.ProtoMessage))
		},
	}

	var msg hagallpb.Msg
	if err := codec.Receive(ws, &msg); err != nil {
		return Msg{}, len(body), errors.New("receiving message failed").
			WithType(ErrTypeMsgReceiveFail).
			Wrap(err)
	}

	if msg.Timestamp == nil {
		return Msg{}, len(body), errors.New("missing message timestamp").
			WithType(ErrTypeMsgMissingTimestamp).
			WithTag("msg_type", msg.Type)
	}

	return Msg{
		Type: msg.Type,
		Time: msg.Timestamp.AsTime(),
		body: body,
	}, len(body), nil
}

// ProtoMsgType returns the type of the protobuf message as a string.
//
// It panics when msg does not have a GetType method that returns a protobuf
// error.
func ProtoMsgType(msg ProtoMsg) string {
	return protoTypes.MsgType(msg)
}

var protoTypes = protoTypeStore{types: make(map[int64]string)}

type protoTypeStore struct {
	mutex sync.RWMutex
	types map[int64]string
}

func (s *protoTypeStore) MsgType(msg ProtoMsg) string {
	t := reflect.Indirect(reflect.ValueOf(msg)).
		FieldByName("Type")

	if !t.CanInt() {
		logs.Debug("type field does not exists in ProtoMsg")
		return ""
	}

	n := t.Int()

	s.mutex.RLock()
	str, ok := s.types[n]
	s.mutex.RUnlock()

	if !ok {
		str = fmt.Sprint(t)

		s.mutex.Lock()
		s.types[n] = str
		s.mutex.Unlock()
	}
	return str
}

func (s *protoTypeStore) Type(e protoreflect.Enum) string {
	n := int64(e.Number())

	s.mutex.RLock()
	str, ok := s.types[n]
	s.mutex.RUnlock()

	if !ok {
		str = fmt.Sprint(e)

		s.mutex.Lock()
		s.types[n] = str
		s.mutex.Unlock()
	}
	return str
}

// Dispatcher represents a message dispatcher that decides if and when a message
// can be consumed.
type Dispatcher interface {
	// Dispatches the given message.
	Dispatch(context.Context, Msg) error

	// The function called when a session frame ends.
	HandleFrame()
}

// Consumer represents a message consumer.
type Consumer interface {
	// Returns the next message to be consumed.
	Consume(context.Context) (Msg, error)

	// Returns the channel that contains the consumable messages.
	Messages() <-chan Msg
}

type scheduler struct {
	queue chan Msg

	mutex                  sync.Mutex
	poseUpdates            map[uint32]Msg
	entityComponentUpdates map[string]Msg
}

func NewScheduler() *scheduler {
	return &scheduler{
		queue:                  make(chan Msg, schedulerQueueSize),
		poseUpdates:            make(map[uint32]Msg),
		entityComponentUpdates: make(map[string]Msg),
	}
}

func (s *scheduler) Close() {
	close(s.queue)
}

func (s *scheduler) Dispatch(ctx context.Context, msg Msg) error {
	switch msg.Type {
	case hagallpb.MsgType_MSG_TYPE_ENTITY_UPDATE_POSE:
		return s.dispatchEntityUpdatePose(ctx, msg)

	case hagallpb.MsgType_MSG_TYPE_ENTITY_COMPONENT_UPDATE:
		return s.dispatchEntityComponentUpdate(ctx, msg)

	default:
		s.queue <- msg
		return nil
	}
}

func (s *scheduler) dispatchEntityUpdatePose(ctx context.Context, msg Msg) error {
	var eup hagallpb.EntityUpdatePose
	if err := msg.DataTo(&eup); err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.poseUpdates[eup.EntityId] = msg
	return nil
}

func (s *scheduler) dispatchEntityComponentUpdate(ctx context.Context, msg Msg) error {
	var req hagallpb.EntityComponentUpdate
	if err := msg.DataTo(&req); err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := fmt.Sprintf("%v:%v", req.EntityComponentTypeId, req.EntityId)
	s.entityComponentUpdates[key] = msg
	return nil
}

func (s *scheduler) HandleFrame() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for id, msg := range s.poseUpdates {
		s.queue <- msg
		delete(s.poseUpdates, id)
	}

	for key, msg := range s.entityComponentUpdates {
		s.queue <- msg
		delete(s.entityComponentUpdates, key)
	}
}

func (s *scheduler) Consume(ctx context.Context) (Msg, error) {
	select {
	case <-ctx.Done():
		return Msg{}, ctx.Err()

	case msg := <-s.Messages():
		return msg, nil
	}
}

func (s *scheduler) Messages() <-chan Msg {
	return s.queue
}
