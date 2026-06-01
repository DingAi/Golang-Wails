<template>
  <!-- 关键：外层 div 必须占满窗口，去掉多余的 p-4 -->
  <div class="flex justify-center items-center bg-gradient-to-br from-teal-900 via-emerald-800 to-slate-900 w-full h-full">
    <!-- 去掉 p-4，让容器可以贴边，内部用 padding 控制内边距 -->
    <div class="bg-slate-900/70 shadow-2xl backdrop-blur-xl p-10 border border-emerald-500/20 rounded-2xl w-full max-w-md">
      <!-- Logo / 标题 -->
      <div class="mb-10 text-center">
        <div class="inline-flex justify-center items-center bg-gradient-to-br from-emerald-600 to-emerald-700 shadow-lg mb-5 rounded-2xl w-18 h-18">
          <svg class="w-9 h-9 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 3h14a2 2 0 012 2v14a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 11h4M12 15h4M8 11h.01M8 15h.01" />
          </svg>
        </div>
        <h1 class="font-bold text-white text-3xl">IoT 协议网关</h1>
        <p class="mt-2 text-emerald-300">设备管理平台</p>
      </div>

      <!-- 登录表单 -->
      <form @submit.prevent="handleLogin" class="space-y-6">
        <div>
          <label class="block mb-2 font-medium text-gray-200 text-sm">用户名</label>
          <input
            v-model="username"
            type="text"
            class="bg-slate-800/60 px-5 py-3 border border-emerald-500/30 focus:border-emerald-400 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-500/50 w-full text-white placeholder-gray-400"
            placeholder="admin"
            autocomplete="off"
          />
        </div>
        <div>
          <label class="block mb-2 font-medium text-gray-200 text-sm">密码</label>
          <input
            v-model="password"
            type="password"
            class="bg-slate-800/60 px-5 py-3 border border-emerald-500/30 focus:border-emerald-400 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-500/50 w-full text-white placeholder-gray-400"
            placeholder="••••••"
          />
        </div>
        <button
          type="submit"
          class="bg-gradient-to-r from-emerald-600 hover:from-emerald-700 to-emerald-700 hover:to-emerald-800 shadow-lg py-3 rounded-xl w-full font-semibold text-white transition duration-200"
        >
          登 录
        </button>
      </form>

      <!-- 错误提示 -->
      <div v-if="errorMessage" class="bg-red-900/30 mt-5 py-3 border border-red-500/20 rounded-xl text-red-400 text-sm text-center">
        {{ errorMessage }}
      </div>

      <!-- 提示 -->
      <div class="mt-8 text-gray-400 text-sm text-center">
        演示默认：admin / 123456
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const username = ref('')
const password = ref('')
const errorMessage = ref('')

const handleLogin = () => {
  if (username.value === 'admin' && password.value === '123456') {
    localStorage.setItem('isLoggedIn', 'true')
    router.push('/dashboard')
  } else {
    errorMessage.value = '用户名或密码错误'
    setTimeout(() => { errorMessage.value = '' }, 3000)
  }
}
</script>