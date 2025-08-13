<template>
  <div class="text-center py-8">
    <div class="inline-flex items-center">
      <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-500"></div>
      <span class="ml-3 text-gray-600">{{ message }}</span>
    </div>
  </div>
</template>

<script setup>
/**
 * 로딩 상태 표시 컴포넌트
 */

defineProps({
  /**
   * 로딩 메시지
   */
  message: {
    type: String,
    default: '로딩 중...'
  }
});
</script>