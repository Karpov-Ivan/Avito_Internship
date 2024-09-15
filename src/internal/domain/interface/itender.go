package _interface

import (
	"avito_2024/src/internal/domain/models"
	"github.com/google/uuid"
)

type TenderService interface {
	CheckUserBelongsToOrganization(orgID uuid.UUID, username string) (bool, error)

	CreateTender(tender *models.Tender) (*models.Tender, error)

	PublishTender(tenderID uuid.UUID) error

	CloseTender(tenderID uuid.UUID) error

	EditTender(updatedTender *models.Tender) error

	GetTenderByID(tenderID uuid.UUID) (*models.Tender, error)

	GetTenders(serviceType string) ([]models.Tender, error)

	GetMyTenders(username string) ([]models.Tender, error)

	RollbackTender(tenderID uuid.UUID, version int) (*models.Tender, error)

	GetTenderStatus(tenderID uuid.UUID) (string, error)
}
