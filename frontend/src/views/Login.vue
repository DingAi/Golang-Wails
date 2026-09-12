<template>
  <!-- 页面背景 mk 主底色 -->
  <div class="flex justify-center items-center bg-mk-bg w-full h-full">
    <!-- 登录卡片 -->
    <div class="bg-mk-surface shadow-2xl p-10 border border-mk-surface rounded-2xl w-full max-w-md">
      <!-- Logo头部 -->
      <div class="mb-10 text-center">
        <div class="inline-flex justify-center items-center bg-mk-surface shadow-lg mb-5 border border-mk-comment rounded-2xl w-16 h-16">
          <svg class="w-8 h-8 text-mk-green" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 3h14a2 2 0 012 2v14a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 11h4M12 15h4M8 11h.01M8 15h.01" />
          </svg>
        </div>
        <h1 class="font-bold text-mk-foreground text-3xl">IoT 协议网关</h1>
        <p class="mt-2 font-medium text-mk-cyan">设备管理平台</p>
      </div>

      <!-- 登录表单 -->
      <form @submit.prevent="handleLogin" class="space-y-6">
        <div>
          <label class="block mb-2 font-medium text-mk-green text-sm">用户名</label>
          <input
            v-model="username"
            type="text"
            class="bg-mk-surface px-5 py-3 border border-mk-comment focus:border-mk-cyan rounded-xl focus:outline-none w-full text-mk-yellow placeholder:text-mk-comment"
            placeholder="admin"
            autocomplete="off"
          />
        </div>
        <div>
          <label class="block mb-2 font-medium text-mk-green text-sm">密码</label>
          <input
            v-model="password"
            type="password"
            class="bg-mk-surface px-5 py-3 border border-mk-comment focus:border-mk-cyan rounded-xl focus:outline-none w-full text-mk-yellow placeholder:text-mk-comment"
            placeholder="••••••"
          />
        </div>
        <!-- 登录按钮 mk玫红 -->
        <button
          type="submit"
          class="bg-mk-pink hover:bg-mk-pinkHover shadow-lg py-3 rounded-xl w-full font-semibold text-white transition duration-200"
        >
          登 录
        </button>
      </form>

      <!-- 错误提示 -->
      <div v-if="errorMessage" class="bg-mk-surface mt-5 py-3 border border-mk-pink/60 rounded-xl text-mk-pink text-sm text-center">
        {{ errorMessage }}
      </div>

      <!-- 底部提示 -->
      <div class="mt-8 text-mk-comment text-sm text-center">
        演示默认：admin / 123456
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const username = ref('admin')
const password = ref('123456')
const errorMessage = ref('')

const handleLogin = () => {
  if (username.value === 'admin' && password.value === '123456') {
    localStorage.setItem('isLoggedIn', 'true')
    router.push('/dashboard/serial')
  } else {
    errorMessage.value = '用户名或密码错误'
    setTimeout(() => { errorMessage.value = '' }, 3000)
  }
}
</script>