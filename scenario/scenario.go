package scenario

import (
	"context"

	"github.com/aukilabs/go-tooling/pkg/errors"
	"github.com/aukilabs/hagall-common/messages/hagallpb"
	hwebsocket "github.com/aukilabs/hagall-common/websocket"
	"golang.org/x/net/websocket"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	// ErrMsgSkip is an error returned when checking a received msg when a msg
	// is not to be checked.
	ErrScenarioMsgSkip = errors.New("msg skipped")
)

// ScenarioCheck represents a function that is used to check a received message
// whithin a scenario.
//
// Returning ErrScenarioMsgSkip skips the current message and the check will be
// reused on the next received message.
//
// Any other returned error stops the scenario.
type ScenarioCheck func(msg hwebsocket.Msg) error

// FilterByType is a check that skips every received message that does not have
// one of the given types.
func FilterByType(types ...protoreflect.Enum) ScenarioCheck {
	return func(msg hwebsocket.Msg) error {
		for _, t := range types {
			if t.Number() == msg.Type.Number() {
				return nil
			}
		}
		return ErrScenarioMsgSkip
	}
}

// FilterByRequestID is a check that skips every received message that does not
// have the given request id.
func FilterByRequestID(v uint32) ScenarioCheck {
	return func(msg hwebsocket.Msg) error {
		var res hagallpb.Response
		if err := msg.DataTo(&res); err != nil {
			return ErrScenarioMsgSkip
		}

		if res.RequestId != v {
			return ErrScenarioMsgSkip
		}

		return nil
	}
}

// Scenario represents a series of steps where messages can be sent, received,
// and checked.
type Scenario struct {
	ws    *websocket.Conn
	steps []any
}

// NewScenario creates a scenario.
func NewScenario(ws *websocket.Conn) *Scenario {
	return &Scenario{
		ws: ws,
	}
}

// Send create a step where the given message is sent.
func (s *Scenario) Send(newMsg func() hwebsocket.ProtoMsg) *Scenario {
	s.steps = append(s.steps, sendStep{newMsg: newMsg})
	return s
}

// Receive creates a step where a message is received and then checked with the
// given checks.
func (s *Scenario) Receive(checks ...ScenarioCheck) *Scenario {
	s.steps = append(s.steps, receiveStep{checks: checks})
	return s
}

// Run launches the scenario.
func (s *Scenario) Run(ctx context.Context) error {
	for len(s.steps) != 0 {
		err := ctx.Err()
		if err != nil {
			return err
		}

		switch step := s.steps[0].(type) {
		case sendStep:
			msg, err := hwebsocket.MsgFromProto(step.newMsg())
			if err != nil {
				err = errors.New("handling send step failed").Wrap(err)
				break
			}

			if err = s.handleSend(msg); err != nil {
				err = errors.New("handling send step failed").Wrap(err)
			}

		case receiveStep:
			if err = s.handleReceive(ctx, step.checks); err != nil {
				err = errors.New("handing receive step failed").Wrap(err)
			}
		}

		if errors.Is(err, ErrScenarioMsgSkip) {
			continue
		}
		if err != nil {
			return err
		}

		copy(s.steps, s.steps[1:])
		s.steps[len(s.steps)-1] = nil
		s.steps = s.steps[:len(s.steps)-1]
	}

	return nil
}

func (s Scenario) handleSend(msg hwebsocket.Msg) error {
	_, err := hwebsocket.Send(s.ws, msg)
	return err
}

func (s Scenario) handleReceive(ctx context.Context, checks []ScenarioCheck) error {
	errChan := make(chan error, 1)
	recvChan := make(chan hwebsocket.Msg, 1)
	defer func() {
		for len(recvChan) != 0 {
			<-recvChan
		}
		for len(errChan) != 0 {
			<-errChan
		}
	}()

	go func() {
		msg, _, err := hwebsocket.Receive(s.ws)
		if err != nil {
			errChan <- err
		} else {
			recvChan <- msg
		}
	}()

	var msg hwebsocket.Msg
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	case msg = <-recvChan:
		for i, c := range checks {
			if err := c(msg); err != nil {
				return errors.New("check failed").
					WithTag("check_index", i).
					Wrap(err)
			}
		}
	}

	return nil
}

type sendStep struct {
	newMsg func() hwebsocket.ProtoMsg
}

type receiveStep struct {
	checks []ScenarioCheck
}
