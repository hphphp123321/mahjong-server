package idservice

type IDService interface {
	GeneratePlayerID() (string, error)
	DeletePlayerID(id string) error
	GenerateRoomID() (string, error)
	DeleteRoomID(id string) error
}
