<template>
    <div class="review-step">
        <div class="step-header">
            <h2 class="step-title">입력 내용을 확인하세요</h2>
            <p class="step-description">모든 정보가 올바른지 확인한 후 저장하세요</p>
        </div>

        <div class="review-content">
            <!-- 요약 카드 -->
            <div class="summary-card">
                <div class="summary-header">
                    <div class="transaction-type">
                        <div class="type-icon" :class="typeIconClass">
                            <component :is="typeIcon" class="w-6 h-6" />
                        </div>
                        <div class="type-info">
                            <h3 class="type-title">{{ modelValue.type === 'out' ? '지출' : '수입' }}</h3>
                            <p class="amount">{{ formatAmount(modelValue.money) }}원</p>
                        </div>
                    </div>
                </div>

                <div class="summary-details">
                    <!-- 사용자 -->
                    <div class="detail-item" @click="jumpToStep('user')">
                        <div class="detail-label">
                            <User class="w-4 h-4 mr-2" />
                            사용자
                        </div>
                        <div class="detail-value">
                            {{ modelValue.user }}
                            <ChevronRight class="w-4 h-4 ml-2 text-gray-400" />
                        </div>
                    </div>

                    <!-- 카테고리 -->
                    <div class="detail-item" @click="jumpToStep('category')">
                        <div class="detail-label">
                            <Folder class="w-4 h-4 mr-2" />
                            카테고리
                        </div>
                        <div class="detail-value">
                            {{ getCategoryName(modelValue.category_id) }}
                            <ChevronRight class="w-4 h-4 ml-2 text-gray-400" />
                        </div>
                    </div>

                    <!-- 키워드 -->
                    <div class="detail-item" @click="jumpToStep('keyword')">
                        <div class="detail-label">
                            <Tag class="w-4 h-4 mr-2" />
                            키워드
                        </div>
                        <div class="detail-value">
                            {{ modelValue.keyword_name }}
                            <ChevronRight class="w-4 h-4 ml-2 text-gray-400" />
                        </div>
                    </div>

                    <!-- 결제수단/입금경로 -->
                    <div class="detail-item" @click="jumpToStep('payment')">
                        <div class="detail-label">
                            <component :is="paymentIcon" class="w-4 h-4 mr-2" />
                            {{ modelValue.type === 'out' ? '결제수단' : '입금경로' }}
                        </div>
                        <div class="detail-value">
                            {{ getPaymentMethodName() }}
                            <ChevronRight class="w-4 h-4 ml-2 text-gray-400" />
                        </div>
                    </div>

                    <!-- 메모 (있는 경우만) -->
                    <div v-if="modelValue.memo" class="detail-item" @click="jumpToStep('memo')">
                        <div class="detail-label">
                            <FileText class="w-4 h-4 mr-2" />
                            메모
                        </div>
                        <div class="detail-value">
                            {{ modelValue.memo }}
                            <ChevronRight class="w-4 h-4 ml-2 text-gray-400" />
                        </div>
                    </div>

                    <!-- 날짜 -->
                    <div class="detail-item" @click="jumpToStep('date')">
                        <div class="detail-label">
                            <Calendar class="w-4 h-4 mr-2" />
                            날짜
                        </div>
                        <div class="detail-value">
                            {{ formatDate(modelValue.date) }}
                            <ChevronRight class="w-4 h-4 ml-2 text-gray-400" />
                        </div>
                    </div>
                </div>
            </div>

            <!-- 저장 옵션 -->
            <div class="save-options">
                <div class="save-option-item">
                    <el-checkbox v-model="continueAdding">
                        저장 후 계속 추가하기
                    </el-checkbox>
                </div>
            </div>

            <!-- 에러 메시지 -->
            <div v-if="hasErrors" class="validation-errors">
                <div class="error-title">
                    <AlertCircle class="w-4 h-4 mr-2 text-red-500" />
                    다음 항목을 확인해주세요
                </div>
                <ul class="error-list">
                    <li v-for="error in errorMessages" :key="error.field" class="error-item">
                        {{ error.message }}
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed } from 'vue';
import { 
    TrendingUp, 
    TrendingDown, 
    User, 
    Folder, 
    Tag, 
    CreditCard, 
    Landmark, 
    FileText, 
    Calendar,
    ChevronRight,
    AlertCircle
} from 'lucide-vue-next';
import { useCategoryStore } from '../../stores/categoryStore';
import { usePaymentMethodStore } from '../../stores/paymentMethodStore';

