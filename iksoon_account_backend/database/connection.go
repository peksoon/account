package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DB 구조체
type DB struct {
	Conn *sql.DB
}

// InitDB 데이터베이스 초기화
func InitDB(dbPath string) (*DB, error) {
	// 환경변수에서 데이터베이스 경로 가져오기
	if dbPath == "" {
		dbPath = os.Getenv("SQLITE_DB_PATH")
		if dbPath == "" {
			dbPath = "./account_app.db"
		}
	}

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// 외래키 제약조건 활성화
	_, err = conn.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("외래키 활성화 오류: %v", err)
	}

	// WAL 모드 활성화 (성능 향상)
	_, err = conn.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		return nil, fmt.Errorf("WAL 모드 활성화 오류: %v", err)
	}

	db := &DB{Conn: conn}

	// 테이블 생성 순서 중요 (외래키 제약조건 때문에)
	if err := db.createUserTable(); err != nil {
		return nil, err
	}

	if err := db.createCategoryTable(); err != nil {
		return nil, err
	}

	if err := db.createKeywordTable(); err != nil {
		return nil, err
	}

	if err := db.createPaymentMethodTable(); err != nil {
		return nil, err
	}

	if err := db.createDepositPathTable(); err != nil {
		return nil, err
	}

	if err := db.createOutAccountTable(); err != nil {
		return nil, err
	}

	if err := db.createInAccountTable(); err != nil {
		return nil, err
	}

	if err := db.createCategoryBudgetTable(); err != nil {
		return nil, err
	}

	return db, nil
}

// 테이블 생성 메서드들

// tableExists 테이블 존재 여부 확인 헬퍼 함수
func (db *DB) tableExists(tableName string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`
	err := db.Conn.QueryRow(query, tableName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 사용자 테이블 생성
func (db *DB) createUserTable() error {
	// 테이블이 이미 존재하는지 확인
	exists, err := db.tableExists("users")
	if err != nil {
		return fmt.Errorf("테이블 존재 여부 확인 오류: %v", err)
	}

	createUserTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(100) NOT NULL UNIQUE,
        email VARCHAR(255),
        is_active BOOLEAN DEFAULT 1,
        created_at TEXT DEFAULT CURRENT_TIMESTAMP,
        updated_at TEXT DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Conn.Exec(createUserTable)
	if err != nil {
		return fmt.Errorf("사용자 테이블 생성 오류: %v", err)
	}

	// 테이블이 새로 생성된 경우에만 기본 데이터 삽입
	if !exists {
		db.insertDefaultUsers()
	}
	return nil
}

func (db *DB) createCategoryTable() error {
	// 테이블이 이미 존재하는지 확인
	exists, err := db.tableExists("categories")
	if err != nil {
		return fmt.Errorf("테이블 존재 여부 확인 오류: %v", err)
	}

	createCategoryTable := `
    CREATE TABLE IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(255) NOT NULL,
        type VARCHAR(10) NOT NULL CHECK (type IN ('out', 'in')),
        is_active BOOLEAN DEFAULT 1,
        created_at TEXT DEFAULT CURRENT_TIMESTAMP,
        updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(name, type)
    );`

	_, err = db.Conn.Exec(createCategoryTable)
	if err != nil {
		return fmt.Errorf("카테고리 테이블 생성 오류: %v", err)
	}

	// 기존 테이블에 is_active 컬럼 추가 (마이그레이션)
	db.addCategoryIsActiveColumn()

	// 테이블이 새로 생성된 경우에만 기본 데이터 삽입
	if !exists {
		db.insertDefaultCategories()
	}
	return nil
}

func (db *DB) createKeywordTable() error {
	createKeywordTable := `
    CREATE TABLE IF NOT EXISTS keywords (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        category_id INTEGER NOT NULL,
        name VARCHAR(255) NOT NULL,
        usage_count INTEGER DEFAULT 1,
        is_active BOOLEAN DEFAULT 1,
        last_used TEXT DEFAULT CURRENT_TIMESTAMP,
        created_at TEXT DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
        UNIQUE(category_id, name)
    );`

	_, err := db.Conn.Exec(createKeywordTable)
	if err != nil {
		return fmt.Errorf("키워드 테이블 생성 오류: %v", err)
	}

	// 기존 테이블에 is_active 컬럼 추가 (마이그레이션)
	db.addKeywordIsActiveColumn()

	return nil
}

func (db *DB) createPaymentMethodTable() error {
	// 테이블이 이미 존재하는지 확인
	exists, err := db.tableExists("payment_methods")
	if err != nil {
		return fmt.Errorf("테이블 존재 여부 확인 오류: %v", err)
	}

	createPaymentMethodTable := `
    CREATE TABLE IF NOT EXISTS payment_methods (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(255) NOT NULL,
        parent_id INTEGER NULL,
        is_active BOOLEAN DEFAULT TRUE,
        created_at TEXT DEFAULT CURRENT_TIMESTAMP,
        updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (parent_id) REFERENCES payment_methods(id),
        UNIQUE(name, parent_id)
    );`

	_, err = db.Conn.Exec(createPaymentMethodTable)
	if err != nil {
		return fmt.Errorf("결제수단 테이블 생성 오류: %v", err)
	}

	// 테이블이 새로 생성된 경우에만 기본 데이터 삽입
	if !exists {
		db.insertDefaultPaymentMethods()
	}
	return nil
}

func (db *DB) createDepositPathTable() error {
	// 테이블이 이미 존재하는지 확인
	exists, err := db.tableExists("deposit_paths")
	if err != nil {
		return fmt.Errorf("테이블 존재 여부 확인 오류: %v", err)
	}

	createDepositPathTable := `
    CREATE TABLE IF NOT EXISTS deposit_paths (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(255) NOT NULL UNIQUE,
        is_active BOOLEAN DEFAULT TRUE,
        created_at TEXT DEFAULT CURRENT_TIMESTAMP,
        updated_at TEXT DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Conn.Exec(createDepositPathTable)
	if err != nil {
		return fmt.Errorf("입금경로 테이블 생성 오류: %v", err)
	}

	// 테이블이 새로 생성된 경우에만 기본 데이터 삽입
	if !exists {
		db.insertDefaultDepositPaths()
	}
	return nil
}

