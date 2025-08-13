<template>
  <div class="text-center py-8">
    <div class="text-gray-400 mb-4">
      <component :is="icon" class="w-16 h-16 mx-auto" />
    </div>
    <p class="text-gray-500 text-lg">{{ message }}</p>
    <p v-if="description" class="text-gray-400 text-sm mt-2">{{ description }}</p>
    
    <!-- 액션 버튼 슬롯 -->
    <div v-if="$slots.action" class="mt-6">
      <slot name="action" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { Box, User, PriceTag, CreditCard } from '@element-plus/icons-vue';

/**
 * 빈 상태 표시 컴포넌트
 */

const props = defineProps({
  /**
   * 빈 상태 메시지
   */
  message: {
    type: String,
    default: '데이터가 없습니다.'
  },
  /**
   * 부가 설명
   */
  description: {
    type: String,
    default: ''
  },
  /**
   * 표시할 아이콘
   * @values 'inbox', 'users', 'tags', 'credit-card'
   */
  iconType: {
    type: String,
    default: 'inbox',
    validator: (value) => ['inbox', 'users', 'tags', 'credit-card'].includes(value)
  }
});

/**
 * 아이콘 타입별 컴포넌트 매핑
 */
const iconComponents = {
  inbox: Box,
  users: User,
  tags: PriceTag,
  'credit-card': CreditCard
};

/**
 * 계산된 아이콘 컴포넌트
 */
const icon = computed(() => iconComponents[props.iconType]);
</script>