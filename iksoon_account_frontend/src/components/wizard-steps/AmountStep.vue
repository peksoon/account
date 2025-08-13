<template>
    <div class="amount-step">
        <div class="step-header">
            <h2 class="step-title">금액을 입력하세요</h2>
            <p class="step-description">{{ modelValue.type === 'out' ? '지출' : '수입' }} 금액을 입력해주세요</p>
        </div>

        <div class="amount-input-container">
            <div class="amount-input-wrapper">
                <el-input
                    ref="amountInputRef"
                    v-model="displayAmount"
                    type="text"
                    placeholder="0"
                    size="large"
                    class="amount-input"
                    :class="{ 'error': hasError }"
                    @input="handleAmountInput"
                    @blur="handleAmountBlur"
                    @keydown="handleKeydown"
                    @paste="handlePaste"
                >
                    <template #suffix>
                        <span class="amount-suffix">원</span>
                    </template>
                </el-input>
                
                <div v-if="hasError" class="error-message">
                    {{ errorMessage }}
                </div>
                
                <div v-if="formattedAmount && !hasError" class="amount-preview">
                    {{ readableAmount }}
                </div>
            </div>
        </div>

        <div class="quick-amounts" v-if="!hasError">
            <p class="quick-amounts-label">빠른 입력</p>
            <div class="quick-amounts-grid">
                <button
                    v-for="amount in quickAmounts"
                    :key="amount"
                    @click="selectQuickAmount(amount)"
                    class="quick-amount-btn"
                >
                    +{{ formatDisplayAmount(amount) }}원
                </button>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, watch, nextTick, onMounted } from 'vue';

export default {
    name: 'AmountStep',
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
        const amountInputRef = ref(null);
        const displayAmount = ref('');
        const autoAdvanceTimer = ref(null);
        const isTyping = ref(false);

        // 빠른 입력 금액들
        const quickAmounts = computed(() => {
            if (props.modelValue.type === 'out') {
                return [1000, 5000, 10000, 50000, 100000];
            } else {
                return [10000, 50000, 100000, 500000, 1000000];
            }
        });

        // 에러 상태
        const hasError = computed(() => !!props.errors.money);
        const errorMessage = computed(() => props.errors.money || '');

        // 숫자 포맷팅 함수들
        const formatDisplayAmount = (amount) => {
            if (!amount) return '';
            return new Intl.NumberFormat('ko-KR').format(amount);
        };

        const parseAmount = (value) => {
            if (!value) return '';
            // 숫자와 소수점만 남기기
            const cleanValue = value.replace(/[^\d.]/g, '');
            // 소수점이 여러개면 첫번째만 유지
            const parts = cleanValue.split('.');
            if (parts.length > 2) {
                return parts[0] + '.' + parts.slice(1).join('');
            }
            return cleanValue;
        };

        const normalizeAmount = (value) => {
            if (!value) return '';
            const parsed = parseFloat(value);
            if (isNaN(parsed)) return '';
            
            // 소수점 2자리까지만 허용
            return parsed.toFixed(2).replace(/\.?0+$/, '');
        };

        // 계산된 값들
        const formattedAmount = computed(() => {
            const cleanAmount = parseAmount(displayAmount.value);
            const parsed = parseFloat(cleanAmount);
            return isNaN(parsed) ? 0 : parsed;
        });

        const readableAmount = computed(() => {
            if (!formattedAmount.value) return '';
            
            const units = ['', '만', '억', '조'];
            let amount = formattedAmount.value;
            let unitIndex = 0;
            
            while (amount >= 10000 && unitIndex < units.length - 1) {
                amount /= 10000;
                unitIndex++;
            }
            
            const formatted = amount % 1 === 0 ? amount.toString() : amount.toFixed(1);
            return `${formatted}${units[unitIndex]}원`;
        });

        // 이벤트 핸들러들
        const handleAmountInput = (value) => {
            isTyping.value = true;
            
            // 자동 포맷팅
            const cleanValue = parseAmount(value);
            const parsed = parseFloat(cleanValue);
            
            if (!isNaN(parsed) && parsed > 0) {
                displayAmount.value = formatDisplayAmount(parsed);
                updateModelValue(parsed);
                
                // 입력 중에는 자동 진행하지 않음
                clearAutoAdvanceTimer();
                
                // 유효성 검사
                emit('validate', 'money', true, '');
            } else if (value === '') {
                displayAmount.value = '';
                updateModelValue('');
                emit('validate', 'money', false, '');
            } else {
                emit('validate', 'money', false, '올바른 금액을 입력해주세요');
            }
        };

        const handleAmountBlur = () => {
            isTyping.value = false;
            
            if (formattedAmount.value > 0) {
                // 소수점 정규화
                const normalized = normalizeAmount(formattedAmount.value);
                displayAmount.value = formatDisplayAmount(normalized);
                updateModelValue(parseFloat(normalized));
                
                // 블러 시에는 자동 진행하지 않음
            }
        };

        const handleKeydown = (e) => {
            // Enter 키 처리
            if (e.key === 'Enter' && formattedAmount.value > 0) {
                e.preventDefault();
                emit('next');
            }
            
            // 숫자, 백스페이스, 삭제, 화살표, 소수점만 허용
            const allowedKeys = ['Backspace', 'Delete', 'ArrowLeft', 'ArrowRight', 'Tab'];
            const isNumber = /^\d$/.test(e.key);
            const isDecimal = e.key === '.' && !displayAmount.value.includes('.');
            
            if (!isNumber && !isDecimal && !allowedKeys.includes(e.key)) {
                e.preventDefault();
            }
        };

        const handlePaste = (e) => {
            e.preventDefault();
            const pastedText = (e.clipboardData || window.clipboardData).getData('text');
            const cleanValue = parseAmount(pastedText);
            
            if (cleanValue && !isNaN(parseFloat(cleanValue))) {
                handleAmountInput(cleanValue);
            }
        };

        const selectQuickAmount = (amount) => {
            const currentAmount = formattedAmount.value || 0;
            const newAmount = currentAmount + amount;
            
            displayAmount.value = formatDisplayAmount(newAmount);
            updateModelValue(newAmount);
            emit('validate', 'money', true, '');
            
            // 빠른 선택 시에는 자동 진행하지 않음
        };

        const updateModelValue = (value) => {
            const updated = { ...props.modelValue, money: value };
            emit('update:modelValue', updated);
        };

        const clearAutoAdvanceTimer = () => {
            if (autoAdvanceTimer.value) {
                clearTimeout(autoAdvanceTimer.value);
                autoAdvanceTimer.value = null;
            }
        };

        // 초기값 설정
        watch(() => props.modelValue.money, (newValue) => {
            if (newValue && newValue !== formattedAmount.value) {
                displayAmount.value = formatDisplayAmount(newValue);
            }
        }, { immediate: true });

        onMounted(() => {
            nextTick(() => {
                if (amountInputRef.value) {
                    amountInputRef.value.focus();
                }
            });
        });

        return {
            amountInputRef,
            displayAmount,
            quickAmounts,
            hasError,
            errorMessage,
            formattedAmount,
            readableAmount,
            formatDisplayAmount,
            handleAmountInput,
            handleAmountBlur,
            handleKeydown,
            handlePaste,
            selectQuickAmount
        };
    }
}
</script>

