<template>
    <div class="payment-method-step">
        <div class="step-header">
            <h2 class="step-title">
                {{ modelValue.type === 'out' ? '결제수단을 선택하세요' : '입금경로를 선택하세요' }}
            </h2>
            <p class="step-description">
                {{ modelValue.type === 'out' ? '사용한 결제수단을 선택해주세요' : '입금받은 경로를 선택해주세요' }}
            </p>
        </div>

        <div class="payment-selection-container">
            <!-- 지출: 결제수단 선택 -->
            <div v-if="modelValue.type === 'out'" class="payment-methods">
                <!-- 카테고리별 결제수단 그룹 -->
                <div v-for="category in validPaymentMethodCategories" :key="category.id" class="payment-category">
                    <h3 class="payment-category-title">{{ category.name }}</h3>
                    <div class="payment-methods-grid">
                        <button v-for="method in category.children" :key="method.id"
                            @click="selectPaymentMethod(method.id)" class="payment-method-btn"
                            :class="{ 'selected': selectedPaymentMethodId === method.id }">
                            <div class="payment-method-icon">
                                <component :is="getPaymentMethodIcon(method.name)" class="w-5 h-5" />
                            </div>
                            <div class="payment-method-info">
                                <div class="payment-method-name">{{ method.name }}</div>
                                <div class="payment-method-description" v-if="method.description">
                                    {{ method.description }}
                                </div>
                            </div>
                        </button>
                    </div>
                </div>

                <!-- 검색 가능한 선택기 (많은 결제수단이 있을 때) -->
                <div class="payment-search" v-if="totalPaymentMethods > 12">
                    <el-select ref="paymentSelectRef" v-model="selectedPaymentMethodId" placeholder="결제수단을 검색하세요"
                        size="large" class="payment-select" :class="{ 'error': hasError }" filterable
                        @change="handlePaymentMethodChange">
                        <el-option-group v-for="category in validPaymentMethodCategories" :key="category.id"
                            :label="category.name">
                            <el-option v-for="method in category.children" :key="method.id" :label="method.name"
                                :value="method.id" class="payment-option">
                                <div class="flex items-center">
                                    <component :is="getPaymentMethodIcon(method.name)" class="w-4 h-4 mr-2" />
                                    <span>{{ method.name }}</span>
                                </div>
                            </el-option>
                        </el-option-group>
                    </el-select>
                </div>
            </div>

            <!-- 수입: 입금경로 선택 -->
            <div v-else class="deposit-paths">
                <div class="deposit-input-wrapper">
                    <el-select ref="depositSelectRef" v-model="selectedDepositPath" placeholder="입금경로를 선택하거나 입력하세요"
                        size="large" class="deposit-select" :class="{ 'error': hasError }" filterable allow-create
                        default-first-option @change="handleDepositPathChange">
                        <el-option v-for="path in depositPathOptions" :key="path.id" :label="path.label"
                            :value="path.value" class="deposit-option">
                            <div class="flex items-center">
                                <Landmark class="w-4 h-4 text-gray-500 mr-2" />
                                <span>{{ path.label }}</span>
                            </div>
                        </el-option>
                    </el-select>
                </div>

                <!-- 빠른 선택 (최근 입금경로) -->
                <div class="quick-deposits" v-if="recentDepositPaths.length > 0 && !selectedDepositPath">
                    <p class="quick-deposits-label">최근 사용한 입금경로</p>
                    <div class="quick-deposits-grid">
                        <button v-for="path in recentDepositPaths" :key="path.value"
                            @click="selectDepositPath(path.value)" class="quick-deposit-btn">
                            <Landmark class="w-4 h-4 mr-2" />
                            {{ path.label }}
                        </button>
                    </div>
                </div>
            </div>

            <div v-if="hasError" class="error-message">
                {{ errorMessage }}
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, watch, nextTick, onMounted } from 'vue';
import {
    CreditCard,
    Smartphone,
    Banknote,
    Landmark,
    Wallet,
    Building2,
    ArrowRightLeft
} from 'lucide-vue-next';
import { usePaymentMethodStore } from '../../stores/paymentMethodStore';
import { useDepositPathStore } from '../../stores/depositPathStore';

