import { defineStore } from 'pinia';
import { getTodayKST } from '../utils';

export const usePopupStore = defineStore('popup', {
  state: () => ({
    showCustomPopup: false,
    showAddPopup: false,
    eventDetail: null,
    isEditMode: false,
    newAccount: {
      money: '',
      category_id: null,
      type: 'out',
      user: '',
      keyword_name: '',
      payment_method_id: null,
      memo: '',
      deposit_path: '',
      date: '',
    },
  }),
  actions: {
    openAddPopup(selectedDate = null) {
      const date = selectedDate || getTodayKST();

      // 새로운 계정 데이터 설정 (AddPopup에서 기대하는 필드명으로)
      this.newAccount = {
        money: '',
        category_id: null,
        type: 'out',
        user: '',
        keyword_name: '',
        payment_method_id: null,
        memo: '',
        deposit_path: '',
        date: date,
      };
      
      // 팝업 표시
      this.showAddPopup = true;
    },
    closeAddPopup() {
      this.showAddPopup = false;
    },
    closePopup() {
      this.showCustomPopup = false;
      this.eventDetail = {};
      this.isEditMode = false;
    },
    showDetailPopup(data) {
      this.showCustomPopup = true;
      this.eventDetail = { ...data };
      this.isEditMode = false;
    },
    openEditMode() {
      this.isEditMode = true;
    },
  },
});
