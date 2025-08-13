<template>
    <div class="summary-panel">
        <div class="summary-header">
            <h3 class="summary-title">입력 요약</h3>
            <div class="summary-progress">
                {{ currentStep + 1 }} / 8 단계
            </div>
        </div>

        <div class="summary-content">
            <!-- 거래 유형 -->
            <div class="summary-item">
                <div class="summary-label">거래 유형</div>
                <div class="summary-value type-value" :class="typeClass">
                    <component :is="typeIcon" class="w-4 h-4 mr-2" />
                    {{ data.type === 'out' ? '지출' : '수입' }}
                </div>
            </div>

            <!-- 금액 -->
            <div class="summary-item" :class="{ 'completed': data.money }">
                <div class="summary-label">금액</div>
                <div class="summary-value">
                    <span v-if="data.money" class="amount">{{ formatAmount(data.money) }}원</span>
                    <span v-else class="empty-value">미입력</span>
                </div>
            </div>

            <!-- 사용자 -->
            <div class="summary-item" :class="{ 'completed': data.user }">
                <div class="summary-label">사용자</div>
                <div class="summary-value">
                    <User class="w-4 h-4 mr-2 text-gray-400" />
                    <span v-if="data.user">{{ data.user }}</span>
                    <span v-else class="empty-value">미선택</span>
                </div>
            </div>

            <!-- 카테고리 -->
            <div class="summary-item" :class="{ 'completed': data.category_id }">
                <div class="summary-label">카테고리</div>
                <div class="summary-value">
                    <Folder class="w-4 h-4 mr-2 text-gray-400" />
                    <span v-if="data.category_id">{{ getCategoryName(data.category_id) }}</span>
                    <span v-else class="empty-value">미선택</span>
                </div>
            </div>

            <!-- 키워드 -->
            <div class="summary-item" :class="{ 'completed': data.keyword_name }">
                <div class="summary-label">키워드</div>
                <div class="summary-value">
                    <Tag class="w-4 h-4 mr-2 text-gray-400" />
                    <span v-if="data.keyword_name">{{ data.keyword_name }}</span>
                    <span v-else class="empty-value">미입력</span>
                </div>
            </div>

            <!-- 결제수단/입금경로 -->
            <div class="summary-item" :class="{ 'completed': getPaymentValue() }">
                <div class="summary-label">
                    {{ data.type === 'out' ? '결제수단' : '입금경로' }}
                </div>
                <div class="summary-value">
                    <component :is="paymentIcon" class="w-4 h-4 mr-2 text-gray-400" />
                    <span v-if="getPaymentValue()">{{ getPaymentMethodName() }}</span>
                    <span v-else class="empty-value">미선택</span>
                </div>
            </div>

            <!-- 메모 -->
            <div class="summary-item" :class="{ 'completed': data.memo, 'optional': true }">
                <div class="summary-label">메모</div>
                <div class="summary-value">
                    <FileText class="w-4 h-4 mr-2 text-gray-400" />
                    <span v-if="data.memo">{{ truncateMemo(data.memo) }}</span>
                    <span v-else class="empty-value optional">선택사항</span>
                </div>
            </div>

            <!-- 날짜 -->
            <div class="summary-item" :class="{ 'completed': data.date }">
                <div class="summary-label">날짜</div>
                <div class="summary-value">
                    <Calendar class="w-4 h-4 mr-2 text-gray-400" />
                    <span v-if="data.date">{{ formatDate(data.date) }}</span>
                    <span v-else class="empty-value">미선택</span>
                </div>
            </div>
        </div>

        <!-- 완료도 표시 -->
        <div class="completion-status">
            <div class="completion-header">
                <span class="completion-label">완료도</span>
                <span class="completion-percentage">{{ completionPercentage }}%</span>
            </div>
            <div class="completion-bar">
                <div 
                    class="completion-progress" 
                    :style="{ width: `${completionPercentage}%` }"
                ></div>
            </div>
            <div class="completion-text">
                {{ completedFields }} / {{ totalRequiredFields }} 항목 완료
            </div>
        </div>
    </div>
</template>

<script>
import { computed } from 'vue';
import { 
    TrendingUp, 
    TrendingDown, 
    User, 
    Folder, 
    Tag, 
    CreditCard, 
    Landmark, 
    FileText, 
    Calendar 
} from 'lucide-vue-next';
import { useCategoryStore } from '../../stores/categoryStore';
import { usePaymentMethodStore } from '../../stores/paymentMethodStore';

