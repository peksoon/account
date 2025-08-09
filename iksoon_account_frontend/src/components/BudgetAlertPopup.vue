<template>
    <div v-if="isVisible" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 w-96 max-w-sm mx-4 shadow-xl">
            <!-- 헤더 -->
            <div class="flex items-center justify-between mb-4">
                <div class="flex items-center">
                    <component :is="headerIcon" :class="headerIconClass" class="w-6 h-6 mr-3" />
                    <h3 class="text-lg font-semibold" :class="headerTextClass">
                        {{ headerTitle }}
                    </h3>
                </div>
                <button @click="closePopup" class="text-gray-400 hover:text-gray-600 transition-colors">
                    <X class="w-5 h-5" />
                </button>
            </div>

            <!-- 메시지 -->
            <div class="mb-6">
                <p class="text-gray-700 mb-2">{{ message }}</p>
                <p class="text-sm text-gray-600">카테고리: {{ budgetUsage?.category_name }}</p>
            </div>

            <!-- 기준치 사용량 정보 -->
            <div v-if="budgetUsage" class="mb-6">
                <BudgetUsageDisplay :usage="budgetUsage" />
            </div>

            <!-- 추가 정보 -->
            <div v-if="showAdditionalInfo" class="mb-6 p-4 bg-gray-50 rounded-lg">
                <h4 class="text-sm font-medium text-gray-700 mb-2">지출 정보</h4>
                <div class="space-y-1 text-sm text-gray-600">
                    <div>금액: {{ formatMoney(expenseAmount) }}원</div>
                    <div>날짜: {{ expenseDate }}</div>
                    <div v-if="expenseKeyword">키워드: {{ expenseKeyword }}</div>
                </div>
            </div>

            <!-- 권장 사항 -->
            <div v-if="recommendations.length > 0" class="mb-6">
                <h4 class="text-sm font-medium text-gray-700 mb-2">권장 사항</h4>
                <ul class="space-y-1">
                    <li v-for="(recommendation, index) in recommendations" :key="index"
                        class="text-sm text-gray-600 flex items-start">
                        <span class="text-blue-500 mr-2">•</span>
                        {{ recommendation }}
                    </li>
                </ul>
            </div>

            <!-- 버튼 -->
            <div class="flex justify-end space-x-3">
                <button v-if="showBudgetManagementButton" @click="openBudgetManagement"
                    class="px-4 py-2 text-blue-600 bg-blue-100 rounded-md hover:bg-blue-200 transition-colors">
                    기준치 관리
                </button>
                <button @click="closePopup"
                    class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 transition-colors">
                    확인
                </button>
            </div>
        </div>
    </div>
</template>

<script>
import { computed, defineComponent } from 'vue'
import {
    AlertTriangle,
    AlertCircle,
    CheckCircle,
    X,
    DollarSign
} from 'lucide-vue-next'
import BudgetUsageDisplay from './BudgetUsageDisplay.vue'

export default defineComponent({
    name: 'BudgetAlertPopup',
    components: {
        AlertTriangle,
        AlertCircle,
        CheckCircle,
        X,
        DollarSign,
        BudgetUsageDisplay
    },
    props: {
        isVisible: {
            type: Boolean,
            default: false
        },
        budgetUsage: {
            type: Object,
            default: null
        },
        expenseAmount: {
            type: Number,
            default: 0
        },
        expenseDate: {
            type: String,
            default: ''
        },
        expenseKeyword: {
            type: String,
            default: ''
        },
        showAdditionalInfo: {
            type: Boolean,
            default: true
        },
        showBudgetManagementButton: {
            type: Boolean,
            default: true
        }
    },
    emits: ['close', 'openBudgetManagement'],
    setup(props, { emit }) {
        // Computed
        const isNearLimit = computed(() => {
            if (!props.budgetUsage) return false

            return (props.budgetUsage.monthly_budget > 0 && props.budgetUsage.monthly_percent >= 80) ||
                (props.budgetUsage.yearly_budget > 0 && props.budgetUsage.yearly_percent >= 80)
        })

        const alertType = computed(() => {
            if (!props.budgetUsage) return 'success'

            if (props.budgetUsage.is_monthly_over || props.budgetUsage.is_yearly_over) {
                return 'error'
            } else if (isNearLimit.value) {
                return 'warning'
            } else {
                return 'success'
            }
        })

        const headerIcon = computed(() => {
            switch (alertType.value) {
                case 'error':
                    return AlertTriangle
                case 'warning':
                    return AlertCircle
                case 'success':
                    return CheckCircle
                default:
                    return DollarSign
            }
        })

        const headerIconClass = computed(() => {
            switch (alertType.value) {
                case 'error':
                    return 'text-red-500'
                case 'warning':
                    return 'text-yellow-500'
                case 'success':
                    return 'text-green-500'
                default:
                    return 'text-blue-500'
            }
        })

        const headerTextClass = computed(() => {
            switch (alertType.value) {
                case 'error':
                    return 'text-red-700'
                case 'warning':
                    return 'text-yellow-700'
                case 'success':
                    return 'text-green-700'
                default:
                    return 'text-gray-700'
            }
        })

        const headerTitle = computed(() => {
            switch (alertType.value) {
                case 'error':
                    return '기준치 초과 알림'
                case 'warning':
                    return '기준치 근접 주의'
                case 'success':
                    return '지출 기록 완료'
                default:
                    return '지출 기록 완료'
            }
        })

        const message = computed(() => {
            if (!props.budgetUsage) {
                return '지출이 성공적으로 기록되었습니다.'
            }

            switch (alertType.value) {
                case 'error':
                    return '설정된 기준치를 초과했습니다. 지출을 검토해 보세요.'
                case 'warning':
                    return '기준치의 80% 이상을 사용했습니다. 주의하세요.'
                case 'success':
                    return '지출이 성공적으로 기록되었습니다. 기준치 내에서 잘 관리하고 있습니다.'
                default:
                    return '지출이 성공적으로 기록되었습니다.'
            }
        })

        const recommendations = computed(() => {
            if (!props.budgetUsage) return []

            const recs = []

            if (props.budgetUsage.is_monthly_over) {
                recs.push('이번 달 남은 기간 동안 해당 카테고리 지출을 줄여보세요.')
                recs.push('다른 카테고리에서 절약할 수 있는 부분을 찾아보세요.')
            } else if (props.budgetUsage.is_yearly_over) {
                recs.push('연간 기준치를 검토하고 필요시 조정해보세요.')
                recs.push('남은 기간 동안 계획적인 지출 관리가 필요합니다.')
            } else if (isNearLimit.value) {
                recs.push('남은 기간 동안 해당 카테고리 지출을 조절해보세요.')
                recs.push('불필요한 지출이 있는지 점검해보세요.')
            } else {
                recs.push('현재 기준치를 잘 지키고 있습니다. 이대로 유지하세요.')
            }

            return recs
        })

        // Methods
        const formatMoney = (amount) => {
            return new Intl.NumberFormat('ko-KR').format(amount || 0)
        }

        const closePopup = () => {
            emit('close')
        }

        const openBudgetManagement = () => {
            emit('openBudgetManagement')
            closePopup()
        }

        return {
            alertType,
            isNearLimit,
            headerIcon,
            headerIconClass,
            headerTextClass,
            headerTitle,
            message,
            recommendations,
            formatMoney,
            closePopup,
            openBudgetManagement
        }
    }
})
</script>

<style scoped>
/* 추가 스타일이 필요한 경우 여기에 작성 */
</style>
