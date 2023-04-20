package player

type Player struct {
	Name   string
	RoomID string
	Seat   int

	ID      string
	IsOwner bool
}

func (p *Player) GetInfo() *Info {
	return &Info{
		Name:    p.Name,
		Seat:    p.Seat,
		IsOwner: p.IsOwner,
	}
}
