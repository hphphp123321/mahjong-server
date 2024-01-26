package player

import "github.com/hphphp123321/mahjong-server/app/errs"

type Player struct {
	Name   string
	RoomID string
	Seat   int

	ID uint

	Ready bool
}

func NewPlayer(id uint, name string) *Player {
	return &Player{
		ID:   id,
		Name: name,
	}
}

func NewRobot(robotType string) *Player {
	return &Player{
		ID:    0,
		Name:  robotType,
		Ready: true,
	}
}

func (p *Player) GetInfo() *Info {
	return &Info{
		ID:    p.ID,
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
