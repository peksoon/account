@echo off
REM Windows용 개발 환경 실행 스크립트
REM config.env.development 파일을 읽어서 환경변수로 설정 후 실행

echo [INFO] 개발 환경 설정 로드 중...

REM config.env.development 파일 확인
if exist "iksoon_account_backend\config.env.development" (
    echo [SUCCESS] config.env.development 파일을 찾았습니다.
    
    REM 환경변수 설정 (Windows는 파일에서 직접 읽기 어려우므로 수동 설정)
    set PORT=8080
    set DB_PATH=./data/account_app_dev.db
    set LOG_LEVEL=DEBUG
    set MAX_CONNECTIONS=50
    
    echo [INFO] 설정된 환경변수:
    echo   PORT: %PORT%
    echo   DB_PATH: %DB_PATH%
    echo   LOG_LEVEL: %LOG_LEVEL%
    echo   MAX_CONNECTIONS: %MAX_CONNECTIONS%
    
) else (
    echo [INFO] config.env.development 파일이 없습니다. 기본값으로 실행합니다.
)

echo.
echo [INFO] 백엔드 서버 시작 중...

REM 백엔드 디렉토리로 이동 후 실행
cd iksoon_account_backend
go run main.go

pause
