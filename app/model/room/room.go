package room

import (
	"github.com/hphphp123321/mahjong-server/app/errs"
	"github.com/hphphp123321/mahjong-server/app/model/player"
	"sync"
)

type Room struct {
	ID      string
	Name    string
	Players map[int]*player.Player

	mu sync.RWMutex
}

func NewRoom(p *player.Player, name string, id string) *Room {
	p.IsOwner = true
	r := &Room{
		ID:      id,
		Name:    name,
		Players: map[int]*player.Player{1: p},
	}
	seat := r.getIdleSeat()
	p.Seat = seat
	return r
}

func (r *Room) Leave(p *player.Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, player2 := range r.Players {
		if player2.ID == p.ID {
			if p.IsOwner {
				r.changeOwner()
			}
			if err := p.LeaveRoom(); err != nil {
				return err
			}
			delete(r.Players, i)
			return nil
		}
	}
	return errs.ErrPlayerNotInRoom
}

func (r *Room) Join(p *player.Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.Players) == 4 {
		return errs.ErrRoomIsFull
	}
	seat := r.getIdleSeat()
	if err := p.JoinRoom(r.ID, seat); err != nil {
		return err
	}
	r.Players[seat] = p
	return nil
}

func (r *Room) IsFull() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Players) == 4
}

func (r *Room) IsEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Players) == 0
}

func (r *Room) GetPlayerBySeat(seat int) (*player.Player, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p, ok := r.Players[seat]; ok {
		return p, nil
	}
	return nil, errs.ErrPlayerSeatInvalid
}

func (r *Room) GetInfo() *Info {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return &Info{
		ID:          r.ID,
		Name:        r.Name,
		PlayerInfos: r.getPlayersInfo(),
	}
}

func (r *Room) ListPlayerIDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var playerIDs []string
	for _, p := range r.Players {
		if p.ID != "" {
			playerIDs = append(playerIDs, p.ID)
		}
	}
	return playerIDs
}

func (r *Room) getPlayersInfo() []*player.Info {
	var playerInfos []*player.Info
	for _, p := range r.Players {
		playerInfos = append(playerInfos, p.GetInfo())
	}
	return playerInfos
}

func (r *Room) getIdleSeat() int {
	for _, seat := range []int{1, 2, 3, 4} {
		if _, ok := r.Players[seat]; !ok {
			return seat
		}
	}
	panic("room is full")
}

func (r *Room) changeOwner() {
	if len(r.Players) == 1 {
		return
	}
	for _, p := range r.Players {
		if p.ID != "" && p.IsOwner == false {
			p.IsOwner = true
			return
		}
	}
	return
}
