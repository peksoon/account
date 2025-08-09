package handlers

import (
	"encoding/json"
	"net/http"

	"iksoon_account_backend/database"
	"iksoon_account_backend/models"
	"iksoon_account_backend/utils"
)

type InAccountHandler struct {
	DB        InAccountRepository
	KeywordDB KeywordRepository
}

type InAccountRepository interface {
	InsertInAccount(date, user string, money, categoryID int, keywordID *int, depositPathID int, memo string) error
	GetInAccountsByDate(date string) ([]models.InAccount, error)
	GetInAccountsForMonth(year, month string) ([]models.InAccount, error)
	UpdateInAccount(uuid, date, user string, money, categoryID int, keywordID *int, depositPathID int, memo string) error
	DeleteInAccount(uuid string) error
	GetInAccountByUUID(uuid string) (*models.InAccount, error)
}

// 새로운 구조의 수입 데이터 삽입 핸들러
func (h *InAccountHandler) InsertInAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	var req struct {
		Date        string `json:"date"`
		User        string `json:"user"`
		Money       int    `json:"money"`
		CategoryID  int    `json:"category_id"`
		KeywordName string `json:"keyword_name,omitempty"`
		DepositPath string `json:"deposit_path"`
		Memo        string `json:"memo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 데이터입니다.")
		return
	}

	// 입력 검증
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
	if req.DepositPath == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "입금 경로를 선택해주세요.")
		return
	}

	// 입금 경로에서 ID 찾기
	var depositPathID int
	// 입금 경로 이름으로 ID 찾기 (DepositPath repository를 사용해야 함)
	// 임시로 간단한 쿼리 사용
	err := h.DB.(*database.DB).Conn.QueryRow("SELECT id FROM deposit_paths WHERE name = ? AND is_active = 1", req.DepositPath).Scan(&depositPathID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "유효하지 않은 입금 경로입니다.")
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

	// 수입 데이터 삽입
	err = h.DB.InsertInAccount(req.Date, req.User, req.Money, req.CategoryID, keywordID, depositPathID, req.Memo)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "데이터 삽입 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "수입 데이터가 성공적으로 저장되었습니다.",
	}

	utils.SendCreatedResponse(w, response)
}

// 새로운 구조의 수입 데이터 조회 핸들러 (특정 날짜)
func (h *InAccountHandler) GetInAccountByDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "날짜가 지정되지 않았습니다.")
		return
	}

	data, err := h.DB.GetInAccountsByDate(date)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "데이터 조회 중 오류 발생")
		return
	}

	utils.SendSuccessResponse(w, data)
}

// 새로운 구조의 월별 수입 데이터 조회 핸들러
func (h *InAccountHandler) GetInAccountByMonthHandler(w http.ResponseWriter, r *http.Request) {
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

	inAccounts, err := h.DB.GetInAccountsForMonth(year, month)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "수입 데이터 조회 중 오류 발생")
		return
	}

	utils.SendSuccessResponse(w, inAccounts)
}

// 새로운 구조의 수입 데이터 업데이트 핸들러
func (h *InAccountHandler) UpdateInAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	var req struct {
		UUID        string `json:"uuid"`
		Date        string `json:"date"`
		User        string `json:"user"`
		Money       int    `json:"money"`
		CategoryID  int    `json:"category_id"`
		KeywordName string `json:"keyword_name,omitempty"`
		DepositPath string `json:"deposit_path"`
		Memo        string `json:"memo"`
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
	if req.DepositPath == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "입금 경로를 선택해주세요.")
		return
	}

	// 입금 경로에서 ID 찾기
	var depositPathID int
	err := h.DB.(*database.DB).Conn.QueryRow("SELECT id FROM deposit_paths WHERE name = ? AND is_active = 1", req.DepositPath).Scan(&depositPathID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "유효하지 않은 입금 경로입니다.")
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

	// 수입 데이터 업데이트
	err = h.DB.UpdateInAccount(req.UUID, req.Date, req.User, req.Money, req.CategoryID, keywordID, depositPathID, req.Memo)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "데이터 업데이트 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "수입 데이터가 성공적으로 업데이트되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}

// 수입 데이터 삭제 핸들러
func (h *InAccountHandler) DeleteInAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "UUID가 필요합니다.")
		return
	}

	err := h.DB.DeleteInAccount(uuid)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "데이터 삭제 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "수입 데이터가 성공적으로 삭제되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}
