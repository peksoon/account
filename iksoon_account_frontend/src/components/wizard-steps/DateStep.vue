<template>
    <div class="date-step">
        <div class="step-header">
            <h2 class="step-title">날짜를 선택하세요</h2>
            <p class="step-description">{{ modelValue.type === 'out' ? '지출' : '수입' }} 날짜를 선택해주세요</p>
        </div>

        <div class="date-selection-container">
            <!-- 날짜 선택기 -->
            <div class="date-picker-wrapper">
                <el-date-picker
                    ref="datePickerRef"
                    v-model="selectedDate"
                    type="date"
                    placeholder="날짜를 선택하세요"
                    format="YYYY년 MM월 DD일"
                    value-format="YYYY-MM-DD"
                    size="large"
                    class="date-picker"
                    :class="{ 'error': hasError }"
                    :editable="false"
                    :clearable="false"
                    @change="handleDateChange"
                />
                
                <div v-if="hasError" class="error-message">
                    {{ errorMessage }}
                </div>
                
                <div v-if="selectedDate" class="selected-date-info">
                    <Calendar class="w-4 h-4 mr-2 text-blue-500" />
                    <span class="text-blue-600 font-medium">{{ formatSelectedDate }}</span>
                </div>
            </div>

            <!-- 빠른 날짜 선택 -->
            <div class="quick-dates">
                <p class="quick-dates-label">빠른 선택</p>
                <div class="quick-dates-grid">
                    <button
                        v-for="quickDate in quickDates"
                        :key="quickDate.value"
                        @click="selectQuickDate(quickDate.value)"
                        class="quick-date-btn"
                        :class="{ 'selected': selectedDate === quickDate.value }"
                    >
                        <component :is="quickDate.icon" class="w-4 h-4 mr-2" />
                        {{ quickDate.label }}
                    </button>
                </div>
            </div>

            <!-- 날짜 정보 -->
            <div class="date-info" v-if="selectedDate">
                <div class="date-info-item">
                    <span class="date-info-label">요일:</span>
                    <span class="date-info-value">{{ dayOfWeek }}</span>
                </div>
                <div class="date-info-item" v-if="daysFromToday !== 0">
                    <span class="date-info-label">기준:</span>
                    <span class="date-info-value">{{ relativeDate }}</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, watch, onMounted } from 'vue';
import { Calendar, Clock, Sun, Sunrise, Sunset } from 'lucide-vue-next';

export default {
    name: 'DateStep',
    components: {
        Calendar,
        Clock,
        Sun,
        Sunrise,
        Sunset
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
        const datePickerRef = ref(null);
        const selectedDate = ref('');

        // 에러 상태
        const hasError = computed(() => !!props.errors.date);
        const errorMessage = computed(() => props.errors.date || '');

        // 오늘 날짜
        const today = new Date().toISOString().slice(0, 10);
        
        // 어제 날짜
        const yesterday = new Date();
        yesterday.setDate(yesterday.getDate() - 1);
        const yesterdayStr = yesterday.toISOString().slice(0, 10);

        // 빠른 날짜 선택 옵션
        const quickDates = computed(() => [
            {
                label: '오늘',
                value: today,
                icon: Sun
            },
            {
                label: '어제',
                value: yesterdayStr,
                icon: Sunset
            },
            {
                label: '그저께',
                value: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString().slice(0, 10),
                icon: Sunrise
            }
        ]);

        // 선택된 날짜 포맷
        const formatSelectedDate = computed(() => {
            if (!selectedDate.value) return '';
            
            const date = new Date(selectedDate.value);
            const options = { 
                year: 'numeric', 
                month: 'long', 
                day: 'numeric',
                weekday: 'long'
            };
            return date.toLocaleDateString('ko-KR', options);
        });

        // 요일
        const dayOfWeek = computed(() => {
            if (!selectedDate.value) return '';
            
            const date = new Date(selectedDate.value);
            const weekdays = ['일요일', '월요일', '화요일', '수요일', '목요일', '금요일', '토요일'];
            return weekdays[date.getDay()];
        });

        // 오늘로부터 몇 일 차이
        const daysFromToday = computed(() => {
            if (!selectedDate.value) return 0;
            
            const selected = new Date(selectedDate.value);
            const todayDate = new Date(today);
            const diffTime = selected.getTime() - todayDate.getTime();
            return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        });

        // 상대적 날짜 표현
        const relativeDate = computed(() => {
            const days = daysFromToday.value;
            
            if (days === 0) return '오늘';
            if (days === -1) return '어제';
            if (days === -2) return '그저께';
            if (days > 0) return `${days}일 후`;
            if (days < 0) return `${Math.abs(days)}일 전`;
            
            return '';
        });

        // 이벤트 핸들러들
        const handleDateChange = (value) => {
            if (value) {
                selectedDate.value = value;
                updateModelValue(value);
                emit('validate', 'date', true, '');
                
                // 날짜 선택 즉시 자동 진행
                setTimeout(() => {
                    emit('auto-advance', 100);
                }, 50);
            } else {
                emit('validate', 'date', false, '날짜를 선택해주세요');
            }
        };

        const selectQuickDate = (dateValue) => {
            selectedDate.value = dateValue;
            handleDateChange(dateValue);
        };

        const updateModelValue = (value) => {
            const updated = { ...props.modelValue, date: value };
            emit('update:modelValue', updated);
        };

        // 초기값 설정
        watch(() => props.modelValue.date, (newValue) => {
            if (newValue !== selectedDate.value) {
                selectedDate.value = newValue || '';
            }
        }, { immediate: true });

        // 컴포넌트 마운트 시 포커스하지 않음 (달력 자동 열림 방지)
        onMounted(() => {
            // 자동 포커스를 제거하여 달력이 자동으로 열리지 않도록 함
        });

        return {
            datePickerRef,
            selectedDate,
            hasError,
            errorMessage,
            quickDates,
            formatSelectedDate,
            dayOfWeek,
            daysFromToday,
            relativeDate,
            handleDateChange,
            selectQuickDate
        };
    }
}
</script>

