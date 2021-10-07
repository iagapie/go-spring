package password

import "golang.org/x/crypto/bcrypt"

type (
	Encoder interface {
		Encode(string) (string, error)
		IsValid(string, string) bool
	}

	BcryptEncoder struct {
		cost int
	}
)

func New(cost int) *BcryptEncoder {
	return &BcryptEncoder{
		cost: cost,
	}
}

func NewDefault() *BcryptEncoder {
	return New(bcrypt.DefaultCost)
}

func (e *BcryptEncoder) Encode(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), e.cost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (e *BcryptEncoder) IsValid(encoded, raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encoded), []byte(raw)) == nil
}
