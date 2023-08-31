package room

import (
	"github.com/hphphp123321/mahjong-server/app/model/player"
	"strconv"
)

type Info struct {
	ID          string
	Name        string
	OwnerSeat   int
	PlayerInfos []*player.Info
}

func (r *Info) String() string {
	var playerInfos string
	for _, v := range r.PlayerInfos {
		playerInfos += "\n" + v.String()
	}
	return "RoomInfo{ID: " + r.ID + ", Name: " + r.Name + ", OwnerSeat: " + strconv.Itoa(r.OwnerSeat) + ", \nPlayerInfos: " + playerInfos + "}"
}
