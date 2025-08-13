<template>
  <div class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
    <div class="flex-1">
      <div class="flex items-center">
        <span class="font-semibold text-gray-900">{{ user.name }}</span>
        <span v-if="user.email" class="ml-2 text-sm text-gray-500">({{ user.email }})</span>
        <el-tag v-if="!user.is_active" type="warning" size="small" class="ml-2">
          비활성
        </el-tag>
      </div>
      <div class="text-xs text-gray-400 mt-1">
        ID: {{ user.id }} | 생성일: {{ formatDate(user.created_at) }}
        <span v-if="user.updated_at !== user.created_at">
          | 수정일: {{ formatDate(user.updated_at) }}
        </span>
      </div>
    </div>
    <div class="flex space-x-2">
      <el-button 
        size="small" 
        @click="$emit('edit', user)" 
        type="primary" 
        :disabled="loading"
        :icon="Edit"
      >
        수정
      </el-button>
      <el-button 
        size="small" 
        @click="$emit('delete', user)" 
        type="danger"
        :disabled="loading"
        :icon="Delete"
      >
        삭제
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { Edit, Delete } from '@element-plus/icons-vue';

/**
 * 사용자 목록 아이템 컴포넌트
 * 개별 사용자 정보를 표시하고 수정/삭제 기능을 제공합니다.
 */

defineProps({
  /**
   * 사용자 객체
   */
  user: {
    type: Object,
    required: true
  },
  /**
   * 로딩 상태
   */
  loading: {
    type: Boolean,
    default: false
  }
});

defineEmits(['edit', 'delete']);

/**
 * 날짜를 한국어 형식으로 포맷팅
 * @param {string} dateString - 날짜 문자열
 * @returns {string} 포맷된 날짜
 */
const formatDate = (dateString) => {
  if (!dateString) return '-';
  
  try {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('ko-KR', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    }).format(date);
  } catch (error) {
    console.error('날짜 포맷 오류:', error);
    return dateString;
  }
};
</script>