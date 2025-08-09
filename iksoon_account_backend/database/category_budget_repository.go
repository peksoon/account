package database

import (
	"fmt"
	"time"

	"iksoon_account_backend/models"
)

// GetCategoryBudgets 카테고리 기준치 목록 조회
func (db *DB) GetCategoryBudgets(userName string, categoryID *int) ([]models.CategoryBudget, error) {
	var query string
	var args []interface{}

	if categoryID != nil && userName != "" {
		// 특정 사용자의 특정 카테고리 기준치 조회
		query = `
			SELECT cb.id, cb.category_id, c.name, cb.user_name, 
			       cb.monthly_budget, cb.yearly_budget,
			       cb.created_at, cb.updated_at
			FROM category_budgets cb
			LEFT JOIN categories c ON cb.category_id = c.id
			WHERE cb.user_name = ? AND cb.category_id = ?`
		args = append(args, userName, *categoryID)
	} else if categoryID != nil {
		// 특정 카테고리의 모든 기준치 조회 (사용자 구분 없음)
		query = `
			SELECT cb.id, cb.category_id, c.name, cb.user_name, 
			       cb.monthly_budget, cb.yearly_budget,
			       cb.created_at, cb.updated_at
			FROM category_budgets cb
			LEFT JOIN categories c ON cb.category_id = c.id
			WHERE cb.category_id = ?
			ORDER BY cb.user_name ASC`
		args = append(args, *categoryID)
	} else if userName != "" {
		// 특정 사용자의 모든 기준치 조회
		query = `
			SELECT cb.id, cb.category_id, c.name, cb.user_name, 
			       cb.monthly_budget, cb.yearly_budget,
			       cb.created_at, cb.updated_at
			FROM category_budgets cb
			LEFT JOIN categories c ON cb.category_id = c.id
			WHERE cb.user_name = ?
			ORDER BY c.name ASC`
		args = append(args, userName)
	} else {
		// 모든 기준치 조회
		query = `
			SELECT cb.id, cb.category_id, c.name, cb.user_name, 
			       cb.monthly_budget, cb.yearly_budget,
			       cb.created_at, cb.updated_at
			FROM category_budgets cb
			LEFT JOIN categories c ON cb.category_id = c.id
			ORDER BY cb.user_name ASC, c.name ASC`
	}

	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("카테고리 기준치 조회 오류: %v", err)
	}
	defer rows.Close()

	var budgets []models.CategoryBudget
	for rows.Next() {
		var budget models.CategoryBudget
		var createdAt, updatedAt string

		err := rows.Scan(&budget.ID, &budget.CategoryID, &budget.CategoryName,
			&budget.UserName, &budget.MonthlyBudget, &budget.YearlyBudget,
			&createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("카테고리 기준치 데이터 읽기 오류: %v", err)
		}

		// 시간 파싱
		if budget.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
			budget.CreatedAt = time.Now()
		}
		if budget.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
			budget.UpdatedAt = time.Now()
		}

		budgets = append(budgets, budget)
	}

	return budgets, nil
}

// CreateCategoryBudget 카테고리 기준치 생성
func (db *DB) CreateCategoryBudget(categoryID int, userName string, monthlyBudget, yearlyBudget int) (int64, error) {
	// 사용자명이 없는 경우 빈 문자열로 처리
	if userName == "" {
		userName = ""
	}

	// 중복 확인
	var count int
	err := db.Conn.QueryRow(`
		SELECT COUNT(*) FROM category_budgets 
		WHERE category_id = ? AND user_name = ?`,
		categoryID, userName).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("기준치 중복 확인 오류: %v", err)
	}
	if count > 0 {
		return 0, fmt.Errorf("이미 설정된 기준치가 있습니다")
	}

	query := `
		INSERT INTO category_budgets (category_id, user_name, monthly_budget, yearly_budget, created_at, updated_at) 
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	result, err := db.Conn.Exec(query, categoryID, userName, monthlyBudget, yearlyBudget)
	if err != nil {
		return 0, fmt.Errorf("기준치 생성 오류: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("기준치 ID 조회 오류: %v", err)
	}

	return id, nil
}

// UpdateCategoryBudget 카테고리 기준치 수정
func (db *DB) UpdateCategoryBudget(id int, monthlyBudget, yearlyBudget int) error {
	query := `
		UPDATE category_budgets 
		SET monthly_budget = ?, yearly_budget = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Conn.Exec(query, monthlyBudget, yearlyBudget, id)
	if err != nil {
		return fmt.Errorf("기준치 수정 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("기준치 수정 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("수정할 기준치를 찾을 수 없습니다")
	}

	return nil
}

