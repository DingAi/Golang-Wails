import { createRouter, createWebHistory } from 'vue-router'
import BlankLayout from '../layouts/BlankLayout.vue'
import Dashboard from '../layouts/Dashboard.vue'
import Login from '../views/Login.vue'
import SerialPort from '../views/SerialPort.vue'
import ModbusParser from '../views/FrameParser/ModbusParser.vue'
import IEC104Parser from '../views/FrameParser/IEC104Parser.vue'
import DLT698Parser from '../views/FrameParser/DLT698Parser.vue'
import DLT645Parser from '../views/FrameParser/DLT645Parser.vue'

const routes = [
  {
    path: '/',
    redirect: (to) => {
      const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
      return isLoggedIn ? '/dashboard' : '/login'
    }
  },
  // ✅ 添加 login 路由
  {
    path: '/login',
    component: BlankLayout,
    children: [{ path: '', component: Login }]
  },
  {
    path: '/dashboard',
    component: Dashboard,
    redirect: '/dashboard/serial',
    children: [
      { path: 'serial', component: SerialPort },
      { path: 'modbus-parser', component: ModbusParser },
      { path: 'iec104-parser', component: IEC104Parser },
      { path: 'dlt698-parser', component: DLT698Parser },
      { path: 'dlt645-parser', component: DLT645Parser }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
  if (to.path !== '/login' && !isLoggedIn) {
    next('/login')
  } else {
    next()
  }
})

export default router