func (db *DB) createOutAccountTable() error {
	createOutAccountTable := `
    CREATE TABLE IF NOT EXISTS out_account_data (
        uuid TEXT PRIMARY KEY,
        date TEXT NOT NULL,
        money INT NOT NULL,
        user VARCHAR(255) NOT NULL,
        category_id INTEGER NOT NULL,
        keyword_id INTEGER NULL,
        payment_method_id INTEGER NOT NULL,
        memo TEXT,
        created_at TEXT DEFAULT CURRENT_TIMESTAMP,
        updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (category_id) REFERENCES categories(id),
        FOREIGN KEY (keyword_id) REFERENCES keywords(id),
        FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
    );`

	_, err := db.Conn.Exec(createOutAccountTable)
	if err != nil {
		return fmt.Errorf("지출 테이블 생성 오류: %v", err)
	}
	return nil
}

func (db *DB) createInAccountTable() error {
	createInAccountTable := `
    CREATE TABLE IF NOT EXISTS in_account_data (
        uuid TEXT PRIMARY KEY,
        date TEXT NOT NULL,
        money INT NOT NULL,
        user VARCHAR(255) NOT NULL,
        category_id INTEGER NOT NULL,
        keyword_id INTEGER NULL,
        deposit_path_id INTEGER NOT NULL,
        memo TEXT,
        created_at TEXT DEFAULT CURRENT_TIMESTAMP,
        updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (category_id) REFERENCES categories(id),
        FOREIGN KEY (keyword_id) REFERENCES keywords(id),
        FOREIGN KEY (deposit_path_id) REFERENCES deposit_paths(id)
    );`

	_, err := db.Conn.Exec(createInAccountTable)
	if err != nil {
		return fmt.Errorf("수입 테이블 생성 오류: %v", err)
	}
	return nil
}

