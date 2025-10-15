<template>
    <div class="keyword-step">
        <div class="step-header">
            <h2 class="step-title">키워드를 추가하세요</h2>
            <p class="step-description">
                {{ modelValue.type === 'out' ? '상품명, 식당명' : '회사명, 수입원' }} 등을 입력해주세요<br>
                <span class="text-xs text-gray-500">입력 시 기존 키워드가 자동완성됩니다</span>
            </p>
        </div>

        <div class="keyword-input-container">
            <!-- 키워드 입력 -->
            <div class="keyword-input-wrapper">
                <el-autocomplete ref="keywordInputRef" v-model="keywordInput" :fetch-suggestions="fetchSuggestions"
                    placeholder="키워드를 입력하세요" size="large" class="keyword-input" :class="{ 'error': hasError }" clearable
                    @select="handleKeywordSelect" @input="handleKeywordInput" @keydown="handleKeydown"
                    @clear="handleClear" @blur="handleBlur">
                    <template #default="{ item }">
                        <div class="flex items-center justify-between w-full">
                            <div class="flex items-center">
                                <Tag class="w-4 h-4 text-gray-500 mr-2" />
                                <span class="text-gray-900">{{ item.name }}</span>
                            </div>
                            <div class="flex items-center text-xs text-gray-500">
                                <span>{{ item.usage_count }}회 사용</span>
                            </div>
                        </div>
                    </template>

                    <template #prefix>
                        <Tag class="w-4 h-4 text-gray-500" />
                    </template>
                </el-autocomplete>

                <div v-if="hasError" class="error-message">
                    {{ errorMessage }}
                </div>

                <!-- 입력된 키워드 표시 -->
                <div v-if="currentKeyword" class="current-keyword">
                    <div class="keyword-tag">
                        <Tag class="w-4 h-4 mr-2" />
                        {{ currentKeyword }}
                        <button @click="removeKeyword" class="keyword-remove">
                            <X class="w-4 h-4" />
                        </button>
                    </div>
                </div>
            </div>

            <!-- 최근 키워드 추천 -->
            <div class="recent-keywords" v-if="recentKeywords.length > 0 && !currentKeyword">
                <p class="recent-keywords-label">최근 사용한 키워드 (상위 10개)</p>
                <div class="recent-keywords-grid">
                    <button v-for="keyword in recentKeywords" :key="keyword.id" @click="selectKeyword(keyword.name)"
                        class="recent-keyword-btn">
                        <Tag class="w-4 h-4 mr-2" />
                        {{ keyword.name }}
                        <span class="keyword-count">{{ keyword.usage_count }}</span>
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, watch, nextTick, onMounted } from 'vue';
import { Tag, X } from 'lucide-vue-next';
import { useKeywordStore } from '../../stores/keywordStore';

