package postgresql

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"avito_2024/src/internal/domain/models"
)

type TenderRepository struct {
	DB *sqlx.DB
}

func NewTenderRepository(db *sqlx.DB) *TenderRepository {
	return &TenderRepository{
		DB: db,
	}
}

func (repo *TenderRepository) CheckUserBelongsToOrganization(orgID uuid.UUID, username string) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM organization_responsible org_res
		WHERE org_res.organization_id = $1 AND org_res.user_id = (
			SELECT id FROM employee WHERE username = $2
		)
	`

	var exists bool
	err := repo.DB.Get(&exists, query, orgID, username)
	if err != nil {
		return false, errors.Wrap(err, "failed to check user organization membership")
	}

	return exists, nil
}

func (repo *TenderRepository) CreateTender(tender *models.Tender) (*models.Tender, error) {
	query := `
		INSERT INTO tender (id, title, description, status, organization_id, version, created_at, updated_at, service_type, creator_username)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	tender.ID = uuid.New()
	tender.Status = "CREATED"
	tender.Version = 1
	tender.CreatedAt = time.Now()
	tender.UpdatedAt = time.Now()

	_, err := repo.DB.Exec(query, tender.ID, tender.Title, tender.Description, tender.Status, tender.OrganizationID, tender.Version, tender.CreatedAt, tender.UpdatedAt, tender.ServiceType, tender.CreatorUsername)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tender")
	}

	return tender, nil
}

func (repo *TenderRepository) PublishTender(tenderID uuid.UUID) error {
	query := `
		UPDATE tender
		SET status = 'PUBLISHED', updated_at = $2
		WHERE id = $1
	`

	_, err := repo.DB.Exec(query, tenderID, time.Now())
	if err != nil {
		return errors.Wrap(err, "failed to publish tender")
	}

	return nil
}

func (repo *TenderRepository) CloseTender(tenderID uuid.UUID) error {
	query := `
		UPDATE tender
		SET status = 'CLOSED', updated_at = $2
		WHERE id = $1
	`

	_, err := repo.DB.Exec(query, tenderID, time.Now())
	if err != nil {
		return errors.Wrap(err, "failed to close tender")
	}

	return nil
}

func (repo *TenderRepository) EditTender(tender *models.Tender) error {
	query := `
		UPDATE tender
		SET title = $2, description = $3, version = version + 1, updated_at = $4
		WHERE id = $1
	`

	_, err := repo.DB.Exec(query, tender.ID, tender.Title, tender.Description, time.Now())
	if err != nil {
		return errors.Wrap(err, "failed to edit tender")
	}

	return nil
}

func (repo *TenderRepository) GetTenderByID(tenderID uuid.UUID) (*models.Tender, error) {
	query := `
		SELECT id, title, description, status, organization_id, version, created_at, updated_at, service_type, creator_username
		FROM tender
		WHERE id = $1
	`

	var tenderRepo models.Tender
	err := repo.DB.Get(&tenderRepo, query, tenderID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tender")
	}

	return &tenderRepo, nil
}

func (repo *TenderRepository) GetTenders(serviceType string) ([]models.Tender, error) {
	var query string
	var args []interface{}

	if serviceType != "" {
		query = `
			SELECT id, title, description, status, organization_id, version, created_at, updated_at, service_type, creator_username
			FROM tender
			WHERE service_type = $1 AND status = 'PUBLISHED'
		`
		args = append(args, serviceType)
	} else {
		query = `
			SELECT id, title, description, status, organization_id, version, created_at, updated_at, service_type, creator_username
			FROM tender
			WHERE status = 'PUBLISHED'
		`
	}

	var tenders []models.Tender
	err := repo.DB.Select(&tenders, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tenders")
	}

	return tenders, nil
}

func (repo *TenderRepository) GetMyTenders(username string) ([]models.Tender, error) {
	query := `
		SELECT t.id, t.title, t.description, t.status, t.organization_id, t.version, t.created_at, t.updated_at, service_type, creator_username
		FROM tender t
		JOIN organization_responsible org_res ON t.organization_id = org_res.organization_id
		JOIN employee e ON org_res.user_id = e.id
		WHERE e.username = $1
	`

	var tenders []models.Tender
	err := repo.DB.Select(&tenders, query, username)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get my tenders")
	}

	return tenders, nil
}

func (repo *TenderRepository) RollbackTender(tenderID uuid.UUID, version int) (*models.Tender, error) {
	query := `
		SELECT id, title, description, status, organization_id, version, created_at, updated_at, service_type, creator_username
		FROM tender
		WHERE id = $1 AND version = $2
	`

	var tenderRepo models.Tender
	err := repo.DB.Get(&tenderRepo, query, tenderID, version)
	if err != nil {
		return nil, errors.Wrap(err, "failed to rollback tender")
	}

	return &tenderRepo, nil
}

func (repo *TenderRepository) GetTenderStatus(tenderID uuid.UUID) (string, error) {
	var status string
	query := `
		SELECT status
		FROM tender
		WHERE id = $1
	`

	err := repo.DB.Get(&status, query, tenderID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get tender status")
	}

	return status, nil
}