export default {
    name: 'ReviewStep',
    components: {
        TrendingUp,
        TrendingDown,
        User,
        Folder,
        Tag,
        CreditCard,
        Landmark,
        FileText,
        Calendar,
        ChevronRight,
        AlertCircle
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
    emits: ['jump-to-step', 'continue-adding'],
    setup(props, { emit }) {
        const categoryStore = useCategoryStore();
        const paymentMethodStore = usePaymentMethodStore();
        
        const continueAdding = ref(false);

        // 거래 유형 아이콘 및 스타일
        const typeIcon = computed(() => {
            return props.modelValue.type === 'out' ? TrendingDown : TrendingUp;
        });

        const typeIconClass = computed(() => {
            return props.modelValue.type === 'out' 
                ? 'bg-gradient-to-r from-red-500 to-red-600'
                : 'bg-gradient-to-r from-green-500 to-green-600';
        });

        const paymentIcon = computed(() => {
            return props.modelValue.type === 'out' ? CreditCard : Landmark;
        });

        // 금액 포맷팅
        const formatAmount = (amount) => {
            if (!amount) return '0';
            return new Intl.NumberFormat('ko-KR').format(amount);
        };

        // 날짜 포맷팅
        const formatDate = (date) => {
            if (!date) return '';
            const dateObj = new Date(date);
            return dateObj.toLocaleDateString('ko-KR', {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
                weekday: 'short'
            });
        };

        // 카테고리 이름 가져오기
        const getCategoryName = (categoryId) => {
            if (!categoryId) return '';
            const category = categoryStore.categories.find(cat => cat.id === categoryId);
            return category ? category.name : '';
        };

        // 결제수단/입금경로 이름 가져오기
        const getPaymentMethodName = () => {
            if (props.modelValue.type === 'out') {
                // 결제수단
                const methodId = props.modelValue.payment_method_id;
                if (!methodId) return '';
                
                for (const category of paymentMethodStore.activePaymentMethods) {
                    if (category.children) {
                        const method = category.children.find(m => m.id === methodId);
                        if (method) return method.name;
                    }
                }
                return '';
            } else {
                // 입금경로
                return props.modelValue.deposit_path || '';
            }
        };

        // 에러 검증
        const hasErrors = computed(() => {
            return Object.keys(props.errors).length > 0;
        });

        const errorMessages = computed(() => {
            const messages = [];
            const fieldNames = {
                user: '사용자',
                money: '금액',
                category_id: '카테고리',
                keyword_name: '키워드',
                payment_method_id: '결제수단',
                deposit_path: '입금경로',
                date: '날짜'
            };

            for (const [field, error] of Object.entries(props.errors)) {
                if (error) {
                    messages.push({
                        field,
                        message: `${fieldNames[field]}: ${error}`
                    });
                }
            }

            return messages;
        });

        // 스텝 이동
        const jumpToStep = (stepType) => {
            const stepMap = {
                'user': 1,
                'category': 2,
                'keyword': 3,
                'payment': 4,
                'memo': 5,
                'date': 6
            };
            
            const stepIndex = stepMap[stepType];
            if (stepIndex !== undefined) {
                emit('jump-to-step', stepIndex);
            }
        };

        return {
            continueAdding,
            typeIcon,
            typeIconClass,
            paymentIcon,
            hasErrors,
            errorMessages,
            formatAmount,
            formatDate,
            getCategoryName,
            getPaymentMethodName,
            jumpToStep
        };
    }
}
</script>

<style scoped>
.review-step {
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

.review-content {
    @apply space-y-6;
}

.summary-card {
    @apply bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden;
}

.summary-header {
    @apply p-6 bg-gray-50 border-b border-gray-200;
}

.transaction-type {
    @apply flex items-center;
}

.type-icon {
    @apply w-12 h-12 rounded-xl flex items-center justify-center mr-4 text-white;
}

.type-info {
    @apply flex-1;
}

.type-title {
    @apply text-lg font-semibold text-gray-900;
}

.amount {
    @apply text-2xl font-bold text-gray-900 mt-1;
}

.summary-details {
    @apply divide-y divide-gray-100;
}

.detail-item {
    @apply flex items-center justify-between p-4 cursor-pointer transition-colors duration-200;
}

.detail-item:hover {
    @apply bg-gray-50;
}

.detail-label {
    @apply flex items-center text-sm font-medium text-gray-600;
}

.detail-value {
    @apply flex items-center text-sm text-gray-900 font-medium;
}

.save-options {
    @apply bg-blue-50 border border-blue-200 rounded-lg p-4;
}

.save-option-item {
    @apply flex items-center;
}

.validation-errors {
    @apply bg-red-50 border border-red-200 rounded-lg p-4;
}

.error-title {
    @apply flex items-center text-sm font-semibold text-red-800 mb-2;
}

.error-list {
    @apply space-y-1;
}

.error-item {
    @apply text-sm text-red-700;
}

/* 모바일 최적화 */
@media (max-width: 768px) {
    .summary-card {
        @apply mx-0;
    }
    
    .type-icon {
        @apply w-10 h-10 mr-3;
    }
    
    .type-title {
        @apply text-base;
    }
    
    .amount {
        @apply text-xl;
    }
    
    .detail-item {
        @apply p-3;
    }
    
    .detail-label,
    .detail-value {
        @apply text-sm;
    }
}

/* 접근성 */
.detail-item:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

.save-option-item:focus-within {
    @apply ring-2 ring-blue-500 ring-offset-2 rounded;
}

/* 애니메이션 */
.summary-card {
    animation: fadeInUp 0.3s ease-out;
}

.save-options {
    animation: fadeInUp 0.3s ease-out 0.1s both;
}

.validation-errors {
    animation: shake 0.5s ease-out, fadeInUp 0.3s ease-out;
}

@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(20px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

@keyframes shake {
    0%, 100% {
        transform: translateX(0);
    }
    10%, 30%, 50%, 70%, 90% {
        transform: translateX(-5px);
    }
    20%, 40%, 60%, 80% {
        transform: translateX(5px);
    }
}

/* 체크박스 스타일 */
:deep(.el-checkbox) {
    @apply text-blue-700;
}

:deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
    @apply bg-blue-500 border-blue-500;
}

:deep(.el-checkbox__label) {
    @apply text-blue-700 font-medium;
}
</style>