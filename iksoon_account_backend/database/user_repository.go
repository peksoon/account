package database

import (
	"database/sql"
	"fmt"
	"time"

	"iksoon_account_backend/models"
)

// GetUsers 사용자 목록 조회
func (db *DB) GetUsers() ([]models.User, error) {
	query := `
		SELECT id, name, COALESCE(email, ''), is_active, created_at, updated_at 
		FROM users 
		WHERE is_active = 1 
		ORDER BY name ASC`

	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("사용자 조회 오류: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var createdAt, updatedAt string

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.IsActive,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("사용자 스캔 오류: %v", err)
		}

		// 시간 파싱
		if user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
			user.CreatedAt = time.Now()
		}
		if user.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
			user.UpdatedAt = time.Now()
		}

		users = append(users, user)
	}

	return users, nil
}

// GetUserByID ID로 사용자 조회
func (db *DB) GetUserByID(id int) (*models.User, error) {
	query := `
		SELECT id, name, COALESCE(email, ''), is_active, created_at, updated_at 
		FROM users 
		WHERE id = ?`

	var user models.User
	var createdAt, updatedAt string

	err := db.Conn.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.IsActive,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("사용자를 찾을 수 없습니다")
		}
		return nil, fmt.Errorf("사용자 조회 오류: %v", err)
	}

	// 시간 파싱
	if user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
		user.CreatedAt = time.Now()
	}
	if user.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
		user.UpdatedAt = time.Now()
	}

	return &user, nil
}

// CreateUser 사용자 생성
func (db *DB) CreateUser(name, email string) (int64, error) {
	// 중복 이름 확인
	var count int
	err := db.Conn.QueryRow("SELECT COUNT(*) FROM users WHERE name = ? AND is_active = 1", name).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("사용자 중복 확인 오류: %v", err)
	}
	if count > 0 {
		return 0, fmt.Errorf("이미 존재하는 사용자 이름입니다")
	}

	query := `
		INSERT INTO users (name, email, updated_at) 
		VALUES (?, ?, CURRENT_TIMESTAMP)`

	result, err := db.Conn.Exec(query, name, email)
	if err != nil {
		return 0, fmt.Errorf("사용자 생성 오류: %v", err)
	}

	return result.LastInsertId()
}

// UpdateUser 사용자 수정
func (db *DB) UpdateUser(id int, name, email string) error {
	// 기존 사용자 이름 조회
	var oldName string
	err := db.Conn.QueryRow("SELECT name FROM users WHERE id = ? AND is_active = 1", id).Scan(&oldName)
	if err != nil {
		return fmt.Errorf("기존 사용자 정보 조회 오류: %v", err)
	}

	// 중복 이름 확인 (자신 제외)
	var count int
	err = db.Conn.QueryRow("SELECT COUNT(*) FROM users WHERE name = ? AND id != ? AND is_active = 1", name, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("사용자 중복 확인 오류: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("이미 존재하는 사용자 이름입니다")
	}

	// 트랜잭션 시작
	tx, err := db.Conn.Begin()
	if err != nil {
		return fmt.Errorf("트랜잭션 시작 오류: %v", err)
	}
	defer tx.Rollback()

	// 사용자 정보 업데이트
	query := `
		UPDATE users 
		SET name = ?, email = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := tx.Exec(query, name, email, id)
	if err != nil {
		return fmt.Errorf("사용자 수정 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("영향받은 행 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("사용자를 찾을 수 없습니다")
	}

	// 이름이 변경된 경우 가계부 정보 업데이트
	if oldName != name {
		// 지출 데이터 업데이트
		_, err = tx.Exec("UPDATE out_account_data SET user = ?, updated_at = CURRENT_TIMESTAMP WHERE user = ?", name, oldName)
		if err != nil {
			return fmt.Errorf("지출 데이터 사용자명 업데이트 오류: %v", err)
		}

		// 수입 데이터 업데이트
		_, err = tx.Exec("UPDATE in_account_data SET user = ?, updated_at = CURRENT_TIMESTAMP WHERE user = ?", name, oldName)
		if err != nil {
			return fmt.Errorf("수입 데이터 사용자명 업데이트 오류: %v", err)
		}
	}

	// 트랜잭션 커밋
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("트랜잭션 커밋 오류: %v", err)
	}

	return nil
}

// DeleteUser 사용자 삭제 (비활성화)
func (db *DB) DeleteUser(id int) error {
	// 사용 중인지 확인
	var count int
	err := db.Conn.QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT 1 FROM out_account_data WHERE user = (SELECT name FROM users WHERE id = ?)
			UNION ALL
			SELECT 1 FROM in_account_data WHERE user = (SELECT name FROM users WHERE id = ?)
		)`, id, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("사용자 사용 확인 오류: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("사용 중인 사용자는 삭제할 수 없습니다")
	}

	query := `
		UPDATE users 
		SET is_active = 0, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("사용자 삭제 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("영향받은 행 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("사용자를 찾을 수 없습니다")
	}

	return nil
}

// ForceDeleteUser 사용자 강제 삭제 (비활성화로 변경하여 기존 가계부 정보 유지)
func (db *DB) ForceDeleteUser(id int) error {
	// 사용 중이어도 비활성화만 하여 기존 가계부 정보 유지
	query := `
		UPDATE users 
		SET is_active = 0, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("사용자 강제 삭제 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("영향받은 행 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("사용자를 찾을 수 없습니다")
	}

	return nil
}

// CheckUserUsage 사용자 사용 여부 확인
func (db *DB) CheckUserUsage(userID int) (bool, error) {
	var count int
	err := db.Conn.QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT 1 FROM out_account_data WHERE user = (SELECT name FROM users WHERE id = ?)
			UNION ALL
			SELECT 1 FROM in_account_data WHERE user = (SELECT name FROM users WHERE id = ?)
		)`, userID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("사용자 사용 확인 오류: %v", err)
	}

	return count > 0, nil
}
