import { defineStore } from 'pinia'
import axios from 'axios'

const API_BASE_URL = process.env.VUE_APP_BACKEND_API_BASE_URL || 'http://localhost:8080'

export const useBudgetStore = defineStore('budget', {
  state: () => ({
    budgets: [], // 카테고리별 기준치 목록
    budgetUsages: [], // 기준치 사용량 목록
    loading: false,
    error: null
  }),

  actions: {
    /**
     * 카테고리 기준치 목록 조회
     * @param {string} userName - 사용자명
     * @param {number} categoryId - 카테고리 ID (선택적)
     */
    async fetchBudgets(userName, categoryId = null) {
      this.loading = true
      this.error = null
      
      try {
        let url = `${API_BASE_URL}/category-budgets?user=${encodeURIComponent(userName)}`
        if (categoryId) {
          url += `&category_id=${categoryId}`
        }
        
        const response = await axios.get(url)
        this.budgets = response.data || []
        
        console.log('기준치 목록 조회 성공:', this.budgets)
      } catch (error) {
        console.error('기준치 목록 조회 실패:', error)
        if (error.code === 'ERR_NETWORK' || error.message.includes('ERR_CONNECTION_REFUSED')) {
          this.error = '서버에 연결할 수 없습니다. 백엔드 서버가 실행 중인지 확인해주세요.'
        } else {
          this.error = error.response?.data?.message || '기준치 조회 중 오류가 발생했습니다.'
        }
        this.budgets = []
      } finally {
        this.loading = false
      }
    },

    /**
     * 카테고리 기준치 생성
     * @param {Object} budgetData - 기준치 데이터
     * @param {number} budgetData.category_id - 카테고리 ID
     * @param {string} budgetData.user_name - 사용자명
     * @param {number} budgetData.monthly_budget - 월 기준치
     * @param {number} budgetData.yearly_budget - 연 기준치
     */
    async createBudget(budgetData) {
      this.loading = true
      this.error = null
      
      try {
        const response = await axios.post(`${API_BASE_URL}/category-budgets/create`, budgetData)
        
        console.log('기준치 생성 성공:', response.data)
        
        // 생성 후 목록 재조회
        await this.fetchBudgets(budgetData.user_name)
        
        return response.data
      } catch (error) {
        console.error('기준치 생성 실패:', error)
        if (error.code === 'ERR_NETWORK' || error.message.includes('ERR_CONNECTION_REFUSED')) {
          this.error = '서버에 연결할 수 없습니다. 백엔드 서버가 실행 중인지 확인해주세요.'
        } else if (error.response?.status === 409) {
          this.error = error.response?.data?.message || '해당 카테고리에 대한 기준치가 이미 존재합니다.'
        } else {
          this.error = error.response?.data?.message || '기준치 생성 중 오류가 발생했습니다.'
        }
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 카테고리 기준치 수정
     * @param {number} id - 기준치 ID
     * @param {Object} updateData - 수정할 데이터
     * @param {number} updateData.monthly_budget - 월 기준치
     * @param {number} updateData.yearly_budget - 연 기준치
     * @param {string} userName - 사용자명 (재조회용)
     */
    async updateBudget(id, updateData, userName) {
      this.loading = true
      this.error = null
      
      try {
        const response = await axios.put(`${API_BASE_URL}/category-budgets/update?id=${id}`, updateData)
        
        console.log('기준치 수정 성공:', response.data)
        
        // 수정 후 목록 재조회
        await this.fetchBudgets(userName)
        
        return response.data
      } catch (error) {
        console.error('기준치 수정 실패:', error)
        this.error = error.response?.data?.message || '기준치 수정 중 오류가 발생했습니다.'
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 카테고리 기준치 삭제
     * @param {number} id - 기준치 ID
     * @param {string} userName - 사용자명 (재조회용)
     */
    async deleteBudget(id, userName) {
      this.loading = true
      this.error = null
      
      try {
        const response = await axios.delete(`${API_BASE_URL}/category-budgets/delete?id=${id}`)
        
        console.log('기준치 삭제 성공:', response.data)
        
        // 삭제 후 목록 재조회
        await this.fetchBudgets(userName)
        
        return response.data
      } catch (error) {
        console.error('기준치 삭제 실패:', error)
        this.error = error.response?.data?.message || '기준치 삭제 중 오류가 발생했습니다.'
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 기준치 사용량 조회
     * @param {string} userName - 사용자명
     * @param {number} categoryId - 카테고리 ID (선택적)
     */
    async fetchBudgetUsage(userName, categoryId = null) {
      this.loading = true
      this.error = null
      
      try {
        let url = `${API_BASE_URL}/category-budgets/usage?user=${encodeURIComponent(userName)}`
        if (categoryId) {
          url += `&category_id=${categoryId}`
        }
        
        const response = await axios.get(url)
        
        if (categoryId) {
          // 특정 카테고리 사용량
          return response.data
        } else {
          // 모든 카테고리 사용량
          this.budgetUsages = response.data || []
          console.log('기준치 사용량 조회 성공:', this.budgetUsages)
          return this.budgetUsages
        }
      } catch (error) {
        console.error('기준치 사용량 조회 실패:', error)
        this.error = error.response?.data?.message || '기준치 사용량 조회 중 오류가 발생했습니다.'
        
        if (categoryId) {
          return null
        } else {
          this.budgetUsages = []
          return []
        }
      } finally {
        this.loading = false
      }
    },

    /**
     * 기준치 정보를 포함한 지출 추가
     * @param {Object} outAccountData - 지출 데이터
     */
    async createOutAccountWithBudget(outAccountData) {
      this.loading = true
      this.error = null
      
      try {
        const response = await axios.post(`${API_BASE_URL}/v2/out-account/insert-with-budget`, outAccountData)
        
        console.log('기준치 포함 지출 추가 성공:', response.data)
        
        // 기준치 사용량 재조회
        if (outAccountData.user) {
          await this.fetchBudgetUsage(outAccountData.user)
        }
        
        return response.data
      } catch (error) {
        console.error('기준치 포함 지출 추가 실패:', error)
        this.error = error.response?.data?.message || '지출 추가 중 오류가 발생했습니다.'
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 특정 카테고리의 기준치 정보 찾기
     * @param {number} categoryId - 카테고리 ID
     * @returns {Object|null} - 기준치 정보
     */
    getBudgetByCategory(categoryId) {
      return this.budgets.find(budget => budget.category_id === categoryId) || null
    },

    /**
     * 특정 카테고리의 기준치 사용량 정보 찾기
     * @param {number} categoryId - 카테고리 ID
     * @returns {Object|null} - 기준치 사용량 정보
     */
    getBudgetUsageByCategory(categoryId) {
      return this.budgetUsages.find(usage => usage.category_id === categoryId) || null
    },

    /**
     * 기준치가 초과된 카테고리 목록 조회
     * @returns {Array} - 초과된 카테고리 목록
     */
    getOverBudgetCategories() {
      return this.budgetUsages.filter(usage => usage.is_monthly_over || usage.is_yearly_over)
    },

    /**
     * 에러 초기화
     */
    clearError() {
      this.error = null
    },

    /**
     * 상태 초기화
     */
    resetState() {
      this.budgets = []
      this.budgetUsages = []
      this.loading = false
      this.error = null
    },

    /**
     * 월별 기준치만 수정
     * @param {number} categoryId - 카테고리 ID
     * @param {string} userName - 사용자명 (선택사항)
     * @param {number} amount - 월별 기준치 금액
     */
    async updateMonthlyBudget(categoryId, userName, amount) {
      this.loading = true
      this.error = null

      try {
        const response = await axios.put(`${API_BASE_URL}/category-budgets/update-monthly`, {
          category_id: categoryId,
          user_name: userName || "",
          amount: amount
        })

        console.log('월별 기준치 수정 성공:', response.data)
        return response.data
      } catch (error) {
        console.error('월별 기준치 수정 실패:', error)
        this.error = error.response?.data?.message || '월별 기준치 수정 중 오류가 발생했습니다.'
        throw error
      } finally {
        this.loading = false
      }
    },

    /**
     * 연별 기준치만 수정
     * @param {number} categoryId - 카테고리 ID
     * @param {string} userName - 사용자명 (선택사항)
     * @param {number} amount - 연별 기준치 금액
     */
    async updateYearlyBudget(categoryId, userName, amount) {
      this.loading = true
      this.error = null

      try {
        const response = await axios.put(`${API_BASE_URL}/category-budgets/update-yearly`, {
          category_id: categoryId,
          user_name: userName || "",
          amount: amount
        })

        console.log('연별 기준치 수정 성공:', response.data)
        return response.data
      } catch (error) {
        console.error('연별 기준치 수정 실패:', error)
        this.error = error.response?.data?.message || '연별 기준치 수정 중 오류가 발생했습니다.'
        throw error
      } finally {
        this.loading = false
      }
    }
  },

  getters: {
    /**
     * 기준치가 설정된 카테고리 수
     */
    budgetCategoriesCount: (state) => state.budgets.length,

    /**
     * 기준치 초과 카테고리 수
     */
    overBudgetCategoriesCount: (state) => {
      return state.budgetUsages.filter(usage => usage.is_monthly_over || usage.is_yearly_over).length
    },

    /**
     * 전체 월 기준치 합계
     */
    totalMonthlyBudget: (state) => {
      return state.budgets.reduce((total, budget) => total + (budget.monthly_budget || 0), 0)
    },

    /**
     * 전체 연 기준치 합계
     */
    totalYearlyBudget: (state) => {
      return state.budgets.reduce((total, budget) => total + (budget.yearly_budget || 0), 0)
    },

    /**
     * 전체 월 사용량 합계
     */
    totalMonthlyUsed: (state) => {
      return state.budgetUsages.reduce((total, usage) => total + (usage.monthly_used || 0), 0)
    },

    /**
     * 전체 연 사용량 합계
     */
    totalYearlyUsed: (state) => {
      return state.budgetUsages.reduce((total, usage) => total + (usage.yearly_used || 0), 0)
    }
  }
})
