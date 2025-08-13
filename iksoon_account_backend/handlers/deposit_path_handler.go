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

// GetDepositPathsHandler 입금경로 목록 조회 핸들러
func (h *DepositPathHandler) GetDepositPathsHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateHTTPMethod(w, r, http.MethodGet) {
		return
	}

	utils.Debug("입금경로 목록 조회 요청")

	depositPaths, err := h.DB.GetDepositPaths()
	if err != nil {
		utils.LogDatabaseError("입금경로 목록 조회", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("입금경로 목록 조회 실패"))
		return
	}

	utils.Debug("입금경로 목록 조회 성공: %d개", len(depositPaths))
	utils.SendSuccessResponse(w, depositPaths)
}

// CreateDepositPathHandler 입금경로 생성 핸들러
func (h *DepositPathHandler) CreateDepositPathHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateHTTPMethod(w, r, http.MethodPost) {
		return
	}

	var req models.DepositPathRequest
	if !utils.ValidateJSONRequest(w, r, &req) {
		return
	}

	utils.Debug("입금경로 생성 요청: %+v", req)

	// 입력 검증
	if req.Name == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("입금경로 이름은 필수입니다"))
		return
	}

	depositPathID, err := h.DB.CreateDepositPath(req.Name)
	if err != nil {
		if strings.Contains(err.Error(), "이미 존재하는") {
			utils.SendError(w, apiErrors.ErrAlreadyExists.WithMessage(err.Error()))
			return
		}
		utils.LogDatabaseError("입금경로 생성", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("입금경로 생성 실패"))
		return
	}

	response := map[string]interface{}{
		"id":      depositPathID,
		"message": "입금경로가 성공적으로 생성되었습니다.",
	}

	utils.Debug("입금경로 생성 성공: ID %d", depositPathID)
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
