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
                    <el-option label="주간" value="week" />
                    <el-option label="월" value="month" />
                    <el-option label="년도" value="year" />
                    <el-option label="전체" value="all" />
                    <el-option label="기간 설정" value="custom" />
                </el-select>

                <!-- 주간 선택 -->
                <div v-if="selectedPeriod === 'week'" class="flex gap-2">
                    <el-select v-model="selectedYear" @change="handleYearChange" class="w-20" placeholder="년도">
                        <el-option v-for="year in availableYears" :key="year" :label="year + '년'" :value="year" />
                    </el-select>
                    <el-select v-model="selectedWeek" @change="handleWeekChange" class="w-24" placeholder="주차">
                        <el-option v-for="week in availableWeeks" :key="week.value" :label="week.label"
                            :value="week.value" />
                    </el-select>
                </div>

                <!-- 월 선택 -->
                <div v-if="selectedPeriod === 'month'" class="flex gap-2">
                    <el-select v-model="selectedYear" @change="handleYearChange" class="w-20" placeholder="년도">
                        <el-option v-for="year in availableYears" :key="year" :label="year + '년'" :value="year" />
                    </el-select>
                    <el-select v-model="selectedMonth" @change="handleMonthChange" class="w-20" placeholder="월">
                        <el-option v-for="month in availableMonths" :key="month.value" :label="month.label"
                            :value="month.value" />
                    </el-select>
                </div>

                <!-- 년도 선택 -->
                <div v-if="selectedPeriod === 'year'" class="flex gap-2">
                    <el-select v-model="selectedYear" @change="handleYearChange" class="w-20" placeholder="년도">
                        <el-option v-for="year in availableYears" :key="year" :label="year + '년'" :value="year" />
                    </el-select>
                </div>

                <!-- 커스텀 날짜 선택 -->
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

        <!-- 키워드 상세 영역 (페이지 내 표시) -->
        <div v-if="selectedCategory" class="mb-8 keyword-detail-section">
            <div class="card">
                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-xl font-bold text-gray-800">{{ selectedCategory.category_name }} 키워드 상세</h2>
                    <el-button @click="closeKeywordDetail" size="small" circle>×</el-button>
                </div>

                <div v-if="keywordStatistics">
                    <!-- 모바일: 세로 레이아웃, 데스크톱: 가로 레이아웃 -->
                    <div :class="isMobile ? 'space-y-6' : 'grid grid-cols-1 md:grid-cols-2 gap-6'" class="mb-6">
                        <div class="keyword-chart" :class="{ 'mobile-chart': isMobile }">
                            <Doughnut v-if="keywordChartData.datasets[0].data.length > 0" :data="keywordChartData"
                                :options="keywordChartOptions" />
                            <div v-else class="empty-chart">
                                <PieChart class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                                <p class="text-gray-500">키워드 데이터가 없습니다</p>
                            </div>
                        </div>
                        <div class="keyword-list" :class="{ 'mobile-keyword-list': isMobile }">
                            <!-- 키워드 리스트 헤더 -->
                            <div class="flex items-center justify-between mb-3">
                                <div class="flex items-center space-x-2">
                                    <!-- 선택된 키워드가 있을 때 전체 보기 버튼 표시 -->
                                    <el-button v-if="selectedKeywordIndex !== null" size="small"
                                        @click="selectedKeywordIndex = null" type="info">
                                        📋 전체 보기
                                    </el-button>
                                    <span v-else class="text-sm font-medium text-gray-700">키워드 목록</span>
                                </div>

                                <!-- 키워드 정렬 버튼 -->
                                <el-button size="small" @click="toggleKeywordSortOrder" type="default">
                                    {{ keywordSortOrder === 'desc' ? '💰 높은순' : '💸 낮은순' }}
                                </el-button>
                            </div>

                            <!-- 키워드 리스트 -->
                            <div v-for="keyword in filteredKeywords" :key="keyword.keyword_id || keyword.keyword_name"
                                class="keyword-item">
                                <div class="flex items-center justify-between">
                                    <span class="font-medium">{{ keyword.keyword_name || '키워드' }}</span>
                                    <span class="font-bold">{{ formatMoney(keyword.total_amount || 0) }}원</span>
                                </div>
                                <div class="flex items-center justify-between text-sm text-gray-500">
                                    <span>{{ keyword.count || 0 }}건</span>
                                    <span>{{ (keyword.percentage || 0).toFixed(1) }}%</span>
                                </div>
                            </div>

                            <!-- 선택된 키워드가 있을 때 안내 텍스트 -->
                            <div v-if="selectedKeywordIndex !== null && filteredKeywords.length > 0"
                                class="mt-3 p-2 bg-blue-50 rounded text-sm text-blue-600">
                                💡 차트를 클릭하여 다른 키워드를 선택하거나 전체 보기를 클릭하세요.
                            </div>
                        </div>
                    </div>
                </div>

                <div v-if="loadingKeywords" class="text-center py-8">
                    <div class="spinner mx-auto"></div>
                    <p class="text-gray-600 mt-2">키워드 데이터를 불러오는 중...</p>
                </div>
            </div>
        </div>
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
        const selectedCategory = ref(null);
        const loadingKeywords = ref(false);
        const selectedKeywordIndex = ref(null);
        const keywordSortOrder = ref('desc');

        // 현재 주차를 계산하는 헬퍼 함수
        function getCurrentWeek() {
            const now = new Date();
            const start = new Date(now.getFullYear(), 0, 1);
            const days = Math.floor((now - start) / (24 * 60 * 60 * 1000));
            return Math.ceil((days + start.getDay() + 1) / 7);
        }

        // 새로운 기간 선택 변수들
        const selectedYear = ref(new Date().getFullYear());
        const selectedMonth = ref(new Date().getMonth() + 1);
        const selectedWeek = ref(getCurrentWeek());

        // 모바일 감지
        const isMobile = computed(() => {
            if (typeof window === 'undefined') return false;
            return window.innerWidth <= 768;
        });

        const statistics = computed(() => statisticsStore.statistics);
        const keywordStatistics = computed(() => statisticsStore.keywordStatistics);
        const users = computed(() => userStore.users || []);
        const budgetUsages = computed(() => statistics.value?.budget_usages || []);

        // 사용 가능한 년도 목록 (현재 년도 기준 ±5년)
        const availableYears = computed(() => {
            const currentYear = new Date().getFullYear();
            const years = [];
            for (let year = currentYear - 5; year <= currentYear + 2; year++) {
                years.push(year);
            }
            return years.reverse(); // 최신 년도부터 표시
        });

        // 사용 가능한 월 목록
        const availableMonths = computed(() => {
            return [
                { value: 1, label: '1월' },
                { value: 2, label: '2월' },
                { value: 3, label: '3월' },
                { value: 4, label: '4월' },
                { value: 5, label: '5월' },
                { value: 6, label: '6월' },
                { value: 7, label: '7월' },
                { value: 8, label: '8월' },
                { value: 9, label: '9월' },
                { value: 10, label: '10월' },
                { value: 11, label: '11월' },
                { value: 12, label: '12월' }
            ];
        });

        // 사용 가능한 주차 목록 (선택된 년도 기준)
        const availableWeeks = computed(() => {
            const weeks = [];
            // 해당 년도의 주차 수 계산 (대략 52-53주)
            const weeksInYear = 53; // 최대 53주
            for (let week = 1; week <= weeksInYear; week++) {
                weeks.push({
                    value: week,
                    label: `${week}주차`
                });
            }
            return weeks;
        });

        // 선택된 키워드 리스트 (차트 클릭 시)
        const filteredKeywords = computed(() => {
            if (!keywordStatistics.value?.keywords) return [];

            let keywords = [...keywordStatistics.value.keywords];

            // 정렬 적용
            keywords.sort((a, b) => {
                const amountA = a.total_amount || 0;
                const amountB = b.total_amount || 0;
                return keywordSortOrder.value === 'desc'
                    ? amountB - amountA
                    : amountA - amountB;
            });

            // 선택된 키워드만 필터링
            if (selectedKeywordIndex.value !== null) {
                // 정렬된 배열에서 원래 선택된 키워드를 찾아야 함
                const originalKeyword = keywordStatistics.value.keywords[selectedKeywordIndex.value];
                const selectedKeyword = keywords.find(k =>
                    (k.keyword_id && k.keyword_id === originalKeyword?.keyword_id) ||
                    (k.keyword_name === originalKeyword?.keyword_name)
                );
                return selectedKeyword ? [selectedKeyword] : [];
            }

            return keywords;
        });

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

        // 키워드 차트 데이터 (정렬된 순서로)
        const keywordChartData = computed(() => {
            if (!keywordStatistics.value?.keywords || !Array.isArray(keywordStatistics.value.keywords)) {
                return {
                    labels: [],
                    datasets: [{
                        data: [],
                        backgroundColor: [],
                        borderWidth: 0
                    }]
                };
            }

            // 키워드를 정렬
            const sortedKeywords = [...keywordStatistics.value.keywords].sort((a, b) => {
                const amountA = a.total_amount || 0;
                const amountB = b.total_amount || 0;
                return keywordSortOrder.value === 'desc'
                    ? amountB - amountA
                    : amountA - amountB;
            });

            // 색상 배열 생성 (다양한 색상)
            const colors = [
                '#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6',
                '#06b6d4', '#84cc16', '#f97316', '#ec4899', '#6366f1'
            ];

            return {
                labels: sortedKeywords.map(keyword => keyword?.keyword_name || '키워드'),
                datasets: [{
                    data: sortedKeywords.map(keyword => keyword?.total_amount || 0),
                    backgroundColor: sortedKeywords.map((_, index) => colors[index % colors.length]),
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

        const keywordChartOptions = computed(() => ({
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'bottom',
                    labels: {
                        padding: isMobile.value ? 8 : 15,
                        usePointStyle: true,
                        font: {
                            size: isMobile.value ? 10 : 11
                        },
                        maxWidth: isMobile.value ? 120 : undefined,
                        generateLabels: function (chart) {
                            const labels = ChartJS.defaults.plugins.legend.labels.generateLabels(chart);
                            if (isMobile.value && labels && Array.isArray(labels)) {
                                return labels.map(label => ({
                                    ...label,
                                    text: (label.text && typeof label.text === 'string' && label.text.length > 8)
                                        ? label.text.substring(0, 8) + '...'
                                        : (label.text || '키워드')
                                }));
                            }
                            return labels || [];
                        }
                    }
                },
                tooltip: {
                    titleFont: {
                        size: isMobile.value ? 12 : 14
                    },
                    bodyFont: {
                        size: isMobile.value ? 11 : 13
                    },
                    padding: isMobile.value ? 8 : 12,
                    callbacks: {
                        label: function (context) {
                            const value = context.parsed || 0;
                            const total = context.dataset.data.reduce((a, b) => (a || 0) + (b || 0), 0);
                            const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : '0.0';
                            const label = context.label || '키워드';
                            return `${label}: ${formatMoney(value)}원 (${percentage}%)`;
                        }
                    }
                }
            },
            layout: {
                padding: {
                    top: isMobile.value ? 10 : 20,
                    bottom: isMobile.value ? 10 : 20,
                    left: isMobile.value ? 10 : 20,
                    right: isMobile.value ? 10 : 20
                }
            },
            onClick: (event, elements) => {
                if (elements.length > 0) {
                    const clickedIndex = elements[0].index;

                    // 정렬된 키워드 배열에서 클릭된 키워드 찾기
                    const sortedKeywords = [...(keywordStatistics.value?.keywords || [])].sort((a, b) => {
                        const amountA = a.total_amount || 0;
                        const amountB = b.total_amount || 0;
                        return keywordSortOrder.value === 'desc'
                            ? amountB - amountA
                            : amountA - amountB;
                    });

                    const clickedKeyword = sortedKeywords[clickedIndex];
                    if (clickedKeyword) {
                        // 원본 배열에서의 인덱스 찾기
                        const originalIndex = keywordStatistics.value?.keywords.findIndex(k =>
                            (k.keyword_id && k.keyword_id === clickedKeyword.keyword_id) ||
                            (k.keyword_name === clickedKeyword.keyword_name)
                        );

                        // 같은 키워드를 다시 클릭하면 전체 보기로 돌아감
                        if (selectedKeywordIndex.value === originalIndex) {
                            selectedKeywordIndex.value = null;
                        } else {
                            selectedKeywordIndex.value = originalIndex;
                        }
                    }
                }
            }
        }));

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

                // 기간별 파라미터 추가
                if (selectedPeriod.value === 'week') {
                    params.year = selectedYear.value;
                    params.week = selectedWeek.value;
                } else if (selectedPeriod.value === 'month') {
                    params.year = selectedYear.value;
                    params.month = selectedMonth.value;
                } else if (selectedPeriod.value === 'year') {
                    params.year = selectedYear.value;
                } else if (selectedPeriod.value === 'custom' && customStartDate.value && customEndDate.value) {
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

            // 기간 타입이 변경될 때 기본값으로 초기화
            if (selectedPeriod.value === 'week') {
                selectedYear.value = new Date().getFullYear();
                selectedWeek.value = getCurrentWeek();
            } else if (selectedPeriod.value === 'month') {
                selectedYear.value = new Date().getFullYear();
                selectedMonth.value = new Date().getMonth() + 1;
            } else if (selectedPeriod.value === 'year') {
                selectedYear.value = new Date().getFullYear();
            }

            loadStatistics();
        };

        // 년도 변경 핸들러
        const handleYearChange = () => {
            loadStatistics();
        };

        // 월 변경 핸들러
        const handleMonthChange = () => {
            loadStatistics();
        };

        // 주차 변경 핸들러
        const handleWeekChange = () => {
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

        // 키워드 정렬 순서 토글
        const toggleKeywordSortOrder = () => {
            keywordSortOrder.value = keywordSortOrder.value === 'desc' ? 'asc' : 'desc';
            // 키워드 정렬이 변경되면 선택된 키워드 해제
            selectedKeywordIndex.value = null;
        };

        // 카테고리 상세 보기
        const showCategoryDetail = async (category) => {
            selectedCategory.value = category;
            loadingKeywords.value = true;

            try {
                const params = {
                    category_id: category.category_id,
                    type: selectedPeriod.value,
                    category: selectedType.value
                };

                // 기간별 파라미터 추가
                if (selectedPeriod.value === 'week') {
                    params.year = selectedYear.value;
                    params.week = selectedWeek.value;
                } else if (selectedPeriod.value === 'month') {
                    params.year = selectedYear.value;
                    params.month = selectedMonth.value;
                } else if (selectedPeriod.value === 'year') {
                    params.year = selectedYear.value;
                } else if (selectedPeriod.value === 'custom' && customStartDate.value && customEndDate.value) {
                    params.start_date = customStartDate.value;
                    params.end_date = customEndDate.value;
                }

                await statisticsStore.fetchKeywordStatistics(params);

                // 스크롤을 키워드 상세 영역으로 이동
                setTimeout(() => {
                    const keywordSection = document.querySelector('.keyword-detail-section');
                    if (keywordSection) {
                        keywordSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
                    }
                }, 100);
            } catch (error) {
                ElMessage.error('키워드 통계를 불러오는데 실패했습니다.');
            } finally {
                loadingKeywords.value = false;
            }
        };

        // 키워드 상세 닫기
        const closeKeywordDetail = () => {
            selectedCategory.value = null;
            selectedKeywordIndex.value = null;
            keywordSortOrder.value = 'desc'; // 정렬 순서도 초기화
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
            selectedCategory,
            loadingKeywords,
            selectedKeywordIndex,
            keywordSortOrder,

            // 새로운 기간 선택 변수들
            selectedYear,
            selectedMonth,
            selectedWeek,

            statistics,
            keywordStatistics,
            users,
            budgetUsages,
            filteredKeywords,
            sortedCategories,
            chartData,
            keywordChartData,
            chartOptions,
            keywordChartOptions,

            // 새로운 computed 속성들
            availableYears,
            availableMonths,
            availableWeeks,

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

            // 새로운 핸들러들
            handleYearChange,
            handleMonthChange,
            handleWeekChange,

            toggleSortOrder,
            toggleKeywordSortOrder,
            showCategoryDetail,
            closeKeywordDetail,
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

/* 모바일 키워드 차트 최적화 */
.mobile-chart {
    height: 200px !important;
    margin: 0 auto;
    max-width: 280px;
}

.mobile-keyword-list {
    margin-top: 1rem;
}

.mobile-keyword-list .keyword-item {
    padding: 0.5rem;
    font-size: 0.875rem;
    margin-bottom: 0.375rem;
}

.mobile-keyword-list .keyword-item .font-medium {
    font-size: 0.875rem;
}

.mobile-keyword-list .keyword-item .font-bold {
    font-size: 0.875rem;
}

/* 키워드 상세 섹션 최적화 */
.keyword-detail-section {
    scroll-margin-top: 2rem;
}

@media (max-width: 768px) {
    .keyword-chart {
        height: 180px !important;
    }

    .mobile-chart {
        height: 180px !important;
        max-width: 100%;
        margin: 0;
    }

    .keyword-detail-section .card {
        margin: 0 1rem;
    }
}
</style>
