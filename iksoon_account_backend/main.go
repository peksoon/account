package main

import (
	"log"
	"net/http"

	"iksoon_account_backend/config"
	"iksoon_account_backend/database"
	"iksoon_account_backend/handlers"
	"iksoon_account_backend/utils"
)

// main 애플리케이션 메인 함수 - 서버 초기화 및 실행
func main() {
	// 애플리케이션 설정 로드 및 초기화
	cfg := config.GetConfig()

	// 설정값 유효성 검사 수행
	if err := cfg.Validate(); err != nil {
		log.Fatalf("설정 오류: %v", err)
	}

	// 설정 파일에서 읽어온 로그 레벨로 전역 로거 설정
	switch cfg.LogLevel {
	case "DEBUG":
		utils.SetLogLevel(utils.DEBUG)
	case "INFO":
		utils.SetLogLevel(utils.INFO)
	case "WARNING", "WARN":
		utils.SetLogLevel(utils.WARNING)
	case "ERROR":
		utils.SetLogLevel(utils.ERROR)
	}

	// 디버그 모드일 때만 설정값 콘솔 출력
	if cfg.LogLevel == "DEBUG" {
		cfg.PrintConfig()
	}

	// SQLite 데이터베이스 연결 및 테이블 초기화
	dbPath := cfg.GetDBPath()
	db, err := database.InitDB(dbPath)
	if err != nil {
		utils.Error("데이터베이스 초기화 오류: %v", err)
		log.Fatalf("데이터베이스 초기화 실패")
	}
	defer db.Conn.Close()

	utils.Info("데이터베이스 연결 성공: %s", dbPath)

	// 각 도메인별 핸들러 인스턴스 생성 및 의존성 주입
	userHandler := &handlers.UserHandler{DB: db}
	categoryHandler := &handlers.CategoryHandler{DB: db}
	keywordHandler := &handlers.KeywordHandler{DB: db}
	paymentMethodHandler := &handlers.PaymentMethodHandler{DB: db}
	depositPathHandler := &handlers.DepositPathHandler{DB: db}
	outAccountHandler := &handlers.OutAccountHandler{DB: db, KeywordDB: db}
	inAccountHandler := &handlers.InAccountHandler{DB: db, KeywordDB: db}
	statisticsHandler := &handlers.StatisticsHandler{DB: db}
	categoryBudgetHandler := handlers.NewCategoryBudgetHandler(db)

	// CORS(Cross-Origin Resource Sharing) 및 HTTP 요청 로깅을 위한 미들웨어
	enableCorsAndLogging := func(next http.Handler) http.Handler {
		return utils.LogHTTPMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 프론트엔드에서의 요청을 허용하기 위한 CORS 헤더 설정
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// 브라우저 preflight 요청 처리
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// 다음 핸들러로 요청 전달
			next.ServeHTTP(w, r)
		}))
	}

	// RESTful API 라우트 설정 - 각 엔드포인트에 미들웨어 적용

	// 사용자 관리 API - 사용자 CRUD 및 사용 여부 확인
	http.Handle("/users", enableCorsAndLogging(http.HandlerFunc(userHandler.GetUsersHandler)))                     // GET: 사용자 목록 조회
	http.Handle("/users/create", enableCorsAndLogging(http.HandlerFunc(userHandler.CreateUserHandler)))            // POST: 신규 사용자 생성
	http.Handle("/users/update", enableCorsAndLogging(http.HandlerFunc(userHandler.UpdateUserHandler)))            // PUT: 사용자 정보 수정
	http.Handle("/users/delete", enableCorsAndLogging(http.HandlerFunc(userHandler.DeleteUserHandler)))            // DELETE: 사용자 삭제 (참조 데이터 있으면 실패)
	http.Handle("/users/force-delete", enableCorsAndLogging(http.HandlerFunc(userHandler.ForceDeleteUserHandler))) // DELETE: 사용자 강제 삭제 (참조 데이터 포함)
	http.Handle("/users/check-usage", enableCorsAndLogging(http.HandlerFunc(userHandler.CheckUserUsageHandler)))   // GET: 사용자 사용 여부 확인

	// 카테고리 관리 API - 지출/수입 카테고리 CRUD
	http.Handle("/categories", enableCorsAndLogging(http.HandlerFunc(categoryHandler.GetCategoriesHandler)))                    // GET: 카테고리 목록 조회 (type 파라미터로 out/in 필터링)
	http.Handle("/categories/create", enableCorsAndLogging(http.HandlerFunc(categoryHandler.CreateCategoryHandler)))            // POST: 신규 카테고리 생성
	http.Handle("/categories/update", enableCorsAndLogging(http.HandlerFunc(categoryHandler.UpdateCategoryHandler)))            // PUT: 카테고리 정보 수정
	http.Handle("/categories/delete", enableCorsAndLogging(http.HandlerFunc(categoryHandler.DeleteCategoryHandler)))            // DELETE: 카테고리 삭제 (사용 중이면 실패) - 기존 방식
	http.Handle("/categories/force-delete", enableCorsAndLogging(http.HandlerFunc(categoryHandler.ForceDeleteCategoryHandler))) // DELETE: 카테고리 강제 삭제 - 기존 방식
	// RESTful API 스타일 추가
	http.Handle("/categories/", enableCorsAndLogging(http.HandlerFunc(categoryHandler.CategoryRESTHandler))) // DELETE: /categories/{id} 또는 /categories/{id}/force-delete

	// 키워드 관리 API
	http.Handle("/keywords/suggestions", enableCorsAndLogging(http.HandlerFunc(keywordHandler.GetKeywordSuggestionsHandler)))
	http.Handle("/keywords/category", enableCorsAndLogging(http.HandlerFunc(keywordHandler.GetKeywordsByCategoryHandler)))
	http.Handle("/keywords/upsert", enableCorsAndLogging(http.HandlerFunc(keywordHandler.UpsertKeywordHandler)))
	http.Handle("/keywords/delete", enableCorsAndLogging(http.HandlerFunc(keywordHandler.DeleteKeywordHandler)))

	// 결제수단 관리 API
	http.Handle("/payment-methods", enableCorsAndLogging(http.HandlerFunc(paymentMethodHandler.GetPaymentMethodsHandler)))
	http.Handle("/payment-methods/create", enableCorsAndLogging(http.HandlerFunc(paymentMethodHandler.CreatePaymentMethodHandler)))
	http.Handle("/payment-methods/update", enableCorsAndLogging(http.HandlerFunc(paymentMethodHandler.UpdatePaymentMethodHandler)))
	http.Handle("/payment-methods/delete", enableCorsAndLogging(http.HandlerFunc(paymentMethodHandler.DeletePaymentMethodHandler)))
	http.Handle("/payment-methods/force-delete", enableCorsAndLogging(http.HandlerFunc(paymentMethodHandler.ForceDeletePaymentMethodHandler)))

	// 입금경로 관리 API
	http.Handle("/deposit-paths", enableCorsAndLogging(http.HandlerFunc(depositPathHandler.GetDepositPathsHandler)))
	http.Handle("/deposit-paths/create", enableCorsAndLogging(http.HandlerFunc(depositPathHandler.CreateDepositPathHandler)))
	http.Handle("/deposit-paths/update", enableCorsAndLogging(http.HandlerFunc(depositPathHandler.UpdateDepositPathHandler)))
	http.Handle("/deposit-paths/delete", enableCorsAndLogging(http.HandlerFunc(depositPathHandler.DeleteDepositPathHandler)))
	http.Handle("/deposit-paths/force-delete", enableCorsAndLogging(http.HandlerFunc(depositPathHandler.ForceDeleteDepositPathHandler)))

	// 지출 관리 API (새로운 구조)
	http.Handle("/v2/out-account/insert", enableCorsAndLogging(http.HandlerFunc(outAccountHandler.InsertOutAccountHandler)))
	http.Handle("/v2/out-account/insert-with-budget", enableCorsAndLogging(http.HandlerFunc(outAccountHandler.InsertOutAccountWithBudgetHandler)))
	http.Handle("/v2/out-account", enableCorsAndLogging(http.HandlerFunc(outAccountHandler.GetOutAccountByDateHandler)))
	http.Handle("/v2/month-out-account", enableCorsAndLogging(http.HandlerFunc(outAccountHandler.GetOutAccountByMonthHandler)))
	http.Handle("/v2/out-accounts", enableCorsAndLogging(http.HandlerFunc(outAccountHandler.GetOutAccountsByDateRangeHandler)))
	http.Handle("/v2/search-keyword-accounts", enableCorsAndLogging(http.HandlerFunc(outAccountHandler.SearchOutAccountsByKeywordHandler)))
	http.Handle("/v2/out-account/update", enableCorsAndLogging(http.HandlerFunc(outAccountHandler.UpdateOutAccountHandler)))
	http.Handle("/v2/out-account/delete", enableCorsAndLogging(http.HandlerFunc(outAccountHandler.DeleteOutAccountHandler)))

	// 수입 관리 API (새로운 구조)
	http.Handle("/v2/in-account/insert", enableCorsAndLogging(http.HandlerFunc(inAccountHandler.InsertInAccountHandler)))
	http.Handle("/v2/in-account", enableCorsAndLogging(http.HandlerFunc(inAccountHandler.GetInAccountByDateHandler)))
	http.Handle("/v2/month-in-account", enableCorsAndLogging(http.HandlerFunc(inAccountHandler.GetInAccountByMonthHandler)))
	http.Handle("/v2/in-accounts", enableCorsAndLogging(http.HandlerFunc(inAccountHandler.GetInAccountsByDateRangeHandler)))
	http.Handle("/v2/in-search-keyword-accounts", enableCorsAndLogging(http.HandlerFunc(inAccountHandler.SearchInAccountsByKeywordHandler)))
	http.Handle("/v2/in-account/update", enableCorsAndLogging(http.HandlerFunc(inAccountHandler.UpdateInAccountHandler)))
	http.Handle("/v2/in-account/delete", enableCorsAndLogging(http.HandlerFunc(inAccountHandler.DeleteInAccountHandler)))

	// 통계 API
	http.Handle("/statistics", enableCorsAndLogging(http.HandlerFunc(statisticsHandler.GetStatisticsHandler)))
	http.Handle("/statistics/category-keywords", enableCorsAndLogging(http.HandlerFunc(statisticsHandler.GetCategoryKeywordStatisticsHandler)))
	http.Handle("/statistics/payment-method-accounts", enableCorsAndLogging(http.HandlerFunc(outAccountHandler.GetOutAccountsByPaymentMethodHandler)))

	// 카테고리 기준치 관리 API
	http.Handle("/category-budgets", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.GetCategoryBudgetsHandler)))
	http.Handle("/category-budgets/create", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.CreateCategoryBudgetHandler)))
	http.Handle("/category-budgets/update", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.UpdateCategoryBudgetHandler)))
	http.Handle("/category-budgets/update-monthly", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.UpdateMonthlyBudgetHandler)))
	http.Handle("/category-budgets/update-yearly", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.UpdateYearlyBudgetHandler)))
	http.Handle("/category-budgets/delete", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.DeleteCategoryBudgetHandler)))
	http.Handle("/category-budgets/usage", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.GetBudgetUsageHandler)))

	// 서비스 상태 확인 API - 로드밸런서 및 모니터링 도구에서 사용
	http.Handle("/health", enableCorsAndLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"iksoon-account-backend"}`))
	})))

	// HTTP 서버 시작 - 설정된 포트에서 요청 대기
	utils.LogStartup(cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
