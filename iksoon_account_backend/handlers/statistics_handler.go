package handlers

import (
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
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// 기본값 설정
	if accountType == "" {
		accountType = "out"
	}
	if statisticsType == "" {
		statisticsType = "month"
	}

	// 날짜 범위 계산
	calculatedStartDate, calculatedEndDate, period := h.calculateDateRange(statisticsType, startDate, endDate)

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

	response := models.StatisticsResponse{
		Period:      period,
		TotalAmount: totalAmount,
		TotalCount:  totalCount,
		Categories:  categories,
		TopCategory: topCategory,
		ChartData:   chartData,
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
	calculatedStartDate, calculatedEndDate, period := h.calculateDateRange(statisticsType, startDate, endDate)

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
	colors := []string{"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7", "#DDA0DD", "#98D8C8", "#F7DC6F"}

	for i, keyword := range keywords {
		chartData[i] = models.ChartData{
			Label:      keyword.KeywordName,
			Value:      keyword.TotalAmount,
			Percentage: keyword.Percentage,
			Color:      colors[i%len(colors)],
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
func (h *StatisticsHandler) calculateDateRange(statisticsType, startDate, endDate string) (string, string, string) {
	now := time.Now()
	var start, end time.Time
	var period string

	switch statisticsType {
	case "week":
		// 이번 주 (월요일부터 일요일까지)
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // 일요일을 7로 변경
		}
		start = now.AddDate(0, 0, -(weekday - 1))
		end = start.AddDate(0, 0, 6)
		period = start.Format("2006년 1월 2일") + " ~ " + end.Format("1월 2일")

	case "month":
		// 이번 달
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, -1)
		period = now.Format("2006년 1월")

	case "year":
		// 올해
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
		period = now.Format("2006년")

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
		// 기본값: 이번 달
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, -1)
		period = now.Format("2006년 1월")
	}

	return start.Format("2006-01-02"), end.Format("2006-01-02"), period
}

// 차트 데이터 생성 헬퍼 함수
func (h *StatisticsHandler) generateChartData(categories []models.CategoryStatistics) []models.ChartData {
	chartData := make([]models.ChartData, len(categories))
	colors := []string{"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7", "#DDA0DD", "#98D8C8", "#F7DC6F"}

	for i, category := range categories {
		chartData[i] = models.ChartData{
			Label:      category.CategoryName,
			Value:      category.TotalAmount,
			Percentage: category.Percentage,
			Color:      colors[i%len(colors)],
		}
	}

	return chartData
}
