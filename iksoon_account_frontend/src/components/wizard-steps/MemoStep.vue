<template>
    <div class="memo-step">
        <div class="step-header">
            <h2 class="step-title">메모를 추가하세요</h2>
            <p class="step-description">추가 설명이나 메모를 입력하세요 (선택사항)</p>
        </div>

        <div class="memo-input-container">
            <div class="memo-input-wrapper">
                <el-input
                    ref="memoInputRef"
                    v-model="memoText"
                    type="textarea"
                    :rows="4"
                    placeholder="예: 점심 회식, 프로젝트 관련 지출, 생일 선물 등..."
                    resize="none"
                    size="large"
                    class="memo-input"
                    maxlength="200"
                    show-word-limit
                    @input="handleMemoInput"
                    @keydown="handleKeydown"
                    @blur="handleBlur"
                />
                
                <div class="memo-hints">
                    <div class="memo-hint-item">
                        <FileText class="w-4 h-4 mr-2 text-gray-400" />
                        <span class="text-sm text-gray-500">Enter 또는 다음 버튼을 눌러 계속하세요</span>
                    </div>
                </div>
            </div>

            <!-- 빠른 메모 템플릿 -->
            <div class="quick-memos" v-if="!memoText.trim()">
                <p class="quick-memos-label">빠른 선택</p>
                <div class="quick-memos-grid">
                    <button
                        v-for="template in memoTemplates"
                        :key="template"
                        @click="selectMemoTemplate(template)"
                        class="quick-memo-btn"
                    >
                        <MessageSquare class="w-4 h-4 mr-2" />
                        {{ template }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, watch, nextTick, onMounted } from 'vue';
import { FileText, MessageSquare } from 'lucide-vue-next';

export default {
    name: 'MemoStep',
    components: {
        FileText,
        MessageSquare
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
        const memoInputRef = ref(null);
        const memoText = ref('');

        // 메모 템플릿 (거래 유형에 따라 다름)
        const memoTemplates = computed(() => {
            if (props.modelValue.type === 'out') {
                return [
                    '회사 점심',
                    '가족 외식',
                    '업무 관련',
                    '생필품 구매',
                    '교통비',
                    '의료비'
                ];
            } else {
                return [
                    '월급',
                    '보너스',
                    '용돈',
                    '부업 수입',
                    '이자 수입',
                    '기타 수입'
                ];
            }
        });

        // 이벤트 핸들러들
        const handleMemoInput = (value) => {
            memoText.value = value;
            updateModelValue(value);
        };

        const handleKeydown = (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                emit('next');
            }
        };

        const handleBlur = () => {
            // 블러 시에는 자동 진행하지 않음
        };

        const selectMemoTemplate = (template) => {
            memoText.value = template;
            updateModelValue(template);
        };

        const updateModelValue = (value) => {
            const updated = { ...props.modelValue, memo: value };
            emit('update:modelValue', updated);
        };


        // 초기값 설정
        watch(() => props.modelValue.memo, (newValue) => {
            if (newValue !== memoText.value) {
                memoText.value = newValue || '';
            }
        }, { immediate: true });

        // 컴포넌트 마운트 시 포커스
        onMounted(() => {
            nextTick(() => {
                if (memoInputRef.value) {
                    memoInputRef.value.focus();
                }
            });
        });

        return {
            memoInputRef,
            memoText,
            memoTemplates,
            handleMemoInput,
            handleKeydown,
            handleBlur,
            selectMemoTemplate
        };
    }
}
</script>

<style scoped>
.memo-step {
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

.memo-input-container {
    @apply space-y-6;
}

.memo-input-wrapper {
    @apply relative;
}

.memo-input {
    @apply w-full;
}

:deep(.memo-input .el-textarea__inner) {
    @apply text-base border-2 rounded-lg p-4 leading-relaxed;
    min-height: 120px;
}

:deep(.memo-input .el-textarea__inner:focus) {
    @apply border-blue-500 ring-2 ring-blue-100;
}

:deep(.memo-input .el-input__count) {
    @apply text-xs text-gray-400;
}

.memo-hints {
    @apply mt-3 space-y-2;
}

.memo-hint-item {
    @apply flex items-center;
}

.quick-memos {
    @apply text-center;
}

.quick-memos-label {
    @apply text-sm text-gray-600 mb-3;
}

.quick-memos-grid {
    @apply grid grid-cols-2 gap-2;
}

.quick-memo-btn {
    @apply flex items-center justify-center px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg transition-colors duration-200;
}

.quick-memo-btn:hover {
    @apply bg-gray-50 border-gray-400;
}

.quick-memo-btn:active {
    @apply bg-gray-100;
}

/* 모바일 최적화 */
@media (max-width: 768px) {
    :deep(.memo-input .el-textarea__inner) {
        @apply text-base p-4;
        font-size: 16px; /* iOS zoom 방지 */
        min-height: 100px;
    }
    
    .quick-memos-grid {
        @apply grid-cols-1 gap-3;
    }
    
    .quick-memo-btn {
        @apply px-4 py-3 text-base justify-start;
        min-height: 44px; /* iOS 권장 터치 영역 */
    }
}

/* 접근성 */
.memo-input:focus-within {
    @apply ring-2 ring-blue-500 ring-offset-2;
}

.quick-memo-btn:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

/* 애니메이션 */
.quick-memos {
    animation: fadeInUp 0.3s ease-out;
}

.memo-hints {
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

/* 텍스트에어리어 플레이스홀더 스타일 */
:deep(.memo-input .el-textarea__inner::placeholder) {
    @apply text-gray-400 italic;
}

/* 워드 카운터 위치 조정 */
:deep(.memo-input .el-input__count) {
    @apply bottom-2 right-3;
}
</style>