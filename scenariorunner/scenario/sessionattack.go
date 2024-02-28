package scenario

import (
	"context"
	"sync"
	"time"

	"github.com/aukilabs/go-tooling/pkg/errors"
	"github.com/aukilabs/go-tooling/pkg/logs"
	"github.com/aukilabs/hagall-common/messages/hagallpb"
	"github.com/aukilabs/hagall-common/scenario"
	hwebsocket "github.com/aukilabs/hagall-common/websocket"
	"golang.org/x/net/websocket"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SessionAttackOptions provides options to SessionAttack scenario
type SessionAttackOptions struct {
	AttackCount    int           `env:"SESSION_ATTACK_COUNT"    help:"The number of attacks on the Hagall server."`
	AttackDuration time.Duration `env:"SESSION_ATTACK_DURATION" help:"The duration while a worker keeps attacking the Hagall server."`
	AttackDefer    time.Duration `env:"SESSION_ATTACK_DEFER"    help:"The duration before a worker starts its attack on the Hagall server."`
	SingleSession  bool          `env:"SESSION_SINGLE_SESSION"  help:"Focus attack on a single session."`
}

// sessionAttack is a scenario to join a hagall server repeatedly using multiple workers
type sessionAttack struct{}

// Run implement scenario's Run method
func (s sessionAttack) Run(ctx context.Context, opts Options) error {
	var sessionID string

	logs.WithTag("attack_count", opts.SessionAttack.AttackCount).
		WithTag("attack_duration", opts.SessionAttack.AttackDuration).
		WithTag("attack_defer", opts.SessionAttack.AttackDefer).
		WithTag("single_session", opts.SessionAttack.SingleSession).
		Info("starting session attack")

	if opts.SessionAttack.SingleSession {
		logs.Debug("creating session")
		id, closeSession, err := createGlobalSession(ctx, opts)
		if err != nil {
			logs.Error(errors.New("creating global session failed").Wrap(err))
			return nil
		}
		defer closeSession()

		sessionID = id
		logs.WithTag("session_id", sessionID).
			Info("session succesfully created")
	}

	var attackCounter attackCounter
	defer attackCounter.LogSummary()

	var wg sync.WaitGroup
	defer wg.Wait()

	for i := 0; i < opts.SessionAttack.AttackCount; i++ {
		wg.Add(1)
		attackCounter.Inc()

		go func(attackID int) {
			defer wg.Done()

			start := time.Now()

			logs.WithTag("attack_id", attackID).
				WithTag("session_id", sessionID).
				Debug("worker is attacking on hagall server")

			if err := attack(ctx, attackID, sessionID, opts); err != nil {
				attackCounter.IncErr(err)
				logs.WithTag("session_id", sessionID).
					WithTag("attack_id", attackID).
					WithTag("attack_total_duration", time.Since(start)).
					Error(errors.Newf("worker attack on hagall server succeeded 😈!").Wrap(err))
				return
			}

			logs.WithTag("attack_id", attackID).
				WithTag("session_id", sessionID).
				WithTag("attack_total_duration", time.Since(start)).
				Debug("worker attack on hagall failed 💪")
		}(i + 1)

		time.Sleep(opts.SessionAttack.AttackDefer)
	}
	return nil
}

// createGlobalSession create a single global session for subsequent tests
func createGlobalSession(ctx context.Context, opts Options) (string, func(), error) {
	conn, err := connectToHagall(ctx, opts)
	if err != nil {
		return "", nil, errors.New("connecting to hagall failed").
			WithTag("hagall", opts.Hagall).
			Wrap(err)
	}

	sessionID, err := joinSession(ctx, conn, "")
	if err != nil {
		conn.Close()
		return "", nil, err
	}

	return sessionID, func() { conn.Close() }, nil
}

// attack start running the session attack
func attack(ctx context.Context, workerID int, sessionID string, opts Options) error {
	timer := time.NewTimer(opts.SessionAttack.AttackDuration)
	defer timer.Stop()

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

	if sessionID, err = joinSession(ctx, clientA, sessionID); err != nil {
		return errors.New("joining client to session failed").
			WithTag("client", "A").
			WithTag("session_id", sessionID).
			Wrap(err)
	}

	clientB, err := connectToHagall(ctx, opts)
	if err != nil {
		return errors.New("connecting client to hagall failed").
			WithTag("client", "B").
			WithTag("hagall", opts.Hagall).
			Wrap(err)
	}
	defer clientB.Close()

	if _, err = joinSession(ctx, clientB, sessionID); err != nil {
		return errors.New("joining client to session failed").
			WithTag("client", "B").
			WithTag("session_id", sessionID).
			Wrap(err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-timer.C:
			return nil
		}
	}
}

type attackCounter struct {
	mutex  sync.Mutex
	count  int
	errors int
}

func (c *attackCounter) Inc() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.count++
}

func (c *attackCounter) IncErr(err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.errors++
}

func (c *attackCounter) LogSummary() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	logs.WithTag("total", c.count).
		WithTag("successes", c.errors).
		WithTag("failures", c.count-c.errors).
		Info("attacks summary")
}

func joinSession(ctx context.Context, conn *websocket.Conn, sessionID string) (string, error) {
	err := scenario.NewScenario(conn).
		Send(func() hwebsocket.ProtoMsg {
			return &hagallpb.ParticipantJoinRequest{
				Type:      hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_REQUEST,
				Timestamp: timestamppb.Now(),
				RequestId: 1,
				SessionId: sessionID,
			}
		}).
		Receive(
			scenario.FilterByType(hagallpb.MsgType_MSG_TYPE_PARTICIPANT_JOIN_RESPONSE),
			scenario.FilterByRequestID(1),
			func(msg hwebsocket.Msg) error {
				var res hagallpb.ParticipantJoinResponse
				if err := msg.DataTo(&res); err != nil {
					return err
				}

				sessionID = res.SessionId
				return nil
			},
		).
		Run(ctx)
	return sessionID, err
}