export default {
    name: 'KeywordStep',
    components: {
        Tag,
        X
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
        const keywordStore = useKeywordStore();
        const keywordInputRef = ref(null);
        const keywordInput = ref('');
        const autoAdvanceTimer = ref(null);

        // 에러 상태
        const hasError = computed(() => !!props.errors.keyword_name);
        const errorMessage = computed(() => props.errors.keyword_name || '');

        // 현재 키워드
        const currentKeyword = computed(() => props.modelValue.keyword_name);

        // 카테고리별 최근 키워드 (최대 10개, 사용 빈도 순)
        const recentKeywords = computed(() => {
            if (!props.modelValue.category_id) return [];
            return keywordStore.getKeywordsByCategory(props.modelValue.category_id)
                .sort((a, b) => (b.usage_count || 0) - (a.usage_count || 0))
                .slice(0, 10); // 최대 10개
        });

        // 키워드 추천 함수
        const fetchSuggestions = async (queryString, callback) => {
            if (!props.modelValue.category_id) {
                callback([]);
                return;
            }

            try {
                // 백엔드 API에서 키워드 추천을 가져옴
                const suggestions = await keywordStore.getKeywordSuggestions(
                    props.modelValue.category_id,
                    (queryString || '').trim(),
                    20 // 자동완성에서는 더 많은 항목 제공
                );
                callback(suggestions);
            } catch (error) {
                console.error('키워드 추천 조회 실패:', error);
                // 실패 시 로컬 키워드에서 필터링
                const keywords = keywordStore.getKeywordsByCategory(props.modelValue.category_id);
                const query = (queryString || '').trim();

                if (!query) {
                    // 빈 쿼리일 때는 최근 사용 순으로 정렬하여 10개만 표시
                    const suggestions = keywords
                        .sort((a, b) => (b.usage_count || 0) - (a.usage_count || 0))
                        .slice(0, 10);
                    callback(suggestions);
                } else {
                    // 쿼리가 있을 때는 전체에서 필터링 (최대 20개)
                    const suggestions = keywords
                        .filter(keyword =>
                            keyword.name && keyword.name.toLowerCase().includes(query.toLowerCase())
                        )
                        .sort((a, b) => (b.usage_count || 0) - (a.usage_count || 0))
                        .slice(0, 20);
                    callback(suggestions);
                }
            }
        };

        // 이벤트 핸들러들
        const handleKeywordInput = (value) => {
            keywordInput.value = value || '';

            // 입력 중일 때는 자동 진행 타이머 취소
            clearAutoAdvanceTimer();

            if (value && value.trim()) {
                // 유효성 검사 통과
                emit('validate', 'keyword_name', true, '');
            } else {
                // 필수 필드이므로 에러
                emit('validate', 'keyword_name', false, '키워드를 입력해주세요');
            }
        };

        const handleKeywordSelect = (item) => {
            selectKeyword(item.name);
        };

        const handleKeydown = (e) => {
            if (e.key === 'Enter' && keywordInput.value && keywordInput.value.trim()) {
                e.preventDefault();
                confirmKeyword();
            }
        };

        const handleBlur = () => {
            if (keywordInput.value && keywordInput.value.trim()) {
                // 입력 필드에서 포커스가 벗어날 때 자동으로 키워드 확정
                setTimeout(() => {
                    if (keywordInput.value && keywordInput.value.trim() && !currentKeyword.value) {
                        confirmKeyword();
                    }
                }, 100);
            }
        };

        const handleClear = () => {
            keywordInput.value = '';
            updateModelValue('');
            emit('validate', 'keyword_name', false, '키워드를 입력해주세요');
        };

        const selectKeyword = (keywordName) => {
            keywordInput.value = keywordName;
            confirmKeyword();
        };

        const confirmKeyword = () => {
            const keyword = keywordInput.value && keywordInput.value.trim();
            if (keyword) {
                updateModelValue(keyword);
                keywordInput.value = ''; // 입력 필드 클리어
                emit('validate', 'keyword_name', true, '');

                // 키워드 확정 후 자동 진행
                clearAutoAdvanceTimer();
                autoAdvanceTimer.value = setTimeout(() => {
                    emit('auto-advance', 300);
                }, 200);
            }
        };

        const removeKeyword = () => {
            updateModelValue('');
            emit('validate', 'keyword_name', false, '키워드를 입력해주세요');

            // 포커스를 입력 필드로 이동
            nextTick(() => {
                if (keywordInputRef.value) {
                    keywordInputRef.value.focus();
                }
            });
        };

        const updateModelValue = (value) => {
            const updated = { ...props.modelValue, keyword_name: value };
            emit('update:modelValue', updated);
        };

        const clearAutoAdvanceTimer = () => {
            if (autoAdvanceTimer.value) {
                clearTimeout(autoAdvanceTimer.value);
                autoAdvanceTimer.value = null;
            }
        };

        // 초기값 설정
        watch(() => props.modelValue.keyword_name, (newValue) => {
            if (!newValue && keywordInput.value) {
                keywordInput.value = '';
            }
        }, { immediate: true });

        // 카테고리 변경 시 키워드 초기화 및 해당 카테고리 키워드 로드
        watch(() => props.modelValue.category_id, async (newCategoryId) => {
            keywordInput.value = '';
            updateModelValue('');

            // 새 카테고리의 키워드 목록 로드
            if (newCategoryId) {
                try {
                    await keywordStore.fetchKeywordsByCategory(newCategoryId);
                } catch (error) {
                    console.error('카테고리 키워드 로드 실패:', error);
                }
            }
        });

        // 컴포넌트 마운트 시 포커스 및 키워드 로드
        onMounted(async () => {
            // 현재 카테고리의 키워드 로드
            if (props.modelValue.category_id) {
                try {
                    await keywordStore.fetchKeywordsByCategory(props.modelValue.category_id);
                } catch (error) {
                    console.error('초기 키워드 로드 실패:', error);
                }
            }

            nextTick(() => {
                if (keywordInputRef.value) {
                    keywordInputRef.value.focus();
                }
            });
        });

        return {
            keywordInputRef,
            keywordInput,
            hasError,
            errorMessage,
            currentKeyword,
            recentKeywords,
            fetchSuggestions,
            handleKeywordInput,
            handleKeywordSelect,
            handleKeydown,
            handleBlur,
            handleClear,
            selectKeyword,
            removeKeyword
        };
    }
}
</script>

