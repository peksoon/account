import { defineStore } from 'pinia';
import axios from 'axios';
import { getApiBaseUrl } from '../config';

const BACKEND_API_BASE_URL = getApiBaseUrl();
console.log('UserStore API URL:', BACKEND_API_BASE_URL);

export const useUserStore = defineStore('user', {
  state: () => ({
    users: [],
    loading: false,
  }),
  actions: {
    async fetchUsers() {
      this.loading = true;
      try {
        const response = await axios.get(`${BACKEND_API_BASE_URL}/users`);
        this.users = response.data || [];
      } catch (error) {
        console.error('사용자 목록 로드 오류:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    async createUser(userData) {
      try {
        const response = await axios.post(`${BACKEND_API_BASE_URL}/users/create`, userData);
        await this.fetchUsers(); // 목록 새로고침
        return response.data;
      } catch (error) {
        console.error('사용자 생성 오류:', error);
        throw error;
      }
    },

    async updateUser(userData) {
      try {
        const response = await axios.put(`${BACKEND_API_BASE_URL}/users/update`, userData);
        await this.fetchUsers(); // 목록 새로고침
        return response.data;
      } catch (error) {
        console.error('사용자 수정 오류:', error);
        throw error;
      }
    },

    async deleteUser(userId) {
      try {
        await axios.delete(`${BACKEND_API_BASE_URL}/users/delete?id=${userId}`);
        await this.fetchUsers(); // 목록 새로고침
      } catch (error) {
        console.error('사용자 삭제 오류:', error);
        throw error;
      }
    },

    async forceDeleteUser(userId) {
      try {
        await axios.delete(`${BACKEND_API_BASE_URL}/users/force-delete?id=${userId}`);
        await this.fetchUsers(); // 목록 새로고침
      } catch (error) {
        console.error('사용자 강제 삭제 오류:', error);
        throw error;
      }
    },

    async checkUserUsage(userId) {
      try {
        const response = await axios.get(`${BACKEND_API_BASE_URL}/users/check-usage?id=${userId}`);
        return response.data.in_use;
      } catch (error) {
        console.error('사용자 사용 확인 오류:', error);
        throw error;
      }
    },

    getUserOptions() {
      return this.users.map(user => ({
        value: user.name,
        label: user.name,
        id: user.id
      }));
    }
  }
});
