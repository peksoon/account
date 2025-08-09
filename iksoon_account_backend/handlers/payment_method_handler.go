package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"iksoon_account_backend/models"
	"iksoon_account_backend/utils"
)

type PaymentMethodHandler struct {
	DB PaymentMethodRepository
}

type PaymentMethodRepository interface {
	GetPaymentMethods() ([]models.PaymentMethod, error)
	CreatePaymentMethod(name string, parentID *int) (int64, error)
	UpdatePaymentMethod(id int, name string) error
	DeletePaymentMethod(id int) error
	ForceDeletePaymentMethod(id int) error
	CheckPaymentMethodExists(id int) (bool, error)
	CheckPaymentMethodUsage(paymentMethodID int) (bool, error)
}

// 결제수단 목록 조회 핸들러
func (h *PaymentMethodHandler) GetPaymentMethodsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	paymentMethods, err := h.DB.GetPaymentMethods()
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "결제수단 조회 중 오류 발생")
		return
	}

	utils.SendSuccessResponse(w, paymentMethods)
}

// 결제수단 생성 핸들러
func (h *PaymentMethodHandler) CreatePaymentMethodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	var req models.PaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 데이터입니다.")
		return
	}

	// 입력 검증
	if req.Name == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "결제수단 이름은 필수입니다.")
		return
	}

	paymentMethodID, err := h.DB.CreatePaymentMethod(req.Name, req.ParentID)
	if err != nil {
		if strings.Contains(err.Error(), "이미 존재하는") {
			utils.SendErrorResponse(w, http.StatusConflict, models.ErrCodeDuplicateEntry, err.Error())
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "결제수단 생성 중 오류 발생")
		return
	}

	response := map[string]interface{}{
		"id":      paymentMethodID,
		"message": "결제수단이 성공적으로 생성되었습니다.",
	}

	utils.SendCreatedResponse(w, response)
}

// 결제수단 수정 핸들러
func (h *PaymentMethodHandler) UpdatePaymentMethodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "결제수단 ID가 필요합니다.")
		return
	}

	paymentMethodID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 결제수단 ID를 입력해주세요.")
		return
	}

	var req models.PaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "잘못된 요청 데이터입니다.")
		return
	}

	if req.Name == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "결제수단 이름은 필수입니다.")
		return
	}

	err = h.DB.UpdatePaymentMethod(paymentMethodID, req.Name)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, "존재하지 않는 결제수단입니다.")
			return
		}
		if strings.Contains(err.Error(), "이미 존재하는") {
			utils.SendErrorResponse(w, http.StatusConflict, models.ErrCodeDuplicateEntry, err.Error())
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "결제수단 수정 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "결제수단이 성공적으로 수정되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}

// 결제수단 삭제 핸들러
func (h *PaymentMethodHandler) DeletePaymentMethodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "결제수단 ID가 필요합니다.")
		return
	}

	paymentMethodID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 결제수단 ID를 입력해주세요.")
		return
	}

	// 결제수단을 사용하는 데이터가 있는지 확인
	hasData, err := h.DB.CheckPaymentMethodUsage(paymentMethodID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "결제수단 사용 여부 확인 중 오류 발생")
		return
	}

	if hasData {
		utils.SendErrorResponse(w, http.StatusConflict, models.ErrCodeCannotDelete, "이 결제수단을 사용하는 데이터가 존재합니다.")
		return
	}

	err = h.DB.DeletePaymentMethod(paymentMethodID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, "존재하지 않는 결제수단입니다.")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "결제수단 삭제 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "결제수단이 성공적으로 삭제되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}

// 결제수단 강제 삭제 핸들러
func (h *PaymentMethodHandler) ForceDeletePaymentMethodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "결제수단 ID가 필요합니다.")
		return
	}

	paymentMethodID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 결제수단 ID를 입력해주세요.")
		return
	}

	err = h.DB.ForceDeletePaymentMethod(paymentMethodID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			utils.SendErrorResponse(w, http.StatusNotFound, models.ErrCodeNotFound, "존재하지 않는 결제수단입니다.")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "결제수단 강제 삭제 중 오류 발생")
		return
	}

	response := map[string]string{
		"message": "결제수단이 성공적으로 삭제되었습니다.",
	}

	utils.SendSuccessResponse(w, response)
}
