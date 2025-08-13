import axios from 'axios';
import { ref } from 'vue';
import { ElMessage } from 'element-plus';
import { getApiBaseUrl } from '../config';

const BACKEND_API_BASE_URL = getApiBaseUrl();

/**
 * API 호출을 위한 컴포저블 함수
 * 공통적인 로딩 상태, 에러 처리, 성공 메시지를 관리합니다.
 */
export function useApi() {
  const loading = ref(false);
  const error = ref(null);

  /**
   * API 요청을 실행하는 공통 함수
   * @param {Function} apiCall - 실행할 API 호출 함수
   * @param {Object} options - 옵션 설정
   * @param {string} options.successMessage - 성공 시 표시할 메시지
   * @param {string} options.errorMessage - 에러 시 표시할 메시지
   * @param {boolean} options.showSuccessMessage - 성공 메시지 표시 여부 (기본: true)
   * @param {boolean} options.showErrorMessage - 에러 메시지 표시 여부 (기본: true)
   * @returns {Promise} API 호출 결과
   */
  const execute = async (apiCall, options = {}) => {
    const {
      successMessage = '',
      errorMessage = '요청 처리 중 오류가 발생했습니다.',
      showSuccessMessage = false,
      showErrorMessage = true
    } = options;

    loading.value = true;
    error.value = null;

    try {
      const result = await apiCall();
      
      if (successMessage && showSuccessMessage) {
        ElMessage.success(successMessage);
      }
      
      return result;
    } catch (err) {
      error.value = err;
      console.error('API 요청 오류:', err);
      
      if (showErrorMessage) {
        const message = err.response?.data?.message || errorMessage;
        ElMessage.error(message);
      }
      
      throw err;
    } finally {
      loading.value = false;
    }
  };

  /**
   * GET 요청을 수행하는 헬퍼 함수
   * @param {string} endpoint - API 엔드포인트
   * @param {Object} params - 쿼리 파라미터
   * @returns {Promise} API 응답 데이터
   */
  const get = (endpoint, params = {}) => {
    return execute(() => axios.get(`${BACKEND_API_BASE_URL}${endpoint}`, { params }));
  };

  /**
   * POST 요청을 수행하는 헬퍼 함수
   * @param {string} endpoint - API 엔드포인트
   * @param {Object} data - 전송할 데이터
   * @param {Object} options - 옵션 설정
   * @returns {Promise} API 응답 데이터
   */
  const post = (endpoint, data = {}, options = {}) => {
    return execute(() => axios.post(`${BACKEND_API_BASE_URL}${endpoint}`, data), options);
  };

  /**
   * PUT 요청을 수행하는 헬퍼 함수
   * @param {string} endpoint - API 엔드포인트
   * @param {Object} data - 전송할 데이터
   * @param {Object} options - 옵션 설정
   * @returns {Promise} API 응답 데이터
   */
  const put = (endpoint, data = {}, options = {}) => {
    return execute(() => axios.put(`${BACKEND_API_BASE_URL}${endpoint}`, data), options);
  };

  /**
   * DELETE 요청을 수행하는 헬퍼 함수
   * @param {string} endpoint - API 엔드포인트
   * @param {Object} params - 쿼리 파라미터
   * @param {Object} options - 옵션 설정
   * @returns {Promise} API 응답 데이터
   */
  const del = (endpoint, params = {}, options = {}) => {
    return execute(() => axios.delete(`${BACKEND_API_BASE_URL}${endpoint}`, { params }), options);
  };

  return {
    loading,
    error,
    execute,
    get,
    post,
    put,
    delete: del
  };
}

/**
 * 공통 CRUD 작업을 위한 컴포저블 함수
 * @param {string} entityName - 엔티티 이름 (예: 'users', 'categories')
 * @param {string} displayName - 표시용 이름 (예: '사용자', '카테고리')
 */
export function useCrudApi(entityName, displayName) {
  const api = useApi();

  /**
   * 목록 조회
   * @param {Object} params - 쿼리 파라미터
   * @returns {Promise} 엔티티 목록
   */
  const fetchList = (params = {}) => {
    return api.get(`/${entityName}`, params);
  };

  /**
   * 단일 항목 조회
   * @param {number|string} id - 엔티티 ID
   * @returns {Promise} 엔티티 데이터
   */
  const fetchById = (id) => {
    return api.get(`/${entityName}/${id}`);
  };

  /**
   * 새 항목 생성
   * @param {Object} data - 생성할 데이터
   * @returns {Promise} 생성된 엔티티 데이터
   */
  const create = (data) => {
    return api.post(`/${entityName}/create`, data, {
      successMessage: `${displayName}이(가) 성공적으로 생성되었습니다.`,
      showSuccessMessage: true
    });
  };

  /**
   * 항목 수정
   * @param {number|string} id - 엔티티 ID
   * @param {Object} data - 수정할 데이터
   * @returns {Promise} 수정된 엔티티 데이터
   */
  const update = (id, data) => {
    return api.put(`/${entityName}/update?id=${id}`, data, {
      successMessage: `${displayName}이(가) 성공적으로 수정되었습니다.`,
      showSuccessMessage: true
    });
  };

  /**
   * 항목 삭제
   * @param {number|string} id - 엔티티 ID
   * @returns {Promise} 삭제 결과
   */
  const remove = (id) => {
    return api.delete(`/${entityName}/delete`, { id }, {
      successMessage: `${displayName}이(가) 성공적으로 삭제되었습니다.`,
      showSuccessMessage: true
    });
  };

  /**
   * 항목 강제 삭제
   * @param {number|string} id - 엔티티 ID
   * @returns {Promise} 삭제 결과
   */
  const forceRemove = (id) => {
    return api.delete(`/${entityName}/force-delete`, { id }, {
      successMessage: `${displayName}이(가) 강제로 삭제되었습니다.`,
      showSuccessMessage: true
    });
  };

  return {
    ...api,
    fetchList,
    fetchById,
    create,
    update,
    remove,
    forceRemove
  };
}