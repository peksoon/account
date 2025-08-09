# 빌드 스테이지
FROM golang:1.19-alpine AS builder

# 작업 디렉토리 설정
WORKDIR /app

# 빌드에 필요한 패키지 설치
RUN apk add --no-cache gcc musl-dev
RUN apk add --no-cache tzdata
ENV TZ=Asia/Seoul

# 의존성 파일 복사 및 다운로드
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 애플리케이션 빌드
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# 실행 스테이지
FROM alpine:latest  

# 필요한 CA 인증서 및 SQLite 라이브러리 설치
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# 빌드 스테이지에서 빌드된 실행 파일 복사
COPY --from=builder /app/main .

# 애플리케이션 실행
CMD ["./main"]
