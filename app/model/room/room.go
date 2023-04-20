package room

import (
	"github.com/hphphp123321/mahjong-server/app/model/player"
	"sync"
)

type Room struct {
	ID      string
	Name    string
	Players []*player.Player

	mu sync.RWMutex
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

func (r *Room) getPlayersInfo() []*player.Info {
	var playerInfos []*player.Info
	for _, p := range r.Players {
		playerInfos = append(playerInfos, p.GetInfo())
	}
	return playerInfos
}