// 기본 데이터 삽입 메서드들
func (db *DB) insertDefaultCategories() error {
	// 기존 카테고리 데이터 존재 여부 확인 (한번만 실행되도록)
	var count int
	err := db.Conn.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	if err != nil {
		return fmt.Errorf("카테고리 개수 확인 오류: %v", err)
	}

	// 이미 기본 데이터가 있으면 스킵
	if count > 0 {
		return nil
	}

	// 지출 카테고리
	outCategories := []string{
		"식비", "교통비", "생활용품", "의료비", "교육비",
		"문화생활", "쇼핑", "여행", "통신비", "주거비",
		"공과금", "보험료", "기타",
	}

	// 수입 카테고리
	inCategories := []string{
		"급여", "용돈", "상여금", "부업", "투자수익",
		"보험금", "환급", "기타수입",
	}

	// 지출 카테고리 삽입
	for _, category := range outCategories {
		_, err := db.Conn.Exec(`
			INSERT OR IGNORE INTO categories (name, type, is_active) 
			VALUES (?, 'out', 1)`, category)
		if err != nil {
			return fmt.Errorf("기본 지출 카테고리 삽입 오류: %v", err)
		}
	}

	// 수입 카테고리 삽입
	for _, category := range inCategories {
		_, err := db.Conn.Exec(`
			INSERT OR IGNORE INTO categories (name, type, is_active) 
			VALUES (?, 'in', 1)`, category)
		if err != nil {
			return fmt.Errorf("기본 수입 카테고리 삽입 오류: %v", err)
		}
	}

	return nil
}

// 기본 사용자 데이터 삽입
func (db *DB) insertDefaultUsers() error {
	// 기존 사용자 데이터 존재 여부 확인 (한번만 실행되도록)
	var count int
	err := db.Conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("사용자 개수 확인 오류: %v", err)
	}

	// 이미 기본 데이터가 있으면 스킵
	if count > 0 {
		return nil
	}

	defaultUsers := []string{"관리자", "손님"}

	for _, userName := range defaultUsers {
		_, err := db.Conn.Exec(`
			INSERT OR IGNORE INTO users (name) 
			VALUES (?)`, userName)
		if err != nil {
			return fmt.Errorf("기본 사용자 삽입 오류: %v", err)
		}
	}

	return nil
}

func (db *DB) insertDefaultPaymentMethods() error {
	// 기존 데이터 존재 여부 확인 (한번만 실행되도록)
	var count int
	err := db.Conn.QueryRow("SELECT COUNT(*) FROM payment_methods WHERE parent_id IS NULL").Scan(&count)
	if err != nil {
		return fmt.Errorf("결제수단 개수 확인 오류: %v", err)
	}

	// 이미 기본 데이터가 있으면 스킵
	if count > 0 {
		return nil
	}

	// 1단계: 기본 카테고리 (parent) 삽입
	defaultCategories := []string{"카드", "계좌이체", "현금", "기타"}

	for _, category := range defaultCategories {
		_, err := db.Conn.Exec(`
			INSERT OR IGNORE INTO payment_methods (name, parent_id, is_active) 
			VALUES (?, NULL, 1)`, category)
		if err != nil {
			return fmt.Errorf("기본 결제수단 카테고리 삽입 오류: %v", err)
		}
	}

	// 2단계: 기본 세부 결제수단 삽입
	// 카드 하위
	cardParentID, err := db.getPaymentMethodIDByName("카드")
	if err == nil {
		cardMethods := []string{"신용카드", "체크카드"}
		for _, method := range cardMethods {
			_, err := db.Conn.Exec(`
				INSERT OR IGNORE INTO payment_methods (name, parent_id, is_active) 
				VALUES (?, ?, 1)`, method, cardParentID)
			if err != nil {
				return fmt.Errorf("카드 세부 결제수단 삽입 오류: %v", err)
			}
		}
	}

	// 계좌이체 하위
	transferParentID, err := db.getPaymentMethodIDByName("계좌이체")
	if err == nil {
		transferMethods := []string{"온라인뱅킹", "ATM"}
		for _, method := range transferMethods {
			_, err := db.Conn.Exec(`
				INSERT OR IGNORE INTO payment_methods (name, parent_id, is_active) 
				VALUES (?, ?, 1)`, method, transferParentID)
			if err != nil {
				return fmt.Errorf("계좌이체 세부 결제수단 삽입 오류: %v", err)
			}
		}
	}

	return nil
}

