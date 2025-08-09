# 🐳 Docker 배포 가이드

스마트 가계부 애플리케이션을 Docker 컨테이너로 배포하는 방법을 설명합니다.

## 📋 목차

- [사전 요구사항](#사전-요구사항)
- [빠른 시작](#빠른-시작)
- [배포 방법](#배포-방법)
- [컨테이너 관리](#컨테이너-관리)
- [문제 해결](#문제-해결)
- [네트워크 구성](#네트워크-구성)

## 🔧 사전 요구사항

### Docker 설치

- **Windows**: [Docker Desktop for Windows](https://docs.docker.com/desktop/windows/)
- **macOS**: [Docker Desktop for Mac](https://docs.docker.com/desktop/mac/)
- **Linux**: [Docker Engine](https://docs.docker.com/engine/install/)

### 시스템 요구사항

- **메모리**: 최소 2GB RAM
- **디스크**: 최소 1GB 여유 공간
- **포트**: 3000, 8080 포트가 사용 가능해야 함

## 🚀 빠른 시작

### 1. 저장소 클론

```bash
git clone https://github.com/peksoon/iksoon_account.git
cd iksoon_account
```

### 2. 배포 실행

#### Linux/macOS

```bash
chmod +x deploy.sh
./deploy.sh
```

#### Windows (PowerShell)

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
.\deploy.ps1
```

### 3. 접속

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080

## 📦 배포 방법

### 일반 배포

기존 이미지를 재사용하여 빠르게 배포합니다.

```bash
# Linux/macOS
./deploy.sh

# Windows
.\deploy.ps1
```

### 완전 새로 빌드 (Clean Build)

모든 이미지를 새로 빌드하여 배포합니다.

```bash
# Linux/macOS
./deploy.sh --clean

# Windows
.\deploy.ps1 -Clean
```

### Docker Compose 직접 사용

```bash
# 이미지 빌드 및 컨테이너 시작
docker-compose up -d --build

# 로그 확인
docker-compose logs -f

# 컨테이너 중지
docker-compose down
```

## 🛠 컨테이너 관리

### 상태 확인

```bash
# Linux/macOS
./deploy.sh --status

# Windows
.\deploy.ps1 -Status

# 또는 직접
docker-compose ps
```

### 서비스 중지

```bash
# Linux/macOS
./deploy.sh --stop

# Windows
.\deploy.ps1 -Stop

# 또는 직접
docker-compose down
```

### 로그 확인

```bash
# 모든 서비스 로그
docker-compose logs -f

# 특정 서비스 로그
docker-compose logs -f iksoon-backend
docker-compose logs -f iksoon-frontend
```

### 컨테이너 재시작

```bash
# 특정 서비스 재시작
docker-compose restart iksoon-backend
docker-compose restart iksoon-frontend

# 모든 서비스 재시작
docker-compose restart
```

## 🔍 문제 해결

### 1. 포트 충돌

다른 애플리케이션이 3000 또는 8080 포트를 사용하는 경우:

```bash
# 포트 사용 확인
netstat -tulpn | grep :3000
netstat -tulpn | grep :8080

# Windows에서
netstat -ano | findstr :3000
netstat -ano | findstr :8080
```

해결방법:

- 다른 애플리케이션 중지
- `docker-compose.yml`에서 포트 변경

### 2. Docker 권한 문제 (Linux)

```bash
# Docker 그룹에 사용자 추가
sudo usermod -aG docker $USER

# 재로그인 후 확인
docker ps
```

### 3. 이미지 빌드 실패

```bash
# Docker 캐시 정리
docker system prune -f

# 완전 새로 빌드
./deploy.sh --clean
```

### 4. 컨테이너 실행 실패

```bash
# 컨테이너 상태 확인
docker-compose ps

# 오류 로그 확인
docker-compose logs

# 개별 컨테이너 확인
docker logs iksoon-account-backend
docker logs iksoon-account-frontend
```

### 5. 데이터베이스 문제

```bash
# 데이터 디렉토리 확인
ls -la ./data/

# 권한 수정 (필요시)
chmod -R 755 ./data/
```

## 🌐 네트워크 구성

### Host Network 모드

이 애플리케이션은 `host` 네트워크 모드로 실행됩니다:

**장점:**

- 네트워크 성능 최적화
- 포트 매핑 불필요
- 간단한 서비스 간 통신

**주의사항:**

- 호스트의 포트를 직접 사용
- 포트 충돌 가능성
- Linux에서만 완전 지원 (Windows/macOS는 제한적)

### 포트 정보

- **Frontend (Nginx)**: 3000
- **Backend (Go)**: 8080
- **Health Check**: 8080/health

### API 프록시

Frontend에서 Backend API 호출 시 Nginx가 프록시 역할:

- Frontend: `http://localhost:3000/api/*`
- Backend: `http://localhost:8080/*`

## 📁 파일 구조

```
iksoon_account/
├── docker-compose.yml           # Docker Compose 설정
├── deploy.sh                    # Linux/macOS 배포 스크립트
├── deploy.ps1                   # Windows 배포 스크립트
├── DOCKER_DEPLOYMENT.md         # 이 문서
├── data/                        # 데이터베이스 영구 저장
│   └── account.db              # SQLite 데이터베이스
├── iksoon_account_backend/
│   ├── Dockerfile              # Backend 컨테이너 설정
│   └── .dockerignore           # Backend 빌드 제외 파일
└── iksoon_account_frontend/
    ├── Dockerfile              # Frontend 컨테이너 설정
    ├── nginx.conf              # Nginx 설정
    └── .dockerignore           # Frontend 빌드 제외 파일
```

## 🔒 보안 고려사항

### 프로덕션 환경 권장사항

1. **환경 변수 설정**

   ```bash
   export DB_PATH=/secure/path/account.db
   export GO_ENV=production
   ```

2. **Nginx 보안 헤더** (이미 적용됨)

   - X-Frame-Options
   - X-Content-Type-Options
   - X-XSS-Protection

3. **방화벽 설정**

   ```bash
   # 필요한 포트만 열기
   sudo ufw allow 3000/tcp
   sudo ufw allow 8080/tcp
   ```

4. **SSL/TLS 적용** (프로덕션 시)
   - Reverse Proxy (예: Traefik, Nginx Proxy Manager)
   - Let's Encrypt 인증서

## 📈 모니터링

### 헬스체크

```bash
# Backend 헬스체크
curl http://localhost:8080/health

# Frontend 헬스체크
curl http://localhost:3000
```

### 리소스 모니터링

```bash
# 컨테이너 리소스 사용량
docker stats

# 특정 컨테이너
docker stats iksoon-account-backend iksoon-account-frontend
```

## ❓ 자주 묻는 질문

### Q: Windows에서 host 네트워크가 동작하지 않습니다.

A: Windows/macOS의 Docker Desktop은 host 네트워크를 완전히 지원하지 않습니다. 대신 `docker-compose.yml`의 `network_mode: host`를 제거하고 포트 매핑을 사용하세요:

```yaml
ports:
  - "3000:3000" # frontend
  - "8080:8080" # backend
```

### Q: 데이터가 컨테이너 재시작 후 사라집니다.

A: `./data` 디렉토리가 볼륨으로 마운트되어 있습니다. 해당 디렉토리의 권한을 확인하세요.

### Q: 업데이트는 어떻게 하나요?

A: 최신 코드를 pull하고 clean build로 재배포하세요:

```bash
git pull origin master
./deploy.sh --clean
```

---

**문제가 해결되지 않으면 GitHub Issues에 로그와 함께 문의해 주세요.**
