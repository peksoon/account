# 🏠 스마트 가계부 (Smart Household Account Book)

Vue.js와 Go로 개발된 현대적인 가계부 관리 시스템입니다.

## 📋 목차

- [주요 기능](#주요-기능)
- [기술 스택](#기술-스택)
- [시스템 요구사항](#시스템-요구사항)
- [설치 및 실행](#설치-및-실행)
- [프로젝트 구조](#프로젝트-구조)
- [API 문서](#api-문서)
- [화면 구성](#화면-구성)
- [데이터베이스 구조](#데이터베이스-구조)
- [개발 정보](#개발-정보)

## ✨ 주요 기능

### 💰 가계부 관리

- **수입/지출 관리**: 일별 수입/지출 내역 기록 및 조회
- **달력 뷰**: 월별 달력으로 직관적인 가계부 확인
- **상세 내역**: 각 거래의 상세 정보 조회 및 편집
- **실시간 업데이트**: 데이터 추가/수정 시 즉시 화면 반영

### 📊 통계 대시보드

- **월별 통계**: 수입/지출 월별 요약 정보
- **카테고리별 분석**: 카테고리별 지출 분석
- **시각적 차트**: 직관적인 그래프로 데이터 시각화

### 🎯 카테고리 및 키워드 관리

- **카테고리 관리**: 수입/지출 카테고리 추가, 수정, 삭제
- **키워드 자동완성**: 입력한 키워드 자동 저장 및 추천
- **계층적 구조**: 체계적인 분류 시스템

### 💳 결제수단 관리

- **계층적 결제수단**: 카드 > 신용카드/체크카드 등 세분화
- **기본 결제수단**: 카드, 계좌이체, 현금, 기타 기본 제공
- **사용자 정의**: 개인 맞춤 결제수단 추가

### 🏦 입금경로 관리

- **입금경로 설정**: 수입원별 입금경로 관리
- **기본 입금경로**: 급여, 용돈, 기타수입 등 기본 제공
- **유연한 관리**: 사용자별 맞춤 입금경로 설정

### 👥 사용자 관리

- **다중 사용자**: 가족 구성원별 가계부 관리
- **데이터 연동**: 사용자명 변경 시 모든 거래 기록 자동 업데이트
- **논리적 삭제**: 사용자 삭제 시 기존 거래 기록 보존

## 🛠 기술 스택

### Frontend

- **Vue.js 3**: Progressive JavaScript Framework
- **Pinia**: 상태 관리
- **Element Plus**: UI 컴포넌트 라이브러리
- **FullCalendar**: 달력 컴포넌트
- **Chart.js**: 차트 라이브러리
- **Axios**: HTTP 클라이언트
- **Lucide Vue**: 아이콘 라이브러리

### Backend

- **Go 1.21+**: 고성능 백엔드 서버
- **SQLite**: 경량 데이터베이스
- **Gorilla Mux**: HTTP 라우터 (선택적)
- **CORS**: Cross-Origin Resource Sharing 지원

### Development Tools

- **npm**: 패키지 매니저
- **Vue CLI**: Vue.js 개발 도구
- **ESLint**: 코드 품질 관리
- **Git**: 버전 관리

## 📋 시스템 요구사항

- **Node.js**: 16.0.0 이상
- **npm**: 8.0.0 이상
- **Go**: 1.21 이상
- **운영체제**: Windows 10+, macOS 10.15+, Ubuntu 20.04+

## 🚀 설치 및 실행

### 1. 프로젝트 클론

```bash
git clone https://github.com/peksoon/iksoon_account.git
cd iksoon_account
```

### 2. 백엔드 실행

```bash
cd iksoon_account_backend
go mod tidy
go run main.go
```

- 서버가 `http://localhost:8080`에서 실행됩니다.

### 3. 프론트엔드 실행

```bash
cd iksoon_account_frontend
npm install
npm run serve
```

- 프론트엔드가 `http://localhost:3000`에서 실행됩니다.

### 4. 접속

브라우저에서 `http://localhost:3000`으로 접속하여 사용합니다.

## 📁 프로젝트 구조

```
iksoon_account/
├── README.md
├── iksoon_account_backend/           # Go 백엔드
│   ├── main.go                      # 메인 서버 파일
│   ├── database/                    # 데이터베이스 관련
│   │   ├── connection.go           # DB 연결 및 초기화
│   │   ├── user_repository.go      # 사용자 데이터 처리
│   │   ├── category_repository.go  # 카테고리 데이터 처리
│   │   ├── keyword_repository.go   # 키워드 데이터 처리
│   │   ├── payment_method_repository.go # 결제수단 데이터 처리
│   │   ├── deposit_path_repository.go   # 입금경로 데이터 처리
│   │   ├── out_account_repository.go    # 지출 데이터 처리
│   │   ├── in_account_repository.go     # 수입 데이터 처리
│   │   └── statistics_repository.go     # 통계 데이터 처리
│   ├── handlers/                    # HTTP 핸들러
│   │   ├── user_handler.go
│   │   ├── category_handler.go
│   │   ├── keyword_handler.go
│   │   ├── payment_method_handler.go
│   │   ├── deposit_path_handler.go
│   │   ├── out_account_handler.go
│   │   ├── in_account_handler.go
│   │   └── statistics_handler.go
│   ├── models/                      # 데이터 모델
│   │   └── types.go
│   ├── utils/                       # 유틸리티 함수
│   │   └── response.go
│   ├── go.mod
│   └── go.sum
└── iksoon_account_frontend/          # Vue.js 프론트엔드
    ├── public/                      # 정적 파일
    ├── src/
    │   ├── components/              # Vue 컴포넌트
    │   │   ├── AccountCalendar.vue  # 메인 달력
    │   │   ├── AddPopup.vue         # 수입/지출 추가
    │   │   ├── DetailPopup.vue      # 상세 정보
    │   │   ├── StatisticsDashboard.vue # 통계 대시보드
    │   │   ├── UserManager.vue      # 사용자 관리
    │   │   ├── CategoryManager.vue  # 카테고리 관리
    │   │   ├── CategoryList.vue     # 카테고리 목록
    │   │   ├── PaymentMethodManager.vue # 결제수단 관리
    │   │   ├── DepositPathManager.vue   # 입금경로 관리
    │   │   └── KeywordAutocomplete.vue  # 키워드 자동완성
    │   ├── stores/                  # Pinia 상태 관리
    │   │   ├── userStore.js
    │   │   ├── categoryStore.js
    │   │   ├── keywordStore.js
    │   │   ├── paymentMethodStore.js
    │   │   ├── depositPathStore.js
    │   │   ├── accountStore.js
    │   │   ├── calendarStore.js
    │   │   └── popupStore.js
    │   ├── App.vue                  # 루트 컴포넌트
    │   └── main.js                  # 앱 진입점
    ├── package.json
    ├── vue.config.js               # Vue 설정
    └── jsconfig.json              # JavaScript 설정
```

## 🔌 API 문서

### 사용자 관리

- `GET /users` - 사용자 목록 조회
- `POST /users/create` - 사용자 생성
- `PUT /users/update?id={id}` - 사용자 수정
- `DELETE /users/delete?id={id}` - 사용자 삭제 (논리적)
- `DELETE /users/force-delete?id={id}` - 사용자 강제 삭제

### 카테고리 관리

- `GET /categories` - 카테고리 목록 조회
- `POST /categories/create` - 카테고리 생성
- `PUT /categories/update?id={id}` - 카테고리 수정
- `DELETE /categories/delete?id={id}` - 카테고리 삭제

### 결제수단 관리

- `GET /payment-methods` - 결제수단 목록 조회 (계층 구조)
- `POST /payment-methods/create` - 결제수단 생성
- `PUT /payment-methods/update?id={id}` - 결제수단 수정
- `DELETE /payment-methods/delete?id={id}` - 결제수단 삭제

### 입금경로 관리

- `GET /deposit-paths` - 입금경로 목록 조회
- `POST /deposit-paths/create` - 입금경로 생성
- `PUT /deposit-paths/update?id={id}` - 입금경로 수정
- `DELETE /deposit-paths/delete?id={id}` - 입금경로 삭제

### 가계부 관리

- `GET /out-accounts?date={YYYY-MM-DD}` - 지출 내역 조회
- `POST /out-accounts/create` - 지출 내역 생성
- `PUT /out-accounts/update` - 지출 내역 수정
- `DELETE /out-accounts/delete?uuid={uuid}` - 지출 내역 삭제

- `GET /in-accounts?date={YYYY-MM-DD}` - 수입 내역 조회
- `POST /in-accounts/create` - 수입 내역 생성
- `PUT /in-accounts/update` - 수입 내역 수정
- `DELETE /in-accounts/delete?uuid={uuid}` - 수입 내역 삭제

### 통계

- `GET /statistics/monthly?year={YYYY}&month={MM}` - 월별 통계

### 키워드

- `GET /keywords` - 키워드 목록 조회

## 🖥 화면 구성

### 메인 화면 (달력 뷰)

- 월별 달력으로 수입/지출 현황 확인
- 일별 합계 금액 표시
- 특정 날짜 클릭 시 상세 내역 확인

### 수입/지출 추가 팝업

- 날짜, 금액, 사용자, 카테고리, 키워드 입력
- 지출: 결제수단 선택
- 수입: 입금경로 선택
- 실시간 키워드 자동완성

### 상세 내역 팝업

- 거래 상세 정보 조회
- 편집 모드로 정보 수정 가능
- 거래 삭제 기능

### 통계 대시보드

- 월별 수입/지출 요약
- 카테고리별 지출 분석
- 차트를 통한 시각적 데이터 표현

### 관리 화면들

- **사용자 관리**: 가족 구성원 추가/수정/삭제
- **카테고리 관리**: 수입/지출 카테고리 관리
- **결제수단 관리**: 계층적 결제수단 관리
- **입금경로 관리**: 수입원별 입금경로 관리

## 🗄 데이터베이스 구조

### 주요 테이블

- `users`: 사용자 정보
- `categories`: 카테고리 정보 (수입/지출 구분)
- `keywords`: 키워드 정보
- `payment_methods`: 결제수단 정보 (계층 구조)
- `deposit_paths`: 입금경로 정보
- `out_account_data`: 지출 내역
- `in_account_data`: 수입 내역

### 데이터 무결성

- 외래키 제약조건으로 데이터 일관성 보장
- 논리적 삭제(`is_active` 플래그)로 기존 거래 기록 보존
- 사용자명 변경 시 모든 관련 거래 기록 자동 업데이트

## 👨‍💻 개발 정보

### 개발 원칙

- **한국 시간(KST) 기준**: 모든 시간 데이터는 KST로 처리
- **논리적 삭제**: 사용자, 카테고리, 키워드 삭제 시 기존 거래 기록 보존
- **에러 코드 기반**: 문자열 에러 대신 구조화된 에러 코드 사용
- **백엔드 중심**: 모든 설정과 로직은 백엔드에 구현, 프론트엔드는 요청만 처리
- **파일 크기 제한**: 각 파일은 1000줄 이하로 역할에 맞게 분할
- **로그 레벨**: 운영 환경에서 불필요한 로그 생성 방지

### 호환성

- 프론트엔드와 백엔드 API 완전 호환
- 데이터 구조 변경 시 양쪽 모두 동기화
- 기존 데이터 마이그레이션 지원

### 코드 품질

- ESLint를 통한 코드 품질 관리
- Go 표준 코딩 컨벤션 준수
- 함수별 주석 및 문서화
- 중복 코드 제거 및 모듈화

## 📝 라이선스

이 프로젝트는 개인 사용 목적으로 개발되었습니다.

## 🤝 기여하기

버그 리포트나 기능 제안은 GitHub Issues를 통해 알려주세요.

---

**개발자**: peksoon  
**연락처**: GitHub Issues를 통해 문의해 주세요.  
**최종 업데이트**: 2024년 12월
