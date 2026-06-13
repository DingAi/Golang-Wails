<template>
  <div class="flex flex-col bg-[#272822] h-screen overflow-hidden">
    <header class="flex justify-between items-center bg-[#3E3D32] px-4 md:px-6 border-[#75715E] border-b h-16">
      <div class="flex items-center gap-4">
        <button @click="foldSidebar"
          class="bg-[#3E3D32] hover:bg-[#4e4d40] px-3 py-1.5 border border-[#75715E] rounded text-[#A6E22E]">
          {{ sidebarFold ? '▶' : '◀' }}
        </button>
        <div class="font-bold text-[#F8F8F2] text-xl">DingVi's Tools</div>
      </div>
      <button @click="logout"
        class="bg-[#3E3D32] hover:bg-[#4e4d40] px-4 py-2 border border-[#75715E] rounded-lg text-[#E6DB74] transition">
        退出登录
      </button>
    </header>

    <div class="flex flex-1 overflow-hidden">
      <aside :class="[
        'border-r border-[#75715E] bg-[#32332E] transition-all duration-300 overflow-hidden shrink-0 pt-5 gap-2',
        sidebarFold ? 'w-18' : 'w-56'
      ]">
        <!-- Logo 区域：折叠时显示小方块，展开时显示长条 -->
        <!-- <div class="pt-5 pb-2" :class="sidebarFold ? 'px-2' : 'px-5'">
          <div class="bg-[#A6E22E] rounded-xl transition-all duration-300"
            :class="sidebarFold ? 'w-10 h-10 mx-auto' : 'w-full h-16'"></div>
        </div> -->

        <nav :class="sidebarFold ? 'px-0' : 'px-5'">
          <router-link to="/dashboard/serial" active-class="bg-[#A6E22E] text-[#272822]" :class="[
            'flex items-center py-3 rounded-lg text-[#F8F8F2] whitespace-nowrap transition',
            sidebarFold ? 'justify-center gap-0 px-0' : 'gap-3 px-4'
          ]">
            <CpuChipIcon class="flex-shrink-0 w-5 h-5" />
            <span v-show="!sidebarFold">串口助手</span>
          </router-link>
          <router-link to="/dashboard/frame-visual-rule" active-class="bg-[#A6E22E] text-[#272822]" :class="[
            'flex items-center py-3 rounded-lg text-[#F8F8F2] whitespace-nowrap transition',
            sidebarFold ? 'justify-center gap-0 px-0' : 'gap-3 px-4'
          ]">
            <InboxArrowDownIcon class="w-5 h-5 shrink-0" />
            <span v-show="!sidebarFold">数据帧规则设计</span>
          </router-link>
        </nav>
      </aside>

      <main class="flex-1 bg-[#272822] p-4 md:p-6 overflow-y-auto">
        <router-view />
      </main>
    </div>
    <!-- 新增全局VSCode风格底部状态栏，绑定到 store 实时数据 -->
    <div class="flex justify-between items-center bg-[#3E3D32] px-4 py-2 border-[#75715E] border-t text-xs md:text-sm">
      <div class="flex gap-5 text-[#F8F8F2]">
        <span>串口：{{ store.statusText }}</span>
        <span>端口：{{ store.portName }}</span>
        <span>波特率：{{ store.baudRate }}</span>
      </div>
      <div class="flex gap-5 text-[#E6DB74]">
        <span>RX：{{ formatBytes(store.rxBytes) }}</span>
        <span>TX：{{ formatBytes(store.txBytes) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { CpuChipIcon, InboxArrowDownIcon } from '@heroicons/vue/24/outline'
import { useSerialStore } from '../stores/serial.js'

const router = useRouter()
const store = useSerialStore()

const sidebarFold = ref(false)

const foldSidebar = () => {
  sidebarFold.value = !sidebarFold.value
}

const logout = () => {
  localStorage.removeItem('isLoggedIn')
  router.push('/login')
}

function formatBytes(bytes) {
  if (bytes === 0) return '0 Byte'
  const units = ['Byte', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  if (i === 0) return bytes + ' ' + units[i]
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}
</script>