<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50">
    <!-- 헤더 영역 -->
    <div class="container-responsive" :class="isMobile ? 'py-4' : 'py-8'">
      <div class="card mb-8" :class="isMobile ? 'p-4' : 'p-6'">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between space-y-4 md:space-y-0">
          <!-- 타이틀 및 현재 월 -->
          <div>
            <div class="flex items-center gap-3 mb-2">
              <div :class="isMobile ? 'w-8 h-8' : 'w-10 h-10'"
                class="bg-gradient-to-r from-yellow-500 to-yellow-600 rounded-lg flex items-center justify-center">
                <DollarSign :class="isMobile ? 'w-4 h-4' : 'w-6 h-6'" class="text-white" />
              </div>
              <h1 :class="isMobile ? 'text-xl' : 'text-3xl'" class="font-bold text-gradient">상윤 가족 가계부</h1>
            </div>
            <p :class="isMobile ? 'text-base' : 'text-lg'" class="text-gray-600 font-medium">{{ currentYear }}년 {{ currentMonth }}월</p>
          </div>

          <!-- 컨트롤 버튼들 -->
          <div class="flex flex-col space-y-3">
            <!-- 달력 모드일 때만 표시 -->
            <template v-if="viewMode === 'calendar'">
              <!-- 첫 번째 줄: 네비게이션 버튼들 (모바일에서 가로 배치) -->
              <div class="flex flex-wrap gap-2 justify-center md:justify-end">
                <el-button @click="goToToday" type="primary" :icon="Calendar" 
                  :size="isMobile ? 'small' : 'large'"
                  :class="isMobile ? 'text-xs px-3 py-1 min-w-0 flex-shrink-0' : 'w-auto'">
                  {{ isMobile ? '오늘' : '오늘' }}
                </el-button>
                <el-button-group :class="isMobile ? 'flex-shrink-0' : ''">
                  <el-button @click="goToPrevMonth" :size="isMobile ? 'small' : 'large'" 
                    :class="isMobile ? 'text-xs px-2 py-1 min-w-0' : ''">
                    {{ isMobile ? '←' : '← 이전' }}
                  </el-button>
                  <el-button @click="goToNextMonth" :size="isMobile ? 'small' : 'large'"
                    :class="isMobile ? 'text-xs px-2 py-1 min-w-0' : ''">
                    {{ isMobile ? '→' : '다음 →' }}
                  </el-button>
                </el-button-group>
                <el-button @click="openAddPopup" type="success" :size="isMobile ? 'small' : 'large'"
                  :class="isMobile ? 'text-xs px-3 py-1 min-w-0 flex-shrink-0' : 'w-auto'">
                  {{ isMobile ? '+' : '+ 추가' }}
                </el-button>
                <el-button @click="openStatistics" type="info" :icon="BarChart" :size="isMobile ? 'small' : 'large'"
                  :class="isMobile ? 'text-xs px-2 py-1 min-w-0 flex-shrink-0' : 'w-auto'">
                  {{ isMobile ? '📊' : '📊 통계' }}
                </el-button>
                <el-button @click="openBudgetManager" type="warning" :size="isMobile ? 'small' : 'large'"
                  :class="isMobile ? 'text-xs px-2 py-1 min-w-0 flex-shrink-0' : 'w-auto'">
                  {{ isMobile ? '💰' : '💰 기준치 관리' }}
                </el-button>
              </div>
            </template>
            <template v-else>
              <div class="flex justify-center md:justify-end">
                <el-button @click="goBackToCalendar" type="primary" :icon="Calendar"
                  :size="isMobile ? 'default' : 'large'" class="w-full sm:w-auto">
                  📅 달력으로 돌아가기
                </el-button>
              </div>
            </template>
          </div>
        </div>

        <!-- 월별 통계 카드 -->
        <div v-if="monthlyStats" :class="isMobile ? 'mt-4' : 'mt-6'">
          <!-- 모바일: 유연한 가로 레이아웃 -->
          <div v-if="isMobile" class="flex gap-2 overflow-x-auto pb-1">
            <div class="bg-gradient-to-r from-green-400 to-green-500 text-white px-3 py-2 rounded-md flex-1 min-w-0">
              <div class="text-center">
                <p class="text-green-100 text-xs leading-tight truncate">총 수입</p>
                <p class="text-xs font-bold leading-tight mt-0.5 truncate">{{ formatMoney(monthlyStats.totalIncome) }}원</p>
              </div>
            </div>

            <div class="bg-gradient-to-r from-red-400 to-red-500 text-white px-3 py-2 rounded-md flex-1 min-w-0">
              <div class="text-center">
                <p class="text-red-100 text-xs leading-tight truncate">총 지출</p>
                <p class="text-xs font-bold leading-tight mt-0.5 truncate">{{ formatMoney(monthlyStats.totalExpense) }}원</p>
              </div>
            </div>

            <div class="bg-gradient-to-r from-blue-400 to-blue-500 text-white px-3 py-2 rounded-md flex-1 min-w-0">
              <div class="text-center">
                <p class="text-blue-100 text-xs leading-tight truncate">잔액</p>
                <p class="text-xs font-bold leading-tight mt-0.5 truncate" :class="monthlyStats.balance >= 0 ? 'text-white' : 'text-yellow-200'">
                  {{ formatMoney(monthlyStats.balance) }}원
                </p>
              </div>
            </div>
          </div>

          <!-- 데스크톱: 기존 그리드 레이아웃 -->
          <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="bg-gradient-to-r from-green-400 to-green-500 text-white p-4 rounded-xl">
              <div class="flex items-center">
                <TrendingUp class="w-8 h-8 mr-3" />
                <div>
                  <p class="text-green-100 text-sm">총 수입</p>
                  <p class="text-xl font-bold">{{ formatMoney(monthlyStats.totalIncome) }}원</p>
                </div>
              </div>
            </div>

            <div class="bg-gradient-to-r from-red-400 to-red-500 text-white p-4 rounded-xl">
              <div class="flex items-center">
                <TrendingDown class="w-8 h-8 mr-3" />
                <div>
                  <p class="text-red-100 text-sm">총 지출</p>
                  <p class="text-xl font-bold">{{ formatMoney(monthlyStats.totalExpense) }}원</p>
                </div>
              </div>
            </div>

            <div class="bg-gradient-to-r from-blue-400 to-blue-500 text-white p-4 rounded-xl">
              <div class="flex items-center">
                <Wallet class="w-8 h-8 mr-3" />
                <div>
                  <p class="text-blue-100 text-sm">잔액</p>
                  <p class="text-xl font-bold" :class="monthlyStats.balance >= 0 ? 'text-white' : 'text-yellow-200'">
                    {{ formatMoney(monthlyStats.balance) }}원
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 달력 모드 -->
      <template v-if="viewMode === 'calendar'">
        <!-- 달력 영역 -->
        <div class="card mb-8" :class="isMobile ? 'p-2' : 'p-6'">
          <FullCalendar ref="calendar" :options="calendarOptions" class="modern-calendar" />
        </div>

        <!-- 선택된 날짜 상세 정보 -->
        <div v-if="selectedDateData.length" class="card animate-slide-up" :class="isMobile ? 'p-4' : 'p-6'">
          <div class="flex items-center mb-4">
            <Calendar class="w-6 h-6 text-primary-500 mr-2" />
            <h2 class="text-xl font-bold text-gray-800">{{ selectedDate }} 가계부</h2>
          </div>

          <div class="grid gap-3">
            <div v-for="(data, index) in selectedDateData" :key="index" @click="showDetailPopup(data)"
              class="p-4 border border-gray-200 rounded-lg hover:border-primary-300 hover:shadow-md transition-all duration-200 cursor-pointer">
              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <div class="flex-shrink-0 mr-3">
                    <div v-if="data.type === 'out'" class="w-3 h-3 bg-red-500 rounded-full"></div>
                    <div v-else class="w-3 h-3 bg-green-500 rounded-full"></div>
                  </div>
                  <div>
                    <p class="font-semibold text-gray-800">{{ getCategoryName(data.category_id) || data.category || '-'
                      }}
                    </p>
                    <p v-if="data.keyword_name || data.keyword" class="text-sm text-gray-600">🏷️ {{ data.keyword_name
                      ||
                      data.keyword }}</p>
                    <p v-else class="text-sm text-gray-500">키워드 없음</p>
                  </div>
                </div>
                <div class="text-right">
                  <p class="font-bold" :class="data.type === 'out' ? 'text-red-600' : 'text-green-600'">
                    {{ data.type === 'out' ? '-' : '+' }}{{ formatMoney(data.money) }}원
                  </p>
                  <p class="text-xs text-gray-500">{{ data.type === 'out' ? '지출' : '수입' }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 로딩 상태 -->
        <div v-if="loading" class="flex justify-center items-center py-12">
          <div class="spinner"></div>
          <span class="ml-3 text-gray-600">데이터를 불러오는 중...</span>
        </div>
      </template>

      <!-- 통계 모드 -->
      <template v-else-if="viewMode === 'statistics'">
        <div class="card p-6">
          <StatisticsDashboard @close="closeStatistics" @open-budget-manager="openBudgetManager" />
        </div>
      </template>

      <!-- 기준치 관리 모드 -->
      <template v-else-if="viewMode === 'budget'">
        <div class="card p-6">
          <BudgetManager @close="closeBudgetManager" />
        </div>
      </template>
    </div>

    <!-- 팝업들 -->
    <AddPopup v-if="showAddPopup" :newAccount="newAccount" @close="closeAddPopup" @save="saveAccount"
      @open-category-manager="openCategoryManager" @open-keyword-manager="openKeywordManager"
      @open-payment-method-manager="openPaymentMethodManager" @open-deposit-path-manager="openDepositPathManager"
      @open-user-manager="openUserManager" @budget-alert="handleBudgetAlert"
      @budget-save-success="handleBudgetSaveSuccess" />

    <DetailPopup v-if="showCustomPopup" :eventDetail="eventDetail" :isEditMode="isEditMode" @close="closePopup"
      @edit="openEditMode" @update="updateAccount" @delete="deleteAccount" @cancel-edit="cancelEdit" />

    <!-- 관리 모달들 -->
    <UserManager v-if="showUserManager" @close="closeUserManager" />
    <CategoryManager v-if="showCategoryManager" @close="closeCategoryManager" />
    <PaymentMethodManager v-if="showPaymentMethodManager" @close="closePaymentMethodManager" />
    <DepositPathManager v-if="showDepositPathManager" @close="closeDepositPathManager" />

    <!-- 키워드 관리 모달 -->
    <KeywordManager v-if="showKeywordManager" :category-id="keywordManagerCategoryId" @close="closeKeywordManager" />

    <!-- 기준치 알림 팝업 -->
    <BudgetAlertPopup :is-visible="showBudgetAlert" :budget-usage="budgetAlertData.budgetUsage"
      :expense-amount="budgetAlertData.expenseAmount" :expense-date="budgetAlertData.expenseDate"
      :expense-keyword="budgetAlertData.expenseKeyword" @close="closeBudgetAlert"
      @open-budget-management="openBudgetManager" />



  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue';
import { storeToRefs } from 'pinia';
import FullCalendar from '@fullcalendar/vue3';
import dayGridPlugin from '@fullcalendar/daygrid';
import interactionPlugin from '@fullcalendar/interaction';
import { useAccountStore } from '../stores/accountStore';
import { useCalendarStore } from '../stores/calendarStore';
import { usePopupStore } from '../stores/popupStore';
import { useCategoryStore } from '../stores/categoryStore';
// import { formatDateToString } from '../utils';
import AddPopup from './AddPopup.vue';
import DetailPopup from './DetailPopup.vue';
import UserManager from './UserManager.vue';
import CategoryManager from './CategoryManager.vue';
import PaymentMethodManager from './PaymentMethodManager.vue';
import DepositPathManager from './DepositPathManager.vue';
import KeywordManager from './KeywordManager.vue';
import StatisticsDashboard from './StatisticsDashboard.vue';
import BudgetAlertPopup from './BudgetAlertPopup.vue';
import BudgetManager from './BudgetManager.vue';
import {
  Calendar,
  BarChart,
  TrendingUp,
  TrendingDown,
  Wallet,
  DollarSign
} from 'lucide-vue-next';


export default {
  components: {
    FullCalendar,
    AddPopup,
    DetailPopup,
    UserManager,
    CategoryManager,
    PaymentMethodManager,
    DepositPathManager,
    KeywordManager,
    StatisticsDashboard,
    BudgetAlertPopup,
    BudgetManager,
    Calendar,
    TrendingUp,
    TrendingDown,
    Wallet,
    DollarSign
  },
  setup() {
    const accountStore = useAccountStore();
    const calendarStore = useCalendarStore();
    const popupStore = usePopupStore();
    const categoryStore = useCategoryStore();
    const { showCustomPopup, eventDetail, isEditMode, showAddPopup } = storeToRefs(popupStore);

    const calendar = ref(null);
    const selectedDate = ref('');
    const selectedDateData = ref([]);
    const loading = ref(false);

    // 관리 모달 상태
    const showUserManager = ref(false);
    const showCategoryManager = ref(false);
    const showPaymentMethodManager = ref(false);
    const showDepositPathManager = ref(false);
    const showStatistics = ref(false);
    const showKeywordManager = ref(false);
    const keywordManagerCategoryId = ref(null);

    // 기준치 관련 상태
    const showBudgetAlert = ref(false);

    const budgetAlertData = ref({
      budgetUsage: null,
      expenseAmount: 0,
      expenseDate: '',
      expenseKeyword: ''
    });

    // 뷰 모드 상태 ('calendar', 'statistics', 또는 'budget')
    const viewMode = ref('calendar');

    // 모바일 체크
    const isMobile = computed(() => {
      if (typeof window === 'undefined') return false;
      return window.innerWidth < 768;
    });

    // 현재 년도와 월
    const currentYear = computed(() => calendarStore.currentYear);
    const currentMonth = computed(() => calendarStore.currentMonth);

    // 새 계정 데이터
    const newAccount = ref(popupStore.newAccount);

    // 월별 통계 계산
    const monthlyStats = computed(() => {
      if (!accountStore.monthlyData.length) return null;

      const totalIncome = accountStore.monthlyData
        .filter(item => item.type === 'in')
        .reduce((sum, item) => sum + item.money, 0);

      const totalExpense = accountStore.monthlyData
        .filter(item => item.type === 'out')
        .reduce((sum, item) => sum + item.money, 0);

      return {
        totalIncome,
        totalExpense,
        balance: totalIncome - totalExpense
      };
    });

    // FullCalendar 옵션
    const calendarOptions = ref({
      plugins: [dayGridPlugin, interactionPlugin],
      headerToolbar: false, // 헤더 툴바 숨김 (직접 구현)
      initialView: 'dayGridMonth',
      locale: 'ko',
      height: 'auto',
      events: [],
      dateClick: handleDateClick,
      eventClick: handleEventClick,
      dayMaxEvents: window.innerWidth < 768 ? 2 : 3, // 모바일에서는 이벤트 수 제한
      moreLinkClick: 'popover',
      dayHeaderContent: (arg) => {
        const dayNames = ['일', '월', '화', '수', '목', '금', '토'];
        return dayNames[arg.date.getDay()];
      },
      // 모바일 최적화
      aspectRatio: window.innerWidth < 768 ? 0.8 : 1.35, // 모바일에서 더 세로로 압축
      eventTextColor: '#fff',
      eventDisplay: 'block',
      dayMaxEventRows: window.innerWidth < 768 ? 2 : 3,
      dayHeaderFormat: window.innerWidth < 768 ? { weekday: 'narrow' } : { weekday: 'short' },
      fixedWeekCount: false, // 주 수를 고정하지 않음 (공간 절약)
      dayCellContent: (arg) => {
        // 모바일에서는 날짜 숫자만 작게 표시
        const fontSize = window.innerWidth < 768 ? 'text-xs' : 'text-sm';
        return { html: `<div class="${fontSize} font-medium">${arg.dayNumberText.replace('일', '')}</div>` };
      }
    });

    // 달력 초기화
    const initCalendar = async () => {
      loading.value = true;
      try {
        const now = new Date();
        calendarStore.currentYear = now.getFullYear();
        calendarStore.currentMonth = String(now.getMonth() + 1).padStart(2, '0');

        await fetchAndUpdateCalendar();
      } catch (error) {
        console.error('Calendar initialization error:', error);
      } finally {
        loading.value = false;
      }
    };

    // 달력 데이터 가져오기 및 업데이트
    const fetchAndUpdateCalendar = async () => {
      await accountStore.fetchMonthAccounts(
        calendarStore.currentYear,
        parseInt(calendarStore.currentMonth)
      );
      updateCalendarEvents();
    };

    // 달력 이벤트 업데이트
    const updateCalendarEvents = () => {
      // 날짜별로 데이터를 그룹화
      const dailyTotals = {};

      accountStore.monthlyData.forEach(item => {
        const dateKey = item.date.split(' ')[0];
        if (!dailyTotals[dateKey]) {
          dailyTotals[dateKey] = { income: 0, expense: 0, items: [] };
        }

        if (item.type === 'out') {
          dailyTotals[dateKey].expense += item.money;
        } else {
          dailyTotals[dateKey].income += item.money;
        }
        dailyTotals[dateKey].items.push(item);
      });

      // 각 날짜별로 총합 이벤트 생성
      const events = [];
      Object.entries(dailyTotals).forEach(([date, totals]) => {
        // 수입 총합 이벤트
        if (totals.income > 0) {
          events.push({
            id: `income-${date}`,
            title: `+${formatMoney(totals.income)}`,
            date: date,
            className: 'income-total',
            display: 'list-item',
            extendedProps: { type: 'income-total', amount: totals.income, items: totals.items.filter(i => i.type === 'in') }
          });
        }

        // 지출 총합 이벤트
        if (totals.expense > 0) {
          events.push({
            id: `expense-${date}`,
            title: `-${formatMoney(totals.expense)}`,
            date: date,
            className: 'expense-total',
            display: 'list-item',
            extendedProps: { type: 'expense-total', amount: totals.expense, items: totals.items.filter(i => i.type === 'out') }
          });
        }
      });

      calendarOptions.value.events = events;

      // 달력 API를 통해 이벤트 다시 렌더링
      if (calendar.value) {
        const calendarApi = calendar.value.getApi();
        calendarApi.removeAllEvents();
        calendarApi.addEventSource(events);
      }
    };

    // 날짜 클릭 핸들러
    function handleDateClick(info) {
      const dateStr = info.dateStr;
      selectedDate.value = dateStr;
      selectedDateData.value = accountStore.fetchDataForDate(dateStr);

      // 날짜 클릭 시 AddPopup 열기 (해당 날짜로 설정)
      popupStore.openAddPopup(dateStr);
    }

    // 이벤트 클릭 핸들러
    function handleEventClick(info) {
      const eventData = info.event.extendedProps;

      // 총합 이벤트인 경우 해당 날짜의 개별 항목들을 표시
      if (eventData.type === 'income-total' || eventData.type === 'expense-total') {
        const dateStr = info.event.startStr;
        selectedDate.value = dateStr;
        selectedDateData.value = accountStore.fetchDataForDate(dateStr);
        return;
      }

      // 개별 이벤트인 경우 상세 팝업 표시
      popupStore.showDetailPopup(eventData);
    }

    // 이전 달로 이동
    const goToPrevMonth = async () => {
      if (calendar.value) {
        const calendarApi = calendar.value.getApi();
        calendarApi.prev();
        updateCurrentDate();
        await fetchAndUpdateCalendar();
        clearSelection();
      }
    };

    // 다음 달로 이동
    const goToNextMonth = async () => {
      if (calendar.value) {
        const calendarApi = calendar.value.getApi();
        calendarApi.next();
        updateCurrentDate();
        await fetchAndUpdateCalendar();
        clearSelection();
      }
    };

    // 오늘로 이동
    const goToToday = async () => {
      if (calendar.value) {
        const calendarApi = calendar.value.getApi();
        calendarApi.today();
        updateCurrentDate();
        await fetchAndUpdateCalendar();
        clearSelection();
      }
    };

    // 현재 날짜 업데이트
    const updateCurrentDate = () => {
      if (calendar.value) {
        const calendarApi = calendar.value.getApi();
        const currentDate = calendarApi.getDate();
        calendarStore.currentYear = currentDate.getFullYear();
        calendarStore.currentMonth = String(currentDate.getMonth() + 1).padStart(2, '0');
      }
    };

    // 선택 초기화
    const clearSelection = () => {
      selectedDate.value = '';
      selectedDateData.value = [];
    };

    // 계정 저장 (수입만 처리, 지출은 budget-save-success에서 처리)
    const saveAccount = async (accountData) => {
      // 지출은 budgetStore.createOutAccountWithBudget에서 이미 처리되므로 수입만 처리
      if (accountData.type === 'out') {
        console.warn('지출은 budgetStore에서 처리되므로 saveAccount에서 무시됩니다.');
        return;
      }

      loading.value = true;
      try {
        await accountStore.saveAccount(accountData);
        popupStore.closeAddPopup();
        await fetchAndUpdateCalendar();

        // 현재 선택된 날짜의 데이터 다시 불러오기
        if (selectedDate.value) {
          selectedDateData.value = accountStore.fetchDataForDate(selectedDate.value);
        }
      } catch (error) {
        console.error('Account save error:', error);
      } finally {
        loading.value = false;
      }
    };

    // 계정 업데이트
    const updateAccount = async (updatedData) => {
      loading.value = true;
      try {
        await accountStore.updateAccount(updatedData);
        popupStore.closePopup();
        await fetchAndUpdateCalendar();

        // 현재 선택된 날짜의 데이터 다시 불러오기
        if (selectedDate.value) {
          selectedDateData.value = accountStore.fetchDataForDate(selectedDate.value);
        }
      } catch (error) {
        console.error('Account update error:', error);
      } finally {
        loading.value = false;
      }
    };

    // 계정 삭제
    const deleteAccount = async (accountData) => {
      loading.value = true;
      try {
        await accountStore.deleteAccount(accountData);
        popupStore.closePopup();
        await fetchAndUpdateCalendar();

        // 현재 선택된 날짜의 데이터 다시 불러오기
        if (selectedDate.value) {
          selectedDateData.value = accountStore.fetchDataForDate(selectedDate.value);
        }
      } catch (error) {
        console.error('Account delete error:', error);
      } finally {
        loading.value = false;
      }
    };

    // 팝업 관련 메서드
    const openAddPopup = () => popupStore.openAddPopup();
    const closeAddPopup = () => popupStore.closeAddPopup();
    const closePopup = () => popupStore.closePopup();
    const openEditMode = () => popupStore.openEditMode();
    const cancelEdit = () => popupStore.closePopup();
    const showDetailPopup = (data) => popupStore.showDetailPopup(data);

    // 관리 모달 관련 메서드
    const openCategoryManager = () => {
      showCategoryManager.value = true;
    };

    const openKeywordManager = (categoryId) => {
      keywordManagerCategoryId.value = categoryId;
      showKeywordManager.value = true;
    };

    const closeKeywordManager = () => {
      showKeywordManager.value = false;
      keywordManagerCategoryId.value = null;
    };

    const openPaymentMethodManager = () => {
      showPaymentMethodManager.value = true;
    };

    const openUserManager = () => {
      showUserManager.value = true;
    };

    const closeCategoryManager = () => {
      showCategoryManager.value = false;
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

    const closeUserManager = () => {
      showUserManager.value = false;
    };

    const openStatistics = () => {
      viewMode.value = 'statistics';
    };

    const closeStatistics = () => {
      viewMode.value = 'calendar';
    };

    const goBackToCalendar = () => {
      viewMode.value = 'calendar';
    };

    // 기준치 알림 처리
    const handleBudgetAlert = (alertData) => {
      budgetAlertData.value = alertData;
      showBudgetAlert.value = true;
      // AddPopup 닫기
      popupStore.closeAddPopup();
    };

    // 기준치 저장 성공 처리 (캘린더 갱신)
    const handleBudgetSaveSuccess = async () => {
      loading.value = true;
      try {
        // 달력 갱신
        await fetchAndUpdateCalendar();

        // 현재 선택된 날짜의 데이터 다시 불러오기
        if (selectedDate.value) {
          selectedDateData.value = accountStore.fetchDataForDate(selectedDate.value);
        }
      } catch (error) {
        console.error('Calendar update error:', error);
      } finally {
        loading.value = false;
      }
    };

    const closeBudgetAlert = () => {
      showBudgetAlert.value = false;
      budgetAlertData.value = {
        budgetUsage: null,
        expenseAmount: 0,
        expenseDate: '',
        expenseKeyword: ''
      };
      // AddPopup도 닫기
      popupStore.closeAddPopup();
    };

    const openBudgetManager = () => {
      viewMode.value = 'budget';
    };

    const closeBudgetManager = () => {
      viewMode.value = 'calendar';
    };

    // 금액 포맷팅 함수
    const formatMoney = (amount) => {
      return new Intl.NumberFormat('ko-KR').format(amount);
    };

    // 카테고리 이름 가져오기
    const getCategoryName = (categoryId) => {
      if (!categoryId) return '';
      const category = categoryStore.getCategoryById(categoryId);
      return category ? category.name : '';
    };

    // 컴포넌트 마운트 시 초기화
    onMounted(async () => {
      await initCalendar();
      // 카테고리 목록 로드
      try {
        await categoryStore.fetchCategories();
      } catch (error) {
        console.error('카테고리 목록 로드 오류:', error);
      }
    });

    return {
      // 상태
      calendar,
      loading,
      selectedDate,
      selectedDateData,
      isMobile,
      currentYear,
      currentMonth,
      monthlyStats,
      newAccount,

      // 달력 옵션
      calendarOptions,

      // 팝업 상태
      showCustomPopup,
      eventDetail,
      isEditMode,
      showAddPopup,

      // 메서드
      goToPrevMonth,
      goToNextMonth,
      goToToday,
      saveAccount,
      updateAccount,
      deleteAccount,
      openAddPopup,
      closeAddPopup,
      closePopup,
      openEditMode,
      cancelEdit,
      showDetailPopup,
      formatMoney,
      getCategoryName,

      // 관리 모달 상태
      showUserManager,
      showCategoryManager,
      showPaymentMethodManager,
      showDepositPathManager,

      // 관리 모달 메서드들
      openCategoryManager,
      openKeywordManager,
      openPaymentMethodManager,
      openDepositPathManager,
      openUserManager,
      closeCategoryManager,
      closePaymentMethodManager,
      closeDepositPathManager,
      closeUserManager,
      closeKeywordManager,
      showKeywordManager,
      keywordManagerCategoryId,
      openStatistics,
      closeStatistics,
      goBackToCalendar,
      showStatistics,
      viewMode,

      // 기준치 관련 상태 및 메서드
      showBudgetAlert,
      budgetAlertData,
      handleBudgetAlert,
      handleBudgetSaveSuccess,
      closeBudgetAlert,
      openBudgetManager,
      closeBudgetManager,

      // 아이콘
      Calendar,
      BarChart,
      TrendingUp,
      TrendingDown,
      Wallet
    };
  },
};
</script>

<style scoped>
/* 커스텀 달력 스타일 */
.modern-calendar :deep(.fc) {
  font-family: 'Inter', sans-serif;
  background: transparent;
}

.modern-calendar :deep(.fc-scrollgrid) {
  border: none;
}

.modern-calendar :deep(.fc-theme-standard td) {
  border: 1px solid #e5e7eb;
}

.modern-calendar :deep(.fc-theme-standard th) {
  border: 1px solid #e5e7eb;
  background: #f9fafb;
  font-weight: 600;
  color: #374151;
  padding: 12px 8px;
}

/* 모바일 최적화 */
@media (max-width: 768px) {
  .modern-calendar :deep(.fc-theme-standard th) {
    padding: 6px 2px;
    font-size: 0.7rem;
    font-weight: 600;
  }

  .modern-calendar :deep(.fc-daygrid-day-number) {
    font-size: 0.75rem;
    padding: 2px;
    line-height: 1;
  }

  .modern-calendar :deep(.fc-daygrid-event) {
    font-size: 0.65rem;
    padding: 1px 3px;
    margin: 0.5px 1px;
    line-height: 1.2;
  }

  .modern-calendar :deep(.fc-event-title) {
    font-size: 0.65rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    line-height: 1.2;
  }

  .modern-calendar :deep(.fc-daygrid-day-frame) {
    min-height: 2.5rem;
  }

  .modern-calendar :deep(.fc-daygrid-day-top) {
    padding: 2px;
  }

  .modern-calendar :deep(.fc-scrollgrid-sync-table) {
    font-size: 0.7rem;
  }

  .modern-calendar :deep(.fc-daygrid-day-events) {
    margin-top: 1px;
  }

  .modern-calendar :deep(.fc-theme-standard td) {
    border-width: 0.5px;
  }
}

.modern-calendar :deep(.fc-daygrid-day) {
  background: white;
  transition: background-color 0.2s ease;
}

.modern-calendar :deep(.fc-daygrid-day:hover) {
  background: #f8fafc;
}

.modern-calendar :deep(.fc-day-today) {
  background: #eff6ff !important;
  border-color: #3b82f6 !important;
}

.modern-calendar :deep(.fc-daygrid-day-number) {
  padding: 8px;
  font-weight: 500;
  color: #1f2937;
}

.modern-calendar :deep(.fc-day-today .fc-daygrid-day-number) {
  color: #1d4ed8;
  font-weight: 700;
}

.modern-calendar :deep(.fc-event) {
  border: none;
  border-radius: 6px;
  padding: 2px 6px;
  margin: 1px 2px;
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.modern-calendar :deep(.fc-event:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.modern-calendar :deep(.expense-event) {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
}

.modern-calendar :deep(.income-event) {
  background: linear-gradient(135deg, #22c55e, #16a34a);
  color: white;
}

.modern-calendar :deep(.fc-more-link) {
  color: #6b7280;
  font-size: 10px;
  padding: 2px 4px;
  border-radius: 4px;
  background: #f3f4f6;
  border: none;
}

.modern-calendar :deep(.fc-more-link:hover) {
  background: #e5e7eb;
  color: #374151;
}

/* 수입/지출 총합 이벤트 스타일 */
.modern-calendar :deep(.income-total) {
  background: #10b981 !important;
  border-color: #10b981 !important;
  color: white !important;
  font-weight: 600 !important;
  border-radius: 3px !important;
  margin: 0.5px 0 !important;
  padding: 1px 3px !important;
  font-size: 0.65rem !important;
  line-height: 1.2 !important;
}

.modern-calendar :deep(.expense-total) {
  background: #ef4444 !important;
  border-color: #ef4444 !important;
  color: white !important;
  font-weight: 600 !important;
  border-radius: 3px !important;
  margin: 0.5px 0 !important;
  padding: 1px 3px !important;
  font-size: 0.65rem !important;
  line-height: 1.2 !important;
}

/* 모바일에서 이벤트 스타일 더 컴팩트하게 */
@media (max-width: 768px) {
  .modern-calendar :deep(.income-total) {
    font-size: 0.6rem !important;
    padding: 0.5px 2px !important;
    margin: 0.2px 0 !important;
  }

  .modern-calendar :deep(.expense-total) {
    font-size: 0.6rem !important;
    padding: 0.5px 2px !important;
    margin: 0.2px 0 !important;
  }
}



/* 애니메이션 클래스 */
.animate-slide-up {
  animation: slideUp 0.4s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 통계 카드 호버 효과 */
.bg-gradient-to-r:hover {
  transform: translateY(-2px);
  transition: transform 0.2s ease;
}

/* 모바일 컨테이너 최적화 */
@media (max-width: 768px) {
  .container-responsive {
    padding-left: 0.75rem;
    padding-right: 0.75rem;
  }
  
  .card {
    border-radius: 0.5rem;
  }
}

/* 매우 작은 화면 최적화 (320px 이하) */
@media (max-width: 320px) {
  .container-responsive {
    padding-left: 0.5rem;
    padding-right: 0.5rem;
  }
  
  /* 통계 카드 간격 더 줄이기 */
  .flex.gap-2 {
    gap: 0.25rem;
  }
  
  /* 카드 패딩 더 줄이기 */
  .px-3.py-2 {
    padding: 0.375rem 0.5rem;
  }
  
  /* 폰트 크기 더 작게 */
  .text-xs {
    font-size: 0.65rem;
  }
}
</style>
