#!/bin/bash

# 개발 환경 실행 스크립트
# config.env.development 파일을 읽어서 환경변수로 설정 후 실행

set -e

# 색깔 정의
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}[INFO]${NC} 개발 환경 설정 로드 중..."

# config.env.development 파일 읽기
if [ -f "./iksoon_account_backend/config.env.development" ]; then
    echo -e "${GREEN}[SUCCESS]${NC} config.env.development 파일을 찾았습니다."
    
    # 환경변수 설정
    export $(grep -v '^#' ./iksoon_account_backend/config.env.development | xargs)
    
    echo -e "${BLUE}[INFO]${NC} 설정된 환경변수:"
    echo "  PORT: $PORT"
    echo "  DB_PATH: $DB_PATH"
    echo "  LOG_LEVEL: $LOG_LEVEL"
    echo "  MAX_CONNECTIONS: $MAX_CONNECTIONS"
    
else
    echo -e "${GREEN}[INFO]${NC} config.env.development 파일이 없습니다. 기본값으로 실행합니다."
fi

echo ""
echo -e "${BLUE}[INFO]${NC} 백엔드 서버 시작 중..."

# 백엔드 디렉토리로 이동 후 실행
cd ./iksoon_account_backend
go run main.go
