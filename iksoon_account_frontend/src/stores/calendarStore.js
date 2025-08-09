import { defineStore } from 'pinia';

export const useCalendarStore = defineStore('calendar', {
  state: () => ({
    currentYear: new Date().getFullYear(),
    currentMonth: String(new Date().getMonth() + 1).padStart(2, '0'),
  }),
  actions: {
    // FullCalendar와 호환성을 위한 더미 메서드들
    createCalendarEvents() {
      // FullCalendar에서는 더 이상 사용하지 않음
      // AccountCalendar 컴포넌트에서 직접 처리
    },
    async changeMonth() {
      // FullCalendar에서는 더 이상 사용하지 않음
      // AccountCalendar 컴포넌트에서 직접 처리
    },
    setCurrentDate(year, month) {
      this.currentYear = year;
      this.currentMonth = String(month).padStart(2, '0');
    },
  },
});