// 결제수단 이름으로 ID 조회 헬퍼 함수
func (db *DB) getPaymentMethodIDByName(name string) (int, error) {
	var id int
	err := db.Conn.QueryRow("SELECT id FROM payment_methods WHERE name = ? AND parent_id IS NULL", name).Scan(&id)
	return id, err
}

func (db *DB) insertDefaultDepositPaths() error {
	// 기존 입금경로 데이터 존재 여부 확인 (한번만 실행되도록)
	var count int
	err := db.Conn.QueryRow("SELECT COUNT(*) FROM deposit_paths").Scan(&count)
	if err != nil {
		return fmt.Errorf("입금경로 개수 확인 오류: %v", err)
	}

	// 이미 기본 데이터가 있으면 스킵
	if count > 0 {
		return nil
	}

	// 기본 입금경로 삽입
	defaultPaths := []string{
		"급여계좌",
		"적금계좌",
		"현금",
		"기타",
	}

	for _, path := range defaultPaths {
		_, err := db.Conn.Exec(`
			INSERT OR IGNORE INTO deposit_paths (name, is_active) 
			VALUES (?, 1)`, path)
		if err != nil {
			return fmt.Errorf("기본 입금경로 삽입 오류: %v", err)
		}
	}

	return nil
}

// addCategoryIsActiveColumn 카테고리 테이블에 is_active 컬럼 추가 (마이그레이션)
func (db *DB) addCategoryIsActiveColumn() {
	// 컬럼이 이미 존재하는지 확인
	checkQuery := `PRAGMA table_info(categories)`
	rows, err := db.Conn.Query(checkQuery)
	if err != nil {
		return
	}
	defer rows.Close()

	hasIsActive := false
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull int
		var defaultValue interface{}
		var pk int

		err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
		if err != nil {
			continue
		}

		if name == "is_active" {
			hasIsActive = true
			break
		}
	}

	// 컬럼이 없으면 추가
	if !hasIsActive {
		alterQuery := `ALTER TABLE categories ADD COLUMN is_active BOOLEAN DEFAULT 1`
		db.Conn.Exec(alterQuery)
	}
}

// addKeywordIsActiveColumn 키워드 테이블에 is_active 컬럼 추가 (마이그레이션)
func (db *DB) addKeywordIsActiveColumn() {
	// 컬럼이 이미 존재하는지 확인
	checkQuery := `PRAGMA table_info(keywords)`
	rows, err := db.Conn.Query(checkQuery)
	if err != nil {
		return
	}
	defer rows.Close()

	hasIsActive := false
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull int
		var defaultValue interface{}
		var pk int

		err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
		if err != nil {
			continue
		}

		if name == "is_active" {
			hasIsActive = true
			break
		}
	}

	// 컬럼이 없으면 추가
	if !hasIsActive {
		alterQuery := `ALTER TABLE keywords ADD COLUMN is_active BOOLEAN DEFAULT 1`
		db.Conn.Exec(alterQuery)
	}
}

// 카테고리 기준치 테이블 생성
func (db *DB) createCategoryBudgetTable() error {
	createCategoryBudgetTable := `
    CREATE TABLE IF NOT EXISTS category_budgets (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        category_id INTEGER NOT NULL,
        user_name VARCHAR(255) DEFAULT '', 
        monthly_budget INTEGER DEFAULT 0,
        yearly_budget INTEGER DEFAULT 0,
        created_at TEXT DEFAULT CURRENT_TIMESTAMP,
        updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (category_id) REFERENCES categories(id),
        UNIQUE(category_id, user_name)
    );`

	_, err := db.Conn.Exec(createCategoryBudgetTable)
	if err != nil {
		return fmt.Errorf("카테고리 기준치 테이블 생성 오류: %v", err)
	}
	return nil
}
