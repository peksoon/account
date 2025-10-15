<template>
    <div class="full-page-container" :class="{ 'mobile-container': isMobile }">
        <div class="page-content" :class="{ 'mobile-content': isMobile }">
            <!-- 헤더 -->
            <div class="flex items-center justify-between p-6 border-b border-gray-200">
                <div class="flex items-center">
                    <button @click="goBack" class="mr-4 p-2 hover:bg-gray-100 rounded-lg transition-colors">
                        <ArrowLeft class="w-5 h-5 text-gray-600" />
                    </button>
                    <div
                        class="w-10 h-10 bg-gradient-to-r from-yellow-500 to-yellow-600 rounded-lg flex items-center justify-center mr-3">
                        <DollarSign class="w-6 h-6 text-white" />
                    </div>
                    <div>
                        <h3 class="text-xl font-bold text-gray-900">새 데이터 추가</h3>
                        <p class="text-sm text-gray-500">단계별로 정보를 입력해보세요</p>
                    </div>
                </div>
                <el-button @click="goBack" size="large" circle>×</el-button>
            </div>

            <!-- 타입 선택 (고정) -->
            <div class="p-6 border-b border-gray-100">
                <div class="grid grid-cols-2 gap-3">
                    <div @click="formData.type = 'out'"
                        class="p-3 border-2 rounded-xl cursor-pointer transition-all duration-200"
                        :class="formData.type === 'out' ? 'border-red-500 bg-red-50' : 'border-gray-200 hover:border-red-300'">
                        <div class="flex items-center justify-center">
                            <div class="w-6 h-6 bg-red-500 rounded-lg flex items-center justify-center mr-2">
                                <TrendingDown class="w-4 h-4 text-white" />
                            </div>
                            <span class="font-semibold text-sm">지출</span>
                        </div>
                    </div>
                    <div @click="formData.type = 'in'"
                        class="p-3 border-2 rounded-xl cursor-pointer transition-all duration-200"
                        :class="formData.type === 'in' ? 'border-green-500 bg-green-50' : 'border-gray-200 hover:border-green-300'">
                        <div class="flex items-center justify-center">
                            <div class="w-6 h-6 bg-green-500 rounded-lg flex items-center justify-center mr-2">
                                <TrendingUp class="w-4 h-4 text-white" />
                            </div>
                            <span class="font-semibold text-sm">수입</span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 스테퍼 -->
            <div class="stepper-container">
                <div class="stepper" :class="{ 'mobile-stepper': isMobile }">
                    <div v-for="(step, index) in steps" :key="step.id" class="step-item" :class="{
                        'active': currentStep === index,
                        'completed': index < currentStep,
                        'clickable': canJumpToStep(index)
                    }" @click="jumpToStep(index)">
                        <div class="step-number">
                            <CheckCircle2 v-if="index < currentStep" class="w-4 h-4" />
                            <span v-else>{{ index + 1 }}</span>
                        </div>
                        <div class="step-label" v-if="!isMobile">{{ step.label }}</div>
                        <div class="step-progress" v-if="index < steps.length - 1"></div>
                    </div>
                </div>
                <div class="step-counter" v-if="isMobile">
                    {{ currentStep + 1 }} / {{ steps.length }} {{ currentStepData.label }}
                </div>
            </div>

            <!-- 메인 컨텐츠 -->
            <div class="wizard-content">
                <!-- Desktop 레이아웃 -->
                <div class="wizard-layout" v-if="!isMobile">
                    <div class="step-content">
                        <component :is="currentStepComponent" v-model="formData" :errors="errors" @next="handleNext"
                            @auto-advance="handleAutoAdvance" @validate="handleValidate" />
                    </div>
                    <div class="summary-panel">
                        <SummaryPanel :data="formData" :current-step="currentStep" @jump-to-step="jumpToStep" />
                    </div>
                </div>

                <!-- Mobile 레이아웃 -->
                <div class="mobile-step-content" v-else>
                    <component :is="currentStepComponent" v-model="formData" :errors="errors" @next="handleNext"
                        @auto-advance="handleAutoAdvance" @validate="handleValidate" />
                </div>
            </div>

            <!-- 하단 액션 바 -->
            <div class="action-bar" :class="{ 'mobile-action-bar': isMobile }">
                <div class="flex justify-between items-center">
                    <el-button v-if="currentStep > 0" @click="handleBack" size="large" class="back-button">
                        <ArrowLeft class="w-4 h-4 mr-2" />
                        이전
                    </el-button>
                    <div v-else></div>

                    <div class="flex space-x-3">
                        <el-button v-if="canSkipCurrentStep" @click="handleSkip" size="large">
                            건너뛰기
                        </el-button>

                        <el-button v-if="currentStep < steps.length - 1" @click="handleNext" type="primary" size="large"
                            :disabled="!canProceed">
                            다음
                            <ArrowRight class="w-4 h-4 ml-2" />
                        </el-button>

                        <el-button v-else @click="handleSave" type="primary" size="large" :loading="saving"
                            :disabled="!isFormValid">
                            <Save class="w-4 h-4 mr-2" />
                            저장하기
                        </el-button>
                    </div>
                </div>
            </div>
        </div>

        <!-- 기준치 알림 팝업 -->
        <BudgetAlertPopup :is-visible="showBudgetAlert" :budget-usage="budgetAlertData.budgetUsage"
            :expense-amount="budgetAlertData.expenseAmount" :expense-date="budgetAlertData.expenseDate"
            :expense-keyword="budgetAlertData.expenseKeyword" @close="closeBudgetAlert"
            @open-budget-management="openBudgetManager" />

    </div>
