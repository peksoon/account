<template>
    <div class="statistics-dashboard">
        <!-- 헤더 -->
        <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between mb-8">
            <div>
                <h1 class="text-3xl font-bold text-gradient mb-2">📊 통계 대시보드</h1>
                <p class="text-gray-600">가계부 통계를 한눈에 확인하세요</p>
            </div>

            <!-- 우측 상단 도구 버튼 -->
            <div class="flex items-center gap-2 mt-4 lg:mt-0">
                <el-button @click="openExportData" type="success" size="small" class="text-xs">
                    <Download class="w-3 h-3 mr-1" />
                    데이터 내보내기
                </el-button>
            </div>
        </div>

        <!-- 필터 컨트롤 -->
        <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between mb-8">
            <div class="flex flex-col sm:flex-row gap-4 w-full">
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

        <!-- 차트 유형 선택 (지출일 때만 표시) -->
        <div v-if="selectedType === 'out'" class="mb-6">
            <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
                <button @click="chartViewType = 'category'" :class="[
                    'chart-type-button',
                    chartViewType === 'category' ? 'active' : ''
                ]">
                    <span class="text-2xl mb-1">📊</span>
                    <span class="text-sm font-medium">카테고리별</span>
                </button>
                <button @click="chartViewType = 'expense_type'" :class="[
                    'chart-type-button',
                    chartViewType === 'expense_type' ? 'active' : ''
                ]">
                    <span class="text-2xl mb-1">💰</span>
                    <span class="text-sm font-medium">고정/변동</span>
                </button>
                <button @click="chartViewType = 'payment_method'" :class="[
                    'chart-type-button',
                    chartViewType === 'payment_method' ? 'active' : ''
                ]">
                    <span class="text-2xl mb-1">💳</span>
                    <span class="text-sm font-medium">결제수단별</span>
                </button>
                <button @click="chartViewType = 'user'" :class="[
                    'chart-type-button',
                    chartViewType === 'user' ? 'active' : ''
                ]">
                    <span class="text-2xl mb-1">👤</span>
                    <span class="text-sm font-medium">사용자별</span>
                </button>
            </div>
        </div>

        <!-- 카테고리별 지출 차트 (수입일 때 항상 표시, 지출일 때는 선택 시만 표시) -->
        <div v-if="selectedType === 'in' || (selectedType === 'out' && chartViewType === 'category')"
            class="grid grid-cols-1 xl:grid-cols-2 gap-8 mb-8">
            <!-- 도넛 차트 -->
            <div class="card">
                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-xl font-bold text-gray-800">카테고리별 {{ selectedType === 'out' ? '지출' : '수입' }}</h2>
                    <span class="text-sm text-gray-500">{{ statistics?.period }}</span>
                </div>

                <div class="chart-container">
                    <Doughnut v-if="chartData.datasets[0].data.length > 0" :key="statistics?.period + selectedType"
                        :data="chartData" :options="chartOptions" @click="handleChartClick" />
                    <div v-else class="empty-chart">
                        <PieChart class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                        <p class="text-gray-500">표시할 데이터가 없습니다</p>
                    </div>
                </div>
            </div>

            <!-- 카테고리별 순위 리스트 -->
            <div class="card">
                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-xl font-bold text-gray-800">카테고리별 {{ selectedType === 'out' ? '지출' : '수입' }} 순위</h2>
                    <el-button size="small" @click="toggleSortOrder">
                        {{ sortOrder === 'desc' ? '↓ 높은순' : '↑ 낮은순' }}
                    </el-button>
                </div>

                <div class="space-y-3">
                    <div v-for="(category, index) in sortedCategories" :key="category.category_id" class="category-item"
                        @click="showCategoryDetail(category)">
                        <div class="flex items-center">
                            <div class="rank-badge"
                                :style="{ backgroundColor: getCategoryColor(category.category_name) }">
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

        <!-- 고정/변동 지출 분석 (지출일 때 선택 시만 표시) -->
        <div v-if="selectedType === 'out' && chartViewType === 'expense_type'" class="mb-8">
            <!-- 고정 vs 변동 비교 차트 -->
            <div class="card mb-6">
                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-xl font-bold text-gray-800">📊 고정 vs 변동 지출 비교</h2>
                    <span class="text-sm text-gray-500">{{ statistics?.period }}</span>
                </div>

                <div class="comparison-chart-container">
                    <Bar v-if="fixedExpenseTotal + variableExpenseTotal > 0" :data="expenseComparisonChartData"
                        :options="expenseComparisonChartOptions" />
                    <div v-else class="empty-chart">
                        <BarChart class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                        <p class="text-gray-500">지출 데이터가 없습니다</p>
                    </div>
                </div>

                <!-- 요약 정보 -->
                <div class="grid grid-cols-2 gap-4 mt-6">
                    <div class="text-center p-4 bg-blue-50 rounded-lg">
                        <p class="text-sm text-blue-600 font-medium mb-1">📌 고정 지출</p>
                        <p class="text-2xl font-bold text-blue-700">{{ formatMoney(fixedExpenseTotal) }}원</p>
                        <p class="text-xs text-blue-500 mt-1">
                            {{ fixedExpenseTotal + variableExpenseTotal > 0
                                ? ((fixedExpenseTotal / (fixedExpenseTotal + variableExpenseTotal)) * 100).toFixed(1)
                                : 0 }}%
                        </p>
                    </div>
                    <div class="text-center p-4 bg-green-50 rounded-lg">
                        <p class="text-sm text-green-600 font-medium mb-1">💳 변동 지출</p>
                        <p class="text-2xl font-bold text-green-700">{{ formatMoney(variableExpenseTotal) }}원</p>
                        <p class="text-xs text-green-500 mt-1">
                            {{ fixedExpenseTotal + variableExpenseTotal > 0
                                ? ((variableExpenseTotal / (fixedExpenseTotal + variableExpenseTotal)) * 100).toFixed(1)
                                : 0 }}%
                        </p>
                    </div>
                </div>
            </div>

            <!-- 카테고리별 상세 분석 -->
            <div class="card">
                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-xl font-bold text-gray-800">💰 카테고리별 상세 분석</h2>
                    <span class="text-sm text-gray-500">{{ statistics?.period }}</span>
                </div>

                <!-- 고정/변동 선택 탭 -->
                <el-tabs v-model="expenseTypeTab" class="mb-4">
                    <el-tab-pane label="💳 변동 지출" name="variable"></el-tab-pane>
                    <el-tab-pane label="📌 고정 지출" name="fixed"></el-tab-pane>
                </el-tabs>

                <div class="grid grid-cols-1 xl:grid-cols-2 gap-6">
                    <!-- 고정/변동 도넛 차트 -->
                    <div>
                        <div class="chart-container">
                            <Doughnut v-if="expenseTypeChartData.datasets[0].data.length > 0"
                                :key="statistics?.period + expenseTypeTab" :data="expenseTypeChartData"
                                :options="expenseTypeChartOptions" @click="handleExpenseTypeChartClick" />
                            <div v-else class="empty-chart">
                                <PieChart class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                                <p class="text-gray-500">{{ expenseTypeTab === 'fixed' ? '고정 지출' : '변동 지출' }} 데이터가 없습니다
                                </p>
                            </div>
                        </div>
                    </div>

                    <!-- 고정/변동 카테고리 리스트 -->
                    <div>
                        <div class="flex items-center justify-between mb-4">
                            <h3 class="font-semibold text-gray-700">
                                {{ expenseTypeTab === 'fixed' ? '고정 지출' : '변동 지출' }} 카테고리
                            </h3>
                            <span class="text-sm text-gray-500">
                                총 {{ formatMoney(expenseTypeTotal) }}원
                            </span>
                        </div>

                        <div class="space-y-3">
                            <div v-for="(category, index) in expenseTypeCategories" :key="category.category_id"
                                class="category-item cursor-pointer" @click="showExpenseTypeCategoryDetail(category)">
                                <div class="flex items-center">
                                    <div class="rank-badge"
                                        :style="{ backgroundColor: getCategoryColor(category.category_name) }">
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
                                            <span class="text-sm font-medium text-red-600">
                                                {{ category.percentage.toFixed(1) }}%
                                            </span>
                                        </div>
                                    </div>
                                </div>
                                <div class="progress-bar mt-2">
                                    <div class="progress-fill bg-red-500" :style="`width: ${category.percentage}%`">
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div v-if="expenseTypeCategories.length === 0" class="empty-state">
                            <Folder class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                            <p class="text-gray-500">{{ expenseTypeTab === 'fixed' ? '고정 지출' : '변동 지출' }} 데이터가 없습니다</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 결제수단별 지출 분석 (지출일 때 선택 시만 표시) -->
        <div v-if="selectedType === 'out' && chartViewType === 'payment_method'" class="mb-8">
            <div class="grid grid-cols-1 xl:grid-cols-2 gap-8">
                <!-- 도넛 차트 -->
                <div class="card">
                    <div class="flex items-center justify-between mb-6">
                        <h2 class="text-xl font-bold text-gray-800">💳 결제수단별 지출</h2>
                        <span class="text-sm text-gray-500">{{ statistics?.period }}</span>
                    </div>

                    <div class="chart-container">
                        <Doughnut v-if="paymentMethodChartData.datasets[0].data.length > 0"
                            :key="statistics?.period + 'payment'" :data="paymentMethodChartData"
                            :options="paymentMethodChartOptions" />
                        <div v-else class="empty-chart">
                            <PieChart class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                            <p class="text-gray-500">결제수단 데이터가 없습니다</p>
                        </div>
                    </div>
                </div>

                <!-- 결제수단별 순위 리스트 -->
                <div class="card">
                    <div class="flex items-center justify-between mb-6">
                        <h2 class="text-xl font-bold text-gray-800">💳 결제수단별 지출 순위</h2>
                        <el-button size="small" @click="togglePaymentMethodSortOrder">
                            {{ paymentMethodSortOrder === 'desc' ? '↓ 높은순' : '↑ 낮은순' }}
                        </el-button>
                    </div>

                    <div class="space-y-3">
                        <div v-for="(payment, index) in sortedPaymentMethods" :key="payment.payment_method_id"
                            class="category-item cursor-pointer hover:shadow-lg transition-shadow"
                            :class="{ 'ring-2 ring-blue-500': selectedPaymentMethod?.payment_method_id === payment.payment_method_id }"
                            @click="showPaymentMethodDetail(payment)">
                            <div class="flex items-center">
                                <div class="rank-badge" :style="{ backgroundColor: getPaymentMethodColor(index) }">
                                    {{ index + 1 }}
                                </div>
                                <div class="ml-3 flex-1">
                                    <div class="flex items-center justify-between">
                                        <span class="font-medium text-gray-900">{{ payment.payment_method_name }}</span>
                                        <span class="font-bold text-gray-800">{{ formatMoney(payment.total_amount)
                                            }}원</span>
                                    </div>
                                    <div class="flex items-center justify-between mt-1">
                                        <span class="text-sm text-gray-500">{{ payment.count }}건</span>
                                        <span class="text-sm font-medium text-red-600">
                                            {{ payment.percentage.toFixed(1) }}%
                                        </span>
                                    </div>
                                </div>
                            </div>
                            <div class="progress-bar mt-2">
                                <div class="progress-fill bg-red-500" :style="`width: ${payment.percentage}%`"></div>
                            </div>
                        </div>
                    </div>

                    <div v-if="sortedPaymentMethods.length === 0" class="empty-state">
                        <Folder class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                        <p class="text-gray-500">해당 기간에 결제수단 데이터가 없습니다</p>
                    </div>
                </div>
            </div>
        </div>

        <!-- 결제수단별 지출 내역 (결제수단 선택 시 표시) -->
        <div v-if="selectedPaymentMethod" id="payment-method-detail" class="mb-8">
            <div class="card">
                <div class="flex items-center justify-between mb-6">
                    <div>
                        <h2 class="text-xl font-bold text-gray-800">
                            💳 {{ selectedPaymentMethod.payment_method_name }} - 지출 내역
                        </h2>
                        <p class="text-sm text-gray-500 mt-1">
                            총 {{ formatMoney(selectedPaymentMethod.total_amount) }}원 ({{ selectedPaymentMethod.count
                            }}건)
                        </p>
                    </div>
                    <el-button size="small" @click="closePaymentMethodDetail">
                        <X class="w-4 h-4" />
                        닫기
                    </el-button>
                </div>

                <div v-if="loadingPaymentMethodAccounts" class="text-center py-8">
                    <div class="spinner"></div>
                    <p class="text-gray-500 mt-4">지출 내역을 불러오는 중...</p>
                </div>

                <div v-else-if="paymentMethodAccounts.length > 0" class="overflow-x-auto">
                    <table class="min-w-full divide-y divide-gray-200">
                        <thead class="bg-gray-50">
                            <tr>
                                <th
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    날짜</th>
                                <th
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    카테고리</th>
                                <th
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    키워드</th>
                                <th
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    메모</th>
                                <th
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    사용자</th>
                                <th
                                    class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    금액</th>
                            </tr>
                        </thead>
                        <tbody class="bg-white divide-y divide-gray-200">
                            <tr v-for="account in paymentMethodAccounts" :key="account.uuid" class="hover:bg-gray-50">
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    {{ formatDate(account.date) }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm">
                                    <span
                                        class="px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
                                        {{ account.category_name }}
                                    </span>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                    {{ account.keyword_name || '-' }}
                                </td>
                                <td class="px-6 py-4 text-sm text-gray-600 max-w-xs truncate">
                                    {{ account.memo || '-' }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                    {{ account.user }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-right font-bold text-red-600">
                                    {{ formatMoney(account.money) }}원
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div v-else class="empty-state">
                    <Folder class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                    <p class="text-gray-500">해당 결제수단으로 결제한 내역이 없습니다</p>
                </div>
            </div>
        </div>

        <!-- 사용자별 지출 분석 (지출일 때 선택 시만 표시) -->
        <div v-if="selectedType === 'out' && chartViewType === 'user'" class="mb-8">
            <div class="grid grid-cols-1 xl:grid-cols-2 gap-8">
                <!-- 도넛 차트 -->
                <div class="card">
                    <div class="flex items-center justify-between mb-6">
                        <h2 class="text-xl font-bold text-gray-800">👤 사용자별 지출</h2>
                        <span class="text-sm text-gray-500">{{ statistics?.period }}</span>
                    </div>

                    <div class="chart-container">
                        <Doughnut v-if="userChartData.datasets[0].data.length > 0" :key="statistics?.period + 'user'"
                            :data="userChartData" :options="userChartOptions" />
                        <div v-else class="empty-chart">
                            <PieChart class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                            <p class="text-gray-500">사용자 데이터가 없습니다</p>
                        </div>
                    </div>
                </div>

                <!-- 사용자별 순위 리스트 -->
                <div class="card">
                    <div class="flex items-center justify-between mb-6">
                        <h2 class="text-xl font-bold text-gray-800">👤 사용자별 지출 순위</h2>
                        <el-button size="small" @click="toggleUserSortOrder">
                            {{ userSortOrder === 'desc' ? '↓ 높은순' : '↑ 낮은순' }}
                        </el-button>
                    </div>

                    <div class="space-y-3">
                        <div v-for="(user, index) in sortedUsers" :key="user.user_name"
                            class="category-item cursor-pointer hover:shadow-lg transition-shadow"
                            :class="{ 'ring-2 ring-blue-500': selectedUserDetail?.user_name === user.user_name }"
                            @click="showUserDetail(user)">
                            <div class="flex items-center">
                                <div class="rank-badge" :style="{ backgroundColor: getUserColor(index) }">
                                    {{ index + 1 }}
                                </div>
                                <div class="ml-3 flex-1">
                                    <div class="flex items-center justify-between">
                                        <span class="font-medium text-gray-900">{{ user.user_name }}</span>
                                        <span class="font-bold text-gray-800">{{ formatMoney(user.total_amount)
                                            }}원</span>
                                    </div>
                                    <div class="flex items-center justify-between mt-1">
                                        <span class="text-sm text-gray-500">{{ user.count }}건</span>
                                        <span class="text-sm font-medium text-red-600">
                                            {{ user.percentage.toFixed(1) }}%
                                        </span>
                                    </div>
                                </div>
                            </div>
                            <div class="progress-bar mt-2">
                                <div class="progress-fill bg-red-500" :style="`width: ${user.percentage}%`"></div>
                            </div>
                        </div>
                    </div>

                    <div v-if="sortedUsers.length === 0" class="empty-state">
                        <Folder class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                        <p class="text-gray-500">해당 기간에 사용자 데이터가 없습니다</p>
                    </div>
                </div>
            </div>
        </div>

        <!-- 사용자별 지출 내역 (사용자 선택 시 표시) -->
        <div v-if="selectedUserDetail" id="user-detail" class="mb-8">
            <div class="card">
                <div class="flex items-center justify-between mb-6">
                    <div>
                        <h2 class="text-xl font-bold text-gray-800">
                            👤 {{ selectedUserDetail.user_name }} - 지출 내역
                        </h2>
                        <p class="text-sm text-gray-500 mt-1">
                            총 {{ formatMoney(selectedUserDetail.total_amount) }}원 ({{ selectedUserDetail.count }}건)
                        </p>
                    </div>
                    <el-button size="small" @click="closeUserDetail">
                        <X class="w-4 h-4" />
                        닫기
                    </el-button>
                </div>

                <div v-if="loadingUserAccounts" class="text-center py-8">
                    <div class="spinner"></div>
                    <p class="text-gray-500 mt-4">지출 내역을 불러오는 중...</p>
                </div>

                <div v-else-if="userAccounts.length > 0" class="overflow-x-auto">
                    <table class="min-w-full divide-y divide-gray-200">
                        <thead class="bg-gray-50">
                            <tr>
                                <th
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    날짜</th>
                                <th
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    카테고리</th>
                                <th
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    키워드</th>
                                <th
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    결제수단</th>
                                <th
                                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    메모</th>
                                <th
                                    class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    금액</th>
                            </tr>
                        </thead>
                        <tbody class="bg-white divide-y divide-gray-200">
                            <tr v-for="account in userAccounts" :key="account.uuid" class="hover:bg-gray-50">
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    {{ formatDate(account.date) }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm">
                                    <span
                                        class="px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
                                        {{ account.category_name }}
                                    </span>
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                    {{ account.keyword_name || '-' }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                    {{ account.payment_method_name || '-' }}
                                </td>
                                <td class="px-6 py-4 text-sm text-gray-600 max-w-xs truncate">
                                    {{ account.memo || '-' }}
                                </td>
                                <td class="px-6 py-4 whitespace-nowrap text-sm text-right font-bold text-red-600">
                                    {{ formatMoney(account.money) }}원
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div v-else class="empty-state">
                    <Folder class="w-16 h-16 text-gray-300 mx-auto mb-4" />
                    <p class="text-gray-500">해당 사용자의 지출 내역이 없습니다</p>
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
    Folder,
    Download,
    X
} from 'lucide-vue-next';
import { ElMessage } from 'element-plus';
import BudgetUsageDisplay from './BudgetUsageDisplay.vue';
import { Doughnut, Bar } from 'vue-chartjs';
import {
    Chart as ChartJS,
    ArcElement,
    Tooltip,
    Legend,
    CategoryScale,
    LinearScale,
    BarElement
} from 'chart.js';
import { useStatisticsStore } from '../stores/statisticsStore';
import { useUserStore } from '../stores/userStore';
import { useCategoryStore } from '../stores/categoryStore';
import { useApi } from '../composables/useApi';
import { useRouter } from 'vue-router';

// Chart.js 등록
ChartJS.register(ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement);

export default {
    name: 'StatisticsDashboard',
    components: {
        TrendingUp,
        TrendingDown,
        BarChart,
        Calculator,
        PieChart,
        Folder,
        Download,
        X,
        Doughnut,
        Bar,
        BudgetUsageDisplay
    },
    emits: ['close', 'open-budget-manager'],
    setup() {
        const statisticsStore = useStatisticsStore();
        const userStore = useUserStore();
        const categoryStore = useCategoryStore();
        const router = useRouter();
        const api = useApi();

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
        const expenseTypeTab = ref('variable'); // 고정/변동 지출 탭
        const chartViewType = ref('category'); // 차트 뷰 타입 (category, expense_type, payment_method, user)
        const paymentMethodSortOrder = ref('desc'); // 결제수단 정렬 순서
        const selectedPaymentMethod = ref(null); // 선택된 결제수단
        const paymentMethodAccounts = ref([]); // 결제수단별 지출 내역
        const loadingPaymentMethodAccounts = ref(false); // 로딩 상태
        const userSortOrder = ref('desc'); // 사용자 정렬 순서
        const selectedUserDetail = ref(null); // 선택된 사용자
        const userAccounts = ref([]); // 사용자별 지출 내역
        const loadingUserAccounts = ref(false); // 로딩 상태

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

        // 고정/변동 지출 카테고리 필터링 (카테고리 정보와 매칭)
        const expenseTypeCategories = computed(() => {
            if (!statistics.value?.categories || !categoryStore.categories) return [];

            const categories = statistics.value.categories
                .map(category => {
                    // category_id로 카테고리 정보 찾기
                    const categoryInfo = categoryStore.categories.find(c => c.id === category.category_id);

                    return {
                        ...category,
                        expense_type: categoryInfo?.expense_type || 'variable'
                    };
                })
                .filter(category => category.expense_type === expenseTypeTab.value);

            // 금액 기준으로 정렬
            return categories.sort((a, b) => b.total_amount - a.total_amount);
        });

        // 고정/변동 지출 총액
        const expenseTypeTotal = computed(() => {
            return expenseTypeCategories.value.reduce((sum, cat) => sum + cat.total_amount, 0);
        });

        // 고정 지출 총액
        const fixedExpenseTotal = computed(() => {
            if (!statistics.value?.categories || !categoryStore.categories) return 0;

            return statistics.value.categories
                .filter(category => {
                    const categoryInfo = categoryStore.categories.find(c => c.id === category.category_id);
                    return categoryInfo?.expense_type === 'fixed';
                })
                .reduce((sum, cat) => sum + cat.total_amount, 0);
        });

        // 변동 지출 총액
        const variableExpenseTotal = computed(() => {
            if (!statistics.value?.categories || !categoryStore.categories) return 0;

            return statistics.value.categories
                .filter(category => {
                    const categoryInfo = categoryStore.categories.find(c => c.id === category.category_id);
                    return categoryInfo?.expense_type === 'variable';
                })
                .reduce((sum, cat) => sum + cat.total_amount, 0);
        });

        // 고정 vs 변동 지출 비교 차트 데이터
        const expenseComparisonChartData = computed(() => {
            const fixed = fixedExpenseTotal.value;
            const variable = variableExpenseTotal.value;

            return {
                labels: ['고정 지출', '변동 지출'],
                datasets: [{
                    data: [fixed, variable],
                    backgroundColor: ['#3b82f6', '#10b981'], // 파란색(고정), 초록색(변동)
                    borderWidth: 0,
                    barThickness: 60,
                    maxBarThickness: 80
                }]
            };
        });

        // 고정 vs 변동 지출 비교 차트 옵션
        const expenseComparisonChartOptions = computed(() => {
            const total = fixedExpenseTotal.value + variableExpenseTotal.value;

            return {
                responsive: true,
                maintainAspectRatio: false,
                indexAxis: 'y', // 가로 막대 그래프
                plugins: {
                    legend: {
                        display: false
                    },
                    tooltip: {
                        callbacks: {
                            label: (context) => {
                                const value = context.parsed.x || 0;
                                const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : 0;
                                return `${formatMoney(value)}원 (${percentage}%)`;
                            }
                        }
                    }
                },
                scales: {
                    x: {
                        beginAtZero: true,
                        ticks: {
                            callback: function (value) {
                                return formatMoney(value) + '원';
                            }
                        }
                    },
                    y: {
                        ticks: {
                            font: {
                                size: 14,
                                weight: 'bold'
                            }
                        }
                    }
                }
            };
        });

        // 고정/변동 지출 차트 데이터
        const expenseTypeChartData = computed(() => {
            if (expenseTypeCategories.value.length === 0) {
                return {
                    labels: [],
                    datasets: [{
                        label: '',
                        data: [],
                        backgroundColor: [],
                        borderWidth: 0
                    }]
                };
            }

            const colors = expenseTypeCategories.value.map(cat => getCategoryColor(cat.category_name));

            return {
                labels: expenseTypeCategories.value.map(cat => cat.category_name),
                datasets: [{
                    label: '',
                    data: expenseTypeCategories.value.map(cat => cat.total_amount),
                    backgroundColor: colors,
                    borderWidth: 0,
                    hoverOffset: 10
                }]
            };
        });

        // 고정/변동 지출 차트 옵션
        const expenseTypeChartOptions = computed(() => {
            return {
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
                            label: (context) => {
                                const label = context.label || '';
                                const value = context.parsed || 0;
                                const total = context.dataset.data.reduce((a, b) => a + b, 0);
                                const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : 0;
                                return `${label}: ${formatMoney(value)}원 (${percentage}%)`;
                            }
                        }
                    }
                }
            };
        });

        // 정렬된 결제수단
        const sortedPaymentMethods = computed(() => {
            if (!statistics.value?.payment_methods) return [];

            const methods = [...statistics.value.payment_methods];
            return methods.sort((a, b) => {
                return paymentMethodSortOrder.value === 'desc'
                    ? b.total_amount - a.total_amount
                    : a.total_amount - b.total_amount;
            });
        });

        // 결제수단 차트 데이터
        const paymentMethodChartData = computed(() => {
            if (!statistics.value?.payment_methods || statistics.value.payment_methods.length === 0) {
                return {
                    labels: [],
                    datasets: [{
                        label: '',
                        data: [],
                        backgroundColor: [],
                        borderWidth: 0
                    }]
                };
            }

            // 결제수단별 색상 생성
            const colors = statistics.value.payment_methods.map((_, index) => getPaymentMethodColor(index));

            return {
                labels: statistics.value.payment_methods.map(item => item.payment_method_name),
                datasets: [{
                    label: '',
                    data: statistics.value.payment_methods.map(item => item.total_amount),
                    backgroundColor: colors,
                    borderWidth: 0,
                    hoverOffset: 10
                }]
            };
        });

        // 결제수단 차트 옵션
        const paymentMethodChartOptions = computed(() => {
            return {
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
                            label: (context) => {
                                const label = context.label || '';
                                const value = context.parsed || 0;
                                const total = context.dataset.data.reduce((a, b) => a + b, 0);
                                const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : 0;
                                return `${label}: ${formatMoney(value)}원 (${percentage}%)`;
                            }
                        }
                    }
                }
            };
        });

        // 차트 데이터
        const chartData = computed(() => {
            if (!statistics.value?.chart_data) {
                return {
                    labels: [],
                    datasets: [{
                        label: '',
                        data: [],
                        backgroundColor: [],
                        borderWidth: 0
                    }]
                };
            }

            // 백엔드에서 받은 색상 데이터 디버깅
            console.log('=== 카테고리 차트 색상 디버깅 ===');
            console.log('전체 statistics 데이터:', statistics.value);
            console.log('chart_data:', statistics.value.chart_data);

            const colors = statistics.value.chart_data.map(item => {
                console.log(`카테고리: ${item.label}, 색상: ${item.color}`);
                return item.color;
            });
            console.log('최종 색상 배열:', colors);
            console.log('=== 디버깅 끝 ===');

            const chartDataResult = {
                labels: statistics.value.chart_data.map(item => item.label),
                datasets: [{
                    label: '',
                    data: statistics.value.chart_data.map(item => item.value),
                    backgroundColor: [...colors], // 배열 복사로 Chart.js 인식 개선
                    borderWidth: 0,
                    hoverOffset: 10
                }]
            };

            console.log('Chart.js로 전달할 최종 데이터:', chartDataResult);
            return chartDataResult;
        });

        // 카테고리별 색상 매핑 함수
        const getCategoryColor = (categoryName) => {
            if (!statistics.value?.chart_data) return '#3b82f6'; // 기본 파란색

            const chartItem = statistics.value.chart_data.find(item => item.label === categoryName);
            return chartItem ? chartItem.color : '#3b82f6';
        };

        // 키워드 차트 데이터 (정렬된 순서로)
        const keywordChartData = computed(() => {
            if (!keywordStatistics.value?.keywords || !Array.isArray(keywordStatistics.value.keywords)) {
                return {
                    labels: [],
                    datasets: [{
                        label: '',
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

            // 백엔드에서 받은 키워드 차트 데이터의 색상 사용
            console.log('키워드 차트 데이터:', keywordStatistics.value.chart_data);
            let colors = [];

            if (keywordStatistics.value.chart_data && Array.isArray(keywordStatistics.value.chart_data)) {
                // 백엔드에서 chart_data로 색상을 받은 경우
                colors = keywordStatistics.value.chart_data.map(item => item.color);
                console.log('백엔드에서 받은 키워드 색상:', colors);
            } else {
                // 백엔드에서 색상을 받지 못한 경우 기본 색상 사용 (임시)
                console.log('백엔드에서 키워드 색상을 받지 못함');
                colors = sortedKeywords.map((_, index) => `hsl(${(index * 40) % 360}, 70%, 60%)`);
            }

            return {
                labels: sortedKeywords.map(keyword => keyword?.keyword_name || '키워드'),
                datasets: [{
                    label: '',
                    data: sortedKeywords.map(keyword => keyword?.total_amount || 0),
                    backgroundColor: colors,
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

        // 날짜 포맷팅
        const formatDate = (dateString) => {
            if (!dateString) return '-';
            const date = new Date(dateString);
            const year = date.getFullYear();
            const month = String(date.getMonth() + 1).padStart(2, '0');
            const day = String(date.getDate()).padStart(2, '0');
            return `${year}-${month}-${day}`;
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

        // 결제수단 정렬 순서 토글
        const togglePaymentMethodSortOrder = () => {
            paymentMethodSortOrder.value = paymentMethodSortOrder.value === 'desc' ? 'asc' : 'desc';
        };

        // 결제수단별 색상 생성 (인덱스 기반)
        const getPaymentMethodColor = (index) => {
            const colors = [
                '#3b82f6', // 파란색
                '#10b981', // 초록색
                '#f59e0b', // 주황색
                '#ef4444', // 빨간색
                '#8b5cf6', // 보라색
                '#ec4899', // 분홍색
                '#14b8a6', // 청록색
                '#f97316', // 진한 주황색
                '#6366f1', // 인디고
                '#84cc16', // 라임
            ];
            return colors[index % colors.length];
        };

        // 사용자별 색상 생성 (인덱스 기반)
        const getUserColor = (index) => {
            const colors = [
                '#f97316', // 주황색
                '#06b6d4', // 시안색
                '#8b5cf6', // 보라색
                '#ec4899', // 분홍색
                '#10b981', // 초록색
                '#eab308', // 노란색
                '#ef4444', // 빨간색
                '#3b82f6', // 파란색
                '#14b8a6', // 청록색
                '#a855f7', // 보라색2
            ];
            return colors[index % colors.length];
        };

        // 결제수단 상세 보기
        const showPaymentMethodDetail = async (paymentMethod) => {
            selectedPaymentMethod.value = paymentMethod;
            loadingPaymentMethodAccounts.value = true;

            try {
                const params = {
                    payment_method_id: paymentMethod.payment_method_id,
                    type: selectedPeriod.value,
                };

                // 개별 기간 선택 파라미터 추가
                if (selectedPeriod.value === 'year' && selectedYear.value) {
                    params.year = selectedYear.value;
                } else if (selectedPeriod.value === 'month' && selectedYear.value && selectedMonth.value) {
                    params.year = selectedYear.value;
                    params.month = selectedMonth.value;
                } else if (selectedPeriod.value === 'week' && selectedYear.value && selectedWeek.value) {
                    params.year = selectedYear.value;
                    params.week = selectedWeek.value;
                } else if (selectedPeriod.value === 'custom') {
                    params.start_date = customStartDate.value;
                    params.end_date = customEndDate.value;
                }

                const response = await api.get('/statistics/payment-method-accounts', params);
                paymentMethodAccounts.value = response.data.accounts || [];

                // 지출 내역 섹션으로 스크롤
                setTimeout(() => {
                    const element = document.getElementById('payment-method-detail');
                    if (element) {
                        element.scrollIntoView({ behavior: 'smooth', block: 'start' });
                    }
                }, 100);
            } catch (error) {
                console.error('결제수단별 지출 내역 조회 오류:', error);
                ElMessage.error('결제수단별 지출 내역 조회 중 오류가 발생했습니다');
                paymentMethodAccounts.value = [];
            } finally {
                loadingPaymentMethodAccounts.value = false;
            }
        };

        // 결제수단 상세 닫기
        const closePaymentMethodDetail = () => {
            selectedPaymentMethod.value = null;
            paymentMethodAccounts.value = [];
        };

        // 정렬된 사용자
        const sortedUsers = computed(() => {
            if (!statistics.value?.users) return [];

            const users = [...statistics.value.users];
            return users.sort((a, b) => {
                return userSortOrder.value === 'desc'
                    ? b.total_amount - a.total_amount
                    : a.total_amount - b.total_amount;
            });
        });

        // 사용자 차트 데이터
        const userChartData = computed(() => {
            if (!statistics.value?.users || statistics.value.users.length === 0) {
                return {
                    labels: [],
                    datasets: [{
                        label: '',
                        data: [],
                        backgroundColor: [],
                        borderWidth: 0
                    }]
                };
            }

            const colors = statistics.value.users.map((_, index) => getUserColor(index));

            return {
                labels: statistics.value.users.map(item => item.user_name),
                datasets: [{
                    label: '',
                    data: statistics.value.users.map(item => item.total_amount),
                    backgroundColor: colors,
                    borderWidth: 0,
                    hoverOffset: 10
                }]
            };
        });

        // 사용자 차트 옵션
        const userChartOptions = computed(() => {
            return {
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
                            label: (context) => {
                                const label = context.label || '';
                                const value = context.parsed || 0;
                                const total = context.dataset.data.reduce((a, b) => a + b, 0);
                                const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : 0;
                                return `${label}: ${formatMoney(value)}원 (${percentage}%)`;
                            }
                        }
                    }
                }
            };
        });

        // 사용자 정렬 순서 토글
        const toggleUserSortOrder = () => {
            userSortOrder.value = userSortOrder.value === 'desc' ? 'asc' : 'desc';
        };

        // 사용자 상세 보기
        const showUserDetail = async (user) => {
            selectedUserDetail.value = user;
            loadingUserAccounts.value = true;

            try {
                const params = {
                    user_name: user.user_name,
                    type: selectedPeriod.value,
                };

                // 개별 기간 선택 파라미터 추가
                if (selectedPeriod.value === 'year' && selectedYear.value) {
                    params.year = selectedYear.value;
                } else if (selectedPeriod.value === 'month' && selectedYear.value && selectedMonth.value) {
                    params.year = selectedYear.value;
                    params.month = selectedMonth.value;
                } else if (selectedPeriod.value === 'week' && selectedYear.value && selectedWeek.value) {
                    params.year = selectedYear.value;
                    params.week = selectedWeek.value;
                } else if (selectedPeriod.value === 'custom') {
                    params.start_date = customStartDate.value;
                    params.end_date = customEndDate.value;
                }

                const response = await api.get('/statistics/user-accounts', params);
                userAccounts.value = response.data.accounts || [];

                // 지출 내역 섹션으로 스크롤
                setTimeout(() => {
                    const element = document.getElementById('user-detail');
                    if (element) {
                        element.scrollIntoView({ behavior: 'smooth', block: 'start' });
                    }
                }, 100);
            } catch (error) {
                console.error('사용자별 지출 내역 조회 오류:', error);
                ElMessage.error('사용자별 지출 내역 조회 중 오류가 발생했습니다');
                userAccounts.value = [];
            } finally {
                loadingUserAccounts.value = false;
            }
        };

        // 사용자 상세 닫기
        const closeUserDetail = () => {
            selectedUserDetail.value = null;
            userAccounts.value = [];
        };

        // 키워드 정렬 순서 토글
        const toggleKeywordSortOrder = () => {
            keywordSortOrder.value = keywordSortOrder.value === 'desc' ? 'asc' : 'desc';
            // 키워드 정렬이 변경되면 선택된 키워드 해제
            selectedKeywordIndex.value = null;
        };

        // 고정/변동 카테고리 상세 보기 (키워드 표시)
        const showExpenseTypeCategoryDetail = async (category) => {
            selectedCategory.value = category;
            loadingKeywords.value = true;

            try {
                const params = {
                    category_id: category.category_id,
                    type: selectedPeriod.value,
                    category: selectedType.value
                };

                // 선택된 기간에 따라 파라미터 추가
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

                // 사용자가 선택된 경우에만 추가
                if (selectedUser.value) {
                    params.user_name = selectedUser.value;
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
                console.error('키워드 통계 로드 오류:', error);
                ElMessage.error('키워드 통계를 불러오는데 실패했습니다.');
            } finally {
                loadingKeywords.value = false;
            }
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

        // 고정/변동 지출 차트 클릭 핸들러
        const handleExpenseTypeChartClick = (event, elements) => {
            if (elements && elements.length > 0) {
                const index = elements[0].index;
                if (expenseTypeCategories.value && expenseTypeCategories.value[index]) {
                    const category = expenseTypeCategories.value[index];
                    showExpenseTypeCategoryDetail(category);
                }
            }
        };

        // 데이터 내보내기 페이지로 이동
        const openExportData = () => {
            router.push('/export-data');
        };

        onMounted(async () => {
            // 사용자 목록 로드
            try {
                await userStore.fetchUsers();
            } catch (error) {
                console.error('사용자 목록 로드 오류:', error);
            }

            // 카테고리 목록 로드 (고정/변동 지출 분석에 사용)
            try {
                await categoryStore.loadCategories();
            } catch (error) {
                console.error('카테고리 목록 로드 오류:', error);
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

            // 고정/변동 지출 관련
            chartViewType,
            expenseTypeTab,
            expenseTypeCategories,
            expenseTypeTotal,
            fixedExpenseTotal,
            variableExpenseTotal,
            expenseComparisonChartData,
            expenseComparisonChartOptions,
            expenseTypeChartData,
            expenseTypeChartOptions,

            // 결제수단 관련
            paymentMethodSortOrder,
            sortedPaymentMethods,
            paymentMethodChartData,
            paymentMethodChartOptions,
            selectedPaymentMethod,
            paymentMethodAccounts,
            loadingPaymentMethodAccounts,

            // 사용자 관련
            userSortOrder,
            sortedUsers,
            userChartData,
            userChartOptions,
            selectedUserDetail,
            userAccounts,
            loadingUserAccounts,

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
            formatDate,
            handleUserChange,
            handleTypeChange,
            handlePeriodChange,
            handleCustomDateChange,

            // 새로운 핸들러들
            handleYearChange,
            handleMonthChange,
            handleWeekChange,

            toggleSortOrder,
            togglePaymentMethodSortOrder,
            toggleUserSortOrder,
            toggleKeywordSortOrder,
            showCategoryDetail,
            showExpenseTypeCategoryDetail,
            showPaymentMethodDetail,
            showUserDetail,
            closeKeywordDetail,
            closePaymentMethodDetail,
            closeUserDetail,
            handleChartClick,
            handleExpenseTypeChartClick,
            openExportData,
            getCategoryColor,
            getPaymentMethodColor,
            getUserColor,

            // 아이콘들
            TrendingUp,
            TrendingDown,
            BarChart,
            Calculator,
            PieChart,
            Folder,
            Download,
            X
        };
    }
};
</script>

<style scoped>
/* 차트 유형 선택 버튼 */
.chart-type-button {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 1rem;
    min-height: 80px;
    background-color: white;
    border: 2px solid #e5e7eb;
    border-radius: 0.75rem;
    transition: all 0.2s ease;
    cursor: pointer;
    box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

.chart-type-button:hover {
    border-color: #3b82f6;
    background-color: #eff6ff;
    transform: translateY(-2px);
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.chart-type-button.active {
    border-color: #3b82f6;
    background-color: #3b82f6;
    color: white;
    font-weight: 600;
}

.chart-type-button.active:hover {
    background-color: #2563eb;
    border-color: #2563eb;
}

.chart-type-button:active {
    transform: translateY(0);
}

/* 모바일 최적화 */
@media (max-width: 768px) {
    .chart-type-button {
        min-height: 70px;
        padding: 0.75rem;
    }

    .chart-type-button span.text-2xl {
        font-size: 1.5rem;
    }

    .chart-type-button span.text-sm {
        font-size: 0.75rem;
    }
}

.chart-container {
    height: 300px;
    position: relative;
}

.comparison-chart-container {
    height: 200px;
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
