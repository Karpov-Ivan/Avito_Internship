package models

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID        uuid.UUID `db:"id" json:"id" binding:"required"`
	Username  string    `db:"username" json:"username" binding:"required"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
