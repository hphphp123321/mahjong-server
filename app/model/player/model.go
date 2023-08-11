package player

import "github.com/hphphp123321/mahjong-go/mahjong"

type Info struct {
	Name  string
	Seat  int
	Ready bool
}

type GameEventChannel struct {
	Events       mahjong.Events
	Err          error
	ValidActions mahjong.Calls
}