<style scoped>
.amount-step {
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

.amount-input-container {
    @apply mb-8;
}

.amount-input-wrapper {
    @apply relative;
}

.amount-input {
    @apply w-full;
}

:deep(.amount-input .el-input__inner) {
    @apply text-center text-3xl font-bold py-6 border-2;
    height: 80px;
}

:deep(.amount-input.error .el-input__inner) {
    @apply border-red-500;
}

:deep(.amount-input .el-input__inner:focus) {
    @apply border-blue-500 ring-2 ring-blue-100;
}

.amount-suffix {
    @apply text-2xl font-semibold text-gray-600 pr-4;
}

.error-message {
    @apply mt-2 text-sm text-red-600 text-center;
}

.amount-preview {
    @apply mt-2 text-center text-lg text-blue-600 font-medium;
}

.quick-amounts {
    @apply text-center;
}

.quick-amounts-label {
    @apply text-sm text-gray-600 mb-3;
}

.quick-amounts-grid {
    @apply grid grid-cols-3 gap-2 max-w-xs mx-auto;
}

.quick-amount-btn {
    @apply px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg transition-colors duration-200;
}

.quick-amount-btn:hover {
    @apply bg-gray-50 border-gray-400;
}

.quick-amount-btn:active {
    @apply bg-gray-100;
}

/* 모바일 최적화 */
@media (max-width: 768px) {
    .step-title {
        @apply text-xl;
    }
    
    .step-description {
        @apply text-sm;
    }
    
    :deep(.amount-input .el-input__inner) {
        @apply text-xl py-3;
        height: 56px;
        font-size: 18px; /* iOS zoom 방지하면서도 읽기 쉬운 크기 */
    }
    
    .amount-suffix {
        @apply text-lg pr-2;
    }
    
    .amount-preview {
        @apply text-base;
    }
    
    .quick-amounts-grid {
        @apply grid-cols-2 gap-3;
    }
    
    .quick-amount-btn {
        @apply px-3 py-3 text-sm;
        min-height: 44px; /* iOS 권장 터치 영역 */
    }
}

/* 접근성 */
.amount-input:focus-within {
    @apply ring-2 ring-blue-500 ring-offset-2;
}

.quick-amount-btn:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

/* 애니메이션 */
.amount-preview {
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
</style>