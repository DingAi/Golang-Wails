<template>
    <div class="flex flex-col bg-monokai-bg w-full h-full overflow-hidden text-monokai-foreground">
        <!-- 顶部标题栏 -->
        <header class="flex justify-between items-center bg-monokai-surface px-4 py-3 border-monokai-comment border-b shrink-0">
            <h2 class="font-medium text-monokai-cyan text-lg">DL/T 698.45 协议解析工具</h2>
            <div class="flex gap-3">
                <button @click="fillExample"
                    class="bg-monokai-input hover:bg-monokai-inputHover px-3 py-1.5 border border-monokai-comment rounded-lg text-monokai-yellow text-sm transition">
                    示例报文
                </button>
                <button @click="parseFrame"
                    class="bg-monokai-green hover:bg-monokai-greenHover px-3 py-1.5 rounded-lg text-monokai-bg text-sm transition">
                    立即解析
                </button>
                <button @click="clearAll"
                    class="bg-monokai-input hover:bg-monokai-inputHover px-3 py-1.5 border border-monokai-comment rounded-lg text-monokai-yellow text-sm transition">
                    清空
                </button>
            </div>
        </header>

        <!-- 输入区 + 字节预览 -->
        <div class="p-4 border-monokai-comment border-b shrink-0">
            <div class="mb-2 text-monokai-comment text-sm">原始报文（十六进制，支持空格分隔）：</div>
            <textarea v-model="rawHex" @input="onInputChange"
                class="bg-monokai-input p-3 border border-monokai-comment focus:border-monokai-cyan rounded-xl focus:outline-none w-full font-mono text-monokai-yellow text-sm resize-none"
                rows="2"
                placeholder="例如：68 1F 1F 68 08 02 02 01 01 00 01 C1 01 80 01 00 00 00 01 00 00 00 01 00 00 00 00 6A 35 16"></textarea>

            <!-- 字节预览区（按帧结构分区） -->
            <div v-if="frameBytes.length"
                class="flex flex-wrap gap-1.5 bg-monokai-input mt-3 p-3 border border-monokai-comment rounded-xl">
                <span v-for="(byte, idx) in frameBytes" :key="idx"
                    class="flex justify-center items-center border border-[#444] rounded w-10 h-8 font-mono text-sm"
                    :style="getByteBgStyle(idx)">
                    {{ byte }}
                </span>
            </div>

            <!-- 图例 -->
            <div class="flex flex-wrap gap-x-4 gap-y-2 mt-3 text-sm">
                <span class="text-monokai-comment">区块图例：</span>
                <div class="flex items-center gap-1"><span class="border border-monokai-comment rounded w-4 h-4"
                        style="background-color: #3b82f6;"></span><span>链路头部</span></div>
                <div class="flex items-center gap-1"><span class="border border-monokai-comment rounded w-4 h-4"
                        style="background-color: #ef4444;"></span><span>地址域</span></div>
                <div class="flex items-center gap-1"><span class="border border-monokai-comment rounded w-4 h-4"
                        style="background-color: #10b981;"></span><span>APDU (应用层)</span></div>
                <div class="flex items-center gap-1"><span class="border border-monokai-comment rounded w-4 h-4"
                        style="background-color: #fbbf24;"></span><span>校验域</span></div>
            </div>
        </div>

        <!-- 解析结果区域（使用卡片+树形结构） -->
        <div class="flex-1 p-4 overflow-y-auto custom-scrollbar">
            <!-- 解析状态提示 -->
            <div v-if="isParsing" class="mb-2 text-monokai-yellow text-sm">解析中...</div>

            <!-- 链路层卡片 -->
            <div v-if="parsed.link" class="bg-monokai-surface mb-4 p-4 border border-monokai-comment/40 rounded-xl">
                <div class="flex items-center gap-2 mb-3">
                    <div class="bg-[#3b82f6] rounded-full w-2 h-5"></div>
                    <h3 class="font-semibold text-monokai-cyan">链路层信息</h3>
                </div>
                <div class="gap-3 grid grid-cols-2 md:grid-cols-3 text-sm">
                    <div><span class="text-monokai-comment">起始符：</span>{{ parsed.link.start }}</div>
                    <div><span class="text-monokai-comment">长度：</span>{{ parsed.link.length }} 字节</div>
                    <div><span class="text-monokai-comment">控制域：</span>{{ parsed.link.control }}</div>
                    <div><span class="text-monokai-comment">地址域：</span>{{ parsed.link.address }}</div>
                    <div><span class="text-monokai-comment">帧头校验(HCS)：</span>{{ parsed.link.hcs }} <span
                            v-if="parsed.link.hcsValid !== undefined"
                            :class="parsed.link.hcsValid ? 'text-monokai-green' : 'text-monokai-pink'">({{ parsed.link.hcsValid
                            ? '通过' : '失败' }})</span></div>
                    <div><span class="text-monokai-comment">整帧校验(FCS)：</span>{{ parsed.link.fcs }} <span
                            v-if="parsed.link.fcsValid !== undefined"
                            :class="parsed.link.fcsValid ? 'text-monokai-green' : 'text-monokai-pink'">({{ parsed.link.fcsValid
                            ? '通过' : '失败' }})</span></div>
                    <div><span class="text-monokai-comment">结束符：</span>{{ parsed.link.end }}</div>
                </div>
            </div>

            <!-- 应用层卡片（APDU） -->
            <div v-if="parsed.apdu" class="bg-monokai-surface mb-4 p-4 border border-monokai-comment/40 rounded-xl">
                <div class="flex items-center gap-2 mb-3">
                    <div class="bg-[#10b981] rounded-full w-2 h-5"></div>
                    <h3 class="font-semibold text-monokai-green">应用层 (APDU)</h3>
                </div>
                <div class="text-sm">
                    <!-- APDU 类型 -->
                    <div class="mb-2"><span class="text-monokai-comment">APDU类型：</span>{{ parsed.apdu.typeName }} ({{ parsed.apdu.typeId }})</div>
                    <!-- 如果是请求/响应，显示 OAD 或数据集 -->
                    <div v-if="parsed.apdu.oads && parsed.apdu.oads.length" class="mt-3">
                        <div class="mb-1 text-monokai-comment">请求对象列表 (OAD)：</div>
                        <div class="space-y-2 ml-2">
                            <div v-for="(oad, idx) in parsed.apdu.oads" :key="idx"
                                class="bg-monokai-bg p-2 border border-monokai-comment/30 rounded">
                                <div><span class="text-monokai-comment">接口类(IC)：</span>{{ oad.ic }} ({{ oad.icName }})</div>
                                <div><span class="text-monokai-comment">对象标识(OI)：</span>{{ oad.oi }}</div>
                                <div><span class="text-monokai-comment">属性标识(PI)：</span>{{ oad.pi }} ({{ oad.piName }})</div>
                                <div v-if="oad.value"><span class="text-monokai-comment">值：</span>{{ oad.value }}</div>
                            </div>
                        </div>
                    </div>
                    <!-- 数据集响应 -->
                    <div v-if="parsed.apdu.dataSet && parsed.apdu.dataSet.length" class="mt-3">
                        <div class="mb-1 text-monokai-comment">响应数据：</div>
                        <div class="space-y-2 ml-2">
                            <div v-for="(item, idx) in parsed.apdu.dataSet" :key="idx"
                                class="bg-monokai-bg p-2 border border-monokai-comment/30 rounded">
                                <div><span class="text-monokai-comment">对象标识：</span>{{ item.oi }}</div>
                                <div><span class="text-monokai-comment">属性标识：</span>{{ item.pi }}</div>
                                <div><span class="text-monokai-comment">值：</span>{{ item.value }}</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 错误提示 -->
            <div v-if="errorMsg" class="bg-monokai-pink/20 p-3 border border-monokai-pink rounded-xl text-monokai-pink text-sm">
                ⚠️ 解析错误：{{ errorMsg }}
            </div>

            <!-- 提示信息：推荐使用后端解析 -->
            <div v-if="!errorMsg && !parsed.link && !parsed.apdu && rawHex" class="py-10 text-monokai-comment text-center">
                报文格式不正确或暂不支持自动解析。<br>
                建议通过后端 Go 服务进行完整解析。
            </div>
            <div v-if="!rawHex" class="py-10 text-monokai-comment text-center">
                请输入 DL/T 698.45 报文进行解析
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