<style scoped>
.keyword-step {
    @apply max-w-md mx-auto;
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

.keyword-input-container {
    @apply space-y-6;
}

.keyword-input-wrapper {
    @apply relative;
}

.keyword-input {
    @apply w-full;
}

:deep(.keyword-input .el-input__inner) {
    @apply h-14 text-lg border-2;
}

:deep(.keyword-input.error .el-input__inner) {
    @apply border-red-500;
}

:deep(.keyword-input .el-input__inner:focus) {
    @apply border-blue-500 ring-2 ring-blue-100;
}

.error-message {
    @apply mt-2 text-sm text-red-600;
}

.current-keyword {
    @apply mt-3;
}

.keyword-tag {
    @apply inline-flex items-center px-4 py-2 bg-blue-100 text-blue-800 rounded-lg border border-blue-200;
}

.keyword-remove {
    @apply ml-2 p-1 text-blue-600 hover:text-blue-800 transition-colors duration-200;
}

.keyword-remove:hover {
    @apply bg-blue-200 rounded;
}

.recent-keywords {
    @apply text-center;
}

.recent-keywords-label {
    @apply text-sm text-gray-600 mb-3;
}

.recent-keywords-grid {
    @apply grid grid-cols-2 gap-2;
}

.recent-keyword-btn {
    @apply flex items-center justify-between px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg transition-colors duration-200;
}

.recent-keyword-btn:hover {
    @apply bg-gray-50 border-gray-400;
}

.recent-keyword-btn:active {
    @apply bg-gray-100;
}

.keyword-count {
    @apply text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded;
}

.keyword-management {
    @apply text-center;
}

.keyword-manage-btn {
    @apply text-blue-600;
}

.keyword-manage-btn:hover {
    @apply text-blue-700;
}

/* 모바일 최적화 */
@media (max-width: 768px) {
    :deep(.keyword-input .el-input__inner) {
        @apply h-12 text-base;
        font-size: 16px;
        /* iOS zoom 방지 */
    }

    .recent-keywords-grid {
        @apply grid-cols-1 gap-3;
    }

    .recent-keyword-btn {
        @apply px-4 py-3 text-base justify-start;
        min-height: 44px;
        /* iOS 권장 터치 영역 */
    }

    .keyword-count {
        @apply ml-auto;
    }
}

/* 접근성 */
.keyword-input:focus-within {
    @apply ring-2 ring-blue-500 ring-offset-2;
}

.recent-keyword-btn:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

.keyword-remove:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-1;
}

.keyword-manage-btn:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

/* Element Plus 자동완성 스타일 */
:deep(.el-autocomplete-suggestion) {
    @apply max-h-60;
}

:deep(.el-autocomplete-suggestion__list) {
    @apply py-2;
}

:deep(.el-autocomplete-suggestion__item) {
    @apply py-3 px-4;
}

:deep(.el-autocomplete-suggestion__item:hover) {
    @apply bg-blue-50;
}

/* 애니메이션 */
.current-keyword {
    animation: fadeInUp 0.3s ease-out;
}

.recent-keywords {
    animation: fadeInUp 0.3s ease-out 0.1s both;
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

/* 키워드 태그 애니메이션 */
.keyword-tag {
    animation: slideInScale 0.3s ease-out;
}

@keyframes slideInScale {
    from {
        opacity: 0;
        transform: scale(0.8) translateY(-10px);
    }

    to {
        opacity: 1;
        transform: scale(1) translateY(0);
    }
}
</style>