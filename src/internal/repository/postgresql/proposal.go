package postgresql

import (
	"time"

	"avito_2024/src/internal/domain/interface"
	"avito_2024/src/internal/domain/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ProposalRepository struct {
	DB *sqlx.DB
}

func NewProposalRepository(db *sqlx.DB) _interface.ProposalRepository {
	return &ProposalRepository{
		DB: db,
	}
}

func (repo *ProposalRepository) CheckUserBelongsToOrganizationByID(orgID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM organization_responsible org_res
		WHERE org_res.organization_id = $1 AND org_res.user_id = $2
	`

	var exists bool
	err := repo.DB.Get(&exists, query, orgID, userID)
	if err != nil {
		return false, errors.Wrap(err, "failed to check user organization membership")
	}

	return exists, nil
}

func (repo *ProposalRepository) CreateProposal(proposal *models.Proposal) error {
	query := `
		INSERT INTO proposal (id, title, description, tender_id, organization_id, author_id, status, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	proposal.ID = uuid.New()
	proposal.Status = "CREATED"
	proposal.Version = 1
	proposal.CreatedAt = time.Now()
	proposal.UpdatedAt = time.Now()

	_, err := repo.DB.Exec(query, proposal.ID, proposal.Title, proposal.Description, proposal.TenderID, proposal.OrganizationID, proposal.AuthorID, proposal.Status, proposal.Version, proposal.CreatedAt, proposal.UpdatedAt)
	if err != nil {
		return errors.Wrap(err, "failed to create proposal")
	}

	return nil
}

func (repo *ProposalRepository) PublishProposal(proposalID uuid.UUID) error {
	query := `
		UPDATE proposal
		SET status = 'PUBLISHED', updated_at = $2
		WHERE id = $1
	`

	_, err := repo.DB.Exec(query, proposalID, time.Now())
	if err != nil {
		return errors.Wrap(err, "failed to publish proposal")
	}

	return nil
}

func (repo *ProposalRepository) CancelProposal(proposalID uuid.UUID) error {
	query := `
		UPDATE proposal
		SET status = 'CANCELED', updated_at = $2
		WHERE id = $1
	`

	_, err := repo.DB.Exec(query, proposalID, time.Now())
	if err != nil {
		return errors.Wrap(err, "failed to cancel proposal")
	}

	return nil
}

func (repo *ProposalRepository) EditProposal(proposal *models.Proposal) error {
	query := `
		UPDATE proposal
		SET title = $2, description = $3, tender_id = $4, organization_id = $5, author_id = $6, version = version + 1, updated_at = $7
		WHERE id = $1
	`

	_, err := repo.DB.Exec(query, proposal.ID, proposal.Title, proposal.Description, proposal.TenderID, proposal.OrganizationID, proposal.AuthorID, time.Now())
	if err != nil {
		return errors.Wrap(err, "failed to edit proposal")
	}

	return nil
}

func (repo *ProposalRepository) AgreeProposal(proposalID uuid.UUID) error {
	query := `
		UPDATE proposal
		SET status = 'AGREED', updated_at = $2
		WHERE id = $1
	`

	_, err := repo.DB.Exec(query, proposalID, time.Now())
	if err != nil {
		return errors.Wrap(err, "failed to agree proposal")
	}

	return nil
}

func (repo *ProposalRepository) DeclineProposal(proposalID uuid.UUID) error {
	query := `
		UPDATE proposal
		SET status = 'DECLINED', updated_at = $2
		WHERE id = $1
	`

	_, err := repo.DB.Exec(query, proposalID, time.Now())
	if err != nil {
		return errors.Wrap(err, "failed to decline proposal")
	}

	return nil
}

func (repo *ProposalRepository) GetProposalByID(proposalID uuid.UUID) (*models.Proposal, error) {
	query := `
		SELECT id, title, description, tender_id, organization_id, author_id, status, version, created_at, updated_at
		FROM proposal
		WHERE id = $1
	`

	var proposal models.Proposal
	err := repo.DB.Get(&proposal, query, proposalID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get proposal")
	}

	return &proposal, nil
}

func (repo *ProposalRepository) GetProposalsByTender(tenderID uuid.UUID) ([]models.Proposal, error) {
	query := `
		SELECT id, title, description, tender_id, organization_id, author_id, status, version, created_at, updated_at
		FROM proposal
		WHERE tender_id = $1 AND status = 'PUBLISHED'
	`

	var proposals []models.Proposal
	err := repo.DB.Select(&proposals, query, tenderID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get proposals by tender")
	}

	return proposals, nil
}

func (repo *ProposalRepository) GetProposalsByUsername(username string) ([]models.Proposal, error) {
	query := `
        SELECT p.id, p.title, p.description, p.tender_id, p.organization_id, p.author_id, p.status, p.version, p.created_at, p.updated_at
        FROM proposal p
        JOIN employee e ON p.author_id = e.id
        WHERE e.username = $1
    `

	var proposals []models.Proposal
	err := repo.DB.Select(&proposals, query, username)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get proposals by username")
	}

	return proposals, nil
}

func (repo *ProposalRepository) RollbackProposal(bidID uuid.UUID, version int) (*models.Proposal, error) {
	query := `
        UPDATE proposal
        SET version = $1
        WHERE id = $2
        RETURNING id, title, description, tender_id, organization_id, author_id, status, version, created_at, updated_at
    `

	var rolledBackProposal models.Proposal
	err := repo.DB.Get(&rolledBackProposal, query, version, bidID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to rollback proposal")
	}

	return &rolledBackProposal, nil
}

func (repo *ProposalRepository) GetProposalStatus(proposalID uuid.UUID) (string, error) {
	var status string
	query := `
		SELECT status
		FROM proposal
		WHERE id = $1
	`

	err := repo.DB.Get(&status, query, proposalID)
	if err != nil {
		return "", errors.Wrap(err, "failed to get proposal status")
	}

	return status, nil
}
