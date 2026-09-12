<template>
  <div class="flex flex-col w-full h-full overflow-hidden">
    <!-- 主双栏容器 -->
    <div class="flex md:flex-row flex-col flex-1 overflow-hidden">
      <!-- 左侧配置面板 -->
      <div
        class="p-4 border-mk-surface md:border-r border-b md:border-b-0 w-full md:w-72 lg:w-80 overflow-y-auto transition-all duration-200 shrink-0 custom-scrollbar"
        :class="{ 'md:hidden': sidebarCollapsed }">
        <!-- 通信模式选择 TCP Client / TCP Server / UDP -->
        <div class="bg-mk-surface mb-4 p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-mk-cyan text-sm">通信类型</h3>
          <div class="space-y-3">
            <label class="flex items-center gap-3 text-mk-foreground text-sm">
              <input type="radio" v-model="commMode" value="tcp-client" class="accent-mk-green" />
              <span>TCP Client</span>
            </label>
            <label class="flex items-center gap-3 text-mk-foreground text-sm">
              <input type="radio" v-model="commMode" value="tcp-server" class="accent-mk-green" />
              <span>TCP Server</span>
            </label>
            <label class="flex items-center gap-3 text-mk-foreground text-sm">
              <input type="radio" v-model="commMode" value="udp" class="accent-mk-green" />
              <span>UDP</span>
            </label>
          </div>
        </div>

        <!-- TCP Client 参数面板 -->
        <div v-if="commMode === 'tcp-client'" class="bg-mk-surface mb-4 p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-mk-cyan text-sm">TCP Client 参数</h3>
          <div class="space-y-3">
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">目标地址</span>
              <input v-model="store.tcpClientHost"
                class="flex-1 bg-mk-surface px-3 py-2 border border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm"
                placeholder="127.0.0.1" />
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">目标端口</span>
              <input v-model.number="store.tcpClientPort" type="number"
                class="flex-1 bg-mk-surface px-3 py-2 border border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm"
                placeholder="8080" />
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">超时(ms)</span>
              <input v-model.number="store.tcpClientTimeout" type="number"
                class="flex-1 bg-mk-surface px-3 py-2 border border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm" />
            </div>
          </div>
        </div>

        <!-- TCP Server 参数面板 -->
        <div v-if="commMode === 'tcp-server'" class="bg-mk-surface mb-4 p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-mk-cyan text-sm">TCP Server 参数</h3>
          <div class="space-y-3">
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">监听地址</span>
              <input v-model="store.tcpServerHost"
                class="flex-1 bg-mk-surface px-3 py-2 border border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm"
                placeholder="0.0.0.0" />
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">监听端口</span>
              <input v-model.number="store.tcpServerPort" type="number"
                class="flex-1 bg-mk-surface px-3 py-2 border border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm" />
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">最大连接</span>
              <input v-model.number="store.tcpServerMaxConn" type="number"
                class="flex-1 bg-mk-surface px-3 py-2 border border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm" />
            </div>
            <div class="text-mk-comment text-xs">
              当前连接数：{{ store.tcpClientConnCount }}
            </div>
          </div>
        </div>

        <!-- UDP 参数面板 -->
        <div v-if="commMode === 'udp'" class="bg-mk-surface mb-4 p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-mk-cyan text-sm">UDP 参数</h3>
          <div class="space-y-3">
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">本地端口</span>
              <input v-model.number="store.udpLocalPort" type="number"
                class="flex-1 bg-mk-surface px-3 py-2 border border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm" />
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">目标地址</span>
              <input v-model="store.udpRemoteHost"
                class="flex-1 bg-mk-surface px-3 py-2 border border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm"
                placeholder="127.0.0.1" />
            </div>
            <div class="flex items-center gap-2">
              <span class="w-16 text-mk-green text-sm shrink-0">目标端口</span>
              <input v-model.number="store.udpRemotePort" type="number"
                class="flex-1 bg-mk-surface px-3 py-2 border border-mk-comment focus:border-mk-cyan rounded-lg focus:outline-none text-mk-foreground text-sm" />
            </div>
          </div>
        </div>

        <!-- 启动/停止操作 -->
        <div class="bg-mk-surface mb-4 p-4 rounded-xl">
          <div class="flex gap-3">
            <button @click="handleStart" :disabled="store.isRunning"
              class="flex-1 bg-mk-green hover:bg-mk-greenHover disabled:opacity-40 py-2 rounded-lg font-medium text-mk-foreground text-sm transition">
              {{ commMode === 'tcp-server' ? '启动服务' : '连接' }}
            </button>
            <button @click="handleStop" :disabled="!store.isRunning"
              class="flex-1 bg-mk-pink hover:bg-mk-pinkHover disabled:opacity-40 py-2 rounded-lg font-medium text-white text-sm transition">
              {{ commMode === 'tcp-server' ? '停止服务' : '断开' }}
            </button>
          </div>
          <div v-if="errorMsg" class="mt-2 text-mk-pink text-xs text-center">
            {{ errorMsg }}
          </div>
        </div>

        <!-- 接收设置 -->
        <div class="bg-mk-surface mb-4 p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-mk-cyan text-sm">接收设置</h3>
          <div class="space-y-2">
            <label class="flex items-center gap-2 text-mk-foreground text-sm">
              <input type="checkbox" v-model="store.autoSplit" @change="store.updateAutoSplit()"
                class="accent-mk-pink" />
              自动断帧
              <input
                class="bg-mk-surface px-2 py-0.5 border border-mk-comment focus:border-mk-cyan rounded w-14 text-mk-yellow text-sm text-center"
                v-model="store.splitMs" @change="store.updateAutoSplit()" />
              <span class="text-mk-comment">ms</span>
            </label>
          </div>
          <div class="flex gap-3 mt-3">
            <button
              class="flex-1 bg-mk-surface hover:bg-mk-surfaceHover py-2 border border-mk-comment rounded-lg text-mk-yellow text-sm">
              保存数据
            </button>
            <button @click="handleClear"
              class="flex-1 bg-mk-surface hover:bg-mk-surfaceHover py-2 border border-mk-comment rounded-lg text-mk-yellow text-sm">
              清空数据
            </button>
          </div>
        </div>

        <!-- 发送设置 -->
        <div class="bg-mk-surface p-4 rounded-xl">
          <h3 class="mb-3 font-medium text-mk-cyan text-sm">发送设置</h3>
          <label class="flex items-center gap-2 text-mk-foreground text-sm">
            <input type="checkbox" v-model="store.autoSend" @change="handleAutoSendToggle" class="accent-mk-pink" />
            定时发送
            <input
              class="bg-mk-surface px-2 py-0.5 border border-mk-comment focus:border-mk-cyan rounded w-14 text-mk-yellow text-sm text-center"
              v-model.number="store.sendSec" />
            <span class="text-mk-comment">秒</span>
          </label>
        </div>
      </div>

      <!-- 右侧收发区域 -->
      <div class="flex flex-col flex-1 ml-3 min-w-0 h-full">
        <!-- 顶部栏 -->
        <div class="flex justify-between items-center px-4 py-2">
          <button @click="sidebarCollapsed = !sidebarCollapsed"
            class="md:hidden bg-mk-surface hover:bg-mk-surfaceHover px-2 py-1 border border-mk-comment rounded text-mk-green text-xs">
            {{ sidebarCollapsed ? '展开设置' : '收起设置' }}
          </button>
          <span class="ml-auto text-mk-comment text-sm">数据帧记录</span>
        </div>

        <!-- 消息列表容器 -->
        <div class="flex flex-col flex-1 border border-mk-comment rounded-xl min-h-0">
          <div ref="msgListRef" class="flex-1 space-y-3 px-4 py-2 min-h-0 overflow-y-auto custom-scrollbar">
            <div v-for="(msg, idx) in store.messages" :key="idx" class="flex flex-col"
              :class="msg.type === 'sent' ? 'items-end' : 'items-start'">
              <!-- 气泡 -->
              <div class="flex flex-col max-w-[85%]"
                :class="msg.type === 'recv' ? 'bg-mk-surface border-mk-comment/40' : 'bg-mk-surface/15 border-mk-cyan/30'"
                style="border-width: 1px; border-radius: 1rem; border-top-left-radius: 0.25rem; border-top-right-radius: 0.25rem;">
                <!-- 消息内容 -->
                <div class="px-4 py-3">
                  <div class="font-mono text-sm break-all whitespace-pre-wrap"
                    :class="msg.type === 'recv' ? 'text-mk-yellow' : 'text-mk-foreground'">
                    {{ msg.text }}
                  </div>
                  <div class="mt-1 text-[10px] text-right"
                    :class="msg.type === 'recv' ? 'text-mk-comment' : 'text-mk-green/60'">
                    {{ msg.time }}
                    <span v-if="msg.remoteAddr" class="ml-2">[{{msg.remoteAddr}}]</span>
                  </div>
                </div>
                <!-- 协议选择器 + 解析结果，复用串口逻辑 -->
                <div class="flex flex-col gap-2 px-4 pt-0 pb-3">
                  <div class="flex items-center gap-2">
                    <span class="text-mk-comment text-[10px]">协议解析:</span>
                    <select v-model="msg.protocol" @change="(e) => onProtocolChange(msg, e.target.value)"
                      class="bg-mk-surface px-2 py-0.5 border border-mk-comment rounded focus:outline-none text-mk-foreground text-xs">
                      <option :value="null">无</option>
                      <option value="modbus">Modbus RTU</option>
                      <option value="iec104">IEC104</option>
                    </select>
                  </div>
                  <div v-if="msg.protocol && msg.parsedResult"
                    class="bg-mk-surface p-2 border border-mk-comment/30 rounded font-mono text-xs">
                    <pre class="text-mk-yellow break-all whitespace-pre-wrap">{{ msg.parsedResult }}</pre>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="store.messages.length === 0"
              class="flex justify-center items-center h-full text-mk-comment text-sm">
              暂无数据，请建立连接/启动服务后收发数据
            </div>
          </div>
          <!-- 接收格式栏，包含JSON -->
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
              {{ store.messages.length }} 帧
            </span>
          </div>
        </div>

        <div class="h-3 shrink-0"></div>
        <!-- 发送区域 -->
        <div class="bg-mk-surface border border-mk-comment rounded-xl overflow-hidden shrink-0">
          <div class="flex gap-3 p-2">
            <textarea v-model="sendData"
              class="flex-1 bg-mk-surface px-4 py-1 border border-mk-comment focus:border-mk-cyan rounded-xl focus:outline-none font-mono text-mk-yellow text-sm resize-none custom-scrollbar"
              :rows="sendRows" placeholder="输入要发送的数据..." @keydown.enter.ctrl="handleSend"
              @keydown.enter.exact="handleSend"></textarea>
            <button @click="handleSend" :disabled="!store.isRunning || !sendData.trim()"
              class="self-stretch bg-mk-red hover:bg-mk-pink disabled:opacity-40 px-6 rounded-xl font-medium text-white text-sm transition shrink-0">
              发送
            </button>
          </div>
          <div class="flex flex-wrap items-center gap-3 bg-mk-surface px-4 py-1 border border-mk-comment border-t">
            <span class="text-mk-comment text-sm">发送格式:</span>
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
import { useTcpUdpStore } from '../stores/tcpudp.js'
const store = useTcpUdpStore()

