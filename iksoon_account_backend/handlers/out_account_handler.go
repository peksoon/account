package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"iksoon_account_backend/database"
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
	GetOutAccountsByDateRange(startDate, endDate string) ([]models.OutAccount, error)
	GetOutAccountsByPaymentMethod(paymentMethodID int, startDate, endDate string) ([]models.OutAccount, error)
	GetOutAccountsByUser(userName, startDate, endDate string) ([]models.OutAccount, error)
	SearchOutAccountsByKeyword(keyword, startDate, endDate string) ([]models.OutAccount, error)
	UpdateOutAccount(uuid, date, user string, money, categoryID int, keywordID *int, paymentMethodID int, memo string) error
	DeleteOutAccount(uuid string) error
	GetOutAccountByUUID(uuid string) (*models.OutAccount, error)
	GetBudgetUsage(categoryID int, userName string, currentDate time.Time) (*models.BudgetUsage, error)
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

	// 외래키 참조 데이터 존재 여부 검증
	if err := h.validateOutAccountReferences(req.CategoryID, req.PaymentMethodID); err != nil {
		utils.LogError("외래키 검증", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, err.Error())
		return
	}

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

// InsertOutAccountWithBudgetHandler 지출 데이터 삽입 후 기준치 정보 반환
func (h *OutAccountHandler) InsertOutAccountWithBudgetHandler(w http.ResponseWriter, r *http.Request) {
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

	utils.Debug("기준치 포함 지출 데이터 삽입 요청: %+v", req)

	// 외래키 참조 데이터 존재 여부 검증
	if err := h.validateOutAccountReferences(req.CategoryID, req.PaymentMethodID); err != nil {
		utils.LogError("외래키 검증", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, err.Error())
		return
	}

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

	// 기준치 정보 조회
	parsedDate, err := utils.ParseDateTimeKST(req.Date)
	if err != nil {
		parsedDate = time.Now()
	}

	budgetUsage, err := h.DB.GetBudgetUsage(req.CategoryID, req.User, parsedDate)
	if err != nil {
		// 기준치 조회 오류는 무시하고 성공 메시지만 반환
		utils.LogError("기준치 조회", err)
		budgetUsage = nil
	}

	response := models.OutAccountWithBudget{
		Message:     "지출 데이터가 성공적으로 저장되었습니다.",
		BudgetUsage: budgetUsage,
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

// GetOutAccountsByDateRangeHandler 기간별 지출 데이터 조회 핸들러
func (h *OutAccountHandler) GetOutAccountsByDateRangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "시작일과 종료일이 필요합니다.")
		return
	}

	outAccounts, err := h.DB.GetOutAccountsByDateRange(startDate, endDate)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "지출 데이터 조회 중 오류 발생")
		return
	}

	utils.SendSuccessResponse(w, outAccounts)
}

// SearchOutAccountsByKeywordHandler 키워드로 지출 데이터 검색 핸들러
func (h *OutAccountHandler) SearchOutAccountsByKeywordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	keyword := r.URL.Query().Get("keyword")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if keyword == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "키워드가 필요합니다.")
		return
	}

	if startDate == "" || endDate == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "시작일과 종료일이 필요합니다.")
		return
	}

	outAccounts, err := h.DB.SearchOutAccountsByKeyword(keyword, startDate, endDate)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "키워드 검색 중 오류 발생")
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
		utils.LogError("지출 업데이트 JSON 디코딩", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 데이터입니다.")
		return
	}

	utils.Debug("지출 업데이트 요청 데이터: %+v", req)

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

	// UUID 존재 여부 먼저 확인
	existingAccount, err := h.DB.GetOutAccountByUUID(req.UUID)
	if err != nil {
		utils.LogError("지출 데이터 존재 확인", err)
		utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, "해당 UUID의 지출 데이터를 찾을 수 없습니다")
		return
	}
	utils.Debug("업데이트 대상 지출 데이터 확인: UUID=%s, 기존 데이터=%+v", req.UUID, existingAccount)

	// 지출 데이터 업데이트
	err = h.DB.UpdateOutAccount(req.UUID, req.Date, req.User, req.Money, req.CategoryID, keywordID, req.PaymentMethodID, req.Memo)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, "업데이트할 지출 데이터를 찾을 수 없습니다")
			return
		}
		utils.LogError("지출 데이터 업데이트", err)
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

