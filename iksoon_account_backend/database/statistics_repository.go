package database

import (
	"fmt"

	"iksoon_account_backend/models"
)

// GetCategoryStatistics 카테고리별 통계 조회
func (db *DB) GetCategoryStatistics(startDate, endDate, accountType string) ([]models.CategoryStatistics, error) {
	var query string

	if accountType == "out" {
		query = `
		SELECT 
			c.id as category_id,
			c.name as category_name,
			COALESCE(SUM(oa.money), 0) as total_amount,
			COALESCE(COUNT(oa.uuid), 0) as count
		FROM categories c
		LEFT JOIN out_account_data oa ON c.id = oa.category_id 
			AND date(oa.date) >= ? AND date(oa.date) <= ?
		WHERE c.type = 'out'
		GROUP BY c.id, c.name
		HAVING total_amount > 0
		ORDER BY total_amount DESC`
	} else {
		query = `
		SELECT 
			c.id as category_id,
			c.name as category_name,
			COALESCE(SUM(ia.money), 0) as total_amount,
			COALESCE(COUNT(ia.uuid), 0) as count
		FROM categories c
		LEFT JOIN in_account_data ia ON c.id = ia.category_id 
			AND date(ia.date) >= ? AND date(ia.date) <= ?
		WHERE c.type = 'in'
		GROUP BY c.id, c.name
		HAVING total_amount > 0
		ORDER BY total_amount DESC`
	}

	rows, err := db.Conn.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("카테고리 통계 조회 오류: %v", err)
	}
	defer rows.Close()

	var statistics []models.CategoryStatistics
	for rows.Next() {
		var stat models.CategoryStatistics
		err := rows.Scan(&stat.CategoryID, &stat.CategoryName, &stat.TotalAmount, &stat.Count)
		if err != nil {
			return nil, fmt.Errorf("카테고리 통계 데이터 읽기 오류: %v", err)
		}
		statistics = append(statistics, stat)
	}

	return statistics, nil
}

// GetKeywordStatistics 키워드별 통계 조회
func (db *DB) GetKeywordStatistics(categoryID int, startDate, endDate, accountType string) ([]models.KeywordStatistics, error) {
	var query string

	if accountType == "out" {
		query = `
		SELECT 
			k.id as keyword_id,
			k.name as keyword_name,
			COALESCE(SUM(oa.money), 0) as total_amount,
			COALESCE(COUNT(oa.uuid), 0) as count
		FROM keywords k
		LEFT JOIN out_account_data oa ON k.id = oa.keyword_id 
			AND oa.category_id = ? 
			AND date(oa.date) >= ? AND date(oa.date) <= ?
		WHERE k.category_id = ?
		GROUP BY k.id, k.name
		HAVING total_amount > 0
		ORDER BY total_amount DESC`
	} else {
		query = `
		SELECT 
			k.id as keyword_id,
			k.name as keyword_name,
			COALESCE(SUM(ia.money), 0) as total_amount,
			COALESCE(COUNT(ia.uuid), 0) as count
		FROM keywords k
		LEFT JOIN in_account_data ia ON k.id = ia.keyword_id 
			AND ia.category_id = ? 
			AND date(ia.date) >= ? AND date(ia.date) <= ?
		WHERE k.category_id = ?
		GROUP BY k.id, k.name
		HAVING total_amount > 0
		ORDER BY total_amount DESC`
	}

	rows, err := db.Conn.Query(query, categoryID, startDate, endDate, categoryID)
	if err != nil {
		return nil, fmt.Errorf("키워드 통계 조회 오류: %v", err)
	}
	defer rows.Close()

	var statistics []models.KeywordStatistics
	for rows.Next() {
		var stat models.KeywordStatistics
		err := rows.Scan(&stat.KeywordID, &stat.KeywordName, &stat.TotalAmount, &stat.Count)
		if err != nil {
			return nil, fmt.Errorf("키워드 통계 데이터 읽기 오류: %v", err)
		}
		statistics = append(statistics, stat)
	}

	return statistics, nil
}

// GetTotalAmount 총 금액과 개수 조회
func (db *DB) GetTotalAmount(startDate, endDate, accountType string) (int, int, error) {
	var query string

	if accountType == "out" {
		query = `
		SELECT 
			COALESCE(SUM(money), 0) as total_amount,
			COALESCE(COUNT(uuid), 0) as total_count
		FROM out_account_data 
		WHERE date(date) >= ? AND date(date) <= ?`
	} else {
		query = `
		SELECT 
			COALESCE(SUM(money), 0) as total_amount,
			COALESCE(COUNT(uuid), 0) as total_count
		FROM in_account_data 
		WHERE date(date) >= ? AND date(date) <= ?`
	}

	var totalAmount, totalCount int
	err := db.Conn.QueryRow(query, startDate, endDate).Scan(&totalAmount, &totalCount)
	if err != nil {
		return 0, 0, fmt.Errorf("총 금액 조회 오류: %v", err)
	}

	return totalAmount, totalCount, nil
}

