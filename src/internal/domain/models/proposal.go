package models

import (
	"time"

	"github.com/google/uuid"
)

type Proposal struct {
	ID             uuid.UUID `db:"id" json:"id" binding:"required"`
	Title          string    `db:"title" json:"title" binding:"required"`
	Description    string    `db:"description" json:"description"`
	TenderID       uuid.UUID `db:"tender_id" json:"tender_id" binding:"required"`
	OrganizationID uuid.UUID `db:"organization_id" json:"organization_id" binding:"required"`
	AuthorID       uuid.UUID `db:"author_id" json:"author_id" binding:"required"`
	Status         string    `db:"status" json:"status" binding:"required"`
	Version        int       `db:"version" json:"version" binding:"required"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
