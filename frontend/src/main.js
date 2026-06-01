// src/main.ts
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'  // 你稍后需要创建 router 实例
import '../public/style.css'  // 引入 Tailwind CSS

// 创建应用
const app = createApp(App)

// 使用 Pinia
const pinia = createPinia()
app.use(pinia)

// 使用 Vue Router
app.use(router)

// 挂载到 #app（Wails 默认容器）
app.mount('#app')