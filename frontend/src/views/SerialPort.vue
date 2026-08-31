<template>
  <div class="flex flex-col w-full h-full overflow-hidden">
    <!-- 主双栏容器 -->
    <div class="flex md:flex-row flex-col flex-1 overflow-hidden">
      <!-- 左侧配置面板（保持不变） -->
      <div
        class="p-4 border-monokai-surface md:border-r border-b md:border-b-0 w-full md:w-72 lg:w-80 overflow-y-auto transition-all duration-200 shrink-0 custom-scrollbar"
        :class="{ 'md:hidden': sidebarCollapsed }">
        <!-- 串口参数 -->
        <div class="bg-monokai-surface mb-4 p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-monokai-cyan text-sm">串口参数</h3>
          <div class="space-y-3">
            <div class="flex items-center gap-2">
              <span class="w-16 text-monokai-green text-sm shrink-0">端口名</span>
              <select v-model="store.portName"
                class="flex-1 bg-monokai-surface px-3 py-2 border border-monokai-comment focus:border-monokai-cyan rounded-lg focus:outline-none text-monokai-foreground text-sm">
                <option v-for="p in store.ports" :key="p.name" :value="p.name">{{ p.name }}</option>
                <option v-if="store.ports.length === 0" value="COM5">COM5</option>
              </select>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-monokai-green text-sm shrink-0">波特率</span>
              <select v-model="store.baudRate"
                class="flex-1 bg-monokai-surface px-3 py-2 border border-monokai-comment focus:border-monokai-cyan rounded-lg focus:outline-none text-monokai-foreground text-sm">
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
              <span class="w-16 text-monokai-green text-sm shrink-0">数据位</span>
              <select v-model="store.dataBits"
                class="flex-1 bg-monokai-surface px-3 py-2 border border-monokai-comment focus:border-monokai-cyan rounded-lg focus:outline-none text-monokai-foreground text-sm">
                <option :value="5">5</option>
                <option :value="6">6</option>
                <option :value="7">7</option>
                <option :value="8">8</option>
              </select>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-monokai-green text-sm shrink-0">校验位</span>
              <select v-model="store.parity"
                class="flex-1 bg-monokai-surface px-3 py-2 border border-monokai-comment focus:border-monokai-cyan rounded-lg focus:outline-none text-monokai-foreground text-sm">
                <option>无</option>
                <option>奇</option>
                <option>偶</option>
                <option>Mark</option>
                <option>Space</option>
              </select>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-monokai-green text-sm shrink-0">停止位</span>
              <select v-model="store.stopBits"
                class="flex-1 bg-monokai-surface px-3 py-2 border border-monokai-comment focus:border-monokai-cyan rounded-lg focus:outline-none text-monokai-foreground text-sm">
                <option :value="1">1</option>
                <option :value="2">2</option>
              </select>
            </div>
          </div>
        </div>

        <!-- 串口操作按钮 -->
        <div class="bg-monokai-surface mb-4 p-4 rounded-xl">
          <div class="flex gap-3">
            <button @click="handleOpen" :disabled="store.isOpen"
              class="flex-1 bg-monokai-green hover:bg-monokai-greenHover disabled:opacity-40 py-2 rounded-lg font-medium text-monokai-foreground text-sm transition">
              打开串口
            </button>
            <button @click="handleClose" :disabled="!store.isOpen"
              class="flex-1 bg-monokai-pink hover:bg-monokai-pinkHover disabled:opacity-40 py-2 rounded-lg font-medium text-white text-sm transition">
              关闭串口
            </button>
          </div>
          <div v-if="errorMsg" class="mt-2 text-monokai-pink text-xs text-center">
            {{ errorMsg }}
          </div>
        </div>

        <!-- 接收设置 -->
        <div class="bg-monokai-surface mb-4 p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-monokai-cyan text-sm">接收设置</h3>
          <div class="space-y-2">
            <label class="flex items-center gap-2 text-monokai-foreground text-sm">
              <input type="checkbox" v-model="store.autoSplit" @change="store.updateAutoSplit()"
                class="accent-monokai-pink" />
              自动断帧
              <input
                class="bg-monokai-surface px-2 py-0.5 border border-monokai-comment focus:border-monokai-cyan rounded w-14 text-monokai-yellow text-sm text-center"
                v-model="store.splitMs" @change="store.updateAutoSplit()" />
              <span class="text-monokai-comment">ms</span>
            </label>
          </div>  
          <div class="flex gap-3 mt-3">
            <button
              class="flex-1 bg-monokai-surface hover:bg-monokai-surfaceHover py-2 border border-monokai-comment rounded-lg text-monokai-yellow text-sm">
              保存数据
            </button>
            <button @click="handleClear"
              class="flex-1 bg-monokai-surface hover:bg-monokai-surfaceHover py-2 border border-monokai-comment rounded-lg text-monokai-yellow text-sm">
              清空数据
            </button>
          </div>
        </div>

        <!-- 发送设置 -->
        <div class="bg-monokai-surface p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-monokai-cyan text-sm">发送设置</h3>
          <label class="flex items-center gap-2 text-monokai-foreground text-sm">
            <input type="checkbox" v-model="store.autoSend" @change="handleAutoSendToggle" class="accent-monokai-pink" />
            定时发送
            <input
              class="bg-monokai-surface px-2 py-0.5 border border-monokai-comment focus:border-monokai-cyan rounded w-14 text-monokai-yellow text-sm text-center"
              v-model.number="store.sendSec" />
            <span class="text-monokai-comment">秒</span>
          </label>
        </div>
      </div>

      <!-- 右侧收发区域 -->
      <div class="flex flex-col flex-1 ml-3 min-w-0 h-full">
        <!-- 顶部栏 -->
        <div class="flex justify-between items-center px-4 py-2">
          <button @click="sidebarCollapsed = !sidebarCollapsed"
            class="md:hidden bg-monokai-surface hover:bg-monokai-surfaceHover px-2 py-1 border border-monokai-comment rounded text-monokai-green text-xs">
            {{ sidebarCollapsed ? '展开设置' : '收起设置' }}
          </button>
          <span class="ml-auto text-monokai-comment text-sm">数据帧记录</span>
        </div>

        <!-- 消息列表容器 -->
        <div class="flex flex-col flex-1 border border-monokai-comment rounded-xl min-h-0">
          <div ref="msgListRef" class="flex-1 space-y-3 px-4 py-2 min-h-0 overflow-y-auto custom-scrollbar">
            <div v-for="(msg, idx) in store.messages" :key="idx" class="flex flex-col"
              :class="msg.type === 'sent' ? 'items-end' : 'items-start'">
              <!-- 气泡 -->
              <div class="flex flex-col max-w-[85%]"
                :class="msg.type === 'recv' ? 'bg-monokai-surface border-monokai-comment/40' : 'bg-monokai-surface/15 border-monokai-cyan/30'"
                style="border-width: 1px; border-radius: 1rem; border-top-left-radius: 0.25rem; border-top-right-radius: 0.25rem;">
                <!-- 消息内容 -->
                <div class="px-4 py-3">
                  <div class="font-mono text-sm break-all whitespace-pre-wrap"
                    :class="msg.type === 'recv' ? 'text-monokai-yellow' : 'text-monokai-foreground'">
                    {{ msg.text }}
                  </div>
                  <div class="mt-1 text-[10px] text-right"
                    :class="msg.type === 'recv' ? 'text-monokai-comment' : 'text-monokai-green/60'">
                    {{ msg.time }}
                  </div>
                </div>
                <!-- 协议选择器 + 解析结果 -->
                <div class="flex flex-col gap-2 px-4 pt-0 pb-3">
                  <div class="flex items-center gap-2">
                    <span class="text-monokai-comment text-[10px]">协议解析:</span>
                    <select v-model="msg.protocol" @change="(e) => onProtocolChange(msg, e.target.value)"
                      class="bg-monokai-surface px-2 py-0.5 border border-monokai-comment rounded focus:outline-none text-monokai-foreground text-xs">
                      <option :value="null">无</option>
                      <option value="modbus">Modbus RTU</option>
                      <option value="iec104">IEC104</option>
                    </select>
                  </div>
                  <div v-if="msg.protocol && msg.parsedResult"
                    class="bg-monokai-surface p-2 border border-monokai-comment/30 rounded font-mono text-xs">
                    <pre class="text-monokai-yellow break-all whitespace-pre-wrap">{{ msg.parsedResult }}</pre>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="store.messages.length === 0"
              class="flex justify-center items-center h-full text-monokai-comment text-sm">
              暂无数据，请打开串口后收发数据
            </div>
          </div>

          <!-- 接收格式栏 -->
          <div class="flex flex-wrap items-center gap-3 bg-monokai-surface px-4 py-1 border border-monokai-comment border-t rounded-b-xl">
            <span class="text-monokai-comment text-sm">显示格式:</span>
            <span v-for="item in recvModeList" :key="item.key" @click="store.setRecvMode(item.key)"
              class="px-2 py-1 rounded text-sm transition cursor-pointer" :class="store.recvMode === item.key
                ? 'bg-monokai-green text-monokai-foreground]'
                : 'text-monokai-foreground] hover:bg-monokai-surface/30'
                ">
              {{ item.name }}
            </span>
            <span class="ml-auto text-monokai-comment text-xs">
              {{ store.messages.length }} 帧
            </span>
          </div>
        </div>

        <div class="h-3 shrink-0"></div>

        <!-- 发送区域 -->
        <div class="bg-monokai-surface border border-monokai-comment rounded-xl overflow-hidden shrink-0">
          <div class="flex gap-3 p-2">
            <textarea v-model="sendData"
              class="flex-1 bg-monokai-surface px-4 py-1 border border-monokai-comment focus:border-monokai-cyan rounded-xl focus:outline-none font-mono text-monokai-yellow text-sm resize-none custom-scrollbar"
              :rows="sendRows" placeholder="输入要发送的数据..." @keydown.enter.ctrl="handleSend"
              @keydown.enter.exact="handleSend"></textarea>
            <button @click="handleSend" :disabled="!store.isOpen || !sendData.trim()"
              class="self-stretch bg-monokai-red hover:bg-monokai-pink disabled:opacity-40 px-6 rounded-xl font-medium text-white text-sm transition shrink-0">
              发送
            </button>
          </div>
          <div class="flex flex-wrap items-center gap-3 bg-monokai-surface px-4 py-1 border border-monokai-comment border-t">
            <span class="text-monokai-comment text-sm">发送格式:</span>
            <span v-for="item in sendModeList" :key="item.key" @click="store.setSendMode(item.key)"
              class="px-2 py-1 rounded text-sm transition cursor-pointer" :class="store.sendMode === item.key
                ? 'bg-monokai-green text-monokai-foreground]'
                : 'text-monokai-foreground] hover:bg-monokai-surface/30'
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
const sidebarCollapsed = ref(false)
const sendRows = ref(2)

