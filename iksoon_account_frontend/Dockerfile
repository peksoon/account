# Frontend Dockerfile
FROM node:18-alpine AS builder

# 작업 디렉토리 설정
WORKDIR /app

# package.json과 package-lock.json 복사
COPY package*.json ./

# 의존성 설치
RUN npm ci --only=production

# 소스 코드 복사
COPY . .

# Vue.js 애플리케이션 빌드
RUN npm run build

# 실행용 Nginx 이미지
FROM nginx:alpine

# Nginx 설정 파일 복사
COPY nginx.conf /etc/nginx/nginx.conf

# 빌드된 Vue.js 앱을 Nginx 서빙 디렉토리로 복사
COPY --from=builder /app/dist /usr/share/nginx/html

# 포트 3000 노출
EXPOSE 3000

# Nginx 실행
CMD ["nginx", "-g", "daemon off;"]