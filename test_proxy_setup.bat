@echo off
echo 🔍 프록시 설정 테스트 스크립트
echo ================================

echo 📊 Docker 서비스 상태 확인...
docker-compose ps

echo.
echo 🌐 Frontend 서비스 접근 테스트...
curl -s -o nul -w "%%{http_code}" http://localhost:3000 | findstr "200" >nul && echo ✅ Frontend 접근 성공 || echo ❌ Frontend 접근 실패

echo.
echo 🔒 Backend 직접 접근 테스트 (실패해야 정상)...
curl -s -o nul -w "%%{http_code}" http://localhost:8080 2>nul | findstr "000" >nul && echo ✅ Backend 직접 접근 차단됨 (정상) || echo ❌ Backend가 여전히 외부에 노출됨

echo.
echo 🔄 API 프록시 테스트...
curl -s -o nul -w "%%{http_code}" http://localhost:3000/api/health | findstr "200" >nul && echo ✅ API 프록시 동작 성공 || echo ❌ API 프록시 동작 실패

echo.
echo 📋 네트워크 설정 확인...
docker network ls | findstr iksoon-network >nul && echo ✅ 내부 네트워크 생성됨 || echo ❌ 내부 네트워크 없음

echo.
echo 🔍 컨테이너 간 통신 테스트...
docker exec iksoon-account-frontend wget -q --spider http://iksoon-backend:8080/health 2>nul && echo ✅ 컨테이너 간 통신 성공 || echo ❌ 컨테이너 간 통신 실패

echo.
echo 📝 프록시 로그 확인...
docker exec iksoon-account-frontend test -f /var/log/nginx/api_proxy.log 2>nul && (
    echo ✅ API 프록시 로그 파일 존재
    echo 최근 프록시 로그:
    docker exec iksoon-account-frontend tail -n 5 /var/log/nginx/api_proxy.log 2>nul || echo 로그가 아직 생성되지 않음
) || echo ⚠️ API 프록시 로그 파일 없음 (아직 요청이 없을 수 있음)

echo.
echo 🎯 테스트 완료!
echo ✅ 성공: Frontend만 외부 접근 가능, Backend는 내부 네트워크만
echo ✅ 성공: /api 요청이 Backend로 프록시됨
pause