// ---------- 协议解析函数 ----------
// 从消息中获取原始字节数组（优先使用 rawData，否则从显示的十六进制文本转换）
function getRawBytes(msg) {
  if (msg.rawData) {
    return new Uint8Array(msg.rawData)
  }
  // 兼容旧消息：假设 msg.text 是 HEX 格式（无空格或带空格）
  let hex = msg.text.replace(/\s/g, '')
  if (hex.length % 2 !== 0) return new Uint8Array(0)
  const bytes = new Uint8Array(hex.length / 2)
  for (let i = 0; i < hex.length; i += 2) {
    bytes[i / 2] = parseInt(hex.substr(i, 2), 16)
  }
  return bytes
}

// Modbus RTU 简易解析（地址、功能码、寄存器地址、数据）
function parseModbusRTU(bytes) {
  if (bytes.length < 4) return '数据过短，无法解析'
  const addr = bytes[0]
  const func = bytes[1]
  let result = `地址: 0x${addr.toString(16)} (${addr})\n功能码: 0x${func.toString(16)}`
  if (func === 0x03 || func === 0x06 || func === 0x10) {
    if (bytes.length >= 4) {
      const regAddr = (bytes[2] << 8) | bytes[3]
      result += `\n寄存器地址: ${regAddr}`
    }
    if (func === 0x03) {
      const byteCount = bytes[2]
      result += `\n字节数: ${byteCount}`
      if (bytes.length >= 4 + byteCount) {
        const values = []
        for (let i = 0; i < byteCount / 2; i++) {
          const val = (bytes[4 + i * 2] << 8) | bytes[5 + i * 2]
          values.push(val)
        }
        result += `\n数据: ${values.join(', ')}`
      }
    } else if (func === 0x06) {
      const data = (bytes[4] << 8) | bytes[5]
      result += `\n写入值: ${data}`
    } else if (func === 0x10) {
      const byteCount = bytes[6]
      result += `\n寄存器数量: ${byteCount / 2}`
      const dataBytes = bytes.slice(7, 7 + byteCount)
      result += `\n数据: ${Array.from(dataBytes).map(b => b.toString(16).padStart(2, '0')).join(' ')}`
    }
  } else if (func >= 0x01 && func <= 0x06) {
    // 其他常用功能码简单处理
    result += `\n未深度解析的数据: ${Array.from(bytes.slice(2)).map(b => b.toString(16).padStart(2, '0')).join(' ')}`
  } else {
    result += `\n原始数据: ${Array.from(bytes).map(b => b.toString(16).padStart(2, '0')).join(' ')}`
  }
  return result
}

