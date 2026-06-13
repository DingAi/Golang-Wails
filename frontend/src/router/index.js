import { createRouter, createWebHistory } from 'vue-router'
import BlankLayout from '../layouts/BlankLayout.vue'
import Dashboard from '../layouts/Dashboard.vue'
import Login from '../views/Login.vue'
import SerialPort from '../views/SerialPort.vue'
import Receiver from '../views/Receiver.vue'
import FrameVisualRule from '../views/FrameVisualRule.vue'

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
      { path: 'receiver', component: Receiver },
      { path: 'frame-visual-rule', component: FrameVisualRule }
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