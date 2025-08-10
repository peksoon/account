#!/bin/bash

# Ubuntu 서버 배포 전 체크 스크립트

echo "=== Ubuntu 서버 배포 체크 ==="

echo "1. go.sum 파일 확인:"
if [ -f "./iksoon_account_backend/go.sum" ]; then
    echo "✅ go.sum 파일이 존재합니다."
    wc -l ./iksoon_account_backend/go.sum
else
    echo "❌ go.sum 파일이 없습니다!"
    echo "해결 방법:"
    echo "1) Windows에서 Ubuntu로 파일 복사 시 go.sum 파일도 포함시키기"
    echo "2) Ubuntu에서 Go 설치 후 go mod tidy 실행"
    
    # Go 설치 여부 확인
    if command -v go &> /dev/null; then
        echo "Go가 설치되어 있습니다. go mod tidy 실행 중..."
        cd ./iksoon_account_backend
        go mod tidy
        cd ..
        echo "✅ go.sum 파일이 생성되었습니다."
    else
        echo "Go가 설치되어 있지 않습니다."
        echo "설치 방법: sudo apt install golang-go -y"
    fi
fi

echo ""
echo "2. 프로젝트 파일 구조 확인:"
echo "Backend 디렉토리:"
ls -la ./iksoon_account_backend/ | head -10

echo ""
echo "3. Docker 및 Docker Compose 확인:"
docker --version
if command -v docker-compose &> /dev/null; then
    echo "Docker Compose v1: $(docker-compose --version)"
elif docker compose version &> /dev/null; then
    echo "Docker Compose v2: $(docker compose version)"
else
    echo "❌ Docker Compose가 설치되어 있지 않습니다."
fi

echo ""
echo "=== 체크 완료 ==="
echo "문제가 없으면 ./deploy.sh 를 실행하세요."