// GetMonthlyTrend 월별 트렌드 조회 (최근 12개월)
func (db *DB) GetMonthlyTrend(accountType string) ([]map[string]interface{}, error) {
	var query string

	if accountType == "out" {
		query = `
		SELECT 
			substr(date, 1, 7) as month,
			SUM(money) as total_amount,
			COUNT(uuid) as total_count
		FROM out_account_data 
		WHERE date >= date('now', '-12 months')
		GROUP BY substr(date, 1, 7)
		ORDER BY month ASC`
	} else {
		query = `
		SELECT 
			substr(date, 1, 7) as month,
			SUM(money) as total_amount,
			COUNT(uuid) as total_count
		FROM in_account_data 
		WHERE date >= date('now', '-12 months')
		GROUP BY substr(date, 1, 7)
		ORDER BY month ASC`
	}

	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("월별 트렌드 조회 오류: %v", err)
	}
	defer rows.Close()

	var trends []map[string]interface{}
	for rows.Next() {
		var month string
		var totalAmount, totalCount int

		err := rows.Scan(&month, &totalAmount, &totalCount)
		if err != nil {
			return nil, fmt.Errorf("월별 트렌드 데이터 읽기 오류: %v", err)
		}

		trend := map[string]interface{}{
			"month":        month,
			"total_amount": totalAmount,
			"total_count":  totalCount,
		}
		trends = append(trends, trend)
	}

	return trends, nil
}

// GetDailyTrend 일별 트렌드 조회 (최근 30일)
func (db *DB) GetDailyTrend(accountType string) ([]map[string]interface{}, error) {
	var query string

	if accountType == "out" {
		query = `
		SELECT 
			date(date) as day,
			SUM(money) as total_amount,
			COUNT(uuid) as total_count
		FROM out_account_data 
		WHERE date >= date('now', '-30 days')
		GROUP BY date(date)
		ORDER BY day ASC`
	} else {
		query = `
		SELECT 
			date(date) as day,
			SUM(money) as total_amount,
			COUNT(uuid) as total_count
		FROM in_account_data 
		WHERE date >= date('now', '-30 days')
		GROUP BY date(date)
		ORDER BY day ASC`
	}

	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("일별 트렌드 조회 오류: %v", err)
	}
	defer rows.Close()

	var trends []map[string]interface{}
	for rows.Next() {
		var day string
		var totalAmount, totalCount int

		err := rows.Scan(&day, &totalAmount, &totalCount)
		if err != nil {
			return nil, fmt.Errorf("일별 트렌드 데이터 읽기 오류: %v", err)
		}

		trend := map[string]interface{}{
			"day":          day,
			"total_amount": totalAmount,
			"total_count":  totalCount,
		}
		trends = append(trends, trend)
	}

	return trends, nil
}

// GetTopCategories 상위 카테고리 조회 (특정 개수)
func (db *DB) GetTopCategories(startDate, endDate, accountType string, limit int) ([]models.CategoryStatistics, error) {
	var query string

	if accountType == "out" {
		query = `
		SELECT 
			c.id as category_id,
			c.name as category_name,
			COALESCE(SUM(oa.money), 0) as total_amount,
			COALESCE(COUNT(oa.uuid), 0) as count
		FROM categories c
		LEFT JOIN out_account_data oa ON c.id = oa.category_id 
			AND date(oa.date) >= ? AND date(oa.date) <= ?
		WHERE c.type = 'out'
		GROUP BY c.id, c.name
		HAVING total_amount > 0
		ORDER BY total_amount DESC
		LIMIT ?`
	} else {
		query = `
		SELECT 
			c.id as category_id,
			c.name as category_name,
			COALESCE(SUM(ia.money), 0) as total_amount,
			COALESCE(COUNT(ia.uuid), 0) as count
		FROM categories c
		LEFT JOIN in_account_data ia ON c.id = ia.category_id 
			AND date(ia.date) >= ? AND date(ia.date) <= ?
		WHERE c.type = 'in'
		GROUP BY c.id, c.name
		HAVING total_amount > 0
		ORDER BY total_amount DESC
		LIMIT ?`
	}

	rows, err := db.Conn.Query(query, startDate, endDate, limit)
	if err != nil {
		return nil, fmt.Errorf("상위 카테고리 조회 오류: %v", err)
	}
	defer rows.Close()

	var statistics []models.CategoryStatistics
	for rows.Next() {
		var stat models.CategoryStatistics
		err := rows.Scan(&stat.CategoryID, &stat.CategoryName, &stat.TotalAmount, &stat.Count)
		if err != nil {
			return nil, fmt.Errorf("상위 카테고리 데이터 읽기 오류: %v", err)
		}
		statistics = append(statistics, stat)
	}

	return statistics, nil
}

