package room

import "github.com/hphphp123321/mahjong-server/app/model/player"

type Info struct {
	ID          string
	Name        string
	PlayerInfos []*player.Info
}
