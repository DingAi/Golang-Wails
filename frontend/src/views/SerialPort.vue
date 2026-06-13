<template>
  <div class="flex flex-col w-full h-full overflow-hidden">
    <!-- 主双栏容器 -->
    <div class="flex md:flex-row flex-col flex-1 overflow-hidden">
      <!-- 左侧配置面板：桌面默认展开 -->
      <div
        class="flex-shrink-0 p-4 border-[#3E3D32] md:border-r border-b md:border-b-0 w-full md:w-72 lg:w-80 overflow-y-auto transition-all duration-200 custom-scrollbar"
        :class="{ 'md:hidden': sidebarCollapsed }">
        <!-- 串口参数 -->
        <div class="bg-[#3E3D32] mb-4 p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-[#66D9EF] text-sm">串口参数</h3>
          <div class="space-y-3">
            <div class="flex items-center gap-2">
              <span class="w-16 text-[#A6E22E] text-sm shrink-0">端口名</span>
              <select v-model="store.portName"
                class="flex-1 bg-[#32332E] px-3 py-2 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none text-[#F8F8F2] text-sm">
                <option v-for="p in store.ports" :key="p.name" :value="p.name">{{ p.name }}</option>
                <option v-if="store.ports.length === 0" value="COM5">COM5</option>
              </select>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-[#A6E22E] text-sm shrink-0">波特率</span>
              <select v-model="store.baudRate"
                class="flex-1 bg-[#32332E] px-3 py-2 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none text-[#F8F8F2] text-sm">
                <option :value="1200">1200</option>
                <option :value="2400">2400</option>
                <option :value="4800">4800</option>
                <option :value="9600">9600</option>
                <option :value="19200">19200</option>
                <option :value="38400">38400</option>
                <option :value="57600">57600</option>
                <option :value="115200">115200</option>
                <option :value="230400">230400</option>
                <option :value="460800">460800</option>
                <option :value="921600">921600</option>
              </select>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-[#A6E22E] text-sm shrink-0">数据位</span>
              <select v-model="store.dataBits"
                class="flex-1 bg-[#32332E] px-3 py-2 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none text-[#F8F8F2] text-sm">
                <option :value="5">5</option>
                <option :value="6">6</option>
                <option :value="7">7</option>
                <option :value="8">8</option>
              </select>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-[#A6E22E] text-sm shrink-0">校验位</span>
              <select v-model="store.parity"
                class="flex-1 bg-[#32332E] px-3 py-2 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none text-[#F8F8F2] text-sm">
                <option>无</option>
                <option>奇</option>
                <option>偶</option>
                <option>Mark</option>
                <option>Space</option>
              </select>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-[#A6E22E] text-sm shrink-0">停止位</span>
              <select v-model="store.stopBits"
                class="flex-1 bg-[#32332E] px-3 py-2 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none text-[#F8F8F2] text-sm">
                <option :value="1">1</option>
                <option :value="2">2</option>
              </select>
            </div>
          </div>
        </div>

        <!-- 串口操作按钮 -->
        <div class="bg-[#3E3D32] mb-4 p-4 rounded-xl">
          <div class="flex gap-3">
            <button @click="handleOpen" :disabled="store.isOpen"
              class="flex-1 bg-[#A6E22E] hover:bg-[#93c725] disabled:opacity-40 py-2 rounded-lg font-medium text-[#272822] text-sm transition">
              打开串口
            </button>
            <button @click="handleClose" :disabled="!store.isOpen"
              class="flex-1 bg-[#F92672] hover:bg-[#dd1e63] disabled:opacity-40 py-2 rounded-lg font-medium text-white text-sm transition">
              关闭串口
            </button>
          </div>
          <div v-if="errorMsg" class="mt-2 text-[#F92672] text-xs text-center">
            {{ errorMsg }}
          </div>
        </div>

        <!-- 接收设置 -->
        <div class="bg-[#3E3D32] mb-4 p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-[#66D9EF] text-sm">接收设置</h3>
          <div class="space-y-2">
            <label class="flex items-center gap-2 text-[#F8F8F2] text-sm">
              <input type="checkbox" v-model="store.autoSplit" @change="store.updateAutoSplit()"
                class="accent-[#F92672]" />
              自动断帧
              <input
                class="bg-[#32332E] px-2 py-0.5 border border-[#75715E] rounded w-14 text-[#E6DB74] text-sm text-center"
                v-model="store.splitMs" @change="store.updateAutoSplit()" />
              <span class="text-[#75715E]">ms</span>
            </label>
            <label class="flex items-center gap-2 text-[#F8F8F2] text-sm">
              <input type="checkbox" v-model="store.saveToFile" class="accent-[#F92672]" />
              将接收保存到文件
            </label>
          </div>
          <div class="flex gap-3 mt-3">
            <button
              class="flex-1 bg-[#32332E] hover:bg-[#4e4d40] py-2 border border-[#75715E] rounded-lg text-[#E6DB74] text-sm">
              保存数据
            </button>
            <button @click="handleClear"
              class="flex-1 bg-[#32332E] hover:bg-[#4e4d40] py-2 border border-[#75715E] rounded-lg text-[#E6DB74] text-sm">
              清空数据
            </button>
          </div>
        </div>

        <!-- 发送设置 -->
        <div class="bg-[#3E3D32] p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-[#66D9EF] text-sm">发送设置</h3>
          <label class="flex items-center gap-2 text-[#F8F8F2] text-sm">
            <input type="checkbox" v-model="store.autoSend" @change="handleAutoSendToggle" class="accent-[#F92672]" />
            定时发送
            <input
              class="bg-[#32332E] px-2 py-0.5 border border-[#75715E] rounded w-14 text-[#E6DB74] text-sm text-center"
              v-model.number="store.sendSec" />
            <span class="text-[#75715E]">秒</span>
          </label>
        </div>
      </div>

      <!-- 右侧收发区域 -->
      <div class="flex flex-col flex-1 ml-3 min-w-0 h-full">
        <!-- 顶部栏 -->
        <div class="flex justify-between items-center px-4 py-2">
          <button @click="sidebarCollapsed = !sidebarCollapsed"
            class="md:hidden bg-[#3E3D32] hover:bg-[#4e4d40] px-2 py-1 border border-[#75715E] rounded text-[#A6E22E] text-xs">
            {{ sidebarCollapsed ? '展开设置' : '收起设置' }}
          </button>
          <span class="ml-auto text-[#75715E] text-sm">数据帧记录</span>
        </div>

        <!-- 消息列表容器（带边框圆角，和截图一致） -->
        <div class="flex flex-col flex-1 border border-[#3E3D32] rounded-xl min-h-0">
          <!-- 消息滚动区 -->
          <div ref="msgListRef" class="flex-1 space-y-3 px-4 py-2 min-h-0 overflow-y-auto custom-scrollbar">
            <div v-for="(msg, idx) in store.messages" :key="idx" class="flex"
              :class="msg.type === 'sent' ? 'justify-end' : 'justify-start'">
              <!-- 接收气泡 -->
              <div v-if="msg.type === 'recv'"
                class="bg-[#3E3D32] shadow-sm px-4 py-3 border border-[#75715E]/40 rounded-2xl rounded-tl-sm max-w-[85%]">
                <div class="font-mono text-[#E6DB74] text-sm break-all whitespace-pre-wrap">
                  {{ msg.text }}
                </div>
                <div class="mt-1 text-[#75715E] text-[10px] text-right">
                  {{ msg.time }}
                </div>
              </div>
              <!-- 发送气泡 -->
              <div v-if="msg.type === 'sent'"
                class="bg-[#A6E22E]/15 shadow-sm px-4 py-3 border border-[#A6E22E]/30 rounded-2xl rounded-tr-sm max-w-[85%]">
                <div class="font-mono text-[#F8F8F2] text-sm break-all whitespace-pre-wrap">
                  {{ msg.text }}
                </div>
                <div class="mt-1 text-[#A6E22E]/60 text-[10px] text-right">
                  {{ msg.time }}
                </div>
              </div>
            </div>
            <!-- 空状态 -->
            <div v-if="store.messages.length === 0"
              class="flex justify-center items-center h-full text-[#75715E] text-sm">
              暂无数据，请打开串口后收发数据
            </div>
          </div>

          <!-- 接收格式栏 -->
          <div class="flex flex-wrap items-center gap-3 bg-[#3E3D32] px-4 py-1 border-[#75715E] border-t rounded-b-xl">
            <span class="text-[#75715E] text-sm">显示格式:</span>
            <span v-for="item in recvModeList" :key="item.key" @click="store.setRecvMode(item.key)"
              class="px-2 py-1 rounded text-sm transition cursor-pointer" :class="store.recvMode === item.key
                  ? 'bg-[#A6E22E] text-[#272822]'
                  : 'text-[#F8F8F2] hover:bg-[#4e4d40]'
                ">
              {{ item.name }}
            </span>
            <span class="ml-auto text-[#75715E] text-xs">
              {{ store.messages.length }} 帧
            </span>
          </div>
        </div>

        <!-- 分隔间距 -->
        <div class="h-3 shrink-0"></div>

        <!-- 发送区域 -->
        <div class="bg-[#272822] border border-[#3E3D32] rounded-xl overflow-hidden shrink-0">
          <!-- 输入框 + 发送按钮 -->
          <div class="flex gap-3 p-2">
            <textarea v-model="sendData"
              class="flex-1 bg-[#32332E] px-4 py-1 border border-[#75715E] focus:border-[#66D9EF] rounded-xl focus:outline-none font-mono text-[#E6DB74] text-sm resize-none custom-scrollbar"
              :rows="sendRows" placeholder="输入要发送的数据..." @keydown.enter.ctrl="handleSend"
              @keydown.enter.exact="handleSend"></textarea>
            <button @click="handleSend" :disabled="!store.isOpen || !sendData.trim()"
              class="self-stretch bg-[#F92672] hover:bg-[#dd1e63] disabled:opacity-40 px-6 rounded-xl font-medium text-white text-sm transition shrink-0">
              发送
            </button>
          </div>
          <!-- 发送格式 移至底部 -->
          <div class="flex flex-wrap items-center gap-3 bg-[#3E3D32] px-4 py-1 border-[#75715E] border-t">
            <span class="text-[#75715E] text-sm">发送格式:</span>
            <span v-for="item in sendModeList" :key="item.key" @click="store.setSendMode(item.key)"
              class="px-2 py-1 rounded text-sm transition cursor-pointer" :class="store.sendMode === item.key
                  ? 'bg-[#A6E22E] text-[#272822]'
                  : 'text-[#F8F8F2] hover:bg-[#4e4d40]'
                ">
              {{ item.name }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useSerialStore } from '../stores/serial.js'

const store = useSerialStore()

const recvModeList = [
  { key: 'ascii', name: 'ASCII' },
  { key: 'hex', name: 'HEX' },
  { key: 'bin', name: 'BIN' },
  { key: 'dec', name: 'DEC' },
  { key: 'bcd', name: 'BCD' },
]
const sendModeList = [
  { key: 'ascii', name: 'ASCII' },
  { key: 'hex', name: 'HEX' },
  { key: 'bin', name: 'BIN' },
  { key: 'dec', name: 'DEC' },
  { key: 'bcd', name: 'BCD' },
]

const sendData = ref('')
const errorMsg = ref('')
const msgListRef = ref(null)
// 默认展开，和截图效果一致
const sidebarCollapsed = ref(false)
const sendRows = ref(2)

onMounted(() => {
  store.initEventListeners()
  store.fetchPorts()
  store.updateAutoSplit()
  updateSendRows()
  window.addEventListener('resize', updateSendRows)
})

onUnmounted(() => {
  store.destroyEventListeners()
  window.removeEventListener('resize', updateSendRows)
})

// 消息自动滚动到底部
watch(
  () => store.messages.length,
  async () => {
    await nextTick()
    if (msgListRef.value) {
      msgListRef.scrollTop = msgListRef.value.scrollHeight
    }
  }
)

// 动态修改输入框行数
function updateSendRows() {
  const h = window.innerHeight
  if (h < 600) {
    sendRows.value = 1
  } else if (h < 800) {
    sendRows.value = 2
  } else {
    sendRows.value = 3
  }
}

// 事件方法
async function handleOpen() {
  errorMsg.value = ''
  try {
    await store.open()
    if (window.innerWidth < 768) sidebarCollapsed.value = true
  } catch (err) {
    errorMsg.value = err?.message || String(err)
  }
}

async function handleClose() {
  errorMsg.value = ''
  try {
    await store.close()
    if (store.autoSend) {
      store.autoSend = false
      await store.stopAutoSend()
    }
  } catch (err) {
    errorMsg.value = err?.message || String(err)
  }
}

async function handleSend() {
  if (!store.isOpen || !sendData.value.trim()) return
  errorMsg.value = ''
  try {
    await store.send(sendData.value)
    sendData.value = ''
  } catch (err) {
    errorMsg.value = err?.message || String(err)
  }
}

async function handleAutoSendToggle() {
  errorMsg.value = ''
  try {
    if (store.autoSend) {
      await store.startAutoSend(sendData.value)
    } else {
      await store.stopAutoSend()
    }
  } catch (err) {
    errorMsg.value = err?.message || String(err)
    store.autoSend = !store.autoSend
  }
}

async function handleClear() {
  errorMsg.value = ''
  try {
    await store.clearBuffer()
  } catch (err) {
    errorMsg.value = err?.message || String(err)
  }
}
</script>

<style scoped>
/* 自定义滚动条 Monokai 风格 */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #75715E;
  border-radius: 3px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #8a876e;
}

.custom-scrollbar {
  scrollbar-width: thin;
  scrollbar-color: #75715E transparent;
}
</style>