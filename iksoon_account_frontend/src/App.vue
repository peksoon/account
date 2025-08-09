<template>
  <div id="app" class="min-h-screen">
    <!-- 메인 컨텐츠 -->
    <main class="relative">
      <AccountCalendar />
    </main>
    
    <!-- 푸터 -->
    <footer class="bg-white border-t border-gray-200 py-6 mt-12">
      <div class="container-responsive">
        <div class="flex flex-col md:flex-row justify-between items-center">
          <div class="text-gray-600 text-sm mb-4 md:mb-0">
            <p>&copy; 2025 상윤 가족 가계부. Made with 익순이</p>
          </div>
          <div class="flex items-center space-x-4 text-sm text-gray-500">
            <span>v1.0.0</span>
            <span>•</span>
            <span>Vue 3 + Element Plus</span>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
import AccountCalendar from "./components/AccountCalendar.vue";

export default {
  name: "App",
  components: {
    AccountCalendar,
  },
};
</script>

<style>
/* 전역 앱 스타일 */
#app {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  position: relative;
}

/* 배경 패턴 */
#app::before {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: 
    radial-gradient(circle at 20% 50%, rgba(120, 119, 198, 0.03) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(255, 119, 198, 0.03) 0%, transparent 50%),
    radial-gradient(circle at 40% 80%, rgba(120, 219, 255, 0.03) 0%, transparent 50%);
  pointer-events: none;
  z-index: -1;
}

/* 부드러운 스크롤 */
html {
  scroll-behavior: smooth;
}

/* 포커스 아웃라인 커스텀 */
*:focus {
  outline: none;
}

*:focus-visible {
  outline: 2px solid #3b82f6;
  outline-offset: 2px;
}

/* 선택 텍스트 색상 */
::selection {
  background-color: #3b82f6;
  color: white;
}

/* 애니메이션 성능 최적화 */
* {
  transform: translateZ(0);
  backface-visibility: hidden;
}

/* 로딩 상태 */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

/* 오류 상태 */
.error-message {
  @apply bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg;
}

/* 성공 상태 */
.success-message {
  @apply bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg;
}

/* 반응형 이미지 */
img {
  max-width: 100%;
  height: auto;
}

/* 접근성 개선 */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

/* 다크모드 대응 (추후 확장 가능) */
@media (prefers-color-scheme: dark) {
  /* 추후 다크모드 지원 시 스타일 추가 */
}

/* 고해상도 디스플레이 최적화 */
@media (-webkit-min-device-pixel-ratio: 2), (min-resolution: 192dpi) {
  /* 레티나 디스플레이용 최적화 */
}

/* 모바일 터치 최적화 */
@media (hover: none) and (pointer: coarse) {
  /* 터치 디바이스용 최적화 */
  .hover\\:scale-105:hover {
    transform: none;
  }
  
  .hover\\:shadow-lg:hover {
    box-shadow: none;
  }
}

/* 프린트 스타일 */
@media print {
  #app {
    background: white !important;
  }
  
  .no-print {
    display: none !important;
  }
  
  .card {
    box-shadow: none !important;
    border: 1px solid #ccc !important;
  }
}
</style>
