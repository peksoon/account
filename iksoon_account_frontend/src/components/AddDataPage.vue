<template>
    <div class="add-data-page">
        <!-- 헤더 -->
        <div class="page-header">
            <div class="container-responsive">
                <div class="flex items-center justify-between p-4">
                    <div class="flex items-center">
                        <button @click="handleBack" class="back-btn">
                            <ArrowLeft class="w-5 h-5" />
                        </button>
                        <div class="ml-3">
                            <h1 class="text-xl font-bold text-gray-900">새 데이터 추가</h1>
                            <p class="text-sm text-gray-600">단계별로 정보를 입력해보세요</p>
                        </div>
                    </div>
                    <button @click="handleClose" class="close-btn">
                        <X class="w-5 h-5" />
                    </button>
                </div>
            </div>
        </div>

        <!-- 타입 선택 (고정) -->
        <div class="type-selection-bar">
            <div class="container-responsive">
                <div class="p-4">
                    <div class="grid grid-cols-2 gap-3 max-w-md mx-auto">
                        <div @click="formData.type = 'out'"
                            class="type-btn" 
                            :class="formData.type === 'out' ? 'selected expense' : ''">
                            <TrendingDown class="w-5 h-5 mr-2" />
                            <span class="font-semibold">지출</span>
                        </div>
                        <div @click="formData.type = 'in'"
                            class="type-btn" 
                            :class="formData.type === 'in' ? 'selected income' : ''">
                            <TrendingUp class="w-5 h-5 mr-2" />
                            <span class="font-semibold">수입</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 스테퍼 -->
        <div class="stepper-section">
            <div class="container-responsive">
                <AddDataStepper 
                    :current-step="currentStep"
                    :steps="steps"
                    @jump-to-step="jumpToStep"
                />
            </div>
        </div>

        <!-- 메인 컨텐츠 -->
        <div class="main-content">
            <div class="container-responsive">
                <div class="content-layout">
                    <!-- Desktop: 두 열 레이아웃 -->
                    <div v-if="!isMobile" class="desktop-layout">
                        <div class="step-content-area">
                            <component 
                                :is="currentStepComponent" 
                                v-model="formData"
                                :errors="errors"
                                @next="handleNext"
                                @auto-advance="handleAutoAdvance"
                                @validate="handleValidate"
                                @open-category-manager="openCategoryManager"
                                @open-payment-method-manager="openPaymentMethodManager"
                                @open-deposit-path-manager="openDepositPathManager"
                                @open-keyword-manager="openKeywordManager"
                                @open-user-manager="openUserManager"
                            />
                        </div>
                        <div class="summary-area">
                            <AddDataSummary 
                                :data="formData" 
                                :current-step="currentStep"
                                @jump-to-step="jumpToStep"
                            />
                        </div>
                    </div>

                    <!-- Mobile: 단일 열 레이아웃 -->
                    <div v-else class="mobile-layout">
                        <component 
                            :is="currentStepComponent" 
                            v-model="formData"
                            :errors="errors"
                            @next="handleNext"
                            @auto-advance="handleAutoAdvance"
                            @validate="handleValidate"
                            @open-category-manager="openCategoryManager"
                            @open-payment-method-manager="openPaymentMethodManager"
                            @open-deposit-path-manager="openDepositPathManager"
                            @open-keyword-manager="openKeywordManager"
                            @open-user-manager="openUserManager"
                        />
                    </div>
                </div>
            </div>
        </div>

        <!-- 하단 액션 바 -->
        <div class="action-section">
            <div class="container-responsive">
                <AddDataActions 
                    :current-step="currentStep"
                    :total-steps="steps.length"
                    :can-proceed="canProceed"
                    :can-skip="canSkipCurrentStep"
                    :is-form-valid="isFormValid"
                    :saving="saving"
                    @back="handleBack"
                    @next="handleNext"
                    @skip="handleSkip"
                    @save="handleSave"
                />
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import {
    ArrowLeft,
    X,
    TrendingUp,
    TrendingDown
} from 'lucide-vue-next';

// Components
import AddDataStepper from './add-data/AddDataStepper.vue';
import AddDataSummary from './add-data/AddDataSummary.vue';
import AddDataActions from './add-data/AddDataActions.vue';

