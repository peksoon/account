<template>
    <div class="page-container">
        <div class="page-content">
            <!-- Header -->
            <div class="flex items-center justify-between p-6 border-b border-gray-200">
                <div class="flex items-center">
                    <div
                        class="w-10 h-10 bg-gradient-to-r from-green-500 to-green-600 rounded-lg flex items-center justify-center mr-3">
                        <Download class="w-6 h-6 text-white" />
                    </div>
                    <div>
                        <h3 class="text-xl font-bold text-gray-900">데이터 내보내기</h3>
                        <p class="text-sm text-gray-500">가계부 데이터를 CSV 파일로 내보내기</p>
                    </div>
                </div>
                <el-button @click="goHome" type="info" size="large">
                    <Calendar class="w-4 h-4 mr-2" />
                    홈으로
                </el-button>
            </div>

            <!-- Export Options -->
            <div class="p-6 border-b border-gray-100">
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
                    <!-- Export Type -->
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">내보내기 유형</label>
                        <el-select v-model="exportType" class="w-full" @change="onExportTypeChange">
                            <el-option label="전체 데이터" value="all" />
                            <el-option label="기간별 데이터" value="period" />
                            <el-option label="카테고리별 데이터" value="category" />
                        </el-select>
                    </div>

                    <!-- Data Type -->
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">데이터 종류</label>
                        <el-select v-model="dataType" class="w-full">
                            <el-option label="전체 (수입+지출)" value="all" />
                            <el-option label="지출만" value="out" />
                            <el-option label="수입만" value="in" />
                        </el-select>
                    </div>

                    <!-- Period Type (기간별 선택 시만 표시) -->
                    <div v-if="exportType === 'period'">
                        <label class="block text-sm font-medium text-gray-700 mb-2">기간</label>
                        <el-select v-model="periodType" class="w-full">
                            <el-option label="이번 주" value="week" />
                            <el-option label="이번 달" value="month" />
                            <el-option label="올해" value="year" />
                            <el-option label="사용자 지정" value="custom" />
                        </el-select>
                    </div>

                    <!-- Category (카테고리별 선택 시만 표시) -->
                    <div v-if="exportType === 'category'">
                        <label class="block text-sm font-medium text-gray-700 mb-2">카테고리</label>
                        <el-select v-model="selectedCategory" class="w-full" placeholder="카테고리 선택">
                            <el-option v-for="category in categories" :key="category.id" :label="category.name"
                                :value="category.id" />
                        </el-select>
                    </div>
                </div>

                <!-- Custom Date Range (사용자 지정 기간 선택 시만 표시) -->
                <div v-if="exportType === 'period' && periodType === 'custom'"
                    class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">시작일</label>
                        <el-date-picker v-model="customStartDate" type="date" placeholder="시작일 선택"
                            format="YYYY년 MM월 DD일" value-format="YYYY-MM-DD" class="w-full" />
                    </div>
                    <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">종료일</label>
                        <el-date-picker v-model="customEndDate" type="date" placeholder="종료일 선택" format="YYYY년 MM월 DD일"
                            value-format="YYYY-MM-DD" class="w-full" />
                    </div>
                </div>

                <!-- Export Options -->
                <div class="border-t border-gray-100 pt-4 mb-6">
                    <h4 class="text-lg font-semibold text-gray-900 mb-4">내보내기 옵션</h4>
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div class="flex items-center">
                            <el-checkbox v-model="includeHeaders">헤더 포함</el-checkbox>
                        </div>
                        <div class="flex items-center">
                            <el-checkbox v-model="includeMemo">메모 포함</el-checkbox>
                        </div>
                        <div class="flex items-center">
                            <el-checkbox v-model="includeStats">통계 정보 포함</el-checkbox>
                        </div>
                    </div>
                </div>

                <!-- Export Button -->
                <div class="flex justify-center">
                    <el-button type="primary" @click="exportData" :loading="exporting" size="large" class="px-8">
                        <Download class="w-4 h-4 mr-2" />
                        {{ getExportButtonText() }}
                    </el-button>
                </div>
            </div>

            <!-- Preview Section -->
            <div v-if="previewData.length > 0" class="p-6">
                <div class="flex items-center justify-between mb-4">
                    <h4 class="text-lg font-semibold text-gray-900">미리보기 (최대 10개 항목)</h4>
                    <p class="text-sm text-gray-600">총 {{ totalCount }}건의 데이터</p>
                </div>

                <div class="overflow-x-auto">
                    <table class="min-w-full bg-white border border-gray-200 rounded-lg">
                        <thead class="bg-gray-50">
                            <tr>
                                <th class="px-4 py-2 border-b text-left text-sm font-medium text-gray-900">날짜</th>
                                <th class="px-4 py-2 border-b text-left text-sm font-medium text-gray-900">유형</th>
                                <th class="px-4 py-2 border-b text-left text-sm font-medium text-gray-900">카테고리</th>
                                <th class="px-4 py-2 border-b text-left text-sm font-medium text-gray-900">키워드</th>
                                <th class="px-4 py-2 border-b text-left text-sm font-medium text-gray-900">금액</th>
                                <th v-if="includeMemo"
                                    class="px-4 py-2 border-b text-left text-sm font-medium text-gray-900">메모</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(item, index) in previewData" :key="index" class="hover:bg-gray-50">
                                <td class="px-4 py-2 border-b text-sm">{{ formatDate(item.date) }}</td>
                                <td class="px-4 py-2 border-b text-sm">
                                    <span :class="item.type === 'out' ? 'text-red-600' : 'text-green-600'">
                                        {{ item.type === 'out' ? '지출' : '수입' }}
                                    </span>
                                </td>
                                <td class="px-4 py-2 border-b text-sm">{{ getCategoryName(item.category_id) }}</td>
                                <td class="px-4 py-2 border-b text-sm">{{ item.keyword_name || '-' }}</td>
                                <td class="px-4 py-2 border-b text-sm font-medium"
                                    :class="item.type === 'out' ? 'text-red-600' : 'text-green-600'">
                                    {{ item.type === 'out' ? '-' : '+' }}{{ formatMoney(item.money) }}원
                                </td>
                                <td v-if="includeMemo" class="px-4 py-2 border-b text-sm">{{ item.memo || '-' }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <!-- Summary Stats (통계 정보 포함 시) -->
            <div v-if="includeStats && previewData.length > 0" class="p-6 bg-gradient-to-r from-blue-50 to-purple-50">
                <h4 class="text-lg font-semibold text-gray-900 mb-4">통계 요약</h4>
                <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
                    <div class="bg-white p-4 rounded-lg shadow-sm">
                        <p class="text-sm text-gray-600">총 건수</p>
                        <p class="text-xl font-bold text-gray-900">{{ totalCount }}건</p>
                    </div>
                    <div class="bg-white p-4 rounded-lg shadow-sm">
                        <p class="text-sm text-gray-600">총 수입</p>
                        <p class="text-xl font-bold text-green-600">+{{ formatMoney(totalIncome) }}원</p>
                    </div>
                    <div class="bg-white p-4 rounded-lg shadow-sm">
                        <p class="text-sm text-gray-600">총 지출</p>
                        <p class="text-xl font-bold text-red-600">-{{ formatMoney(totalExpense) }}원</p>
                    </div>
                    <div class="bg-white p-4 rounded-lg shadow-sm">
                        <p class="text-sm text-gray-600">순 수익</p>
                        <p class="text-xl font-bold" :class="netAmount >= 0 ? 'text-green-600' : 'text-red-600'">
                            {{ netAmount >= 0 ? '+' : '' }}{{ formatMoney(netAmount) }}원
                        </p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { Download, Calendar } from 'lucide-vue-next';
import { useAccountStore } from '../stores/accountStore';
import { useCategoryStore } from '../stores/categoryStore';
import { getTodayKST } from '../utils';

export default {
    name: 'ExportDataPage',
    components: {
        Download,
        Calendar
    },
    setup() {
        const router = useRouter();
        const accountStore = useAccountStore();
        const categoryStore = useCategoryStore();

        // State
        const exportType = ref('all');
        const dataType = ref('all');
        const periodType = ref('month');
        const selectedCategory = ref('');
        const customStartDate = ref('');
        const customEndDate = ref('');
        const includeHeaders = ref(true);
        const includeMemo = ref(true);
        const includeStats = ref(true);
        const exporting = ref(false);
        const previewData = ref([]);
        const totalCount = ref(0);

        // Computed
        const categories = computed(() => categoryStore.categories || []);

        const totalIncome = computed(() => {
            return previewData.value
                .filter(item => item.type === 'in')
                .reduce((sum, item) => sum + item.money, 0);
        });

        const totalExpense = computed(() => {
            return previewData.value
                .filter(item => item.type === 'out')
                .reduce((sum, item) => sum + item.money, 0);
        });

        const netAmount = computed(() => totalIncome.value - totalExpense.value);

        // Methods
        const goHome = () => {
            router.push('/');
        };

        const onExportTypeChange = () => {
            // 내보내기 유형 변경 시 관련 설정 초기화
            if (exportType.value !== 'period') {
                periodType.value = 'month';
                customStartDate.value = '';
                customEndDate.value = '';
            }
            if (exportType.value !== 'category') {
                selectedCategory.value = '';
            }
            loadPreviewData();
        };

        const getDateRange = () => {
            const now = new Date();
            let startDate, endDate;

            if (exportType.value === 'period') {
                switch (periodType.value) {
                    case 'week': {
                        const startOfWeek = new Date(now);
                        startOfWeek.setDate(now.getDate() - now.getDay());
                        startDate = startOfWeek.toISOString().slice(0, 10);

                        const endOfWeek = new Date(startOfWeek);
                        endOfWeek.setDate(startOfWeek.getDate() + 6);
                        endDate = endOfWeek.toISOString().slice(0, 10);
                        break;
                    }
                    case 'month':
                        startDate = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().slice(0, 10);
                        endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0).toISOString().slice(0, 10);
                        break;
                    case 'year':
                        startDate = new Date(now.getFullYear(), 0, 1).toISOString().slice(0, 10);
                        endDate = new Date(now.getFullYear(), 11, 31).toISOString().slice(0, 10);
                        break;
                    case 'custom':
                        startDate = customStartDate.value;
                        endDate = customEndDate.value;
                        break;
                }
            } else {
                // 전체 데이터 또는 카테고리별 데이터의 경우
                startDate = '2020-01-01'; // 충분히 과거 날짜
                endDate = new Date(now.getFullYear() + 1, 11, 31).toISOString().slice(0, 10); // 내년까지
            }

            return { startDate, endDate };
        };

        const loadPreviewData = async () => {
            try {
                const { startDate, endDate } = getDateRange();

                if (exportType.value === 'period' && periodType.value === 'custom') {
                    if (!customStartDate.value || !customEndDate.value) {
                        previewData.value = [];
                        totalCount.value = 0;
                        return;
                    }
                }

                let allData = [];

                // 데이터 조회
                if (exportType.value === 'all' || exportType.value === 'period') {
                    allData = await accountStore.fetchAccountsInDateRange(startDate, endDate);
                } else if (exportType.value === 'category' && selectedCategory.value) {
                    // TODO: 카테고리별 조회 API 구현 필요
                    allData = await accountStore.fetchAccountsInDateRange(startDate, endDate);
                    allData = allData.filter(item => item.category_id === selectedCategory.value);
                }

                // 데이터 타입에 따라 필터링
                if (dataType.value !== 'all') {
                    allData = allData.filter(item => item.type === dataType.value);
                }

                // 날짜순 정렬 (최신순)
                allData.sort((a, b) => new Date(b.date) - new Date(a.date));

                totalCount.value = allData.length;
                previewData.value = allData.slice(0, 10); // 미리보기용 최대 10개

            } catch (error) {
                console.error('데이터 로드 오류:', error);
                ElMessage.error('데이터를 불러오는 중 오류가 발생했습니다.');
                previewData.value = [];
                totalCount.value = 0;
            }
        };

        const exportData = async () => {
            if (totalCount.value === 0) {
                ElMessage.warning('내보낼 데이터가 없습니다.');
                return;
            }

            exporting.value = true;
            try {
                const { startDate, endDate } = getDateRange();
                let allData = [];

                // 전체 데이터 조회 (미리보기가 아닌 전체)
                if (exportType.value === 'all' || exportType.value === 'period') {
                    allData = await accountStore.fetchAccountsInDateRange(startDate, endDate);
                } else if (exportType.value === 'category' && selectedCategory.value) {
                    allData = await accountStore.fetchAccountsInDateRange(startDate, endDate);
                    allData = allData.filter(item => item.category_id === selectedCategory.value);
                }

                // 데이터 타입에 따라 필터링
                if (dataType.value !== 'all') {
                    allData = allData.filter(item => item.type === dataType.value);
                }

                // 날짜순 정렬 (최신순)
                allData.sort((a, b) => new Date(b.date) - new Date(a.date));

                if (allData.length === 0) {
                    ElMessage.warning('내보낼 데이터가 없습니다.');
                    return;
                }

                // CSV 데이터 생성
                const csvData = allData.map(item => {
                    const row = {
                        날짜: formatDate(item.date),
                        유형: item.type === 'out' ? '지출' : '수입',
                        카테고리: getCategoryName(item.category_id),
                        키워드: item.keyword_name || '',
                        금액: item.money,
                        사용자: item.user || '',
                        결제수단: item.payment_method_name || ''
                    };

                    if (includeMemo.value) {
                        row.메모 = item.memo || '';
                    }

                    return row;
                });

                // 통계 정보 추가
                if (includeStats.value) {
                    const stats = [
                        {},
                        { 날짜: '=== 통계 요약 ===' },
                        { 날짜: '총 건수', 유형: `${allData.length}건` },
                        { 날짜: '총 수입', 유형: `${formatMoney(totalIncome.value)}원` },
                        { 날짜: '총 지출', 유형: `${formatMoney(totalExpense.value)}원` },
                        { 날짜: '순 수익', 유형: `${formatMoney(netAmount.value)}원` },
                        {}
                    ];
                    csvData.push(...stats);
                }

                // CSV 파일 생성 및 다운로드
                const headers = Object.keys(csvData[0]);
                let csvContent = '';

                if (includeHeaders.value) {
                    csvContent += headers.join(',') + '\n';
                }

                csvContent += csvData.map(row =>
                    headers.map(header => {
                        const value = row[header] || '';
                        // CSV에서 쉼표와 따옴표 처리
                        if (typeof value === 'string' && (value.includes(',') || value.includes('"') || value.includes('\n'))) {
                            return `"${value.replace(/"/g, '""')}"`;
                        }
                        return value;
                    }).join(',')
                ).join('\n');

                // UTF-8 BOM 추가 (한글 지원)
                const BOM = '\uFEFF';
                const blob = new Blob([BOM + csvContent], { type: 'text/csv;charset=utf-8;' });

                const link = document.createElement('a');
                link.href = URL.createObjectURL(blob);

                // 파일명 생성
                const today = getTodayKST();
                let filename = '';
                switch (exportType.value) {
                    case 'all':
                        filename = `가계부_전체데이터_${today}.csv`;
                        break;
                    case 'period':
                        filename = `가계부_기간별데이터_${today}.csv`;
                        break;
                    case 'category': {
                        const categoryName = getCategoryName(selectedCategory.value);
                        filename = `가계부_${categoryName}_${today}.csv`;
                        break;
                    }
                }

                link.download = filename;
                link.click();

                ElMessage.success(`${allData.length}건의 데이터를 내보냈습니다.`);
            } catch (error) {
                console.error('데이터 내보내기 오류:', error);
                ElMessage.error('데이터 내보내기 중 오류가 발생했습니다.');
            } finally {
                exporting.value = false;
            }
        };

        const getExportButtonText = () => {
            if (exporting.value) return '내보내는 중...';
            return `${totalCount.value}건 데이터 내보내기`;
        };

        const formatMoney = (amount) => {
            return new Intl.NumberFormat('ko-KR').format(amount);
        };

        const formatDate = (dateString) => {
            const date = new Date(dateString);
            return date.toLocaleDateString('ko-KR', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit'
            });
        };

        const getCategoryName = (categoryId) => {
            if (!categoryId) return '';
            const category = categoryStore.getCategoryById(categoryId);
            return category ? category.name : '';
        };

        // Watchers
        watch([dataType, periodType, selectedCategory, customStartDate, customEndDate], () => {
            loadPreviewData();
        });

        // Initialize
        onMounted(async () => {
            // 카테고리 로드
            await categoryStore.fetchCategories();

            // 기본 날짜 설정
            const now = new Date();
            customStartDate.value = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().slice(0, 10);
            customEndDate.value = new Date(now.getFullYear(), now.getMonth() + 1, 0).toISOString().slice(0, 10);

            // 초기 데이터 로드
            await loadPreviewData();
        });

        return {
            // State
            exportType,
            dataType,
            periodType,
            selectedCategory,
            customStartDate,
            customEndDate,
            includeHeaders,
            includeMemo,
            includeStats,
            exporting,
            previewData,
            totalCount,

            // Computed
            categories,
            totalIncome,
            totalExpense,
            netAmount,

            // Methods
            goHome,
            onExportTypeChange,
            loadPreviewData,
            exportData,
            getExportButtonText,
            formatMoney,
            formatDate,
            getCategoryName,

            // Icons
            Download,
            Calendar
        };
    }
};
</script>

<style scoped>
.page-container {
    @apply min-h-screen bg-gray-50;
}

.page-content {
    @apply bg-white w-full h-full max-w-none max-h-none overflow-hidden flex flex-col min-h-screen;
}

/* 테이블 스타일 */
table {
    font-size: 0.875rem;
}

/* 모바일 최적화 */
@media (max-width: 768px) {
    .page-content {
        @apply text-sm;
    }

    table {
        font-size: 0.75rem;
    }

    .grid {
        @apply grid-cols-1;
    }
}
</style>
