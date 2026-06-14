<template>
  <div class="flex flex-col bg-[#272822] w-full h-full overflow-hidden text-[#F8F8F2]">
    <!-- 顶部标题栏 -->
    <header class="flex justify-between items-center bg-[#3E3D32] px-4 py-3 border-[#75715E] border-b shrink-0">
      <h2 class="font-medium text-[#66D9EF] text-lg">DL/T 645 协议解析工具</h2>
      <div class="flex gap-3">
        <button
          @click="fillExample"
          class="bg-[#32332E] hover:bg-[#4e4d40] px-3 py-1.5 border border-[#75715E] rounded-lg text-[#E6DB74] text-sm transition"
        >
          示例报文
        </button>
        <button
          @click="parseFrame"
          class="bg-[#A6E22E] hover:bg-[#93c725] px-3 py-1.5 rounded-lg text-[#272822] text-sm transition"
        >
          立即解析
        </button>
        <button
          @click="clearAll"
          class="bg-[#32332E] hover:bg-[#4e4d40] px-3 py-1.5 border border-[#75715E] rounded-lg text-[#E6DB74] text-sm transition"
        >
          清空
        </button>
      </div>
    </header>

    <!-- 输入区 + 版本选择 -->
    <div class="p-4 border-[#75715E] border-b shrink-0">
      <div class="flex justify-between items-center mb-2">
        <div class="text-[#75715E] text-sm">原始报文（十六进制，支持空格分隔）：</div>
        <div class="flex gap-2">
          <label class="flex items-center gap-1 text-sm">
            <input type="radio" v-model="protocolVersion" value="1997" class="accent-[#A6E22E]" />
            <span>DL/T 645-1997</span>
          </label>
          <label class="flex items-center gap-1 text-sm">
            <input type="radio" v-model="protocolVersion" value="2007" class="accent-[#A6E22E]" />
            <span>DL/T 645-2007</span>
          </label>
        </div>
      </div>
      <textarea
        v-model="rawHex"
        @input="onInputChange"
        class="bg-[#32332E] p-3 border border-[#75715E] focus:border-[#66D9EF] rounded-xl focus:outline-none w-full font-mono text-[#E6DB74] text-sm resize-none"
        rows="2"
        placeholder="例如：68 12 90 78 56 34 12 68 01 02 43 43 33 33 33 33 5D 16"
      ></textarea>

      <!-- 字节预览区 -->
      <div v-if="frameBytes.length" class="flex flex-wrap gap-1.5 bg-[#32332E] mt-3 p-3 border border-[#75715E] rounded-xl">
        <span
          v-for="(byte, idx) in frameBytes"
          :key="idx"
          class="flex justify-center items-center border border-[#444] rounded w-10 h-8 font-mono text-sm"
          :style="getByteBgStyle(idx)"
        >
          {{ byte }}
        </span>
      </div>
      <div class="flex flex-wrap gap-x-4 gap-y-2 mt-3 text-sm">
        <span class="text-[#75715E]">区块图例：</span>
        <div class="flex items-center gap-1"><span class="border border-[#75715E] rounded w-4 h-4" style="background-color: #3b82f6;"></span><span>帧头/起始</span></div>
        <div class="flex items-center gap-1"><span class="border border-[#75715E] rounded w-4 h-4" style="background-color: #ef4444;"></span><span>地址域</span></div>
        <div class="flex items-center gap-1"><span class="border border-[#75715E] rounded w-4 h-4" style="background-color: #10b981;"></span><span>控制码/长度</span></div>
        <div class="flex items-center gap-1"><span class="border border-[#75715E] rounded w-4 h-4" style="background-color: #fbbf24;"></span><span>数据标识+数据域</span></div>
        <div class="flex items-center gap-1"><span class="border border-[#75715E] rounded w-4 h-4" style="background-color: #a78bfa;"></span><span>校验+结束符</span></div>
      </div>
    </div>

    <!-- 解析结果区域 -->
    <div class="flex-1 p-4 overflow-y-auto custom-scrollbar">
      <div v-if="isParsing" class="mb-2 text-[#E6DB74] text-sm">解析中...</div>

      <!-- 解析结果卡片 -->
      <div v-if="parsed.frame" class="bg-[#3E3D32] mb-4 p-4 border border-[#75715E]/40 rounded-xl">
        <h3 class="mb-3 font-semibold text-[#66D9EF]">帧结构解析</h3>
        <div class="gap-3 grid grid-cols-1 md:grid-cols-2 text-sm">
          <div><span class="text-[#75715E]">起始符：</span>{{ parsed.frame.start }}</div>
          <div><span class="text-[#75715E]">地址域 (倒序)：</span>{{ parsed.frame.address }}</div>
          <div><span class="text-[#75715E]">起始符2：</span>{{ parsed.frame.start2 }}</div>
          <div><span class="text-[#75715E]">控制码：</span>{{ parsed.frame.control }} ({{ parsed.frame.controlDesc }})</div>
          <div><span class="text-[#75715E]">数据长度 (L)：</span>{{ parsed.frame.dataLen }} 字节</div>
          <div v-if="parsed.frame.dataId"><span class="text-[#75715E]">数据标识 (DI)：</span>{{ parsed.frame.dataId }} ({{ parsed.frame.dataIdDesc }})</div>
          <div class="col-span-2"><span class="text-[#75715E]">数据域：</span>{{ parsed.frame.dataField || '无' }}</div>
          <div><span class="text-[#75715E]">校验和 (CS)：</span>{{ parsed.frame.cs }} <span :class="parsed.frame.csValid ? 'text-[#A6E22E]' : 'text-[#F92672]'">({{ parsed.frame.csValid ? '正确' : '错误' }})</span></div>
          <div><span class="text-[#75715E]">结束符：</span>{{ parsed.frame.end }}</div>
        </div>
      </div>

      <!-- 翻译后的数据值 -->
      <div v-if="parsed.value" class="bg-[#3E3D32] mb-4 p-4 border border-[#75715E]/40 rounded-xl">
        <h3 class="mb-3 font-semibold text-[#A6E22E]">数据解析</h3>
        <div class="text-sm">
          <div><span class="text-[#75715E]">物理量：</span>{{ parsed.value.quantity }}</div>
          <div><span class="text-[#75715E]">数值：</span>{{ parsed.value.value }} {{ parsed.value.unit }}</div>
          <div v-if="parsed.value.raw"><span class="text-[#75715E]">原始字节：</span>{{ parsed.value.raw }}</div>
        </div>
      </div>

      <div v-if="errorMsg" class="bg-[#F92672]/20 p-3 border border-[#F92672] rounded-xl text-[#F92672] text-sm">
        ⚠️ {{ errorMsg }}
      </div>
      <div v-if="!rawHex" class="py-10 text-[#75715E] text-center">
        请输入 DL/T 645 报文进行解析
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const rawHex = ref('')
const frameBytes = ref([])
const protocolVersion = ref('2007')   // 默认 2007 版
const parsed = ref({ frame: null, value: null })
const errorMsg = ref('')
const isParsing = ref(false)

