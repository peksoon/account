<template>
    <div v-if="loading" class="text-center py-8">
        <div class="inline-flex items-center">
            <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-green-500"></div>
            <span class="ml-3 text-gray-600">로딩 중...</span>
        </div>
    </div>
    <div v-else-if="categories.length === 0" class="text-center py-8 text-gray-500">
        등록된 {{ type === 'out' ? '지출' : '수입' }} 카테고리가 없습니다.
    </div>
    <div v-else class="space-y-3">
        <div v-for="category in categories" :key="category.id"
            class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
            <div class="flex-1">
                <div class="flex items-center">
                    <span class="font-semibold text-gray-900">{{ category.name }}</span>
                    <span class="ml-2 px-2 py-1 text-xs rounded-full"
                        :class="category.type === 'out' ? 'bg-red-100 text-red-600' : 'bg-blue-100 text-blue-600'">
                        {{ category.type === 'out' ? '지출' : '수입' }}
                    </span>
                    <!-- 지출 유형 표시 (지출 카테고리만) -->
                    <span v-if="category.type === 'out'" class="ml-2 px-2 py-1 text-xs rounded-full"
                        :class="category.expense_type === 'fixed' ? 'bg-blue-100 text-blue-600' : 'bg-green-100 text-green-600'">
                        {{ category.expense_type === 'fixed' ? '📌 고정' : '💳 변동' }}
                    </span>
                </div>
                <div class="text-xs text-gray-400 mt-1">
                    ID: {{ category.id }} | 생성일: {{ formatDate(category.created_at) }}
                </div>
            </div>
            <div class="flex space-x-2">
                <el-button size="small" @click="$emit('edit', category)" type="primary">
                    수정
                </el-button>
                <el-button size="small" @click="$emit('delete', category)" type="danger">
                    삭제
                </el-button>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    props: ['categories', 'loading', 'type'],
    emits: ['edit', 'delete'],
    setup() {
        const formatDate = (dateString) => {
            if (!dateString) return '';
            return new Date(dateString).toLocaleString('ko-KR', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit'
            });
        };

        return { formatDate };
    }
};
</script>
