package handlers

import (
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
	CreateCategory(name, categoryType, expenseType string) (int64, error)
	UpdateCategory(id int, name string, categoryType string, expenseType string) error
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

// CreateCategoryHandler 카테고리 생성 핸들러
func (h *CategoryHandler) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateHTTPMethod(w, r, http.MethodPost) {
		return
	}

	var req models.CategoryRequest
	if !utils.ValidateJSONRequest(w, r, &req) {
		return
	}

	utils.Debug("카테고리 생성 요청: %+v", req)

	// 입력 검증
	if req.Name == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("카테고리 이름은 필수입니다"))
		return
	}

	if req.Type != "out" && req.Type != "in" {
		utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("타입은 'out' 또는 'in'이어야 합니다"))
		return
	}

	// expenseType 기본값 처리 및 검증
	expenseType := req.ExpenseType
	if expenseType == "" {
		expenseType = "variable"
	}
	if expenseType != "fixed" && expenseType != "variable" {
		utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("expense_type은 'fixed' 또는 'variable'이어야 합니다"))
		return
	}

	categoryID, err := h.DB.CreateCategory(req.Name, req.Type, expenseType)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.SendError(w, apiErrors.ErrAlreadyExists.WithMessage("이미 존재하는 카테고리입니다"))
			return
		}
		utils.LogDatabaseError("카테고리 생성", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("카테고리 생성 실패"))
		return
	}

	response := map[string]interface{}{
		"id":      categoryID,
		"message": "카테고리가 성공적으로 생성되었습니다.",
	}

	utils.Debug("카테고리 생성 성공: ID %d", categoryID)
	utils.SendCreatedResponse(w, response)
}

// UpdateCategoryHandler 카테고리 수정 핸들러
func (h *CategoryHandler) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateHTTPMethod(w, r, http.MethodPut) {
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("카테고리 ID는 필수입니다"))
		return
	}

	categoryID, err := strconv.Atoi(idStr)
	if err != nil || categoryID <= 0 {
		utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("올바르지 않은 카테고리 ID입니다"))
		return
	}

	var req models.CategoryRequest
	if !utils.ValidateJSONRequest(w, r, &req) {
		return
	}

	utils.Debug("카테고리 수정 요청: ID %d, %+v", categoryID, req)

	// 입력 검증
	if req.Name == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("카테고리 이름은 필수입니다"))
		return
	}

	if req.Type != "out" && req.Type != "in" {
		utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("타입은 'out' 또는 'in'이어야 합니다"))
		return
	}

	// expenseType 기본값 처리 및 검증
	expenseType := req.ExpenseType
	if expenseType == "" {
		expenseType = "variable"
	}
	if expenseType != "fixed" && expenseType != "variable" {
		utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("expense_type은 'fixed' 또는 'variable'이어야 합니다"))
		return
	}

	err = h.DB.UpdateCategory(categoryID, req.Name, req.Type, expenseType)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			utils.SendError(w, apiErrors.ErrNotFound.WithMessage("카테고리를 찾을 수 없습니다"))
			return
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.SendError(w, apiErrors.ErrAlreadyExists.WithMessage("이미 존재하는 카테고리 이름입니다"))
			return
		}
		utils.LogDatabaseError("카테고리 수정", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("카테고리 수정 실패"))
		return
	}

	utils.Debug("카테고리 수정 성공: ID %d", categoryID)
	utils.SendSuccessResponse(w, utils.CreateSuccessMessage("카테고리가 성공적으로 수정되었습니다"))
}

// DeleteCategoryHandler 카테고리 삭제 핸들러
func (h *CategoryHandler) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateHTTPMethod(w, r, http.MethodDelete) {
		return
	}

	// 디버깅: 전체 쿼리 파라미터 확인
	utils.Debug("전체 쿼리 파라미터: %v", r.URL.Query())
	utils.Debug("요청 URL: %s", r.URL.String())

	idStr := r.URL.Query().Get("id")
	utils.Debug("파싱된 ID 문자열: '%s'", idStr)

	if idStr == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("카테고리 ID는 필수입니다"))
		return
	}

	categoryID, err := strconv.Atoi(idStr)
	if err != nil || categoryID <= 0 {
		utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("올바르지 않은 카테고리 ID입니다"))
		return
	}

	utils.Debug("카테고리 삭제 요청: ID %d", categoryID)

	// 카테고리를 사용하는 데이터가 있는지 확인
	hasData, err := h.DB.CheckCategoryUsage(categoryID)
	if err != nil {
		utils.LogDatabaseError("카테고리 사용 여부 확인", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("카테고리 사용 여부 확인 실패"))
		return
	}

	utils.Debug("카테고리 사용 여부 확인: ID %d, 사용 중 %v", categoryID, hasData)

	if hasData {
		utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("이 카테고리를 사용하는 데이터가 존재합니다"))
		return
	}

	err = h.DB.DeleteCategory(categoryID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			utils.SendError(w, apiErrors.ErrNotFound.WithMessage("카테고리를 찾을 수 없습니다"))
			return
		}
		utils.LogDatabaseError("카테고리 삭제", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("카테고리 삭제 실패"))
		return
	}

	utils.Debug("카테고리 삭제 성공: ID %d", categoryID)
	utils.SendSuccessResponse(w, utils.CreateSuccessMessage("카테고리가 성공적으로 삭제되었습니다"))
}