// UpdateMonthlyBudget 월별 기준치만 수정
func (db *DB) UpdateMonthlyBudget(categoryID int, userName string, monthlyBudget int) error {
	// 사용자명이 없는 경우 빈 문자열로 처리
	if userName == "" {
		userName = ""
	}

	query := `
		UPDATE category_budgets
		SET monthly_budget = ?, updated_at = CURRENT_TIMESTAMP
		WHERE category_id = ? AND user_name = ?`
	args := []interface{}{monthlyBudget, categoryID, userName}

	result, err := db.Conn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("월별 기준치 수정 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("월별 기준치 수정 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("수정할 기준치를 찾을 수 없습니다")
	}

	return nil
}

// UpdateYearlyBudget 연별 기준치만 수정
func (db *DB) UpdateYearlyBudget(categoryID int, userName string, yearlyBudget int) error {
	// 사용자명이 없는 경우 빈 문자열로 처리
	if userName == "" {
		userName = ""
	}

	query := `
		UPDATE category_budgets
		SET yearly_budget = ?, updated_at = CURRENT_TIMESTAMP
		WHERE category_id = ? AND user_name = ?`
	args := []interface{}{yearlyBudget, categoryID, userName}

	result, err := db.Conn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("연별 기준치 수정 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("연별 기준치 수정 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("수정할 기준치를 찾을 수 없습니다")
	}

	return nil
}

// DeleteCategoryBudget 카테고리 기준치 삭제 (물리적 삭제)
func (db *DB) DeleteCategoryBudget(id int) error {
	query := `DELETE FROM category_budgets WHERE id = ?`

	result, err := db.Conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("기준치 삭제 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("기준치 삭제 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("삭제할 기준치를 찾을 수 없습니다")
	}

	return nil
}

