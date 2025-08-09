package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"iksoon_account_backend/database"
	"iksoon_account_backend/models"
	"iksoon_account_backend/utils"
)

// CategoryBudgetHandler 카테고리 기준치 핸들러
type CategoryBudgetHandler struct {
	DB *database.DB
}

// NewCategoryBudgetHandler 카테고리 기준치 핸들러 생성자
func NewCategoryBudgetHandler(db *database.DB) *CategoryBudgetHandler {
	return &CategoryBudgetHandler{DB: db}
}

// GetCategoryBudgetsHandler 카테고리 기준치 목록 조회
func (h *CategoryBudgetHandler) GetCategoryBudgetsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	userName := r.URL.Query().Get("user")
	categoryIDStr := r.URL.Query().Get("category_id")

	var categoryID *int
	if categoryIDStr != "" {
		if id, err := strconv.Atoi(categoryIDStr); err == nil {
			categoryID = &id
		} else {
			utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 카테고리 ID입니다.")
			return
		}
	}

	budgets, err := h.DB.GetCategoryBudgets(userName, categoryID)
	if err != nil {
		utils.LogError("기준치 목록 조회", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "기준치 조회 중 오류 발생")
		return
	}

	utils.SendSuccessResponse(w, budgets)
}

// CreateCategoryBudgetHandler 카테고리 기준치 생성
func (h *CategoryBudgetHandler) CreateCategoryBudgetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	var req models.CategoryBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LogError("JSON 디코드", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 형식입니다.")
		return
	}

	utils.Debug("기준치 생성 요청: %+v", req)

	// 입력 검증
	if req.CategoryID <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "카테고리를 선택해주세요.")
		return
	}
	// 사용자명은 선택사항 (빈 문자열 허용)
	if req.MonthlyBudget < 0 || req.YearlyBudget < 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "기준치는 0 이상이어야 합니다.")
		return
	}
	if req.MonthlyBudget == 0 && req.YearlyBudget == 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "월별 또는 연별 기준치 중 하나는 설정해야 합니다.")
		return
	}

	// 기준치 생성
	id, err := h.DB.CreateCategoryBudget(req.CategoryID, req.UserName, req.MonthlyBudget, req.YearlyBudget)
	if err != nil {
		utils.LogError("기준치 생성", err)
		errorMsg := err.Error()

		if errorMsg == "이미 설정된 기준치가 있습니다" {
			utils.SendErrorResponse(w, http.StatusConflict, models.ErrCodeDuplicateEntry, "해당 카테고리에 대한 기준치가 이미 존재합니다.")
		} else if strings.Contains(errorMsg, "UNIQUE constraint failed") {
			utils.SendErrorResponse(w, http.StatusConflict, models.ErrCodeDuplicateEntry, "해당 카테고리에 대한 기준치가 이미 존재합니다.")
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "기준치 생성 중 오류 발생")
		}
		return
	}

	response := map[string]interface{}{
		"message": "기준치가 성공적으로 생성되었습니다.",
		"id":      id,
	}

	utils.SendCreatedResponse(w, response)
}