// 数据标识映射表 (2007版 4字节 DI，97版 2字节 DI，这里以2007为例)
const diMap = {
  '00000000': { name: '正向有功总电能', unit: 'kWh', factor: 0.01 },
  '00010000': { name: '反向有功总电能', unit: 'kWh', factor: 0.01 },
  '00020000': { name: '正向无功总电能', unit: 'kvarh', factor: 0.01 },
  '00030000': { name: '反向无功总电能', unit: 'kvarh', factor: 0.01 },
  '02010100': { name: 'A相电压', unit: 'V', factor: 0.1 },
  '02010200': { name: 'B相电压', unit: 'V', factor: 0.1 },
  '02010300': { name: 'C相电压', unit: 'V', factor: 0.1 },
  '02020100': { name: 'A相电流', unit: 'A', factor: 0.01 },
  '02020200': { name: 'B相电流', unit: 'A', factor: 0.01 },
  '02020300': { name: 'C相电流', unit: 'A', factor: 0.01 },
  '02030000': { name: '瞬时总有功功率', unit: 'kW', factor: 0.01 },
  '02030100': { name: 'A相有功功率', unit: 'kW', factor: 0.01 },
  '02030200': { name: 'B相有功功率', unit: 'kW', factor: 0.01 },
  '02030300': { name: 'C相有功功率', unit: 'kW', factor: 0.01 },
  '00000001': { name: '当前时间', unit: '', factor: 1 },
  // 可继续添加更多DI
}

