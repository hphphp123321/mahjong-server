package playerid

type IDService interface {
	GenerateID() (string, error)
}