// IEC104 简易解析（复用之前的逻辑，但简化输出为文本）
function parseIEC104(bytes) {
  if (bytes.length < 6) return '报文太短'
  if (bytes[0] !== 0x68) return '起始字节不是 68'
  const apduLen = bytes[1]
  if (bytes.length < apduLen + 2) return '长度不匹配'

  const ctrl1 = bytes[2], ctrl2 = bytes[3], ctrl3 = bytes[4], ctrl4 = bytes[5]
  let frameType = '未知'
  if ((ctrl1 & 0x01) === 0 && (ctrl2 & 0x01) === 0) frameType = 'I-帧'
  else if ((ctrl1 & 0x01) === 1 && (ctrl2 & 0x01) === 0 && (ctrl1 & 0x02) === 0) frameType = 'S-帧'
  else frameType = 'U-帧'

  let result = `帧类型: ${frameType}\n`
  if (frameType === 'I-帧') {
    const sendSeq = ((ctrl2 & 0xFE) >> 1) | ((ctrl1 & 0x01) << 7)
    const recvSeq = ((ctrl4 & 0xFE) >> 1) | ((ctrl3 & 0x01) << 7)
    result += `发送序号: ${sendSeq}, 接收序号: ${recvSeq}\n`
  } else if (frameType === 'S-帧') {
    const recvSeq = ((ctrl4 & 0xFE) >> 1) | ((ctrl3 & 0x01) << 7)
    result += `接收序号: ${recvSeq}\n`
  } else {
    if (ctrl1 === 0x07 && ctrl2 === 0x00) result += '原因: 启动帧\n'
    else if (ctrl1 === 0x0B && ctrl2 === 0x00) result += '原因: 停止帧\n'
    else if (ctrl1 === 0x03 && ctrl2 === 0x00) result += '原因: 测试帧\n'
  }

  if (frameType === 'I-帧' && bytes.length > 6) {
    let pos = 6
    const typeId = bytes[pos++]
    const vsq = bytes[pos++]
    const causeTx = bytes[pos++]
    const originAddr = bytes[pos++]
    const commonAddrLow = bytes[pos++]
    const commonAddrHigh = bytes[pos++]
    const commonAddr = (commonAddrHigh << 8) | commonAddrLow
    const infoCount = vsq & 0x7F
    result += `ASDU: 类型=${typeId.toString(16).padStart(2, '0')}, 信息体个数=${infoCount}, 传送原因=${causeTx}, 公共地址=${commonAddr}\n`
    if (infoCount > 0 && pos + 3 < bytes.length) {
      result += `信息体1地址: ${bytes[pos] | (bytes[pos + 1] << 8) | (bytes[pos + 2] << 16)}\n`
    }
  }
  return result
}

