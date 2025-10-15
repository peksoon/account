import { createRouter, createWebHistory } from 'vue-router'
import AccountCalendar from '../components/AccountCalendar.vue'
import AddDataWizard from '../components/AddDataWizard.vue'
import KeywordSearchPage from '../components/KeywordSearchPage.vue'
import ExportDataPage from '../components/ExportDataPage.vue'
import DetailPage from '../components/DetailPage.vue'
import ManagementPage from '../components/ManagementPage.vue'

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
  },
  {
    path: '/detail',
    name: 'Detail',
    component: DetailPage,
    props: route => ({
      eventDetail: route.query.data ? JSON.parse(route.query.data) : {}
    })
  },
  {
    path: '/keyword-search',
    name: 'KeywordSearch',
    component: KeywordSearchPage
  },
  {
    path: '/export-data',
    name: 'ExportData',
    component: ExportDataPage
  },
  {
    path: '/management',
    name: 'Management',
    component: ManagementPage
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router