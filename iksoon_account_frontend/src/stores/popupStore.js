import { defineStore } from 'pinia';

export const usePopupStore = defineStore('popup', {
  state: () => ({
    showCustomPopup: false,
    showAddPopup: false,
    eventDetail: null,
    isEditMode: false,
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
    openAddPopup(selectedDate = null) {
      this.showAddPopup = true;
      const date = selectedDate || new Date().toISOString().slice(0, 10);

      this.newAccount = {
        money: '',
        category: '',
        type: 'out',
        user: '',
        keyword: '',
        payment: '',
        memo: '',
        account_number: '',
        date: date,
      };
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
