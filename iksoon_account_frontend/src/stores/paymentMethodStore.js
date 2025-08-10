import { defineStore } from 'pinia';
import axios from 'axios';

const BACKEND_API_BASE_URL = process.env.VUE_APP_BACKEND_API_BASE_URL || 
  (process.env.NODE_ENV === 'production' ? `http://${window.location.hostname}:8080` : 'http://localhost:8080');

export const usePaymentMethodStore = defineStore('paymentMethod', {
  state: () => ({
    paymentMethods: [],
    loading: false,
  }),
  
  getters: {
    // 활성화된 결제수단만 필터링
    activePaymentMethods: (state) => 
      state.paymentMethods.filter(method => method.is_active),
    
    // ID로 결제수단 찾기 (계층구조 포함)
    getPaymentMethodById: (state) => (id) => {
      // 부모에서 먼저 찾기
      const parent = state.paymentMethods.find(method => method.id === id);
      if (parent) return parent;
      
      // 자식에서 찾기
      for (const parent of state.paymentMethods) {
        if (parent.children) {
          const child = parent.children.find(child => child.id === id);
          if (child) return child;
        }
      }
      return null;
    },
    
    // 부모 결제수단만 (카테고리)
    parentPaymentMethods: (state) => 
      state.paymentMethods.filter(method => method.is_active && !method.parent_id),
    
    // 모든 결제수단을 플랫 리스트로 (하위 호환성)
    flatPaymentMethods: (state) => {
      const flat = [];
      state.paymentMethods.forEach(parent => {
        flat.push(parent);
        if (parent.children) {
          flat.push(...parent.children);
        }
      });
      return flat.filter(method => method.is_active);
    },
  },
  
  actions: {
    // 결제수단 목록 조회
    async fetchPaymentMethods() {
      try {
        this.loading = true;
        
        const response = await axios.get(`${BACKEND_API_BASE_URL}/payment-methods`);
        this.paymentMethods = response.data || [];
        
        return this.paymentMethods;
      } catch (error) {
        console.error('결제수단 조회 오류:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    // 결제수단 생성
    async createPaymentMethod(paymentMethodData) {
      try {
        const response = await axios.post(
          `${BACKEND_API_BASE_URL}/payment-methods/create`, 
          {
            name: paymentMethodData.name,
            parent_id: paymentMethodData.parent_id || null
          }
        );
        
        // 생성 후 목록 새로고침
        await this.fetchPaymentMethods();
        
        return response.data;
      } catch (error) {
        console.error('결제수단 생성 오류:', error);
        throw error;
      }
    },
    
    // 결제수단 수정
    async updatePaymentMethod(paymentMethodId, updateData) {
      try {
        const response = await axios.put(
          `${BACKEND_API_BASE_URL}/payment-methods/update?id=${paymentMethodId}`, 
          updateData
        );
        
        // 수정 후 목록 새로고침
        await this.fetchPaymentMethods();
        
        return response.data;
      } catch (error) {
        console.error('결제수단 수정 오류:', error);
        throw error;
      }
    },
    
    // 결제수단 삭제
    async deletePaymentMethod(paymentMethodId) {
      try {
        const response = await axios.delete(
          `${BACKEND_API_BASE_URL}/payment-methods/delete?id=${paymentMethodId}`
        );
        
        // 삭제 후 목록 새로고침
        await this.fetchPaymentMethods();
        
        return response.data;
      } catch (error) {
        console.error('결제수단 삭제 오류:', error);
        throw error;
      }
    },
    
    // 결제수단 강제 삭제
    async forceDeletePaymentMethod(paymentMethodId) {
      try {
        const response = await axios.delete(
          `${BACKEND_API_BASE_URL}/payment-methods/force-delete?id=${paymentMethodId}`
        );
        
        // 삭제 후 목록 새로고침
        await this.fetchPaymentMethods();
        
        return response.data;
      } catch (error) {
        console.error('결제수단 강제 삭제 오류:', error);
        throw error;
      }
    },
    
    // 결제수단 목록 초기화
    clearPaymentMethods() {
      this.paymentMethods = [];
    },
    
    // 선택 옵션용 리스트 생성
    getPaymentMethodOptions() {
      return this.activePaymentMethods.map(method => ({
        id: method.id,
        name: method.name,
        fullName: method.name,
        isActive: method.is_active
      }));
    },
  },
});