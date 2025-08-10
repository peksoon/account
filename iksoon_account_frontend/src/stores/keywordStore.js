import { defineStore } from 'pinia';
import axios from 'axios';

const BACKEND_API_BASE_URL = process.env.VUE_APP_BACKEND_API_BASE_URL || 
  (process.env.NODE_ENV === 'production' ? `http://${window.location.hostname}:8080` : 'http://localhost:8080');

export const useKeywordStore = defineStore('keyword', {
  state: () => ({
    keywords: [],
    suggestions: [],
    loading: false,
  }),
  
  getters: {
    // 카테고리별 키워드 필터링
    getKeywordsByCategory: (state) => (categoryId) => 
      state.keywords.filter(keyword => keyword.category_id === categoryId),
    
    // ID로 키워드 찾기
    getKeywordById: (state) => (id) => 
      state.keywords.find(keyword => keyword.id === id),
    
    // 이름으로 키워드 찾기
    getKeywordByName: (state) => (categoryId, name) => 
      state.keywords.find(keyword => 
        keyword.category_id === categoryId && keyword.name === name
      ),
  },
  
  actions: {
    // 키워드 자동완성 제안 조회
    async getKeywordSuggestions(categoryId, query = '', limit = 10) {
      try {
        this.loading = true;
        
        let url = `${BACKEND_API_BASE_URL}/keywords/suggestions?category_id=${categoryId}&limit=${limit}`;
        if (query) {
          url += `&q=${encodeURIComponent(query)}`;
        }
        
        const response = await axios.get(url);
        this.suggestions = response.data || [];
        
        return this.suggestions;
      } catch (error) {
        console.error('키워드 제안 조회 오류:', error);
        return [];
      } finally {
        this.loading = false;
      }
    },
    
    // 카테고리별 키워드 목록 조회
    async fetchKeywordsByCategory(categoryId) {
      try {
        this.loading = true;
        
        const response = await axios.get(
          `${BACKEND_API_BASE_URL}/keywords/category?category_id=${categoryId}`
        );
        
        const keywords = response.data || [];
        
        // 기존 키워드 목록에서 해당 카테고리 키워드 제거 후 새로 추가
        this.keywords = this.keywords.filter(k => k.category_id !== categoryId);
        this.keywords.push(...keywords);
        
        return keywords;
      } catch (error) {
        console.error('카테고리별 키워드 조회 오류:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    // 키워드 생성 또는 업데이트 (Upsert)
    async upsertKeyword(categoryId, keywordName) {
      try {
        const response = await axios.post(`${BACKEND_API_BASE_URL}/keywords/upsert`, {
          category_id: categoryId,
          name: keywordName
        });
        
        // 해당 카테고리의 키워드 목록 새로고침
        await this.fetchKeywordsByCategory(categoryId);
        
        return response.data;
      } catch (error) {
        console.error('키워드 생성/업데이트 오류:', error);
        throw error;
      }
    },
    
    // 키워드 삭제
    async deleteKeyword(keywordId) {
      try {
        const response = await axios.delete(
          `${BACKEND_API_BASE_URL}/keywords/delete?id=${keywordId}`
        );
        
        // 로컬 상태에서 삭제된 키워드 제거
        this.keywords = this.keywords.filter(k => k.id !== keywordId);
        
        return response.data;
      } catch (error) {
        console.error('키워드 삭제 오류:', error);
        throw error;
      }
    },
    
    // 키워드 목록 초기화
    clearKeywords() {
      this.keywords = [];
      this.suggestions = [];
    },
    
    // 특정 카테고리의 키워드 클리어
    clearKeywordsByCategory(categoryId) {
      this.keywords = this.keywords.filter(k => k.category_id !== categoryId);
    },
    
    // 키워드 사용 (사용 횟수 증가를 위한 upsert 호출)
    async useKeyword(categoryId, keywordName) {
      if (!keywordName || !categoryId) return;
      
      try {
        await this.upsertKeyword(categoryId, keywordName);
      } catch (error) {
        // 키워드 사용 기록 실패는 무시 (주요 기능에 영향 주지 않음)
        console.warn('키워드 사용 기록 실패:', error);
      }
    },
  },
});
