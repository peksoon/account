<template>
    <div class="statistics-dashboard">
        <!-- 헤더 -->
        <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between mb-8">
            <div>
                <h1 class="text-3xl font-bold text-gradient mb-2">📊 통계 대시보드</h1>
                <p class="text-gray-600">가계부 통계를 한눈에 확인하세요</p>
            </div>

            <!-- 필터 컨트롤 -->
            <div class="flex flex-col sm:flex-row gap-4 mt-4 lg:mt-0">
                <el-select v-model="selectedUser" @change="handleUserChange" class="w-full sm:w-40"
                    placeholder="사용자 선택">
                    <el-option label="전체" value="" />
                    <el-option v-for="user in users" :key="user.id" :label="user.name" :value="user.name" />
                </el-select>

                <el-select v-model="selectedType" @change="handleTypeChange" class="w-full sm:w-40">
                    <el-option label="지출" value="out" />
                    <el-option label="수입" value="in" />
                </el-select>

                <el-select v-model="selectedPeriod" @change="handlePeriodChange" class="w-full sm:w-40">
                    <el-option label="이번 주" value="week" />
                    <el-option label="이번 달" value="month" />
                    <el-option label="올해" value="year" />
                    <el-option label="전체" value="all" />
                    <el-option label="기간 설정" value="custom" />
                </el-select>

                <div v-if="selectedPeriod === 'custom'" class="flex gap-2">
                    <el-date-picker v-model="customStartDate" type="date" placeholder="시작일" format="YYYY-MM-DD"
                        @change="handleCustomDateChange" />
                    <el-date-picker v-model="customEndDate" type="date" placeholder="종료일" format="YYYY-MM-DD"
                        @change="handleCustomDateChange" />
                </div>
            </div>
        </div>

        <!-- 통계 카드 -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div class="stat-card" :class="selectedType === 'out' ? 'stat-card-red' : 'stat-card-green'">
                <div class="flex items-center">
                    <div class="stat-icon">
                        <TrendingDown v-if="selectedType === 'out'" class="w-8 h-8" />
                        <TrendingUp v-else class="w-8 h-8" />
                    </div>
                    <div class="ml-4">
                        <p class="stat-label">총 {{ selectedType === 'out' ? '지출' : '수입' }}</p>
                        <p class="stat-value">{{ formatMoney(statistics?.total_amount || 0) }}원</p>
                    </div>
                </div>
            </div>

            <div class="stat-card stat-card-blue">
                <div class="flex items-center">
                    <div class="stat-icon bg-blue-100 text-blue-600">
                        <BarChart class="w-8 h-8" />
                    </div>
                    <div class="ml-4">
                        <p class="stat-label">총 건수</p>
                        <p class="stat-value">{{ statistics?.total_count || 0 }}건</p>
                    </div>
                </div>
            </div>

            <div class="stat-card stat-card-purple">
                <div class="flex items-center">
                    <div class="stat-icon bg-purple-100 text-purple-600">
                        <Calculator class="w-8 h-8" />
                    </div>
                    <div class="ml-4">
                        <p class="stat-label">평균 금액</p>
                        <p class="stat-value">
                            {{ formatMoney(statistics?.total_count ? Math.round(statistics.total_amount /
                                statistics.total_count) : 0) }}원
                        </p>
                    </div>
                </div>
            </div>
        </div>

        <!-- 차트 영역 -->
        <div class="grid grid-cols-1 xl:grid-cols-2 gap-8 mb-8">
            <!-- 도넛 차트 -->
            <div class="card">
                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-xl font-bold text-gray-800">카테고리별 {{ selectedType === 'out' ? '지출' : '수입' }}</h2>
                    <span class="text-sm text-gray-500">{{ statistics?.period }}</span>
                </div>

                <div class="chart-container">
                    <Doughnut v-if="chartData.datasets[0].data.length > 0" :data="chartData" :options="chartOptions"
                        @click="handleChartClick" />
                    <div v-else class="empty-chart">
                        <PieChart class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                        <p class="text-gray-500">표시할 데이터가 없습니다</p>
                    </div>
                </div>
            </div>

            <!-- 상위 카테고리 리스트 -->
            <div class="card">
                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-xl font-bold text-gray-800">상위 카테고리</h2>
                    <el-button size="small" @click="toggleSortOrder">
                        {{ sortOrder === 'desc' ? '↓ 높은순' : '↑ 낮은순' }}
                    </el-button>
                </div>

                <div class="space-y-3">
                    <div v-for="(category, index) in sortedCategories" :key="category.category_id" class="category-item"
                        @click="showCategoryDetail(category)">
                        <div class="flex items-center">
                            <div class="rank-badge" :class="`rank-${index + 1}`">
                                {{ index + 1 }}
                            </div>
                            <div class="ml-3 flex-1">
                                <div class="flex items-center justify-between">
                                    <span class="font-medium text-gray-900">{{ category.category_name }}</span>
                                    <span class="font-bold text-gray-800">{{ formatMoney(category.total_amount)
                                        }}원</span>
                                </div>
                                <div class="flex items-center justify-between mt-1">
                                    <span class="text-sm text-gray-500">{{ category.count }}건</span>
                                    <span class="text-sm font-medium"
                                        :class="selectedType === 'out' ? 'text-red-600' : 'text-green-600'">
                                        {{ category.percentage.toFixed(1) }}%
                                    </span>
                                </div>
                            </div>
                        </div>
                        <div class="progress-bar mt-2">
                            <div class="progress-fill" :class="selectedType === 'out' ? 'bg-red-500' : 'bg-green-500'"
                                :style="`width: ${category.percentage}%`"></div>
                        </div>
                    </div>
                </div>

                <div v-if="sortedCategories.length === 0" class="empty-state">
                    <Folder class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                    <p class="text-gray-500">해당 기간에 데이터가 없습니다</p>
                </div>
            </div>
        </div>

        <!-- 기준치 정보 섹션 (지출일 때만 표시) -->
        <div v-if="selectedType === 'out' && selectedUser && budgetUsages && budgetUsages.length > 0" class="mb-8">
            <div class="card">
                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-xl font-bold text-gray-800">📊 기준치 사용량</h2>
                    <span class="text-sm text-gray-500">{{ selectedUser }}님의 기준치 현황</span>
                </div>

                <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    <div v-for="budget in budgetUsages" :key="budget.category_id"
                        class="budget-card p-4 border border-gray-200 rounded-lg" :class="{
                            'border-red-300 bg-red-50': budget.is_monthly_over || budget.is_yearly_over,
                            'border-yellow-300 bg-yellow-50': !budget.is_monthly_over && !budget.is_yearly_over && isNearLimit(budget),
                            'border-green-300 bg-green-50': !budget.is_monthly_over && !budget.is_yearly_over && !isNearLimit(budget)
                        }">
                        <h3 class="font-semibold text-gray-800 mb-3">{{ budget.category_name }}</h3>
                        <BudgetUsageDisplay :usage="budget" />
                    </div>
                </div>

                <!-- 전체 기준치 요약 -->
                <div class="mt-6 pt-6 border-t border-gray-200">
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div class="text-center">
                            <p class="text-sm text-gray-600">전체 월 기준치</p>
                            <p class="text-lg font-bold text-blue-600">{{ formatMoney(totalMonthlyBudget) }}원</p>
                        </div>
                        <div class="text-center">
                            <p class="text-sm text-gray-600">전체 월 사용량</p>
                            <p class="text-lg font-bold"
                                :class="totalMonthlyUsed > totalMonthlyBudget ? 'text-red-600' : 'text-green-600'">
                                {{ formatMoney(totalMonthlyUsed) }}원
                            </p>
                        </div>
                        <div class="text-center">
                            <p class="text-sm text-gray-600">초과 카테고리</p>
                            <p class="text-lg font-bold text-red-600">{{ overBudgetCount }}개</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 기준치 설정 안내 (지출이고 사용자가 선택됐지만 기준치가 없는 경우) -->
        <div v-if="selectedType === 'out' && selectedUser && (!budgetUsages || budgetUsages.length === 0)" class="mb-8">
            <div class="card">
                <div class="text-center py-8">
                    <Calculator class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                    <h3 class="text-lg font-semibold text-gray-700 mb-2">기준치가 설정되지 않았습니다</h3>
                    <p class="text-gray-500 mb-4">카테고리별 기준치를 설정하여 지출을 관리해보세요.</p>
                    <el-button type="primary" @click="$emit('open-budget-manager')">
                        기준치 설정하기
                    </el-button>
                </div>
            </div>
        </div>

        <!-- 키워드 상세 다이얼로그 -->
        <el-dialog v-model="keywordDialogVisible" :title="`${selectedCategory?.category_name} 키워드 상세`" width="600px"
            destroy-on-close>
            <div v-if="keywordStatistics">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                    <div class="keyword-chart">
                        <Doughnut v-if="keywordChartData.datasets[0].data.length > 0" :data="keywordChartData"
                            :options="keywordChartOptions" />
                    </div>
                    <div class="keyword-list">
                        <div v-for="keyword in keywordStatistics.keywords" :key="keyword.keyword_id"
                            class="keyword-item">
                            <div class="flex items-center justify-between">
                                <span class="font-medium">{{ keyword.keyword_name }}</span>
                                <span class="font-bold">{{ formatMoney(keyword.total_amount) }}원</span>
                            </div>
                            <div class="flex items-center justify-between text-sm text-gray-500">
                                <span>{{ keyword.count }}건</span>
                                <span>{{ keyword.percentage.toFixed(1) }}%</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="loadingKeywords" class="text-center py-8">
                <div class="spinner mx-auto"></div>
                <p class="text-gray-600 mt-2">키워드 데이터를 불러오는 중...</p>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import {
    TrendingUp,
    TrendingDown,
    BarChart,
    Calculator,
    PieChart,
    Folder
} from 'lucide-vue-next';
import { ElMessage } from 'element-plus';
import BudgetUsageDisplay from './BudgetUsageDisplay.vue';
import { Doughnut } from 'vue-chartjs';
import {
    Chart as ChartJS,
    ArcElement,
    Tooltip,
    Legend
} from 'chart.js';
import { useStatisticsStore } from '../stores/statisticsStore';
import { useUserStore } from '../stores/userStore';

