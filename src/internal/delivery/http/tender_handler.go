package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"

	"avito_2024/src/internal/domain/interface"
	"avito_2024/src/internal/domain/models"
	"github.com/google/uuid"
)

type TenderHandler struct {
	TenderService _interface.TenderService
}

func NewTenderHandler(tenderService _interface.TenderService) *TenderHandler {
	return &TenderHandler{TenderService: tenderService}
}

// Ping проверяет состояние сервиса.
// @Summary Проверка состояния сервиса
// @Description Возвращает "ok" если сервис работает
// @Tags Health Check
// @Accept  json
// @Produce  json
// @Success 200 {string} string "ok"
// @Router /api/ping [get]
func (h *TenderHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// GetTenders получает список тендеров по типу сервиса.
// @Summary Получить список тендеров
// @Description Возвращает список тендеров по переданному типу сервиса
// @Tags Tenders
// @Accept  json
// @Produce  json
// @Param serviceType query string false "Тип сервиса"
// @Success 200 {array} models.Tender "Список тендеров"
// @Failure 500 {string} string "Ошибка сервиса"
// @Router /api/tenders [get]
func (h *TenderHandler) GetTenders(w http.ResponseWriter, r *http.Request) {
	serviceType := r.URL.Query().Get("serviceType")

	tenders, err := h.TenderService.GetTenders(serviceType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenders)
}

// CreateTender создает новый тендер.
// @Summary Создать новый тендер
// @Description Создает новый тендер на основе переданных данных
// @Tags Tenders
// @Accept  json
// @Produce  json
// @Param tender body models.Tender true "Тендер"
// @Success 200 {object} models.Tender "Созданный тендер"
// @Failure 400 {string} string "Ошибка валидации"
// @Failure 500 {string} string "Ошибка сервиса"
// @Router /api/tenders/new [post]
func (h *TenderHandler) CreateTender(w http.ResponseWriter, r *http.Request) {
	var tender models.Tender
	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := h.TenderService.CheckUserBelongsToOrganization(tender.OrganizationID, tender.CreatorUsername)
	if err != nil || !status {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tenderResult, err := h.TenderService.CreateTender(&tender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenderResult)
}

// GetMyTenders получает список тендеров для конкретного пользователя по его имени.
// @Summary Получить мои тендеры
// @Description Возвращает список тендеров, созданных пользователем с указанным именем
// @Tags Tenders
// @Accept  json
// @Produce  json
// @Param username query string true "Имя пользователя"
// @Success 200 {array} models.Tender "Список тендеров"
// @Failure 500 {string} string "Ошибка сервиса"
// @Router /api/tenders/my [get]
func (h *TenderHandler) GetMyTenders(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	tenders, err := h.TenderService.GetMyTenders(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenders)
}

// EditTender редактирует существующий тендер по его ID.
// @Summary Редактировать тендер
// @Description Обновляет информацию о тендере по переданным данным и ID
// @Tags Tenders
// @Accept  json
// @Produce  json
// @Param tenderID path string true "ID тендера"  // Передаем ID тендера через URL
// @Param updatedTender body models.Tender true "Обновленный тендер"
// @Success 200 {object} models.Tender "Обновленный тендер"
// @Failure 400 {string} string "Ошибка валидации"
// @Failure 500 {string} string "Ошибка сервиса"
// @Router /api/tenders/{tenderID}/edit [patch]
func (h *TenderHandler) EditTender(w http.ResponseWriter, r *http.Request) {
	tenderID := r.URL.Path[len("/api/tenders/") : len("/api/tenders/")+36]
	var updatedTender models.Tender

	if err := json.NewDecoder(r.Body).Decode(&updatedTender); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(tenderID)
	if err != nil {
		http.Error(w, "invalid tender ID", http.StatusBadRequest)
		return
	}

	updatedTender.ID = id
	err = h.TenderService.EditTender(&updatedTender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTender)
}

// RollbackTender возвращает тендер к указанной версии.
// @Summary Откатить тендер до указанной версии
// @Description Откат тендера до определенной версии по его ID
// @Tags Tenders
// @Accept  json
// @Produce  json
// @Param tenderId path string true "ID тендера"
// @Param version path int true "Версия тендера для отката"
// @Success 200 {object} models.Tender "Откатанный тендер"
// @Failure 400 {string} string "Неверный ID тендера или версия"
// @Failure 500 {string} string "Ошибка сервиса"
// @Router /api/tenders/{tenderId}/rollback/{version} [put]
func (h *TenderHandler) RollbackTender(w http.ResponseWriter, r *http.Request) {
	tenderID := r.URL.Path[len("/api/tenders/") : len("/api/tenders/")+36]
	versionStr := r.URL.Path[len("/api/tenders/")+37:]

	id, err := uuid.Parse(tenderID)
	if err != nil {
		http.Error(w, "invalid tender ID", http.StatusBadRequest)
		return
	}

	versionStr = strings.Replace(versionStr, "rollback/", "", -1)
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "invalid version", http.StatusBadRequest)
		return
	}

	rolledBackTender, err := h.TenderService.RollbackTender(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rolledBackTender)
}

// PublishTender публикует тендер, делая его доступным для всех пользователей.
// @Summary Публикация тендера
// @Description Публикация тендера, чтобы он стал доступен всем пользователям
// @Tags Tenders
// @Param tenderId path string true "ID тендера"
// @Success 200 {string} string "Тендер успешно опубликован"
// @Failure 400 {string} string "Неверный ID тендера"
// @Failure 500 {string} string "Ошибка при публикации тендера"
// @Router /api/tenders/{tenderId}/publish [put]
func (h *TenderHandler) PublishTender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID := vars["tenderId"]

	id, err := uuid.Parse(tenderID)
	if err != nil {
		http.Error(w, "invalid tender ID", http.StatusBadRequest)
		return
	}

	err = h.TenderService.PublishTender(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Тендер успешно опубликован"))
}

// CloseTender закрывает тендер, делая его недоступным для всех, кроме ответственных лиц.
// @Summary Закрытие тендера
// @Description Закрытие тендера, чтобы он стал недоступен для всех пользователей, кроме ответственных
// @Tags Tenders
// @Param tenderId path string true "ID тендера"
// @Success 200 {string} string "Тендер успешно закрыт"
// @Failure 400 {string} string "Неверный ID тендера"
// @Failure 500 {string} string "Ошибка при закрытии тендера"
// @Router /api/tenders/{tenderId}/close [put]
func (h *TenderHandler) CloseTender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID := vars["tenderId"]

	id, err := uuid.Parse(tenderID)
	if err != nil {
		http.Error(w, "invalid tender ID", http.StatusBadRequest)
		return
	}

	err = h.TenderService.CloseTender(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Тендер успешно закрыт"))
}

// GetTenderStatus возвращает статус тендера по его ID.
// @Summary Получение статуса тендера
// @Description Возвращает текущий статус тендера
// @Tags Tenders
// @Param tenderId query string true "ID тендера"
// @Success 200 {string} string "Текущий статус тендера"
// @Failure 400 {string} string "Неверный ID тендера"
// @Failure 500 {string} string "Ошибка при получении статуса тендера"
// @Router /api/tenders/status [get]
func (h *TenderHandler) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	tenderIDStr := r.URL.Query().Get("tenderId")
	if tenderIDStr == "" {
		http.Error(w, "tenderId is required", http.StatusBadRequest)
		return
	}

	tenderID, err := uuid.Parse(tenderIDStr)
	if err != nil {
		http.Error(w, "invalid tender ID", http.StatusBadRequest)
		return
	}

	status, err := h.TenderService.GetTenderStatus(tenderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(status))
}
