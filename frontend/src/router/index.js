import { createRouter, createWebHistory } from 'vue-router'
import BlankLayout from '../layouts/BlankLayout.vue'
import DefaultLayout from '../layouts/DefaultLayout.vue'
import Login from '../views/Login.vue'
// import Dashboard from '../views/Dashboard.vue'
// import SerialPort from '../views/SerialPort.vue'  // 串口助手页面，稍后实现

const routes = [
  {
    path: '/login',
    component: BlankLayout,
    children: [
      { path: '', component: Login }
    ]
  },
  {
    path: '/',
    component: DefaultLayout,
    children: [
    //   { path: 'dashboard', component: Dashboard },
    //   { path: 'serial', component: SerialPort },
    //   { path: '', redirect: '/dashboard' }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局路由守卫：未登录时跳转登录页
router.beforeEach((to, from, next) => {
  const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true'
  if (to.path !== '/login' && !isLoggedIn) {
    next('/login')
  } else {
    next()
  }
})

export default router