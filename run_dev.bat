@echo off
REM 개발 환경 실행 스크립트 (Windows)
REM 환경별 설정 파일 사용

echo [INFO] 개발 환경 백엔드 서버 시작 중...

REM 백엔드 디렉토리로 이동
cd iksoon_account_backend

REM config.env (개발용) 파일 확인
if exist "config.env" (
    echo [SUCCESS] config.env ^(개발용^) 파일을 찾았습니다.
    echo [CONFIG] 개발 환경 설정 로드 중...
    
    echo === 개발 환경 설정 ===
    type config.env | findstr /v "^#" | findstr /v "^$"
    echo =====================
    echo.
) else (
    echo [WARNING] config.env 파일이 없습니다. 기본 설정으로 실행합니다.
)

echo [INFO] Go 백엔드 서버 시작...
go run main.go

pause