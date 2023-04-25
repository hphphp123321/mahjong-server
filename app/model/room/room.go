package room

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"github.com/hphphp123321/mahjong-server/app/errs"
	"github.com/hphphp123321/mahjong-server/app/model/player"
	"math/rand"
	"sync"
)

type Room struct {
	ID      string
	Name    string
	Players map[int]*player.Player

	OwnerSeat int

	mu sync.RWMutex

	seatsOrder  []int
	game        *mahjong.Game
	gamePlayers map[int]*mahjong.Player
}

func NewRoom(p *player.Player, name string, id string) (*Room, error) {
	r := &Room{
		ID:        id,
		Name:      name,
		Players:   map[int]*player.Player{1: p},
		OwnerSeat: 1,
	}
	if err := p.JoinRoom(r.ID, 1); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Room) LeaveRoomBySeat(seat int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.Players[seat]
	if !ok {
		return errs.ErrPlayerNotInRoom
	}
	if p.Seat == r.OwnerSeat {
		r.changeOwner()
	}
	if err := p.LeaveRoom(); err != nil {
		return err
	}
	delete(r.Players, seat)
	return nil
}

func (r *Room) Leave(p *player.Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, player2 := range r.Players {
		if player2.ID == p.ID {
			if p.Seat == r.OwnerSeat {
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
		return errs.ErrRoomFull
	}
	seat := r.getIdleSeat()
	if err := p.JoinRoom(r.ID, seat); err != nil {
		return err
	}
	r.Players[seat] = p
	return nil
}

func (r *Room) AddRobot(robotType string, seat int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.Players) == 4 {
		return errs.ErrRoomFull
	}
	if _, ok := r.Players[seat]; ok {
		return errs.ErrPlayerSeatOccupied
	}
	p := player.NewRobot(robotType)
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
	if len(r.Players) == 0 {
		return true
	}
	for _, p := range r.Players {
		if p.ID != "" {
			return false
		}
	}
	return true
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
		OwnerSeat:   r.OwnerSeat,
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

func (r *Room) CheckAllReady() bool {
	if len(r.Players) < 4 {
		return false
	}
	for _, p := range r.Players {
		if !p.Ready {
			return false
		}
	}
	return true
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
		if p.ID != "" && p.Seat != r.OwnerSeat {
			r.OwnerSeat = p.Seat
			return
		}
	}
	return
}

func (r *Room) StartGame(rule *mahjong.Rule, seed int64) ([]int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.IsFull() {
		return nil, errs.ErrRoomNotFull
	}
	if !r.CheckAllReady() {
		return nil, errs.ErrRoomNotAllReady
	}
	r.game = mahjong.NewMahjongGame(seed, rule)
	seatsToShuffle := []int{1, 2, 3, 4}
	rand.Shuffle(4, func(i, j int) {
		seatsToShuffle[i], seatsToShuffle[j] = seatsToShuffle[j], seatsToShuffle[i]
	})
	for seat, _ := range r.Players {
		r.gamePlayers[seat] = mahjong.NewMahjongPlayer()
	}
	r.seatsOrder = seatsToShuffle
	return seatsToShuffle, nil
}
