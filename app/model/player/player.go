package player

import "github.com/hphphp123321/mahjong-server/app/errs"

type Player struct {
	Name   string
	RoomID string
	Seat   int

	ID string

	Ready bool
}

func NewPlayer(id string, name string) *Player {
	return &Player{
		ID:   id,
		Name: name,
	}
}

func NewRobot(robotType string) *Player {
	return &Player{
		ID:    "",
		Name:  robotType,
		Ready: true,
	}
}

func (p *Player) GetInfo() *Info {
	return &Info{
		Name:  p.Name,
		Seat:  p.Seat,
		Ready: p.Ready,
	}
}

func (p *Player) GetReady() error {
	if p.Ready {
		return errs.ErrPlayerReady
	}
	p.Ready = true
	return nil
}

func (p *Player) CancelReady() error {
	if !p.Ready {
		return errs.ErrPlayerNotReady
	}
	p.Ready = false
	return nil
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