// Chart.js 등록
ChartJS.register(ArcElement, Tooltip, Legend);

export default {
    name: 'StatisticsDashboard',
    components: {
        TrendingUp,
        TrendingDown,
        BarChart,
        Calculator,
        PieChart,
        Folder,
        Doughnut,
        BudgetUsageDisplay
    },
    emits: ['close', 'open-budget-manager'],
    setup() {
        const statisticsStore = useStatisticsStore();
        const userStore = useUserStore();

        const selectedType = ref('out');
        const selectedPeriod = ref('month');
        const selectedUser = ref('');
        const customStartDate = ref(null);
        const customEndDate = ref(null);
        const sortOrder = ref('desc');
        const keywordDialogVisible = ref(false);
        const selectedCategory = ref(null);
        const loadingKeywords = ref(false);

        const statistics = computed(() => statisticsStore.statistics);
        const keywordStatistics = computed(() => statisticsStore.keywordStatistics);
        const users = computed(() => userStore.users || []);
        const budgetUsages = computed(() => statistics.value?.budget_usages || []);

        // 정렬된 카테고리
        const sortedCategories = computed(() => {
            if (!statistics.value?.categories) return [];

            const categories = [...statistics.value.categories];
            return categories.sort((a, b) => {
                return sortOrder.value === 'desc'
                    ? b.total_amount - a.total_amount
                    : a.total_amount - b.total_amount;
            });
        });

        // 차트 데이터
        const chartData = computed(() => {
            if (!statistics.value?.chart_data) {
                return {
                    labels: [],
                    datasets: [{
                        data: [],
                        backgroundColor: [],
                        borderWidth: 0
                    }]
                };
            }

            return {
                labels: statistics.value.chart_data.map(item => item.label),
                datasets: [{
                    data: statistics.value.chart_data.map(item => item.value),
                    backgroundColor: statistics.value.chart_data.map(item => item.color),
                    borderWidth: 0,
                    hoverOffset: 10
                }]
            };
        });

        // 키워드 차트 데이터
        const keywordChartData = computed(() => {
            if (!keywordStatistics.value?.chart_data) {
                return {
                    labels: [],
                    datasets: [{
                        data: [],
                        backgroundColor: [],
                        borderWidth: 0
                    }]
                };
            }

            return {
                labels: keywordStatistics.value.chart_data.map(item => item.label),
                datasets: [{
                    data: keywordStatistics.value.chart_data.map(item => item.value),
                    backgroundColor: keywordStatistics.value.chart_data.map(item => item.color),
                    borderWidth: 0,
                    hoverOffset: 10
                }]
            };
        });

        // 차트 옵션
        const chartOptions = {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'bottom',
                    labels: {
                        padding: 20,
                        usePointStyle: true
                    }
                },
                tooltip: {
                    callbacks: {
                        label: function (context) {
                            const value = context.parsed;
                            const total = context.dataset.data.reduce((a, b) => a + b, 0);
                            const percentage = ((value / total) * 100).toFixed(1);
                            return `${context.label}: ${formatMoney(value)}원 (${percentage}%)`;
                        }
                    }
                }
            },
            onClick: (event, elements) => {
                if (elements.length > 0) {
                    const index = elements[0].index;
                    const category = statistics.value.categories[index];
                    showCategoryDetail(category);
                }
            }
        };

        const keywordChartOptions = {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'bottom',
                    labels: {
                        padding: 15,
                        usePointStyle: true,
                        font: {
                            size: 11
                        }
                    }
                },
                tooltip: {
                    callbacks: {
                        label: function (context) {
                            const value = context.parsed;
                            const total = context.dataset.data.reduce((a, b) => a + b, 0);
                            const percentage = ((value / total) * 100).toFixed(1);
                            return `${context.label}: ${formatMoney(value)}원 (${percentage}%)`;
                        }
                    }
                }
            }
        };

        // 금액 포맷팅
        const formatMoney = (amount) => {
            if (!amount) return '0';
            return new Intl.NumberFormat('ko-KR').format(amount);
        };

        // 기준치 관련 computed
        const totalMonthlyBudget = computed(() => {
            return budgetUsages.value.reduce((total, usage) => total + (usage.monthly_budget || 0), 0);
        });

        const totalMonthlyUsed = computed(() => {
            return budgetUsages.value.reduce((total, usage) => total + (usage.monthly_used || 0), 0);
        });

        const overBudgetCount = computed(() => {
            return budgetUsages.value.filter(usage => usage.is_monthly_over || usage.is_yearly_over).length;
        });

        // 기준치 근접 여부 확인
        const isNearLimit = (budget) => {
            return (budget.monthly_budget > 0 && budget.monthly_percent >= 80) ||
                (budget.yearly_budget > 0 && budget.yearly_percent >= 80);
        };

        // 통계 데이터 로드
        const loadStatistics = async () => {
            try {
                const params = {
                    type: selectedPeriod.value,
                    category: selectedType.value
                };

                // 사용자가 선택된 경우 파라미터에 추가
                if (selectedUser.value) {
                    params.user = selectedUser.value;
                }

                if (selectedPeriod.value === 'custom' && customStartDate.value && customEndDate.value) {
                    params.start_date = customStartDate.value;
                    params.end_date = customEndDate.value;
                }

                await statisticsStore.fetchStatistics(params);
            } catch (error) {
                ElMessage.error('통계 데이터를 불러오는데 실패했습니다.');
            }
        };

        // 사용자 변경 핸들러
        const handleUserChange = () => {
            loadStatistics();
        };

        // 타입 변경 핸들러
        const handleTypeChange = () => {
            loadStatistics();
        };

        // 기간 변경 핸들러
        const handlePeriodChange = () => {
            if (selectedPeriod.value !== 'custom') {
                customStartDate.value = null;
                customEndDate.value = null;
            }
            loadStatistics();
        };

        // 커스텀 날짜 변경 핸들러
        const handleCustomDateChange = () => {
            if (customStartDate.value && customEndDate.value) {
                loadStatistics();
            }
        };

        // 정렬 순서 토글
        const toggleSortOrder = () => {
            sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc';
        };

        // 카테고리 상세 보기
        const showCategoryDetail = async (category) => {
            selectedCategory.value = category;
            keywordDialogVisible.value = true;
            loadingKeywords.value = true;

            try {
                const params = {
                    category_id: category.category_id,
                    type: selectedPeriod.value,
                    category: selectedType.value
                };

                if (selectedPeriod.value === 'custom' && customStartDate.value && customEndDate.value) {
                    params.start_date = customStartDate.value;
                    params.end_date = customEndDate.value;
                }

                await statisticsStore.fetchKeywordStatistics(params);
            } catch (error) {
                ElMessage.error('키워드 통계를 불러오는데 실패했습니다.');
            } finally {
                loadingKeywords.value = false;
            }
        };

        // 차트 클릭 핸들러
        const handleChartClick = (event, elements) => {
            if (elements && elements.length > 0) {
                const index = elements[0].index;
                if (statistics.value?.categories && statistics.value.categories[index]) {
                    const category = statistics.value.categories[index];
                    showCategoryDetail(category);
                }
            }
        };

        onMounted(async () => {
            // 사용자 목록 로드
            try {
                await userStore.fetchUsers();
            } catch (error) {
                console.error('사용자 목록 로드 오류:', error);
            }

            loadStatistics();
        });

        return {
            selectedType,
            selectedPeriod,
            selectedUser,
            customStartDate,
            customEndDate,
            sortOrder,
            keywordDialogVisible,
            selectedCategory,
            loadingKeywords,

            statistics,
            keywordStatistics,
            users,
            budgetUsages,
            sortedCategories,
            chartData,
            keywordChartData,
            chartOptions,
            keywordChartOptions,

            // 기준치 관련
            totalMonthlyBudget,
            totalMonthlyUsed,
            overBudgetCount,
            isNearLimit,

            formatMoney,
            handleUserChange,
            handleTypeChange,
            handlePeriodChange,
            handleCustomDateChange,
            toggleSortOrder,
            showCategoryDetail,
            handleChartClick,

            // 아이콘들
            TrendingUp,
            TrendingDown,
            BarChart,
            Calculator,
            PieChart,
            Folder
        };
    }
};
</script>