// 接收格式列表，增加 JSON
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

const commMode = ref('tcp-client')
const sendData = ref('')
const errorMsg = ref('')
const msgListRef = ref(null)
const sidebarCollapsed = ref(false)
const sendRows = ref(2)

// ============ 协议解析函数【完全复用，不变】============
function getRawBytes(msg) {
  if (msg.rawData) {
    return new Uint8Array(msg.rawData)
  }
  let hex = msg.text.replace(/\s/g, '')
  if (hex.length % 2 !== 0) return new Uint8Array(0)
  const bytes = new Uint8Array(hex.length / 2)
  for (let i = 0; i < hex.length; i += 2) {
    bytes[i / 2] = parseInt(hex.substr(i, 2), 16)
  }
  return bytes
}
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
    result += `\n未深度解析的数据: ${Array.from(bytes.slice(2)).map(b => b.toString(16).padStart(2, '0')).join(' ')}`
  } else {
    result += `\n原始数据: ${Array.from(bytes).map(b => b.toString(16).padStart(2, '0')).join(' ')}`
  }
  return result
}
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
function onProtocolChange(msg, protocolValue) {
  msg.protocol = protocolValue
  if (protocolValue) {
    msg.parsedResult = parseMessageByProtocol(msg, protocolValue)
  } else {
    msg.parsedResult = ''
  }
}
// =============================================

watch(() => store.messages, (newMsgs) => {
  newMsgs.forEach(msg => {
    if (msg.protocol === undefined) {
      msg.protocol = null
      msg.parsedResult = ''
    }
  })
}, { deep: true, immediate: true })

onMounted(() => {
  store.initEventListeners()
  store.updateAutoSplit()
  updateSendRows()
  window.addEventListener('resize', updateSendRows)
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

async function handleStart() {
  errorMsg.value = ''
  try {
    await store.start(commMode.value)
    if (window.innerWidth < 768) sidebarCollapsed.value = true
  } catch (err) {
    errorMsg.value = err?.message || String(err)
  }
}
async function handleStop() {
  errorMsg.value = ''
  try {
    await store.stop()
    if (store.autoSend) {
      store.autoSend = false
      await store.stopAutoSend()
    }
  } catch (err) {
    errorMsg.value = err?.message || String(err)
  }
}
async function handleSend() {
  if (!store.isRunning || !sendData.value.trim()) return
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
