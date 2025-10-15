package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"iksoon_account_backend/models"
	"iksoon_account_backend/utils"
)

type StatisticsHandler struct {
	DB StatisticsRepository
}

type StatisticsRepository interface {
	GetCategoryStatistics(startDate, endDate, accountType string) ([]models.CategoryStatistics, error)
	GetKeywordStatistics(categoryID int, startDate, endDate, accountType string) ([]models.KeywordStatistics, error)
	GetTotalAmount(startDate, endDate, accountType string) (int, int, error) // total, count
	GetAllBudgetUsages(userName string, currentDate time.Time) ([]models.BudgetUsage, error)
	GetPaymentMethodStatistics(startDate, endDate string) ([]models.PaymentMethodStatistics, error)
	GetPaymentMethodCategoryStatistics(paymentMethodID int, startDate, endDate string) ([]models.CategoryStatistics, error)
}

// 통계 조회 핸들러
func (h *StatisticsHandler) GetStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	// 쿼리 파라미터 파싱
	statisticsType := r.URL.Query().Get("type")  // 'week', 'month', 'year', 'custom', 'all'
	accountType := r.URL.Query().Get("category") // 'out' 또는 'in'
	userName := r.URL.Query().Get("user")        // 사용자명 (기준치 조회용)
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// 새로운 개별 기간 선택 파라미터
	yearStr := r.URL.Query().Get("year")   // 선택한 년도
	monthStr := r.URL.Query().Get("month") // 선택한 월 (1-12)
	weekStr := r.URL.Query().Get("week")   // 선택한 주차 (1-53)

	// 기본값 설정
	if accountType == "" {
		accountType = "out"
	}
	if statisticsType == "" {
		statisticsType = "month"
	}

	// 날짜 범위 계산
	calculatedStartDate, calculatedEndDate, period := h.calculateDateRange(statisticsType, startDate, endDate, yearStr, monthStr, weekStr)

	// 카테고리별 통계 조회
	categories, err := h.DB.GetCategoryStatistics(calculatedStartDate, calculatedEndDate, accountType)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "통계 조회 중 오류 발생")
		return
	}

	// 총합 계산
	totalAmount, totalCount, err := h.DB.GetTotalAmount(calculatedStartDate, calculatedEndDate, accountType)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "총합 계산 중 오류 발생")
		return
	}

	// 퍼센테지 계산
	for i := range categories {
		if totalAmount > 0 {
			categories[i].Percentage = float64(categories[i].TotalAmount) / float64(totalAmount) * 100
		}
	}

	// 차트 데이터 생성
	chartData := h.generateChartData(categories)

	// 가장 큰 카테고리 찾기
	var topCategory *models.CategoryStatistics
	if len(categories) > 0 {
		topCategory = &categories[0]
	}

	// 기준치 정보 조회 (지출 통계이고 사용자명이 있는 경우에만)
	var budgetUsages []models.BudgetUsage
	if accountType == "out" && userName != "" {
		currentDate := time.Now()
		budgetUsages, err = h.DB.GetAllBudgetUsages(userName, currentDate)
		if err != nil {
			utils.LogError("기준치 사용량 조회", err)
			// 기준치 조회 오류는 무시하고 계속 진행
			budgetUsages = nil
		}
	}

	// 결제수단별 통계 조회 (지출 통계인 경우에만)
	var paymentMethods []models.PaymentMethodStatistics
	if accountType == "out" {
		paymentMethods, err = h.DB.GetPaymentMethodStatistics(calculatedStartDate, calculatedEndDate)
		if err != nil {
			utils.LogError("결제수단 통계 조회", err)
			// 결제수단 통계 조회 오류는 무시하고 계속 진행
			paymentMethods = nil
		} else {
			// 퍼센테지 계산
			for i := range paymentMethods {
				if totalAmount > 0 {
					paymentMethods[i].Percentage = float64(paymentMethods[i].TotalAmount) / float64(totalAmount) * 100
				}
			}
		}
	}

	response := models.StatisticsResponse{
		Period:         period,
		TotalAmount:    totalAmount,
		TotalCount:     totalCount,
		Categories:     categories,
		TopCategory:    topCategory,
		ChartData:      chartData,
		BudgetUsages:   budgetUsages,
		PaymentMethods: paymentMethods,
	}

	utils.SendSuccessResponse(w, response)
}

