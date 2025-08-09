package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	apiErrors "iksoon_account_backend/errors"
	"iksoon_account_backend/models"
	"iksoon_account_backend/utils"
)

type CategoryHandler struct {
	DB CategoryRepository
}

type CategoryRepository interface {
	GetCategories(categoryType string) ([]models.Category, error)
	CreateCategory(name, categoryType string) (int64, error)
	UpdateCategory(id int, name string, categoryType string) error
	CheckCategoryUsage(categoryID int) (bool, error)
	DeleteCategory(id int) error
	ForceDeleteCategory(id int) error
}

// GetCategoriesHandler 카테고리 목록 조회 핸들러
func (h *CategoryHandler) GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendError(w, apiErrors.ErrInvalidRequest.WithMessage("지원되지 않는 메소드입니다"))
		return
	}

	categoryType := r.URL.Query().Get("type") // 'out' 또는 'in'
	utils.Debug("카테고리 조회 요청: type=%s", categoryType)

	categories, err := h.DB.GetCategories(categoryType)
	if err != nil {
		utils.LogDatabaseError("카테고리 조회", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("카테고리 조회 실패"))
		return
	}

	utils.Debug("카테고리 조회 성공: %d개", len(categories))
	utils.SendSuccessResponse(w, categories)
}

// 카테고리 생성 핸들러
func (h *CategoryHandler) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	var req models.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 데이터입니다.")
		return
	}

	// 입력 검증
	if req.Name == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "카테고리 이름은 필수입니다.")
		return
	}

	if req.Type != "out" && req.Type != "in" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "타입은 'out' 또는 'in'이어야 합니다.")
		return
	}

	categoryID, err := h.DB.CreateCategory(req.Name, req.Type)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.SendErrorResponse(w, http.StatusConflict, models.ErrCodeDuplicateEntry, "이미 존재하는 카테고리입니다.")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "카테고리 생성 중 오류 발생")
		return
	}

	response := map[string]interface{}{
		"id":      categoryID,
		"message": "카테고리가 성공적으로 생성되었습니다.",
	}

	utils.SendCreatedResponse(w, response)
}

// 카테고리 수정 핸들러
func (h *CategoryHandler) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "카테고리 ID가 필요합니다.")
		return
	}

	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 카테고리 ID를 입력해주세요.")
		return
	}

	var req models.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 데이터입니다.")
		return
	}

	if req.Name == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "카테고리 이름은 필수입니다.")
		return
	}

	if req.Type != "out" && req.Type != "in" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 카테고리 타입을 선택해주세요.")
		return
	}

	err = h.DB.UpdateCategory(categoryID, req.Name, req.Type)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, "존재하지 않는 카테고리입니다.")
			return
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.SendErrorResponse(w, http.StatusConflict, models.ErrCodeDuplicateEntry, "이미 존재하는 카테고리 이름입니다.")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "카테고리 수정 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "카테고리가 성공적으로 수정되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}

// 카테고리 삭제 핸들러
func (h *CategoryHandler) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "카테고리 ID가 필요합니다.")
		return
	}

	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 카테고리 ID를 입력해주세요.")
		return
	}

	// 카테고리를 사용하는 데이터가 있는지 확인
	hasData, err := h.DB.CheckCategoryUsage(categoryID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "카테고리 사용 여부 확인 중 오류 발생")
		return
	}

	if hasData {
		utils.SendErrorResponse(w, http.StatusConflict, models.ErrCodeCannotDelete, "이 카테고리를 사용하는 데이터가 존재합니다. 삭제하시겠습니까?")
		return
	}

	err = h.DB.DeleteCategory(categoryID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, "존재하지 않는 카테고리입니다.")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "카테고리 삭제 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "카테고리가 성공적으로 삭제되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}

// 카테고리 강제 삭제 핸들러
func (h *CategoryHandler) ForceDeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "카테고리 ID가 필요합니다.")
		return
	}

	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 카테고리 ID를 입력해주세요.")
		return
	}

	err = h.DB.ForceDeleteCategory(categoryID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, "존재하지 않는 카테고리입니다.")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "카테고리 삭제 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "카테고리와 관련 데이터가 모두 삭제되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}
