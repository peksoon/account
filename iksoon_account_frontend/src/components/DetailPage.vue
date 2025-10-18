<template>
    <div class="detail-page min-h-screen bg-gray-50">
        <!-- 헤더 -->
        <div class="bg-white shadow-sm border-b border-gray-200">
            <div class="max-w-4xl mx-auto px-4 py-4">
                <div class="flex items-center justify-between">
                    <div class="flex items-center">
                        <el-button @click="goBack" type="default" circle class="mr-4">
                            ←
                        </el-button>
                        <div class="flex items-center">
                            <div class="w-12 h-12 rounded-lg flex items-center justify-center mr-4" :class="localEventDetail.type === 'out'
                                ? 'bg-gradient-to-r from-red-500 to-red-600'
                                : 'bg-gradient-to-r from-green-500 to-green-600'">
                                <TrendingDown v-if="localEventDetail.type === 'out'" class="w-7 h-7 text-white" />
                                <TrendingUp v-else class="w-7 h-7 text-white" />
                            </div>
                            <div>
                                <h1 class="text-2xl font-bold text-gray-900">
                                    {{ isEditMode ? '데이터 편집' : (localEventDetail.type === 'out' ? '지출 상세' : '수입 상세') }}
                                </h1>
                                <p class="text-sm text-gray-500">
                                    {{ isEditMode ? '정보를 수정하세요' : '거래 내역을 확인하세요' }}
                                </p>
                            </div>
                        </div>
                    </div>
                    <el-button @click="goBack" size="large" circle>×</el-button>
                </div>
            </div>
        </div>

        <!-- 내용 -->
        <div class="max-w-4xl mx-auto px-4 py-8">
            <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-8">
                <!-- 편집 모드 -->
                <div v-if="isEditMode">
                    <el-form ref="formRef" :model="localEventDetail" :rules="rules" label-position="top">
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <!-- 사용자 -->
                            <el-form-item label="사용자" prop="user">
                                <el-input v-model="localEventDetail.user" placeholder="사용자명을 입력하세요" size="large" />
                            </el-form-item>

                            <!-- 금액 -->
                            <el-form-item label="금액" prop="money">
                                <el-input v-model.number="localEventDetail.money" type="number" placeholder="금액을 입력하세요"
                                    size="large">
                                    <template #suffix>원</template>
                                </el-input>
                            </el-form-item>
                        </div>

                        <!-- 카테고리 -->
                        <el-form-item label="카테고리" prop="category_name" class="mb-6">
                            <el-select v-model="localEventDetail.category_name" placeholder="카테고리를 선택하세요" size="large"
                                class="w-full" filterable allow-create :loading="categoryStore.loading">
                                <el-option v-for="category in availableCategories" :key="category" :label="category"
                                    :value="category" />
                            </el-select>
                            <div v-if="categoryStore.loading" class="mt-1 text-xs text-gray-500">
                                카테고리 목록을 불러오는 중...
                            </div>
                            <div v-else-if="availableCategories.length === 0" class="mt-1 text-xs text-orange-500">
                                {{ localEventDetail.type === 'out' ? '지출' : '수입' }} 카테고리가 없습니다.
                            </div>
                            <div v-else class="mt-1 text-xs text-gray-500">
                                {{ availableCategories.length }}개의 카테고리를 사용할 수 있습니다.
                            </div>
                        </el-form-item>

                        <!-- 지출 전용 필드 -->
                        <div v-if="localEventDetail.type === 'out'" class="space-y-6">
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                                <!-- 키워드 -->
                                <el-form-item label="키워드">
                                    <el-input v-model="localEventDetail.keyword_name" :prefix-icon="Tag"
                                        placeholder="키워드를 입력하세요" size="large" />
                                </el-form-item>

                                <!-- 결제 방법 -->
                                <el-form-item label="결제 방법">
                                    <el-select v-model="localEventDetail.payment_method_id" placeholder="결제 방법을 선택하세요"
                                        size="large" class="w-full">
                                        <el-option v-for="method in paymentMethodStore.flatPaymentMethods"
                                            :key="method.id" :label="method.name" :value="method.id" />
                                    </el-select>
                                </el-form-item>
                            </div>
                        </div>

                        <!-- 수입 전용 필드 -->
                        <div v-if="localEventDetail.type === 'in'" class="space-y-6">
                            <el-form-item label="입금 경로">
                                <el-select v-model="localEventDetail.deposit_path" placeholder="입금 경로를 선택하세요"
                                    size="large" class="w-full">
                                    <el-option v-for="path in depositPathOptions" :key="path.id" :label="path.name"
                                        :value="path.value" />
                                </el-select>
                            </el-form-item>
                        </div>

                        <!-- 공통 필드 -->
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <!-- 날짜 -->
                            <el-form-item label="날짜" prop="date">
                                <el-date-picker v-model="localEventDetail.date" type="date" placeholder="날짜를 선택하세요"
                                    size="large" format="YYYY-MM-DD" value-format="YYYY-MM-DD" />
                            </el-form-item>

                            <!-- 메모 -->
                            <el-form-item label="메모">
                                <el-input v-model="localEventDetail.memo" type="textarea" :rows="3"
                                    placeholder="메모를 입력하세요" size="large" />
                            </el-form-item>
                        </div>

                        <!-- 편집 모드 버튼 -->
                        <div class="flex justify-end space-x-3 pt-6 border-t border-gray-200">
                            <el-button @click="cancelEdit" size="large">취소</el-button>
                            <el-button @click="updateAccount" type="primary" size="large" :loading="updating"
                                :disabled="!isFormValid">
                                <Save class="w-4 h-4 mr-2" />
                                저장
                            </el-button>
                        </div>
                    </el-form>
                </div>

                <!-- 상세 보기 모드 -->
                <div v-else>
                    <!-- 기본 정보 -->
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
                        <div class="space-y-4">
                            <!-- 금액 -->
                            <div class="flex items-center p-4 bg-white border border-gray-200 rounded-lg">
                                <div class="w-10 h-10 rounded-lg flex items-center justify-center mr-4" :class="localEventDetail.type === 'out'
                                    ? 'bg-red-100 text-red-600'
                                    : 'bg-green-100 text-green-600'">
                                    <TrendingDown v-if="localEventDetail.type === 'out'" class="w-5 h-5" />
                                    <TrendingUp v-else class="w-5 h-5" />
                                </div>
                                <div>
                                    <p class="text-sm text-gray-500">금액</p>
                                    <p class="text-2xl font-bold"
                                        :class="localEventDetail.type === 'out' ? 'text-red-600' : 'text-green-600'">
                                        {{ formatMoney(localEventDetail.money) }}원
                                    </p>
                                </div>
                            </div>

                            <!-- 사용자 -->
                            <div class="flex items-center p-4 bg-white border border-gray-200 rounded-lg">
                                <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center mr-4">
                                    <span class="text-blue-600 font-semibold">{{ localEventDetail.user ?
                                        localEventDetail.user[0] : 'U'
                                    }}</span>
                                </div>
                                <div>
                                    <p class="text-sm text-gray-500">사용자</p>
                                    <p class="font-semibold text-gray-900">{{ localEventDetail.user || '-' }}</p>
                                </div>
                            </div>

                            <!-- 날짜 -->
                            <div class="flex items-center p-4 bg-white border border-gray-200 rounded-lg">
                                <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center mr-4">
                                    <Calendar class="w-5 h-5 text-purple-600" />
                                </div>
                                <div>
                                    <p class="text-sm text-gray-500">날짜</p>
                                    <p class="font-semibold text-gray-900">{{ formatDate(localEventDetail.date) }}</p>
                                </div>
                            </div>
                        </div>

                        <!-- 세부 정보 -->
                        <div class="space-y-4">
                            <!-- 카테고리 -->
                            <div class="flex items-center p-4 bg-white border border-gray-200 rounded-lg">
                                <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center mr-4">
                                    <Folder class="w-5 h-5 text-blue-600" />
                                </div>
                                <div>
                                    <p class="text-sm text-gray-500">카테고리</p>
                                    <p class="font-semibold text-gray-900">{{
                                        getCategoryName(localEventDetail.category_id) ||
                                        localEventDetail.category || '-' }}</p>
                                </div>
                            </div>

                            <!-- 키워드 (지출만) -->
                            <div v-if="(localEventDetail.keyword_name || localEventDetail.keyword) && localEventDetail.type === 'out'"
                                class="flex items-center p-4 bg-white border border-gray-200 rounded-lg">
                                <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center mr-4">
                                    <Tag class="w-5 h-5 text-purple-600" />
                                </div>
                                <div>
                                    <p class="text-sm text-gray-500">키워드</p>
                                    <p class="font-semibold text-gray-900">{{ localEventDetail.keyword_name ||
                                        localEventDetail.keyword
                                        ||
                                        '-' }}
                                    </p>
                                </div>
                            </div>

                            <!-- 결제 방법 (지출만) -->
                            <div v-if="localEventDetail.payment_method_id && localEventDetail.type === 'out'"
                                class="flex items-center p-4 bg-white border border-gray-200 rounded-lg">
                                <div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center mr-4">
                                    <CreditCard class="w-5 h-5 text-green-600" />
                                </div>
                                <div>
                                    <p class="text-sm text-gray-500">결제 방법</p>
                                    <p class="font-semibold text-gray-900">{{
                                        getPaymentMethodName(localEventDetail.payment_method_id)
                                        ||
                                        '-' }}</p>
                                </div>
                            </div>

                            <!-- 입금 경로 (수입만) -->
                            <div v-if="localEventDetail.deposit_path && localEventDetail.type === 'in'"
                                class="flex items-center p-4 bg-white border border-gray-200 rounded-lg">
                                <div class="w-10 h-10 bg-yellow-100 rounded-lg flex items-center justify-center mr-4">
                                    <CreditCard class="w-5 h-5 text-yellow-600" />
                                </div>
                                <div>
                                    <p class="text-sm text-gray-500">입금 경로</p>
                                    <p class="font-semibold text-gray-900">{{ localEventDetail.deposit_path || '-' }}
                                    </p>
                                </div>
                            </div>

                            <!-- 메모 -->
                            <div v-if="localEventDetail.memo"
                                class="flex items-start p-4 bg-white border border-gray-200 rounded-lg">
                                <div
                                    class="w-10 h-10 bg-gray-100 rounded-lg flex items-center justify-center mr-4 mt-0">
                                    <FileText class="w-5 h-5 text-gray-600" />
                                </div>
                                <div class="flex-1">
                                    <p class="text-sm text-gray-500">메모</p>
                                    <p class="font-semibold text-gray-900 whitespace-pre-wrap">{{ localEventDetail.memo
                                    }}</p>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- 상세 보기 모드 버튼 -->
                    <div class="flex justify-between pt-6 border-t border-gray-200">
                        <el-button @click="confirmDelete" type="danger" size="large">삭제</el-button>
                        <el-button @click="openEditMode" type="primary" size="large">편집</el-button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import {
    Tag,
    CreditCard,
    Save,
    TrendingUp,
    TrendingDown,
    Calendar,
    Folder,
    FileText
} from 'lucide-vue-next';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useAccountStore } from '../stores/accountStore';
import { useCategoryStore } from '../stores/categoryStore';
import { usePaymentMethodStore } from '../stores/paymentMethodStore';
import { useDepositPathStore } from '../stores/depositPathStore';

