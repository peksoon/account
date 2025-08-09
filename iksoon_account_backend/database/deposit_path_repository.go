package database

import (
	"fmt"
	"time"

	"iksoon_account_backend/models"
)

// GetDepositPaths 입금경로 목록 조회
func (db *DB) GetDepositPaths() ([]models.DepositPath, error) {
	query := `
		SELECT id, name, is_active, created_at, updated_at
		FROM deposit_paths 
		WHERE is_active = TRUE
		ORDER BY name ASC`

	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("입금경로 조회 오류: %v", err)
	}
	defer rows.Close()

	var paths []models.DepositPath
	for rows.Next() {
		var path models.DepositPath
		var createdAt, updatedAt string

		err := rows.Scan(&path.ID, &path.Name, &path.IsActive, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("입금경로 데이터 읽기 오류: %v", err)
		}

		// 시간 파싱
		if path.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
			path.CreatedAt = time.Now()
		}
		if path.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
			path.UpdatedAt = time.Now()
		}

		paths = append(paths, path)
	}

	return paths, nil
}

// CreateDepositPath 입금경로 생성
func (db *DB) CreateDepositPath(name string) (int64, error) {
	// 중복 이름 확인
	checkQuery := `SELECT COUNT(*) FROM deposit_paths WHERE name = ? AND is_active = 1`
	var count int
	err := db.Conn.QueryRow(checkQuery, name).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("입금경로 중복 확인 오류: %v", err)
	}

	if count > 0 {
		return 0, fmt.Errorf("이미 존재하는 입금경로 이름입니다")
	}

	query := `
		INSERT INTO deposit_paths (name, is_active, created_at, updated_at) 
		VALUES (?, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	result, err := db.Conn.Exec(query, name)
	if err != nil {
		return 0, fmt.Errorf("입금경로 생성 오류: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("입금경로 ID 조회 오류: %v", err)
	}

	return id, nil
}

// UpdateDepositPath 입금경로 수정
func (db *DB) UpdateDepositPath(id int, name string) error {
	// 중복 이름 확인 (자기 자신 제외)
	checkQuery := `SELECT COUNT(*) FROM deposit_paths WHERE name = ? AND id != ? AND is_active = 1`
	var count int
	err := db.Conn.QueryRow(checkQuery, name, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("입금경로 중복 확인 오류: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("이미 존재하는 입금경로 이름입니다")
	}

	query := `
		UPDATE deposit_paths 
		SET name = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Conn.Exec(query, name, id)
	if err != nil {
		return fmt.Errorf("입금경로 수정 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("입금경로 수정 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// CheckDepositPathExists 입금경로 존재 여부 확인
func (db *DB) CheckDepositPathExists(id int) (bool, error) {
	query := `SELECT COUNT(*) FROM deposit_paths WHERE id = ? AND is_active = 1`

	var count int
	err := db.Conn.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("입금경로 존재 여부 확인 오류: %v", err)
	}

	return count > 0, nil
}

// CheckDepositPathUsage 입금경로 사용 여부 확인
func (db *DB) CheckDepositPathUsage(depositPathID int) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM in_account_data 
		WHERE deposit_path_id = ?`

	var count int
	err := db.Conn.QueryRow(query, depositPathID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("입금경로 사용 여부 확인 오류: %v", err)
	}

	return count > 0, nil
}

// DeleteDepositPath 입금경로 논리 삭제
func (db *DB) DeleteDepositPath(id int) error {
	query := `
		UPDATE deposit_paths 
		SET is_active = 0, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("입금경로 삭제 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("입금경로 삭제 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// ForceDeleteDepositPath 입금경로 강제 삭제 (사용 중인 경우)
func (db *DB) ForceDeleteDepositPath(id int) error {
	query := `
		UPDATE deposit_paths 
		SET is_active = 0, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("입금경로 강제 삭제 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("입금경로 강제 삭제 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// GetDepositPathByID ID로 입금경로 조회
func (db *DB) GetDepositPathByID(id int) (*models.DepositPath, error) {
	query := `
		SELECT id, name, is_active, created_at, updated_at
		FROM deposit_paths 
		WHERE id = ? AND is_active = 1`

	var path models.DepositPath
	var createdAt, updatedAt string

	err := db.Conn.QueryRow(query, id).Scan(&path.ID, &path.Name,
		&path.IsActive, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("입금경로 조회 오류: %v", err)
	}

	// 시간 파싱
	if path.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
		path.CreatedAt = time.Now()
	}
	if path.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
		path.UpdatedAt = time.Now()
	}

	return &path, nil
}
