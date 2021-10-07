package user

import (
	"github.com/google/uuid"
	"time"
)

type CreateUserDTO struct {
	Name           string `json:"name,omitempty" validate:"required,min=2,max=100"`
	Email          string `json:"email,omitempty" validate:"required,email,min=3,max=255"`
	Password       string `json:"password,omitempty" validate:"required,min=8,max=64"`
	RepeatPassword string `json:"repeat_password,omitempty" validate:"eqfield=Password"`
}

type User struct {
	UUID      string    `json:"uuid,omitempty" gorm:"primaryKey;size:36"`
	Name      string    `json:"name,omitempty" gorm:"size:100"`
	Email     string    `json:"email,omitempty" gorm:"uniqueIndex;size:255"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"index"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"index"`
}

type ListResponse struct {
	Users []User `json:"users,omitempty"`
}

func NewUser(dto CreateUserDTO) User {
	return User{
		UUID:      uuid.NewString(),
		Name:      dto.Name,
		Email:     dto.Email,
		Password:  dto.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
