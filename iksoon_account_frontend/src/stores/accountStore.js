import { defineStore } from 'pinia';
import axios from 'axios';
import { formatDateToLocalISOString } from '../utils';
import { getApiBaseUrl } from '../config';

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
        await axios.put(url, dataToSend);
        await this.fetchMonthAccounts(
          new Date(eventDetail.date).getFullYear(),
          new Date(eventDetail.date).getMonth() + 1
        );
      } catch (error) {
        console.error('데이터 업데이트 오류:', error);
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
