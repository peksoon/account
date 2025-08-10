@echo off
setlocal enabledelayedexpansion

echo 🔍 운영 환경 네트워크 통신 테스트
echo ========================================

REM 1. 컨테이너 상태 확인
echo 📊 컨테이너 상태:
docker-compose ps

echo.
echo 🔍 네트워크 정보:
docker network inspect iksoon-network --format "{{range .Containers}}{{.Name}}: {{.IPv4Address}}{{end}}" 2>nul || echo 네트워크가 존재하지 않습니다.

echo.
echo 🌐 외부 접근 테스트:

REM 2. Frontend 외부 접근 테스트
echo|set /p="Frontend 접근 (http://localhost:3000): "
curl -s -o nul -w "%%{http_code}" http://localhost:3000 | findstr "200" >nul && echo ✅ 성공 || echo ❌ 실패

REM 3. Backend 직접 접근 테스트 (실패해야 정상)
echo|set /p="Backend 직접 접근 (http://localhost:8080): "
curl -s -o nul -w "%%{http_code}" http://localhost:8080/health 2>nul | findstr "200" >nul && echo ❌ 노출됨 (보안 문제!) || echo ✅ 차단됨 (정상)

REM 4. API 프록시 테스트
echo|set /p="API 프록시 접근 (http://localhost:3000/api/health): "
for /f %%i in ('curl -s -o nul -w "%%{http_code}" http://localhost:3000/api/health 2^>nul') do set PROXY_STATUS=%%i
if "!PROXY_STATUS!"=="200" (
    echo ✅ 성공 (프록시 정상)
) else if "!PROXY_STATUS!"=="502" (
    echo ⚠️ 502 Bad Gateway (Backend 연결 문제)
) else if "!PROXY_STATUS!"=="404" (
    echo ⚠️ 404 Not Found (엔드포인트 없음)
) else (
    echo ❌ 실패 (HTTP !PROXY_STATUS!)
)

echo.
echo 🔄 컨테이너 간 내부 통신 테스트:

REM 5. Frontend에서 Backend로 내부 통신 테스트
echo|set /p="Frontend → Backend 내부 통신: "
docker exec iksoon-account-frontend wget -q --spider http://iksoon-backend:8080/health 2>nul && echo ✅ 성공 || echo ❌ 실패

REM 6. DNS 해석 테스트
echo|set /p="Docker 내부 DNS (iksoon-backend): "
docker exec iksoon-account-frontend nslookup iksoon-backend 2>nul | findstr "Address" >nul && (
    echo ✅ 해석됨
    docker exec iksoon-account-frontend nslookup iksoon-backend 2>nul | findstr "Address"
) || echo ❌ 해석 실패

echo.
echo 📋 프록시 로그 확인:
docker exec iksoon-account-frontend test -f /var/log/nginx/api_proxy.log 2>nul && (
    echo 최근 API 프록시 로그 (최대 5줄):
    docker exec iksoon-account-frontend tail -n 5 /var/log/nginx/api_proxy.log 2>nul || echo 로그가 비어있음
) || echo API 프록시 로그 파일이 없습니다.

echo.
echo 🔍 포트 매핑 정보:
echo Frontend 포트 매핑:
docker port iksoon-account-frontend 2>nul || echo 포트 매핑 없음
echo Backend 포트 매핑:
docker port iksoon-account-backend 2>nul || echo 포트 매핑 없음 (정상 - 보안)

echo.
echo 📝 네트워크 구조 요약:
echo 🌐 외부 → Frontend(3000) ✅
echo 🔒 외부 → Backend(8080) ❌ (차단됨)
echo 🔄 Frontend → Backend (내부) ✅
echo 📁 Backend → SQLite DB ✅

echo.
echo 🎯 테스트 완료!
pause