</template>

<script>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import {
    DollarSign,
    TrendingUp,
    TrendingDown,
    ArrowLeft,
    ArrowRight,
    Save,
    CheckCircle2
} from 'lucide-vue-next';

// Step Components
import AmountStep from './wizard-steps/AmountStep.vue';
import UserStep from './wizard-steps/UserStep.vue';
import CategoryStep from './wizard-steps/CategoryStep.vue';
import KeywordStep from './wizard-steps/KeywordStep.vue';
import PaymentMethodStep from './wizard-steps/PaymentMethodStep.vue';
import MemoStep from './wizard-steps/MemoStep.vue';
import DateStep from './wizard-steps/DateStep.vue';
import ReviewStep from './wizard-steps/ReviewStep.vue';
import SummaryPanel from './wizard-steps/SummaryPanel.vue';
import BudgetAlertPopup from './BudgetAlertPopup.vue';
import CategoryManager from './CategoryManager.vue';
import PaymentMethodManager from './PaymentMethodManager.vue';
import DepositPathManager from './DepositPathManager.vue';
import KeywordManager from './KeywordManager.vue';
import UserManager from './UserManager.vue';

// Stores
import { useCategoryStore } from '../stores/categoryStore';
import { usePaymentMethodStore } from '../stores/paymentMethodStore';
import { useDepositPathStore } from '../stores/depositPathStore';
import { useKeywordStore } from '../stores/keywordStore';
import { useUserStore } from '../stores/userStore';
import { useBudgetStore } from '../stores/budgetStore';
import { useAccountStore } from '../stores/accountStore';
import { getTodayKST } from '../utils';

