<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="handleOverlayClick">
    <div 
      :class="modalClasses"
      @click.stop
    >
      <!-- 헤더 -->
      <div :class="headerClasses">
        <div class="flex justify-between items-center">
          <h2 class="text-xl font-bold">{{ title }}</h2>
          <button 
            @click="$emit('close')" 
            class="text-white hover:text-gray-200 transition-colors"
            aria-label="닫기"
          >
            <span class="text-2xl">×</span>
          </button>
        </div>
      </div>

      <!-- 컨텐츠 -->
      <div :class="contentClasses">
        <slot />
      </div>

      <!-- 푸터 (선택사항) -->
      <div v-if="$slots.footer" class="px-6 py-4 bg-gray-50 border-t">
        <slot name="footer" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

/**
 * 공통 모달 컴포넌트
 * 일관된 모달 레이아웃과 스타일을 제공합니다.
 */

const props = defineProps({
  /**
   * 모달 제목
   */
  title: {
    type: String,
    required: true
  },
  /**
   * 모달 크기
   * @values 'sm', 'md', 'lg', 'xl', '2xl', '3xl'
   */
  size: {
    type: String,
    default: '2xl',
    validator: (value) => ['sm', 'md', 'lg', 'xl', '2xl', '3xl'].includes(value)
  },
  /**
   * 헤더 색상 테마
   * @values 'blue', 'green', 'red', 'yellow', 'purple', 'gray'
   */
  theme: {
    type: String,
    default: 'blue',
    validator: (value) => ['blue', 'green', 'red', 'yellow', 'purple', 'gray'].includes(value)
  },
  /**
   * 오버레이 클릭으로 닫기 허용 여부
   */
  closeOnOverlay: {
    type: Boolean,
    default: true
  }
});

const emit = defineEmits(['close']);

/**
 * 모달 크기별 클래스 매핑
 */
const sizeClasses = {
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-xl',
  '2xl': 'max-w-2xl',
  '3xl': 'max-w-3xl'
};

/**
 * 테마별 헤더 색상 클래스 매핑
 */
const themeClasses = {
  blue: 'bg-gradient-to-r from-blue-500 to-blue-600',
  green: 'bg-gradient-to-r from-green-500 to-green-600',
  red: 'bg-gradient-to-r from-red-500 to-red-600',
  yellow: 'bg-gradient-to-r from-yellow-500 to-yellow-600',
  purple: 'bg-gradient-to-r from-purple-500 to-purple-600',
  gray: 'bg-gradient-to-r from-gray-500 to-gray-600'
};

/**
 * 계산된 모달 클래스
 */
const modalClasses = computed(() => [
  'bg-white rounded-xl shadow-2xl w-full max-h-[90vh] overflow-hidden',
  sizeClasses[props.size]
]);

/**
 * 계산된 헤더 클래스
 */
const headerClasses = computed(() => [
  'text-white p-6',
  themeClasses[props.theme]
]);

/**
 * 계산된 컨텐츠 클래스
 */
const contentClasses = computed(() => [
  'p-6',
  'max-h-[calc(90vh-120px)] overflow-y-auto'
]);

/**
 * 오버레이 클릭 핸들러
 */
const handleOverlayClick = () => {
  if (props.closeOnOverlay) {
    emit('close');
  }
};
</script>