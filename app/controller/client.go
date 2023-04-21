package controller

import (
	pb "github.com/hphphp123321/mahjong-server/app/api/v1"
	"time"
)

type Client struct {
	LastTime time.Time
	Online   bool

	ID          string
	readyStream pb.Mahjong_ReadyServer
	gameStream  pb.Mahjong_GameServer
}

func NewClient(id string) *Client {
	return &Client{
		LastTime: time.Now(),
		Online:   true,
		ID:       id,
	}
}

func (c *Client) Login() {
	c.LastTime = time.Now()
	c.Online = true
}
