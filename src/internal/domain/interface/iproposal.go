package _interface

import (
	"avito_2024/src/internal/domain/models"
	"github.com/google/uuid"
)

type ProposalRepository interface {
	CheckUserBelongsToOrganizationByID(orgID uuid.UUID, userID uuid.UUID) (bool, error)

	CreateProposal(proposal *models.Proposal) error

	PublishProposal(proposalID uuid.UUID) error

	CancelProposal(proposalID uuid.UUID) error

	EditProposal(proposal *models.Proposal) error

	AgreeProposal(proposalID uuid.UUID) error

	DeclineProposal(proposalID uuid.UUID) error

	GetProposalByID(proposalID uuid.UUID) (*models.Proposal, error)

	GetProposalsByTender(tenderID uuid.UUID) ([]models.Proposal, error)

	GetProposalsByUsername(username string) ([]models.Proposal, error)

	RollbackProposal(bidID uuid.UUID, version int) (*models.Proposal, error)

	GetProposalStatus(proposalID uuid.UUID) (string, error)
}
