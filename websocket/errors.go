package websocket

import "github.com/aukilabs/go-tooling/pkg/errors"

const (
	// Error for when sending a message fails.
	ErrTypeMsgSendfail = "msg_send_fail"

	// Error for when receiving a message fails.
	ErrTypeMsgReceiveFail = "msg_receive_fail"

	// Error for when a participant sends a request that requires that it
	// already joined a session.
	ErrTypeSessionNotJoined = "session_not_joined"

	// Error for when a message is received without a timestamp set.
	ErrTypeMsgMissingTimestamp = "msg_missing_timestamp"

	// Error for when a message is skipped by a module.
	ErrTypeMsgSkip = "module_msg_skip"

	// Error for when an entity does not exist.
	ErrEntityComponentTypeNotAdded = "entity-component-type-not-added"

	// Error for when an entity component type already exists.
	ErrEntityComponentTypeAlreadyAdded = "entity-component-type-already-added"
)

var (
	// Error returned when a module skipped handling a message.
	ErrModuleMsgSkip = errors.New("handling message is skipped").WithType(ErrTypeMsgSkip)
)
