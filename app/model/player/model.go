package player

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"strconv"
)

type Info struct {
	Name  string
	Seat  int
	Ready bool
}

func (i Info) String() string {
	return "PlayerInfo{Name: " + i.Name + ", Seat: " + strconv.Itoa(i.Seat) + ", Ready: " + strconv.FormatBool(i.Ready) + "}"
}

type GameEventChannel struct {
	Events       mahjong.Events
	Err          error
	ValidActions mahjong.Calls
}
