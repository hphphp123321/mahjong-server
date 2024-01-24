package password

type Password interface {
	Encrypt(password []byte) ([]byte, error)
	Compare(hashedPassword []byte, password []byte) error
}