// GetPaymentMethodStatistics 결제수단별 통계 조회 (지출만)
func (db *DB) GetPaymentMethodStatistics(startDate, endDate string) ([]models.PaymentMethodStatistics, error) {
	query := `
	SELECT 
		pm.id as payment_method_id,
		pm.name as payment_method_name,
		COALESCE(SUM(oa.money), 0) as total_amount,
		COALESCE(COUNT(oa.uuid), 0) as count
	FROM payment_methods pm
	LEFT JOIN out_account_data oa ON pm.id = oa.payment_method_id 
		AND date(oa.date) >= ? AND date(oa.date) <= ?
	WHERE pm.is_active = 1
	GROUP BY pm.id, pm.name
	HAVING total_amount > 0
	ORDER BY total_amount DESC`

	rows, err := db.Conn.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("결제수단 통계 조회 오류: %v", err)
	}
	defer rows.Close()

	var statistics []models.PaymentMethodStatistics
	for rows.Next() {
		var stat models.PaymentMethodStatistics
		err := rows.Scan(&stat.PaymentMethodID, &stat.PaymentMethodName, &stat.TotalAmount, &stat.Count)
		if err != nil {
			return nil, fmt.Errorf("결제수단 통계 데이터 읽기 오류: %v", err)
		}
		statistics = append(statistics, stat)
	}

	return statistics, nil
}

// GetPaymentMethodCategoryStatistics 결제수단별 카테고리 통계 조회
func (db *DB) GetPaymentMethodCategoryStatistics(paymentMethodID int, startDate, endDate string) ([]models.CategoryStatistics, error) {
	query := `
	SELECT 
		c.id as category_id,
		c.name as category_name,
		COALESCE(SUM(oa.money), 0) as total_amount,
		COALESCE(COUNT(oa.uuid), 0) as count
	FROM categories c
	LEFT JOIN out_account_data oa ON c.id = oa.category_id 
		AND oa.payment_method_id = ?
		AND date(oa.date) >= ? AND date(oa.date) <= ?
	WHERE c.type = 'out'
	GROUP BY c.id, c.name
	HAVING total_amount > 0
	ORDER BY total_amount DESC`

	rows, err := db.Conn.Query(query, paymentMethodID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("결제수단별 카테고리 통계 조회 오류: %v", err)
	}
	defer rows.Close()

	var statistics []models.CategoryStatistics
	for rows.Next() {
		var stat models.CategoryStatistics
		err := rows.Scan(&stat.CategoryID, &stat.CategoryName, &stat.TotalAmount, &stat.Count)
		if err != nil {
			return nil, fmt.Errorf("결제수단별 카테고리 통계 데이터 읽기 오류: %v", err)
		}
		statistics = append(statistics, stat)
	}

	return statistics, nil
}

// GetUserStatistics 사용자별 통계 조회 (지출만)
func (db *DB) GetUserStatistics(startDate, endDate string) ([]models.UserStatistics, error) {
	query := `
	SELECT 
		oa.user,
		COALESCE(SUM(oa.money), 0) as total_amount,
		COALESCE(COUNT(oa.uuid), 0) as count
	FROM out_account_data oa
	WHERE date(oa.date) >= ? AND date(oa.date) <= ?
	GROUP BY oa.user
	HAVING total_amount > 0
	ORDER BY total_amount DESC`

	rows, err := db.Conn.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("사용자별 통계 조회 오류: %v", err)
	}
	defer rows.Close()

	var statistics []models.UserStatistics
	for rows.Next() {
		var stat models.UserStatistics
		err := rows.Scan(&stat.UserName, &stat.TotalAmount, &stat.Count)
		if err != nil {
			return nil, fmt.Errorf("사용자별 통계 데이터 읽기 오류: %v", err)
		}
		statistics = append(statistics, stat)
	}

	return statistics, nil
}
