<template>
  <div class="flex flex-col w-full h-full overflow-hidden">
    <!-- 主双栏容器 -->
    <div class="flex md:flex-row flex-col flex-1 overflow-hidden">
      <!-- 左侧配置面板 -->
      <div
        class="p-4 border-mk-surface md:border-r border-b md:border-b-0 w-full md:w-72 lg:w-80 overflow-y-auto transition-all duration-200 shrink-0 custom-scrollbar"
        :class="{ 'md:hidden': sidebarCollapsed }">
        <!-- MQTT 连接参数 -->
        <div class="bg-mk-surface mb-4 p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-mk-cyan text-sm">连接参数</h3>
          <div class="space-y-3">
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">地址</span>
              <input v-model="store.brokerUrl"
                class="flex-1 bg-mk-surface px-3 py-2 border w-28 border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm"
                placeholder="mqtt://127.0.0.1:1883" />
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">ClientID</span>
              <input v-model="store.clientId"
                class="flex-1 bg-mk-surface px-3 py-2 border w-28 border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm"
                placeholder="随机生成" />
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">用户名</span>
              <input v-model="store.username"
                class="flex-1 bg-mk-surface px-3 py-2 border w-28 border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm" />
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">密码</span>
              <input v-model="store.password" type="password"
                class="flex-1 bg-mk-surface px-3 py-2 border w-28 border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm" />
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">心跳(s)</span>
              <input v-model.number="store.keepalive" type="number" min="1"
                class="flex-1 bg-mk-surface px-3 py-2 border w-28 border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm" />
            </div>
          </div>
        </div>
        <!-- 连接操作按钮 -->
        <div class="bg-mk-surface mb-4 p-4 rounded-xl">
          <div class="flex gap-3">
            <button @click="handleConnect" :disabled="store.isConnected"
              class="flex-1 bg-mk-green hover:bg-mk-greenHover disabled:opacity-40 py-2 rounded-lg font-medium text-mk-foreground text-sm transition">
              连接
            </button>
            <button @click="handleDisconnect" :disabled="!store.isConnected"
              class="flex-1 bg-mk-pink hover:bg-mk-pinkHover disabled:opacity-40 py-2 rounded-lg font-medium text-white text-sm transition">
              断开
            </button>
          </div>
          <div v-if="errorMsg" class="mt-2 text-mk-pink text-xs text-center">
            {{ errorMsg }}
          </div>
        </div>
        <!-- 订阅主题管理（支持多主题） -->
        <div class="bg-mk-surface mb-4 p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-mk-cyan text-sm">订阅主题</h3>
          <div class="space-y-2">
            <div class="flex gap-2">
              <input v-model="newSubTopic" placeholder="topic/#"
                class="flex-1 bg-mk-surface px-2 py-1 border border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm" />
              <button @click="addSubscribeTopic"
                class="bg-mk-green hover:bg-mk-greenHover px-2 rounded-lg text-sm text-mk-foreground">+</button>
            </div>
            <div class="max-h-32 overflow-y-auto custom-scrollbar space-y-1">
              <div v-for="(item, idx) in store.subList" :key="idx" class="flex items-center gap-1">
                <span class="flex-1 text-mk-foreground text-xs truncate">{{ item.topic }}</span>
                <span class="text-mk-comment text-xs">QoS{{ item.qos }}</span>
                <button @click="removeSubTopic(idx)" class="text-mk-pink text-xs px-1">×</button>
              </div>
            </div>
          </div>
        </div>
        <!-- 发布设置 -->
        <div class="bg-mk-surface p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-mk-cyan text-sm">发布设置</h3>
          <div class="space-y-2">
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">主题</span>
              <input v-model="pubTopic" placeholder="pub/topic"
                class="flex-1 bg-mk-surface px-2 py-1 border w-28 border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm" />
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">QoS</span>
              <select v-model.number="pubQos"
                class="flex-1 bg-mk-surface px-2 py-1 border w-28 border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm">
                <option :value="0">0</option>
                <option :value="1">1</option>
                <option :value="2">2</option>
              </select>
            </div>
            <label class="flex items-center gap-2 text-mk-foreground text-sm mt-2">
              <input type="checkbox" v-model="store.autoPub" @change="handleAutoPubToggle" class="accent-mk-pink" />
              定时发布
              <input
                class="bg-mk-surface px-2 py-0.5 border border-mk-comment focus:border-mk-cyan rounded w-14 text-mk-yellow text-sm text-center"
                v-model.number="store.pubSec" />
              <span class="text-mk-comment">秒</span>
            </label>
          </div>
        </div>
      </div>
      <!-- 右侧消息展示区域 -->
      <div class="flex flex-col flex-1 ml-3 min-w-0 h-full">
        <!-- 顶部栏 -->
        <div class="flex justify-between items-center px-4 py-2">
          <button @click="sidebarCollapsed = !sidebarCollapsed"
            class="md:hidden bg-mk-surface hover:bg-mk-surfaceHover px-2 py-1 border border-mk-comment rounded text-mk-green text-xs">
            {{ sidebarCollapsed ? '展开设置' : '收起设置' }}
          </button>
          <span class="ml-auto text-mk-comment text-sm">消息记录</span>
        </div>
        <!-- 消息列表容器 -->
        <div class="flex flex-col flex-1 border border-mk-comment rounded-xl min-h-0">
          <div ref="msgListRef" class="flex-1 space-y-3 px-4 py-2 min-h-0 overflow-y-auto custom-scrollbar">
            <div v-for="(msg, idx) in store.messages" :key="idx" class="flex flex-col"
              :class="msg.type === 'pub' ? 'items-end' : 'items-start'">
              <!-- 气泡 -->
              <div class="flex flex-col max-w-[85%]"
                :class="msg.type === 'sub' ? 'bg-mk-surface border-mk-comment/40' : 'bg-mk-surface/15 border-mk-cyan/30'"
                style="border-width: 1px; border-radius: 1rem; border-top-left-radius: 0.25rem; border-top-right-radius: 0.25rem;">
                <!-- 消息头部：主题+QoS -->
                <div class="px-4 pt-2">
                  <div class="text-[10px] text-mk-comment flex justify-between">
                    <span>{{ msg.topic }} [QoS{{ msg.qos }}]</span>
                  </div>
                </div>
                <!-- 消息内容 -->
                <div class="px-4 py-2">
                  <div class="font-mono text-sm break-all whitespace-pre-wrap"
                    :class="msg.type === 'sub' ? 'text-mk-yellow' : 'text-mk-foreground'">
                    {{ msg.text }}
                  </div>
                  <div class="mt-1 text-[10px] text-right"
                    :class="msg.type === 'sub' ? 'text-mk-comment' : 'text-mk-green/60'">
                    {{ msg.time }}
                  </div>
                </div>
              </div>
            </div>
            <div v-if="store.messages.length === 0"
              class="flex justify-center items-center h-full text-mk-comment text-sm">
              暂无消息，请连接Broker后订阅主题
            </div>
          </div>
          <!-- 接收显示格式栏，新增 JSON -->
          <div class="flex flex-wrap items-center gap-3 bg-mk-surface px-4 py-1 border border-mk-comment border-t rounded-b-xl">
            <span class="text-mk-comment text-sm">显示格式:</span>
            <span v-for="item in recvModeList" :key="item.key" @click="store.setRecvMode(item.key)"
              class="px-2 py-1 rounded text-sm transition cursor-pointer" :class="store.recvMode === item.key
                ? 'bg-mk-green text-mk-foreground'
                : 'text-mk-foreground hover:bg-mk-surface/30'
                ">
              {{ item.name }}
            </span>
            <span class="ml-auto text-mk-comment text-xs">
              {{ store.messages.length }} 条
            </span>
            <button @click="handleClearMsg" class="text-mk-pink text-xs hover:underline">清空记录</button>
          </div>
        </div>
        <div class="h-3 shrink-0"></div>
        <!-- 发布数据输入区 -->
        <div class="bg-mk-surface border border-mk-comment rounded-xl overflow-hidden shrink-0">
          <div class="flex gap-3 p-2">
            <textarea v-model="pubPayload"
              class="flex-1 bg-mk-surface px-4 py-1 border border-mk-comment focus:border-mk-cyan rounded-xl focus:outline-none font-mono text-mk-yellow text-sm resize-none custom-scrollbar"
              :rows="pubRows" placeholder="输入发布负载payload..." @keydown.enter.ctrl="handlePublish"
              @keydown.enter.exact="handlePublish"></textarea>
            <button @click="handlePublish" :disabled="!store.isConnected || !pubTopic.trim() || !pubPayload.trim()"
              class="self-stretch bg-mk-red hover:bg-mk-pink disabled:opacity-40 px-6 rounded-xl font-medium text-white text-sm transition shrink-0">
              发布
            </button>
          </div>
          <div class="flex flex-wrap items-center gap-3 bg-mk-surface px-4 py-1 border border-mk-comment border-t">
            <span class="text-mk-comment text-sm">负载格式:</span>
            <span v-for="item in sendModeList" :key="item.key" @click="store.setSendMode(item.key)"
              class="px-2 py-1 rounded text-sm transition cursor-pointer" :class="store.sendMode === item.key
                ? 'bg-mk-green text-mk-foreground'
                : 'text-mk-foreground hover:bg-mk-surface/30'
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
import { useMqttStore } from '../stores/mqtt.js'
const store = useMqttStore()
// 接收格式，新增 json
const recvModeList = [
  { key: 'ascii', name: 'ASCII' },
  { key: 'hex', name: 'HEX' },
  { key: 'bin', name: 'BIN' },
  { key: 'dec', name: 'DEC' },
  { key: 'bcd', name: 'BCD' },
  { key: 'json', name: 'JSON' },
]
const sendModeList = [
  { key: 'ascii', name: 'ASCII' },
  { key: 'hex', name: 'HEX' },
  { key: 'bin', name: 'BIN' },
  { key: 'dec', name: 'DEC' },
  { key: 'bcd', name: 'BCD' },
  { key: 'json', name: 'JSON' },
]
// 发布相关
const pubPayload = ref('')
const pubTopic = ref('')
const pubQos = ref(0)
const newSubTopic = ref('')
const errorMsg = ref('')
const msgListRef = ref(null)
const sidebarCollapsed = ref(false)
const pubRows = ref(2)