// Step Components
import AmountStep from './wizard-steps/AmountStep.vue';
import UserStep from './wizard-steps/UserStep.vue';
import CategoryStep from './wizard-steps/CategoryStep.vue';
import KeywordStep from './wizard-steps/KeywordStep.vue';
import PaymentMethodStep from './wizard-steps/PaymentMethodStep.vue';
import MemoStep from './wizard-steps/MemoStep.vue';
import DateStep from './wizard-steps/DateStep.vue';
import ReviewStep from './wizard-steps/ReviewStep.vue';

// Stores
import { useCategoryStore } from '../stores/categoryStore';
import { usePaymentMethodStore } from '../stores/paymentMethodStore';
import { useDepositPathStore } from '../stores/depositPathStore';
import { useKeywordStore } from '../stores/keywordStore';
import { useUserStore } from '../stores/userStore';
import { useBudgetStore } from '../stores/budgetStore';

export default {
    name: 'AddDataPage',
    components: {
        ArrowLeft,
        X,
        TrendingUp,
        TrendingDown,
        AddDataStepper,
        AddDataSummary,
        AddDataActions,
        AmountStep,
        UserStep,
        CategoryStep,
        KeywordStep,
        PaymentMethodStep,
        MemoStep,
        DateStep,
        ReviewStep
    },
    props: {
        initialData: {
            type: Object,
            default: () => ({})
        },
        selectedDate: {
            type: String,
            default: null
        }
    },
    emits: ['save', 'close', 'open-category-manager', 'open-payment-method-manager', 'open-deposit-path-manager', 'open-keyword-manager', 'open-user-manager', 'budget-alert', 'budget-save-success'],
    setup(props, { emit }) {
        const router = useRouter();
        
        // Stores
        const categoryStore = useCategoryStore();
        const paymentMethodStore = usePaymentMethodStore();
        const depositPathStore = useDepositPathStore();
        const keywordStore = useKeywordStore();
        const userStore = useUserStore();
        const budgetStore = useBudgetStore();

        // State
        const currentStep = ref(0);
        const saving = ref(false);
        const errors = ref({});
        const autoAdvanceTimeouts = ref({});

        // 폼 데이터 초기화
        const initializeFormData = (initialData) => {
            const defaultDate = new Date().toISOString().slice(0, 10);
            let dateToUse = defaultDate;
            
            if (props.selectedDate && props.selectedDate !== '') {
                dateToUse = props.selectedDate;
            } else if (initialData && initialData.date && initialData.date !== '') {
                dateToUse = initialData.date;
            }

            return {
                type: 'out',
                user: '',
                money: '',
                category_id: null,
                keyword_name: '',
                payment_method_id: null,
                memo: '',
                deposit_path: '',
                ...initialData,
                date: dateToUse,
            };
        };

        const formData = ref(initializeFormData(props.initialData));

        // 모바일 체크
        const isMobile = computed(() => {
            if (typeof window === 'undefined') return false;
            return window.innerWidth < 768;
        });

        // 스텝 정의
        const steps = computed(() => [
            { id: 'amount', label: '금액', component: 'AmountStep', required: true },
            { id: 'user', label: '사용자', component: 'UserStep', required: true },
            { id: 'category', label: '카테고리', component: 'CategoryStep', required: true },
            { id: 'keyword', label: '키워드', component: 'KeywordStep', required: true },
            { 
                id: 'payment', 
                label: formData.value.type === 'out' ? '결제수단' : '입금경로', 
                component: 'PaymentMethodStep', 
                required: true 
            },
            { id: 'memo', label: '메모', component: 'MemoStep', required: false },
            { id: 'date', label: '날짜', component: 'DateStep', required: true },
            { id: 'review', label: '검토 및 저장', component: 'ReviewStep', required: false }
        ]);

        const currentStepData = computed(() => steps.value[currentStep.value]);
        const currentStepComponent = computed(() => currentStepData.value.component);

        // 다음 단계로 진행 가능 여부
        const canProceed = computed(() => {
            if (!currentStepData.value.required) return true;
            
            switch (currentStepData.value.id) {
                case 'amount':
                    return formData.value.money && parseFloat(formData.value.money) > 0;
                case 'user':
                    return formData.value.user;
                case 'category':
                    return formData.value.category_id;
                case 'keyword':
                    return formData.value.keyword_name;
                case 'payment':
                    return formData.value.type === 'out' 
                        ? formData.value.payment_method_id 
                        : formData.value.deposit_path;
                case 'date':
                    return formData.value.date;
                default:
                    return true;
            }
        });

        const canSkipCurrentStep = computed(() => {
            return !currentStepData.value.required;
        });

        const isFormValid = computed(() => {
            return formData.value.user &&
                formData.value.money > 0 &&
                formData.value.category_id &&
                formData.value.keyword_name &&
                (formData.value.type === 'in' ? formData.value.deposit_path : formData.value.payment_method_id) &&
                formData.value.date;
        });

        // 이벤트 핸들러들
        const handleNext = async () => {
            if (!canProceed.value) return;
            
            clearAutoAdvanceTimeout();
            
            if (currentStep.value < steps.value.length - 1) {
                currentStep.value++;
                await nextTick();
                focusCurrentStepInput();
            }
        };

        const handleBack = () => {
            if (currentStep.value > 0) {
                clearAutoAdvanceTimeout();
                currentStep.value--;
            } else {
                handleClose();
            }
        };

        const handleSkip = () => {
            if (canSkipCurrentStep.value) {
                handleNext();
            }
        };

        const handleAutoAdvance = (delay = 600) => {
            if (!canProceed.value) return;
            
            clearAutoAdvanceTimeout();
            autoAdvanceTimeouts.value[currentStep.value] = setTimeout(() => {
                handleNext();
            }, delay);
        };

        const clearAutoAdvanceTimeout = () => {
            if (autoAdvanceTimeouts.value[currentStep.value]) {
                clearTimeout(autoAdvanceTimeouts.value[currentStep.value]);
                delete autoAdvanceTimeouts.value[currentStep.value];
            }
        };

        const jumpToStep = (stepIndex) => {
            if (stepIndex <= currentStep.value) {
                clearAutoAdvanceTimeout();
                currentStep.value = stepIndex;
            }
        };

        const handleValidate = (field, isValid, errorMessage) => {
            if (isValid) {
                delete errors.value[field];
            } else {
                errors.value[field] = errorMessage;
            }
        };

        const handleSave = async () => {
            if (!isFormValid.value) return;

            saving.value = true;
            try {
                const accountData = { ...formData.value };

                if (formData.value.type === 'out') {
                    const response = await budgetStore.createOutAccountWithBudget(accountData);
                    
                    if (response.budget_usage) {
                        emit('budget-alert', {
                            budgetUsage: response.budget_usage,
                            expenseAmount: accountData.money,
                            expenseDate: accountData.date,
                            expenseKeyword: accountData.keyword_name
                        });
                    }
                    
                    emit('budget-save-success');
                    ElMessage.success('지출이 성공적으로 기록되었습니다.');
                } else {
                    if (formData.value.keyword_name && formData.value.category_id) {
                        await keywordStore.useKeyword(formData.value.category_id, formData.value.keyword_name);
                    }
                    
                    emit('save', accountData);
                    ElMessage.success('수입이 성공적으로 기록되었습니다.');
                }

                // 저장 후 홈으로 이동
                router.push('/');
            } catch (error) {
                console.error('저장 실패:', error);
                ElMessage.error('저장 중 오류가 발생했습니다.');
            } finally {
                saving.value = false;
            }
        };

        const handleClose = async () => {
            const hasUnsavedChanges = formData.value.money || formData.value.user || formData.value.memo;
            
            if (hasUnsavedChanges) {
                try {
                    await ElMessageBox.confirm(
                        '저장되지 않은 변경사항이 있습니다. 정말로 나가시겠습니까?',
                        '확인',
                        {
                            confirmButtonText: '나가기',
                            cancelButtonText: '계속 편집',
                            type: 'warning',
                        }
                    );
                } catch {
                    return;
                }
            }
            
            router.push('/');
        };

        const focusCurrentStepInput = () => {
            nextTick(() => {
                const firstInput = document.querySelector('.step-content-area input, .step-content-area .el-input__inner, .step-content-area .el-select');
                if (firstInput) {
                    firstInput.focus();
                }
            });
        };

        // 관리자 모달 열기
        const openCategoryManager = () => {
            emit('open-category-manager');
        };

        const openPaymentMethodManager = () => {
            emit('open-payment-method-manager');
        };

        const openDepositPathManager = () => {
            emit('open-deposit-path-manager');
        };

        const openUserManager = () => {
            emit('open-user-manager');
        };

        const openKeywordManager = (categoryId) => {
            emit('open-keyword-manager', categoryId);
        };

        // 초기 데이터 로드
        const loadInitialData = async () => {
            try {
                await Promise.all([
                    categoryStore.fetchCategories(),
                    paymentMethodStore.fetchPaymentMethods(),
                    depositPathStore.fetchDepositPaths(),
                    userStore.fetchUsers()
                ]);
            } catch (error) {
                console.error('초기 데이터 로드 오류:', error);
            }
        };

        // Watchers
        watch(() => props.initialData, (newVal) => {
            const initialized = initializeFormData(newVal);
            formData.value = { ...initialized };
        }, { deep: true, immediate: true });

        watch(() => formData.value.type, (newType) => {
            formData.value.category_id = null;
            formData.value.keyword_name = '';
            formData.value.payment_method_id = null;
            if (newType === 'out') {
                formData.value.deposit_path = '';
            }
        });

        // Lifecycle
        onMounted(() => {
            loadInitialData();
            focusCurrentStepInput();
        });

        onUnmounted(() => {
            Object.values(autoAdvanceTimeouts.value).forEach(timeout => clearTimeout(timeout));
        });

        return {
            // State
            currentStep,
            saving,
            errors,
            formData,
            
            // Computed
            isMobile,
            steps,
            currentStepData,
            currentStepComponent,
            canProceed,
            canSkipCurrentStep,
            isFormValid,
            
            // Methods
            handleNext,
            handleBack,
            handleSkip,
            handleAutoAdvance,
            jumpToStep,
            handleValidate,
            handleSave,
            handleClose,
            openCategoryManager,
            openPaymentMethodManager,
            openDepositPathManager,
            openUserManager,
            openKeywordManager
        };
    }
}
</script>

