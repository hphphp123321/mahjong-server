package service

import "time"

type Client struct {
	LastTime time.Time
	Online   bool
}
