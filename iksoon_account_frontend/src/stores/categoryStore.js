import { defineStore } from 'pinia';
import axios from 'axios';
import { getApiBaseUrl } from '../config';

const BACKEND_API_BASE_URL = getApiBaseUrl();
console.log('CategoryStore API URL:', BACKEND_API_BASE_URL);

export const useCategoryStore = defineStore('category', {
  state: () => ({
    categories: [],
    loading: false,
  }),
  
  getters: {
    // 지출 카테고리만 필터링
    outCategories: (state) => state.categories.filter(cat => cat.type === 'out'),
    
    // 수입 카테고리만 필터링
    inCategories: (state) => state.categories.filter(cat => cat.type === 'in'),
    
    // ID로 카테고리 찾기
    getCategoryById: (state) => (id) => state.categories.find(cat => cat.id === id),
    
    // 이름으로 카테고리 찾기
    getCategoryByName: (state) => (name, type) => 
      state.categories.find(cat => cat.name === name && cat.type === type),
  },
  
  actions: {
    // 카테고리 목록 조회
    async fetchCategories(type = '') {
      try {
        this.loading = true;
        let url = `${BACKEND_API_BASE_URL}/categories`;
        if (type) {
          url += `?type=${type}`;
        }
        
        const response = await axios.get(url);
        this.categories = response.data || [];
        
        return this.categories;
      } catch (error) {
        console.error('카테고리 조회 오류:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    // 카테고리 생성
    async createCategory(categoryData) {
      try {
        const response = await axios.post(`${BACKEND_API_BASE_URL}/categories/create`, categoryData);
        
        // 생성 후 목록 새로고침
        await this.fetchCategories();
        
        return response.data;
      } catch (error) {
        console.error('카테고리 생성 오류:', error);
        throw error;
      }
    },
    
    // 카테고리 수정
    async updateCategory(updateData) {
      try {
        const response = await axios.put(
          `${BACKEND_API_BASE_URL}/categories/update?id=${updateData.id}`, 
          {
            name: updateData.name,
            type: updateData.type
          }
        );
        
        // 수정 후 목록 새로고침
        await this.fetchCategories();
        
        return response.data;
      } catch (error) {
        console.error('카테고리 수정 오류:', error);
        throw error;
      }
    },
    
    // 카테고리 삭제
    async deleteCategory(categoryId) {
      try {
        if (!categoryId || categoryId <= 0) {
          throw new Error('유효하지 않은 카테고리 ID입니다');
        }
        
        const url = `${BACKEND_API_BASE_URL}/categories/${encodeURIComponent(categoryId)}`;
        
        const response = await axios.delete(url);
        
        // 삭제 후 목록 새로고침
        await this.fetchCategories();
        
        return response.data;
      } catch (error) {
        console.error('카테고리 삭제 오류:', error);
        console.error('삭제 오류 상세:', error.response?.data);
        throw error;
      }
    },
    
    // 카테고리 강제 삭제 (관련 데이터도 함께 삭제)
    async forceDeleteCategory(categoryId) {
      try {
        if (!categoryId || categoryId <= 0) {
          throw new Error('유효하지 않은 카테고리 ID입니다');
        }
        
        const url = `${BACKEND_API_BASE_URL}/categories/${encodeURIComponent(categoryId)}/force-delete`;
        
        const response = await axios.delete(url);
        
        // 삭제 후 목록 새로고침
        await this.fetchCategories();
        
        return response.data;
      } catch (error) {
        console.error('카테고리 강제 삭제 오류:', error);
        console.error('오류 상세:', error.response?.data);
        throw error;
      }
    },
    
    // 카테고리 목록 초기화
    clearCategories() {
      this.categories = [];
    },
    
    // 특정 타입의 카테고리만 가져오기
    async fetchCategoriesByType(type) {
      return await this.fetchCategories(type);
    },
  },
});
