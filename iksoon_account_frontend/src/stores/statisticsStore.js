import { defineStore } from 'pinia';
import axios from 'axios';
import { getApiBaseUrl } from '../config';

const BACKEND_API_BASE_URL = getApiBaseUrl();
console.log('StatisticsStore API URL:', BACKEND_API_BASE_URL);

export const useStatisticsStore = defineStore('statistics', {
  state: () => ({
    statistics: null,
    keywordStatistics: null,
    loading: false,
    keywordLoading: false,
  }),
  
  getters: {
    // 카테고리별 총합
    categoryTotals: (state) => {
      if (!state.statistics?.categories) return {};
      
      const totals = {};
      state.statistics.categories.forEach(category => {
        totals[category.category_id] = category.total_amount;
      });
      return totals;
    },
    
    // 상위 카테고리 (상위 N개)
    topCategories: (state) => (limit = 5) => {
      if (!state.statistics?.categories) return [];
      
      return state.statistics.categories
        .slice(0, limit)
        .map(category => ({
          ...category,
          percentage: parseFloat(category.percentage.toFixed(1))
        }));
    },
    
    // 통계 요약 정보
    summaryInfo: (state) => {
      if (!state.statistics) return null;
      
      return {
        total_amount: state.statistics.total_amount,
        total_count: state.statistics.total_count,
        average_amount: state.statistics.total_count > 0 
          ? Math.round(state.statistics.total_amount / state.statistics.total_count)
          : 0,
        period: state.statistics.period,
        top_category: state.statistics.top_category
      };
    },
    
    // 차트용 색상 팔레트
    chartColors: () => [
      '#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7',
      '#DDA0DD', '#98D8C8', '#F7DC6F', '#BB8FCE', '#85C1E9',
      '#F8C471', '#82E0AA', '#F1948A', '#85DCDB', '#D7DBDD'
    ],
  },
  
  actions: {
    // 통계 데이터 조회
    async fetchStatistics(params = {}) {
      try {
        this.loading = true;
        
        const queryParams = new URLSearchParams();
        
        // 기본 파라미터 설정
        queryParams.append('type', params.type || 'month');
        queryParams.append('category', params.category || 'out');
        
        // 사용자 파라미터 추가 (기준치 정보용)
        if (params.user) {
          queryParams.append('user', params.user);
        }
        
        // 커스텀 날짜 범위가 있는 경우
        if (params.start_date && params.end_date) {
          queryParams.append('start_date', params.start_date);
          queryParams.append('end_date', params.end_date);
        }
        
        const response = await axios.get(
          `${BACKEND_API_BASE_URL}/statistics?${queryParams.toString()}`
        );
        
        this.statistics = response.data;
        
        return this.statistics;
      } catch (error) {
        console.error('통계 조회 오류:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    // 카테고리별 키워드 통계 조회
    async fetchKeywordStatistics(params = {}) {
      try {
        this.keywordLoading = true;
        
        const queryParams = new URLSearchParams();
        
        // 필수 파라미터
        if (!params.category_id) {
          throw new Error('카테고리 ID가 필요합니다.');
        }
        queryParams.append('category_id', params.category_id);
        
        // 기본 파라미터 설정
        queryParams.append('type', params.type || 'month');
        queryParams.append('category', params.category || 'out');
        
        // 커스텀 날짜 범위가 있는 경우
        if (params.start_date && params.end_date) {
          queryParams.append('start_date', params.start_date);
          queryParams.append('end_date', params.end_date);
        }
        
        const response = await axios.get(
          `${BACKEND_API_BASE_URL}/statistics/category-keywords?${queryParams.toString()}`
        );
        
        this.keywordStatistics = response.data;
        
        return this.keywordStatistics;
      } catch (error) {
        console.error('키워드 통계 조회 오류:', error);
        throw error;
      } finally {
        this.keywordLoading = false;
      }
    },
    
    // 특정 기간의 통계 조회
    async fetchPeriodStatistics(type, category = 'out', startDate = null, endDate = null) {
      const params = {
        type,
        category
      };
      
      if (startDate && endDate) {
        params.start_date = this.formatDate(startDate);
        params.end_date = this.formatDate(endDate);
      }
      
      return await this.fetchStatistics(params);
    },
    
    // 월별 통계 조회
    async fetchMonthlyStatistics(category = 'out') {
      return await this.fetchPeriodStatistics('month', category);
    },
    
    // 주간 통계 조회
    async fetchWeeklyStatistics(category = 'out') {
      return await this.fetchPeriodStatistics('week', category);
    },
    
    // 연간 통계 조회
    async fetchYearlyStatistics(category = 'out') {
      return await this.fetchPeriodStatistics('year', category);
    },
    
    // 전체 통계 조회
    async fetchAllStatistics(category = 'out') {
      return await this.fetchPeriodStatistics('all', category);
    },
    
    // 통계 데이터 초기화
    clearStatistics() {
      this.statistics = null;
      this.keywordStatistics = null;
    },
    
    // 키워드 통계 초기화
    clearKeywordStatistics() {
      this.keywordStatistics = null;
    },
    
    // 날짜 포맷팅 (YYYY-MM-DD)
    formatDate(date) {
      if (!date) return null;
      
      if (typeof date === 'string') return date;
      
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      
      return `${year}-${month}-${day}`;
    },
    
    // 통계 데이터 내보내기용 포맷
    exportStatistics() {
      if (!this.statistics) return null;
      
      return {
        period: this.statistics.period,
        summary: {
          total_amount: this.statistics.total_amount,
          total_count: this.statistics.total_count,
          average_amount: this.summaryInfo.average_amount
        },
        categories: this.statistics.categories.map(category => ({
          name: category.category_name,
          amount: category.total_amount,
          count: category.count,
          percentage: category.percentage
        })),
        exported_at: new Date().toISOString()
      };
    },
    
    // 비교 통계 생성 (이전 기간과 비교)
    async fetchComparisonStatistics(currentParams, previousParams) {
      try {
        const [currentStats, previousStats] = await Promise.all([
          this.fetchStatistics(currentParams),
          this.fetchStatistics(previousParams)
        ]);
        
        const comparison = {
          current: currentStats,
          previous: previousStats,
          changes: {
            amount_change: currentStats.total_amount - previousStats.total_amount,
            amount_change_percent: previousStats.total_amount > 0 
              ? ((currentStats.total_amount - previousStats.total_amount) / previousStats.total_amount * 100)
              : 0,
            count_change: currentStats.total_count - previousStats.total_count,
            count_change_percent: previousStats.total_count > 0
              ? ((currentStats.total_count - previousStats.total_count) / previousStats.total_count * 100)
              : 0
          }
        };
        
        return comparison;
      } catch (error) {
        console.error('비교 통계 조회 오류:', error);
        throw error;
      }
    },
  },
});
