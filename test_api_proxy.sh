#!/bin/bash

echo "=== API 프록시 테스트 ==="

# Frontend 컨테이너가 실행 중인지 확인
if ! docker ps | grep -q "iksoon-account-frontend"; then
    echo "❌ Frontend 컨테이너가 실행되고 있지 않습니다."
    exit 1
fi

# Backend 컨테이너가 실행 중인지 확인
if ! docker ps | grep -q "iksoon-account-backend"; then
    echo "❌ Backend 컨테이너가 실행되고 있지 않습니다."
    exit 1
fi

echo "✅ 컨테이너 상태 정상"

# Frontend를 통한 API 프록시 테스트
echo ""
echo "1. Frontend 직접 접근 테스트:"
curl -s "http://localhost:3000" > /dev/null && echo "✅ Frontend 정상" || echo "❌ Frontend 접근 실패"

echo ""
echo "2. Backend 직접 접근 테스트:"
curl -s "http://localhost:8080/health" > /dev/null && echo "✅ Backend 정상" || echo "❌ Backend 접근 실패"

echo ""
echo "3. Frontend nginx 프록시를 통한 Backend 접근 테스트:"
curl -s "http://localhost:3000/health" > /dev/null && echo "✅ 프록시 정상" || echo "❌ 프록시 실패"

echo ""
echo "4. v2 API 프록시 테스트:"
curl -s "http://localhost:3000/v2/month-out-account?year=2025&month=08" > /dev/null && echo "✅ v2 API 프록시 정상" || echo "❌ v2 API 프록시 실패"

echo ""
echo "=== 테스트 완료 ==="
