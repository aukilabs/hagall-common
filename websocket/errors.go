package websocket

import "github.com/aukilabs/hagall-common/errors"

const (
	// An error type for when sending a message fails.
	ErrTypeMsgSendfail = "msg_send_fail"

	// An error type for when receiving a message fails.
	ErrTypeMsgReceiveFail = "msg_receive_fail"

	// An error type for when a participant sends a request that requires it
	// already joined a session.
	ErrTypeSessionNotJoined = "session_not_joined"

	// An error type for when a message is received without a timestamp set.
	ErrTypeMsgMissingTimestamp = "msg_missing_timestamp"

	// Error when message skipped by module
	ErrTypeMsgSkip = "module_msg_skip"

	// Error when entity not exists
	ErrEntityComponentTypeNotAdded = "entity-component-type-not-added"

	// Error when entity component type already exist
	ErrEntityComponentTypeAlreadyAdded = "entity-component-type-already-added"
)

var (
	// ErrModuleMsgSkip is an error returned when a module skipped handling a
	// message.
	ErrModuleMsgSkip = errors.New("handling message is skipped").WithType(ErrTypeMsgSkip)
)
