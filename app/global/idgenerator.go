package global

import "github.com/hphphp123321/mahjong-server/app/component/idservice"

var (
	IDGenerator idservice.IDService = idservice.UUIDGenerator{
		RoomIDs: make(map[string]struct{}),
	}
)
