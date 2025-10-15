<template>
    <div class="management-page min-h-screen bg-gray-50 py-6 px-4">
        <div class="max-w-7xl mx-auto">
            <!-- 헤더 -->
            <div class="bg-white rounded-xl shadow-md p-6 mb-6">
                <div class="flex items-center justify-between">
                    <div>
                        <h1 class="text-3xl font-bold text-gray-800 mb-2">⚙️ 관리</h1>
                        <p class="text-gray-600">가계부 설정 및 데이터를 관리합니다</p>
                    </div>
                    <el-button @click="goBack" type="primary" size="large">
                        <span class="text-lg">← 달력으로 돌아가기</span>
                    </el-button>
                </div>
            </div>

            <!-- 관리 메뉴 그리드 -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">

                <!-- 사용자 관리 -->
                <div @click="openUserManager"
                    class="bg-white rounded-xl shadow-md p-8 cursor-pointer hover:shadow-xl transition-all duration-200 hover:scale-105">
                    <div class="flex flex-col items-center text-center">
                        <div class="w-20 h-20 bg-blue-100 rounded-full flex items-center justify-center mb-4">
                            <span class="text-4xl">👥</span>
                        </div>
                        <h2 class="text-2xl font-bold text-gray-800 mb-2">사용자 관리</h2>
                        <p class="text-gray-600 mb-4">가계부 사용자를 추가하고 관리합니다</p>
                        <span class="text-sm text-blue-600 font-semibold">관리하기 →</span>
                    </div>
                </div>

                <!-- 카테고리 관리 -->
                <div @click="openCategoryManager"
                    class="bg-white rounded-xl shadow-md p-8 cursor-pointer hover:shadow-xl transition-all duration-200 hover:scale-105">
                    <div class="flex flex-col items-center text-center">
                        <div class="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mb-4">
                            <span class="text-4xl">🏷️</span>
                        </div>
                        <h2 class="text-2xl font-bold text-gray-800 mb-2">카테고리 관리</h2>
                        <p class="text-gray-600 mb-4">수입/지출 카테고리를 관리합니다</p>
                        <span class="text-sm text-green-600 font-semibold">관리하기 →</span>
                    </div>
                </div>

                <!-- 키워드 관리 -->
                <div @click="openKeywordManager"
                    class="bg-white rounded-xl shadow-md p-8 cursor-pointer hover:shadow-xl transition-all duration-200 hover:scale-105">
                    <div class="flex flex-col items-center text-center">
                        <div class="w-20 h-20 bg-purple-100 rounded-full flex items-center justify-center mb-4">
                            <span class="text-4xl">🔖</span>
                        </div>
                        <h2 class="text-2xl font-bold text-gray-800 mb-2">키워드 관리</h2>
                        <p class="text-gray-600 mb-4">카테고리별 키워드를 관리합니다</p>
                        <span class="text-sm text-purple-600 font-semibold">관리하기 →</span>
                    </div>
                </div>

                <!-- 결제수단 관리 -->
                <div @click="openPaymentMethodManager"
                    class="bg-white rounded-xl shadow-md p-8 cursor-pointer hover:shadow-xl transition-all duration-200 hover:scale-105">
                    <div class="flex flex-col items-center text-center">
                        <div class="w-20 h-20 bg-orange-100 rounded-full flex items-center justify-center mb-4">
                            <span class="text-4xl">💳</span>
                        </div>
                        <h2 class="text-2xl font-bold text-gray-800 mb-2">결제수단 관리</h2>
                        <p class="text-gray-600 mb-4">카드, 현금 등 결제수단을 관리합니다</p>
                        <span class="text-sm text-orange-600 font-semibold">관리하기 →</span>
                    </div>
                </div>

                <!-- 입금경로 관리 -->
                <div @click="openDepositPathManager"
                    class="bg-white rounded-xl shadow-md p-8 cursor-pointer hover:shadow-xl transition-all duration-200 hover:scale-105">
                    <div class="flex flex-col items-center text-center">
                        <div class="w-20 h-20 bg-teal-100 rounded-full flex items-center justify-center mb-4">
                            <span class="text-4xl">🏦</span>
                        </div>
                        <h2 class="text-2xl font-bold text-gray-800 mb-2">입금경로 관리</h2>
                        <p class="text-gray-600 mb-4">수입 입금경로를 관리합니다</p>
                        <span class="text-sm text-teal-600 font-semibold">관리하기 →</span>
                    </div>
                </div>

                <!-- 예산(기준치) 관리 -->
                <div @click="openBudgetManager"
                    class="bg-white rounded-xl shadow-md p-8 cursor-pointer hover:shadow-xl transition-all duration-200 hover:scale-105">
                    <div class="flex flex-col items-center text-center">
                        <div class="w-20 h-20 bg-yellow-100 rounded-full flex items-center justify-center mb-4">
                            <span class="text-4xl">💰</span>
                        </div>
                        <h2 class="text-2xl font-bold text-gray-800 mb-2">기준치 관리</h2>
                        <p class="text-gray-600 mb-4">카테고리별 예산 기준치를 설정합니다</p>
                        <span class="text-sm text-yellow-600 font-semibold">관리하기 →</span>
                    </div>
                </div>

            </div>

            <!-- 관리 팁 -->
            <div class="bg-blue-50 border border-blue-200 rounded-xl p-6 mt-6">
                <h3 class="text-lg font-semibold text-blue-800 mb-3">💡 관리 팁</h3>
                <ul class="space-y-2 text-blue-700">
                    <li>• <strong>사용자 관리</strong>: 가족 구성원별로 가계부를 분리 관리할 수 있습니다</li>
                    <li>• <strong>카테고리 관리</strong>: 지출 카테고리에서 고정/변동 유형을 설정하세요</li>
                    <li>• <strong>키워드 관리</strong>: 자주 사용하는 항목을 키워드로 등록하면 빠른 입력이 가능합니다</li>
                    <li>• <strong>기준치 관리</strong>: 카테고리별 예산을 설정하면 초과 시 알림을 받을 수 있습니다</li>
                </ul>
            </div>

        </div>

        <!-- 관리 모달들 -->
        <UserManager v-if="showUserManager" @close="closeUserManager" />
        <CategoryManager v-if="showCategoryManager" @close="closeCategoryManager" />
        <KeywordManager v-if="showKeywordManager" :category-id="keywordManagerCategoryId"
            @close="closeKeywordManager" />
        <PaymentMethodManager v-if="showPaymentMethodManager" @close="closePaymentMethodManager" />
        <DepositPathManager v-if="showDepositPathManager" @close="closeDepositPathManager" />
        <BudgetManager v-if="showBudgetManager" @close="closeBudgetManager" />
    </div>