export default {
    name: 'PaymentMethodStep',
    components: {
        CreditCard,
        Smartphone,
        Banknote,
        Landmark,
        Wallet,
        Building2,
        ArrowRightLeft
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
        const paymentMethodStore = usePaymentMethodStore();
        const depositPathStore = useDepositPathStore();
        const paymentSelectRef = ref(null);
        const depositSelectRef = ref(null);
        const selectedPaymentMethodId = ref(null);
        const selectedDepositPath = ref('');

        // 에러 상태
        const hasError = computed(() => {
            if (props.modelValue.type === 'out') {
                return !!props.errors.payment_method_id;
            } else {
                return !!props.errors.deposit_path;
            }
        });

        const errorMessage = computed(() => {
            if (props.modelValue.type === 'out') {
                return props.errors.payment_method_id || '';
            } else {
                return props.errors.deposit_path || '';
            }
        });

        // 결제수단 관련
        const paymentMethodCategories = computed(() => {
            return paymentMethodStore.activePaymentMethods;
        });

        const validPaymentMethodCategories = computed(() => {
            return paymentMethodCategories.value.filter(category =>
                category.children && category.children.length > 0
            );
        });

        const totalPaymentMethods = computed(() => {
            return paymentMethodCategories.value.reduce((total, category) => {
                return total + (category.children ? category.children.length : 0);
            }, 0);
        });

        // 입금경로 관련
        const depositPathOptions = computed(() => {
            return depositPathStore.activeDepositPaths.map(path => ({
                id: path.id,
                label: path.name,
                value: path.name
            }));
        });

        const recentDepositPaths = computed(() => {
            return depositPathOptions.value.slice(0, 4);
        });

        // 아이콘 매핑
        const paymentMethodIconMap = {
            '신용카드': CreditCard,
            '체크카드': CreditCard,
            '현금': Banknote,
            '계좌이체': Building2,
            '간편결제': Smartphone,
            '기타': Wallet,
        };

        const getPaymentMethodIcon = (methodName) => {
            for (const [key, icon] of Object.entries(paymentMethodIconMap)) {
                if (methodName.includes(key)) {
                    return icon;
                }
            }
            return CreditCard;
        };

        // 이벤트 핸들러들
        const selectPaymentMethod = (methodId) => {
            selectedPaymentMethodId.value = methodId;
            updateModelValue({ payment_method_id: methodId });
            emit('validate', 'payment_method_id', true, '');

            // 선택 즉시 자동 진행
            setTimeout(() => {
                emit('auto-advance', 100);
            }, 50);
        };

        const handlePaymentMethodChange = (value) => {
            if (value) {
                selectedPaymentMethodId.value = value;
                updateModelValue({ payment_method_id: value });
                emit('validate', 'payment_method_id', true, '');

                setTimeout(() => {
                    emit('auto-advance', 200);
                }, 100);
            } else {
                emit('validate', 'payment_method_id', false, '결제수단을 선택해주세요');
            }
        };

        const selectDepositPath = (path) => {
            selectedDepositPath.value = path;
            handleDepositPathChange(path);
        };

        const handleDepositPathChange = (value) => {
            if (value) {
                selectedDepositPath.value = value;
                updateModelValue({ deposit_path: value });
                emit('validate', 'deposit_path', true, '');

                setTimeout(() => {
                    emit('auto-advance', 100);
                }, 50);
            } else {
                emit('validate', 'deposit_path', false, '입금경로를 선택해주세요');
            }
        };

        const updateModelValue = (updates) => {
            const updated = { ...props.modelValue, ...updates };
            emit('update:modelValue', updated);
        };

        // 초기값 설정
        watch(() => props.modelValue.payment_method_id, (newValue) => {
            if (newValue !== selectedPaymentMethodId.value) {
                selectedPaymentMethodId.value = newValue;
            }
        }, { immediate: true });

        watch(() => props.modelValue.deposit_path, (newValue) => {
            if (newValue !== selectedDepositPath.value) {
                selectedDepositPath.value = newValue || '';
            }
        }, { immediate: true });

        // 타입 변경 시 초기화
        watch(() => props.modelValue.type, () => {
            selectedPaymentMethodId.value = null;
            selectedDepositPath.value = '';
        });

        // 컴포넌트 마운트 시 포커스
        onMounted(() => {
            nextTick(() => {
                if (props.modelValue.type === 'out') {
                    if (totalPaymentMethods.value <= 12) {
                        const firstBtn = document.querySelector('.payment-method-btn');
                        if (firstBtn) {
                            firstBtn.focus();
                        }
                    } else if (paymentSelectRef.value) {
                        paymentSelectRef.value.focus();
                    }
                } else if (depositSelectRef.value) {
                    depositSelectRef.value.focus();
                }
            });
        });

        return {
            paymentSelectRef,
            depositSelectRef,
            selectedPaymentMethodId,
            selectedDepositPath,
            hasError,
            errorMessage,
            paymentMethodCategories,
            validPaymentMethodCategories,
            totalPaymentMethods,
            depositPathOptions,
            recentDepositPaths,
            getPaymentMethodIcon,
            selectPaymentMethod,
            handlePaymentMethodChange,
            selectDepositPath,
            handleDepositPathChange
        };
    }
}
</script>

