# 스마트 가계부 시스템

Vue.js + Go + SQLite 기반의 개인 가계부 관리 시스템입니다.

## 🏗️ 시스템 아키텍처

### 개발 환경

```
브라우저 → Frontend(localhost:8081) → Backend(localhost:8080)
```

### 운영 환경 (Docker 프록시)

```
브라우저 → Frontend(Nginx:3000) → /api 프록시 → Backend(내부 네트워크:8080)
```

**보안 강화**: 운영 환경에서 Backend는 내부 네트워크에서만 동작하며, Frontend(Nginx)가 프록시 역할을 수행합니다.

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
# Docker 컨테이너 빌드 및 실행 (프록시 구조)
chmod +x deploy.sh
./deploy.sh
```

### 2. 배포 상태 확인

```bash
# 서비스 상태 확인
./deploy.sh --status

# 로그 확인
docker compose logs -f

# 프록시 설정 테스트
./test_proxy_setup.sh  # Linux
test_proxy_setup.bat   # Windows
```

### 3. 서비스 중지

```bash
./deploy.sh --stop
```

### 4. 네트워크 구조 확인

```bash
# 내부 네트워크 확인
docker network ls | grep iksoon-network

# 컨테이너 간 통신 확인
docker exec iksoon-account-frontend wget -q --spider http://iksoon-backend:8080/health
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
│   ├── nginx.conf                   # Nginx 프록시 설정 (/api → backend)
│   └── Dockerfile                   # 컨테이너 설정
├── docker-compose.yml               # 내부 네트워크 & 프록시 구성
├── deploy.sh                        # 운영 배포 스크립트
├── run_dev.sh                       # 개발 환경 실행 (Linux/Mac)
└── run_dev.bat                      # 개발 환경 실행 (Windows)
```

## ⚙️ 환경 설정

### 개발 환경 (`config.env.development`)

- Backend 포트: 8080 (직접 접근)
- Frontend 포트: 8081 (npm run serve)
- DB 경로: `./data/account_app_dev.db`
- 로그 레벨: DEBUG
- API 통신: Frontend → Backend 직접 연결

### 운영 환경 (`config.env.production`)

- Frontend 포트: 3000 (Nginx, 외부 접근 가능)
- Backend 포트: 8080 (내부 네트워크만, 외부 차단)
- DB 경로: `/db/account_app.db` (Docker 볼륨)
- 로그 레벨: ERROR
- API 통신: Frontend(Nginx) /api 프록시 → Backend
- 네트워크: iksoon-network (172.20.0.0/16)

## 🌐 접속 정보

### 개발 환경

- **Frontend**: http://localhost:8081 (npm run serve)
- **Backend API**: http://localhost:8080 (직접 접근)
- **API 통신**: Frontend에서 Backend로 직접 호출

### 운영 환경 (Docker 프록시)

- **Frontend**: http://133.186.153.179:3000 (Nginx, 외부 접근)
- **API 통신**: http://133.186.153.179:3000/api → Backend (프록시)
- **Backend**: 내부 네트워크만 (외부 접근 차단)
- **데이터 저장**: `/db/` 디렉토리 (Docker 볼륨)

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

### 프록시 설정 문제 해결

```bash
# 프록시 설정 테스트
./test_proxy_setup.sh   # Linux
test_proxy_setup.bat    # Windows

# 컨테이너 로그 확인
docker logs iksoon-account-frontend
docker logs iksoon-account-backend

# Nginx 프록시 로그 확인
docker exec iksoon-account-frontend tail -f /var/log/nginx/api_proxy.log
```

### 개발 환경 설정 확인

```bash
# Windows
run_dev.bat

# Linux/Mac
./run_dev.sh
```

### 네트워크 문제 해결

```bash
# 내부 네트워크 재생성
docker-compose down
docker network prune
docker-compose up -d

# 컨테이너 간 통신 확인
docker exec iksoon-account-frontend ping iksoon-backend
```

---

**Made with ❤️ using Vue.js, Go & SQLite**
