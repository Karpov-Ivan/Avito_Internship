package models

import (
	"github.com/google/uuid"
)

type OrganizationResponsible struct {
	ID             uuid.UUID `db:"id" json:"id" binding:"required"`
	OrganizationID uuid.UUID `db:"organization_id" json:"organization_id" binding:"required"`
	UserID         uuid.UUID `db:"user_id" json:"user_id" binding:"required"`
}
