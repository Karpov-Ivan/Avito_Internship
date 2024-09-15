package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID        `db:"id" json:"id" binding:"required"`
	Name        string           `db:"name" json:"name" binding:"required"`
	Description string           `db:"description" json:"description"`
	Type        OrganizationType `db:"type" json:"type" binding:"required"`
	CreatedAt   time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time        `db:"updated_at" json:"updated_at"`
}