<style scoped>
.chart-container {
    height: 300px;
    position: relative;
}

.stat-card {
    background-color: white;
    padding: 1.5rem;
    border-radius: 0.75rem;
    border: 1px solid #e5e7eb;
    box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

.stat-card-red {
    border-color: #fecaca;
    background: linear-gradient(to right, #fef2f2, #fee2e2);
}

.stat-card-green {
    border-color: #bbf7d0;
    background: linear-gradient(to right, #f0fdf4, #dcfce7);
}

.stat-card-blue {
    border-color: #bfdbfe;
    background: linear-gradient(to right, #eff6ff, #dbeafe);
}

.stat-card-purple {
    border-color: #e9d5ff;
    background: linear-gradient(to right, #faf5ff, #f3e8ff);
}

.stat-icon {
    width: 3rem;
    height: 3rem;
    border-radius: 0.75rem;
    display: flex;
    align-items: center;
    justify-content: center;
}

.stat-card-red .stat-icon {
    background-color: #fee2e2;
    color: #dc2626;
}

.stat-card-green .stat-icon {
    background-color: #dcfce7;
    color: #16a34a;
}

.stat-label {
    font-size: 0.875rem;
    font-weight: 500;
    color: #6b7280;
}

.stat-value {
    font-size: 1.5rem;
    font-weight: 700;
    color: #111827;
}

.category-item {
    padding: 1rem;
    background-color: #f9fafb;
    border-radius: 0.5rem;
    cursor: pointer;
    transition: background-color 0.2s;
}

.category-item:hover {
    background-color: #f3f4f6;
}

.rank-badge {
    width: 2rem;
    height: 2rem;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.875rem;
    font-weight: 700;
    color: white;
}

.rank-1 {
    background-color: #eab308;
}

.rank-2 {
    background-color: #9ca3af;
}

.rank-3 {
    background-color: #f97316;
}

.rank-badge:not(.rank-1):not(.rank-2):not(.rank-3) {
    background-color: #3b82f6;
}

.progress-bar {
    width: 100%;
    background-color: #e5e7eb;
    border-radius: 9999px;
    height: 0.375rem;
}

.progress-fill {
    height: 0.375rem;
    border-radius: 9999px;
    transition: all 0.3s;
}

.empty-chart {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    text-align: center;
}

.empty-state {
    text-align: center;
    padding: 2rem 0;
}

.keyword-chart {
    height: 250px;
}

.keyword-item {
    padding: 0.75rem;
    background-color: #f9fafb;
    border-radius: 0.5rem;
    margin-bottom: 0.5rem;
}
</style>
