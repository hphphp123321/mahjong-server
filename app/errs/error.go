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
)
