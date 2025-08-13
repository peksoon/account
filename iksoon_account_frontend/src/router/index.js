import { createRouter, createWebHistory } from 'vue-router'
import AccountCalendar from '../components/AccountCalendar.vue'
import AddDataWizard from '../components/AddDataWizard.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: AccountCalendar
  },
  {
    path: '/add-data',
    name: 'AddData',
    component: AddDataWizard,
    props: route => ({
      newAccount: route.query.data ? JSON.parse(route.query.data) : {},
      selectedDate: route.query.date || null
    })
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router