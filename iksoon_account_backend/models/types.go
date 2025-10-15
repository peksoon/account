package models

import "time"

// User 구조체 - 사용자 관리
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Category 구조체 - 카테고리 관리
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`         // 'out' 또는 'in'
	ExpenseType string    `json:"expense_type"` // 'fixed' 또는 'variable' (지출 카테고리만 해당)
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Keyword 구조체 - 키워드 관리
type Keyword struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	Name       string    `json:"name"`
	UsageCount int       `json:"usage_count"`
	LastUsed   time.Time `json:"last_used"`
	CreatedAt  time.Time `json:"created_at"`
}

// PaymentMethod 구조체 - 결제수단 관리
type PaymentMethod struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	ParentID  *int            `json:"parent_id"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Children  []PaymentMethod `json:"children,omitempty"`
}

// DepositPath 구조체 - 입금 경로 관리
type DepositPath struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OutAccount 구조체 - 지출 데이터
type OutAccount struct {
	UUID              string `json:"uuid"`
	Date              string `json:"date"`
	User              string `json:"user"`
	Money             int    `json:"money"`
	CategoryID        int    `json:"category_id"`
	CategoryName      string `json:"category_name,omitempty"`
	KeywordID         *int   `json:"keyword_id,omitempty"`
	KeywordName       string `json:"keyword_name,omitempty"`
	PaymentMethodID   int    `json:"payment_method_id"`
	PaymentMethodName string `json:"payment_method_name,omitempty"`
	Memo              string `json:"memo"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// InAccount 구조체 - 수입 데이터
type InAccount struct {
	UUID            string `json:"uuid"`
	Date            string `json:"date"`
	User            string `json:"user"`
	Money           int    `json:"money"`
	CategoryID      int    `json:"category_id"`
	CategoryName    string `json:"category_name,omitempty"`
	KeywordID       *int   `json:"keyword_id,omitempty"`
	KeywordName     string `json:"keyword_name,omitempty"`
	DepositPathID   int    `json:"deposit_path_id"`
	DepositPathName string `json:"deposit_path_name,omitempty"`
	Memo            string `json:"memo"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// KeywordSuggestion 구조체 - 키워드 자동완성용
type KeywordSuggestion struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	UsageCount int    `json:"usage_count"`
}

// StatisticsRequest 구조체 - 통계 요청
type StatisticsRequest struct {
	Type      string `json:"type"`       // 'week', 'month', 'year', 'custom', 'all'
	StartDate string `json:"start_date"` // 'custom' 타입일 때 사용
	EndDate   string `json:"end_date"`   // 'custom' 타입일 때 사용
	Category  string `json:"category"`   // 'out' 또는 'in'
}

// CategoryStatistics 구조체 - 카테고리별 통계
type CategoryStatistics struct {
	CategoryID   int                 `json:"category_id"`
	CategoryName string              `json:"category_name"`
	TotalAmount  int                 `json:"total_amount"`
	Percentage   float64             `json:"percentage"`
	Count        int                 `json:"count"`
	Keywords     []KeywordStatistics `json:"keywords,omitempty"`
}

// KeywordStatistics 구조체 - 키워드별 통계
type KeywordStatistics struct {
	KeywordID   int     `json:"keyword_id"`
	KeywordName string  `json:"keyword_name"`
	TotalAmount int     `json:"total_amount"`
	Percentage  float64 `json:"percentage"`
	Count       int     `json:"count"`
}

// PaymentMethodStatistics 구조체 - 결제수단별 통계
type PaymentMethodStatistics struct {
	PaymentMethodID   int     `json:"payment_method_id"`
	PaymentMethodName string  `json:"payment_method_name"`
	TotalAmount       int     `json:"total_amount"`
	Percentage        float64 `json:"percentage"`
	Count             int     `json:"count"`
}

