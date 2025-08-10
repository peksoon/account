#!/bin/bash

# Docker 캐시 제거 기능 테스트 스크립트

set -e

echo "🧪 Docker 캐시 제거 기능 테스트"
echo "================================"

# 색깔 정의
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}📊 테스트 전 Docker 상태:${NC}"
echo "이미지 개수: $(docker images -q | wc -l)"
echo "컨테이너 개수: $(docker ps -aq | wc -l)"
echo "네트워크 개수: $(docker network ls -q | wc -l)"
echo "볼륨 개수: $(docker volume ls -q | wc -l)"

echo ""
echo -e "${YELLOW}🧹 캐시 제거 테스트 시작...${NC}"

# --force-clean 옵션 테스트
./deploy.sh --force-clean > /dev/null 2>&1 || true

echo ""
echo -e "${BLUE}📊 테스트 후 Docker 상태:${NC}"
echo "이미지 개수: $(docker images -q | wc -l)"
echo "컨테이너 개수: $(docker ps -aq | wc -l)"
echo "네트워크 개수: $(docker network ls -q | wc -l)"
echo "볼륨 개수: $(docker volume ls -q | wc -l)"

echo ""
echo -e "${GREEN}✅ 캐시 제거 기능 테스트 완료${NC}"

# 서비스 상태 확인
if docker ps | grep -q "iksoon-account"; then
    echo -e "${GREEN}✅ 서비스가 정상적으로 실행 중입니다${NC}"
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
else
    echo -e "${YELLOW}⚠️ 서비스가 실행되지 않았습니다${NC}"
fi
