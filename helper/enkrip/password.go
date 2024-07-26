package enkrip

import "golang.org/x/crypto/bcrypt"

type HashInterface interface {
	Compare(hashed, input string) error
	HashPassword(input string) (string, error)
}

type Hash struct{}

func New() HashInterface { return &Hash{} }

func (h *Hash) Compare(hashed, input string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input))
}
func (h *Hash) HashPassword(input string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input), 10)
	if err != nil {
		return "", nil
	}
	return string(hashPassword), nil
}