export default {
    name: 'DetailPage',
    components: {
        Tag,
        CreditCard,
        Save,
        TrendingUp,
        TrendingDown,
        Calendar,
        Folder,
        FileText
    },
    setup() {
        const router = useRouter();
        const route = useRoute();
        const formRef = ref(null);
        const updating = ref(false);
        const isEditMode = ref(false);

        // Stores
        const accountStore = useAccountStore();
        const categoryStore = useCategoryStore();
        const paymentMethodStore = usePaymentMethodStore();
        const depositPathStore = useDepositPathStore();

        // 라우트에서 데이터 받아오기
        const eventDetailFromRoute = route.query.data ? JSON.parse(route.query.data) : {};
        const localEventDetail = ref({ ...eventDetailFromRoute });

        // 폼 검증 규칙
        const rules = {
            user: [
                { required: true, message: '사용자명을 입력해주세요', trigger: 'blur' }
            ],
            money: [
                { required: true, message: '금액을 입력해주세요', trigger: 'blur' },
                { type: 'number', min: 1, message: '금액은 1원 이상이어야 합니다', trigger: 'blur' }
            ],
            category_name: [
                { required: true, message: '카테고리를 선택해주세요', trigger: 'change' }
            ],
            date: [
                { required: true, message: '날짜를 선택해주세요', trigger: 'change' }
            ]
        };

        // 사용 가능한 카테고리 목록 (백엔드에서 가져오기)
        const availableCategories = computed(() => {
            if (!categoryStore.categories || categoryStore.categories.length === 0) {
                return [];
            }

            const filtered = categoryStore.categories
                .filter(category => category.type === localEventDetail.value.type)
                .map(category => category.name);

            return filtered;
        });

        // 입금경로 옵션 목록
        const depositPathOptions = computed(() => {
            return depositPathStore.activeDepositPaths.map(path => ({
                id: path.id,
                name: path.name,
                value: path.name
            }));
        });

        // 폼 유효성 검사
        const isFormValid = computed(() => {
            return localEventDetail.value.user &&
                localEventDetail.value.money > 0 &&
                (localEventDetail.value.category_name || localEventDetail.value.category) &&
                localEventDetail.value.date;
        });

        // 뒤로 가기
        const goBack = () => {
            router.push('/');
        };

        // 계정 업데이트
        const updateAccount = async () => {
            if (!formRef.value) return;

            try {
                // 폼 검증
                await formRef.value.validate();

                updating.value = true;

                // 데이터 정리
                const accountData = { ...localEventDetail.value };

                // category_name이 있으면 category 필드에도 설정
                if (accountData.category_name) {
                    accountData.category = accountData.category_name;
                }

                console.log('DetailPage에서 업데이트 이벤트 발생:', accountData);

                // 실제 업데이트 API 호출
                await accountStore.updateAccount(accountData);

                // 컴포넌트가 여전히 마운트된 상태인지 확인
                if (isEditMode.value !== undefined) {
                    isEditMode.value = false;
                    ElMessage.success('데이터가 성공적으로 업데이트되었습니다.');

                    // 업데이트된 데이터로 로컬 상태 갱신
                    localEventDetail.value = { ...accountData };
                }

            } catch (error) {
                console.error('업데이트 실패:', error);
                if (isEditMode.value !== undefined) {
                    ElMessage.error('데이터 업데이트 중 오류가 발생했습니다.');
                }
            } finally {
                if (updating.value !== undefined) {
                    updating.value = false;
                }
            }
        };

        // 삭제 확인
        const confirmDelete = async () => {
            try {
                await ElMessageBox.confirm(
                    '정말로 이 항목을 삭제하시겠습니까? 삭제된 데이터는 복구할 수 없습니다.',
                    '삭제 확인',
                    {
                        confirmButtonText: '삭제',
                        cancelButtonText: '취소',
                        type: 'warning',
                        confirmButtonClass: 'el-button--danger'
                    }
                );

                // 실제 삭제 API 호출
                await accountStore.deleteAccount(localEventDetail.value);

                ElMessage.success('데이터가 삭제되었습니다.');

                // 삭제 후 메인 페이지로 이동
                goBack();

            } catch (error) {
                // 사용자가 취소했거나 오류 발생
                if (error !== 'cancel') {
                    console.error('삭제 실패:', error);
                    ElMessage.error('데이터 삭제 중 오류가 발생했습니다.');
                }
            }
        };

        // 편집 모드 열기
        const openEditMode = async () => {
            try {
                // 항상 최신 카테고리 데이터를 로드 (캐시 무시)
                await categoryStore.fetchCategories();

                isEditMode.value = true;
                // 편집 모드 시작 시 현재 데이터를 다시 로드하여 최신 상태 반영
                localEventDetail.value = { ...eventDetailFromRoute };

                // 카테고리 이름 설정 (ID가 있는 경우 이름으로 변환)
                if (localEventDetail.value.category_id && !localEventDetail.value.category_name) {
                    const categoryName = getCategoryName(localEventDetail.value.category_id);
                    if (categoryName) {
                        localEventDetail.value.category_name = categoryName;
                    }
                }
                // category 필드가 있고 category_name이 없는 경우
                else if (localEventDetail.value.category && !localEventDetail.value.category_name) {
                    localEventDetail.value.category_name = localEventDetail.value.category;
                }
            } catch (error) {
                console.error('편집 모드 시작 중 오류:', error);
                ElMessage.error('편집 모드를 시작하는데 실패했습니다.');
            }
        };

        // 편집 취소
        const cancelEdit = () => {
            isEditMode.value = false;
            // 원본 데이터로 되돌리기
            localEventDetail.value = { ...eventDetailFromRoute };
        };

        // 금액 포맷팅
        const formatMoney = (amount) => {
            if (!amount) return '0';
            return new Intl.NumberFormat('ko-KR').format(amount);
        };

        // 날짜 포맷팅 (시간 제거)
        const formatDate = (dateString) => {
            if (!dateString) return '';

            try {
                const date = new Date(dateString);
                const options = {
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                    timeZone: 'Asia/Seoul'
                };

                return date.toLocaleDateString('ko-KR', options);
            } catch {
                return dateString;
            }
        };

        // 카테고리 이름 가져오기
        const getCategoryName = (categoryId) => {
            if (!categoryId) return '';
            const category = categoryStore.categories.find(cat => cat.id === categoryId);
            return category ? category.name : '';
        };

        // 결제수단 이름 가져오기
        const getPaymentMethodName = (paymentMethodId) => {
            if (!paymentMethodId) return '';
            const paymentMethod = paymentMethodStore.flatPaymentMethods.find(pm => pm.id === paymentMethodId);
            return paymentMethod ? paymentMethod.name : '';
        };

        // 컴포넌트 마운트 시 데이터 로드
        onMounted(async () => {
            try {
                // 모든 필요한 데이터를 병렬로 로드
                await Promise.all([
                    categoryStore.fetchCategories(),
                    paymentMethodStore.fetchPaymentMethods(),
                    depositPathStore.fetchDepositPaths()
                ]);

                // 초기 카테고리 이름 설정
                if (!localEventDetail.value.category_name) {
                    if (localEventDetail.value.category_id) {
                        const categoryName = getCategoryName(localEventDetail.value.category_id);
                        if (categoryName) {
                            localEventDetail.value.category_name = categoryName;
                        }
                    } else if (localEventDetail.value.category) {
                        localEventDetail.value.category_name = localEventDetail.value.category;
                    }
                }
            } catch (error) {
                console.error('데이터 로드 오류:', error);
                ElMessage.error('데이터를 불러오는데 실패했습니다.');
            }
        });

        return {
            formRef,
            updating,
            isEditMode,
            localEventDetail,
            rules,
            availableCategories,
            depositPathOptions,
            isFormValid,
            updateAccount,
            confirmDelete,
            goBack,
            openEditMode,
            cancelEdit,
            formatMoney,
            formatDate,
            getCategoryName,
            getPaymentMethodName,

            // Stores
            paymentMethodStore,
            categoryStore,

            // 아이콘들
            Tag,
            CreditCard,
            Save,
            TrendingUp,
            TrendingDown,
            Calendar,
            Folder,
            FileText
        };
    }
};
</script>

