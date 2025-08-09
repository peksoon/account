package database

import (
	"fmt"

	"iksoon_account_backend/models"
	"iksoon_account_backend/utils"

	"github.com/google/uuid"
)

// InsertOutAccount 지출 데이터 삽입
func (db *DB) InsertOutAccount(date, user string, money, categoryID int, keywordID *int, paymentMethodID int, memo string) error {
	uuidStr := uuid.New().String()
	parsedDate, err := utils.ParseDateTimeKST(date)
	if err != nil {
		return fmt.Errorf("날짜 파싱 오류: %v", err)
	}
	formattedDate := utils.FormatDateTimeKST(parsedDate)

	insertQuery := `
    INSERT INTO out_account_data (uuid, date, money, user, category_id, keyword_id, payment_method_id, memo, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err = db.Conn.Exec(insertQuery, uuidStr, formattedDate, money, user, categoryID, keywordID, paymentMethodID, memo)
	if err != nil {
		return fmt.Errorf("지출 데이터 삽입 오류: %v", err)
	}
	return nil
}

// GetOutAccountsByDate 일별 지출 데이터 조회
func (db *DB) GetOutAccountsByDate(date string) ([]models.OutAccount, error) {
	query := `
    SELECT oa.uuid, oa.date, oa.user, oa.money, oa.category_id, oa.keyword_id, oa.payment_method_id, oa.memo, oa.created_at, oa.updated_at,
           c.name as category_name,
           COALESCE(k.name, '') as keyword_name,
           pm.name as payment_method_name
    FROM out_account_data oa
    LEFT JOIN categories c ON oa.category_id = c.id
    LEFT JOIN keywords k ON oa.keyword_id = k.id
    LEFT JOIN payment_methods pm ON oa.payment_method_id = pm.id
    WHERE date(oa.date) = date(?)`

	rows, err := db.Conn.Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("지출 데이터 조회 오류: %v", err)
	}
	defer rows.Close()

	var outAccounts []models.OutAccount
	for rows.Next() {
		var outAccount models.OutAccount
		var keywordID *int

		err := rows.Scan(&outAccount.UUID, &outAccount.Date, &outAccount.User, &outAccount.Money,
			&outAccount.CategoryID, &keywordID, &outAccount.PaymentMethodID, &outAccount.Memo,
			&outAccount.CreatedAt, &outAccount.UpdatedAt,
			&outAccount.CategoryName, &outAccount.KeywordName, &outAccount.PaymentMethodName)
		if err != nil {
			return nil, fmt.Errorf("지출 데이터 읽기 오류: %v", err)
		}

		outAccount.KeywordID = keywordID
		outAccounts = append(outAccounts, outAccount)
	}

	return outAccounts, nil
}

// GetOutAccountsForMonth 월별 지출 데이터 조회
func (db *DB) GetOutAccountsForMonth(year, month string) ([]models.OutAccount, error) {
	query := `
    SELECT oa.uuid, oa.date, oa.user, oa.money, oa.category_id, oa.keyword_id, oa.payment_method_id, oa.memo, oa.created_at, oa.updated_at,
           c.name as category_name,
           COALESCE(k.name, '') as keyword_name,
           pm.name as payment_method_name
    FROM out_account_data oa
    LEFT JOIN categories c ON oa.category_id = c.id
    LEFT JOIN keywords k ON oa.keyword_id = k.id
    LEFT JOIN payment_methods pm ON oa.payment_method_id = pm.id
    WHERE substr(oa.date, 1, 7) = ?`

	rows, err := db.Conn.Query(query, year+"-"+month)
	if err != nil {
		return nil, fmt.Errorf("월별 지출 데이터 조회 오류: %v", err)
	}
	defer rows.Close()

	var accounts []models.OutAccount
	for rows.Next() {
		var account models.OutAccount
		var keywordID *int

		err := rows.Scan(&account.UUID, &account.Date, &account.User, &account.Money,
			&account.CategoryID, &keywordID, &account.PaymentMethodID, &account.Memo,
			&account.CreatedAt, &account.UpdatedAt,
			&account.CategoryName, &account.KeywordName, &account.PaymentMethodName)
		if err != nil {
			return nil, fmt.Errorf("지출 데이터 읽기 오류: %v", err)
		}

		account.KeywordID = keywordID
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// UpdateOutAccount 지출 데이터 업데이트
func (db *DB) UpdateOutAccount(uuidStr, date, user string, money, categoryID int, keywordID *int, paymentMethodID int, memo string) error {
	parsedDate, err := utils.ParseDateTimeKST(date)
	if err != nil {
		return fmt.Errorf("날짜 파싱 오류: %v", err)
	}
	formattedDate := utils.FormatDateTimeKST(parsedDate)

	updateQuery := `
    UPDATE out_account_data
    SET date = ?, money = ?, user = ?, category_id = ?, keyword_id = ?, payment_method_id = ?, memo = ?, updated_at = CURRENT_TIMESTAMP
    WHERE uuid = ?`

	result, err := db.Conn.Exec(updateQuery, formattedDate, money, user, categoryID, keywordID, paymentMethodID, memo, uuidStr)
	if err != nil {
		return fmt.Errorf("지출 데이터 업데이트 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("지출 업데이트 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// DeleteOutAccount 지출 데이터 삭제
func (db *DB) DeleteOutAccount(uuidStr string) error {
	deleteQuery := `DELETE FROM out_account_data WHERE uuid = ?`
	result, err := db.Conn.Exec(deleteQuery, uuidStr)
	if err != nil {
		return fmt.Errorf("지출 데이터 삭제 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("지출 삭제 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// GetOutAccountByUUID UUID로 지출 데이터 조회
func (db *DB) GetOutAccountByUUID(uuidStr string) (*models.OutAccount, error) {
	query := `
    SELECT oa.uuid, oa.date, oa.user, oa.money, oa.category_id, oa.keyword_id, oa.payment_method_id, oa.memo, oa.created_at, oa.updated_at,
           c.name as category_name,
           COALESCE(k.name, '') as keyword_name,
           pm.name as payment_method_name
    FROM out_account_data oa
    LEFT JOIN categories c ON oa.category_id = c.id
    LEFT JOIN keywords k ON oa.keyword_id = k.id
    LEFT JOIN payment_methods pm ON oa.payment_method_id = pm.id
    WHERE oa.uuid = ?`

	var outAccount models.OutAccount
	var keywordID *int

	err := db.Conn.QueryRow(query, uuidStr).Scan(&outAccount.UUID, &outAccount.Date, &outAccount.User, &outAccount.Money,
		&outAccount.CategoryID, &keywordID, &outAccount.PaymentMethodID, &outAccount.Memo,
		&outAccount.CreatedAt, &outAccount.UpdatedAt,
		&outAccount.CategoryName, &outAccount.KeywordName, &outAccount.PaymentMethodName)
	if err != nil {
		return nil, fmt.Errorf("지출 데이터 조회 오류: %v", err)
	}

	outAccount.KeywordID = keywordID
	return &outAccount, nil
}
