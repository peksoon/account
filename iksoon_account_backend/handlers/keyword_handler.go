package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"iksoon_account_backend/models"
	"iksoon_account_backend/utils"
)

type KeywordHandler struct {
	DB KeywordRepository
}

type KeywordRepository interface {
	GetKeywordSuggestions(categoryID int, query string, limit int) ([]models.KeywordSuggestion, error)
	GetKeywordsByCategory(categoryID int) ([]models.Keyword, error)
	UpsertKeyword(categoryID int, name string) (int64, error)
	CheckKeywordUsage(keywordID int) (bool, error)
	DeleteKeyword(id int) error
}

// 키워드 자동완성 핸들러
func (h *KeywordHandler) GetKeywordSuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	categoryIDStr := r.URL.Query().Get("category_id")
	if categoryIDStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "카테고리 ID가 필요합니다.")
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 카테고리 ID를 입력해주세요.")
		return
	}

	query := r.URL.Query().Get("q") // 검색어
	limit := 10                     // 기본 제한 개수

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	suggestions, err := h.DB.GetKeywordSuggestions(categoryID, query, limit)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "키워드 조회 중 오류 발생")
		return
	}

	utils.SendSuccessResponse(w, suggestions)
}

// 카테고리별 키워드 목록 조회 핸들러
func (h *KeywordHandler) GetKeywordsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	categoryIDStr := r.URL.Query().Get("category_id")
	if categoryIDStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "카테고리 ID가 필요합니다.")
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 카테고리 ID를 입력해주세요.")
		return
	}

	keywords, err := h.DB.GetKeywordsByCategory(categoryID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "키워드 조회 중 오류 발생")
		return
	}

	utils.SendSuccessResponse(w, keywords)
}

// 키워드 생성 또는 업데이트 핸들러
func (h *KeywordHandler) UpsertKeywordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	var req struct {
		CategoryID int    `json:"category_id"`
		Name       string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 데이터입니다.")
		return
	}

	// 입력 검증
	if req.CategoryID == 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "카테고리 ID는 필수입니다.")
		return
	}

	if req.Name == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "키워드 이름은 필수입니다.")
		return
	}

	keywordID, err := h.DB.UpsertKeyword(req.CategoryID, req.Name)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "키워드 처리 중 오류 발생")
		return
	}

	response := map[string]interface{}{
		"id":      keywordID,
		"message": "키워드가 성공적으로 처리되었습니다.",
	}

	utils.SendCreatedResponse(w, response)
}

// 키워드 삭제 핸들러
func (h *KeywordHandler) DeleteKeywordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "키워드 ID가 필요합니다.")
		return
	}

	keywordID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 키워드 ID를 입력해주세요.")
		return
	}

	// 키워드를 사용하는 데이터가 있는지 확인
	hasData, err := h.DB.CheckKeywordUsage(keywordID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "키워드 사용 여부 확인 중 오류 발생")
		return
	}

	if hasData {
		utils.SendErrorResponse(w, http.StatusConflict, models.ErrCodeCannotDelete, "이 키워드를 사용하는 데이터가 존재합니다.")
		return
	}

	err = h.DB.DeleteKeyword(keywordID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "키워드 삭제 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "키워드가 성공적으로 삭제되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}
