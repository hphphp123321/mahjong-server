package playerid

import "github.com/google/uuid"

type UUIDGenerator struct {
}

func (U UUIDGenerator) GenerateID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
