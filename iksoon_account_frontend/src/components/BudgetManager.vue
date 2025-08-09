<template>
  <div class="budget-manager">
    <!-- 헤더 -->
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold text-gray-800">카테고리별 기준치 관리</h2>
      <!-- 서버 연결 상태 표시 -->
      <div v-if="budgetStore.error && budgetStore.error.includes('서버에 연결할 수 없습니다')"
        class="px-3 py-1 bg-red-100 text-red-700 text-sm rounded-md">
        ⚠️ 서버 연결 오류
      </div>
    </div>
    <!-- 액션 버튼 -->
    <div class="flex justify-end mb-6">
      <button @click="showCreateModal = true"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
        <Plus class="w-4 h-4 inline-block mr-2" />
        기준치 추가
      </button>
    </div>

    <!-- 사용자 선택 -->
    <div class="mb-6">
      <label class="block text-sm font-medium text-gray-700 mb-2">사용자 필터 (선택사항)</label>
      <select v-model="selectedUser" @change="loadBudgets"
        class="w-64 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
        <option value="">전체 (사용자 구분 없음)</option>
        <option v-for="user in users" :key="user.id" :value="user.name">
          {{ user.name }}
        </option>
      </select>
    </div>

    <!-- 로딩 상태 -->
    <div v-if="budgetStore.loading" class="text-center py-8">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <p class="mt-2 text-gray-600">로딩 중...</p>
    </div>

    <!-- 에러 메시지 -->
    <div v-if="budgetStore.error" class="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
      {{ budgetStore.error }}
    </div>

    <!-- 기준치 목록 -->
    <div v-if="!budgetStore.loading" class="space-y-4">
      <div v-if="budgetStore.budgets.length === 0" class="text-center py-8 text-gray-500">
        설정된 기준치가 없습니다.
      </div>

      <div v-for="budget in budgetStore.budgets" :key="budget.id"
        class="bg-white border border-gray-200 rounded-lg p-6 shadow-sm hover:shadow-md transition-shadow cursor-pointer"
        @click="selectCategory(budget)">
        <div class="flex justify-between items-start mb-4">
          <div>
            <h3 class="text-lg font-semibold text-gray-800 hover:text-blue-600 transition-colors">
              {{ budget.category_name }}
            </h3>
            <p class="text-sm text-gray-600">카테고리 ID: {{ budget.category_id }}</p>
          </div>
          <div class="flex space-x-2">
            <button @click.stop="editBudget(budget)"
              class="px-3 py-1 text-blue-600 bg-blue-100 rounded hover:bg-blue-200 transition-colors">
              <Edit class="w-4 h-4" />
            </button>
            <button @click.stop="deleteBudget(budget)"
              class="px-3 py-1 text-red-600 bg-red-100 rounded hover:bg-red-200 transition-colors">
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- 월 기준치 -->
          <div class="bg-blue-50 p-4 rounded-lg">
            <h4 class="text-sm font-medium text-blue-800 mb-2">월 기준치</h4>
            <p class="text-2xl font-bold text-blue-600">
              {{ formatMoney(budget.monthly_budget) }}원
            </p>
          </div>

          <!-- 연 기준치 -->
          <div class="bg-green-50 p-4 rounded-lg">
            <h4 class="text-sm font-medium text-green-800 mb-2">연 기준치</h4>
            <p class="text-2xl font-bold text-green-600">
              {{ formatMoney(budget.yearly_budget) }}원
            </p>
          </div>
        </div>

        <!-- 사용량 정보 (있는 경우) -->
        <div v-if="getBudgetUsage(budget.category_id)" class="mt-4 pt-4 border-t border-gray-200">
          <h4 class="text-sm font-medium text-gray-700 mb-3">현재 사용량</h4>
          <BudgetUsageDisplay :usage="getBudgetUsage(budget.category_id)" />
        </div>
      </div>
    </div>

    <!-- 기준치 비교 그래프 섹션 -->
    <div v-if="selectedCategory && budgetStore.budgets.length > 0" class="mt-8">
      <div class="bg-white border border-gray-200 rounded-lg p-6 shadow-sm">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-xl font-bold text-gray-800">📊 {{ selectedCategory.category_name }} 기준치 비교 분석</h3>
          <button @click="selectedCategory = null"
            class="px-3 py-1 text-gray-600 bg-gray-100 rounded hover:bg-gray-200 transition-colors">
            ✕ 닫기
          </button>
        </div>

        <!-- 월 선택 -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">비교할 월 선택</label>
          <select v-model="selectedCompareMonth" @change="updateBudgetChart"
            class="w-48 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
            <option v-for="month in availableMonths" :key="month.value" :value="month.value">
              {{ month.label }}
            </option>
          </select>
        </div>

        <!-- 차트 영역 -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- 월별 기준치 vs 사용량 -->
          <div class="bg-gray-50 p-4 rounded-lg">
            <h4 class="text-lg font-semibold text-gray-800 mb-4">월별 기준치 비교</h4>
            <div class="relative h-80">
              <canvas ref="monthlyChartCanvas"></canvas>
            </div>
          </div>

          <!-- 연별 기준치 vs 사용량 -->
          <div class="bg-gray-50 p-4 rounded-lg">
            <h4 class="text-lg font-semibold text-gray-800 mb-4">연별 기준치 비교</h4>
            <div class="relative h-80">
              <canvas ref="yearlyChartCanvas"></canvas>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 기준치 생성/수정 모달 -->
    <div v-if="showCreateModal || showEditModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-96 max-w-sm mx-4">
        <h3 class="text-lg font-semibold mb-4">
          {{ showEditModal ? '기준치 수정' : '기준치 추가' }}
        </h3>

        <form @submit.prevent="submitBudget">
          <!-- 카테고리 선택 (생성 시에만) -->
          <div v-if="showCreateModal" class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">카테고리</label>
            <el-autocomplete v-model="categorySearchText" :fetch-suggestions="fetchCategorySuggestions"
              placeholder="카테고리를 입력하세요" @select="handleCategorySelect" @clear="handleCategoryClear" clearable
              value-key="name" class="w-full">
              <template #default="{ item }">
                <div class="flex justify-between items-center">
                  <span>{{ item.name }}</span>
                  <span class="text-gray-500 text-sm">{{ item.type === 'out' ? '지출' : '수입' }}</span>
                </div>
              </template>
            </el-autocomplete>
          </div>

          <!-- 사용자 선택 (생성 시에만, 선택사항) -->
          <div v-if="showCreateModal" class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">사용자 (선택사항)</label>
            <select v-model="budgetForm.user_name" @change="handleUserChange"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
              <option value="">전체 (사용자 구분 없음)</option>
              <option v-for="user in users" :key="user.id" :value="user.name">
                {{ user.name }}
              </option>
            </select>
          </div>

          <!-- 월 기준치 -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">월 기준치 (원)</label>
            <input v-model.number="budgetForm.monthly_budget" type="number" min="0"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="0" />
          </div>

          <!-- 연 기준치 -->
          <div class="mb-6">
            <label class="block text-sm font-medium text-gray-700 mb-2">연 기준치 (원)</label>
            <input v-model.number="budgetForm.yearly_budget" type="number" min="0"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="0" />
          </div>

          <!-- 버튼 -->
          <div class="flex justify-end space-x-3">
            <button type="button" @click="closeModals"
              class="px-4 py-2 text-gray-600 bg-gray-200 rounded-md hover:bg-gray-300 transition-colors">
              취소
            </button>
            <button type="submit" :disabled="!isFormValid"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
              {{ showEditModal ? '수정' : '생성' }}
            </button>
          </div>
        </form>
      </div>
    </div>

  </div> <!-- budget-manager 닫기 -->
