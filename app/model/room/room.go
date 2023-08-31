package room

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"github.com/hphphp123321/mahjong-server/app/errs"
	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/model/game"
	"github.com/hphphp123321/mahjong-server/app/model/player"
	"math/rand"
	"sync"
)

type Room struct {
	ID      string
	Name    string
	Players map[int]*player.Player

	OwnerSeat int

	gameRoom *game.GameRoom
	mu       sync.RWMutex
}

func NewRoom(p *player.Player, name string, id string) (*Room, error) {
	r := &Room{
		ID:        id,
		Name:      name,
		Players:   map[int]*player.Player{0: p},
		OwnerSeat: 0,
	}
	if err := p.JoinRoom(r.ID, 0); err != nil {
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
	for _, seat := range []int{0, 1, 2, 3} {
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

func (r *Room) StartGame(rule *mahjong.Rule, mode int) ([]int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.Players) != 4 {
		return nil, errs.ErrRoomNotFull
	}
	if !r.CheckAllReady() {
		return nil, errs.ErrRoomNotAllReady
	}

	var seed int64
	seatsToShuffle := []int{0, 1, 2, 3}
	if mode >= 0 {
		seed = rand.Int63()
		rand.Shuffle(4, func(i, j int) {
			seatsToShuffle[i], seatsToShuffle[j] = seatsToShuffle[j], seatsToShuffle[i]
		})
	}

	// Create game
	g := mahjong.NewMahjongGame(seed, rule)
	r.gameRoom = game.NewGameRoom(g, seatsToShuffle)

	// Start robot
	for seat, p := range r.Players {
		if p.ID != "" {
			continue
		}
		if robot, ok := global.RobotRegistry.GetRobot(p.Name); !ok {
			continue
		} else {
			actionChan := player.InitRobotStream()
			ech := r.gameRoom.RegisterSeat(seat, actionChan)
			player.StartRobotStream(robot, ech, actionChan)
		}
	}

	// Start Game
	r.gameRoom.StartGame(mode, r.CancelHumanPlayersReady)
	return seatsToShuffle, nil
}

func (r *Room) StartGameStream(p *player.Player, action chan *mahjong.Call) chan *player.GameEventChannel {
	return r.gameRoom.RegisterSeat(p.Seat, action)
}

func (r *Room) GetBoardState(p *player.Player) (*mahjong.BoardState, error) {
	if r.gameRoom == nil {
		return nil, errs.ErrGameNotStart
	}
	if _, ok := r.Players[p.Seat]; !ok {
		return nil, errs.ErrPlayerNotInRoom
	}
	return r.gameRoom.GetBoardStateBySeat(p.Seat), nil
}

func (r *Room) CancelHumanPlayersReady() {
	for _, p := range r.Players {
		if p.ID != "" {
			p.Ready = false
		}
	}
}
