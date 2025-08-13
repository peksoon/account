import { defineStore } from 'pinia';
import { useCrudApi } from '../composables/useApi';

/**
 * 사용자 관리를 위한 Pinia 스토어
 * 사용자 CRUD 작업과 상태 관리를 담당합니다.
 */
export const useUserStore = defineStore('user', {
  state: () => ({
    users: [],
    loading: false,
  }),

  getters: {
    /**
     * 사용자 옵션 목록 (셀렉트 박스용)
     * @returns {Array} 사용자 옵션 배열
     */
    userOptions: (state) => state.users.map(user => ({
      value: user.name,
      label: user.name,
      id: user.id
    })),

    /**
     * ID로 사용자 찾기
     * @param {number} id - 사용자 ID
     * @returns {Object|undefined} 사용자 객체
     */
    getUserById: (state) => (id) => state.users.find(user => user.id === id),

    /**
     * 이름으로 사용자 찾기
     * @param {string} name - 사용자 이름
     * @returns {Object|undefined} 사용자 객체
     */
    getUserByName: (state) => (name) => state.users.find(user => user.name === name),

    /**
     * 활성 사용자 목록
     * @returns {Array} 활성 사용자 배열
     */
    activeUsers: (state) => state.users.filter(user => user.is_active)
  },

  actions: {
    /**
     * 사용자 목록을 서버에서 가져옵니다.
     */
    async fetchUsers() {
      const api = useCrudApi('users', '사용자');
      this.loading = true;
      
      try {
        const response = await api.fetchList();
        this.users = response.data || [];
      } catch (error) {
        console.error('사용자 목록 로드 오류:', error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    /**
     * 새 사용자를 생성합니다.
     * @param {Object} userData - 사용자 데이터
     * @param {string} userData.name - 사용자 이름
     * @param {string} userData.email - 사용자 이메일 (선택사항)
     * @returns {Promise<Object>} 생성된 사용자 데이터
     */
    async createUser(userData) {
      const api = useCrudApi('users', '사용자');
      
      try {
        const response = await api.create(userData);
        await this.fetchUsers(); // 목록 새로고침
        return response.data;
      } catch (error) {
        console.error('사용자 생성 오류:', error);
        throw error;
      }
    },

    /**
     * 사용자 정보를 수정합니다.
     * @param {Object} userData - 수정할 사용자 데이터
     * @param {number} userData.id - 사용자 ID
     * @param {string} userData.name - 사용자 이름
     * @param {string} userData.email - 사용자 이메일 (선택사항)
     * @returns {Promise<Object>} 수정된 사용자 데이터
     */
    async updateUser(userData) {
      const api = useCrudApi('users', '사용자');
      
      try {
        const response = await api.update(userData.id, userData);
        await this.fetchUsers(); // 목록 새로고침
        return response.data;
      } catch (error) {
        console.error('사용자 수정 오류:', error);
        throw error;
      }
    },

    /**
     * 사용자를 삭제합니다.
     * @param {number} userId - 삭제할 사용자 ID
     */
    async deleteUser(userId) {
      const api = useCrudApi('users', '사용자');
      
      try {
        await api.remove(userId);
        await this.fetchUsers(); // 목록 새로고침
      } catch (error) {
        console.error('사용자 삭제 오류:', error);
        throw error;
      }
    },

    /**
     * 사용자를 강제로 삭제합니다. (참조 데이터 포함)
     * @param {number} userId - 강제 삭제할 사용자 ID
     */
    async forceDeleteUser(userId) {
      const api = useCrudApi('users', '사용자');
      
      try {
        await api.forceRemove(userId);
        await this.fetchUsers(); // 목록 새로고침
      } catch (error) {
        console.error('사용자 강제 삭제 오류:', error);
        throw error;
      }
    },

    /**
     * 사용자 사용 여부를 확인합니다.
     * @param {number} userId - 확인할 사용자 ID
     * @returns {Promise<boolean>} 사용 중이면 true
     */
    async checkUserUsage(userId) {
      const api = useCrudApi('users', '사용자');
      
      try {
        const response = await api.get(`/users/check-usage?id=${userId}`);
        return response.data.in_use;
      } catch (error) {
        console.error('사용자 사용 확인 오류:', error);
        throw error;
      }
    },

    /**
     * 사용자 옵션 목록을 반환합니다. (Deprecated: userOptions getter 사용 권장)
     * @returns {Array} 사용자 옵션 배열
     * @deprecated 대신 userOptions getter를 사용하세요.
     */
    getUserOptions() {
      return this.userOptions;
    }
  }
});