</template>

<script>
import { ref, computed, onMounted, onUnmounted, watch, defineComponent, nextTick } from 'vue'
import { Plus, Edit, Trash2 } from 'lucide-vue-next'
import { ElMessage } from 'element-plus'
import { Chart, registerables } from 'chart.js'
import { useBudgetStore } from '@/stores/budgetStore'
import { useUserStore } from '@/stores/userStore'
import { useCategoryStore } from '@/stores/categoryStore'
import BudgetUsageDisplay from './BudgetUsageDisplay.vue'

Chart.register(...registerables)

export default defineComponent({
  name: 'BudgetManager',
  components: {
    BudgetUsageDisplay,
    Plus,
    Edit,
    Trash2
  },
  emits: ['close'],
  setup(props, { emit }) {

    // 스토어
    const budgetStore = useBudgetStore()
    const userStore = useUserStore()
    const categoryStore = useCategoryStore()

    // 반응형 데이터
    const selectedUser = ref('')
    const showCreateModal = ref(false)
    const showEditModal = ref(false)
    const editingBudget = ref(null)
    const categorySearchText = ref('')
    const selectedCategory = ref(null) // 선택된 카테고리 추가

    // 차트 관련 데이터
    const selectedCompareMonth = ref(new Date().toISOString().slice(0, 7)) // YYYY-MM 형식
    const monthlyChartCanvas = ref(null)
    const yearlyChartCanvas = ref(null)

    // Chart.js 인스턴스는 반응형으로 만들지 않음 (무한 루프 방지)
    let monthlyChart = null
    let yearlyChart = null

    // 폼 데이터
    const budgetForm = ref({
      category_id: '',
      user_name: '',
      monthly_budget: 0,
      yearly_budget: 0
    })

    // 사용자 목록
    const users = computed(() => userStore.users || [])

    // 사용 가능한 카테고리 목록 (지출 카테고리만, 현재 선택된 사용자에 대해 이미 기준치가 설정되지 않은 것만)
    const availableCategories = computed(() => {
      const outCategories = categoryStore.categories.filter(cat => cat.type === 'out')

      // 현재 폼에서 선택된 사용자에 대해 이미 기준치가 설정된 카테고리 ID들
      const currentUserName = budgetForm.value.user_name || ""
      const usedCategoryIds = budgetStore.budgets
        .filter(budget => (budget.user_name || "") === currentUserName)
        .map(budget => budget.category_id)

      return outCategories.filter(cat => !usedCategoryIds.includes(cat.id))
    })

    // 사용 가능한 월 목록 (현재 월 기준으로 과거 12개월)
    const availableMonths = computed(() => {
      const months = []
      const currentDate = new Date()

      for (let i = 0; i < 12; i++) {
        const date = new Date(currentDate.getFullYear(), currentDate.getMonth() - i, 1)
        const value = date.toISOString().slice(0, 7) // YYYY-MM
        const label = `${date.getFullYear()}년 ${date.getMonth() + 1}월`
        months.push({ value, label })
      }

      return months
    })

    // 폼 유효성 검사
    const isFormValid = computed(() => {
      if (showCreateModal.value) {
        return budgetForm.value.category_id &&
          (budgetForm.value.monthly_budget > 0 || budgetForm.value.yearly_budget > 0)
      } else {
        return budgetForm.value.monthly_budget > 0 || budgetForm.value.yearly_budget > 0
      }
    })

    // 메소드
    const formatMoney = (amount) => {
      return new Intl.NumberFormat('ko-KR').format(amount || 0)
    }

    const loadBudgets = async () => {
      budgetStore.clearError()
      // 사용자가 선택되지 않은 경우 빈 문자열로 전달하여 모든 기준치 조회
      await budgetStore.fetchBudgets(selectedUser.value || "")
      if (selectedUser.value) {
        await budgetStore.fetchBudgetUsage(selectedUser.value)
      }
    }

    // 카테고리 자동완성 제안 가져오기
    const fetchCategorySuggestions = (queryString, callback) => {
      const suggestions = availableCategories.value.filter(category => {
        return category.name.toLowerCase().includes(queryString.toLowerCase())
      })
      callback(suggestions)
    }

    // 카테고리 선택 핸들러
    const handleCategorySelect = (item) => {
      budgetForm.value.category_id = item.id
      categorySearchText.value = item.name
    }

    // 사용자 선택 변경 시 카테고리 초기화
    const handleUserChange = () => {
      budgetForm.value.category_id = ''
      categorySearchText.value = ''
    }

    // 카테고리 클리어 핸들러
    const handleCategoryClear = () => {
      budgetForm.value.category_id = ''
      categorySearchText.value = ''
    }

    // 카테고리 선택 함수
    const selectCategory = async (budget) => {
      selectedCategory.value = budget

      console.log('카테고리 선택됨:', budget.category_name, '기준치 사용자:', budget.user_name)

      // 선택된 기준치의 사용자명을 기준으로 사용량 조회
      // 전체 기준치(user_name='')인 경우 현재 선택된 사용자 또는 빈 문자열 사용
      const targetUser = budget.user_name || selectedUser.value || ''

      console.log('사용량 조회 대상 사용자:', targetUser)

      try {
        await budgetStore.fetchBudgetUsage(targetUser)
        console.log('사용량 데이터 로딩 후:', budgetStore.budgetUsages?.length || 0)
      } catch (error) {
        console.error('사용량 데이터 로딩 오류:', error)
        // 빈 문자열로 다시 시도
        if (targetUser !== '') {
          console.log('빈 문자열로 재시도')
          await budgetStore.fetchBudgetUsage('')
        }
      }

      await nextTick()

      // 차트 영역이 DOM에 렌더링될 때까지 충분히 대기
      setTimeout(() => {
        initCharts()
        setTimeout(() => {
          updateBudgetChart()
        }, 300)
      }, 300)
    }

    // 차트 초기화
    const initCharts = async () => {
      console.log('차트 초기화 시작')

      // 기존 차트가 있으면 먼저 제거
      destroyCharts()

      await nextTick()

      // DOM이 완전히 준비될 때까지 대기
      await new Promise(resolve => setTimeout(resolve, 300))

      try {
        console.log('월별 차트 캔버스 확인:', !!monthlyChartCanvas.value)
        console.log('월별 차트 부모 노드 확인:', !!monthlyChartCanvas.value?.parentNode)

        if (monthlyChartCanvas.value && monthlyChartCanvas.value.parentNode) {
          const ctx = monthlyChartCanvas.value.getContext('2d')

          // 캔버스 크기 강제 설정
          monthlyChartCanvas.value.style.width = '100%'
          monthlyChartCanvas.value.style.height = '400px'

          monthlyChart = new Chart(ctx, {
            type: 'bar',
            data: {
              labels: [],
              datasets: [
                {
                  label: '월 기준치',
                  data: [],
                  backgroundColor: 'rgba(59, 130, 246, 0.5)',
                  borderColor: 'rgb(59, 130, 246)',
                  borderWidth: 1
                },
                {
                  label: '월 사용량',
                  data: [],
                  backgroundColor: 'rgba(239, 68, 68, 0.5)',
                  borderColor: 'rgb(239, 68, 68)',
                  borderWidth: 1
                }
              ]
            },
            options: {
              responsive: true,
              maintainAspectRatio: false,
              animation: false,
              events: [],
              scales: {
                y: {
                  beginAtZero: true,
                  ticks: {
                    callback: function (value) {
                      return new Intl.NumberFormat('ko-KR').format(value) + '원'
                    }
                  }
                }
              },
              plugins: {
                legend: {
                  display: true,
                  position: 'top'
                },
                tooltip: {
                  enabled: false
                }
              }
            }
          })
        }

        console.log('연별 차트 캔버스 확인:', !!yearlyChartCanvas.value)
        console.log('연별 차트 부모 노드 확인:', !!yearlyChartCanvas.value?.parentNode)

        if (yearlyChartCanvas.value && yearlyChartCanvas.value.parentNode) {
          const ctx = yearlyChartCanvas.value.getContext('2d')

          // 캔버스 크기 강제 설정
          yearlyChartCanvas.value.style.width = '100%'
          yearlyChartCanvas.value.style.height = '400px'

          yearlyChart = new Chart(ctx, {
            type: 'bar',
            data: {
              labels: [],
              datasets: [
                {
                  label: '연 기준치',
                  data: [],
                  backgroundColor: 'rgba(16, 185, 129, 0.5)',
                  borderColor: 'rgb(16, 185, 129)',
                  borderWidth: 1
                },
                {
                  label: '연 사용량',
                  data: [],
                  backgroundColor: 'rgba(245, 158, 11, 0.5)',
                  borderColor: 'rgb(245, 158, 11)',
                  borderWidth: 1
                }
              ]
            },
            options: {
              responsive: true,
              maintainAspectRatio: false,
              animation: false,
              events: [],
              scales: {
                y: {
                  beginAtZero: true,
                  ticks: {
                    callback: function (value) {
                      return new Intl.NumberFormat('ko-KR').format(value) + '원'
                    }
                  }
                }
              },
              plugins: {
                legend: {
                  display: true,
                  position: 'top'
                },
                tooltip: {
                  enabled: false
                }
              }
            }
          })
        }

        console.log('차트 초기화 완료 - 월별:', !!monthlyChart, '연별:', !!yearlyChart)

      } catch (error) {
        console.error('차트 초기화 오류:', error)
      }
    }

    // 차트 업데이트
    const updateBudgetChart = async () => {
      try {
        // 차트가 초기화되지 않았거나 DOM에 연결되지 않았으면 리턴
        if (!monthlyChart || !yearlyChart ||
          !monthlyChart.canvas || !yearlyChart.canvas ||
          !monthlyChartCanvas.value || !yearlyChartCanvas.value ||
          !monthlyChartCanvas.value.parentNode || !yearlyChartCanvas.value.parentNode) {
          console.log('차트가 초기화되지 않음 또는 DOM에 연결되지 않음')
          return
        }

        // 선택된 카테고리가 없으면 리턴
        if (!selectedCategory.value) {
          console.log('선택된 카테고리가 없음')
          return
        }

        // 선택된 카테고리의 사용량 데이터 다시 로딩
        if (selectedCategory.value) {
          const targetUser = selectedCategory.value.user_name || selectedUser.value || ''
          try {
            await budgetStore.fetchBudgetUsage(targetUser)
          } catch (error) {
            console.error('사용량 데이터 로딩 오류:', error)
            if (targetUser !== '') {
              await budgetStore.fetchBudgetUsage('')
            }
          }
        }

        const budgets = budgetStore.budgets || []
        const usageData = budgetStore.budgetUsages || []

        console.log('기준치 데이터:', budgets)
        console.log('사용량 데이터:', usageData)

        // 선택된 카테고리의 데이터 처리
        const budget = selectedCategory.value
        const usage = usageData.find(u => u && u.category_id === budget.category_id)

        console.log('선택된 기준치:', budget)
        console.log('찾은 사용량 데이터:', usage)

        // 완전히 비반응형 배열 생성
        let monthlyLabels = []
        let monthlyBudgetData = []
        let monthlyUsageData = []
        let yearlyLabels = []
        let yearlyBudgetData = []
        let yearlyUsageData = []

        // 월별 기준치가 있는 경우만 추가
        if (budget.monthly_budget && budget.monthly_budget > 0) {
          monthlyLabels = [String(budget.category_name)]
          monthlyBudgetData = [Number(budget.monthly_budget) || 0]
          monthlyUsageData = [Number(usage?.monthly_used) || 0]
          console.log('월별 데이터:', { monthlyLabels, monthlyBudgetData, monthlyUsageData })
        }

        // 연별 기준치가 있는 경우만 추가  
        if (budget.yearly_budget && budget.yearly_budget > 0) {
          yearlyLabels = [String(budget.category_name)]
          yearlyBudgetData = [Number(budget.yearly_budget) || 0]
          yearlyUsageData = [Number(usage?.yearly_used) || 0]
          console.log('연별 데이터:', { yearlyLabels, yearlyBudgetData, yearlyUsageData })
        }

        // 월별 차트 업데이트 (데이터가 있는 경우만)
        if (monthlyChart && monthlyChart.canvas && monthlyLabels.length > 0) {
          try {
            monthlyChart.data.labels = monthlyLabels
            monthlyChart.data.datasets[0].data = monthlyBudgetData
            monthlyChart.data.datasets[1].data = monthlyUsageData
            monthlyChart.update('none')
          } catch (error) {
            console.error('월별 차트 업데이트 오류:', error)
          }
        } else if (monthlyChart && monthlyChart.canvas) {
          try {
            // 데이터가 없는 경우 빈 상태로 설정
            monthlyChart.data.labels = []
            monthlyChart.data.datasets[0].data = []
            monthlyChart.data.datasets[1].data = []
            monthlyChart.update('none')
          } catch (error) {
            console.error('월별 차트 클리어 오류:', error)
          }
        }

        // 연별 차트 업데이트 (데이터가 있는 경우만)
        if (yearlyChart && yearlyChart.canvas && yearlyLabels.length > 0) {
          try {
            yearlyChart.data.labels = yearlyLabels
            yearlyChart.data.datasets[0].data = yearlyBudgetData
            yearlyChart.data.datasets[1].data = yearlyUsageData
            yearlyChart.update('none')
          } catch (error) {
            console.error('연별 차트 업데이트 오류:', error)
          }
        } else if (yearlyChart && yearlyChart.canvas) {
          try {
            // 데이터가 없는 경우 빈 상태로 설정
            yearlyChart.data.labels = []
            yearlyChart.data.datasets[0].data = []
            yearlyChart.data.datasets[1].data = []
            yearlyChart.update('none')
          } catch (error) {
            console.error('연별 차트 클리어 오류:', error)
          }
        }

      } catch (error) {
        console.error('차트 업데이트 오류:', error)
      }
    }

    // 차트 제거
    const destroyCharts = () => {
      try {
        if (monthlyChart && typeof monthlyChart.destroy === 'function') {
          monthlyChart.destroy()
        }
      } catch (error) {
        console.error('월별 차트 제거 오류:', error)
      }
      monthlyChart = null

      try {
        if (yearlyChart && typeof yearlyChart.destroy === 'function') {
          yearlyChart.destroy()
        }
      } catch (error) {
        console.error('연별 차트 제거 오류:', error)
      }
      yearlyChart = null
    }

    const getBudgetUsage = (categoryId) => {
      return budgetStore.getBudgetUsageByCategory(categoryId)
    }

    const editBudget = (budget) => {
      editingBudget.value = budget
      budgetForm.value = {
        category_id: budget.category_id,
        user_name: budget.user_name,
        monthly_budget: budget.monthly_budget,
        yearly_budget: budget.yearly_budget
      }
      showEditModal.value = true
    }

    const deleteBudget = async (budget) => {
      if (!confirm(`'${budget.category_name}' 카테고리의 기준치를 삭제하시겠습니까?`)) {
        return
      }

      try {
        await budgetStore.deleteBudget(budget.id, selectedUser.value)
        // 사용량 정보도 재조회
        await budgetStore.fetchBudgetUsage(selectedUser.value)
      } catch (error) {
        console.error('기준치 삭제 실패:', error)
      }
    }

    const submitBudget = async () => {
      try {
        if (showEditModal.value) {
          // 수정
          const updateData = {
            monthly_budget: budgetForm.value.monthly_budget,
            yearly_budget: budgetForm.value.yearly_budget
          }
          await budgetStore.updateBudget(editingBudget.value.id, updateData, editingBudget.value.user_name || "")
        } else {
          // 생성 - 폼에서 입력받은 사용자 정보 사용
          await budgetStore.createBudget(budgetForm.value)
        }

        // 기준치 목록 재조회
        await loadBudgets()

        closeModals()
      } catch (error) {
        console.error('기준치 저장 실패:', error)
        // 사용자에게 오류 메시지 표시
        if (budgetStore.error) {
          ElMessage.error(budgetStore.error)
        } else {
          ElMessage.error('기준치 저장 중 오류가 발생했습니다.')
        }
      }
    }

    const resetForm = () => {
      budgetForm.value = {
        category_id: '',
        user_name: '',
        monthly_budget: 0,
        yearly_budget: 0
      }
      categorySearchText.value = ''
    }

    const closeModals = () => {
      showCreateModal.value = false
      showEditModal.value = false
      editingBudget.value = null
      resetForm()
    }

    // 라이프사이클
    onMounted(async () => {
      try {
        await userStore.fetchUsers()
        await categoryStore.fetchCategories()
        // 초기 기준치 목록 로드
        await loadBudgets()

        // DOM이 완전히 렌더링된 후 차트 초기화
        await nextTick()
        setTimeout(async () => {
          try {
            await initCharts()
            // 차트가 완전히 초기화된 후 데이터 업데이트
            setTimeout(() => {
              updateBudgetChart()
            }, 100)
          } catch (chartError) {
            console.error('차트 초기화 오류:', chartError)
          }
        }, 100)

      } catch (error) {
        console.error('BudgetManager 초기화 오류:', error)
        // 서버 연결 오류 시에도 컴포넌트는 정상 렌더링되도록 함
      }
    })

    // 컴포넌트 언마운트 시 차트 정리
    onUnmounted(() => {
      destroyCharts()
    })

    // 사용자 변경 시 기준치 목록 초기화
    watch(selectedUser, async (newUser) => {
      try {
        if (!newUser) {
          budgetStore.resetState()
          // 사용자가 선택되지 않은 경우 차트 비우기
          if (monthlyChart.value && yearlyChart.value) {
            updateBudgetChart()
          }
        } else {
          await loadBudgets()
          // 차트가 초기화된 경우에만 업데이트
          if (monthlyChart.value && yearlyChart.value) {
            setTimeout(() => {
              updateBudgetChart()
            }, 100)
          }
        }
      } catch (error) {
        console.error('사용자 변경 오류:', error)
      }
    })

    return {
      // 스토어
      budgetStore,
      userStore,
      categoryStore,

      // 반응형 상태
      selectedUser,
      showCreateModal,
      showEditModal,
      budgetForm,
      editingBudget,
      categorySearchText,
      selectedCompareMonth,
      selectedCategory,
      monthlyChartCanvas,
      yearlyChartCanvas,

      // 계산된 속성
      users,
      availableCategories,
      availableMonths,
      isFormValid,

      // 메소드
      formatMoney,
      loadBudgets,
      getBudgetUsage,
      submitBudget,
      editBudget,
      deleteBudget,
      closeModals,
      resetForm,
      fetchCategorySuggestions,
      handleCategorySelect,
      handleUserChange,
      handleCategoryClear,
      selectCategory,
      updateBudgetChart,

      // 이벤트
      emit
    }
  }
})
</script>

<style scoped>
.budget-manager {
  @apply max-w-6xl mx-auto p-6;
}
</style>
