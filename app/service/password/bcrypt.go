package password

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct {
}

func (b Bcrypt) Encrypt(data []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(data, bcrypt.DefaultCost)
}

func (b Bcrypt) Compare(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
