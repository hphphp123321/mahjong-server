package room

import "github.com/hphphp123321/mahjong-server/app/model/player"

type Info struct {
	ID          string
	Name        string
	OwnerSeat   int
	PlayerInfos []*player.Info
}
