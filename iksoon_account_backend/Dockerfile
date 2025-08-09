# Backend Dockerfile
FROM golang:1.21-alpine AS builder

# 작업 디렉토리 설정
WORKDIR /app

# 의존성 파일 복사
COPY go.mod go.sum ./

# 의존성 다운로드
RUN go mod download

# 소스 코드 복사
COPY . .

# 애플리케이션 빌드
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# 실행용 최소 이미지
FROM alpine:latest

# 필수 패키지 설치 (SQLite 지원)
RUN apk --no-cache add ca-certificates sqlite

# 작업 디렉토리 설정
WORKDIR /root/

# 빌드된 바이너리 복사
COPY --from=builder /app/main .

# 데이터베이스 디렉토리 생성
RUN mkdir -p /root/data

# 포트 8080 노출
EXPOSE 8080

# 실행 명령
CMD ["./main"]