// 控制码含义
const controlMap = {
  0x81: '读数据',
  0x01: '读数据响应',
  0x82: '读后续数据',
  0x02: '读后续数据响应',
  0x83: '写数据',
  0x03: '写数据响应',
  0x84: '广播校时',
  0x04: '广播校时响应',
  0x85: '更改通信速率',
  0x05: '更改通信速率响应',
  0x86: '修改密码',
  0x06: '修改密码响应',
  0x87: '最大需量清零',
  0x07: '最大需量清零响应',
  0x88: '电表清零',
  0x08: '电表清零响应',
  0x89: '事件清零',
  0x09: '事件清零响应',
  0x8A: '冻结',
  0x0A: '冻结响应'
}

// 更新字节数组
function updateFrameBytes() {
  const hexStr = rawHex.value.replace(/\s/g, '')
  const bytes = []
  for (let i = 0; i < hexStr.length; i += 2) {
    const byte = hexStr.substr(i, 2).toUpperCase()
    if (/^[0-9A-F]{2}$/.test(byte)) bytes.push(byte)
  }
  frameBytes.value = bytes
}

// 字节索引背景色（基于645帧结构的大致分区）
function getByteBgStyle(idx) {
  if (idx === 0) return { backgroundColor: '#3b82f6' }        // 起始符68
  if (idx >= 1 && idx <= 6) return { backgroundColor: '#ef4444' } // 地址域 6字节
  if (idx === 7) return { backgroundColor: '#3b82f6' }        // 第二个68
  if (idx === 8) return { backgroundColor: '#10b981' }        // 控制码
  if (idx === 9) return { backgroundColor: '#10b981' }        // 数据长度L
  const dataLen = frameBytes.value[9] ? parseInt(frameBytes.value[9], 16) : 0
  const diLen = protocolVersion.value === '2007' ? 4 : 2
  const dataStart = 10 + diLen
  if (idx >= 10 && idx < dataStart) return { backgroundColor: '#fbbf24' } // 数据标识区
  if (idx >= dataStart && idx < dataStart + dataLen) return { backgroundColor: '#fbbf24' } // 数据域
  if (idx === dataStart + dataLen) return { backgroundColor: '#a78bfa' } // 校验和
  if (idx === dataStart + dataLen + 1) return { backgroundColor: '#a78bfa' } // 结束符16
  return { backgroundColor: 'transparent' }
}

// 计算校验和（累加和取低字节）
function calculateCS(bytes, start, end) {
  let sum = 0
  for (let i = start; i <= end; i++) {
    sum += bytes[i]
  }
  return (sum & 0xFF).toString(16).toUpperCase().padStart(2, '0')
}

// 将BCD码字节数组转换为数值（支持小数位）
function bcdToNumber(bcdBytes, factor = 0.01) {
  let hexStr = ''
  for (let b of bcdBytes) {
    hexStr += b.toString(16).padStart(2, '0')
  }
  // 去除可能的尾随F（填充字节）
  hexStr = hexStr.replace(/F+$/, '')
  const intVal = parseInt(hexStr, 10)
  return intVal * factor
}

// BCD时间转换 (秒分时日周月年)
function bcdToDateTime(bcdBytes) {
  if (bcdBytes.length < 6) return ''
  const sec = ((bcdBytes[0] >> 4) * 10 + (bcdBytes[0] & 0x0F)).toString().padStart(2,'0')
  const min = ((bcdBytes[1] >> 4) * 10 + (bcdBytes[1] & 0x0F)).toString().padStart(2,'0')
  const hour = ((bcdBytes[2] >> 4) * 10 + (bcdBytes[2] & 0x0F)).toString().padStart(2,'0')
  const day = ((bcdBytes[3] >> 4) * 10 + (bcdBytes[3] & 0x0F)).toString().padStart(2,'0')
  const month = ((bcdBytes[4] >> 4) * 10 + (bcdBytes[4] & 0x0F)).toString().padStart(2,'0')
  const year = ((bcdBytes[5] >> 4) * 10 + (bcdBytes[5] & 0x0F)).toString().padStart(4,'20')
  return `${year}-${month}-${day} ${hour}:${min}:${sec}`
}

