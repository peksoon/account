package main

import (
	"log"
	"net/http"
	"os"

	"iksoon_account_backend/database"
	"iksoon_account_backend/handlers"
	"iksoon_account_backend/utils"
)

func main() {
	// 데이터베이스 초기화
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		dbPath = "./account_app.db"
	}

	db, err := database.InitDB(dbPath)
	if err != nil {
		utils.Error("데이터베이스 초기화 오류: %v", err)
		log.Fatalf("데이터베이스 초기화 실패")
	}
	defer db.Conn.Close()

	utils.Info("데이터베이스 연결 성공: %s", dbPath)

	// 핸들러 초기화
	userHandler := &handlers.UserHandler{DB: db}
	categoryHandler := &handlers.CategoryHandler{DB: db}
	keywordHandler := &handlers.KeywordHandler{DB: db}
	paymentMethodHandler := &handlers.PaymentMethodHandler{DB: db}
	depositPathHandler := &handlers.DepositPathHandler{DB: db}
	outAccountHandler := &handlers.OutAccountHandler{DB: db, KeywordDB: db}
	inAccountHandler := &handlers.InAccountHandler{DB: db, KeywordDB: db}
	statisticsHandler := &handlers.StatisticsHandler{DB: db}
	categoryBudgetHandler := handlers.NewCategoryBudgetHandler(db)

	// CORS 및 로깅 미들웨어
	enableCorsAndLogging := func(next http.Handler) http.Handler {
		return utils.LogHTTPMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		}))
	}

	// 라우트 설정

	// 사용자 관리 API
	http.Handle("/users", enableCorsAndLogging(http.HandlerFunc(userHandler.GetUsersHandler)))
	http.Handle("/users/create", enableCorsAndLogging(http.HandlerFunc(userHandler.CreateUserHandler)))
	http.Handle("/users/update", enableCorsAndLogging(http.HandlerFunc(userHandler.UpdateUserHandler)))
	http.Handle("/users/delete", enableCorsAndLogging(http.HandlerFunc(userHandler.DeleteUserHandler)))
	http.Handle("/users/force-delete", enableCorsAndLogging(http.HandlerFunc(userHandler.ForceDeleteUserHandler)))
	http.Handle("/users/check-usage", enableCorsAndLogging(http.HandlerFunc(userHandler.CheckUserUsageHandler)))

	// 카테고리 관리 API
	http.Handle("/categories", enableCorsAndLogging(http.HandlerFunc(categoryHandler.GetCategoriesHandler)))
	http.Handle("/categories/create", enableCorsAndLogging(http.HandlerFunc(categoryHandler.CreateCategoryHandler)))
	http.Handle("/categories/update", enableCorsAndLogging(http.HandlerFunc(categoryHandler.UpdateCategoryHandler)))
	http.Handle("/categories/delete", enableCorsAndLogging(http.HandlerFunc(categoryHandler.DeleteCategoryHandler)))
	http.Handle("/categories/force-delete", enableCorsAndLogging(http.HandlerFunc(categoryHandler.ForceDeleteCategoryHandler)))

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
	http.Handle("/v2/out-account/update", enableCorsAndLogging(http.HandlerFunc(outAccountHandler.UpdateOutAccountHandler)))
	http.Handle("/v2/out-account/delete", enableCorsAndLogging(http.HandlerFunc(outAccountHandler.DeleteOutAccountHandler)))

	// 수입 관리 API (새로운 구조)
	http.Handle("/v2/in-account/insert", enableCorsAndLogging(http.HandlerFunc(inAccountHandler.InsertInAccountHandler)))
	http.Handle("/v2/in-account", enableCorsAndLogging(http.HandlerFunc(inAccountHandler.GetInAccountByDateHandler)))
	http.Handle("/v2/month-in-account", enableCorsAndLogging(http.HandlerFunc(inAccountHandler.GetInAccountByMonthHandler)))
	http.Handle("/v2/in-account/update", enableCorsAndLogging(http.HandlerFunc(inAccountHandler.UpdateInAccountHandler)))
	http.Handle("/v2/in-account/delete", enableCorsAndLogging(http.HandlerFunc(inAccountHandler.DeleteInAccountHandler)))

	// 통계 API
	http.Handle("/statistics", enableCorsAndLogging(http.HandlerFunc(statisticsHandler.GetStatisticsHandler)))
	http.Handle("/statistics/category-keywords", enableCorsAndLogging(http.HandlerFunc(statisticsHandler.GetCategoryKeywordStatisticsHandler)))

	// 카테고리 기준치 관리 API
	http.Handle("/category-budgets", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.GetCategoryBudgetsHandler)))
	http.Handle("/category-budgets/create", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.CreateCategoryBudgetHandler)))
	http.Handle("/category-budgets/update", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.UpdateCategoryBudgetHandler)))
	http.Handle("/category-budgets/update-monthly", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.UpdateMonthlyBudgetHandler)))
	http.Handle("/category-budgets/update-yearly", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.UpdateYearlyBudgetHandler)))
	http.Handle("/category-budgets/delete", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.DeleteCategoryBudgetHandler)))
	http.Handle("/category-budgets/usage", enableCorsAndLogging(http.HandlerFunc(categoryBudgetHandler.GetBudgetUsageHandler)))

	// Health Check API
	http.Handle("/health", enableCorsAndLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"iksoon-account-backend"}`))
	})))

	// 서버 시작
	utils.LogStartup("8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
