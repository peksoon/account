# 1. Build 단계
FROM node:lts AS build

# 작업 디렉토리 설정
WORKDIR /app

# package.json과 lock 파일 복사
COPY package*.json ./

# 의존성 설치
RUN npm install

# 소스 코드 복사
COPY . .

# Vue.js 애플리케이션 빌드
RUN npm run build

# 2. Production 단계
FROM nginx:alpine

# 빌드된 파일을 Nginx의 기본 HTML 디렉토리로 복사
COPY --from=build /app/dist /usr/share/nginx/html

# 컨테이너가 실행될 때 Nginx 시작
CMD ["nginx", "-g", "daemon off;"]