<style scoped>
/* 페이지 전체 스타일 */
.detail-page {
    min-height: 100vh;
}

/* 커스텀 입력 필드 스타일 */
:deep(.el-form-item__label) {
    @apply font-semibold text-gray-700 text-sm;
}

:deep(.el-input__wrapper) {
    @apply border-gray-200 rounded-lg;
    transition: all 0.2s ease;
}

:deep(.el-input__wrapper:hover) {
    @apply border-blue-300;
}

:deep(.el-input__wrapper.is-focus) {
    @apply border-blue-500 shadow-sm;
}

:deep(.el-select) {
    width: 100%;
}

:deep(.el-select .el-input__wrapper) {
    @apply border-gray-200 rounded-lg;
}

:deep(.el-select .el-input__wrapper:hover) {
    @apply border-blue-300;
}

:deep(.el-select .el-input__wrapper.is-focus) {
    @apply border-blue-500;
}

:deep(.el-textarea__inner) {
    @apply border-gray-200 rounded-lg;
    transition: all 0.2s ease;
}

:deep(.el-textarea__inner:hover) {
    @apply border-blue-300;
}

:deep(.el-textarea__inner:focus) {
    @apply border-blue-500 shadow-sm;
}

:deep(.el-date-editor) {
    width: 100%;
}

:deep(.el-date-editor .el-input__wrapper) {
    @apply border-gray-200 rounded-lg;
}

/* 반응형 스타일 */
@media (max-width: 768px) {
    .detail-page {
        padding: 0;
    }

    .max-w-4xl {
        margin: 0;
        padding: 1rem;
    }

    .bg-white.rounded-lg {
        border-radius: 0;
        margin: 0;
    }
}
</style>
