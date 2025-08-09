package database

import (
	"fmt"
	"time"

	"iksoon_account_backend/models"
)

// GetKeywordSuggestions 키워드 자동완성 목록 조회
func (db *DB) GetKeywordSuggestions(categoryID int, query string, limit int) ([]models.KeywordSuggestion, error) {
	var sqlQuery string
	var args []interface{}

	if query != "" {
		sqlQuery = `
			SELECT id, name, usage_count 
			FROM keywords 
			WHERE category_id = ? AND name LIKE ? AND is_active = 1
			ORDER BY usage_count DESC, last_used DESC, name ASC
			LIMIT ?`
		args = []interface{}{categoryID, "%" + query + "%", limit}
	} else {
		sqlQuery = `
			SELECT id, name, usage_count 
			FROM keywords 
			WHERE category_id = ? AND is_active = 1
			ORDER BY usage_count DESC, last_used DESC, name ASC
			LIMIT ?`
		args = []interface{}{categoryID, limit}
	}

	rows, err := db.Conn.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("키워드 조회 오류: %v", err)
	}
	defer rows.Close()

	var suggestions []models.KeywordSuggestion
	for rows.Next() {
		var suggestion models.KeywordSuggestion
		err := rows.Scan(&suggestion.ID, &suggestion.Name, &suggestion.UsageCount)
		if err != nil {
			return nil, fmt.Errorf("키워드 데이터 읽기 오류: %v", err)
		}
		suggestions = append(suggestions, suggestion)
	}

	return suggestions, nil
}

// GetKeywordsByCategory 카테고리별 키워드 목록 조회
func (db *DB) GetKeywordsByCategory(categoryID int) ([]models.Keyword, error) {
	query := `
		SELECT id, category_id, name, usage_count, last_used, created_at
		FROM keywords 
		WHERE category_id = ? AND is_active = 1
		ORDER BY usage_count DESC, last_used DESC, name ASC`

	rows, err := db.Conn.Query(query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("키워드 조회 오류: %v", err)
	}
	defer rows.Close()

	var keywords []models.Keyword
	for rows.Next() {
		var keyword models.Keyword
		var lastUsed, createdAt string

		err := rows.Scan(&keyword.ID, &keyword.CategoryID, &keyword.Name,
			&keyword.UsageCount, &lastUsed, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("키워드 데이터 읽기 오류: %v", err)
		}

		// 시간 파싱
		if keyword.LastUsed, err = time.Parse("2006-01-02 15:04:05", lastUsed); err != nil {
			keyword.LastUsed = time.Now()
		}
		if keyword.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
			keyword.CreatedAt = time.Now()
		}

		keywords = append(keywords, keyword)
	}

	return keywords, nil
}

// UpsertKeyword 키워드 생성 또는 업데이트 (이미 존재하면 사용 횟수 증가)
func (db *DB) UpsertKeyword(categoryID int, name string) (int64, error) {
	// 기존 키워드 확인
	var existingID int64
	var usageCount int

	checkQuery := `SELECT id, usage_count FROM keywords WHERE category_id = ? AND name = ?`
	err := db.Conn.QueryRow(checkQuery, categoryID, name).Scan(&existingID, &usageCount)

	if err == nil {
		// 기존 키워드가 있으면 사용 횟수와 마지막 사용 시간 업데이트
		updateQuery := `
			UPDATE keywords 
			SET usage_count = usage_count + 1, last_used = CURRENT_TIMESTAMP 
			WHERE id = ?`

		_, err = db.Conn.Exec(updateQuery, existingID)
		if err != nil {
			return 0, fmt.Errorf("키워드 업데이트 오류: %v", err)
		}

		return existingID, nil
	}

	// 새로운 키워드 생성
	insertQuery := `
		INSERT INTO keywords (category_id, name, usage_count, last_used, created_at) 
		VALUES (?, ?, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	result, err := db.Conn.Exec(insertQuery, categoryID, name)
	if err != nil {
		return 0, fmt.Errorf("키워드 생성 오류: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("키워드 ID 조회 오류: %v", err)
	}

	return id, nil
}

// CheckKeywordUsage 키워드 사용 여부 확인
func (db *DB) CheckKeywordUsage(keywordID int) (bool, error) {
	// 지출 데이터에서 사용 여부 확인
	outQuery := `SELECT COUNT(*) FROM out_account_data WHERE keyword_id = ?`
	var outCount int
	err := db.Conn.QueryRow(outQuery, keywordID).Scan(&outCount)
	if err != nil {
		return false, fmt.Errorf("지출 데이터에서 키워드 사용 여부 확인 오류: %v", err)
	}

	// 수입 데이터에서 사용 여부 확인
	inQuery := `SELECT COUNT(*) FROM in_account_data WHERE keyword_id = ?`
	var inCount int
	err = db.Conn.QueryRow(inQuery, keywordID).Scan(&inCount)
	if err != nil {
		return false, fmt.Errorf("수입 데이터에서 키워드 사용 여부 확인 오류: %v", err)
	}

	return (outCount + inCount) > 0, nil
}

// DeleteKeyword 키워드 삭제 (비활성화로 변경하여 기존 가계부 정보 유지)
func (db *DB) DeleteKeyword(id int) error {
	// 비활성화로 변경하여 기존 가계부 정보 유지
	query := `
		UPDATE keywords 
		SET is_active = 0, last_used = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("키워드 삭제 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("키워드 삭제 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// GetKeywordByID ID로 키워드 조회
func (db *DB) GetKeywordByID(id int) (*models.Keyword, error) {
	query := `
		SELECT id, category_id, name, usage_count, last_used, created_at
		FROM keywords 
		WHERE id = ?`

	var keyword models.Keyword
	var lastUsed, createdAt string

	err := db.Conn.QueryRow(query, id).Scan(&keyword.ID, &keyword.CategoryID,
		&keyword.Name, &keyword.UsageCount, &lastUsed, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("키워드 조회 오류: %v", err)
	}

	// 시간 파싱
	if keyword.LastUsed, err = time.Parse("2006-01-02 15:04:05", lastUsed); err != nil {
		keyword.LastUsed = time.Now()
	}
	if keyword.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
		keyword.CreatedAt = time.Now()
	}

	return &keyword, nil
}

// GetKeywordByName 이름으로 키워드 조회
func (db *DB) GetKeywordByName(categoryID int, name string) (*models.Keyword, error) {
	query := `
		SELECT id, category_id, name, usage_count, last_used, created_at
		FROM keywords 
		WHERE category_id = ? AND name = ?`

	var keyword models.Keyword
	var lastUsed, createdAt string

	err := db.Conn.QueryRow(query, categoryID, name).Scan(&keyword.ID, &keyword.CategoryID,
		&keyword.Name, &keyword.UsageCount, &lastUsed, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("키워드 조회 오류: %v", err)
	}

	// 시간 파싱
	if keyword.LastUsed, err = time.Parse("2006-01-02 15:04:05", lastUsed); err != nil {
		keyword.LastUsed = time.Now()
	}
	if keyword.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
		keyword.CreatedAt = time.Now()
	}

	return &keyword, nil
}
