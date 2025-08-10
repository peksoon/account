import { defineStore } from 'pinia';
import axios from 'axios';

const BACKEND_API_BASE_URL = process.env.VUE_APP_BACKEND_API_BASE_URL || 'http://133.186.153.179:8080';

export const useDepositPathStore = defineStore('depositPath', {
  state: () => ({
    depositPaths: [],
    loading: false,
  }),
  
  getters: {
    // 활성화된 입금경로만 필터링
    activeDepositPaths: (state) => 
      state.depositPaths.filter(path => path.is_active),
    
    // ID로 입금경로 찾기
    getDepositPathById: (state) => (id) => 
      state.depositPaths.find(path => path.id === id),
  },
  
  actions: {
    // 입금경로 목록 조회
    async fetchDepositPaths() {
      try {
        this.loading = true;
        
        const response = await axios.get(`${BACKEND_API_BASE_URL}/deposit-paths`);
        this.depositPaths = response.data || [];
        
        return this.depositPaths;
      } catch (error) {
        console.error('입금경로 조회 오류:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    // 입금경로 생성
    async createDepositPath(depositPathData) {
      try {
        const response = await axios.post(
          `${BACKEND_API_BASE_URL}/deposit-paths/create`, 
          depositPathData
        );
        
        // 생성 후 목록 새로고침
        await this.fetchDepositPaths();
        
        return response.data;
      } catch (error) {
        console.error('입금경로 생성 오류:', error);
        throw error;
      }
    },
    
    // 입금경로 수정
    async updateDepositPath(depositPathId, updateData) {
      try {
        const response = await axios.put(
          `${BACKEND_API_BASE_URL}/deposit-paths/update?id=${depositPathId}`, 
          updateData
        );
        
        // 수정 후 목록 새로고침
        await this.fetchDepositPaths();
        
        return response.data;
      } catch (error) {
        console.error('입금경로 수정 오류:', error);
        throw error;
      }
    },
    
    // 입금경로 삭제
    async deleteDepositPath(depositPathId) {
      try {
        const response = await axios.delete(
          `${BACKEND_API_BASE_URL}/deposit-paths/delete?id=${depositPathId}`
        );
        
        // 삭제 후 목록 새로고침
        await this.fetchDepositPaths();
        
        return response.data;
      } catch (error) {
        console.error('입금경로 삭제 오류:', error);
        throw error;
      }
    },
    
    // 입금경로 강제 삭제
    async forceDeleteDepositPath(depositPathId) {
      try {
        const response = await axios.delete(
          `${BACKEND_API_BASE_URL}/deposit-paths/force-delete?id=${depositPathId}`
        );
        
        // 삭제 후 목록 새로고침
        await this.fetchDepositPaths();
        
        return response.data;
      } catch (error) {
        console.error('입금경로 강제 삭제 오류:', error);
        throw error;
      }
    },
    
    // 입금경로 목록 초기화
    clearDepositPaths() {
      this.depositPaths = [];
    },
    
    // 선택 옵션용 리스트 생성
    getDepositPathOptions() {
      return this.activeDepositPaths.map(path => ({
        id: path.id,
        name: path.name,
        fullName: path.name,
        isActive: path.is_active
      }));
    },
  },
});
