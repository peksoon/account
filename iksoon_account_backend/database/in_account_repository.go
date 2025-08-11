package database

import (
	"fmt"

	"iksoon_account_backend/models"
	"iksoon_account_backend/utils"

	"github.com/google/uuid"
)

// InsertInAccount 수입 데이터 삽입
func (db *DB) InsertInAccount(date, user string, money, categoryID int, keywordID *int, depositPathID int, memo string) error {
	uuidStr := uuid.New().String()
	parsedDate, err := utils.ParseDateTimeKST(date)
	if err != nil {
		utils.LogError("수입 데이터 날짜 파싱", err)
		return fmt.Errorf("날짜 파싱 오류: %v", err)
	}
	formattedDate := utils.FormatDateTimeKST(parsedDate)

	insertQuery := `
    INSERT INTO in_account_data (uuid, date, money, user, category_id, keyword_id, deposit_path_id, memo, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	// 디버그용 상세 로깅
	utils.Debug("수입 데이터 삽입 시도: UUID=%s, Date=%s, User=%s, Money=%d, CategoryID=%d, KeywordID=%v, DepositPathID=%d, Memo=%s",
		uuidStr, formattedDate, user, money, categoryID, keywordID, depositPathID, memo)

	_, err = db.Conn.Exec(insertQuery, uuidStr, formattedDate, money, user, categoryID, keywordID, depositPathID, memo)
	if err != nil {
		utils.LogError("수입 데이터 SQL 실행", err)
		utils.Debug("실패한 SQL: %s", insertQuery)
		utils.Debug("실패한 파라미터: [%s, %s, %d, %s, %d, %v, %d, %s]",
			uuidStr, formattedDate, money, user, categoryID, keywordID, depositPathID, memo)
		return fmt.Errorf("수입 데이터 삽입 오류: %v", err)
	}
	utils.Debug("수입 데이터 삽입 성공: UUID=%s", uuidStr)
	return nil
}

// GetInAccountsByDate 일별 수입 데이터 조회
func (db *DB) GetInAccountsByDate(date string) ([]models.InAccount, error) {
	query := `
    SELECT ia.uuid, ia.date, ia.user, ia.money, ia.category_id, ia.keyword_id, ia.deposit_path_id, ia.memo, ia.created_at, ia.updated_at,
           c.name as category_name,
           COALESCE(k.name, '') as keyword_name,
           COALESCE(dp.name, '') as deposit_path_name
    FROM in_account_data ia
    LEFT JOIN categories c ON ia.category_id = c.id
    LEFT JOIN keywords k ON ia.keyword_id = k.id
    LEFT JOIN deposit_paths dp ON ia.deposit_path_id = dp.id
    WHERE date(ia.date) = date(?)`

	rows, err := db.Conn.Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("수입 데이터 조회 오류: %v", err)
	}
	defer rows.Close()

	var inAccounts []models.InAccount
	for rows.Next() {
		var inAccount models.InAccount
		var keywordID *int

		err := rows.Scan(&inAccount.UUID, &inAccount.Date, &inAccount.User, &inAccount.Money,
			&inAccount.CategoryID, &keywordID, &inAccount.DepositPathID, &inAccount.Memo,
			&inAccount.CreatedAt, &inAccount.UpdatedAt,
			&inAccount.CategoryName, &inAccount.KeywordName, &inAccount.DepositPathName)
		if err != nil {
			return nil, fmt.Errorf("수입 데이터 읽기 오류: %v", err)
		}

		inAccount.KeywordID = keywordID
		inAccounts = append(inAccounts, inAccount)
	}

	return inAccounts, nil
}

// GetInAccountsForMonth 월별 수입 데이터 조회
func (db *DB) GetInAccountsForMonth(year, month string) ([]models.InAccount, error) {
	query := `
    SELECT ia.uuid, ia.date, ia.user, ia.money, ia.category_id, ia.keyword_id, ia.deposit_path_id, ia.memo, ia.created_at, ia.updated_at,
           c.name as category_name,
           COALESCE(k.name, '') as keyword_name,
           COALESCE(dp.name, '') as deposit_path_name
    FROM in_account_data ia
    LEFT JOIN categories c ON ia.category_id = c.id
    LEFT JOIN keywords k ON ia.keyword_id = k.id
    LEFT JOIN deposit_paths dp ON ia.deposit_path_id = dp.id
    WHERE substr(ia.date, 1, 7) = ?`

	rows, err := db.Conn.Query(query, year+"-"+month)
	if err != nil {
		return nil, fmt.Errorf("월별 수입 데이터 조회 오류: %v", err)
	}
	defer rows.Close()

	var accounts []models.InAccount
	for rows.Next() {
		var account models.InAccount
		var keywordID *int

		err := rows.Scan(&account.UUID, &account.Date, &account.User, &account.Money,
			&account.CategoryID, &keywordID, &account.DepositPathID, &account.Memo,
			&account.CreatedAt, &account.UpdatedAt,
			&account.CategoryName, &account.KeywordName, &account.DepositPathName)
		if err != nil {
			return nil, fmt.Errorf("수입 데이터 읽기 오류: %v", err)
		}

		account.KeywordID = keywordID
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// UpdateInAccount 수입 데이터 업데이트
func (db *DB) UpdateInAccount(uuidStr, date, user string, money, categoryID int, keywordID *int, depositPathID int, memo string) error {
	parsedDate, err := utils.ParseDateTimeKST(date)
	if err != nil {
		return fmt.Errorf("날짜 파싱 오류: %v", err)
	}
	formattedDate := utils.FormatDateTimeKST(parsedDate)

	updateQuery := `
    UPDATE in_account_data
    SET date = ?, money = ?, user = ?, category_id = ?, keyword_id = ?, deposit_path_id = ?, memo = ?, updated_at = CURRENT_TIMESTAMP
    WHERE uuid = ?`

	result, err := db.Conn.Exec(updateQuery, formattedDate, money, user, categoryID, keywordID, depositPathID, memo, uuidStr)
	if err != nil {
		return fmt.Errorf("수입 데이터 업데이트 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("수입 업데이트 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// DeleteInAccount 수입 데이터 삭제
func (db *DB) DeleteInAccount(uuidStr string) error {
	deleteQuery := `DELETE FROM in_account_data WHERE uuid = ?`
	result, err := db.Conn.Exec(deleteQuery, uuidStr)
	if err != nil {
		return fmt.Errorf("수입 데이터 삭제 오류: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("수입 삭제 결과 확인 오류: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

// GetInAccountByUUID UUID로 수입 데이터 조회
func (db *DB) GetInAccountByUUID(uuidStr string) (*models.InAccount, error) {
	query := `
    SELECT ia.uuid, ia.date, ia.user, ia.money, ia.category_id, ia.keyword_id, ia.deposit_path_id, ia.memo, ia.created_at, ia.updated_at,
           c.name as category_name,
           COALESCE(k.name, '') as keyword_name,
           COALESCE(dp.name, '') as deposit_path_name
    FROM in_account_data ia
    LEFT JOIN categories c ON ia.category_id = c.id
    LEFT JOIN keywords k ON ia.keyword_id = k.id
    LEFT JOIN deposit_paths dp ON ia.deposit_path_id = dp.id
    WHERE ia.uuid = ?`

	var inAccount models.InAccount
	var keywordID *int

	err := db.Conn.QueryRow(query, uuidStr).Scan(&inAccount.UUID, &inAccount.Date, &inAccount.User, &inAccount.Money,
		&inAccount.CategoryID, &keywordID, &inAccount.DepositPathID, &inAccount.Memo,
		&inAccount.CreatedAt, &inAccount.UpdatedAt,
		&inAccount.CategoryName, &inAccount.KeywordName, &inAccount.DepositPathName)
	if err != nil {
		return nil, fmt.Errorf("수입 데이터 조회 오류: %v", err)
	}

	inAccount.KeywordID = keywordID
	return &inAccount, nil
}
