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

# 완전 캐시 제거 후 빌드 (권장: 문제 해결용)
chmod +x clean_deploy.sh
./clean_deploy.sh
```

#### **배포 옵션**

- `./deploy.sh` - 일반 배포 (기존 캐시 활용)
- `./deploy.sh --clean` - 이미지 제거 후 빌드
- `./deploy.sh --force-clean` - 🧹 **모든 Docker 캐시 완전 제거 후 빌드**
- `./clean_deploy.sh` - 간편 완전 캐시 제거 배포

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
# 네트워크 통신 종합 테스트 (추천)
./test_network_communication.sh    # Linux/Mac
test_network_communication.bat     # Windows

# 수동 네트워크 확인
docker network ls | grep iksoon-network
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
├── deploy.sh                        # 운영 배포 스크립트 (캐시 옵션 지원)
├── deploy.ps1                       # 운영 배포 스크립트 (Windows)
├── clean_deploy.sh                  # 완전 캐시 제거 배포 (Linux/Mac)
├── clean_deploy.bat                 # 완전 캐시 제거 배포 (Windows)
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
- ✅ 통합 관리 페이지
  - 사용자, 카테고리, 키워드, 결제수단, 입금경로 통합 관리
  - 카테고리별 고정/변동 지출 유형 설정
  - 직관적인 UI로 모든 설정을 한 곳에서 관리
- ✅ 카테고리별 예산 관리
  - 카테고리 삭제 시 사용 중인 데이터 자동 확인
  - 명확한 안내: 삭제 불가 시 구체적인 해결 방법 제시
  - 데이터 무결성 보장: 연결된 데이터가 있는 카테고리는 삭제 불가
- ✅ 키워드 자동완성
- ✅ 통계 및 리포트
  - **✨ 유연한 기간 설정**: 주간, 월, 년도 개별 선택 가능
  - 주간: 특정 년도의 특정 주차 선택 (예: 2024년 1주차, 2주차...)
  - 월: 특정 년도의 특정 월 선택 (예: 2024년 1월, 2월...)
  - 년도: 특정 년도 선택 (예: 2023년, 2024년...)
  - 커스텀 기간: 시작일-종료일 직접 설정
- ✅ 결제수단 관리
- ✅ 입금경로 관리

### 📊 통계 대시보드 기간 설정

통계 대시보드에서는 다음과 같은 유연한 기간 설정이 가능합니다:

#### 1. 주간 선택

- **년도 선택**: 현재 년도 ±5년 범위에서 선택
- **주차 선택**: 선택한 년도의 1-53주차 중 선택
- **예시**: "2024년 12주차" 형태로 특정 주간 통계 조회

#### 2. 월 선택

- **년도 선택**: 현재 년도 ±5년 범위에서 선택
- **월 선택**: 1월~12월 중 선택
- **예시**: "2024년 3월" 형태로 특정 월 통계 조회

#### 3. 년도 선택

- **년도 선택**: 현재 년도 ±5년 범위에서 선택
- **예시**: "2023년" 형태로 연간 통계 조회

#### 4. 기간 설정 (커스텀)

- **시작일/종료일**: 달력에서 직접 날짜 범위 선택
- **예시**: "2024-01-15 ~ 2024-02-28" 형태로 임의 기간 통계 조회

#### 5. 전체 기간

- 2020년부터 현재까지 모든 데이터 조회

#### 💡 사용법

1. 통계 대시보드에서 원하는 기간 타입 선택 (주간/월/년도/기간 설정/전체)
2. 선택한 타입에 따라 나타나는 드롭다운에서 구체적인 기간 선택
3. 자동으로 해당 기간의 통계가 조회됨
4. 사용자별, 수입/지출별 필터링도 함께 적용 가능

## 📋 사전 요구사항

### 개발 환경

- Go 1.21+
- Node.js 18+
- npm

### 운영 환경 (Ubuntu)

- Docker
- Docker Compose v2

## 🔧 트러블슈팅

### Docker 캐시 문제 해결

```bash
# 🧹 모든 캐시 제거 후 재배포 (문제 해결 1순위)
./clean_deploy.sh       # Linux/Mac
clean_deploy.bat        # Windows

# 또는 수동으로 캐시 제거
./deploy.sh --force-clean    # Linux/Mac
.\deploy.ps1 -ForceClean     # Windows
```

### 프록시 설정 문제 해결

```bash
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
