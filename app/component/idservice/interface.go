package idservice

type IDService interface {
	GeneratePlayerID() (uint32, error)
	DeletePlayerID(id string) error
	GenerateRoomID() (string, error)
	DeleteRoomID(id string) error
	GenerateGameID() (string, error)
}