// GetBudgetUsage 카테고리별 기준치 사용량 계산
func (db *DB) GetBudgetUsage(categoryID int, userName string, currentDate time.Time) (*models.BudgetUsage, error) {
	// 기준치 조회 (사용자별 기준치가 없으면 전체 기준치 조회)
	var budget models.CategoryBudget
	var categoryName string
	var createdAt, updatedAt string

	// 먼저 해당 사용자의 기준치 조회
	err := db.Conn.QueryRow(`
		SELECT cb.id, cb.category_id, c.name, cb.user_name, 
		       cb.monthly_budget, cb.yearly_budget,
		       cb.created_at, cb.updated_at
		FROM category_budgets cb
		LEFT JOIN categories c ON cb.category_id = c.id
		WHERE cb.category_id = ? AND cb.user_name = ?`,
		categoryID, userName).Scan(
		&budget.ID, &budget.CategoryID, &categoryName, &budget.UserName,
		&budget.MonthlyBudget, &budget.YearlyBudget,
		&createdAt, &updatedAt)

	if err != nil {
		// 사용자별 기준치가 없으면 전체 기준치 조회 (user_name = "")
		err = db.Conn.QueryRow(`
			SELECT cb.id, cb.category_id, c.name, cb.user_name, 
			       cb.monthly_budget, cb.yearly_budget,
			       cb.created_at, cb.updated_at
			FROM category_budgets cb
			LEFT JOIN categories c ON cb.category_id = c.id
			WHERE cb.category_id = ? AND cb.user_name = ''`,
			categoryID).Scan(
			&budget.ID, &budget.CategoryID, &categoryName, &budget.UserName,
			&budget.MonthlyBudget, &budget.YearlyBudget,
			&createdAt, &updatedAt)

		if err != nil {
			// 기준치가 전혀 설정되지 않은 경우
			return nil, nil
		}
	}

	budget.CategoryName = categoryName

	// 시간 파싱
	if budget.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
		budget.CreatedAt = time.Now()
	}
	if budget.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
		budget.UpdatedAt = time.Now()
	}

	// 월별 사용량 계산 (현재 월)
	monthStart := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())
	monthEnd := monthStart.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	var monthlyUsed int

	// 전체 기준치(user_name = "")인 경우 모든 사용자의 지출 합산
	if budget.UserName == "" {
		err = db.Conn.QueryRow(`
			SELECT COALESCE(SUM(money), 0) FROM out_account_data 
			WHERE category_id = ? 
			AND date >= ? AND date <= ?`,
			categoryID,
			monthStart.Format("2006-01-02 15:04:05"),
			monthEnd.Format("2006-01-02 15:04:05")).Scan(&monthlyUsed)
	} else {
		// 특정 사용자의 지출만 계산
		err = db.Conn.QueryRow(`
			SELECT COALESCE(SUM(money), 0) FROM out_account_data 
			WHERE category_id = ? AND user = ? 
			AND date >= ? AND date <= ?`,
			categoryID, userName,
			monthStart.Format("2006-01-02 15:04:05"),
			monthEnd.Format("2006-01-02 15:04:05")).Scan(&monthlyUsed)
	}

	if err != nil {
		return nil, fmt.Errorf("월별 사용량 계산 오류: %v", err)
	}

	// 연별 사용량 계산 (현재 연도)
	yearStart := time.Date(currentDate.Year(), 1, 1, 0, 0, 0, 0, currentDate.Location())
	yearEnd := time.Date(currentDate.Year(), 12, 31, 23, 59, 59, 0, currentDate.Location())

	var yearlyUsed int

	// 전체 기준치(user_name = "")인 경우 모든 사용자의 지출 합산
	if budget.UserName == "" {
		err = db.Conn.QueryRow(`
			SELECT COALESCE(SUM(money), 0) FROM out_account_data 
			WHERE category_id = ? 
			AND date >= ? AND date <= ?`,
			categoryID,
			yearStart.Format("2006-01-02 15:04:05"),
			yearEnd.Format("2006-01-02 15:04:05")).Scan(&yearlyUsed)
	} else {
		// 특정 사용자의 지출만 계산
		err = db.Conn.QueryRow(`
			SELECT COALESCE(SUM(money), 0) FROM out_account_data 
			WHERE category_id = ? AND user = ? 
			AND date >= ? AND date <= ?`,
			categoryID, userName,
			yearStart.Format("2006-01-02 15:04:05"),
			yearEnd.Format("2006-01-02 15:04:05")).Scan(&yearlyUsed)
	}

	if err != nil {
		return nil, fmt.Errorf("연별 사용량 계산 오류: %v", err)
	}

	// 사용량 정보 계산
	usage := &models.BudgetUsage{
		CategoryID:       categoryID,
		CategoryName:     categoryName,
		MonthlyBudget:    budget.MonthlyBudget,
		YearlyBudget:     budget.YearlyBudget,
		MonthlyUsed:      monthlyUsed,
		YearlyUsed:       yearlyUsed,
		MonthlyRemaining: budget.MonthlyBudget - monthlyUsed,
		YearlyRemaining:  budget.YearlyBudget - yearlyUsed,
	}

	// 퍼센티지 계산
	if budget.MonthlyBudget > 0 {
		usage.MonthlyPercent = float64(monthlyUsed) / float64(budget.MonthlyBudget) * 100
		usage.IsMonthlyOver = monthlyUsed > budget.MonthlyBudget
	}

	if budget.YearlyBudget > 0 {
		usage.YearlyPercent = float64(yearlyUsed) / float64(budget.YearlyBudget) * 100
		usage.IsYearlyOver = yearlyUsed > budget.YearlyBudget
	}

	return usage, nil
}

// GetAllBudgetUsages 사용자의 모든 카테고리 기준치 사용량 조회
func (db *DB) GetAllBudgetUsages(userName string, currentDate time.Time) ([]models.BudgetUsage, error) {
	// 사용자의 모든 기준치 조회
	budgets, err := db.GetCategoryBudgets(userName, nil)
	if err != nil {
		return nil, fmt.Errorf("기준치 목록 조회 오류: %v", err)
	}

	var usages []models.BudgetUsage
	for _, budget := range budgets {
		usage, err := db.GetBudgetUsage(budget.CategoryID, userName, currentDate)
		if err != nil {
			continue // 오류가 있는 항목은 건너뜀
		}
		if usage != nil {
			usages = append(usages, *usage)
		}
	}

	return usages, nil
}
