import { defineStore } from 'pinia';
import axios from 'axios';
import { formatDateToLocalISOString } from '../utils';
import { getApiBaseUrl } from '../config';
import { useCategoryStore } from './categoryStore';
import { usePaymentMethodStore } from './paymentMethodStore';

// 설정 파일에서 API URL 가져오기 (확실한 운영 서버 연결)
const BACKEND_API_BASE_URL = getApiBaseUrl();
console.log('AccountStore API URL:', BACKEND_API_BASE_URL);

export const useAccountStore = defineStore('account', {
  state: () => ({
    monthlyData: [],
    newAccount: {
      money: '',
      category: '',
      type: 'out',
      user: '',
      keyword: '',
      payment: '',
      memo: '',
      account_number: '',
      date: '',
    },
  }),
  actions: {
    async fetchMonthAccounts(year, month) {
      try {
        const formattedMonth = String(month).padStart(2, '0');
        const outRes = await axios.get(
          `${BACKEND_API_BASE_URL}/v2/month-out-account?year=${year}&month=${formattedMonth}`
        );
        const inRes = await axios.get(
          `${BACKEND_API_BASE_URL}/v2/month-in-account?year=${year}&month=${formattedMonth}`
        );
        const outData = outRes.data || [];
        const inData = inRes.data || [];
        this.monthlyData = [
          ...outData.map((item) => ({ ...item, type: 'out' })),
          ...inData.map((item) => ({ ...item, type: 'in' })),
        ];
      } catch (error) {
        console.error('데이터 로드 오류:', error);
      }
    },

    async saveAccount(accountData) {
      try {
        const url = accountData.type === 'out'
          ? `${BACKEND_API_BASE_URL}/v2/out-account/insert`
          : `${BACKEND_API_BASE_URL}/v2/in-account/insert`;
        
        const date = new Date(accountData.date);

        const dataToSend = {
          ...accountData,
          date: formatDateToLocalISOString(date),
        };
        
        await axios.post(url, dataToSend);
        await this.fetchMonthAccounts(date.getFullYear(), date.getMonth() + 1);
      } catch (error) {
        console.error('데이터 저장 오류:', error);
        throw error;
      }
    },

    async updateAccount(eventDetail) {
      try {
        const url = eventDetail.type === 'out'
          ? `${BACKEND_API_BASE_URL}/v2/out-account/update`
          : `${BACKEND_API_BASE_URL}/v2/in-account/update`;
        const date = new Date(eventDetail.date);

        const dataToSend = {
          ...eventDetail,
          date: formatDateToLocalISOString(date),
        };

        // 카테고리 이름이 있는 경우 항상 ID로 변환 (기존 category_id가 있어도 덮어씀)
        if (dataToSend.category_name) {
          // 카테고리 스토어에서 카테고리 ID 찾기
          const categoryStore = useCategoryStore();
          const categories = categoryStore.categories;
          const targetType = eventDetail.type === 'out' ? 'out' : 'in';
          
          if (process.env.NODE_ENV === 'development') {
            console.log(`카테고리 변환 시도: "${dataToSend.category_name}" (타입: ${targetType})`);
            console.log('전체 카테고리 목록:', categories);
            console.log('해당 타입 카테고리:', categories.filter(cat => cat.type === targetType));
          }
          
          const category = categories.find(cat => 
            cat.name === dataToSend.category_name && 
            cat.type === targetType
          );
          if (category) {
            const oldCategoryId = dataToSend.category_id;
            dataToSend.category_id = category.id;
            if (process.env.NODE_ENV === 'development') {
              console.log(`✅ 카테고리 변환 성공: "${dataToSend.category_name}" → ID ${oldCategoryId} → ${category.id}`);
            }
          } else {
            console.warn(`❌ 카테고리를 찾을 수 없음: ${dataToSend.category_name} (타입: ${targetType})`);
            console.warn('사용 가능한 카테고리:', categories.filter(cat => cat.type === targetType));
          }
        }

        // UUID 필드 매핑 (백엔드는 uuid를 요구함)
        if (dataToSend.id && !dataToSend.uuid) {
          dataToSend.uuid = dataToSend.id;
        }

        // 키워드 필드 매핑
        if (dataToSend.keyword && !dataToSend.keyword_name) {
          dataToSend.keyword_name = dataToSend.keyword;
        }

        // 지출/수입별 필드 매핑
        if (eventDetail.type === 'out') {
          // 지출의 경우 payment_method_id 필요
          if (dataToSend.payment_method && !dataToSend.payment_method_id) {
            // 결제수단 이름을 ID로 변환
            const paymentMethodStore = usePaymentMethodStore();
            const paymentMethod = paymentMethodStore.flatPaymentMethods.find(pm => 
              pm.name === dataToSend.payment_method
            );
            if (paymentMethod) {
              dataToSend.payment_method_id = paymentMethod.id;
              if (process.env.NODE_ENV === 'development') {
                console.log(`결제수단 이름 "${dataToSend.payment_method}"을 ID ${paymentMethod.id}로 변환`);
              }
            } else {
              dataToSend.payment_method_id = 1; // 기본값 (첫 번째 결제수단)
              console.warn(`결제수단을 찾을 수 없음: ${dataToSend.payment_method}`);
            }
          }
          if (!dataToSend.payment_method_id) {
            dataToSend.payment_method_id = 1; // 기본값
          }
        } else {
          // 수입의 경우 deposit_path가 필수
          if (!dataToSend.deposit_path && dataToSend.payment_method) {
            dataToSend.deposit_path = dataToSend.payment_method;
          }
          // 기본값 설정 (백엔드에 실제 존재하는 값으로)
          if (!dataToSend.deposit_path) {
            dataToSend.deposit_path = '급여계좌';
          }
        }
        
        // 개발 모드에서만 로그 출력
        if (process.env.NODE_ENV === 'development') {
          console.log('업데이트 요청 URL:', url);
          console.log('업데이트 요청 데이터:', dataToSend);
        }
        
        const response = await axios.put(url, dataToSend);
        
        // 개발 모드에서만 로그 출력
        if (process.env.NODE_ENV === 'development') {
          console.log('업데이트 응답:', response.data);
          console.log('데이터 새로고침 시작...');
        }
        
        await this.fetchMonthAccounts(
          new Date(eventDetail.date).getFullYear(),
          new Date(eventDetail.date).getMonth() + 1
        );
        
        // 개발 모드에서만 로그 출력
        if (process.env.NODE_ENV === 'development') {
          console.log('데이터 새로고침 완료');
        }
      } catch (error) {
        console.error('데이터 업데이트 오류:', error);
        console.error('에러 응답:', error.response?.data);
        console.error('에러 상태:', error.response?.status);
        throw error;
      }
    },

    async deleteAccount(eventDetail) {
      try {
        const url = eventDetail.type === 'out'
          ? `${BACKEND_API_BASE_URL}/v2/out-account/delete?uuid=${eventDetail.uuid}`
          : `${BACKEND_API_BASE_URL}/v2/in-account/delete?uuid=${eventDetail.uuid}`;
        await axios.delete(url);
        await this.fetchMonthAccounts(
          new Date(eventDetail.date).getFullYear(),
          new Date(eventDetail.date).getMonth() + 1
        );
      } catch (error) {
        console.error('데이터 삭제 오류:', error);
      }
    },

    fetchDataForDate(date) {
      return this.monthlyData.filter((data) => {
        const dataDate = data.date.split(' ')[0];
        return dataDate === date;
      });
    },

    async searchByKeyword(keyword, startDate, endDate) {
      try {
        const results = [];
        
        // Search in expense accounts
        const outRes = await axios.get(
          `${BACKEND_API_BASE_URL}/v2/search-keyword-accounts`, {
            params: {
              keyword,
              start_date: startDate,
              end_date: endDate
            }
          }
        );
        
        // Search in income accounts
        const inRes = await axios.get(
          `${BACKEND_API_BASE_URL}/v2/in-search-keyword-accounts`, {
            params: {
              keyword,
              start_date: startDate,
              end_date: endDate
            }
          }
        );

        const outData = outRes.data || [];
        const inData = inRes.data || [];
        
        results.push(
          ...outData.map(item => ({ ...item, type: 'out' })),
          ...inData.map(item => ({ ...item, type: 'in' }))
        );
        
        return results;
      } catch (error) {
        console.error('키워드 검색 오류:', error);
        // Fallback to local search if API fails
        return this.searchByKeywordLocal(keyword, startDate, endDate);
      }
    },

    searchByKeywordLocal(keyword, startDate, endDate) {
      if (!keyword) return [];
      
      const lowerKeyword = keyword.toLowerCase();
      
      return this.monthlyData.filter(item => {
        // Check if item is within date range
        const itemDate = item.date.split(' ')[0];
        if (itemDate < startDate || itemDate > endDate) return false;
        
        // Check if keyword matches (case-insensitive)
        const keywordMatch = item.keyword_name && 
          item.keyword_name.toLowerCase().includes(lowerKeyword);
        const memoMatch = item.memo && 
          item.memo.toLowerCase().includes(lowerKeyword);
        
        return keywordMatch || memoMatch;
      });
    },

    async fetchAccountsInDateRange(startDate, endDate) {
      try {
        const results = [];
        
        // Fetch expense accounts in date range
        const outRes = await axios.get(
          `${BACKEND_API_BASE_URL}/v2/out-accounts`, {
            params: {
              start_date: startDate,
              end_date: endDate
            }
          }
        );
        
        // Fetch income accounts in date range
        const inRes = await axios.get(
          `${BACKEND_API_BASE_URL}/v2/in-accounts`, {
            params: {
              start_date: startDate,
              end_date: endDate
            }
          }
        );

        const outData = outRes.data || [];
        const inData = inRes.data || [];
        
        results.push(
          ...outData.map(item => ({ ...item, type: 'out' })),
          ...inData.map(item => ({ ...item, type: 'in' }))
        );
        
        return results;
      } catch (error) {
        console.error('기간별 계정 조회 오류:', error);
        // Fallback to local search if API fails
        return this.fetchAccountsInDateRangeLocal(startDate, endDate);
      }
    },

    fetchAccountsInDateRangeLocal(startDate, endDate) {
      return this.monthlyData.filter(item => {
        const itemDate = item.date.split(' ')[0];
        return itemDate >= startDate && itemDate <= endDate;
      });
    },
  },
});
