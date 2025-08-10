#!/bin/bash

# 개발 환경 실행 스크립트 (Linux/macOS)
# 환경별 설정 파일 사용

set -e

# 색깔 정의
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}[INFO]${NC} 개발 환경 백엔드 서버 시작 중..."

# 백엔드 디렉토리로 이동
cd ./iksoon_account_backend

# config.env (개발용) 파일 확인
if [ -f "config.env" ]; then
    echo -e "${GREEN}[SUCCESS]${NC} config.env (개발용) 파일을 찾았습니다."
    echo -e "${YELLOW}[CONFIG]${NC} 개발 환경 설정 로드 중..."
    
    # 설정 내용 표시
    echo "=== 개발 환경 설정 ==="
    cat config.env | grep -v '^#' | grep -v '^$'
    echo "====================="
    echo ""
else
    echo -e "${YELLOW}[WARNING]${NC} config.env 파일이 없습니다. 기본 설정으로 실행합니다."
fi

echo -e "${BLUE}[INFO]${NC} Go 백엔드 서버 시작..."
go run main.go