export default {
    name: 'SummaryPanel',
    components: {
        TrendingUp,
        TrendingDown,
        User,
        Folder,
        Tag,
        CreditCard,
        Landmark,
        FileText,
        Calendar
    },
    props: {
        data: {
            type: Object,
            required: true
        },
        currentStep: {
            type: Number,
            required: true
        }
    },
    emits: ['jump-to-step'],
    setup(props) {
        const categoryStore = useCategoryStore();
        const paymentMethodStore = usePaymentMethodStore();

        // 거래 유형 관련
        const typeIcon = computed(() => {
            return props.data.type === 'out' ? TrendingDown : TrendingUp;
        });

        const typeClass = computed(() => {
            return props.data.type === 'out' ? 'expense' : 'income';
        });

        const paymentIcon = computed(() => {
            return props.data.type === 'out' ? CreditCard : Landmark;
        });

        // 포맷팅 함수들
        const formatAmount = (amount) => {
            if (!amount) return '';
            return new Intl.NumberFormat('ko-KR').format(amount);
        };

        const formatDate = (date) => {
            if (!date) return '';
            const dateObj = new Date(date);
            const today = new Date();
            const diffTime = dateObj.getTime() - today.getTime();
            const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
            
            if (diffDays === 0) return '오늘';
            if (diffDays === -1) return '어제';
            if (diffDays === 1) return '내일';
            
            return dateObj.toLocaleDateString('ko-KR', {
                month: 'short',
                day: 'numeric'
            });
        };

        const truncateMemo = (memo) => {
            if (!memo) return '';
            return memo.length > 20 ? memo.substring(0, 20) + '...' : memo;
        };

        // 데이터 조회 함수들
        const getCategoryName = (categoryId) => {
            if (!categoryId) return '';
            const category = categoryStore.categories.find(cat => cat.id === categoryId);
            return category ? category.name : '';
        };

        const getPaymentValue = () => {
            return props.data.type === 'out' 
                ? props.data.payment_method_id 
                : props.data.deposit_path;
        };

        const getPaymentMethodName = () => {
            if (props.data.type === 'out') {
                const methodId = props.data.payment_method_id;
                if (!methodId) return '';
                
                for (const category of paymentMethodStore.activePaymentMethods) {
                    if (category.children) {
                        const method = category.children.find(m => m.id === methodId);
                        if (method) return method.name;
                    }
                }
                return '';
            } else {
                return props.data.deposit_path || '';
            }
        };

        // 완료도 계산
        const requiredFields = computed(() => [
            'money',
            'user', 
            'category_id',
            'keyword_name',
            props.data.type === 'out' ? 'payment_method_id' : 'deposit_path',
            'date'
        ]);

        const completedFields = computed(() => {
            return requiredFields.value.filter(field => {
                const value = props.data[field];
                return value !== null && value !== undefined && value !== '';
            }).length;
        });

        const totalRequiredFields = computed(() => requiredFields.value.length);

        const completionPercentage = computed(() => {
            if (totalRequiredFields.value === 0) return 0;
            return Math.round((completedFields.value / totalRequiredFields.value) * 100);
        });

        return {
            typeIcon,
            typeClass,
            paymentIcon,
            formatAmount,
            formatDate,
            truncateMemo,
            getCategoryName,
            getPaymentValue,
            getPaymentMethodName,
            completedFields,
            totalRequiredFields,
            completionPercentage
        };
    }
}
</script>

<style scoped>
.summary-panel {
    @apply h-full flex flex-col;
}

.summary-header {
    @apply p-6 border-b border-gray-200 bg-white;
}

.summary-title {
    @apply text-lg font-semibold text-gray-900 mb-2;
}

.summary-progress {
    @apply text-sm text-gray-500;
}

.summary-content {
    @apply flex-1 p-4 space-y-4 overflow-y-auto;
}

.summary-item {
    @apply p-3 rounded-lg border border-gray-100 transition-all duration-200;
}

.summary-item.completed {
    @apply bg-green-50 border-green-200;
}

.summary-item.optional {
    @apply border-gray-100;
}

.summary-label {
    @apply text-xs font-medium text-gray-600 mb-2;
}

.summary-value {
    @apply flex items-center text-sm;
}

.summary-item.completed .summary-value {
    @apply text-green-800;
}

.type-value {
    @apply px-2 py-1 rounded-md text-white text-xs font-semibold;
}

.type-value.expense {
    @apply bg-red-500;
}

.type-value.income {
    @apply bg-green-500;
}

.amount {
    @apply font-bold text-base;
}

.empty-value {
    @apply text-gray-400 italic text-xs;
}

.empty-value.optional {
    @apply text-gray-300;
}

.completion-status {
    @apply p-4 border-t border-gray-200 bg-gray-50;
}

.completion-header {
    @apply flex justify-between items-center mb-2;
}

.completion-label {
    @apply text-sm font-medium text-gray-700;
}

.completion-percentage {
    @apply text-sm font-bold text-blue-600;
}

.completion-bar {
    @apply w-full bg-gray-200 rounded-full h-2 mb-2;
}

.completion-progress {
    @apply bg-gradient-to-r from-blue-500 to-blue-600 h-2 rounded-full transition-all duration-500;
}

.completion-text {
    @apply text-xs text-gray-600 text-center;
}

/* 스크롤바 스타일 */
.summary-content::-webkit-scrollbar {
    @apply w-1;
}

.summary-content::-webkit-scrollbar-track {
    @apply bg-gray-100;
}

.summary-content::-webkit-scrollbar-thumb {
    @apply bg-gray-300 rounded-full;
}

.summary-content::-webkit-scrollbar-thumb:hover {
    @apply bg-gray-400;
}

/* 애니메이션 */
.summary-item {
    animation: fadeInRight 0.3s ease-out;
}

.summary-item.completed {
    animation: fadeInRight 0.3s ease-out, pulse 0.5s ease-out 0.3s;
}

@keyframes fadeInRight {
    from {
        opacity: 0;
        transform: translateX(10px);
    }
    to {
        opacity: 1;
        transform: translateX(0);
    }
}

@keyframes pulse {
    0%, 100% {
        transform: scale(1);
    }
    50% {
        transform: scale(1.02);
    }
}

/* 완료 상태 강조 */
.summary-item.completed::before {
    content: '';
    @apply absolute left-1 top-1/2 w-1 h-4 bg-green-500 rounded-full transform -translate-y-1/2;
}

.summary-item {
    @apply relative;
}

/* 아이콘 색상 조정 */
.summary-item.completed svg {
    @apply text-green-600;
}
</style>