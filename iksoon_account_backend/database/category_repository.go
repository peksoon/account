package database

import (
	"fmt"
	"time"

	"iksoon_account_backend/models"
)

// GetCategories 카테고리 목록 조회
func (db *DB) GetCategories(categoryType string) ([]models.Category, error) {
	var query string
	var args []interface{}

	if categoryType != "" {
		query = `
			SELECT id, name, type, created_at, updated_at 
			FROM categories 
			WHERE type = ? AND is_active = 1
			ORDER BY name ASC`
		args = append(args, categoryType)
	} else {
		query = `
			SELECT id, name, type, created_at, updated_at 
			FROM categories 
			WHERE is_active = 1
			ORDER BY type ASC, name ASC`
	}

	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("카테고리 조회 오류: %v", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		var createdAt, updatedAt string

		err := rows.Scan(&category.ID, &category.Name, &category.Type, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("카테고리 데이터 읽기 오류: %v", err)
		}

		// 시간 파싱
		if category.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
			category.CreatedAt = time.Now()
		}
		if category.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
			category.UpdatedAt = time.Now()
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// CreateCategory 카테고리 생성
func (db *DB) CreateCategory(name, categoryType string) (int64, error) {
	// 중복 확인 (같은 타입에서 같은 이름의 활성 카테고리)
	var count int
	err := db.Conn.QueryRow(`
		SELECT COUNT(*) FROM categories 
		WHERE name = ? AND type = ? AND is_active = 1`, name, categoryType).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("카테고리 중복 확인 오류: %v", err)
	}
	if count > 0 {
		return 0, fmt.Errorf("이미 존재하는 카테고리입니다")
	}

	query := `
		INSERT INTO categories (name, type, created_at, updated_at) 
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	result, err := db.Conn.Exec(query, name, categoryType)
	if err != nil {
		return 0, fmt.Errorf("카테고리 생성 오류: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("카테고리 ID 조회 오류: %v", err)
	}

	return id, nil
}

// UpdateCategory 카테고리 수정
func (db *DB) UpdateCategory(id int, name string, categoryType string) error {
	// 중복 확인 (자신 제외)
	var count int
	err := db.Conn.QueryRow(`
		SELECT COUNT(*) FROM categories 
		WHERE name = ? AND type = ? AND id != ? AND is_active = 1`,
		name, categoryType, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("카테고리 중복 확인 오류: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("이미 존재하는 카테고리입니다")
	}

	query := `
		UPDATE categories 
		SET name = ?, type = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Conn.Exec(query, name, categoryType, id)
	if err != nil {
		return fmt.Errorf("카테고리 수정 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("카테고리 수정 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// CheckCategoryUsage 카테고리 사용 여부 확인
func (db *DB) CheckCategoryUsage(categoryID int) (bool, error) {
	// 지출 데이터에서 사용 여부 확인
	outQuery := `SELECT COUNT(*) FROM out_account_data WHERE category_id = ?`
	var outCount int
	err := db.Conn.QueryRow(outQuery, categoryID).Scan(&outCount)
	if err != nil {
		return false, fmt.Errorf("지출 데이터에서 카테고리 사용 여부 확인 오류: %v", err)
	}

	// 수입 데이터에서 사용 여부 확인
	inQuery := `SELECT COUNT(*) FROM in_account_data WHERE category_id = ?`
	var inCount int
	err = db.Conn.QueryRow(inQuery, categoryID).Scan(&inCount)
	if err != nil {
		return false, fmt.Errorf("수입 데이터에서 카테고리 사용 여부 확인 오류: %v", err)
	}

	// 개발 시에만 로그 출력
	if outCount+inCount > 0 {
		fmt.Printf("카테고리 %d 사용 여부 확인: 지출 %d건, 수입 %d건, 총 %d건\n", categoryID, outCount, inCount, outCount+inCount)
	}

	return (outCount + inCount) > 0, nil
}

// DeleteCategory 카테고리 삭제 (비활성화로 변경하여 기존 가계부 정보 유지)
func (db *DB) DeleteCategory(id int) error {
	// 비활성화로 변경하여 기존 가계부 정보 유지
	query := `
		UPDATE categories 
		SET is_active = 0, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("카테고리 삭제 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("카테고리 삭제 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// ForceDeleteCategory 카테고리 강제 삭제 (비활성화로 변경하여 기존 가계부 정보 유지)
func (db *DB) ForceDeleteCategory(id int) error {
	// 사용 중이어도 비활성화만 하여 기존 가계부 정보 유지
	query := `
		UPDATE categories 
		SET is_active = 0, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("카테고리 강제 삭제 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("카테고리 삭제 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// GetCategoryByID ID로 카테고리 조회
func (db *DB) GetCategoryByID(id int) (*models.Category, error) {
	query := `
		SELECT id, name, type, created_at, updated_at 
		FROM categories 
		WHERE id = ?`

	var category models.Category
	var createdAt, updatedAt string

	err := db.Conn.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.Type, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("카테고리 조회 오류: %v", err)
	}

	// 시간 파싱
	if category.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
		category.CreatedAt = time.Now()
	}
	if category.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
		category.UpdatedAt = time.Now()
	}

	return &category, nil
}
