<template>
    <div class="category-step">
        <div class="step-header">
            <h2 class="step-title">카테고리를 선택하세요</h2>
            <p class="step-description">{{ modelValue.type === 'out' ? '지출' : '수입' }} 카테고리를 선택해주세요</p>
        </div>

        <div class="category-selection-container">
            <!-- 카테고리 그리드 -->
            <div class="category-grid" v-if="availableCategories.length > 0">
                <button v-for="category in availableCategories" :key="category.id" @click="selectCategory(category)"
                    class="category-btn" :class="{ 'selected': selectedCategoryId === category.id }">
                    <div class="category-info">
                        <div class="category-name">{{ category.name }}</div>
                        <div class="category-description" v-if="category.description">
                            {{ category.description }}
                        </div>
                    </div>
                </button>
            </div>

            <!-- 검색 가능한 선택기 (많은 카테고리가 있을 때) -->
            <div class="category-search" v-if="availableCategories.length > 12">
                <el-select ref="categorySelectRef" v-model="selectedCategoryId" placeholder="카테고리를 검색하세요" size="large"
                    class="category-select" :class="{ 'error': hasError }" filterable @change="handleCategoryChange">
                    <el-option v-for="category in availableCategories" :key="category.id" :label="category.name"
                        :value="category.id" class="category-option">
                        <span>{{ category.name }}</span>
                    </el-option>
                </el-select>
            </div>

            <div v-if="hasError" class="error-message">
                {{ errorMessage }}
            </div>
        </div>

        <!-- 카테고리가 없을 때 -->
        <div v-if="availableCategories.length === 0" class="no-categories">
            <div class="no-categories-icon">
                <FolderPlus class="w-12 h-12 text-gray-400" />
            </div>
            <p class="no-categories-message">
                {{ modelValue.type === 'out' ? '지출' : '수입' }} 카테고리가 없습니다.
            </p>
            <p class="text-sm text-gray-500 mt-2">
                💡 달력 화면의 <strong>⚙️ 관리</strong> 버튼에서 카테고리를 추가하세요
            </p>
        </div>
    </div>
</template>

<script>
import { ref, computed, watch, nextTick, onMounted } from 'vue';
import {
    FolderPlus
} from 'lucide-vue-next';
import { useCategoryStore } from '../../stores/categoryStore';

export default {
    name: 'CategoryStep',
    components: {
        FolderPlus
    },
    props: {
        modelValue: {
            type: Object,
            required: true
        },
        errors: {
            type: Object,
            default: () => ({})
        }
    },
    emits: ['update:modelValue', 'next', 'auto-advance', 'validate'],
    setup(props, { emit }) {
        const categoryStore = useCategoryStore();
        const categorySelectRef = ref(null);
        const selectedCategoryId = ref(null);

        // 에러 상태
        const hasError = computed(() => !!props.errors.category_id);
        const errorMessage = computed(() => props.errors.category_id || '');

        // 사용 가능한 카테고리
        const availableCategories = computed(() => {
            return categoryStore.categories.filter(cat => cat.type === props.modelValue.type);
        });


        // 이벤트 핸들러들
        const selectCategory = (category) => {
            selectedCategoryId.value = category.id;
            updateModelValue(category.id);
            emit('validate', 'category_id', true, '');

            // 선택 즉시 자동 진행
            setTimeout(() => {
                emit('auto-advance', 100);
            }, 50);
        };

        const handleCategoryChange = (value) => {
            if (value) {
                selectedCategoryId.value = value;
                updateModelValue(value);
                emit('validate', 'category_id', true, '');

                // 드롭다운 선택 시 자동 진행
                setTimeout(() => {
                    emit('auto-advance', 200);
                }, 100);
            } else {
                emit('validate', 'category_id', false, '카테고리를 선택해주세요');
            }
        };

        const updateModelValue = (value) => {
            const updated = {
                ...props.modelValue,
                category_id: value,
                keyword_name: '' // 카테고리 변경 시 키워드 초기화
            };
            emit('update:modelValue', updated);
        };

        // 초기값 설정
        watch(() => props.modelValue.category_id, (newValue) => {
            if (newValue !== selectedCategoryId.value) {
                selectedCategoryId.value = newValue;
            }
        }, { immediate: true });

        // 타입 변경 시 선택 초기화
        watch(() => props.modelValue.type, () => {
            selectedCategoryId.value = null;
            updateModelValue(null);
        });

        // 컴포넌트 마운트 시 포커스
        onMounted(() => {
            nextTick(() => {
                if (availableCategories.value.length <= 12) {
                    // 그리드 방식일 때 첫 번째 버튼에 포커스
                    const firstBtn = document.querySelector('.category-btn');
                    if (firstBtn) {
                        firstBtn.focus();
                    }
                } else if (categorySelectRef.value) {
                    // 드롭다운 방식일 때 select에 포커스
                    categorySelectRef.value.focus();
                }
            });
        });

        return {
            categorySelectRef,
            selectedCategoryId,
            hasError,
            errorMessage,
            availableCategories,
            selectCategory,
            handleCategoryChange
        };
    }
}
</script>

