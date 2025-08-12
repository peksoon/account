<template>
    <div class="modal-backdrop" @click.self="closeAddPopup">
        <div class="modal-content" :class="{ 'mobile-modal': isMobile }">
            <!-- 헤더 -->
            <div class="flex items-center justify-between p-6 border-b border-gray-200">
                <div class="flex items-center">
                    <div
                        class="w-10 h-10 bg-gradient-to-r from-yellow-500 to-yellow-600 rounded-lg flex items-center justify-center mr-3">
                        <DollarSign class="w-6 h-6 text-white" />
                    </div>
                    <div>
                        <h3 class="text-xl font-bold text-gray-900">새 데이터 추가</h3>
                        <p class="text-sm text-gray-500">수입 또는 지출 내역을 추가해보세요</p>
                    </div>
                </div>
                <el-button @click="closeAddPopup" size="large" circle>×</el-button>
            </div>

            <!-- 폼 내용 -->
            <div class="p-6">
                <el-form ref="formRef" :model="localAccount" :rules="rules" label-position="top"
                    @submit.prevent="saveAccount">
                    <!-- 타입 선택 -->
                    <div class="mb-6">
                        <label class="form-label">거래 유형</label>
                        <div class="grid grid-cols-2 gap-3">
                            <div @click="localAccount.type = 'out'"
                                class="p-4 border-2 rounded-xl cursor-pointer transition-all duration-200" :class="localAccount.type === 'out'
                                    ? 'border-red-500 bg-red-50'
                                    : 'border-gray-200 hover:border-red-300'">
                                <div class="flex items-center">
                                    <div class="w-8 h-8 bg-red-500 rounded-lg flex items-center justify-center mr-3">
                                        <TrendingDown class="w-5 h-5 text-white" />
                                    </div>
                                    <div>
                                        <p class="font-semibold text-gray-900">지출</p>
                                        <p class="text-sm text-gray-500">돈이 나간 내역</p>
                                    </div>
                                </div>
                            </div>

                            <div @click="localAccount.type = 'in'"
                                class="p-4 border-2 rounded-xl cursor-pointer transition-all duration-200" :class="localAccount.type === 'in'
                                    ? 'border-green-500 bg-green-50'
                                    : 'border-gray-200 hover:border-green-300'">
                                <div class="flex items-center">
                                    <div class="w-8 h-8 bg-green-500 rounded-lg flex items-center justify-center mr-3">
                                        <TrendingUp class="w-5 h-5 text-white" />
                                    </div>
                                    <div>
                                        <p class="font-semibold text-gray-900">수입</p>
                                        <p class="text-sm text-gray-500">돈이 들어온 내역</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <!-- 사용자 -->
                        <el-form-item label="사용자" prop="user">
                            <el-select v-model="localAccount.user" placeholder="사용자를 선택하세요" size="large" class="w-full"
                                filterable allow-create>
                                <el-option v-for="user in userOptions" :key="user.id" :label="user.label"
                                    :value="user.value" />
                            </el-select>
                            <div class="mt-2 flex items-center space-x-2">
                                <el-button size="small" text @click="openUserManager">
                                    사용자 관리
                                </el-button>
                            </div>
                        </el-form-item>

                        <!-- 금액 -->
                        <el-form-item label="금액" prop="money">
                            <el-input v-model.number="localAccount.money" type="number" placeholder="금액을 입력하세요"
                                size="large">
                                <template #suffix>원</template>
                            </el-input>
                        </el-form-item>
                    </div>

                    <!-- 카테고리 -->
                    <el-form-item label="카테고리" prop="category_id" class="mb-6">
                        <el-select v-model="localAccount.category_id" placeholder="카테고리를 선택하세요" size="large"
                            class="w-full" filterable @change="handleCategoryChange">
                            <el-option v-for="category in availableCategories" :key="category.id" :label="category.name"
                                :value="category.id" />
                        </el-select>
                        <div class="mt-2 flex items-center space-x-2">
                            <el-button size="small" text @click="openCategoryManager">
                                카테고리 관리
                            </el-button>
                        </div>
                    </el-form-item>

                    <!-- 지출 전용 필드 -->
                    <div v-if="localAccount.type === 'out'" class="space-y-6">
                        <!-- 키워드 (자동완성) -->
                        <el-form-item label="키워드 (선택사항)">
                            <KeywordAutocomplete v-model="localAccount.keyword_name"
                                :category-id="localAccount.category_id" placeholder="키워드를 입력하세요 (예: 식당명, 상품명 등)"
                                size="large" show-manage-button @select="handleKeywordSelect"
                                @open-keyword-manager="openKeywordManager" />
                        </el-form-item>

                        <!-- 결제 방법 -->
                        <el-form-item label="결제 방법" prop="payment_method_id">
                            <el-select v-model="localAccount.payment_method_id" placeholder="결제 방법을 선택하세요" size="large"
                                class="w-full" filterable>
                                <el-option-group v-for="category in paymentMethodStore.activePaymentMethods"
                                    :key="category.id" :label="category.name">
                                    <el-option v-for="method in category.children" :key="method.id" :label="method.name"
                                        :value="method.id">
                                        <div class="flex items-center">
                                            <CreditCard class="w-4 h-4 mr-2 text-purple-600" />
                                            {{ method.name }}
                                        </div>
                                    </el-option>
                                </el-option-group>
                            </el-select>
                            <div class="mt-2 flex items-center space-x-2">
                                <el-button size="small" text @click="openPaymentMethodManager">
                                    결제수단 관리
                                </el-button>
                            </div>
                        </el-form-item>
                    </div>

                    <!-- 수입 전용 필드 -->
                    <div v-if="localAccount.type === 'in'" class="mb-6 space-y-4">
                        <!-- 키워드 (수입용) -->
                        <el-form-item label="키워드 (선택사항)">
                            <KeywordAutocomplete v-model="localAccount.keyword_name"
                                :category-id="localAccount.category_id" placeholder="키워드를 입력하세요 (예: 회사명, 수입원 등)"
                                size="large" show-manage-button @select="handleKeywordSelect"
                                @open-keyword-manager="openKeywordManager" />
                        </el-form-item>

                        <!-- 입금경로 -->
                        <el-form-item label="입금경로" prop="deposit_path">
                            <el-select v-model="localAccount.deposit_path" placeholder="입금경로를 선택하세요" size="large"
                                class="w-full" filterable allow-create>
                                <el-option v-for="path in depositPathOptions" :key="path.id" :label="path.label"
                                    :value="path.value" />
                            </el-select>
                            <div class="mt-2 flex items-center space-x-2">
                                <el-button size="small" text @click="openDepositPathManager">
                                    입금경로 관리
                                </el-button>
                            </div>
                        </el-form-item>
                    </div>

                    <!-- 메모 -->
                    <el-form-item label="메모" class="mb-6">
                        <el-input v-model="localAccount.memo" type="textarea" :rows="3" placeholder="추가 메모를 입력하세요"
                            resize="none" />
                    </el-form-item>

                    <!-- 날짜 -->
                    <el-form-item label="날짜" prop="date" class="mb-6">
                        <el-date-picker v-model="localAccount.date" type="date" placeholder="날짜를 선택하세요"
                            format="YYYY-MM-DD" value-format="YYYY-MM-DD" size="large" class="w-full" />
                    </el-form-item>
                </el-form>
            </div>

            <!-- 하단 버튼 -->
            <div class="flex justify-end space-x-3 p-6 border-t border-gray-200 bg-gray-50">
                <el-button @click="closeAddPopup" size="large">
                    취소
                </el-button>
                <el-button type="primary" @click="saveAccount" size="large" :loading="saving" :disabled="!isFormValid">
                    <Save class="w-4 h-4 mr-2" />
                    저장하기
                </el-button>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, watch, computed, onMounted } from 'vue';
