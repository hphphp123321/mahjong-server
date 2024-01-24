package idservice

import (
	"github.com/google/uuid"
	"github.com/hphphp123321/mahjong-server/app/errs"
	"math/rand"
	"strconv"
	"sync"
)

type UUIDGenerator struct {
	RoomIDs map[string]struct{}
	lock    sync.Mutex
}

func (U *UUIDGenerator) DeletePlayerID(id string) error {
	return nil
}

func (U *UUIDGenerator) DeleteRoomID(id string) error {
	if _, exists := U.RoomIDs[id]; !exists {
		return errs.ErrRoomNotFound
	}
	U.lock.Lock()
	delete(U.RoomIDs, id)
	U.lock.Unlock()
	return nil
}

func (U *UUIDGenerator) GenerateRoomID() (string, error) {
	var id string
	if len(U.RoomIDs) >= 800000 {
		return "", errs.ErrRoomIDInSufficient
	}
	u := uuid.New()
	seed := u.ID()
	r := rand.New(rand.NewSource(int64(seed)))
	for {
		tempID := r.Intn(900000) + 100000 // 100000 - 999999
		id = strconv.Itoa(tempID)

		if _, exists := U.RoomIDs[id]; !exists {
			break
		}
	}
	return id, nil
}

func (U *UUIDGenerator) GeneratePlayerID() (uint32, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return 0, err
	}
	return id.ID(), nil
}

func (U *UUIDGenerator) GenerateGameID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
