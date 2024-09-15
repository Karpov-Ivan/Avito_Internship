package models

import (
	"time"

	"github.com/google/uuid"
)

type Tender struct {
	ID              uuid.UUID `db:"id" json:"id" binding:"required"`
	Title           string    `db:"title" json:"title" binding:"required"`
	Description     string    `db:"description" json:"description"`
	Status          string    `db:"status" json:"status" binding:"required"`
	OrganizationID  uuid.UUID `db:"organization_id" json:"organizationId" binding:"required"`
	Version         int       `db:"version" json:"version" binding:"required"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
	ServiceType     string    `db:"service_type" json:"serviceType" binding:"required"`
	CreatorUsername string    `db:"creator_username" json:"creatorUsername" binding:"required"`
}