import {
    Save,
    TrendingUp,
    TrendingDown,
    CreditCard,
    DollarSign
} from 'lucide-vue-next';
import { ElMessage } from 'element-plus';
import KeywordAutocomplete from './KeywordAutocomplete.vue';
import { useCategoryStore } from '../stores/categoryStore';
import { usePaymentMethodStore } from '../stores/paymentMethodStore';
import { useDepositPathStore } from '../stores/depositPathStore';
import { useKeywordStore } from '../stores/keywordStore';
import { useUserStore } from '../stores/userStore';
import { useBudgetStore } from '../stores/budgetStore';

export default {
    components: {
        Save,
        TrendingUp,
        TrendingDown,
        CreditCard,
        DollarSign,
        KeywordAutocomplete
    },
    props: {
        newAccount: {
            type: Object,
            required: true
        },
        selectedDate: {
            type: String,
            default: null
        }
    },
    emits: ['save', 'close', 'update:newAccount', 'open-category-manager', 'open-payment-method-manager', 'open-deposit-path-manager', 'open-keyword-manager', 'open-user-manager', 'budget-alert', 'budget-save-success'],
    setup(props, { emit }) {
        const categoryStore = useCategoryStore();
        const paymentMethodStore = usePaymentMethodStore();
        const depositPathStore = useDepositPathStore();
        const keywordStore = useKeywordStore();
        const userStore = useUserStore();
        const budgetStore = useBudgetStore();

        const formRef = ref(null);
        const saving = ref(false);

        // 로컬 계정 데이터 초기화 함수
        const initializeLocalAccount = (accountData) => {
            const defaultDate = new Date().toISOString().slice(0, 10);

            console.log('initializeLocalAccount - 입력:', accountData);
            console.log('initializeLocalAccount - props.selectedDate:', props.selectedDate);
            console.log('initializeLocalAccount - accountData.date:', accountData?.date);

            // 우선순위: props.selectedDate > accountData.date > defaultDate
            let dateToUse = defaultDate;
            if (props.selectedDate && props.selectedDate !== '') {
                dateToUse = props.selectedDate;
            } else if (accountData && accountData.date && accountData.date !== '') {
                dateToUse = accountData.date;
            }

            console.log('initializeLocalAccount - 사용할 날짜:', dateToUse);

            const result = {
                type: 'out',
                user: '',
                money: '',
                category_id: null,
                keyword_name: '',
                payment_method_id: null,
                memo: '',
                deposit_path: '',
                ...accountData, // props 값으로 덮어쓰기
                date: dateToUse, // date는 마지막에 확실히 설정
            };

            console.log('initializeLocalAccount - 최종 결과:', result);

            return result;
        };

        // 로컬 계정 데이터
        const localAccount = ref(initializeLocalAccount(props.newAccount));

        // 모바일 체크
        const isMobile = computed(() => {
            if (typeof window === 'undefined') return false;
            return window.innerWidth < 768;
        });

        // 폼 검증 규칙
        const rules = {
            user: [
                { required: true, message: '사용자명을 입력해주세요', trigger: 'blur' }
            ],
            money: [
                { required: true, message: '금액을 입력해주세요', trigger: 'blur' },
                { type: 'number', min: 1, message: '금액은 1원 이상이어야 합니다', trigger: 'blur' }
            ],
            category_id: [
                { required: true, message: '카테고리를 선택해주세요', trigger: 'change' }
            ],
            payment_method_id: [
                { required: true, message: '결제수단을 선택해주세요', trigger: 'change' }
            ],
            date: [
                { required: true, message: '날짜를 선택해주세요', trigger: 'change' }
            ]
        };

        // 사용 가능한 카테고리 목록
        const availableCategories = computed(() => {
            return categoryStore.categories.filter(cat => cat.type === localAccount.value.type);
        });

        // 사용자 옵션 목록
        const userOptions = computed(() => {
            return userStore.getUserOptions();
        });

        // 결제수단 관련 computed는 store getter를 직접 사용

        // 입금경로 옵션 생성
        const depositPathOptions = computed(() => {
            return depositPathStore.activeDepositPaths.map(path => ({
                id: path.id,
                label: path.name,
                value: path.name
            }));
        });

        // 폼 유효성 검사
        const isFormValid = computed(() => {
            return localAccount.value.user &&
                localAccount.value.money > 0 &&
                localAccount.value.category_id &&
                (localAccount.value.type === 'in' ? localAccount.value.deposit_path : localAccount.value.payment_method_id) &&
                localAccount.value.date;
        });

        // 결제수단 아이콘 매핑
        // 결제 아이콘은 템플릿에서 CreditCard로 통일

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

        // props 변경 감지
        watch(() => props.newAccount, (newVal) => {
            console.log('AddPopup props.newAccount 변경됨:', newVal);
            const initialized = initializeLocalAccount(newVal);
            localAccount.value = { ...initialized };
            console.log('AddPopup localAccount 초기화 후:', localAccount.value);
        }, { deep: true, immediate: true });

        // 날짜가 변경될 때 추가 로그
        watch(() => localAccount.value.date, (newDate) => {
            console.log('AddPopup localAccount.date 변경됨:', newDate);
        });

        // 타입 변경 시 관련 필드 초기화
        watch(() => localAccount.value.type, (newType) => {
            localAccount.value.category_id = null;
            localAccount.value.keyword_name = '';
            localAccount.value.payment_method_id = null;
            if (newType === 'out') {
                localAccount.value.account_number = '';
            }
        });

        // 카테고리 변경 핸들러
        const handleCategoryChange = () => {
            // 카테고리 변경 시 키워드 초기화
            localAccount.value.keyword_name = '';
        };

        // 키워드 선택 핸들러
        const handleKeywordSelect = (keyword) => {
            localAccount.value.keyword_name = keyword.name;
        };

        // 계정 저장
        const saveAccount = async () => {
            if (!formRef.value) return;

            try {
                // 폼 검증
                await formRef.value.validate();

                saving.value = true;

                // 데이터 정리
                const accountData = { ...localAccount.value };

                // 지출인 경우 기준치 정보를 포함한 API 사용 (키워드 처리는 백엔드에서)
                if (localAccount.value.type === 'out') {
                    try {
                        const response = await budgetStore.createOutAccountWithBudget(accountData);

                        console.log('기준치 포함 지출 응답:', response);
                        console.log('기준치 사용량 데이터:', response.budget_usage);

                        // 기준치 정보가 있는 경우 알림 팝업 표시
                        if (response.budget_usage) {
                            console.log('기준치 경고 팝업 표시 중...');
                            emit('budget-alert', {
                                budgetUsage: response.budget_usage,
                                expenseAmount: accountData.money,
                                expenseDate: accountData.date,
                                expenseKeyword: accountData.keyword_name
                            });
                        } else {
                            console.log('기준치 정보가 없어서 팝업 닫기');
                            // 기준치 정보가 없으면 팝업 닫기
                            emit('close');
                        }

                        // 캘린더 갱신을 위해 save 이벤트 발생
                        emit('budget-save-success');

                        ElMessage.success('지출이 성공적으로 기록되었습니다.');
                    } catch (error) {
                        console.error('기준치 포함 지출 저장 실패:', error);
                        ElMessage.error('지출 저장 중 오류가 발생했습니다.');
                        return;
                    }
                } else {
                    // 수입인 경우 키워드 사용 기록 후 기존 방식 사용
                    if (localAccount.value.keyword_name && localAccount.value.category_id) {
                        await keywordStore.useKeyword(localAccount.value.category_id, localAccount.value.keyword_name);
                    }

                    emit('save', accountData);
                    emit('close');
                    ElMessage.success('수입이 성공적으로 기록되었습니다.');
                }

            } catch (error) {
                console.error('Form validation failed:', error);
                ElMessage.error('입력 정보를 확인해주세요.');
            } finally {
                saving.value = false;
            }
        };

        // 팝업 닫기
        const closeAddPopup = () => {
            // 폼 초기화
            if (formRef.value) {
                formRef.value.resetFields();
            }
            emit('close');
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

        // 금액 포맷팅
        const formatMoney = (amount) => {
            if (!amount) return '';
            return new Intl.NumberFormat('ko-KR').format(amount);
        };

        // 컴포넌트 마운트 시 초기 데이터 로드
        onMounted(() => {
            loadInitialData();
        });

        return {
            formRef,
            saving,
            localAccount,
            rules,
            availableCategories,
            depositPathOptions,
            userOptions,
            isFormValid,
            isMobile,

            saveAccount,
            closeAddPopup,
            handleCategoryChange,
            handleKeywordSelect,
            openCategoryManager,
            openPaymentMethodManager,
            openDepositPathManager,
            openKeywordManager,
            openUserManager,
            formatMoney,

            // Stores
            paymentMethodStore,

            // 아이콘들
            Save,
            TrendingUp,
            TrendingDown,
            CreditCard,
            DollarSign
        };
    }
}
</script>

<style scoped>
/* 커스텀 입력 필드 스타일 */
:deep(.el-form-item__label) {
    @apply font-semibold text-gray-700 text-sm;
}

:deep(.el-input__wrapper) {
    @apply border-gray-200 rounded-lg;
    transition: all 0.2s ease;
}

:deep(.el-input__wrapper:hover) {
    @apply border-gray-300;
}

:deep(.el-input__wrapper.is-focus) {
    @apply border-blue-500 ring-2 ring-blue-100;
}

:deep(.el-select .el-input__wrapper) {
    @apply border-gray-200 rounded-lg;
}

:deep(.el-textarea__inner) {
    @apply border-gray-200 rounded-lg;
    resize: none;
}

:deep(.el-date-editor.el-input) {
    @apply w-full;
}

/* 버튼 커스텀 스타일 */
:deep(.el-button--primary) {
    @apply bg-gradient-to-r from-blue-500 to-blue-600 border-blue-500;
}

:deep(.el-button--primary:hover) {
    @apply from-blue-600 to-blue-700 border-blue-600;
}

/* 모달 애니메이션 */
.modal-backdrop {
    animation: fadeIn 0.3s ease-out;
}

.modal-content {
    animation: slideIn 0.3s ease-out;
}

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

/* 타입 선택 카드 애니메이션 */
.transition-all {
    transition: all 0.2s ease;
}

/* 모바일 최적화 */
.mobile-modal {
    position: fixed !important;
    top: 0 !important;
    left: 0 !important;
    right: 0 !important;
    bottom: 0 !important;
    width: 100vw !important;
    height: 100vh !important;
    max-width: none !important;
    max-height: none !important;
    margin: 0 !important;
    border-radius: 0 !important;
    display: flex !important;
    flex-direction: column !important;
    overflow: hidden !important;
}

.mobile-modal .p-6 {
    @apply p-4;
}

/* 모바일에서 폼 내용 영역 스크롤 처리 */
.mobile-modal .p-6:nth-child(2) {
    flex: 1;
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
}

/* 모바일에서 하단 버튼 영역 고정 */
.mobile-modal .border-t {
    flex-shrink: 0;
}

@media (max-width: 768px) {
    .modal-content:not(.mobile-modal) {
        @apply max-w-none w-full mx-4 max-h-[90vh] overflow-y-auto;
    }

    .grid.grid-cols-2 {
        @apply grid-cols-1;
    }

    .grid.md\\:grid-cols-2 {
        @apply grid-cols-1;
    }
}

/* 폼 레이블 스타일 */
.form-label {
    @apply block text-sm font-semibold text-gray-700 mb-3;
}

/* 호버 효과 */
.cursor-pointer:hover {
    transform: translateY(-1px);
}

/* 저장 버튼 로딩 상태 */
:deep(.el-button.is-loading) {
    @apply pointer-events-none;
}

/* 아이콘 색상 커스텀 */
.payment-method-icon {
    color: inherit !important;
}

.icon-card {
    color: #3b82f6 !important;
}

.icon-transfer {
    color: #10b981 !important;
}

.icon-cash {
    color: #f59e0b !important;
}

.icon-other {
    color: #6b7280 !important;
}

/* DollarSign 아이콘 색상 강조 */
.lucide-dollar-sign {
    color: #ffffff !important;
    filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
}

/* 모바일에서 Element Plus 드롭다운 z-index 수정 */
@media (max-width: 768px) {
    :deep(.el-select-dropdown) {
        z-index: 9999 !important;
    }

    :deep(.el-popper) {
        z-index: 9999 !important;
    }

    :deep(.el-date-picker__popper) {
        z-index: 9999 !important;
    }

    /* 모바일에서 select 옵션이 잘 보이도록 */
    :deep(.el-select-dropdown__item) {
        padding: 12px 20px !important;
        font-size: 16px !important;
        line-height: 1.5 !important;
    }
}
</style>