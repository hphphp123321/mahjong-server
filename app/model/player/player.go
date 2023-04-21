package player

import "github.com/hphphp123321/mahjong-server/app/errs"

type Player struct {
	Name   string
	RoomID string
	Seat   int

	ID string

	isReady bool
}

func NewPlayer(id string, name string) *Player {
	return &Player{
		ID:   id,
		Name: name,
	}
}

func (p *Player) GetInfo() *Info {
	return &Info{
		Name: p.Name,
		Seat: p.Seat,
	}
}

func (p *Player) GetReady() error {
	if p.isReady {
		return errs.ErrPlayerReady
	}
	p.isReady = true
	return nil
}

func (p *Player) CancelReady() error {
	if !p.isReady {
		return errs.ErrPlayerNotReady
	}
	p.isReady = false
	return nil
}

func (p *Player) IsReady() bool {
	return p.isReady
}

func (p *Player) JoinRoom(roomID string, seat int) error {
	if p.RoomID != "" {
		return errs.ErrPlayerAlreadyInRoom
	}
	p.RoomID = roomID
	p.Seat = seat
	return nil
}

func (p *Player) LeaveRoom() error {
	if p.RoomID == "" {
		return errs.ErrPlayerNotInRoom
	}
	p.RoomID = ""
	p.Seat = 0
	return nil
}