<style scoped>
.add-data-page {
    @apply min-h-screen bg-gray-50 flex flex-col;
}

.container-responsive {
    @apply max-w-6xl mx-auto px-4;
}

.page-header {
    @apply bg-white border-b border-gray-200 sticky top-0 z-10;
}

.back-btn {
    @apply p-2 rounded-lg hover:bg-gray-100 transition-colors duration-200;
}

.close-btn {
    @apply p-2 rounded-lg hover:bg-gray-100 transition-colors duration-200 text-gray-500;
}

.type-selection-bar {
    @apply bg-white border-b border-gray-100;
}

.type-btn {
    @apply flex items-center justify-center p-3 border-2 border-gray-200 rounded-lg cursor-pointer transition-all duration-200;
}

.type-btn:hover {
    @apply border-gray-300 bg-gray-50;
}

.type-btn.selected.expense {
    @apply border-red-500 bg-red-50 text-red-700;
}

.type-btn.selected.income {
    @apply border-green-500 bg-green-50 text-green-700;
}

.stepper-section {
    @apply bg-white border-b border-gray-100 py-4;
}

.main-content {
    @apply flex-1 py-6;
}

.desktop-layout {
    @apply grid grid-cols-3 gap-6;
}

.step-content-area {
    @apply col-span-2 bg-white rounded-xl p-6 shadow-sm;
}

.summary-area {
    @apply col-span-1 bg-white rounded-xl shadow-sm;
}

.mobile-layout {
    @apply bg-white rounded-xl p-4 shadow-sm;
}

.action-section {
    @apply bg-white border-t border-gray-200 py-4 sticky bottom-0;
}

/* 모바일 최적화 */
@media (max-width: 768px) {
    .container-responsive {
        @apply px-3;
    }
    
    .main-content {
        @apply py-4;
    }
    
    .mobile-layout {
        @apply p-3 rounded-lg;
    }
}
</style>