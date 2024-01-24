package game

import "github.com/hphphp123321/mahjong-go/mahjong"

type Result struct {
	AllEvents  mahjong.Events
	Seat2Order map[int]int // seat -> order
	Err        error
}