// validateOutAccountReferences 지출 데이터의 외래키 참조 검증
func (h *OutAccountHandler) validateOutAccountReferences(categoryID, paymentMethodID int) error {
	db := h.DB.(*database.DB)

	// 카테고리 존재 여부 확인
	var categoryExists bool
	err := db.Conn.QueryRow("SELECT 1 FROM categories WHERE id = ? AND type = 'out' AND is_active = 1", categoryID).Scan(&categoryExists)
	if err != nil {
		utils.Debug("카테고리 검증 실패: categoryID=%d, err=%v", categoryID, err)
		return fmt.Errorf("존재하지 않거나 비활성화된 지출 카테고리입니다 (ID: %d)", categoryID)
	}

	// 결제수단 존재 여부 확인
	var paymentMethodExists bool
	err = db.Conn.QueryRow("SELECT 1 FROM payment_methods WHERE id = ? AND is_active = 1", paymentMethodID).Scan(&paymentMethodExists)
	if err != nil {
		utils.Debug("결제수단 검증 실패: paymentMethodID=%d, err=%v", paymentMethodID, err)
		return fmt.Errorf("존재하지 않거나 비활성화된 결제수단입니다 (ID: %d)", paymentMethodID)
	}

	utils.Debug("외래키 검증 성공: categoryID=%d, paymentMethodID=%d", categoryID, paymentMethodID)
	return nil
}

// GetOutAccountsByPaymentMethodHandler 결제수단별 지출 내역 조회 핸들러
func (h *OutAccountHandler) GetOutAccountsByPaymentMethodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendError(w, apiErrors.ErrInvalidRequest.WithMessage("지원되지 않는 메소드입니다"))
		return
	}

	// 쿼리 파라미터 파싱
	paymentMethodIDStr := r.URL.Query().Get("payment_method_id")
	statisticsType := r.URL.Query().Get("type")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")
	weekStr := r.URL.Query().Get("week")

	if paymentMethodIDStr == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("결제수단 ID가 필요합니다"))
		return
	}

	paymentMethodID, ok := utils.ParseIDFromQuery(w, r, "payment_method_id")
	if !ok {
		return
	}

	// 기본값 설정
	if statisticsType == "" {
		statisticsType = "month"
	}

	// 날짜 범위 계산 (statistics_handler의 calculateDateRange 로직과 동일)
	now := time.Now()
	var start, end time.Time

	switch statisticsType {
	case "week":
		// 주간 통계
		if yearStr != "" && weekStr != "" {
			year, err := strconv.Atoi(yearStr)
			if err == nil {
				week, err := strconv.Atoi(weekStr)
				if err == nil && week >= 1 && week <= 53 {
					jan1 := time.Date(year, 1, 1, 0, 0, 0, 0, now.Location())
					jan1Weekday := int(jan1.Weekday())
					if jan1Weekday == 0 {
						jan1Weekday = 7
					}
					daysToFirstMonday := 8 - jan1Weekday
					if jan1Weekday == 1 {
						daysToFirstMonday = 0
					}
					firstMonday := jan1.AddDate(0, 0, daysToFirstMonday)
					start = firstMonday.AddDate(0, 0, (week-1)*7)
					end = start.AddDate(0, 0, 6)
				}
			}
		}
		if start.IsZero() {
			start = now.AddDate(0, 0, -int(now.Weekday())+1)
			end = start.AddDate(0, 0, 6)
		}
	case "month":
		// 월간 통계
		if yearStr != "" && monthStr != "" {
			year, err := strconv.Atoi(yearStr)
			if err == nil {
				month, err := strconv.Atoi(monthStr)
				if err == nil && month >= 1 && month <= 12 {
					start = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())
					end = start.AddDate(0, 1, -1)
				}
			}
		}
		if start.IsZero() {
			start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			end = start.AddDate(0, 1, -1)
		}
	case "year":
		// 연간 통계
		if yearStr != "" {
			year, err := strconv.Atoi(yearStr)
			if err == nil {
				start = time.Date(year, 1, 1, 0, 0, 0, 0, now.Location())
				end = time.Date(year, 12, 31, 23, 59, 59, 0, now.Location())
			}
		}
		if start.IsZero() {
			start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
			end = time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
		}
	case "custom":
		// 커스텀 기간
		if startDate != "" && endDate != "" {
			parsedStart, err1 := utils.ParseDateTimeKST(startDate)
			parsedEnd, err2 := utils.ParseDateTimeKST(endDate)
			if err1 == nil && err2 == nil {
				start = parsedStart
				end = parsedEnd
			}
		}
	default:
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, -1)
	}

	// 날짜 포맷팅
	calculatedStartDate := utils.FormatDateKST(start)
	calculatedEndDate := utils.FormatDateKST(end)

	// 지출 내역 조회
	accounts, err := h.DB.GetOutAccountsByPaymentMethod(paymentMethodID, calculatedStartDate, calculatedEndDate)
	if err != nil {
		utils.LogError("결제수단별 지출 내역 조회", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "DATABASE_ERROR", "결제수단별 지출 내역 조회 중 오류 발생")
		return
	}

	response := map[string]interface{}{
		"payment_method_id": paymentMethodID,
		"start_date":        calculatedStartDate,
		"end_date":          calculatedEndDate,
		"accounts":          accounts,
		"total_count":       len(accounts),
	}

	utils.SendSuccessResponse(w, response)
}