// 카테고리별 키워드 통계 조회 핸들러
func (h *StatisticsHandler) GetCategoryKeywordStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	// 쿼리 파라미터 파싱
	categoryIDStr := r.URL.Query().Get("category_id")
	statisticsType := r.URL.Query().Get("type")
	accountType := r.URL.Query().Get("category")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// 새로운 개별 기간 선택 파라미터
	yearStr := r.URL.Query().Get("year")   // 선택한 년도
	monthStr := r.URL.Query().Get("month") // 선택한 월 (1-12)
	weekStr := r.URL.Query().Get("week")   // 선택한 주차 (1-53)

	if categoryIDStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "카테고리 ID가 필요합니다.")
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 카테고리 ID를 입력해주세요.")
		return
	}

	// 기본값 설정
	if accountType == "" {
		accountType = "out"
	}
	if statisticsType == "" {
		statisticsType = "month"
	}

	// 날짜 범위 계산
	calculatedStartDate, calculatedEndDate, period := h.calculateDateRange(statisticsType, startDate, endDate, yearStr, monthStr, weekStr)

	// 키워드별 통계 조회
	keywords, err := h.DB.GetKeywordStatistics(categoryID, calculatedStartDate, calculatedEndDate, accountType)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "키워드 통계 조회 중 오류 발생")
		return
	}

	// 카테고리 총합 계산
	var categoryTotal int
	for _, keyword := range keywords {
		categoryTotal += keyword.TotalAmount
	}

	// 퍼센테지 계산
	for i := range keywords {
		if categoryTotal > 0 {
			keywords[i].Percentage = float64(keywords[i].TotalAmount) / float64(categoryTotal) * 100
		}
	}

	// 차트 데이터 생성 (키워드용)
	chartData := make([]models.ChartData, len(keywords))
	colors := h.generateColors(len(keywords))

	for i, keyword := range keywords {
		chartData[i] = models.ChartData{
			Label:      keyword.KeywordName,
			Value:      keyword.TotalAmount,
			Percentage: keyword.Percentage,
			Color:      colors[i],
		}
	}

	response := map[string]interface{}{
		"period":         period,
		"category_id":    categoryID,
		"category_total": categoryTotal,
		"keywords":       keywords,
		"chart_data":     chartData,
	}

	utils.SendSuccessResponse(w, response)
}

