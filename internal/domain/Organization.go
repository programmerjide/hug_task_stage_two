package domain

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type Organisation struct {
	OrgID       string    `json:"orgId" gorm:"PrimaryKey;unique" validate:"required"`
	Name        string    `json:"name" gorm:"not null" validate:"required"`
	Description string    `json:"description"`
	Users       []User    `gorm:"many2many:user_organisations;"`
	CreatedAt   time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"default:current_timestamp"`
}

// Validate validates the Organisation fields
func (o *Organisation) Validate() error {
	validate := validator.New()
	return validate.Struct(o)
}
