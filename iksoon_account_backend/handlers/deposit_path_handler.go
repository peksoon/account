package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"iksoon_account_backend/models"
	"iksoon_account_backend/utils"
)

type DepositPathHandler struct {
	DB DepositPathRepository
}

type DepositPathRepository interface {
	GetDepositPaths() ([]models.DepositPath, error)
	CreateDepositPath(name string) (int64, error)
	UpdateDepositPath(id int, name string) error
	DeleteDepositPath(id int) error
	ForceDeleteDepositPath(id int) error
	CheckDepositPathExists(id int) (bool, error)
	CheckDepositPathUsage(depositPathID int) (bool, error)
}

// 입금경로 목록 조회 핸들러
func (h *DepositPathHandler) GetDepositPathsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	depositPaths, err := h.DB.GetDepositPaths()
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "입금경로 조회 중 오류 발생")
		return
	}

	utils.SendSuccessResponse(w, depositPaths)
}

// 입금경로 생성 핸들러
func (h *DepositPathHandler) CreateDepositPathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	var req models.DepositPathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 데이터입니다.")
		return
	}

	// 입력 검증
	if req.Name == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "입금경로 이름은 필수입니다.")
		return
	}

	depositPathID, err := h.DB.CreateDepositPath(req.Name)
	if err != nil {
		if strings.Contains(err.Error(), "이미 존재하는") {
			utils.SendErrorResponse(w, http.StatusConflict, models.ErrCodeDuplicateEntry, err.Error())
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "입금경로 생성 중 오류 발생")
		return
	}

	response := map[string]interface{}{
		"id":      depositPathID,
		"message": "입금경로가 성공적으로 생성되었습니다.",
	}

	utils.SendCreatedResponse(w, response)
}

// 입금경로 수정 핸들러
func (h *DepositPathHandler) UpdateDepositPathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "입금경로 ID가 필요합니다.")
		return
	}

	depositPathID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 입금경로 ID를 입력해주세요.")
		return
	}

	var req models.DepositPathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 데이터입니다.")
		return
	}

	if req.Name == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "입금경로 이름은 필수입니다.")
		return
	}

	err = h.DB.UpdateDepositPath(depositPathID, req.Name)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, "존재하지 않는 입금경로입니다.")
			return
		}
		if strings.Contains(err.Error(), "이미 존재하는") {
			utils.SendErrorResponse(w, http.StatusConflict, models.ErrCodeDuplicateEntry, err.Error())
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "입금경로 수정 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "입금경로가 성공적으로 수정되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}

// 입금경로 삭제 핸들러
func (h *DepositPathHandler) DeleteDepositPathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "입금경로 ID가 필요합니다.")
		return
	}

	depositPathID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 입금경로 ID를 입력해주세요.")
		return
	}

	// 입금경로를 사용하는 데이터가 있는지 확인
	hasData, err := h.DB.CheckDepositPathUsage(depositPathID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "입금경로 사용 여부 확인 중 오류 발생")
		return
	}

	if hasData {
		utils.SendErrorResponse(w, http.StatusConflict, models.ErrCodeCannotDelete, "이 입금경로를 사용하는 데이터가 존재합니다.")
		return
	}

	err = h.DB.DeleteDepositPath(depositPathID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, "존재하지 않는 입금경로입니다.")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "입금경로 삭제 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "입금경로가 성공적으로 삭제되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}

// 입금경로 강제 삭제 핸들러
func (h *DepositPathHandler) ForceDeleteDepositPathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "입금경로 ID가 필요합니다.")
		return
	}

	depositPathID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 입금경로 ID를 입력해주세요.")
		return
	}

	err = h.DB.ForceDeleteDepositPath(depositPathID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, "존재하지 않는 입금경로입니다.")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "입금경로 강제 삭제 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "입금경로가 성공적으로 삭제되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}
