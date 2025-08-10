@echo off
setlocal enabledelayedexpansion

REM 스마트 가계부 완전 캐시 제거 후 배포 스크립트
REM 모든 Docker 캐시를 제거하고 깨끗하게 새로 빌드

echo 🧹 완전 캐시 제거 후 배포 시작
echo ========================================

REM 현재 실행 중인 컨테이너 확인
echo 📊 현재 Docker 상태:
docker ps --format "table {{.Names}}	{{.Status}}	{{.Ports}}"

echo.
echo ⚠️  주의: 모든 Docker 캐시가 제거됩니다!
echo    - 빌드 캐시
echo    - 이미지 캐시
echo    - 네트워크 캐시
echo    - 볼륨 캐시
echo.

set /p confirm="계속하시겠습니까? (y/N): "
if /i not "%confirm%"=="y" (
    echo 취소되었습니다.
    pause
    exit /b 0
)

echo.
echo 🚀 완전 캐시 제거 배포를 시작합니다...

REM 메인 배포 스크립트 실행
powershell -ExecutionPolicy Bypass -File ".\deploy.ps1" -ForceClean

echo.
echo ✅ 완전 캐시 제거 배포가 완료되었습니다!
echo 📱 브라우저에서 http://localhost:3000 에 접속하세요.
pause
