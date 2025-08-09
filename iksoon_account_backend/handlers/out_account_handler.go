package handlers

import (
	"encoding/json"
	"net/http"

	apiErrors "iksoon_account_backend/errors"
	"iksoon_account_backend/models"
	"iksoon_account_backend/utils"
)

type OutAccountHandler struct {
	DB        OutAccountRepository
	KeywordDB KeywordRepository
}

type OutAccountRepository interface {
	InsertOutAccount(date, user string, money, categoryID int, keywordID *int, paymentMethodID int, memo string) error
	GetOutAccountsByDate(date string) ([]models.OutAccount, error)
	GetOutAccountsForMonth(year, month string) ([]models.OutAccount, error)
	UpdateOutAccount(uuid, date, user string, money, categoryID int, keywordID *int, paymentMethodID int, memo string) error
	DeleteOutAccount(uuid string) error
	GetOutAccountByUUID(uuid string) (*models.OutAccount, error)
}

// InsertOutAccountHandler 새로운 구조의 지출 데이터 삽입 핸들러
func (h *OutAccountHandler) InsertOutAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendError(w, apiErrors.ErrInvalidRequest.WithMessage("지원되지 않는 메소드입니다"))
		return
	}

	var req struct {
		Date            string `json:"date"`
		User            string `json:"user"`
		Money           int    `json:"money"`
		CategoryID      int    `json:"category_id"`
		KeywordName     string `json:"keyword_name,omitempty"`
		PaymentMethodID int    `json:"payment_method_id"`
		Memo            string `json:"memo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LogError("JSON 디코드", err)
		utils.SendError(w, apiErrors.ErrInvalidJSON)
		return
	}

	utils.Debug("지출 데이터 삽입 요청: %+v", req)

	// 입력 검증
	if req.Date == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("날짜는 필수입니다"))
		return
	}
	if req.User == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("사용자는 필수입니다"))
		return
	}
	if req.Money <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "금액은 0보다 커야 합니다.")
		return
	}
	if req.CategoryID <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "카테고리를 선택해주세요.")
		return
	}
	if req.PaymentMethodID <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "결제수단을 선택해주세요.")
		return
	}

	// 키워드 처리 (있는 경우)
	var keywordID *int
	if req.KeywordName != "" {
		id, err := h.KeywordDB.UpsertKeyword(req.CategoryID, req.KeywordName)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "키워드 처리 중 오류 발생")
			return
		}
		keywordIDValue := int(id)
		keywordID = &keywordIDValue
	}

	// 지출 데이터 삽입
	err := h.DB.InsertOutAccount(req.Date, req.User, req.Money, req.CategoryID, keywordID, req.PaymentMethodID, req.Memo)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "데이터 삽입 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "지출 데이터가 성공적으로 저장되었습니다.",
	}

	utils.SendCreatedResponse(w, response)
}

// 새로운 구조의 지출 데이터 조회 핸들러 (특정 날짜)
func (h *OutAccountHandler) GetOutAccountByDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "날짜가 지정되지 않았습니다.")
		return
	}

	data, err := h.DB.GetOutAccountsByDate(date)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "데이터 조회 중 오류 발생")
		return
	}

	utils.SendSuccessResponse(w, data)
}

// 새로운 구조의 월별 지출 데이터 조회 핸들러
func (h *OutAccountHandler) GetOutAccountByMonthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	year := r.URL.Query().Get("year")
	month := r.URL.Query().Get("month")

	if year == "" || month == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "년도와 월이 필요합니다.")
		return
	}

	outAccounts, err := h.DB.GetOutAccountsForMonth(year, month)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "지출 데이터 조회 중 오류 발생")
		return
	}

	utils.SendSuccessResponse(w, outAccounts)
}

// 새로운 구조의 지출 데이터 업데이트 핸들러
func (h *OutAccountHandler) UpdateOutAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	var req struct {
		UUID            string `json:"uuid"`
		Date            string `json:"date"`
		User            string `json:"user"`
		Money           int    `json:"money"`
		CategoryID      int    `json:"category_id"`
		KeywordName     string `json:"keyword_name,omitempty"`
		PaymentMethodID int    `json:"payment_method_id"`
		Memo            string `json:"memo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 데이터입니다.")
		return
	}

	// 입력 검증
	if req.UUID == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "UUID는 필수입니다.")
		return
	}
	if req.Date == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "날짜는 필수입니다.")
		return
	}
	if req.User == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "사용자는 필수입니다.")
		return
	}
	if req.Money <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "금액은 0보다 커야 합니다.")
		return
	}
	if req.CategoryID <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "카테고리를 선택해주세요.")
		return
	}
	if req.PaymentMethodID <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "결제수단을 선택해주세요.")
		return
	}

	// 키워드 처리 (있는 경우)
	var keywordID *int
	if req.KeywordName != "" {
		id, err := h.KeywordDB.UpsertKeyword(req.CategoryID, req.KeywordName)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "키워드 처리 중 오류 발생")
			return
		}
		keywordIDValue := int(id)
		keywordID = &keywordIDValue
	}

	// 지출 데이터 업데이트
	err := h.DB.UpdateOutAccount(req.UUID, req.Date, req.User, req.Money, req.CategoryID, keywordID, req.PaymentMethodID, req.Memo)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "데이터 업데이트 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "지출 데이터가 성공적으로 업데이트되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}

// 지출 데이터 삭제 핸들러
func (h *OutAccountHandler) DeleteOutAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "UUID가 필요합니다.")
		return
	}

	err := h.DB.DeleteOutAccount(uuid)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "데이터 삭제 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "지출 데이터가 성공적으로 삭제되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}
