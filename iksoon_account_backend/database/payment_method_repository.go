package database

import (
	"fmt"
	"time"

	"iksoon_account_backend/models"
)

// GetPaymentMethods 결제수단 목록 조회 (계층구조)
func (db *DB) GetPaymentMethods() ([]models.PaymentMethod, error) {
	// 1단계: 부모 결제수단들 조회
	parentQuery := `
		SELECT id, name, parent_id, is_active, created_at, updated_at
		FROM payment_methods 
		WHERE parent_id IS NULL AND is_active = TRUE
		ORDER BY name ASC`

	parentRows, err := db.Conn.Query(parentQuery)
	if err != nil {
		return nil, fmt.Errorf("부모 결제수단 조회 오류: %v", err)
	}
	defer parentRows.Close()

	var methods []models.PaymentMethod
	for parentRows.Next() {
		var method models.PaymentMethod
		var createdAt, updatedAt string

		err := parentRows.Scan(&method.ID, &method.Name, &method.ParentID, &method.IsActive, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("부모 결제수단 데이터 읽기 오류: %v", err)
		}

		// 시간 파싱
		if method.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
			method.CreatedAt = time.Now()
		}
		if method.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
			method.UpdatedAt = time.Now()
		}

		// 2단계: 각 부모의 자식 결제수단들 조회
		childQuery := `
			SELECT id, name, parent_id, is_active, created_at, updated_at
			FROM payment_methods 
			WHERE parent_id = ? AND is_active = TRUE
			ORDER BY name ASC`

		childRows, err := db.Conn.Query(childQuery, method.ID)
		if err != nil {
			return nil, fmt.Errorf("자식 결제수단 조회 오류: %v", err)
		}

		var children []models.PaymentMethod
		for childRows.Next() {
			var child models.PaymentMethod
			var childCreatedAt, childUpdatedAt string

			err := childRows.Scan(&child.ID, &child.Name, &child.ParentID, &child.IsActive, &childCreatedAt, &childUpdatedAt)
			if err != nil {
				childRows.Close()
				return nil, fmt.Errorf("자식 결제수단 데이터 읽기 오류: %v", err)
			}

			// 시간 파싱
			if child.CreatedAt, err = time.Parse("2006-01-02 15:04:05", childCreatedAt); err != nil {
				child.CreatedAt = time.Now()
			}
			if child.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", childUpdatedAt); err != nil {
				child.UpdatedAt = time.Now()
			}

			children = append(children, child)
		}
		childRows.Close()

		method.Children = children
		methods = append(methods, method)
	}

	return methods, nil
}

// CreatePaymentMethod 결제수단 생성
func (db *DB) CreatePaymentMethod(name string, parentID *int) (int64, error) {
	// 중복 이름 확인 (같은 부모 하에서)
	checkQuery := `SELECT COUNT(*) FROM payment_methods WHERE name = ? AND parent_id = ? AND is_active = 1`
	var count int
	err := db.Conn.QueryRow(checkQuery, name, parentID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("결제수단 중복 확인 오류: %v", err)
	}

	if count > 0 {
		return 0, fmt.Errorf("이미 존재하는 결제수단 이름입니다")
	}

	query := `
		INSERT INTO payment_methods (name, parent_id, is_active, created_at, updated_at) 
		VALUES (?, ?, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	result, err := db.Conn.Exec(query, name, parentID)
	if err != nil {
		return 0, fmt.Errorf("결제수단 생성 오류: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("결제수단 ID 조회 오류: %v", err)
	}

	return id, nil
}

// UpdatePaymentMethod 결제수단 수정
func (db *DB) UpdatePaymentMethod(id int, name string) error {
	// 중복 이름 확인 (자기 자신 제외)
	checkQuery := `SELECT COUNT(*) FROM payment_methods WHERE name = ? AND id != ? AND is_active = 1`
	var count int
	err := db.Conn.QueryRow(checkQuery, name, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("결제수단 중복 확인 오류: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("이미 존재하는 결제수단 이름입니다")
	}

	query := `
		UPDATE payment_methods 
		SET name = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Conn.Exec(query, name, id)
	if err != nil {
		return fmt.Errorf("결제수단 수정 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("결제수단 수정 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// CheckPaymentMethodExists 결제수단 존재 여부 확인
func (db *DB) CheckPaymentMethodExists(id int) (bool, error) {
	query := `SELECT COUNT(*) FROM payment_methods WHERE id = ? AND is_active = 1`

	var count int
	err := db.Conn.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("결제수단 존재 여부 확인 오류: %v", err)
	}

	return count > 0, nil
}

// CheckPaymentMethodUsage 결제수단 사용 여부 확인
func (db *DB) CheckPaymentMethodUsage(paymentMethodID int) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM out_account_data 
		WHERE payment_method_id = ?`

	var count int
	err := db.Conn.QueryRow(query, paymentMethodID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("결제수단 사용 여부 확인 오류: %v", err)
	}

	return count > 0, nil
}

// DeletePaymentMethod 결제수단 논리 삭제
func (db *DB) DeletePaymentMethod(id int) error {
	query := `
		UPDATE payment_methods 
		SET is_active = 0, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("결제수단 삭제 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("결제수단 삭제 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// ForceDeletePaymentMethod 결제수단 강제 삭제 (사용 중인 경우)
func (db *DB) ForceDeletePaymentMethod(id int) error {
	query := `
		UPDATE payment_methods 
		SET is_active = 0, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("결제수단 강제 삭제 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("결제수단 강제 삭제 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// GetPaymentMethodByID ID로 결제수단 조회
func (db *DB) GetPaymentMethodByID(id int) (*models.PaymentMethod, error) {
	query := `
		SELECT id, name, is_active, created_at, updated_at
		FROM payment_methods 
		WHERE id = ? AND is_active = 1`

	var method models.PaymentMethod
	var createdAt, updatedAt string

	err := db.Conn.QueryRow(query, id).Scan(&method.ID, &method.Name,
		&method.IsActive, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("결제수단 조회 오류: %v", err)
	}

	// 시간 파싱
	if method.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
		method.CreatedAt = time.Now()
	}
	if method.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
		method.UpdatedAt = time.Now()
	}

	return &method, nil
}
