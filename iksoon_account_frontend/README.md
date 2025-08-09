# 💰 스마트 가계부 프론트엔드

현대적이고 직관적인 가계부 관리 웹 애플리케이션입니다.

## ✨ 주요 기능

- 📅 **모던한 달력 인터페이스**: FullCalendar 기반의 아름다운 달력
- 💸 **수입/지출 관리**: 직관적인 카테고리별 거래 관리
- 📊 **실시간 통계**: 월별 수입, 지출, 잔액 현황
- 📱 **반응형 디자인**: 모든 디바이스에서 완벽한 사용 경험
- 🎨 **모던 UI/UX**: Element Plus와 Tailwind CSS 기반의 세련된 디자인
- ⚡ **빠른 성능**: Vue 3 Composition API와 최적화된 컴포넌트

## 🛠 기술 스택

### Frontend

- **Vue 3** - 프레임워크
- **Element Plus** - UI 컴포넌트 라이브러리
- **FullCalendar** - 달력 컴포넌트
- **Tailwind CSS** - 유틸리티 우선 CSS 프레임워크
- **Pinia** - 상태 관리
- **Axios** - HTTP 클라이언트
- **Lucide Vue** - 아이콘 라이브러리

### Development

- **Vue CLI** - 빌드 도구
- **ESLint** - 코드 품질 관리
- **Babel** - JavaScript 컴파일러

## 🚀 설치 및 실행

### 프로젝트 설정

```bash
npm install
```

### 개발 서버 실행 (핫 리로드)

```bash
npm run serve
```

**주의**: 프론트엔드가 정상적으로 작동하려면 백엔드 서버가 `http://localhost:8080`에서 실행되고 있어야 합니다.

### 프로덕션 빌드

```bash
npm run build
```

### 코드 린팅 및 수정

```bash
npm run lint
```

### 문제 해결

**의존성 오류 발생 시:**
```bash
rm -rf node_modules package-lock.json
npm install
```

**ESLint 오류 발생 시:**
```bash
npm run lint
```

## 📦 프로젝트 구조

```
src/
├── components/          # Vue 컴포넌트
│   ├── AccountCalendar.vue    # 메인 달력 컴포넌트
│   ├── AddPopup.vue          # 데이터 추가 팝업
│   └── DetailPopup.vue       # 상세보기/편집 팝업
├── stores/             # Pinia 스토어
│   ├── accountStore.js       # 계정 데이터 관리
│   ├── calendarStore.js      # 달력 상태 관리
│   └── popupStore.js         # 팝업 상태 관리
├── assets/             # 정적 자원
│   └── styles.css           # 전역 스타일
├── App.vue             # 루트 컴포넌트
├── main.js             # 애플리케이션 진입점
└── utils.js            # 유틸리티 함수
```

## 🎨 디자인 시스템

### 색상 팔레트

- **Primary**: Blue (수입)
- **Danger**: Red (지출)
- **Success**: Green (수입)
- **Gray**: 중성 색상

### 컴포넌트

- **카드**: 그림자와 둥근 모서리를 가진 컨테이너
- **버튼**: 그라데이션과 호버 효과
- **모달**: 부드러운 애니메이션과 백드롭
- **폼**: 일관된 입력 필드 스타일

## 📱 반응형 브레이크포인트

- **Mobile**: < 768px
- **Tablet**: 768px - 1024px
- **Desktop**: > 1024px

## 🔧 환경 변수

```env
VUE_APP_BACKEND_API_BASE_URL=http://localhost:8080
```

## 🌟 주요 개선사항 (v1.0.0)

### UI/UX 개선

- ✅ Toast UI Calendar → FullCalendar 마이그레이션
- ✅ Element Plus 통합으로 일관된 디자인
- ✅ Tailwind CSS로 모던한 스타일링
- ✅ 반응형 디자인 완전 지원
- ✅ 부드러운 애니메이션 및 전환 효과

### 기능 개선

- ✅ 실시간 월별 통계 대시보드
- ✅ 향상된 폼 검증 및 사용자 피드백
- ✅ 카테고리 자동 제안 기능
- ✅ 모바일 친화적 터치 인터페이스
- ✅ 접근성 개선

### 성능 최적화

- ✅ Vue 3 Composition API 활용
- ✅ 컴포넌트 지연 로딩
- ✅ 효율적인 상태 관리
- ✅ 최적화된 빌드 설정

## 🔗 백엔드 연동

이 프론트엔드는 Go 언어로 작성된 REST API 백엔드와 연동됩니다.

### API 엔드포인트

- `GET /month-out-account` - 월별 지출 조회
- `GET /month-in-account` - 월별 수입 조회
- `POST /out-account/insert` - 지출 데이터 추가
- `POST /in-account/insert` - 수입 데이터 추가
- `PUT /out-account/update` - 지출 데이터 수정
- `PUT /in-account/update` - 수입 데이터 수정
- `DELETE /out-account/delete` - 지출 데이터 삭제
- `DELETE /in-account/delete` - 수입 데이터 삭제

## 🚀 배포

### Nginx 설정 예시

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        root /path/to/dist;
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 📄 라이선스

이 프로젝트는 MIT 라이선스 하에 배포됩니다.

## 🤝 기여하기

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

**Made with ❤️ using Vue 3 + Element Plus**
