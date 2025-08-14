<template>
  <div class="page-container">
    <div class="page-content">
      <!-- Header -->
      <div class="flex items-center justify-between p-6 border-b border-gray-200">
        <div class="flex items-center">
          <div
            class="w-10 h-10 bg-gradient-to-r from-purple-500 to-purple-600 rounded-lg flex items-center justify-center mr-3">
            <Search class="w-6 h-6 text-white" />
          </div>
          <div>
            <h3 class="text-xl font-bold text-gray-900">키워드 검색</h3>
            <p class="text-sm text-gray-500">키워드 검색 또는 기간별 전체 내역 조회</p>
          </div>
        </div>
        <el-button @click="goHome" type="info" size="large">
          <Calendar class="w-4 h-4 mr-2" />
          홈으로
        </el-button>
      </div>

      <!-- Search Form -->
      <div class="p-6 border-b border-gray-100">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4 mb-4">
          <!-- Keyword Input -->
          <div class="md:col-span-2 lg:col-span-1">
            <label class="block text-sm font-medium text-gray-700 mb-2">키워드</label>
            <el-autocomplete v-model="searchQuery" :fetch-suggestions="queryKeywordSuggestions" placeholder="키워드 입력..."
              :prefix-icon="Search" @keyup.enter="handleSearch" @select="handleKeywordSelect" clearable
              class="w-full" />
          </div>

          <!-- Search Type -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">검색 유형</label>
            <el-select v-model="searchType" class="w-full">
              <el-option label="전체" value="all" />
              <el-option label="지출" value="out" />
              <el-option label="수입" value="in" />
            </el-select>
          </div>

          <!-- Date Range -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">기간</label>
            <el-select v-model="dateRange" class="w-full">
              <el-option label="이번 주" value="week" />
              <el-option label="이번 달" value="month" />
              <el-option label="올해" value="year" />
              <el-option label="사용자 지정" value="custom" />
            </el-select>
          </div>

          <!-- Sort By -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">정렬</label>
            <el-select v-model="sortBy" class="w-full">
              <el-option label="날짜 (최신순)" value="date-desc" />
              <el-option label="날짜 (오래된순)" value="date-asc" />
              <el-option label="금액 (높은순)" value="amount-desc" />
              <el-option label="금액 (낮은순)" value="amount-asc" />
            </el-select>
          </div>

          <!-- Search Button -->
          <div class="flex items-end">
            <el-button type="primary" @click="handleSearch" :loading="searching" class="w-full">
              <Search class="w-4 h-4 mr-2" />
              검색
            </el-button>
          </div>
        </div>

        <!-- Custom Date Range -->
        <div v-if="dateRange === 'custom'" class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">시작일</label>
            <el-date-picker v-model="customStartDate" type="date" placeholder="시작일 선택" format="YYYY년 MM월 DD일"
              value-format="YYYY-MM-DD" class="w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">종료일</label>
            <el-date-picker v-model="customEndDate" type="date" placeholder="종료일 선택" format="YYYY년 MM월 DD일"
              value-format="YYYY-MM-DD" class="w-full" />
          </div>
        </div>
      </div>

      <!-- Results Summary -->
      <div v-if="searchResults.length > 0"
        class="p-6 bg-gradient-to-r from-blue-50 to-purple-50 border-b border-gray-100">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between">
          <div>
            <h4 class="text-lg font-semibold text-gray-900">검색 결과: {{ searchResults.length }}건</h4>
            <p class="text-sm text-gray-600">{{ formatDateRange() }}</p>
          </div>
          <div class="mt-2 md:mt-0">
            <div class="text-right">
              <template v-if="searchType === 'out'">
                <p class="text-sm text-gray-600">총 지출 금액</p>
                <p class="text-2xl font-bold text-red-600">{{ formatMoney(totalExpenseAmount) }}원</p>
              </template>
              <template v-else-if="searchType === 'in'">
                <p class="text-sm text-gray-600">총 수입 금액</p>
                <p class="text-2xl font-bold text-green-600">{{ formatMoney(totalIncomeAmount) }}원</p>
              </template>
              <template v-else>
                <div class="space-y-1">
                  <div>
                    <p class="text-xs text-gray-500">총 수입</p>
                    <p class="text-lg font-bold text-green-600">+{{ formatMoney(totalIncomeAmount) }}원</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500">총 지출</p>
                    <p class="text-lg font-bold text-red-600">-{{ formatMoney(totalExpenseAmount) }}원</p>
                  </div>
                </div>
              </template>
            </div>
          </div>
        </div>
      </div>

      <!-- Results List -->
      <div class="flex-1 overflow-y-auto p-6">
        <div v-if="searching" class="flex justify-center items-center py-12">
          <div class="text-center">
            <div class="spinner mb-4"></div>
            <p class="text-gray-600">검색 중...</p>
          </div>
        </div>

        <div v-else-if="searchResults.length === 0 && hasSearched" class="text-center py-12">
          <div class="text-gray-500 mb-4">
            <Search class="w-16 h-16 mx-auto mb-4 opacity-50" />
            <p class="text-lg">검색 결과가 없습니다</p>
            <p class="text-sm">다른 조건으로 검색해보세요</p>
          </div>
        </div>

        <div v-else-if="searchResults.length === 0 && !searching" class="text-center py-12">
          <div class="text-gray-500">
            <Search class="w-16 h-16 mx-auto mb-4 opacity-50" />
            <p class="text-lg">선택한 기간에 데이터가 없습니다</p>
            <p class="text-sm">다른 기간을 선택해보세요</p>
          </div>
        </div>

        <div v-else class="space-y-3">
          <div v-for="(item, index) in searchResults" :key="index" @click="openDetailModal(item)"
            class="result-item p-4 border border-gray-200 rounded-lg hover:border-purple-300 hover:shadow-md transition-all duration-200 cursor-pointer">
            <div class="flex items-center justify-between">
              <div class="flex-1">
                <div class="flex items-center mb-2">
                  <div class="flex-shrink-0 mr-3">
                    <div v-if="item.type === 'out'" class="w-3 h-3 bg-red-500 rounded-full"></div>
                    <div v-else class="w-3 h-3 bg-green-500 rounded-full"></div>
                  </div>
                  <div>
                    <h5 class="text-lg font-bold text-gray-900">{{ item.keyword_name || '키워드 없음' }}</h5>
                    <p class="text-sm text-gray-600">{{ getCategoryName(item.category_id) || '카테고리 없음' }}</p>
                  </div>
                </div>
                <div class="flex items-center text-xs text-gray-500 space-x-4">
                  <span class="flex items-center">
                    <Calendar class="w-3 h-3 mr-1" />
                    {{ formatDate(item.date) }}
                  </span>
                  <span v-if="item.memo" class="flex items-center">
                    <FileText class="w-3 h-3 mr-1" />
                    {{ item.memo.length > 20 ? item.memo.substring(0, 20) + '...' : item.memo }}
                  </span>
                </div>
              </div>
              <div class="text-right ml-4">
                <p class="text-lg font-bold" :class="item.type === 'out' ? 'text-red-600' : 'text-green-600'">
                  {{ item.type === 'out' ? '-' : '+' }}{{ formatMoney(item.money) }}원
                </p>
                <p class="text-xs text-gray-500">{{ item.type === 'out' ? '지출' : '수입' }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="flex justify-end p-6 border-t border-gray-200 space-x-3">
        <el-button v-if="searchResults.length > 0" @click="exportResults" type="success" size="large">
          <Download class="w-4 h-4 mr-2" />
          결과 내보내기
        </el-button>
      </div>
    </div>

    <!-- Detail Modal -->
    <DetailPopup v-if="showDetailModal" :eventDetail="selectedItem" :isEditMode="false" @close="closeDetailModal" />
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { Search, Calendar, FileText, Download } from 'lucide-vue-next';
import { useRouter } from 'vue-router';
import { useAccountStore } from '../stores/accountStore';
import { useCategoryStore } from '../stores/categoryStore';
import DetailPopup from './DetailPopup.vue';
import { getTodayKST } from '../utils';

export default {
  name: 'KeywordSearchPage',
  components: {
    Search,
    Calendar,
    FileText,
    Download,
    DetailPopup
  },
  setup() {
    const router = useRouter();
    const accountStore = useAccountStore();
    const categoryStore = useCategoryStore();

    // State
    const searchQuery = ref('');
    const searchType = ref('all');
    const dateRange = ref('month');
    const sortBy = ref('date-desc');
    const customStartDate = ref('');
    const customEndDate = ref('');
    const searching = ref(false);
    const searchResults = ref([]);
    const showDetailModal = ref(false);
    const selectedItem = ref(null);
    const hasSearched = ref(false);
    const keywordSuggestions = ref([]);


    // Total expense amount
    const totalExpenseAmount = computed(() => {
      return searchResults.value
        .filter(item => item.type === 'out')
        .reduce((sum, item) => sum + item.money, 0);
    });

    // Total income amount
    const totalIncomeAmount = computed(() => {
      return searchResults.value
        .filter(item => item.type === 'in')
        .reduce((sum, item) => sum + item.money, 0);
    });

    // Get date range for search
    const getDateRange = () => {
      const now = new Date();
      let startDate, endDate;

      switch (dateRange.value) {
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

        default:
          startDate = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().slice(0, 10);
          endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0).toISOString().slice(0, 10);
      }

      return { startDate, endDate };
    };

    // Format date range for display
    const formatDateRange = () => {
      const { startDate, endDate } = getDateRange();
      const start = new Date(startDate).toLocaleDateString('ko-KR');
      const end = new Date(endDate).toLocaleDateString('ko-KR');
      return `${start} ~ ${end}`;
    };

    // Search function
    const handleSearch = async (showMessages = true) => {
      searching.value = true;
      try {
        const { startDate, endDate } = getDateRange();

        let allData = [];

        if (!searchQuery.value.trim()) {
          // 키워드가 없으면 해당 기간의 모든 데이터 조회
          allData = await accountStore.fetchAccountsInDateRange(startDate, endDate);
        } else {
          // 키워드가 있으면 키워드 검색
          allData = await accountStore.searchByKeyword(
            searchQuery.value.trim(),
            startDate,
            endDate
          );
        }

        // 검색 유형에 따라 필터링
        if (searchType.value !== 'all') {
          allData = allData.filter(item => item.type === searchType.value);
        }

        // Sort results
        let sortedResults = [...allData];
        switch (sortBy.value) {
          case 'date-desc':
            sortedResults.sort((a, b) => new Date(b.date) - new Date(a.date));
            break;
          case 'date-asc':
            sortedResults.sort((a, b) => new Date(a.date) - new Date(b.date));
            break;
          case 'amount-desc':
            sortedResults.sort((a, b) => b.money - a.money);
            break;
          case 'amount-asc':
            sortedResults.sort((a, b) => a.money - b.money);
            break;
        }

        searchResults.value = sortedResults;
        hasSearched.value = true;

        if (showMessages) {
          if (sortedResults.length === 0) {
            if (searchQuery.value.trim()) {
              ElMessage.info(`"${searchQuery.value}" 키워드에 대한 검색 결과가 없습니다.`);
            } else {
              ElMessage.info('선택한 기간에 데이터가 없습니다.');
            }
          } else {
            if (searchQuery.value.trim()) {
              ElMessage.success(`"${searchQuery.value}" 키워드로 ${sortedResults.length}건을 찾았습니다.`);
            } else {
              ElMessage.success(`선택한 기간에 ${sortedResults.length}건의 데이터를 찾았습니다.`);
            }
          }
        }
      } catch (error) {
        console.error('키워드 검색 오류:', error);
        ElMessage.error('검색 중 오류가 발생했습니다.');
        searchResults.value = [];
      } finally {
        searching.value = false;
      }
    };

    // Format money
    const formatMoney = (amount) => {
      return new Intl.NumberFormat('ko-KR').format(amount);
    };

    // Format date
    const formatDate = (dateString) => {
      const date = new Date(dateString);
      return date.toLocaleDateString('ko-KR', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        weekday: 'short'
      });
    };

    // Get category name
    const getCategoryName = (categoryId) => {
      if (!categoryId) return '';
      const category = categoryStore.getCategoryById(categoryId);
      return category ? category.name : '';
    };

    // Open detail modal
    const openDetailModal = (item) => {
      selectedItem.value = item;
      showDetailModal.value = true;
    };

    // Close detail modal
    const closeDetailModal = () => {
      showDetailModal.value = false;
      selectedItem.value = null;
    };

    // Export results
    const exportResults = () => {
      if (searchResults.value.length === 0) return;

      const csvData = searchResults.value.map(item => ({
        날짜: formatDate(item.date),
        유형: item.type === 'out' ? '지출' : '수입',
        카테고리: getCategoryName(item.category_id),
        키워드: item.keyword_name || '',
        금액: item.money,
        메모: item.memo || ''
      }));

      // Simple CSV export
      const csvContent = [
        Object.keys(csvData[0]).join(','),
        ...csvData.map(row => Object.values(row).join(','))
      ].join('\n');

      const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = `키워드검색결과_${searchQuery.value}_${getTodayKST()}.csv`;
      link.click();

      ElMessage.success('검색 결과를 내보냈습니다.');
    };

    // Keyword autocomplete
    const queryKeywordSuggestions = async (queryString, callback) => {
      if (!queryString || queryString.length < 1) {
        callback([]);
        return;
      }

      try {
        // 기존 검색 결과에서 키워드 추출
        const localKeywords = [];
        searchResults.value.forEach(item => {
          if (item.keyword_name && !localKeywords.includes(item.keyword_name)) {
            localKeywords.push(item.keyword_name);
          }
        });

        // 로컬 키워드에서 매칭되는 항목 찾기
        const matches = localKeywords
          .filter(keyword => keyword.toLowerCase().includes(queryString.toLowerCase()))
          .slice(0, 10) // 최대 10개로 제한
          .map(keyword => ({ value: keyword }));

        callback(matches);
      } catch (error) {
        console.error('키워드 자동완성 오류:', error);
        callback([]);
      }
    };

    // Handle keyword selection
    const handleKeywordSelect = (item) => {
      searchQuery.value = item.value;
      handleSearch();
    };

    // Go to home
    const goHome = () => {
      router.push('/');
    };

    // Initialize
    onMounted(async () => {
      // Set default custom dates to this month
      const now = new Date();
      customStartDate.value = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().slice(0, 10);
      customEndDate.value = new Date(now.getFullYear(), now.getMonth() + 1, 0).toISOString().slice(0, 10);

      // 최초 로드 시 현재 기간의 모든 데이터 조회 (메시지 표시 안함)
      await handleSearch(false);
    });

    return {
      // State
      searchQuery,
      searchType,
      dateRange,
      sortBy,
      customStartDate,
      customEndDate,
      searching,
      searchResults,
      showDetailModal,
      selectedItem,
      hasSearched,
      keywordSuggestions,

      // Computed
      totalExpenseAmount,
      totalIncomeAmount,

      // Methods
      handleSearch,
      formatDateRange,
      formatMoney,
      formatDate,
      getCategoryName,
      openDetailModal,
      closeDetailModal,
      exportResults,
      queryKeywordSuggestions,
      handleKeywordSelect,
      goHome,

      // Icons
      Search,
      Calendar,
      FileText,
      Download
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

.result-item {
  transition: all 0.2s ease;
}

.result-item:hover {
  transform: translateY(-2px);
}

.spinner {
  @apply inline-block w-8 h-8 border-4 border-purple-200 border-t-purple-600 rounded-full;
  animation: spin 1s linear infinite;
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
    transform: translateY(-20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Mobile optimizations */
@media (max-width: 768px) {
  .result-item {
    @apply p-3;
  }

  .result-item h5 {
    @apply text-sm;
  }

  .result-item p {
    @apply text-xs;
  }
}
</style>