// 날짜 범위 계산 헬퍼 함수
func (h *StatisticsHandler) calculateDateRange(statisticsType, startDate, endDate, yearStr, monthStr, weekStr string) (string, string, string) {
	now := time.Now()
	var start, end time.Time
	var period string

	switch statisticsType {
	case "week":
		if yearStr != "" && weekStr != "" {
			// 특정 년도와 주차 선택
			year, yearErr := strconv.Atoi(yearStr)
			week, weekErr := strconv.Atoi(weekStr)

			if yearErr == nil && weekErr == nil && week >= 1 && week <= 53 {
				// 해당 년도 1월 1일부터 주차 계산
				jan1 := time.Date(year, 1, 1, 0, 0, 0, 0, now.Location())
				// 1월 1일이 무슨 요일인지 확인
				jan1Weekday := int(jan1.Weekday())
				if jan1Weekday == 0 {
					jan1Weekday = 7 // 일요일을 7로 변경
				}

				// 첫 번째 월요일을 찾기
				daysToFirstMonday := 8 - jan1Weekday
				if jan1Weekday == 1 {
					daysToFirstMonday = 0 // 1월 1일이 월요일인 경우
				}

				firstMonday := jan1.AddDate(0, 0, daysToFirstMonday)
				start = firstMonday.AddDate(0, 0, (week-1)*7)
				end = start.AddDate(0, 0, 6)

				period = fmt.Sprintf("%d년 %d주차 (%s ~ %s)",
					year, week,
					start.Format("1월 2일"),
					end.Format("1월 2일"))
			} else {
				// 잘못된 파라미터인 경우 현재 주로 대체
				h.setCurrentWeek(&start, &end, &period, now)
			}
		} else {
			// 파라미터가 없는 경우 현재 주
			h.setCurrentWeek(&start, &end, &period, now)
		}

	case "month":
		if yearStr != "" && monthStr != "" {
			// 특정 년도와 월 선택
			year, yearErr := strconv.Atoi(yearStr)
			month, monthErr := strconv.Atoi(monthStr)

			if yearErr == nil && monthErr == nil && month >= 1 && month <= 12 {
				start = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())
				end = start.AddDate(0, 1, -1)
				period = fmt.Sprintf("%d년 %d월", year, month)
			} else {
				// 잘못된 파라미터인 경우 현재 월로 대체
				h.setCurrentMonth(&start, &end, &period, now)
			}
		} else {
			// 파라미터가 없는 경우 현재 월
			h.setCurrentMonth(&start, &end, &period, now)
		}

	case "year":
		if yearStr != "" {
			// 특정 년도 선택
			year, yearErr := strconv.Atoi(yearStr)

			if yearErr == nil && year >= 2020 && year <= now.Year()+5 {
				start = time.Date(year, 1, 1, 0, 0, 0, 0, now.Location())
				end = time.Date(year, 12, 31, 23, 59, 59, 0, now.Location())
				period = fmt.Sprintf("%d년", year)
			} else {
				// 잘못된 파라미터인 경우 현재 년도로 대체
				h.setCurrentYear(&start, &end, &period, now)
			}
		} else {
			// 파라미터가 없는 경우 현재 년도
			h.setCurrentYear(&start, &end, &period, now)
		}

	case "custom":
		// 사용자 지정 기간
		if startDate != "" && endDate != "" {
			if s, err := time.Parse("2006-01-02", startDate); err == nil {
				start = s
			} else {
				start = now.AddDate(0, -1, 0)
			}
			if e, err := time.Parse("2006-01-02", endDate); err == nil {
				end = e
			} else {
				end = now
			}
			period = start.Format("2006년 1월 2일") + " ~ " + end.Format("1월 2일")
		} else {
			// 기본값: 지난 한 달
			start = now.AddDate(0, -1, 0)
			end = now
			period = "지난 한 달"
		}

	case "all":
		// 전체 기간
		start = time.Date(2020, 1, 1, 0, 0, 0, 0, now.Location())
		end = now
		period = "전체 기간"

	default:
		// 기본값: 현재 월
		h.setCurrentMonth(&start, &end, &period, now)
	}

	return start.Format("2006-01-02"), end.Format("2006-01-02"), period
}

// 현재 주 설정 헬퍼 함수
func (h *StatisticsHandler) setCurrentWeek(start, end *time.Time, period *string, now time.Time) {
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 일요일을 7로 변경
	}
	*start = now.AddDate(0, 0, -(weekday - 1))
	*end = start.AddDate(0, 0, 6)
	*period = start.Format("2006년 1월 2일") + " ~ " + end.Format("1월 2일")
}

// 현재 월 설정 헬퍼 함수
func (h *StatisticsHandler) setCurrentMonth(start, end *time.Time, period *string, now time.Time) {
	*start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	*end = start.AddDate(0, 1, -1)
	*period = now.Format("2006년 1월")
}

// 현재 년도 설정 헬퍼 함수
func (h *StatisticsHandler) setCurrentYear(start, end *time.Time, period *string, now time.Time) {
	*start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	*end = time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
	*period = now.Format("2006년")
}

