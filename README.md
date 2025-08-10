# 스마트 가계부 시스템

Vue.js + Go + SQLite 기반의 개인 가계부 관리 시스템입니다.

## 🚀 개발 환경 실행 방법

### 1. Windows 개발 환경

```cmd
# config.env.development 설정으로 백엔드 실행
run_dev.bat

# 또는 수동으로
cd iksoon_account_backend
go run main.go
```

### 2. Linux/macOS 개발 환경

```bash
# config.env.development 설정으로 백엔드 실행
chmod +x run_dev.sh
./run_dev.sh

# 또는 수동으로
cd iksoon_account_backend
go run main.go
```

### 3. 프론트엔드 개발 서버 (별도 터미널)

```bash
cd iksoon_account_frontend
npm install
npm run serve
```

## 🐳 운영 환경 배포 (Ubuntu)

### 1. 전체 시스템 배포

```bash
# config.env.production 설정으로 Docker 배포
chmod +x deploy.sh
./deploy.sh
```

### 2. 배포 상태 확인

```bash
# 서비스 상태 확인
./deploy.sh --status

# 로그 확인
docker compose logs -f
```

### 3. 서비스 중지

```bash
./deploy.sh --stop
```

## 📁 프로젝트 구조

```
00.account/
├── iksoon_account_backend/          # Go 백엔드
│   ├── config.env.development       # 개발 환경 설정
│   ├── config.env.production        # 운영 환경 설정
│   ├── config/                      # 설정 관리 패키지
│   ├── handlers/                    # API 핸들러
│   ├── database/                    # DB 레포지토리
│   └── main.go                      # 메인 엔트리포인트
├── iksoon_account_frontend/         # Vue.js 프론트엔드
│   ├── src/                         # 소스 코드
│   ├── nginx.conf                   # Nginx 설정
│   └── Dockerfile                   # 컨테이너 설정
├── docker-compose.yml               # 운영 환경 컨테이너 구성
├── deploy.sh                        # 운영 배포 스크립트
├── run_dev.sh                       # 개발 환경 실행 (Linux/Mac)
└── run_dev.bat                      # 개발 환경 실행 (Windows)
```

## ⚙️ 환경 설정

### 개발 환경 (`config.env.development`)

- 포트: 8080
- DB 경로: `./data/account_app_dev.db`
- 로그 레벨: DEBUG
- 최대 연결: 50

### 운영 환경 (`config.env.production`)

- 포트: 8080
- DB 경로: `/db/account_app.db` (Docker 볼륨)
- 로그 레벨: ERROR
- 최대 연결: 200

## 🌐 접속 정보

### 개발 환경

- 백엔드 API: http://localhost:8080
- 프론트엔드: http://localhost:8081 (npm run serve)

### 운영 환경 (Docker)

- 백엔드 API: http://localhost:8080
- 프론트엔드: http://localhost:3000
- 데이터 저장: `./data/` 디렉토리

## 🛠️ 주요 기능

- ✅ 수입/지출 관리
- ✅ 카테고리별 예산 관리
- ✅ 키워드 자동완성
- ✅ 통계 및 리포트
- ✅ 결제수단 관리
- ✅ 입금경로 관리

## 📋 사전 요구사항

### 개발 환경

- Go 1.21+
- Node.js 18+
- npm

### 운영 환경 (Ubuntu)

- Docker
- Docker Compose v2

## 🔧 트러블슈팅

### Frontend 빌드 에러 해결

```bash
# nginx.conf 파일 문제 해결됨 (Dockerfile 수정)
# npm deprecated 경고는 정상 (빌드에 영향 없음)
```

### 개발 환경 설정 확인

```bash
# Windows
run_dev.bat

# Linux/Mac
./run_dev.sh
```

---

**Made with ❤️ using Vue.js, Go & SQLite**