// StatisticsResponse 구조체 - 통계 응답 (기준치 정보 포함)
type StatisticsResponse struct {
	Period         string                    `json:"period"`
	TotalAmount    int                       `json:"total_amount"`
	TotalCount     int                       `json:"total_count"`
	Categories     []CategoryStatistics      `json:"categories"`
	TopCategory    *CategoryStatistics       `json:"top_category,omitempty"`
	ChartData      []ChartData               `json:"chart_data"`
	BudgetUsages   []BudgetUsage             `json:"budget_usages,omitempty"`   // 기준치 사용량 정보
	PaymentMethods []PaymentMethodStatistics `json:"payment_methods,omitempty"` // 결제수단 통계
}

// ChartData 구조체 - 차트 데이터
type ChartData struct {
	Label      string  `json:"label"`
	Value      int     `json:"value"`
	Percentage float64 `json:"percentage"`
	Color      string  `json:"color"`
}

// Request 구조체들
type CategoryRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	ExpenseType string `json:"expense_type"` // 'fixed' 또는 'variable' (지출 카테고리만 해당)
}

type PaymentMethodRequest struct {
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
}

type DepositPathRequest struct {
	Name string `json:"name"`
}

type BankAccountRequest struct {
	BankName   string `json:"bank_name"`
	AccountNum string `json:"account_num"`
}

// ErrorResponse 구조체 - 에러 응답
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CategoryBudget 구조체 - 카테고리별 기준치 관리
type CategoryBudget struct {
	ID            int       `json:"id"`
	CategoryID    int       `json:"category_id"`
	CategoryName  string    `json:"category_name,omitempty"`
	UserName      string    `json:"user_name"`
	MonthlyBudget int       `json:"monthly_budget"` // 월 기준치
	YearlyBudget  int       `json:"yearly_budget"`  // 연 기준치
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BudgetUsage 구조체 - 기준치 사용량 정보
type BudgetUsage struct {
	CategoryID       int     `json:"category_id"`
	CategoryName     string  `json:"category_name"`
	MonthlyBudget    int     `json:"monthly_budget"`
	YearlyBudget     int     `json:"yearly_budget"`
	MonthlyUsed      int     `json:"monthly_used"`      // 월 사용량
	YearlyUsed       int     `json:"yearly_used"`       // 연 사용량
	MonthlyRemaining int     `json:"monthly_remaining"` // 월 잔여
	YearlyRemaining  int     `json:"yearly_remaining"`  // 연 잔여
	MonthlyPercent   float64 `json:"monthly_percent"`   // 월 사용 퍼센티지
	YearlyPercent    float64 `json:"yearly_percent"`    // 연 사용 퍼센티지
	IsMonthlyOver    bool    `json:"is_monthly_over"`   // 월 기준치 초과 여부
	IsYearlyOver     bool    `json:"is_yearly_over"`    // 연 기준치 초과 여부
}

// CategoryBudgetRequest 구조체 - 기준치 요청
type CategoryBudgetRequest struct {
	CategoryID    int    `json:"category_id"`
	UserName      string `json:"user_name"`
	MonthlyBudget int    `json:"monthly_budget"`
	YearlyBudget  int    `json:"yearly_budget"`
}

// MonthlyBudgetRequest 구조체 - 월별 기준치 요청
type MonthlyBudgetRequest struct {
	CategoryID int    `json:"category_id"`
	UserName   string `json:"user_name"`
	Amount     int    `json:"amount"`
}

// YearlyBudgetRequest 구조체 - 연별 기준치 요청
type YearlyBudgetRequest struct {
	CategoryID int    `json:"category_id"`
	UserName   string `json:"user_name"`
	Amount     int    `json:"amount"`
}

// OutAccountWithBudget 구조체 - 기준치 정보 포함 지출 응답
type OutAccountWithBudget struct {
	Message     string       `json:"message"`
	BudgetUsage *BudgetUsage `json:"budget_usage,omitempty"`
}

// 에러 코드 상수
const (
	ErrCodeInvalidInput    = "INVALID_INPUT"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeDuplicateEntry  = "DUPLICATE_ENTRY"
	ErrCodeForeignKeyError = "FOREIGN_KEY_ERROR"
	ErrCodeDatabaseError   = "DATABASE_ERROR"
	ErrCodeCannotDelete    = "CANNOT_DELETE"
	ErrCodeInternalError   = "INTERNAL_ERROR"
)
