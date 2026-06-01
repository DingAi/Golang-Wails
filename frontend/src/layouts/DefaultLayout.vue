<template>
  <div class="flex flex-col h-screen">
    <!-- 顶部导航栏 -->
    <header class="flex justify-between items-center bg-gray-800 shadow-md px-6 h-14 text-white">
      <div class="font-bold text-xl">IoT 网关平台</div>
      <div class="flex gap-4">
        <router-link to="/dashboard" class="hover:text-gray-300">仪表盘</router-link>
        <router-link to="/settings" class="hover:text-gray-300">设置</router-link>
        <button @click="logout" class="bg-red-600 hover:bg-red-700 px-3 py-1 rounded">退出</button>
      </div>
    </header>

    <!-- 主体：侧边栏 + 内容区，flex-1 自动占满剩余高度 -->
    <div class="flex flex-1 overflow-hidden">
      <!-- 侧边栏 -->
      <aside class="bg-gray-100 p-4 border-gray-200 border-r w-64 overflow-y-auto">
        <nav class="space-y-2">
          <router-link to="/dashboard/overview" class="block hover:bg-gray-200 p-2 rounded">
            📊 概览
          </router-link>
          <router-link to="/dashboard/devices" class="block hover:bg-gray-200 p-2 rounded">
            🔌 设备管理
          </router-link>
          <router-link to="/dashboard/serial" class="block hover:bg-gray-200 p-2 rounded">
            📡 串口助手
          </router-link>
        </nav>
      </aside>

      <!-- 主要内容区域，独立滚动 -->
      <main class="flex-1 bg-white p-6 overflow-y-auto">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()
const logout = () => {
  // 清除登录状态，跳转到登录页
  router.push('/login')
}
</script>