// 主解析函数
async function parseFrame() {
  errorMsg.value = ''
  parsed.value = { frame: null, value: null }
  if (frameBytes.value.length === 0) return

  isParsing.value = true
  try {
    // 优先调用后端解析（如果存在）
    if (window.go && window.go.main && window.go.main.App && window.go.main.App.Parse645) {
      const hex = rawHex.value.replace(/\s/g, '')
      const result = await window.go.main.App.Parse645(hex, protocolVersion.value)
      parsed.value = result
      return
    }

    // 前端简易解析（用于演示）
    const bytes = frameBytes.value.map(b => parseInt(b, 16))
    if (bytes[0] !== 0x68) throw new Error('起始符错误，应为68')
    if (bytes[7] !== 0x68) throw new Error('第二个起始符错误，应为68')
    const address = bytes.slice(1, 7).reverse().map(b => b.toString(16).padStart(2,'0')).join('')
    const control = bytes[8]
    const controlDesc = controlMap[control] || '未知'
    const dataLen = bytes[9]   // 数据域长度（不包括校验和与结束符）
    const diLen = protocolVersion.value === '2007' ? 4 : 2
    if (bytes.length < 10 + diLen + dataLen + 2) throw new Error('报文长度不足')
    let dataId = ''
    let dataIdDesc = ''
    const diBytes = bytes.slice(10, 10 + diLen)
    if (protocolVersion.value === '2007') {
      // DI 顺序：DI0 DI1 DI2 DI3
      dataId = diBytes.map(b => b.toString(16).padStart(2,'0')).join('')
    } else {
      // 97版 DI 为两个字节，顺序为 DI1 DI0
      dataId = diBytes[1].toString(16).padStart(2,'0') + diBytes[0].toString(16).padStart(2,'0')
    }
    const diInfo = diMap[dataId] || { name: `未知标识 ${dataId}`, unit: '', factor: 1 }
    dataIdDesc = diInfo.name

    const dataFieldBytes = bytes.slice(10 + diLen, 10 + diLen + dataLen)
    let dataFieldHex = dataFieldBytes.map(b => b.toString(16).padStart(2,'0')).join(' ')
    let valueObj = null
    // 尝试解析数据值
    if (dataLen > 0) {
      if (dataId === '00000001' && protocolVersion.value === '2007') {
        // 时间数据 BCD 6字节 秒分时日周月年
        const dateTime = bcdToDateTime(dataFieldBytes)
        valueObj = { quantity: '当前时间', value: dateTime, unit: '', raw: dataFieldHex }
      } else if (diInfo.name) {
        const numericVal = bcdToNumber(dataFieldBytes, diInfo.factor)
        valueObj = { quantity: diInfo.name, value: numericVal, unit: diInfo.unit, raw: dataFieldHex }
      } else {
        valueObj = { quantity: dataIdDesc, value: dataFieldHex, unit: '', raw: dataFieldHex }
      }
    }

    // 校验和
    const csPos = 10 + diLen + dataLen
    const csByte = bytes[csPos]
    const csHex = csByte.toString(16).toUpperCase().padStart(2,'0')
    const computedCS = calculateCS(bytes, 0, csPos - 1)
    const csValid = csHex === computedCS
    const endByte = bytes[csPos + 1]
    if (endByte !== 0x16) throw new Error('结束符错误，应为16')

    parsed.value.frame = {
      start: '68',
      address,
      start2: '68',
      control: `0x${control.toString(16).padStart(2,'0')}`,
      controlDesc,
      dataLen,
      dataId,
      dataIdDesc,
      dataField: dataFieldHex || '无',
      cs: csHex,
      csValid,
      end: '16'
    }
    if (valueObj) parsed.value.value = valueObj
  } catch (err) {
    errorMsg.value = err.message || String(err)
  } finally {
    isParsing.value = false
  }
}

function onInputChange() {
  updateFrameBytes()
  if (rawHex.value.trim()) parseFrame()
  else clearAll()
}

function fillExample() {
  // 读正向有功总电能请求报文 (2007版)
  // 68 AA AA AA AA AA AA 68 01 02 43 43 33 33 33 33 16
  // 这里提供一个更标准的示例：读正向有功总电能请求，地址为 123456789012
  rawHex.value = '68 12 90 78 56 34 12 68 01 02 43 43 33 33 33 33 5D 16'
  protocolVersion.value = '2007'
  onInputChange()
}

function clearAll() {
  rawHex.value = ''
  frameBytes.value = []
  parsed.value = { frame: null, value: null }
  errorMsg.value = ''
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