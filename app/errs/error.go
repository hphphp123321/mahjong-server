package errs

import "errors"

var (
	ErrClientNotFound      = errors.New("client not found")
	ErrMetaDataNotFound    = errors.New("no metadata in context")
	ErrHeaderIDNotFound    = errors.New("no id in header")
	ErrPlayerNotFound      = errors.New("player not found")
	ErrRoomIDInSufficient  = errors.New("room id insufficient")
	ErrRoomNotFound        = errors.New("room not found")
	ErrPlayerNotInRoom     = errors.New("player not in room")
	ErrPlayerAlreadyInRoom = errors.New("player already in room")
	ErrPlayerSeatInvalid   = errors.New("player seat invalid")
	ErrRoomFull            = errors.New("room is full")
	ErrPlayerNotOwner      = errors.New("player is not owner")
	ErrPlayerReady         = errors.New("player is ready")
	ErrPlayerNotReady      = errors.New("player is not ready")
	ErrPlayerSeatOccupied  = errors.New("player seat occupied")
	ErrStreamNotFound      = errors.New("stream not found")
	ErrRoomNotFull         = errors.New("room is not full")
	ErrRoomNotAllReady     = errors.New("room is not all ready")
	ErrGameEnd             = errors.New("game end")
	ErrGameNotStart        = errors.New("game not start")
	ErrGameEndUnexpect     = errors.New("game end unexpect")
	ErrRobotNotFound       = errors.New("robot not found")
	ErrUserNameExist       = errors.New("user name exist")
)