// ForceDeleteCategoryHandler 카테고리 강제 삭제 핸들러
func (h *CategoryHandler) ForceDeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateHTTPMethod(w, r, http.MethodDelete) {
		return
	}

	// 디버깅: 전체 쿼리 파라미터 확인
	utils.Debug("강제삭제 전체 쿼리 파라미터: %v", r.URL.Query())
	utils.Debug("강제삭제 요청 URL: %s", r.URL.String())

	idStr := r.URL.Query().Get("id")
	utils.Debug("강제삭제 파싱된 ID 문자열: '%s'", idStr)

	if idStr == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("카테고리 ID는 필수입니다"))
		return
	}

	categoryID, err := strconv.Atoi(idStr)
	if err != nil || categoryID <= 0 {
		utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("올바르지 않은 카테고리 ID입니다"))
		return
	}

	utils.Debug("카테고리 강제 삭제 요청: ID %d", categoryID)

	err = h.DB.ForceDeleteCategory(categoryID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			utils.SendError(w, apiErrors.ErrNotFound.WithMessage("카테고리를 찾을 수 없습니다"))
			return
		}
		utils.LogDatabaseError("카테고리 강제 삭제", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("카테고리 강제 삭제 실패"))
		return
	}

	utils.Debug("카테고리 강제 삭제 성공: ID %d", categoryID)
	utils.SendSuccessResponse(w, utils.CreateSuccessMessage("카테고리와 관련 데이터가 모두 삭제되었습니다"))
}

// CategoryRESTHandler RESTful 스타일 카테고리 핸들러
// URL 패턴: /categories/{id} 또는 /categories/{id}/force-delete
func (h *CategoryHandler) CategoryRESTHandler(w http.ResponseWriter, r *http.Request) {
	utils.Debug("RESTful 카테고리 요청: %s %s", r.Method, r.URL.Path)

	// URL 경로에서 ID 추출
	path := strings.TrimPrefix(r.URL.Path, "/categories/")
	utils.Debug("경로 분석: '%s'", path)

	if path == "" || path == "/" {
		utils.SendError(w, apiErrors.ErrInvalidRequest.WithMessage("카테고리 ID가 필요합니다"))
		return
	}

	// force-delete 확인
	isForceDelete := false
	var idStr string

	if strings.HasSuffix(path, "/force-delete") {
		isForceDelete = true
		idStr = strings.TrimSuffix(path, "/force-delete")
	} else {
		idStr = path
		// 추가 경로가 있는지 확인 (예: /categories/5/something)
		if strings.Contains(idStr, "/") {
			utils.SendError(w, apiErrors.ErrInvalidRequest.WithMessage("잘못된 요청 경로입니다"))
			return
		}
	}

	utils.Debug("ID 문자열: '%s', 강제삭제: %v", idStr, isForceDelete)

	// ID 파싱
	categoryID, err := strconv.Atoi(idStr)
	if err != nil || categoryID <= 0 {
		utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("올바르지 않은 카테고리 ID입니다"))
		return
	}

	utils.Debug("파싱된 카테고리 ID: %d", categoryID)

	// DELETE 메소드만 허용
	if r.Method != http.MethodDelete {
		utils.SendError(w, apiErrors.ErrInvalidRequest.WithMessage("DELETE 메소드만 지원됩니다"))
		return
	}

	if isForceDelete {
		// 강제 삭제 로직
		utils.Debug("RESTful 카테고리 강제 삭제 요청: ID %d", categoryID)

		err = h.DB.ForceDeleteCategory(categoryID)
		if err != nil {
			if strings.Contains(err.Error(), "no rows affected") {
				utils.SendError(w, apiErrors.ErrNotFound.WithMessage("카테고리를 찾을 수 없습니다"))
				return
			}
			utils.LogDatabaseError("카테고리 강제 삭제", err)
			utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("카테고리 강제 삭제 실패"))
			return
		}

		utils.Debug("RESTful 카테고리 강제 삭제 성공: ID %d", categoryID)
		utils.SendSuccessResponse(w, utils.CreateSuccessMessage("카테고리가 강제 삭제되었습니다"))
	} else {
		// 일반 삭제 로직
		utils.Debug("RESTful 카테고리 삭제 요청: ID %d", categoryID)

		// 카테고리를 사용하는 데이터가 있는지 확인
		hasData, err := h.DB.CheckCategoryUsage(categoryID)
		if err != nil {
			utils.LogDatabaseError("카테고리 사용 여부 확인", err)
			utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("카테고리 사용 여부 확인 실패"))
			return
		}

		utils.Debug("RESTful 카테고리 사용 여부 확인: ID %d, 사용 중 %v", categoryID, hasData)

		if hasData {
			utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("이 카테고리를 사용하는 데이터가 존재합니다"))
			return
		}

		err = h.DB.DeleteCategory(categoryID)
		if err != nil {
			if strings.Contains(err.Error(), "no rows affected") {
				utils.SendError(w, apiErrors.ErrNotFound.WithMessage("카테고리를 찾을 수 없습니다"))
				return
			}
			utils.LogDatabaseError("카테고리 삭제", err)
			utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("카테고리 삭제 실패"))
			return
		}

		utils.Debug("RESTful 카테고리 삭제 성공: ID %d", categoryID)
		utils.SendSuccessResponse(w, utils.CreateSuccessMessage("카테고리가 성공적으로 삭제되었습니다"))
	}
}