// 색상 생성 헬퍼 함수
func (h *StatisticsHandler) generateColors(count int) []string {
	colors := make([]string, count)

	// 미리 정의된 색상 팔레트 (더 많은 색상)
	predefinedColors := []string{
		"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7", "#DDA0DD", "#98D8C8", "#F7DC6F",
		"#FF8A80", "#80CBC4", "#81C784", "#FFB74D", "#F06292", "#9575CD", "#64B5F6", "#4DB6AC",
		"#AED581", "#FFD54F", "#FF8A65", "#A1887F", "#90A4AE", "#FFAB91", "#CE93D8", "#80DEEA",
		"#C5E1A5", "#FFF176", "#BCAAA4", "#B39DDB", "#81D4FA", "#A5D6A7", "#FFCC02", "#FF7043",
	}

	for i := 0; i < count; i++ {
		if i < len(predefinedColors) {
			colors[i] = predefinedColors[i]
		} else {
			// HSL 색상 공간을 사용하여 고유한 색상 생성
			hue := (i * 137) % 360      // 황금각도를 사용하여 색상 분산
			saturation := 60 + (i%3)*15 // 60%, 75%, 90% 채도
			lightness := 50 + (i%4)*10  // 50%, 60%, 70%, 80% 명도
			colors[i] = fmt.Sprintf("hsl(%d, %d%%, %d%%)", hue, saturation, lightness)
		}
	}

	return colors
}

// 차트 데이터 생성 헬퍼 함수
func (h *StatisticsHandler) generateChartData(categories []models.CategoryStatistics) []models.ChartData {
	chartData := make([]models.ChartData, len(categories))
	colors := h.generateColors(len(categories))

	for i, category := range categories {
		chartData[i] = models.ChartData{
			Label:      category.CategoryName,
			Value:      category.TotalAmount,
			Percentage: category.Percentage,
			Color:      colors[i],
		}
	}

	return chartData
}

// 결제수단별 카테고리 통계 조회 핸들러
func (h *StatisticsHandler) GetPaymentMethodCategoryStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, models.ErrCodeInvalidInput, "지원되지 않는 메소드입니다.")
		return
	}

	// 쿼리 파라미터 파싱
	paymentMethodIDStr := r.URL.Query().Get("payment_method_id")
	statisticsType := r.URL.Query().Get("type")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// 새로운 개별 기간 선택 파라미터
	yearStr := r.URL.Query().Get("year")   // 선택한 년도
	monthStr := r.URL.Query().Get("month") // 선택한 월 (1-12)
	weekStr := r.URL.Query().Get("week")   // 선택한 주차 (1-53)

	if paymentMethodIDStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "결제수단 ID가 필요합니다.")
		return
	}

	paymentMethodID, err := strconv.Atoi(paymentMethodIDStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, models.ErrCodeInvalidInput, "올바른 결제수단 ID를 입력해주세요.")
		return
	}

	// 기본값 설정
	if statisticsType == "" {
		statisticsType = "month"
	}

	// 날짜 범위 계산
	calculatedStartDate, calculatedEndDate, period := h.calculateDateRange(statisticsType, startDate, endDate, yearStr, monthStr, weekStr)

	// 카테고리별 통계 조회
	categories, err := h.DB.GetPaymentMethodCategoryStatistics(paymentMethodID, calculatedStartDate, calculatedEndDate)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, models.ErrCodeDatabaseError, "결제수단별 카테고리 통계 조회 중 오류 발생")
		return
	}

	// 총합 계산
	var categoryTotal int
	for _, category := range categories {
		categoryTotal += category.TotalAmount
	}

	// 퍼센테지 계산
	for i := range categories {
		if categoryTotal > 0 {
			categories[i].Percentage = float64(categories[i].TotalAmount) / float64(categoryTotal) * 100
		}
	}

	// 차트 데이터 생성
	chartData := make([]models.ChartData, len(categories))
	colors := h.generateColors(len(categories))

	for i, category := range categories {
		chartData[i] = models.ChartData{
			Label:      category.CategoryName,
			Value:      category.TotalAmount,
			Percentage: category.Percentage,
			Color:      colors[i],
		}
	}

	response := map[string]interface{}{
		"period":            period,
		"payment_method_id": paymentMethodID,
		"payment_total":     categoryTotal,
		"categories":        categories,
		"chart_data":        chartData,
	}

	utils.SendSuccessResponse(w, response)
}