// 订阅主题添加删除
function addSubscribeTopic() {
  if(!newSubTopic.value.trim()) return
  store.subList.push({
    topic: newSubTopic.value.trim(),
    qos:0
  })
  newSubTopic.value = ''
}
function removeSubTopic(idx) {
  store.subList.splice(idx,1)
}

// 连接断开
async function handleConnect() {
  errorMsg.value = ''
  try {
    await store.connect()
    if (window.innerWidth < 768) sidebarCollapsed.value = true
  } catch (err) {
    errorMsg.value = err?.message || String(err)
  }
}
async function handleDisconnect() {
  errorMsg.value = ''
  try {
    await store.disconnect()
    if (store.autoPub) {
      store.autoPub = false
      await store.stopAutoPub()
    }
  } catch (err) {
    errorMsg.value = err?.message || String(err)
  }
}
// 发布消息
async function handlePublish() {
  if (!store.isConnected || !pubTopic.value.trim() || !pubPayload.value.trim()) return
  errorMsg.value = ''
  try {
    await store.publish(pubTopic.value, pubPayload.value, pubQos)
    pubPayload.value = ''
  } catch (err) {
    errorMsg.value = err?.message || String(err)
  }
}
async function handleAutoPubToggle() {
  errorMsg.value = ''
  try {
    if (store.autoPub) {
      await store.startAutoPub(pubPayload.value, pubTopic.value, pubQos)
    } else {
      await store.stopAutoPub()
    }
  } catch (err) {
    errorMsg.value = err?.message || String(err)
    store.autoPub = !store.autoPub
  }
}
async function handleClearMsg() {
  errorMsg.value = ''
  try {
    await store.clearMessages()
  } catch (err) {
    errorMsg.value = err?.message || String(err)
  }
}

// 自适应文本框行数
function updatePubRows() {
  const h = window.innerHeight
  if (h < 600) pubRows.value = 1
  else if (h < 800) pubRows.value = 2
  else pubRows.value = 3
}

onMounted(() => {
  store.initEventListeners()
  updatePubRows()
  window.addEventListener('resize', updatePubRows)
})
onUnmounted(() => {
  store.destroyEventListeners()
  window.removeEventListener('resize', updatePubRows)
})

// 消息列表自动滚动到底
watch(
  () => store.messages.length,
  async () => {
    await nextTick()
    if (msgListRef.value) {
      msgListRef.value.scrollTop = msgListRef.value.scrollHeight
    }
  }
)
</script>

<style scoped>
/* 滚动条样式完全沿用串口页面，保持统一 */
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
