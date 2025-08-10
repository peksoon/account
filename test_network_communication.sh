#!/bin/bash

# 운영 환경 네트워크 통신 테스트 스크립트

set -e

# 색깔 정의
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔍 운영 환경 네트워크 통신 테스트${NC}"
echo "========================================"

# 1. 컨테이너 상태 확인
echo -e "${YELLOW}📊 컨테이너 상태:${NC}"
docker-compose ps

echo ""
echo -e "${YELLOW}🔍 네트워크 정보:${NC}"
docker network inspect iksoon-network --format '{{range .Containers}}{{.Name}}: {{.IPv4Address}}{{"\n"}}{{end}}' 2>/dev/null || echo "네트워크가 존재하지 않습니다."

echo ""
echo -e "${YELLOW}🌐 외부 접근 테스트:${NC}"

# 2. Frontend 외부 접근 테스트
echo -n "Frontend 접근 (http://localhost:3000): "
if curl -s -o /dev/null -w "%{http_code}" http://localhost:3000 | grep -q "200"; then
    echo -e "${GREEN}✅ 성공${NC}"
else
    echo -e "${RED}❌ 실패${NC}"
fi

# 3. Backend 직접 접근 테스트 (실패해야 정상)
echo -n "Backend 직접 접근 (http://localhost:8080): "
if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health 2>/dev/null | grep -q "200"; then
    echo -e "${RED}❌ 노출됨 (보안 문제!)${NC}"
else
    echo -e "${GREEN}✅ 차단됨 (정상)${NC}"
fi

# 4. API 프록시 테스트
echo -n "API 프록시 접근 (http://localhost:3000/api/health): "
PROXY_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/api/health 2>/dev/null)
if [ "$PROXY_STATUS" = "200" ]; then
    echo -e "${GREEN}✅ 성공 (프록시 정상)${NC}"
elif [ "$PROXY_STATUS" = "502" ]; then
    echo -e "${YELLOW}⚠️ 502 Bad Gateway (Backend 연결 문제)${NC}"
elif [ "$PROXY_STATUS" = "404" ]; then
    echo -e "${YELLOW}⚠️ 404 Not Found (엔드포인트 없음)${NC}"
else
    echo -e "${RED}❌ 실패 (HTTP $PROXY_STATUS)${NC}"
fi

echo ""
echo -e "${YELLOW}🔄 컨테이너 간 내부 통신 테스트:${NC}"

# 5. Frontend에서 Backend로 내부 통신 테스트
echo -n "Frontend → Backend 내부 통신: "
if docker exec iksoon-account-frontend wget -q --spider http://iksoon-backend:8080/health 2>/dev/null; then
    echo -e "${GREEN}✅ 성공${NC}"
else
    echo -e "${RED}❌ 실패${NC}"
fi

# 6. DNS 해석 테스트
echo -n "Docker 내부 DNS (iksoon-backend): "
if docker exec iksoon-account-frontend nslookup iksoon-backend 2>/dev/null | grep -q "Address"; then
    echo -e "${GREEN}✅ 해석됨${NC}"
    docker exec iksoon-account-frontend nslookup iksoon-backend 2>/dev/null | grep "Address" | tail -1
else
    echo -e "${RED}❌ 해석 실패${NC}"
fi

echo ""
echo -e "${YELLOW}📋 프록시 로그 확인:${NC}"
if docker exec iksoon-account-frontend test -f /var/log/nginx/api_proxy.log; then
    echo "최근 API 프록시 로그 (최대 5줄):"
    docker exec iksoon-account-frontend tail -n 5 /var/log/nginx/api_proxy.log 2>/dev/null || echo "로그가 비어있음"
else
    echo "API 프록시 로그 파일이 없습니다."
fi

echo ""
echo -e "${YELLOW}🔍 포트 매핑 정보:${NC}"
echo "Frontend 포트 매핑:"
docker port iksoon-account-frontend 2>/dev/null || echo "포트 매핑 없음"
echo "Backend 포트 매핑:"
docker port iksoon-account-backend 2>/dev/null || echo "포트 매핑 없음 (정상 - 보안)"

echo ""
echo -e "${BLUE}📝 네트워크 구조 요약:${NC}"
echo "🌐 외부 → Frontend(3000) ✅"
echo "🔒 외부 → Backend(8080) ❌ (차단됨)"
echo "🔄 Frontend → Backend (내부) ✅"
echo "📁 Backend → SQLite DB ✅"

echo ""
echo -e "${GREEN}🎯 테스트 완료!${NC}"
