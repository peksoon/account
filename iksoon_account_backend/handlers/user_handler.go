package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	apiErrors "iksoon_account_backend/errors"
	"iksoon_account_backend/models"
	"iksoon_account_backend/utils"
)

type UserHandler struct {
	DB UserRepository
}

type UserRepository interface {
	GetUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(name, email string) (int64, error)
	UpdateUser(id int, name, email string) error
	DeleteUser(id int) error
	ForceDeleteUser(id int) error
	CheckUserUsage(userID int) (bool, error)
}

// GetUsersHandler 사용자 목록 조회 핸들러
func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendError(w, apiErrors.ErrInvalidRequest.WithMessage("지원되지 않는 메소드입니다"))
		return
	}

	utils.Debug("사용자 목록 조회 요청")

	users, err := h.DB.GetUsers()
	if err != nil {
		utils.LogDatabaseError("사용자 목록 조회", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("사용자 목록 조회 실패"))
		return
	}

	utils.Debug("사용자 목록 조회 성공: %d개", len(users))
	utils.SendSuccessResponse(w, users)
}

// CreateUserHandler 사용자 생성 핸들러
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendError(w, apiErrors.ErrInvalidRequest.WithMessage("지원되지 않는 메소드입니다"))
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LogError("JSON 디코드", err)
		utils.SendError(w, apiErrors.ErrInvalidJSON)
		return
	}

	utils.Debug("사용자 생성 요청: %+v", req)

	// 입력 검증
	if req.Name == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("사용자 이름은 필수입니다"))
		return
	}

	userID, err := h.DB.CreateUser(req.Name, req.Email)
	if err != nil {
		if err.Error() == "이미 존재하는 사용자 이름입니다" {
			utils.SendError(w, apiErrors.ErrAlreadyExists.WithMessage("이미 존재하는 사용자 이름입니다"))
			return
		}
		utils.LogDatabaseError("사용자 생성", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("사용자 생성 실패"))
		return
	}

	// 생성된 사용자 조회
	user, err := h.DB.GetUserByID(int(userID))
	if err != nil {
		utils.LogDatabaseError("생성된 사용자 조회", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("생성된 사용자 조회 실패"))
		return
	}

	utils.Debug("사용자 생성 성공: ID %d", userID)
	utils.SendCreatedResponse(w, user)
}

// UpdateUserHandler 사용자 수정 핸들러
func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.SendError(w, apiErrors.ErrInvalidRequest.WithMessage("지원되지 않는 메소드입니다"))
		return
	}

	var req struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LogError("JSON 디코드", err)
		utils.SendError(w, apiErrors.ErrInvalidJSON)
		return
	}

	utils.Debug("사용자 수정 요청: %+v", req)

	// 입력 검증
	if req.ID <= 0 {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("사용자 ID는 필수입니다"))
		return
	}
	if req.Name == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("사용자 이름은 필수입니다"))
		return
	}

	err := h.DB.UpdateUser(req.ID, req.Name, req.Email)
	if err != nil {
		if err.Error() == "이미 존재하는 사용자 이름입니다" {
			utils.SendError(w, apiErrors.ErrAlreadyExists.WithMessage("이미 존재하는 사용자 이름입니다"))
			return
		}
		if err.Error() == "사용자를 찾을 수 없습니다" {
			utils.SendError(w, apiErrors.ErrNotFound.WithMessage("사용자를 찾을 수 없습니다"))
			return
		}
		utils.LogDatabaseError("사용자 수정", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("사용자 수정 실패"))
		return
	}

	// 수정된 사용자 조회
	user, err := h.DB.GetUserByID(req.ID)
	if err != nil {
		utils.LogDatabaseError("수정된 사용자 조회", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("수정된 사용자 조회 실패"))
		return
	}

	utils.Debug("사용자 수정 성공: ID %d", req.ID)
	utils.SendSuccessResponse(w, user)
}

// DeleteUserHandler 사용자 삭제 핸들러
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendError(w, apiErrors.ErrInvalidRequest.WithMessage("지원되지 않는 메소드입니다"))
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("사용자 ID는 필수입니다"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("올바르지 않은 사용자 ID입니다"))
		return
	}

	utils.Debug("사용자 삭제 요청: ID %d", id)

	err = h.DB.DeleteUser(id)
	if err != nil {
		if err.Error() == "사용 중인 사용자는 삭제할 수 없습니다" {
			utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("사용 중인 사용자는 삭제할 수 없습니다"))
			return
		}
		if err.Error() == "사용자를 찾을 수 없습니다" {
			utils.SendError(w, apiErrors.ErrNotFound.WithMessage("사용자를 찾을 수 없습니다"))
			return
		}
		utils.LogDatabaseError("사용자 삭제", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("사용자 삭제 실패"))
		return
	}

	utils.Debug("사용자 삭제 성공: ID %d", id)
	utils.SendSuccessResponse(w, map[string]string{"message": "사용자가 성공적으로 삭제되었습니다"})
}

// ForceDeleteUserHandler 사용자 강제 삭제 핸들러
func (h *UserHandler) ForceDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendError(w, apiErrors.ErrInvalidRequest.WithMessage("지원되지 않는 메소드입니다"))
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("사용자 ID는 필수입니다"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("올바르지 않은 사용자 ID입니다"))
		return
	}

	utils.Debug("사용자 강제 삭제 요청: ID %d", id)

	err = h.DB.ForceDeleteUser(id)
	if err != nil {
		if err.Error() == "사용자를 찾을 수 없습니다" {
			utils.SendError(w, apiErrors.ErrNotFound.WithMessage("사용자를 찾을 수 없습니다"))
			return
		}
		utils.LogDatabaseError("사용자 강제 삭제", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("사용자 강제 삭제 실패"))
		return
	}

	utils.Debug("사용자 강제 삭제 성공: ID %d", id)
	utils.SendSuccessResponse(w, map[string]string{"message": "사용자가 강제로 삭제되었습니다"})
}

// CheckUserUsageHandler 사용자 사용 여부 확인 핸들러
func (h *UserHandler) CheckUserUsageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendError(w, apiErrors.ErrInvalidRequest.WithMessage("지원되지 않는 메소드입니다"))
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendError(w, apiErrors.ErrMissingRequired.WithMessage("사용자 ID는 필수입니다"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("올바르지 않은 사용자 ID입니다"))
		return
	}

	utils.Debug("사용자 사용 여부 확인 요청: ID %d", id)

	inUse, err := h.DB.CheckUserUsage(id)
	if err != nil {
		utils.LogDatabaseError("사용자 사용 여부 확인", err)
		utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("사용자 사용 여부 확인 실패"))
		return
	}

	utils.Debug("사용자 사용 여부 확인 완료: ID %d, InUse %t", id, inUse)
	utils.SendSuccessResponse(w, map[string]bool{"in_use": inUse})
}
