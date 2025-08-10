import { defineStore } from 'pinia';
import axios from 'axios';
import { formatDateToLocalISOString } from '../utils';

// 환경변수에서 API 주소를 가져오고, 기본값은 운영 서버 IP
const BACKEND_API_BASE_URL = process.env.VUE_APP_BACKEND_API_BASE_URL || 'http://133.186.153.179:8080';

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
  },
});