<style scoped>
.date-step {
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

.date-selection-container {
    @apply space-y-6;
}

.date-picker-wrapper {
    @apply relative;
}

.date-picker {
    @apply w-full;
}

:deep(.date-picker .el-input__inner) {
    @apply h-14 text-lg border-2 text-center;
}

:deep(.date-picker.error .el-input__inner) {
    @apply border-red-500;
}

:deep(.date-picker .el-input__inner:focus) {
    @apply border-blue-500 ring-2 ring-blue-100;
}

:deep(.date-picker .el-input__prefix) {
    @apply left-4;
}

:deep(.date-picker .el-input__suffix) {
    @apply right-4;
}

.error-message {
    @apply mt-2 text-sm text-red-600 text-center;
}

.selected-date-info {
    @apply mt-3 flex items-center justify-center text-sm;
}

.quick-dates {
    @apply text-center;
}

.quick-dates-label {
    @apply text-sm text-gray-600 mb-3;
}

.quick-dates-grid {
    @apply grid grid-cols-3 gap-2;
}

.quick-date-btn {
    @apply flex items-center justify-center px-3 py-3 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg transition-all duration-200;
}

.quick-date-btn:hover {
    @apply bg-gray-50 border-gray-400 transform scale-105;
}

.quick-date-btn:active {
    @apply bg-gray-100 transform scale-95;
}

.quick-date-btn.selected {
    @apply bg-blue-50 border-blue-500 text-blue-700 ring-2 ring-blue-100;
}

.date-info {
    @apply bg-gray-50 rounded-lg p-4 space-y-2;
}

.date-info-item {
    @apply flex justify-between items-center;
}

.date-info-label {
    @apply text-sm text-gray-600;
}

.date-info-value {
    @apply text-sm font-medium text-gray-900;
}

/* 모바일 최적화 */
@media (max-width: 768px) {
    :deep(.date-picker .el-input__inner) {
        @apply h-12 text-base;
        font-size: 16px; /* iOS zoom 방지 */
    }
    
    .quick-dates-grid {
        @apply grid-cols-1 gap-3;
    }
    
    .quick-date-btn {
        @apply px-4 py-4 text-base;
        min-height: 44px; /* iOS 권장 터치 영역 */
    }
}

/* 접근성 */
.date-picker:focus-within {
    @apply ring-2 ring-blue-500 ring-offset-2;
}

.quick-date-btn:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

/* Element Plus 날짜 선택기 커스터마이징 */
:deep(.el-date-picker__popper) {
    @apply border-2 border-gray-200 rounded-xl shadow-xl;
}

:deep(.el-date-picker__header) {
    @apply bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-t-lg;
}

:deep(.el-date-picker__header-label) {
    @apply text-white font-semibold;
}

:deep(.el-date-picker__prev-btn),
:deep(.el-date-picker__next-btn) {
    @apply text-white;
}

:deep(.el-date-table td.current) {
    @apply bg-blue-500;
}

:deep(.el-date-table td.today) {
    @apply text-blue-600 font-bold;
}

/* 애니메이션 */
.quick-dates {
    animation: fadeInUp 0.3s ease-out;
}

.selected-date-info {
    animation: fadeInUp 0.3s ease-out;
}

.date-info {
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

/* 선택된 빠른 날짜 강조 애니메이션 */
.quick-date-btn.selected {
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

/* 모바일에서 그리드 간격 조정 */
@media (max-width: 480px) {
    .quick-dates-grid {
        @apply gap-2;
    }
}
</style>