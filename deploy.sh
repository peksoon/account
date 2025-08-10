#!/bin/bash

# 스마트 가계부 Docker 배포 스크립트 (Linux/macOS)
# Host Network 모드로 실행

set -e  # 에러 발생 시 스크립트 종료

# 색깔 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 로그 함수들
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 프로젝트 정보
PROJECT_NAME="iksoon-account"
BACKEND_IMAGE="${PROJECT_NAME}-backend"
FRONTEND_IMAGE="${PROJECT_NAME}-frontend"
DATA_DIR="./data"

# 함수: Docker 설치 확인
check_docker() {
    log_info "Docker 설치 상태 확인 중..."
    if ! command -v docker &> /dev/null; then
        log_error "Docker가 설치되어 있지 않습니다."
        log_info "Ubuntu 24.04 Docker 설치 방법:"
        log_info "  sudo apt update"
        log_info "  sudo apt install docker.io docker-compose-v2 -y"
        log_info "  sudo systemctl enable docker"
        log_info "  sudo usermod -aG docker \$USER"
        log_info "  # 로그아웃 후 재로그인 필요"
        exit 1
    fi
    
    # Docker Compose v2 확인 (Ubuntu 24.04 기본)
    if command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    elif docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
    else
        log_error "Docker Compose가 설치되어 있지 않습니다."
        log_info "설치 명령: sudo apt install docker-compose-v2 -y"
        exit 1
    fi
    
    # Go 설치 확인 (go.sum 생성용)
    if ! command -v go &> /dev/null; then
        log_warning "Go가 설치되어 있지 않습니다. go.sum 자동 생성을 건너뜁니다."
        GO_AVAILABLE=false
    else
        GO_AVAILABLE=true
    fi
    
    log_success "Docker 및 Docker Compose가 설치되어 있습니다."
}

# 함수: 기존 컨테이너 정리
cleanup_containers() {
    log_info "기존 컨테이너 정리 중..."
    
    # 실행 중인 컨테이너 확인 및 정지
    if docker ps | grep -q "${PROJECT_NAME}"; then
        log_warning "실행 중인 ${PROJECT_NAME} 컨테이너를 정지합니다."
        $COMPOSE_CMD down
    fi
    
    # 기존 이미지 제거 (선택적)
    if [[ "$1" == "--clean" ]]; then
        log_info "기존 이미지 제거 중..."
        docker rmi ${BACKEND_IMAGE} 2>/dev/null || true
        docker rmi ${FRONTEND_IMAGE} 2>/dev/null || true
        docker system prune -f
    fi
    
    log_success "컨테이너 정리 완료"
}

# 함수: 데이터 디렉토리 생성
create_data_directory() {
    log_info "데이터 디렉토리 생성 중..."
    mkdir -p ${DATA_DIR}
    chmod 755 ${DATA_DIR}
    log_success "데이터 디렉토리 생성 완료: ${DATA_DIR}"
}

# 함수: 이미지 빌드
build_images() {
    log_info "Docker 이미지 빌드 시작..."
    
    # Backend 빌드 준비 (go.sum 파일 확인)
    log_info "Backend 빌드 환경 준비 중..."
    if [ ! -f "./iksoon_account_backend/go.sum" ] && [ "$GO_AVAILABLE" = true ]; then
        log_warning "go.sum 파일이 없습니다. 생성 중..."
        cd ./iksoon_account_backend
        go mod tidy
        cd ..
        log_success "go.sum 파일 생성 완료"
    elif [ ! -f "./iksoon_account_backend/go.sum" ]; then
        log_error "go.sum 파일이 없고 Go도 설치되어 있지 않습니다."
        log_info "해결 방법:"
        log_info "1. Go 설치: sudo apt install golang-go -y"
        log_info "2. 또는 개발 환경에서 go.sum 파일을 생성하여 복사"
        exit 1
    fi
    
    # Backend 이미지 빌드
    log_info "Backend 이미지 빌드 중..."
    docker build -t ${BACKEND_IMAGE} ./iksoon_account_backend/
    log_success "Backend 이미지 빌드 완료"
    
    # Frontend 이미지 빌드
    log_info "Frontend 이미지 빌드 중..."
    docker build -t ${FRONTEND_IMAGE} ./iksoon_account_frontend/
    log_success "Frontend 이미지 빌드 완료"
    
    log_success "모든 이미지 빌드 완료"
}

# 함수: 컨테이너 실행
start_containers() {
    log_info "컨테이너 시작 중..."
    $COMPOSE_CMD up -d
    
    # 헬스체크 대기
    log_info "서비스 시작 대기 중..."
    sleep 10
    
    # Backend 헬스체크
    log_info "Backend 서비스 확인 중..."
    for i in {1..30}; do
        if curl -s http://localhost:8080/health >/dev/null 2>&1; then
            log_success "Backend 서비스가 정상적으로 시작되었습니다."
            break
        fi
        if [ $i -eq 30 ]; then
            log_error "Backend 서비스 시작 실패"
            exit 1
        fi
        sleep 2
    done
    
    # Frontend 헬스체크
    log_info "Frontend 서비스 확인 중..."
    for i in {1..30}; do
        if curl -s http://localhost:3000 >/dev/null 2>&1; then
            log_success "Frontend 서비스가 정상적으로 시작되었습니다."
            break
        fi
        if [ $i -eq 30 ]; then
            log_error "Frontend 서비스 시작 실패"
            exit 1
        fi
        sleep 2
    done
}

# 함수: 컨테이너 상태 확인
check_status() {
    log_info "컨테이너 상태 확인 중..."
    $COMPOSE_CMD ps
    
    echo ""
    log_info "서비스 접속 정보:"
    echo "  🌐 Frontend: http://localhost:3000"
    echo "  🔧 Backend API: http://localhost:8080"
    echo ""
    log_info "로그 확인: $COMPOSE_CMD logs -f"
    log_info "서비스 중지: $COMPOSE_CMD down"
}

# 함수: 사용법 출력
usage() {
    echo "사용법: $0 [옵션]"
    echo ""
    echo "옵션:"
    echo "  --clean    기존 이미지를 모두 제거하고 새로 빌드"
    echo "  --stop     실행 중인 컨테이너만 중지"
    echo "  --status   현재 실행 상태 확인"
    echo "  --help     이 도움말 출력"
    echo ""
    echo "예시:"
    echo "  $0              # 일반 배포"
    echo "  $0 --clean      # 완전 새로 빌드 후 배포"
    echo "  $0 --stop       # 서비스 중지"
    echo "  $0 --status     # 상태 확인"
}

# 메인 실행 로직
main() {
    echo ""
    log_info "===== 스마트 가계부 Docker 배포 스크립트 ====="
    echo ""
    
    case "${1:-}" in
        --help|-h)
            usage
            exit 0
            ;;
        --stop)
            cleanup_containers
            log_success "서비스가 중지되었습니다."
            exit 0
            ;;
        --status)
            check_status
            exit 0
            ;;
        --clean)
            check_docker
            cleanup_containers --clean
            create_data_directory
            build_images
            start_containers
            check_status
            ;;
        "")
            check_docker
            cleanup_containers
            create_data_directory
            build_images
            start_containers
            check_status
            ;;
        *)
            log_error "알 수 없는 옵션: $1"
            usage
            exit 1
            ;;
    esac
    
    echo ""
    log_success "배포가 완료되었습니다! 🎉"
    log_info "브라우저에서 http://localhost:3000 에 접속하세요."
}

# 스크립트 실행
main "$@"
