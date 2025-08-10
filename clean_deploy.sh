#!/bin/bash

# 스마트 가계부 완전 캐시 제거 후 배포 스크립트
# 모든 Docker 캐시를 제거하고 깨끗하게 새로 빌드

set -e

# 색깔 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🧹 완전 캐시 제거 후 배포 시작${NC}"
echo "========================================"

# 현재 실행 중인 컨테이너 확인
echo -e "${YELLOW}📊 현재 Docker 상태:${NC}"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo -e "${RED}⚠️  주의: 모든 Docker 캐시가 제거됩니다!${NC}"
echo "   - 빌드 캐시"
echo "   - 이미지 캐시"
echo "   - 네트워크 캐시"
echo "   - 볼륨 캐시"
echo ""

read -p "계속하시겠습니까? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}취소되었습니다.${NC}"
    exit 0
fi

echo ""
echo -e "${GREEN}🚀 완전 캐시 제거 배포를 시작합니다...${NC}"

# 메인 배포 스크립트 실행
./deploy.sh --force-clean

echo ""
echo -e "${GREEN}✅ 완전 캐시 제거 배포가 완료되었습니다!${NC}"
echo -e "${BLUE}📱 브라우저에서 http://localhost:3000 에 접속하세요.${NC}"