// GetOutAccountsByUserHandler 사용자별 지출 내역 조회 핸들러
func (h *OutAccountHandler) GetOutAccountsByUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendError(w, apiErrors.ErrInvalidRequest.WithMessage("지원되지 않는 메소드입니다"))
		return
	}

	// 쿼리 파라미터 파싱
	userName := r.URL.Query().Get("user_name")
	statisticsType := r.URL.Query().Get("type")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")
	weekStr := r.URL.Query().Get("week")

	if userName == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("사용자 이름이 필요합니다"))
		return
	}

	// 기본값 설정
	if statisticsType == "" {
		statisticsType = "month"
	}

	// 날짜 범위 계산
	now := time.Now()
	var start, end time.Time

	switch statisticsType {
	case "week":
		if yearStr != "" && weekStr != "" {
			year, err := strconv.Atoi(yearStr)
			if err == nil {
				week, err := strconv.Atoi(weekStr)
				if err == nil && week >= 1 && week <= 53 {
					jan1 := time.Date(year, 1, 1, 0, 0, 0, 0, now.Location())
					jan1Weekday := int(jan1.Weekday())
					if jan1Weekday == 0 {
						jan1Weekday = 7
					}
					daysToFirstMonday := 8 - jan1Weekday
					if jan1Weekday == 1 {
						daysToFirstMonday = 0
					}
					firstMonday := jan1.AddDate(0, 0, daysToFirstMonday)
					start = firstMonday.AddDate(0, 0, (week-1)*7)
					end = start.AddDate(0, 0, 6)
				}
			}
		}
		if start.IsZero() {
			start = now.AddDate(0, 0, -int(now.Weekday())+1)
			end = start.AddDate(0, 0, 6)
		}
	case "month":
		if yearStr != "" && monthStr != "" {
			year, err := strconv.Atoi(yearStr)
			if err == nil {
				month, err := strconv.Atoi(monthStr)
				if err == nil && month >= 1 && month <= 12 {
					start = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())
					end = start.AddDate(0, 1, -1)
				}
			}
		}
		if start.IsZero() {
			start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			end = start.AddDate(0, 1, -1)
		}
	case "year":
		if yearStr != "" {
			year, err := strconv.Atoi(yearStr)
			if err == nil {
				start = time.Date(year, 1, 1, 0, 0, 0, 0, now.Location())
				end = time.Date(year, 12, 31, 23, 59, 59, 0, now.Location())
			}
		}
		if start.IsZero() {
			start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
			end = time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
		}
	case "custom":
		if startDate != "" && endDate != "" {
			parsedStart, err1 := utils.ParseDateTimeKST(startDate)
			parsedEnd, err2 := utils.ParseDateTimeKST(endDate)
			if err1 == nil && err2 == nil {
				start = parsedStart
				end = parsedEnd
			}
		}
	default:
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, -1)
	}

	// 날짜 포맷팅
	calculatedStartDate := utils.FormatDateKST(start)
	calculatedEndDate := utils.FormatDateKST(end)

	// 지출 내역 조회
	accounts, err := h.DB.GetOutAccountsByUser(userName, calculatedStartDate, calculatedEndDate)
	if err != nil {
		utils.LogError("사용자별 지출 내역 조회", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "DATABASE_ERROR", "사용자별 지출 내역 조회 중 오류 발생")
		return
	}

	response := map[string]interface{}{
		"user_name":   userName,
		"start_date":  calculatedStartDate,
		"end_date":    calculatedEndDate,
		"accounts":    accounts,
		"total_count": len(accounts),
	}

	utils.SendSuccessResponse(w, response)
}
