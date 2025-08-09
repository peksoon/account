# 💰 스마트 가계부 백엔드 API

Go 언어로 작성된 가계부 관리 시스템의 백엔드 API 서버입니다.

## ✨ 주요 기능

- 📊 **수입/지출 관리**: 카테고리별 거래 데이터 관리
- 📅 **월별 통계**: 실시간 수입/지출 통계 제공
- 🏷 **카테고리 관리**: 커스텀 카테고리 시스템
- 💳 **결제수단 관리**: 다양한 결제수단 지원
- 🏦 **은행계좌 관리**: 계좌별 거래 관리
- 🔍 **키워드 검색**: 스마트 키워드 자동완성
- 📈 **통계 분석**: 상세한 거래 패턴 분석

## 🛠 기술 스택

### Backend

- **Go 1.21+** - 메인 언어
- **SQLite** - 데이터베이스
- **Standard Library** - HTTP 서버 (net/http)
- **CORS** - Cross-Origin Resource Sharing 지원

### 개발 도구

- **Go Modules** - 의존성 관리
- **구조화된 로깅** - 레벨별 로그 시스템
- **에러 코드 시스템** - 표준화된 에러 처리

## 🚀 설치 및 실행

### 사전 요구사항

- Go 1.21 이상
- Git

### 프로젝트 설정

```bash
# 프로젝트 클론
git clone <repository-url>
cd iksoon_account_backend

# 의존성 설치
go mod tidy
```

### 개발 서버 실행

```bash
# 개발 환경으로 실행
set ENV=development  # Windows
export ENV=development  # Linux/Mac

go run main.go
```

### 프로덕션 빌드

```bash
# 빌드
go build -o account_server .

# 실행
set ENV=production  # Windows
./account_server  # Linux/Mac
```

## 📝 환경 변수

```bash
# 환경 설정
ENV=development          # 실행 환경 (development/production)
LOG_LEVEL=DEBUG         # 로그 레벨 (DEBUG/INFO/WARNING/ERROR)
SQLITE_DB_PATH=./account_app.db  # 데이터베이스 파일 경로
```

## 📊 로그 시스템

### 로그 레벨

- **DEBUG**: 개발 시 상세 정보 (개발환경에서만)
- **INFO**: 일반 정보 (요청/응답 로그)
- **WARNING**: 경고 메시지
- **ERROR**: 에러 메시지 (운영환경에서 필수)

### 로그 설정

```bash
# 개발 환경 - 모든 로그 출력
set ENV=development
set LOG_LEVEL=DEBUG

# 운영 환경 - 에러 로그만 출력
set ENV=production
set LOG_LEVEL=ERROR
```

## 🚨 에러 처리 시스템

### 에러 코드 구조

```json
{
  "error": {
    "code": "CATEGORY_NOT_FOUND",
    "message": "카테고리를 찾을 수 없습니다",
    "status": 404
  }
}
```

### 주요 에러 코드

- `INTERNAL_SERVER_ERROR`: 내부 서버 오류
- `INVALID_REQUEST`: 잘못된 요청
- `NOT_FOUND`: 데이터를 찾을 수 없음
- `ALREADY_EXISTS`: 이미 존재하는 데이터
- `DATABASE_CONNECTION_ERROR`: 데이터베이스 연결 오류

## 🔌 API 엔드포인트

### 카테고리 관리

```
GET    /categories                 # 카테고리 목록 조회
POST   /categories/create          # 카테고리 생성
PUT    /categories/update          # 카테고리 수정
DELETE /categories/delete          # 카테고리 삭제
DELETE /categories/force-delete    # 카테고리 강제 삭제
```

### 키워드 관리

```
GET    /keywords/suggestions       # 키워드 제안
GET    /keywords/category          # 카테고리별 키워드
POST   /keywords/upsert           # 키워드 생성/수정
DELETE /keywords/delete           # 키워드 삭제
```

### 결제수단 관리

```
GET    /payment-methods           # 결제수단 목록
GET    /payment-methods/type      # 타입별 결제수단
POST   /payment-methods/create    # 결제수단 생성
PUT    /payment-methods/update    # 결제수단 수정
DELETE /payment-methods/delete    # 결제수단 삭제
PUT    /payment-methods/toggle    # 결제수단 활성화/비활성화
```

### 은행계좌 관리

```
GET    /bank-accounts             # 은행계좌 목록
POST   /bank-accounts/create      # 은행계좌 생성
PUT    /bank-accounts/update      # 은행계좌 수정
DELETE /bank-accounts/delete      # 은행계좌 삭제
PUT    /bank-accounts/toggle      # 은행계좌 활성화/비활성화
```

### 지출 관리 (v2)

```
POST   /v2/out-account/insert     # 지출 데이터 추가
GET    /v2/out-account            # 일별 지출 조회
GET    /v2/month-out-account      # 월별 지출 조회
PUT    /v2/out-account/update     # 지출 데이터 수정
DELETE /v2/out-account/delete     # 지출 데이터 삭제
```

### 수입 관리 (v2)

```
POST   /v2/in-account/insert      # 수입 데이터 추가
GET    /v2/in-account             # 일별 수입 조회
GET    /v2/month-in-account       # 월별 수입 조회
PUT    /v2/in-account/update      # 수입 데이터 수정
DELETE /v2/in-account/delete      # 수입 데이터 삭제
```

