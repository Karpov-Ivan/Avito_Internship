package http

import (
	"avito_2024/src/internal/domain/interface"
	"avito_2024/src/internal/domain/models"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type ProposalHandler struct {
	ProposalRepo _interface.ProposalRepository
}

func NewProposalHandler(proposalService _interface.ProposalRepository) *ProposalHandler {
	return &ProposalHandler{ProposalRepo: proposalService}
}

// CreateProposal создает новое предложение.
// @Summary Создание предложения
// @Description Создает новое предложение от имени пользователя, проверяя принадлежность пользователя к организации
// @Tags Proposals
// @Accept  json
// @Produce  json
// @Param proposal body models.Proposal true "Данные предложения"
// @Success 200 {object} models.Proposal "Предложение успешно создано"
// @Failure 400 {string} string "Неверные данные или пользователь не принадлежит организации"
// @Failure 500 {string} string "Ошибка при создании предложения"
// @Router /api/bids/new [post]
func (h *ProposalHandler) CreateProposal(w http.ResponseWriter, r *http.Request) {
	var proposal models.Proposal
	if err := json.NewDecoder(r.Body).Decode(&proposal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := h.ProposalRepo.CheckUserBelongsToOrganizationByID(proposal.OrganizationID, proposal.AuthorID)
	if err != nil && !status {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.ProposalRepo.CreateProposal(&proposal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(proposal)
}

// GetMyProposals возвращает список предложений для конкретного пользователя.
// @Summary Получение предложений пользователя
// @Description Возвращает предложения, связанные с указанным пользователем
// @Tags Proposals
// @Produce  json
// @Param username query string true "Имя пользователя"
// @Success 200 {array} models.Proposal "Список предложений пользователя"
// @Failure 400 {string} string "Имя пользователя отсутствует"
// @Failure 500 {string} string "Ошибка при получении предложений"
// @Router /api/bids/my [get]
func (h *ProposalHandler) GetMyProposals(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	proposals, err := h.ProposalRepo.GetProposalsByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proposals)
}

// GetProposalsByTender возвращает список предложений для указанного тендера.
// @Summary Получение предложений по тендеру
// @Description Возвращает список всех предложений, связанных с указанным тендером
// @Tags Proposals
// @Produce  json
// @Param tenderId path string true "ID тендера"
// @Success 200 {array} models.Proposal "Список предложений для указанного тендера"
// @Failure 400 {string} string "Неверный ID тендера"
// @Failure 500 {string} string "Ошибка при получении предложений"
// @Router /api/bids/{tenderId}/list [get]
func (h *ProposalHandler) GetProposalsByTender(w http.ResponseWriter, r *http.Request) {
	tenderIDStr := mux.Vars(r)["tenderId"]
	tenderID, err := uuid.Parse(tenderIDStr)
	if err != nil {
		http.Error(w, "invalid tender ID", http.StatusBadRequest)
		return
	}

	proposals, err := h.ProposalRepo.GetProposalsByTender(tenderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proposals)
}

// EditProposal редактирует существующее предложение по его ID.
// @Summary Редактирование предложения
// @Description Редактирует предложение по указанному ID
// @Tags Proposals
// @Accept json
// @Produce json
// @Param bidId path string true "ID предложения"
// @Param proposal body models.Proposal true "Данные для обновления предложения"
// @Success 200 {object} models.Proposal "Обновленное предложение"
// @Failure 400 {string} string "Неверный ID предложения или некорректные данные"
// @Failure 500 {string} string "Ошибка при редактировании предложения"
// @Router /api/bids/{bidId}/edit [patch]
func (h *ProposalHandler) EditProposal(w http.ResponseWriter, r *http.Request) {
	bidIDStr := mux.Vars(r)["bidId"]
	bidID, err := uuid.Parse(bidIDStr)
	if err != nil {
		http.Error(w, "invalid bid ID", http.StatusBadRequest)
		return
	}

	var updatedProposal models.Proposal
	if err := json.NewDecoder(r.Body).Decode(&updatedProposal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedProposal.ID = bidID

	if err := h.ProposalRepo.EditProposal(&updatedProposal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProposal)
}

// RollbackProposal возвращает предложение к предыдущей версии.
// @Summary Откат версии предложения
// @Description Откатывает предложение к указанной версии по ID предложения
// @Tags Proposals
// @Produce json
// @Param bidId path string true "ID предложения"
// @Param version path int true "Версия предложения"
// @Success 200 {object} models.Proposal "Откатанное предложение"
// @Failure 400 {string} string "Неверный ID предложения или версия"
// @Failure 500 {string} string "Ошибка при откате предложения"
// @Router /api/bids/{bidId}/rollback/{version} [put]
func (h *ProposalHandler) RollbackProposal(w http.ResponseWriter, r *http.Request) {
	bidIDStr := r.URL.Path[len("/api/bids/") : len("/api/bids/")+36]
	bidID, err := uuid.Parse(bidIDStr)
	if err != nil {
		http.Error(w, "invalid bid ID", http.StatusBadRequest)
		return
	}

	versionStr := r.URL.Path[len("/api/bids/")+37:]
	versionStr = strings.Replace(versionStr, "rollback/", "", -1)
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "invalid version", http.StatusBadRequest)
		return
	}

	rolledBackProposal, err := h.ProposalRepo.RollbackProposal(bidID, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rolledBackProposal)
}

// PublishProposal публикует предложение, делая его доступным для ответственных и автора.
// @Summary Публикация предложения
// @Description Делает предложение доступным для ответственных за организацию и автора
// @Tags Proposals
// @Param bidId path string true "ID предложения"
// @Success 200 {string} string "Предложение успешно опубликовано"
// @Failure 400 {string} string "Неверный ID предложения"
// @Failure 500 {string} string "Ошибка при публикации предложения"
// @Router /api/bids/{bidId}/publish [put]
func (h *ProposalHandler) PublishProposal(w http.ResponseWriter, r *http.Request) {
	proposalIDStr := mux.Vars(r)["bidId"]
	proposalID, err := uuid.Parse(proposalIDStr)
	if err != nil {
		http.Error(w, "invalid proposal ID", http.StatusBadRequest)
		return
	}

	err = h.ProposalRepo.PublishProposal(proposalID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Предложение успешно опубликовано"))
}

// CancelProposal отменяет предложение, делая его видимым только автору и ответственным.
// @Summary Отмена предложения
// @Description Делает предложение видимым только автору и ответственным за организацию
// @Tags Proposals
// @Param bidId path string true "ID предложения"
// @Success 200 {string} string "Предложение успешно отменено"
// @Failure 400 {string} string "Неверный ID предложения"
// @Failure 500 {string} string "Ошибка при отмене предложения"
// @Router /api/bids/{bidId}/cancel [put]
func (h *ProposalHandler) CancelProposal(w http.ResponseWriter, r *http.Request) {
	proposalIDStr := mux.Vars(r)["bidId"]
	proposalID, err := uuid.Parse(proposalIDStr)
	if err != nil {
		http.Error(w, "invalid proposal ID", http.StatusBadRequest)
		return
	}

	err = h.ProposalRepo.CancelProposal(proposalID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Предложение успешно отменено"))
}

// GetProposalStatus возвращает статус предложения по его ID.
// @Summary Получение статуса предложения
// @Description Возвращает текущий статус предложения
// @Tags Proposals
// @Param bidId query string true "ID предложения"
// @Success 200 {string} string "Текущий статус предложения"
// @Failure 400 {string} string "Неверный ID предложения"
// @Failure 500 {string} string "Ошибка при получении статуса предложения"
// @Router /api/bids/status [get]
func (h *ProposalHandler) GetProposalStatus(w http.ResponseWriter, r *http.Request) {
	proposalIDStr := r.URL.Query().Get("bidId")
	if proposalIDStr == "" {
		http.Error(w, "proposalId is required", http.StatusBadRequest)
		return
	}

	proposalID, err := uuid.Parse(proposalIDStr)
	if err != nil {
		http.Error(w, "invalid proposal ID", http.StatusBadRequest)
		return
	}

	status, err := h.ProposalRepo.GetProposalStatus(proposalID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(status))
}
