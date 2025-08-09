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
import { Doughnut } from 'vue-chartjs';
import {
    Chart as ChartJS,
    ArcElement,
    Tooltip,
    Legend
} from 'chart.js';
import { useStatisticsStore } from '../stores/statisticsStore';

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
        Doughnut
    },
    setup() {
        const statisticsStore = useStatisticsStore();

        const selectedType = ref('out');
        const selectedPeriod = ref('month');
        const customStartDate = ref(null);
        const customEndDate = ref(null);
        const sortOrder = ref('desc');
        const keywordDialogVisible = ref(false);
        const selectedCategory = ref(null);
        const loadingKeywords = ref(false);

        const statistics = computed(() => statisticsStore.statistics);
        const keywordStatistics = computed(() => statisticsStore.keywordStatistics);

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

        // 통계 데이터 로드
        const loadStatistics = async () => {
            try {
                const params = {
                    type: selectedPeriod.value,
                    category: selectedType.value
                };

                if (selectedPeriod.value === 'custom' && customStartDate.value && customEndDate.value) {
                    params.start_date = customStartDate.value;
                    params.end_date = customEndDate.value;
                }

                await statisticsStore.fetchStatistics(params);
            } catch (error) {
                ElMessage.error('통계 데이터를 불러오는데 실패했습니다.');
            }
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

        onMounted(() => {
            loadStatistics();
        });

        return {
            selectedType,
            selectedPeriod,
            customStartDate,
            customEndDate,
            sortOrder,
            keywordDialogVisible,
            selectedCategory,
            loadingKeywords,

            statistics,
            keywordStatistics,
            sortedCategories,
            chartData,
            keywordChartData,
            chartOptions,
            keywordChartOptions,

            formatMoney,
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
    @apply bg-white p-6 rounded-xl border border-gray-200 shadow-sm;
}

.stat-card-red {
    @apply border-red-200 bg-gradient-to-r from-red-50 to-red-100;
}

.stat-card-green {
    @apply border-green-200 bg-gradient-to-r from-green-50 to-green-100;
}

.stat-card-blue {
    @apply border-blue-200 bg-gradient-to-r from-blue-50 to-blue-100;
}

.stat-card-purple {
    @apply border-purple-200 bg-gradient-to-r from-purple-50 to-purple-100;
}

.stat-icon {
    @apply w-12 h-12 rounded-xl flex items-center justify-center;
}

.stat-card-red .stat-icon {
    @apply bg-red-100 text-red-600;
}

.stat-card-green .stat-icon {
    @apply bg-green-100 text-green-600;
}

.stat-label {
    @apply text-sm font-medium text-gray-600;
}

.stat-value {
    @apply text-2xl font-bold text-gray-900;
}

.category-item {
    @apply p-4 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors duration-200;
}

.rank-badge {
    @apply w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold text-white;
}

.rank-1 {
    @apply bg-yellow-500;
}

.rank-2 {
    @apply bg-gray-400;
}

.rank-3 {
    @apply bg-orange-500;
}

.rank-badge:not(.rank-1):not(.rank-2):not(.rank-3) {
    @apply bg-blue-500;
}

.progress-bar {
    @apply w-full bg-gray-200 rounded-full h-1.5;
}

.progress-fill {
    @apply h-1.5 rounded-full transition-all duration-300;
}

.empty-chart {
    @apply flex flex-col items-center justify-center h-full text-center;
}

.empty-state {
    @apply text-center py-8;
}

.keyword-chart {
    height: 250px;
}

.keyword-item {
    @apply p-3 bg-gray-50 rounded-lg mb-2;
}
</style>
