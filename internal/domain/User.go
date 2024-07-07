package domain

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type User struct {
	UserID    string         `json:"userId" gorm:"PrimaryKey;unique" validate:"required"`
	FirstName string         `json:"firstName" gorm:"not null" validate:"required"`
	LastName  string         `json:"lastName" gorm:"not null" validate:"required"`
	Email     string         `json:"email" gorm:"index;unique;not null" validate:"required,email"`
	Password  string         `json:"password" gorm:"not null" validate:"required,min=6,max=30"`
	Phone     string         `json:"phone"`
	Orgs      []Organisation `gorm:"many2many:user_organisations;"`
	CreatedAt time.Time      `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"default:current_timestamp"`
}

// Validate validates the User fields
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