### 통계

```
GET    /statistics                        # 기본 통계
GET    /statistics/category-keywords      # 카테고리-키워드 통계
```

## 📦 프로젝트 구조

```
iksoon_account_backend/
├── main.go                    # 애플리케이션 진입점
├── handlers/                  # HTTP 핸들러
│   ├── category_handler.go   # 카테고리 관리
│   ├── keyword_handler.go    # 키워드 관리
│   ├── payment_method_handler.go  # 결제수단 관리
│   ├── bank_account_handler.go    # 은행계좌 관리
│   ├── out_account_handler.go     # 지출 관리
│   ├── in_account_handler.go      # 수입 관리
│   └── statistics_handler.go     # 통계
├── database/                  # 데이터베이스 레이어
│   ├── connection.go         # DB 연결 관리
│   ├── category_repository.go # 카테고리 저장소
│   ├── keyword_repository.go  # 키워드 저장소
│   ├── payment_method_repository.go  # 결제수단 저장소
│   ├── bank_account_repository.go    # 은행계좌 저장소
│   ├── out_account_repository.go     # 지출 저장소
│   ├── in_account_repository.go      # 수입 저장소
│   └── statistics_repository.go     # 통계 저장소
├── models/                    # 데이터 모델
│   └── types.go              # 공통 타입 정의
├── errors/                    # 에러 관리
│   └── error_codes.go        # 에러 코드 정의
├── utils/                     # 유틸리티
│   ├── logger.go             # 로깅 시스템
│   ├── response.go           # HTTP 응답 유틸
│   └── time.go               # 시간 유틸리티
└── go.mod                     # Go 모듈 정의
```

## 🔧 개발 가이드

### 새로운 핸들러 추가

1. `handlers/` 디렉토리에 새 핸들러 파일 생성
2. 인터페이스 정의 및 구현
3. `main.go`에 라우트 등록
4. 새로운 에러 코드 정의 (필요 시)

### 로깅 사용법

```go
// 디버그 로그 (개발환경에서만)
utils.Debug("사용자 요청: %s", userID)

// 정보 로그
utils.Info("서버 시작됨: %s", port)

// 경고 로그
utils.Warning("사용자 인증 실패: %s", userID)

// 에러 로그
utils.Error("데이터베이스 연결 실패: %v", err)

// 데이터베이스 에러 전용
utils.LogDatabaseError("사용자 조회", err)
```

### 에러 처리

```go
// 새로운 에러 코드 사용
utils.SendError(w, apiErrors.ErrCategoryNotFound)

// 에러 메시지 커스터마이징
utils.SendError(w, apiErrors.ErrInvalidData.WithMessage("잘못된 카테고리 ID"))

// 상세 정보 추가
utils.SendError(w, apiErrors.ErrDatabaseConnection.WithDetails("카테고리 조회 실패"))
```

## 📈 성능 최적화

- **HTTP 미들웨어**: 요청/응답 로깅 및 CORS 처리
- **구조화된 에러 처리**: 일관된 에러 응답
- **레벨별 로깅**: 운영환경에서 불필요한 로그 제거
- **리포지토리 패턴**: 데이터 접근 계층 분리

## 🔒 보안 고려사항

- CORS 설정으로 크로스 오리진 요청 제어
- SQL 인젝션 방지를 위한 Prepared Statement 사용
- 입력 검증 및 사니타이징
- 에러 메시지에서 민감한 정보 노출 방지

## 🚀 배포

### Docker 배포 (예정)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o account_server .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/account_server .
EXPOSE 8080
CMD ["./account_server"]
```

## 📝 로그 예시

### 개발 환경

```
2024-01-08 14:30:25 [INFO] [ACCOUNT_API] [main.go:27] 데이터베이스 연결 성공: ./account_app.db
2024-01-08 14:30:25 [INFO] [ACCOUNT_API] [logger.go:145] 서버가 8080 포트에서 실행 중입니다...
2024-01-08 14:30:25 [INFO] [ACCOUNT_API] [logger.go:146] 새로운 구조의 API가 적용되었습니다.
2024-01-08 14:30:25 [INFO] [ACCOUNT_API] [logger.go:147] 로그 레벨: DEBUG
2024-01-08 14:30:30 [INFO] [ACCOUNT_API] [logger.go:165] HTTP GET /categories from 127.0.0.1:52483
2024-01-08 14:30:30 [DEBUG] [ACCOUNT_API] [category_handler.go:35] 카테고리 조회 요청: type=out
2024-01-08 14:30:30 [DEBUG] [ACCOUNT_API] [category_handler.go:44] 카테고리 조회 성공: 5개
```

### 운영 환경

```
2024-01-08 14:30:25 [INFO] [ACCOUNT_API] [main.go:27] 데이터베이스 연결 성공: ./account_app.db
2024-01-08 14:30:25 [INFO] [ACCOUNT_API] [logger.go:145] 서버가 8080 포트에서 실행 중입니다...
2024-01-08 14:35:12 [ERROR] [ACCOUNT_API] [logger.go:198] Database 카테고리 조회 failed: database connection lost
```

## 📄 라이선스

이 프로젝트는 MIT 라이선스 하에 배포됩니다.

---

**Made with ❤️ using Go & SQLite**