<style scoped>
.payment-method-step {
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

.payment-selection-container {
    @apply space-y-6;
}

.payment-category {
    @apply mb-6;
}

.payment-category-title {
    @apply text-lg font-semibold text-gray-800 mb-3;
}

.payment-methods-grid {
    @apply grid grid-cols-2 gap-3;
}

.payment-method-btn {
    @apply flex items-center p-3 bg-white border-2 border-gray-200 rounded-xl transition-all duration-200 text-left;
}

.payment-method-btn:hover {
    @apply border-gray-300 bg-gray-50 transform scale-105;
}

.payment-method-btn:active {
    @apply transform scale-95;
}

.payment-method-btn.selected {
    @apply border-blue-500 bg-blue-50 ring-2 ring-blue-100;
}

.payment-method-icon {
    @apply flex-shrink-0 w-10 h-10 bg-gradient-to-br from-purple-500 to-purple-600 rounded-lg flex items-center justify-center mr-3;
}

.payment-method-btn.selected .payment-method-icon {
    @apply from-blue-500 to-blue-600;
}

.payment-method-icon svg {
    @apply text-white;
}

.payment-method-info {
    @apply flex-1 min-w-0;
}

.payment-method-name {
    @apply font-semibold text-gray-900 text-sm mb-1;
}

.payment-method-btn.selected .payment-method-name {
    @apply text-blue-700;
}

.payment-method-description {
    @apply text-xs text-gray-500 line-clamp-1;
}

.payment-search,
.deposit-input-wrapper {
    @apply mb-4;
}

.payment-select,
.deposit-select {
    @apply w-full;
}

:deep(.payment-select .el-input__inner),
:deep(.deposit-select .el-input__inner) {
    @apply h-14 text-lg border-2;
}

:deep(.payment-select.error .el-input__inner),
:deep(.deposit-select.error .el-input__inner) {
    @apply border-red-500;
}

:deep(.payment-select .el-input__inner:focus),
:deep(.deposit-select .el-input__inner:focus) {
    @apply border-blue-500 ring-2 ring-blue-100;
}

.payment-option,
.deposit-option {
    @apply py-2;
}

.quick-deposits {
    @apply text-center;
}

.quick-deposits-label {
    @apply text-sm text-gray-600 mb-3;
}

.quick-deposits-grid {
    @apply grid grid-cols-2 gap-2;
}

.quick-deposit-btn {
    @apply flex items-center justify-center px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg transition-colors duration-200;
}

.quick-deposit-btn:hover {
    @apply bg-gray-50 border-gray-400;
}

.quick-deposit-btn:active {
    @apply bg-gray-100;
}

.error-message {
    @apply text-sm text-red-600 text-center;
}

.management-links {
    @apply text-center;
}

.manage-btn {
    @apply text-blue-600;
}

.manage-btn:hover {
    @apply text-blue-700;
}

/* 모바일 최적화 */
@media (max-width: 768px) {
    .payment-methods-grid {
        @apply grid-cols-1 gap-4;
    }

    .payment-method-btn {
        @apply p-4;
        min-height: 60px;
    }

    .payment-method-icon {
        @apply w-10 h-10 mr-3;
    }

    .payment-method-name {
        @apply text-base;
    }

    :deep(.payment-select .el-input__inner),
    :deep(.deposit-select .el-input__inner) {
        @apply h-12 text-base;
        font-size: 16px;
        /* iOS zoom 방지 */
    }

    .quick-deposits-grid {
        @apply grid-cols-1 gap-3;
    }

    .quick-deposit-btn {
        @apply px-4 py-3 text-base;
        min-height: 44px;
    }
}

/* 접근성 */
.payment-method-btn:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

.payment-select:focus-within,
.deposit-select:focus-within {
    @apply ring-2 ring-blue-500 ring-offset-2;
}

.quick-deposit-btn:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

.manage-btn:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

/* 애니메이션 */
.payment-category {
    animation: fadeInUp 0.3s ease-out;
}

.quick-deposits {
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

/* 선택된 결제수단 강조 애니메이션 */
.payment-method-btn.selected {
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

/* 그리드 반응형 조정 */
@media (min-width: 640px) {
    .payment-methods-grid {
        @apply grid-cols-2;
    }
}

@media (min-width: 768px) {
    .payment-methods-grid {
        @apply grid-cols-2;
    }
}
</style>