</template>

<script>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import UserManager from './UserManager.vue';
import CategoryManager from './CategoryManager.vue';
import KeywordManager from './KeywordManager.vue';
import PaymentMethodManager from './PaymentMethodManager.vue';
import DepositPathManager from './DepositPathManager.vue';
import BudgetManager from './BudgetManager.vue';

export default {
    name: 'ManagementPage',
    components: {
        UserManager,
        CategoryManager,
        KeywordManager,
        PaymentMethodManager,
        DepositPathManager,
        BudgetManager
    },
    setup() {
        const router = useRouter();

        // 관리 모달 상태
        const showUserManager = ref(false);
        const showCategoryManager = ref(false);
        const showKeywordManager = ref(false);
        const showPaymentMethodManager = ref(false);
        const showDepositPathManager = ref(false);
        const showBudgetManager = ref(false);
        const keywordManagerCategoryId = ref(null);

        // 뒤로가기
        const goBack = () => {
            router.push('/');
        };

        // 관리 모달 열기/닫기 메서드
        const openUserManager = () => {
            showUserManager.value = true;
        };

        const closeUserManager = () => {
            showUserManager.value = false;
        };

        const openCategoryManager = () => {
            showCategoryManager.value = true;
        };

        const closeCategoryManager = () => {
            showCategoryManager.value = false;
        };

        const openKeywordManager = () => {
            keywordManagerCategoryId.value = null;
            showKeywordManager.value = true;
        };

        const closeKeywordManager = () => {
            showKeywordManager.value = false;
            keywordManagerCategoryId.value = null;
        };

        const openPaymentMethodManager = () => {
            showPaymentMethodManager.value = true;
        };

        const closePaymentMethodManager = () => {
            showPaymentMethodManager.value = false;
        };

        const openDepositPathManager = () => {
            showDepositPathManager.value = true;
        };

        const closeDepositPathManager = () => {
            showDepositPathManager.value = false;
        };

        const openBudgetManager = () => {
            showBudgetManager.value = true;
        };

        const closeBudgetManager = () => {
            showBudgetManager.value = false;
        };

        return {
            goBack,
            showUserManager,
            showCategoryManager,
            showKeywordManager,
            showPaymentMethodManager,
            showDepositPathManager,
            showBudgetManager,
            keywordManagerCategoryId,
            openUserManager,
            closeUserManager,
            openCategoryManager,
            closeCategoryManager,
            openKeywordManager,
            closeKeywordManager,
            openPaymentMethodManager,
            closePaymentMethodManager,
            openDepositPathManager,
            closeDepositPathManager,
            openBudgetManager,
            closeBudgetManager
        };
    }
};
</script>

<style scoped>
.management-page {
    min-height: 100vh;
}
</style>
