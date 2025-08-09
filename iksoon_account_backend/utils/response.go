package utils

import (
	"encoding/json"
	"net/http"
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
