package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	apiErrors "iksoon_account_backend/errors"
	"iksoon_account_backend/models"
)

// SendErrorResponse 에러 응답 전송 헬퍼 함수 (기존 호환성 유지)
func SendErrorResponse(w http.ResponseWriter, statusCode int, errorCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := models.ErrorResponse{
		Code:    errorCode,
		Message: message,
	}

	Error("Error Response: %s - %s", errorCode, message)
	json.NewEncoder(w).Encode(errorResponse)
}

// SendError 새로운 에러 코드 시스템을 사용한 에러 응답 전송
func SendError(w http.ResponseWriter, err apiErrors.ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)

	errorResponse := apiErrors.NewErrorResponse(err)
	Error("API Error: %s", err.Error())
	json.NewEncoder(w).Encode(errorResponse)
}

// SendSuccessResponse 성공 응답 전송 헬퍼 함수
func SendSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if IsDebugEnabled() {
		Debug("Success Response sent")
	}
	json.NewEncoder(w).Encode(data)
}

// SendCreatedResponse 생성 성공 응답 전송 헬퍼 함수
func SendCreatedResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if IsDebugEnabled() {
		Debug("Created Response sent")
	}
	json.NewEncoder(w).Encode(data)
}

// LogHTTPMiddleware HTTP 요청/응답 로깅 미들웨어
func LogHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 요청 로깅
		LogHTTPRequest(r.Method, r.URL.Path, r.RemoteAddr)

		// ResponseWriter 래핑하여 상태 코드 캡처
		wrapped := &responseWriter{ResponseWriter: w}

		// 다음 핸들러 실행
		next.ServeHTTP(wrapped, r)

		// 응답 로깅
		duration := time.Since(start)
		LogHTTPResponse(r.Method, r.URL.Path, wrapped.statusCode, duration)
	})
}

// responseWriter 상태 코드를 캡처하기 위한 래퍼
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = 200
	}
	return rw.ResponseWriter.Write(b)
}

// ValidateHTTPMethod HTTP 메소드 유효성 검증
func ValidateHTTPMethod(w http.ResponseWriter, r *http.Request, expectedMethod string) bool {
	if r.Method != expectedMethod {
		SendError(w, apiErrors.ErrInvalidRequest.WithMessage("지원되지 않는 메소드입니다"))
		return false
	}
	return true
}

// ParseIDFromQuery URL 쿼리에서 ID 파라미터를 파싱하고 검증
func ParseIDFromQuery(w http.ResponseWriter, r *http.Request, paramName string) (int, bool) {
	idStr := r.URL.Query().Get(paramName)
	if idStr == "" {
		SendError(w, apiErrors.ErrMissingRequired.WithMessage(paramName+" ID는 필수입니다"))
		return 0, false
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		SendError(w, apiErrors.ErrInvalidData.WithMessage("올바르지 않은 "+paramName+" ID입니다"))
		return 0, false
	}

	return id, true
}

// ValidateJSONRequest JSON 요청 데이터를 디코드하고 검증
func ValidateJSONRequest(w http.ResponseWriter, r *http.Request, req interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		LogError("JSON 디코드", err)
		SendError(w, apiErrors.ErrInvalidJSON)
		return false
	}
	return true
}

// CreateSuccessMessage 성공 메시지 응답 생성
func CreateSuccessMessage(message string) map[string]string {
	return map[string]string{"message": message}
}
