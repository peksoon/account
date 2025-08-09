package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 구조체 정의
type ErrorCode struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Error 인터페이스 구현
func (e ErrorCode) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// 공통 에러 코드 정의
var (
	// 시스템 에러
	ErrInternalServer = ErrorCode{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "내부 서버 오류가 발생했습니다",
		Status:  http.StatusInternalServerError,
	}

	ErrDatabaseConnection = ErrorCode{
		Code:    "DATABASE_CONNECTION_ERROR",
		Message: "데이터베이스 연결에 실패했습니다",
		Status:  http.StatusInternalServerError,
	}

	// 요청 관련 에러
	ErrInvalidRequest = ErrorCode{
		Code:    "INVALID_REQUEST",
		Message: "잘못된 요청입니다",
		Status:  http.StatusBadRequest,
	}

	ErrInvalidJSON = ErrorCode{
		Code:    "INVALID_JSON",
		Message: "JSON 형식이 올바르지 않습니다",
		Status:  http.StatusBadRequest,
	}

	ErrMissingRequired = ErrorCode{
		Code:    "MISSING_REQUIRED_FIELD",
		Message: "필수 필드가 누락되었습니다",
		Status:  http.StatusBadRequest,
	}

	// 데이터 관련 에러
	ErrNotFound = ErrorCode{
		Code:    "NOT_FOUND",
		Message: "요청한 데이터를 찾을 수 없습니다",
		Status:  http.StatusNotFound,
	}

	ErrAlreadyExists = ErrorCode{
		Code:    "ALREADY_EXISTS",
		Message: "이미 존재하는 데이터입니다",
		Status:  http.StatusConflict,
	}

	ErrInvalidData = ErrorCode{
		Code:    "INVALID_DATA",
		Message: "잘못된 데이터입니다",
		Status:  http.StatusBadRequest,
	}

	// 계좌 관련 에러
	ErrAccountNotFound = ErrorCode{
		Code:    "ACCOUNT_NOT_FOUND",
		Message: "계좌 정보를 찾을 수 없습니다",
		Status:  http.StatusNotFound,
	}

	ErrInvalidAccountData = ErrorCode{
		Code:    "INVALID_ACCOUNT_DATA",
		Message: "계좌 정보가 올바르지 않습니다",
		Status:  http.StatusBadRequest,
	}

	// 카테고리 관련 에러
	ErrCategoryNotFound = ErrorCode{
		Code:    "CATEGORY_NOT_FOUND",
		Message: "카테고리를 찾을 수 없습니다",
		Status:  http.StatusNotFound,
	}

	ErrInvalidCategoryData = ErrorCode{
		Code:    "INVALID_CATEGORY_DATA",
		Message: "카테고리 정보가 올바르지 않습니다",
		Status:  http.StatusBadRequest,
	}

	// 결제수단 관련 에러
	ErrPaymentMethodNotFound = ErrorCode{
		Code:    "PAYMENT_METHOD_NOT_FOUND",
		Message: "결제수단을 찾을 수 없습니다",
		Status:  http.StatusNotFound,
	}

	ErrInvalidPaymentMethodData = ErrorCode{
		Code:    "INVALID_PAYMENT_METHOD_DATA",
		Message: "결제수단 정보가 올바르지 않습니다",
		Status:  http.StatusBadRequest,
	}

	// 키워드 관련 에러
	ErrKeywordNotFound = ErrorCode{
		Code:    "KEYWORD_NOT_FOUND",
		Message: "키워드를 찾을 수 없습니다",
		Status:  http.StatusNotFound,
	}

	ErrInvalidKeywordData = ErrorCode{
		Code:    "INVALID_KEYWORD_DATA",
		Message: "키워드 정보가 올바르지 않습니다",
		Status:  http.StatusBadRequest,
	}

	// 은행계좌 관련 에러
	ErrBankAccountNotFound = ErrorCode{
		Code:    "BANK_ACCOUNT_NOT_FOUND",
		Message: "은행계좌를 찾을 수 없습니다",
		Status:  http.StatusNotFound,
	}

	ErrInvalidBankAccountData = ErrorCode{
		Code:    "INVALID_BANK_ACCOUNT_DATA",
		Message: "은행계좌 정보가 올바르지 않습니다",
		Status:  http.StatusBadRequest,
	}

	// 통계 관련 에러
	ErrStatisticsCalculation = ErrorCode{
		Code:    "STATISTICS_CALCULATION_ERROR",
		Message: "통계 계산 중 오류가 발생했습니다",
		Status:  http.StatusInternalServerError,
	}

	ErrInvalidDateRange = ErrorCode{
		Code:    "INVALID_DATE_RANGE",
		Message: "잘못된 날짜 범위입니다",
		Status:  http.StatusBadRequest,
	}
)

// ErrorResponse 응답 구조체
type ErrorResponse struct {
	Error ErrorCode `json:"error"`
}

// NewErrorResponse 에러 응답 생성
func NewErrorResponse(err ErrorCode) ErrorResponse {
	return ErrorResponse{
		Error: err,
	}
}

// CustomError 커스텀 에러 생성
func CustomError(code, message string, status int) ErrorCode {
	return ErrorCode{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// WithMessage 기존 에러 코드에 메시지 추가
func (e ErrorCode) WithMessage(message string) ErrorCode {
	return ErrorCode{
		Code:    e.Code,
		Message: message,
		Status:  e.Status,
	}
}

// WithDetails 기존 에러 코드에 상세 정보 추가
func (e ErrorCode) WithDetails(details string) ErrorCode {
	return ErrorCode{
		Code:    e.Code,
		Message: fmt.Sprintf("%s: %s", e.Message, details),
		Status:  e.Status,
	}
}