// UpdateCategoryBudgetHandler 카테고리 기준치 수정
func (h *CategoryBudgetHandler) UpdateCategoryBudgetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "기준치 ID가 지정되지 않았습니다.")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 기준치 ID입니다.")
		return
	}

	var req struct {
		MonthlyBudget int `json:"monthly_budget"`
		YearlyBudget  int `json:"yearly_budget"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LogError("JSON 디코드", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 형식입니다.")
		return
	}

	utils.Debug("기준치 수정 요청: ID=%d, %+v", id, req)

	// 입력 검증
	if req.MonthlyBudget < 0 || req.YearlyBudget < 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "기준치는 0 이상이어야 합니다.")
		return
	}
	if req.MonthlyBudget == 0 && req.YearlyBudget == 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "월별 또는 연별 기준치 중 하나는 설정해야 합니다.")
		return
	}

	// 기준치 수정
	err = h.DB.UpdateCategoryBudget(id, req.MonthlyBudget, req.YearlyBudget)
	if err != nil {
		utils.LogError("기준치 수정", err)
		if err.Error() == "수정할 기준치를 찾을 수 없습니다" {
			utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, err.Error())
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "기준치 수정 중 오류 발생")
		}
		return
	}

	response := map[string]string{
		"message": "기준치가 성공적으로 수정되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}

// DeleteCategoryBudgetHandler 카테고리 기준치 삭제
func (h *CategoryBudgetHandler) DeleteCategoryBudgetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "기준치 ID가 지정되지 않았습니다.")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 기준치 ID입니다.")
		return
	}

	utils.Debug("기준치 삭제 요청: ID=%d", id)

	// 기준치 삭제 (논리적 삭제)
	err = h.DB.DeleteCategoryBudget(id)
	if err != nil {
		utils.LogError("기준치 삭제", err)
		if err.Error() == "삭제할 기준치를 찾을 수 없습니다" {
			utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, err.Error())
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "기준치 삭제 중 오류 발생")
		}
		return
	}

	response := map[string]string{
		"message": "기준치가 성공적으로 삭제되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}

// GetBudgetUsageHandler 카테고리 기준치 사용량 조회
func (h *CategoryBudgetHandler) GetBudgetUsageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	userName := r.URL.Query().Get("user")
	categoryIDStr := r.URL.Query().Get("category_id")

	// userName이 빈 문자열이어도 허용 (전체 사용량 조회를 위해)
	// if userName == "" {
	// 	utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "사용자명은 필수입니다.")
	// 	return
	// }

	currentDate := time.Now()

	if categoryIDStr != "" {
		// 특정 카테고리의 기준치 사용량 조회
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 카테고리 ID입니다.")
			return
		}

		usage, err := h.DB.GetBudgetUsage(categoryID, userName, currentDate)
		if err != nil {
			utils.LogError("기준치 사용량 조회", err)
			utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "기준치 사용량 조회 중 오류 발생")
			return
		}

		if usage == nil {
			utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, "설정된 기준치가 없습니다.")
			return
		}

		utils.SendSuccessResponse(w, usage)
	} else {
		// 사용자의 모든 카테고리 기준치 사용량 조회
		usages, err := h.DB.GetAllBudgetUsages(userName, currentDate)
		if err != nil {
			utils.LogError("전체 기준치 사용량 조회", err)
			utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "기준치 사용량 조회 중 오류 발생")
			return
		}

		utils.SendSuccessResponse(w, usages)
	}
}

// UpdateMonthlyBudgetHandler 월별 기준치만 수정
func (h *CategoryBudgetHandler) UpdateMonthlyBudgetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	var req models.MonthlyBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LogError("JSON 디코드", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 형식입니다.")
		return
	}

	// 입력 검증
	if req.CategoryID <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "유효한 카테고리 ID가 필요합니다.")
		return
	}

	if req.Amount < 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "금액은 0 이상이어야 합니다.")
		return
	}

	err := h.DB.UpdateMonthlyBudget(req.CategoryID, req.UserName, req.Amount)
	if err != nil {
		utils.LogError("월별 기준치 수정", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "월별 기준치 수정 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "월별 기준치가 성공적으로 수정되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}

// UpdateYearlyBudgetHandler 연별 기준치만 수정
func (h *CategoryBudgetHandler) UpdateYearlyBudgetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	var req models.YearlyBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LogError("JSON 디코드", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 형식입니다.")
		return
	}

	// 입력 검증
	if req.CategoryID <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "유효한 카테고리 ID가 필요합니다.")
		return
	}

	if req.Amount < 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "금액은 0 이상이어야 합니다.")
		return
	}

	err := h.DB.UpdateYearlyBudget(req.CategoryID, req.UserName, req.Amount)
	if err != nil {
		utils.LogError("연별 기준치 수정", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "연별 기준치 수정 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "연별 기준치가 성공적으로 수정되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}