export default {
    name: 'AddDataWizard',
    components: {
        DollarSign,
        TrendingUp,
        TrendingDown,
        ArrowLeft,
        ArrowRight,
        Save,
        CheckCircle2,
        AmountStep,
        UserStep,
        CategoryStep,
        KeywordStep,
        PaymentMethodStep,
        MemoStep,
        DateStep,
        ReviewStep,
        SummaryPanel,
        BudgetAlertPopup,
        CategoryManager,
        PaymentMethodManager,
        DepositPathManager,
        KeywordManager,
        UserManager
    },
    props: {
        newAccount: {
            type: Object,
            default: () => ({})
        },
        selectedDate: {
            type: String,
            default: null
        }
    },
    emits: ['save', 'close', 'update:newAccount', 'open-category-manager', 'open-payment-method-manager', 'open-deposit-path-manager', 'open-keyword-manager', 'open-user-manager', 'budget-alert', 'budget-save-success'],
    setup(props, { emit }) {
        const router = useRouter();
        const route = useRoute();
        // Stores
        const categoryStore = useCategoryStore();
        const paymentMethodStore = usePaymentMethodStore();
        const depositPathStore = useDepositPathStore();
        const keywordStore = useKeywordStore();
        const userStore = useUserStore();
        const budgetStore = useBudgetStore();
        const accountStore = useAccountStore();

        // State
        const currentStep = ref(0);
        const saving = ref(false);
        const errors = ref({});
        const autoAdvanceTimeouts = ref({});

        // 기준치 알림 상태
        const showBudgetAlert = ref(false);
        const budgetAlertData = ref({
            budgetUsage: null,
            expenseAmount: 0,
            expenseDate: '',
            expenseKeyword: ''
        });

        // 관리 모달 상태
        const showCategoryManager = ref(false);
        const showPaymentMethodManager = ref(false);
        const showDepositPathManager = ref(false);
        const showKeywordManager = ref(false);
        const showUserManager = ref(false);
        const keywordManagerCategoryId = ref(null);

        // 뒤로 가기 함수
        const goBack = () => {
            router.push('/');
        };

        // 폼 데이터 초기화 - route 파라미터 사용
        const initializeFormData = (accountData) => {
            const defaultDate = getTodayKST();
            let dateToUse = defaultDate;

            // Route 쿼리에서 날짜 가져오기
            if (route.query.date) {
                dateToUse = route.query.date;
            } else if (props.selectedDate && props.selectedDate !== '') {
                dateToUse = props.selectedDate;
            } else if (accountData && accountData.date && accountData.date !== '') {
                dateToUse = accountData.date;
            }

            // Route 쿼리에서 초기 데이터 가져오기
            let initialData = {};
            if (route.query.data) {
                try {
                    initialData = JSON.parse(route.query.data);
                } catch (error) {
                    console.warn('Failed to parse route data:', error);
                }
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
                ...accountData,
                ...initialData,
                date: dateToUse,
            };
        };

        const formData = ref(initializeFormData(props.newAccount));

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

        // 현재 스텝 건너뛰기 가능 여부
        const canSkipCurrentStep = computed(() => {
            return !currentStepData.value.required;
        });

        // 특정 스텝으로 점프 가능 여부
        const canJumpToStep = (stepIndex) => {
            return stepIndex <= currentStep.value;
        };

        // 폼 전체 유효성
        const isFormValid = computed(() => {
            return formData.value.user &&
                formData.value.money > 0 &&
                formData.value.category_id &&
                formData.value.keyword_name &&
                (formData.value.type === 'in' ? formData.value.deposit_path : formData.value.payment_method_id) &&
                formData.value.date;
        });

        // 키보드 이벤트 핸들러
        const handleKeyboardNavigation = (e) => {
            // textarea, input 등에서 발생한 이벤트는 무시
            if (e.target.tagName === 'TEXTAREA' || e.target.tagName === 'INPUT') {
                return;
            }

            if (e.altKey) {
                if (e.key === 'ArrowRight') {
                    e.preventDefault();
                    handleNext();
                } else if (e.key === 'ArrowLeft') {
                    e.preventDefault();
                    handleBack();
                }
            } else if (e.key === 'Enter') {
                e.preventDefault();
                if (currentStep.value < steps.value.length - 1) {
                    handleNext();
                } else {
                    handleSave();
                }
            } else if (e.key === 'Escape') {
                e.preventDefault();
                handleClose();
            }
        };

        // 이벤트 핸들러들
        const handleNext = async () => {
            if (!canProceed.value) return;

            clearAutoAdvanceTimeout();

            if (currentStep.value < steps.value.length - 1) {
                currentStep.value++;
                announceStepChange();
                await nextTick();
                focusCurrentStepInput();
            }
        };

        const handleBack = () => {
            if (currentStep.value > 0) {
                clearAutoAdvanceTimeout();
                currentStep.value--;
                announceStepChange();
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
            if (canJumpToStep(stepIndex)) {
                clearAutoAdvanceTimeout();
                currentStep.value = stepIndex;
                announceStepChange();
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
                        // 기준치 알림 표시
                        budgetAlertData.value = {
                            budgetUsage: response.budget_usage,
                            expenseAmount: accountData.money,
                            expenseDate: accountData.date,
                            expenseKeyword: accountData.keyword_name
                        };
                        showBudgetAlert.value = true;
                    } else {
                        goBack();
                    }

                    emit('budget-save-success');
                    ElMessage.success('지출이 성공적으로 기록되었습니다.');
                } else {
                    // 수입 데이터를 accountStore를 통해 저장
                    await accountStore.saveAccount(accountData);

                    if (formData.value.keyword_name && formData.value.category_id) {
                        await keywordStore.useKeyword(formData.value.category_id, formData.value.keyword_name);
                    }

                    emit('save', accountData);
                    goBack();
                    ElMessage.success('수입이 성공적으로 기록되었습니다.');
                }
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
                        '저장되지 않은 변경사항이 있습니다. 정말로 닫으시겠습니까?',
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

            goBack();
        };

        const handleOverlayClick = () => {
            handleClose();
        };

        // 관리자 모달 열기
        const openCategoryManager = () => {
            showCategoryManager.value = true;
        };

        const openPaymentMethodManager = () => {
            showPaymentMethodManager.value = true;
        };

        const openDepositPathManager = () => {
            showDepositPathManager.value = true;
        };

        const openUserManager = () => {
            showUserManager.value = true;
        };

        // 관리자 모달 닫기
        const closeCategoryManager = () => {
            showCategoryManager.value = false;
        };

        const closePaymentMethodManager = () => {
            showPaymentMethodManager.value = false;
        };

        const closeDepositPathManager = () => {
            showDepositPathManager.value = false;
        };

        const closeUserManager = () => {
            showUserManager.value = false;
        };

        const closeKeywordManager = () => {
            showKeywordManager.value = false;
            keywordManagerCategoryId.value = null;
        };

        // 기준치 알림 관련 메서드
        const closeBudgetAlert = () => {
            showBudgetAlert.value = false;
            budgetAlertData.value = {
                budgetUsage: null,
                expenseAmount: 0,
                expenseDate: '',
                expenseKeyword: ''
            };
            // 알림 닫으면 페이지에서 나가기
            goBack();
        };

        const openBudgetManager = () => {
            // 기준치 관리는 메인 페이지에서 처리
            closeBudgetAlert();
        };

        const openKeywordManager = (categoryId) => {
            keywordManagerCategoryId.value = categoryId;
            showKeywordManager.value = true;
        };

        // 접근성 관련
        const announceStepChange = () => {
            // 스크린 리더를 위한 live region 업데이트는 각 스텝 컴포넌트에서 처리
            // 추후 aria-live 영역에 메시지를 표시할 수 있음
        };

        const focusCurrentStepInput = () => {
            nextTick(() => {
                const firstInput = document.querySelector('.step-content input, .step-content .el-input__inner, .step-content .el-select');
                if (firstInput) {
                    firstInput.focus();
                }
            });
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
        watch(() => props.newAccount, (newVal) => {
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
            document.addEventListener('keydown', handleKeyboardNavigation);
            focusCurrentStepInput();
        });

        onUnmounted(() => {
            document.removeEventListener('keydown', handleKeyboardNavigation);
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
            canJumpToStep,
            isFormValid,

            // Methods
            goBack,
            handleNext,
            handleBack,
            handleSkip,
            handleAutoAdvance,
            jumpToStep,
            handleValidate,
            handleSave,
            handleClose,
            handleOverlayClick,
            openCategoryManager,
            openPaymentMethodManager,
            openDepositPathManager,
            openUserManager,
            openKeywordManager,

            // Budget Alert
            showBudgetAlert,
            budgetAlertData,
            closeBudgetAlert,
            openBudgetManager,

            // Management Modals
            showCategoryManager,
            showPaymentMethodManager,
            showDepositPathManager,
            showKeywordManager,
            showUserManager,
            keywordManagerCategoryId,
            closeCategoryManager,
            closePaymentMethodManager,
            closeDepositPathManager,
            closeKeywordManager,
            closeUserManager,

            // Icons
            DollarSign,
            TrendingUp,
            TrendingDown,
            ArrowLeft,
            ArrowRight,
            Save,
            CheckCircle2
        };
    }
}
</script>

<style scoped>
.full-page-container {
    @apply min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50;
    animation: fadeIn 0.3s ease-out;
}

.page-content {
    @apply bg-white shadow-lg w-full max-w-6xl mx-auto flex flex-col;
    animation: slideIn 0.3s ease-out;
}

.mobile-container {
    @apply min-h-screen;
}

.mobile-content {
    @apply w-full max-w-none rounded-none;
}

/* 스테퍼 스타일 */
.stepper-container {
    @apply p-4 border-b border-gray-100 bg-gray-50;
}

.stepper {
    @apply flex items-center justify-center space-x-4;
}

.mobile-stepper {
    @apply flex space-x-2 overflow-x-auto pb-2;
}

.step-item {
    @apply flex items-center relative cursor-pointer transition-all duration-200;
}

.step-item.clickable:hover {
    @apply transform scale-105;
}

.step-number {
    @apply w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold border-2 transition-all duration-200;
    position: relative;
    z-index: 2;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.step-item.active .step-number {
    @apply bg-blue-500 text-white border-blue-500;
}

.step-item.completed .step-number {
    @apply bg-green-500 text-white border-green-500;
}

.step-item:not(.active):not(.completed) .step-number {
    @apply bg-white text-gray-400 border-gray-300;
}

.step-label {
    @apply ml-2 text-sm font-medium;
    text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
    background: rgba(255, 255, 255, 0.9);
    padding: 2px 6px;
    border-radius: 4px;
    position: relative;
    z-index: 1;
}

.step-item.active .step-label {
    @apply text-blue-600;
}

.step-item.completed .step-label {
    @apply text-green-600;
}

.step-item:not(.active):not(.completed) .step-label {
    @apply text-gray-400;
}

.step-progress {
    @apply absolute left-10 w-12 h-0.5 bg-gray-300;
    z-index: -1;
    top: 50%;
    transform: translateY(-50%);
}

.step-item.completed .step-progress {
    @apply bg-green-500;
}

.step-counter {
    @apply text-center text-sm font-medium text-gray-600 mt-2;
    background: rgba(255, 255, 255, 0.9);
    padding: 4px 8px;
    border-radius: 6px;
    margin: 8px auto;
    display: inline-block;
    text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

/* 레이아웃 */
.wizard-content {
    @apply flex-1 overflow-hidden;
}

.wizard-layout {
    @apply flex h-full;
}

.step-content {
    @apply flex-1 p-6 overflow-y-auto;
}

.summary-panel {
    @apply w-80 border-l border-gray-200 bg-gray-50;
}

.mobile-step-content {
    @apply flex-1 p-4 overflow-y-auto;
    /* 모바일에서 더 나은 스크롤 경험 */
    -webkit-overflow-scrolling: touch;
}

/* 액션 바 */
.action-bar {
    @apply p-6 border-t border-gray-200 bg-gray-50 flex-shrink-0;
}

.mobile-action-bar {
    @apply p-4 sticky bottom-0 bg-white border-t border-gray-200;
    /* 모바일에서 액션 바가 잘 보이도록 */
    box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
}

.back-button {
    @apply text-gray-600 border-gray-300;
}

.back-button:hover {
    @apply text-gray-800 border-gray-400;
}

/* 애니메이션 */
@keyframes fadeIn {
    from {
        opacity: 0;
    }

    to {
        opacity: 1;
    }
}

@keyframes slideIn {
    from {
        opacity: 0;
        transform: translateY(-20px) scale(0.95);
    }

    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

/* 모바일 최적화 */
@media (max-width: 768px) {
    .step-item:not(.active) {
        @apply opacity-50;
    }

    .step-label {
        @apply hidden;
    }

    .step-progress {
        @apply hidden;
    }

    .stepper {
        @apply justify-start;
    }

    .step-header h2 {
        @apply text-lg;
    }

    .step-header p {
        @apply text-sm;
    }

    .modal-content {
        @apply text-sm;
    }

    .action-bar .el-button {
        @apply text-sm px-4 py-2;
        min-height: 44px;
    }

    /* 스테퍼 카운터 텍스트 크기 조정 */
    .step-counter {
        @apply text-sm;
    }
}

/* 접근성 */
.step-item:focus {
    @apply outline-none ring-2 ring-blue-500 ring-offset-2;
}

/* Element Plus 커스터마이징 */
:deep(.el-button--primary) {
    @apply bg-blue-500 border-blue-500;
}

:deep(.el-button--primary:hover) {
    @apply bg-blue-600 border-blue-600;
}
</style>