<template>
    <div class="budget-usage-display">
        <!-- 월별 사용량 -->
        <div v-if="usage.monthly_budget > 0" class="mb-4">
            <div class="flex justify-between items-center mb-2">
                <span class="text-sm font-medium text-gray-700">월별 사용량</span>
                <span class="text-sm text-gray-600">
                    {{ formatMoney(usage.monthly_used) }} / {{ formatMoney(usage.monthly_budget) }}원
                </span>
            </div>

            <!-- 월별 진행률 바 -->
            <div class="w-full bg-gray-200 rounded-full h-3 mb-2">
                <div class="h-3 rounded-full transition-all duration-300" :class="monthlyProgressBarClass"
                    :style="{ width: `${Math.min(usage.monthly_percent, 100)}%` }"></div>
            </div>

            <div class="flex justify-between items-center text-xs">
                <span :class="monthlyPercentClass">
                    {{ usage.monthly_percent.toFixed(1) }}% 사용
                </span>
                <span v-if="usage.monthly_remaining > 0" class="text-green-600">
                    {{ formatMoney(usage.monthly_remaining) }}원 남음
                </span>
                <span v-else class="text-red-600 font-medium">
                    {{ formatMoney(Math.abs(usage.monthly_remaining)) }}원 초과
                </span>
            </div>
        </div>

        <!-- 연별 사용량 -->
        <div v-if="usage.yearly_budget > 0">
            <div class="flex justify-between items-center mb-2">
                <span class="text-sm font-medium text-gray-700">연별 사용량</span>
                <span class="text-sm text-gray-600">
                    {{ formatMoney(usage.yearly_used) }} / {{ formatMoney(usage.yearly_budget) }}원
                </span>
            </div>

            <!-- 연별 진행률 바 -->
            <div class="w-full bg-gray-200 rounded-full h-3 mb-2">
                <div class="h-3 rounded-full transition-all duration-300" :class="yearlyProgressBarClass"
                    :style="{ width: `${Math.min(usage.yearly_percent, 100)}%` }"></div>
            </div>

            <div class="flex justify-between items-center text-xs">
                <span :class="yearlyPercentClass">
                    {{ usage.yearly_percent.toFixed(1) }}% 사용
                </span>
                <span v-if="usage.yearly_remaining > 0" class="text-green-600">
                    {{ formatMoney(usage.yearly_remaining) }}원 남음
                </span>
                <span v-else class="text-red-600 font-medium">
                    {{ formatMoney(Math.abs(usage.yearly_remaining)) }}원 초과
                </span>
            </div>
        </div>

        <!-- 경고 메시지 -->
        <div v-if="usage.is_monthly_over || usage.is_yearly_over"
            class="mt-3 p-3 bg-red-50 border border-red-200 rounded-lg">
            <div class="flex items-center">
                <AlertTriangle class="w-4 h-4 text-red-500 mr-2" />
                <span class="text-sm text-red-700 font-medium">기준치 초과 경고</span>
            </div>
            <div class="mt-1 text-xs text-red-600">
                <div v-if="usage.is_monthly_over">• 월 기준치를 초과했습니다.</div>
                <div v-if="usage.is_yearly_over">• 연 기준치를 초과했습니다.</div>
            </div>
        </div>

        <!-- 양호 메시지 -->
        <div v-else-if="!isNearLimit" class="mt-3 p-3 bg-green-50 border border-green-200 rounded-lg">
            <div class="flex items-center">
                <CheckCircle class="w-4 h-4 text-green-500 mr-2" />
                <span class="text-sm text-green-700 font-medium">기준치 내 사용 중</span>
            </div>
        </div>

        <!-- 주의 메시지 (80% 이상 사용 시) -->
        <div v-else class="mt-3 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
            <div class="flex items-center">
                <AlertCircle class="w-4 h-4 text-yellow-500 mr-2" />
                <span class="text-sm text-yellow-700 font-medium">기준치 근접 주의</span>
            </div>
            <div class="mt-1 text-xs text-yellow-600">
                기준치의 80% 이상을 사용했습니다.
            </div>
        </div>
    </div>
</template>

<script>
import { computed, defineComponent } from 'vue'
import { AlertTriangle, AlertCircle, CheckCircle } from 'lucide-vue-next'

export default defineComponent({
    name: 'BudgetUsageDisplay',
    components: {
        AlertTriangle,
        AlertCircle,
        CheckCircle
    },
    props: {
        usage: {
            type: Object,
            required: true
        }
    },
    setup(props) {
        // Computed
        const monthlyProgressBarClass = computed(() => {
            if (props.usage.is_monthly_over) {
                return 'bg-red-500'
            } else if (props.usage.monthly_percent >= 80) {
                return 'bg-yellow-500'
            } else {
                return 'bg-blue-500'
            }
        })

        const yearlyProgressBarClass = computed(() => {
            if (props.usage.is_yearly_over) {
                return 'bg-red-500'
            } else if (props.usage.yearly_percent >= 80) {
                return 'bg-yellow-500'
            } else {
                return 'bg-green-500'
            }
        })

        const monthlyPercentClass = computed(() => {
            if (props.usage.is_monthly_over) {
                return 'text-red-600 font-medium'
            } else if (props.usage.monthly_percent >= 80) {
                return 'text-yellow-600 font-medium'
            } else {
                return 'text-blue-600'
            }
        })

        const yearlyPercentClass = computed(() => {
            if (props.usage.is_yearly_over) {
                return 'text-red-600 font-medium'
            } else if (props.usage.yearly_percent >= 80) {
                return 'text-yellow-600 font-medium'
            } else {
                return 'text-green-600'
            }
        })

        const isNearLimit = computed(() => {
            return (props.usage.monthly_budget > 0 && props.usage.monthly_percent >= 80) ||
                (props.usage.yearly_budget > 0 && props.usage.yearly_percent >= 80)
        })

        // Methods
        const formatMoney = (amount) => {
            return new Intl.NumberFormat('ko-KR').format(amount || 0)
        }

        return {
            monthlyProgressBarClass,
            yearlyProgressBarClass,
            monthlyPercentClass,
            yearlyPercentClass,
            isNearLimit,
            formatMoney
        }
    }
})
</script>

<style scoped>
.budget-usage-display {
    margin-bottom: 0.5rem;
}
</style>