// 根据协议解析消息，返回可读字符串
function parseMessageByProtocol(msg, protocol) {
  if (!protocol) return ''
  const bytes = getRawBytes(msg)
  if (bytes.length === 0) return '无有效数据'
  try {
    if (protocol === 'modbus') return parseModbusRTU(bytes)
    if (protocol === 'iec104') return parseIEC104(bytes)
    return '不支持的协议'
  } catch (e) {
    return `解析错误: ${e.message}`
  }
}

// 协议改变时的处理函数
function onProtocolChange(msg, protocolValue) {
  msg.protocol = protocolValue
  if (protocolValue) {
    msg.parsedResult = parseMessageByProtocol(msg, protocolValue)
  } else {
    msg.parsedResult = ''
  }
  // 触发视图更新（Vue 会自动响应，因为 msg 是响应式对象）
}

// 当新消息加入时，自动为消息对象增加协议相关字段（如果不存在）
// 注意：需要在 store 中保证每条消息是响应式的。如果 store 使用 push 添加，Vue 3 的 reactive 会自动处理。
// 我们在这里使用一个 watcher 或修改 store 的添加逻辑，但简单的方式是直接在模板中初始化。
// 更好的办法：在接收消息时对 msg 进行初始化。我们可以在 onMounted 时对已有消息初始化，并在每次新消息到来时也做初始化。
// 由于 store 是外部提供，我们可以 watch store.messages 并处理新增项。
watch(() => store.messages, (newMsgs) => {
  newMsgs.forEach(msg => {
    if (msg.protocol === undefined) {
      msg.protocol = null
      msg.parsedResult = ''
    }
  })
}, { deep: true, immediate: true })

// 原有生命周期和函数保持不变
onMounted(() => {
  store.initEventListeners()
  store.fetchPorts()
  store.updateAutoSplit()
  updateSendRows()
  window.addEventListener('resize', updateSendRows)
  // 初始化现有消息
  store.messages.forEach(msg => {
    if (msg.protocol === undefined) {
      msg.protocol = null
      msg.parsedResult = ''
    }
  })
})

onUnmounted(() => {
  store.destroyEventListeners()
  window.removeEventListener('resize', updateSendRows)
})

watch(
  () => store.messages.length,
  async () => {
    await nextTick()
    if (msgListRef.value) {
      msgListRef.value.scrollTop = msgListRef.value.scrollHeight
    }
  }
)

function updateSendRows() {
  const h = window.innerHeight
  if (h < 600) sendRows.value = 1
  else if (h < 800) sendRows.value = 2
  else sendRows.value = 3
}

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
/* 滚动条样式保持不变 */
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