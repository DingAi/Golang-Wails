<template>
  <div class="flex flex-col bg-mk-bg h-screen overflow-hidden">
    <header class="flex justify-between items-center bg-mk-surface px-4 md:px-6 border-mk-comment border-b h-16">
      <div class="flex items-center gap-4">
        <button @click="foldSidebar"
          class="bg-mk-surface hover:bg-mk-inputHover px-3 py-1.5 border border-mk-comment rounded text-mk-green">
          {{ sidebarFold ? '▶' : '◀' }}
        </button>
        <div class="font-bold text-mk-foreground text-xl">DingVi's Tools</div>
      </div>
      <button @click="logout"
        class="bg-mk-surface hover:bg-mk-inputHover px-4 py-2 border border-mk-comment rounded-lg text-mk-yellow transition">
        退出登录
      </button>
    </header>

    <div class="flex flex-1 overflow-hidden">
      <aside :class="[
        'border-r border-mk-comment bg-mk-surface transition-all duration-300 overflow-hidden shrink-0 pt-5 gap-2',
        sidebarFold ? 'w-18' : 'w-56'
      ]">
        <!-- Logo 区域：折叠时显示小方块，展开时显示长条 -->
        <!-- <div class="pt-5 pb-2" :class="sidebarFold ? 'px-2' : 'px-5'">
          <div class="bg-mk-green rounded-xl transition-all duration-300"
            :class="sidebarFold ? 'w-10 h-10 mx-auto' : 'w-full h-16'"></div>
        </div> -->

        <nav :class="sidebarFold ? 'px-0' : 'px-5'">
          <!-- 串口助手（一级菜单） -->
          <router-link to="/dashboard/serial" active-class="bg-mk-green text-[#272822]" :class="[
            'flex items-center py-3 rounded-lg text-mk-foreground whitespace-nowrap transition',
            sidebarFold ? 'justify-center gap-0 px-0' : 'gap-3 px-4'
          ]">
            <CpuChipIcon class="w-5 h-5 shrink-0" />
            <span v-show="!sidebarFold" class="font-bold">串口助手</span>
          </router-link>

          <!-- MQTT（一级菜单） -->
          <router-link to="/dashboard/mqtt" active-class="bg-mk-green text-[#272822]" :class="[
            'flex items-center py-3 rounded-lg text-mk-foreground whitespace-nowrap transition',
            sidebarFold ? 'justify-center gap-0 px-0' : 'gap-3 px-4'
          ]">
            <CpuChipIcon class="w-5 h-5 shrink-0" />
            <span v-show="!sidebarFold" class="font-bold">MQTT</span>
          </router-link>

          <!-- 协议解析（二级菜单父级） -->
          <div v-if="!sidebarFold">
            <!-- 父级按钮：点击展开/收起子菜单 -->
            <div @click="toggleProtocolMenu"
              class="flex justify-between items-center hover:bg-mk-surface px-4 py-3 rounded-lg text-mk-foreground transition cursor-pointer">
              <div class="flex items-center gap-3">
                <RectangleGroupIcon class="w-5 h-5 shrink-0" />
                <span>协议解析</span>
              </div>
              <ChevronDownIcon v-if="protocolMenuOpen" class="w-4 h-4" />
              <ChevronRightIcon v-else class="w-4 h-4" />
            </div>
            <!-- 子菜单列表 -->
            <div v-show="protocolMenuOpen" class="flex flex-col gap-1 mt-1 ml-6">
              <router-link to="/dashboard/modbus-parser" active-class="bg-mk-green text-[#272822]"
                class="flex items-center gap-3 hover:bg-mk-surface px-4 py-2 rounded-lg text-mk-foreground transition">
                <span class="font-bold">Modbus规则设计</span>
              </router-link>
              <router-link to="/dashboard/iec104-parser" active-class="bg-mk-green text-[#272822]"
                class="flex items-center gap-3 hover:bg-mk-surface px-4 py-2 rounded-lg text-mk-foreground transition">
                <span class="font-bold">IEC104解析器</span>
              </router-link>
              <router-link to="/dashboard/dlt698-parser" active-class="bg-mk-green text-[#272822]"
                class="flex items-center gap-3 hover:bg-mk-surface px-4 py-2 rounded-lg text-mk-foreground transition">
                <span class="font-bold">DLT698解析器</span>
              </router-link>
              <router-link to="/dashboard/dlt645-parser" active-class="bg-mk-green text-[#272822]"
                class="flex items-center gap-3 hover:bg-mk-surface px-4 py-2 rounded-lg text-mk-foreground transition">
                <span class="font-bold">DLT645解析器</span>
              </router-link>
            </div>
            <router-link to="/dashboard/card" active-class="bg-mk-green text-[#272822]" :class="[
              'flex items-center py-3 rounded-lg text-mk-foreground whitespace-nowrap transition',
              sidebarFold ? 'justify-center gap-0 px-0' : 'gap-3 px-4'
            ]">
              <CpuChipIcon class="w-5 h-5 shrink-0" />
              <span v-show="!sidebarFold" class="font-bold">Card</span>
            </router-link>
          </div>



          <!-- 侧边栏折叠时：只显示一个可点击的父级图标（不带文字和子菜单） -->
          <div v-else>
            <div @click="toggleProtocolMenu"
              class="flex justify-center hover:bg-mk-surface py-3 rounded-lg text-mk-foreground transition cursor-pointer">
              <RectangleGroupIcon class="w-5 h-5 shrink-0" />
            </div>
            <!-- 折叠时如果要显示子菜单悬浮窗会比较复杂，此处省略，保持简洁 -->
          </div>
        </nav>
      </aside>

      <main class="flex-1 bg-mk-bg p-4 md:p-6 overflow-y-auto">
        <router-view />
      </main>
    </div>
    <!-- 新增全局VSCode风格底部状态栏，绑定到 store 实时数据 -->
    <div
      class="flex justify-between items-center bg-mk-surface px-4 py-2 border-mk-comment border-t text-xs md:text-sm">
      <div class="flex gap-5 text-mk-foreground">
        <span>串口：{{ store.statusText }}</span>
        <span>端口：{{ store.portName }}</span>
        <span>波特率：{{ store.baudRate }}</span>
      </div>
      <div class="flex gap-5 text-mk-yellow">
        <span>RX：{{ formatBytes(store.rxBytes) }}</span>
        <span>TX：{{ formatBytes(store.txBytes) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  CpuChipIcon,
  InboxArrowDownIcon,
  RectangleGroupIcon,
  ChevronDownIcon,
  ChevronRightIcon,
} from '@heroicons/vue/24/outline'
import { useSerialStore } from '../stores/serial.js'

const router = useRouter()
const store = useSerialStore()

const sidebarFold = ref(false)

const foldSidebar = () => {
  sidebarFold.value = !sidebarFold.value
}

// 控制协议解析二级菜单的展开/收起
const protocolMenuOpen = ref(true)  // 默认展开
const toggleProtocolMenu = () => {
  if (!sidebarFold.value) {
    protocolMenuOpen.value = !protocolMenuOpen.value
  }
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