<style scoped>
.category-step {
    @apply max-w-lg mx-auto;
}

.step-header {
    @apply text-center mb-8;
}

.step-title {
    @apply text-2xl font-bold text-gray-900 mb-2;
}

.step-description {
    @apply text-gray-600;
}

.category-selection-container {
    @apply space-y-6;
}

.category-grid {
    @apply grid gap-2;
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
}

@media (max-width: 640px) {
    .category-grid {
        grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
        @apply gap-2;
    }
}

.category-btn {
    @apply flex items-center justify-center p-3 bg-white border-2 border-gray-200 rounded-lg transition-all duration-200 text-center min-h-[60px];
}

.category-btn:hover {
    @apply border-gray-300 bg-gray-50 transform scale-105;
}

.category-btn:active {
    @apply transform scale-95;
}

.category-btn.selected {
    @apply border-blue-500 bg-blue-50 ring-2 ring-blue-100;
}

.category-info {
    @apply w-full;
}

.category-name {
    @apply font-semibold text-gray-900 text-sm;
}

.category-btn.selected .category-name {
    @apply text-blue-700;
}

.category-description {
    @apply text-xs text-gray-500 mt-1;
}

.category-search {
    @apply mb-4;
}

.category-select {
    @apply w-full;
}

:deep(.category-select .el-input__inner) {
    @apply h-14 text-lg border-2;
}

:deep(.category-select.error .el-input__inner) {
    @apply border-red-500;
}

:deep(.category-select .el-input__inner:focus) {
    @apply border-blue-500 ring-2 ring-blue-100;
}

.category-option {
    @apply py-2;
}

.error-message {
    @apply text-sm text-red-600 text-center;
}

.category-management {
    @apply text-center;
}

.category-manage-btn {
    @apply text-blue-600;
}

.category-manage-btn:hover {
    @apply text-blue-700;
}

.no-categories {
    @apply text-center py-8;
}

.no-categories-icon {
    @apply mb-4;
}

.no-categories-message {
    @apply text-gray-600 mb-4;
}

/* 모바일 최적화 */
@media (max-width: 768px) {
    .category-btn {
        @apply p-3;
        min-height: 50px;
        /* iOS 권장 터치 영역 */
    }

    .category-name {
        @apply text-sm;
    }

    .category-description {
        @apply text-xs;
    }

    :deep(.category-select .el-input__inner) {
        @apply h-12 text-base;
        font-size: 16px;
        /* iOS zoom 방지 */
    }
}

/* 접근성 */
.category-btn:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

.category-select:focus-within {
    @apply ring-2 ring-blue-500 ring-offset-2;
}

.category-manage-btn:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

/* 애니메이션 */
.category-grid {
    animation: fadeInUp 0.3s ease-out;
}

@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(10px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}


/* 선택된 카테고리 강조 애니메이션 */
.category-btn.selected {
    animation: selectPulse 0.3s ease-out;
}

@keyframes selectPulse {
    0% {
        transform: scale(1);
    }

    50% {
        transform: scale(1.05);
    }

    100% {
        transform: scale(1);
    }
}
</style>