<template>
  <div class="flex flex-col bg-monokai-bg w-full h-full overflow-hidden text-monokai-foreground">
    <!-- 顶部标题栏 -->
    <header
      class="flex justify-between items-center bg-monokai-surface px-4 py-3 border-monokai-comment border-b shrink-0">
      <h2 class="font-medium text-monokai-cyan text-lg">IEC104 协议解析工具</h2>
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

    <!-- 输入区 + 字节预览（类似 Modbus） -->
    <div class="p-4 border-monokai-comment border-b shrink-0">
      <div class="mb-2 text-monokai-comment text-sm">原始报文（十六进制，支持空格分隔）：</div>
      <textarea v-model="rawHex" @input="onInputChange"
        class="bg-monokai-input p-3 border border-monokai-comment focus:border-monokai-cyan rounded-xl focus:outline-none w-full font-mono text-monokai-yellow text-sm resize-none"
        rows="2" placeholder="例如：68 12 00 00 00 00 02 01 00 01 00 00 01 00 01 00 00 00 00 79"></textarea>

      <!-- 字节预览区（独立显示每个字节，颜色分区） -->
      <div v-if="frameBytes.length"
        class="flex flex-wrap gap-1.5 bg-monokai-input mt-3 p-3 border border-monokai-comment rounded-xl">
        <span v-for="(byte, idx) in frameBytes" :key="idx"
          class="flex justify-center items-center border border-[#444] rounded w-10 h-8 font-mono text-sm"
          :style="getByteBgStyle(idx)">
          {{ byte }}
        </span>
      </div>

      <!-- 图例说明 -->
      <div class="flex flex-wrap gap-x-4 gap-y-2 mt-3 text-sm">
        <span class="text-monokai-comment">区块图例：</span>
        <div class="flex items-center gap-1"><span class="border border-monokai-comment rounded w-4 h-4"
            style="background-color: #3b82f6;"></span><span>APCI头部</span></div>
        <div class="flex items-center gap-1"><span class="border border-monokai-comment rounded w-4 h-4"
            style="background-color: #ef4444;"></span><span>ASDU头部</span></div>
        <div class="flex items-center gap-1"><span class="border border-monokai-comment rounded w-4 h-4"
            style="background-color: #fbbf24;"></span><span>信息体地址</span></div>
        <div class="flex items-center gap-1"><span class="border border-monokai-comment rounded w-4 h-4"
            style="background-color: #10b981;"></span><span>信息体数据</span></div>
      </div>
    </div>

    <!-- 解析结果区域（卡片式） -->
    <div class="flex-1 p-4 overflow-y-auto custom-scrollbar">
      <!-- APCI 卡片 -->
      <div v-if="parsed.apci" class="bg-monokai-surface mb-4 p-4 border border-monokai-comment/40 rounded-xl">
        <div class="flex items-center gap-2 mb-3">
          <div class="bg-[#3b82f6] rounded-full w-2 h-5"></div>
          <h3 class="font-semibold text-monokai-cyan">APCI (应用规约控制信息)</h3>
        </div>
        <div class="gap-3 grid grid-cols-2 md:grid-cols-4 text-sm">
          <div><span class="text-monokai-comment">帧类型：</span>{{ parsed.apci.frameType }}</div>
          <div v-if="parsed.apci.sendSeq !== undefined"><span class="text-monokai-comment">发送序号：</span>{{
            parsed.apci.sendSeq
            }}</div>
          <div v-if="parsed.apci.recvSeq !== undefined"><span class="text-monokai-comment">接收序号：</span>{{
            parsed.apci.recvSeq
            }}</div>
          <div v-if="parsed.apci.cause"><span class="text-monokai-comment">原因：</span>{{ parsed.apci.cause }}</div>
        </div>
      </div>

      <!-- ASDU 卡片 -->
      <div v-if="parsed.asdu" class="bg-monokai-surface mb-4 p-4 border border-monokai-comment/40 rounded-xl">
        <div class="flex items-center gap-2 mb-3">
          <div class="bg-[#ef4444] rounded-full w-2 h-5"></div>
          <h3 class="font-semibold text-monokai-pink">ASDU (应用服务数据单元)</h3>
        </div>
        <div class="gap-3 grid grid-cols-2 md:grid-cols-3 mb-4 text-sm">
          <div><span class="text-monokai-comment">类型标识：</span>{{ parsed.asdu.typeId }} ({{ parsed.asdu.typeName }})
          </div>
          <div><span class="text-monokai-comment">可变结构限定词：</span>{{ parsed.asdu.vsq }} (个数={{ parsed.asdu.infoCount }},
            顺序={{
              parsed.asdu.isSequence ? '顺序' : '非顺序' }})</div>
          <div><span class="text-monokai-comment">传送原因：</span>{{ parsed.asdu.causeTx }} ({{ parsed.asdu.causeName }})
          </div>
          <div><span class="text-monokai-comment">原发地址：</span>{{ parsed.asdu.originAddr }}</div>
          <div><span class="text-monokai-comment">公共地址：</span>{{ parsed.asdu.commonAddr }}</div>
        </div>

        <!-- 信息体列表（卡片嵌套） -->
        <div v-if="parsed.asdu.infos && parsed.asdu.infos.length">
          <div class="mb-2 text-monokai-comment text-sm">信息体 (共 {{ parsed.asdu.infos.length }} 个)</div>
          <div class="space-y-3">
            <div v-for="(info, idx) in parsed.asdu.infos" :key="idx"
              class="bg-monokai-input p-3 border border-monokai-comment/30 rounded-lg">
              <div class="flex items-center gap-2 mb-2">
                <div class="bg-[#fbbf24] rounded-full w-2 h-4"></div>
                <span class="font-medium text-monokai-yellow">信息体 #{{ idx + 1 }}</span>
              </div>
              <div class="gap-2 grid grid-cols-1 md:grid-cols-3 text-sm">
                <div><span class="text-monokai-comment">信息体地址：</span>{{ info.address }}</div>
                <div><span class="text-monokai-comment">原始数据：</span><span class="font-mono text-monokai-green">{{
                  info.dataHex
                    }}</span></div>
                <div><span class="text-monokai-comment">解析值：</span>{{ info.interpreted }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 错误提示 -->
      <div v-if="errorMsg"
        class="bg-monokai-pink/20 p-3 border border-monokai-pink rounded-xl text-monokai-pink text-sm">
        ⚠️ {{ errorMsg }}
      </div>

      <!-- 无数据提示 -->
      <div v-if="!parsed.apci && !parsed.asdu && !errorMsg && rawHex" class="py-10 text-monokai-comment text-center">
        等待解析...
      </div>
      <div v-if="!rawHex" class="py-10 text-monokai-comment text-center">
        请输入 IEC104 报文进行解析
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'

const rawHex = ref('')
const parsed = ref({ apci: null, asdu: null })
const errorMsg = ref('')
const frameBytes = ref([])          // 存储字节数组用于预览
let highlightRanges = reactive([])  // 存储各区块的字节索引范围

// 清空所有
function clearAll() {
  rawHex.value = ''
  parsed.value = { apci: null, asdu: null }
  errorMsg.value = ''
  frameBytes.value = []
  highlightRanges = []
}

// 填充示例报文（单点信息遥测）
function fillExample() {
  rawHex.value = '68 12 00 00 00 00 02 01 00 01 00 00 01 00 01 00 00 00 00 79'
  onInputChange()
}

function onInputChange() {
  if (!rawHex.value.trim()) {
    clearAll()
  } else {
    updateFrameBytes()
    parseFrame()
  }
}

// 更新字节数组和颜色分区
function updateFrameBytes() {
  const hexStr = rawHex.value.replace(/\s/g, '')
  const bytes = []
  for (let i = 0; i < hexStr.length; i += 2) {
    const byte = hexStr.substr(i, 2).toUpperCase()
    if (/^[0-9A-F]{2}$/.test(byte)) {
      bytes.push(byte)
    }
  }
  frameBytes.value = bytes
  // 计算高亮区域（基于最新解析结果，但解析尚未执行，因此需要解析后重新调用）
  // 实际在 parseFrame 成功后重新调用 computeHighlights
}

// 根据字节索引获取背景色
function getByteBgStyle(idx) {
  if (!highlightRanges.length) return { backgroundColor: 'transparent' }
  for (const range of highlightRanges) {
    if (idx >= range.start && idx <= range.end) {
      return { backgroundColor: range.color }
    }
  }
  return { backgroundColor: 'transparent' }
}

// 计算字节高亮区域（基于解析过程中得到的位置信息）
function computeHighlights(apciLen, asduStart, infoPositions) {
  const ranges = []
  // APCI: 0 ~ apciLen-1（通常6字节，但实际是控制域4+起始2? 实际上字节0-1是启动和长度，2-5是控制域，所以APCI是0-5）
  ranges.push({ start: 0, end: 5, color: '#3b82f6' })
  if (asduStart !== undefined) {
    // ASDU头部：从 asduStart 到第一个信息体地址之前（6字节头部）
    const asduHeaderEnd = asduStart + 6 - 1   // 类型、VSQ、传送原因、原发地址、公共地址(2字节) -> 6字节
    ranges.push({ start: asduStart, end: asduHeaderEnd, color: '#ef4444' })
    // 信息体地址和数据
    for (const pos of infoPositions) {
      // 地址3字节
      ranges.push({ start: pos.addrStart, end: pos.addrStart + 2, color: '#fbbf24' })
      // 数据部分
      if (pos.dataStart !== undefined && pos.dataLen > 0) {
        ranges.push({ start: pos.dataStart, end: pos.dataStart + pos.dataLen - 1, color: '#10b981' })
      }
    }
  }
  highlightRanges = ranges
}

// 主解析函数（增强版，返回更详细的位置信息）
function parseFrame() {
  errorMsg.value = ''
  parsed.value = { apci: null, asdu: null }
  if (!frameBytes.value.length) return

  const bytes = frameBytes.value.map(b => parseInt(b, 16))
  try {
    if (bytes.length < 6) throw new Error('报文太短，无法解析 APCI')
    if (bytes[0] !== 0x68) throw new Error('报文起始字节不是 68H')
    const apduLen = bytes[1]
    if (bytes.length < apduLen + 2) throw new Error(`报文长度不足`)

    // APCI 解析
    const ctrl1 = bytes[2], ctrl2 = bytes[3], ctrl3 = bytes[4], ctrl4 = bytes[5]
    const frameType = getFrameType(ctrl1, ctrl2)
    const apci = { frameType }
    if (frameType === 'I-Format') {
      apci.sendSeq = ((ctrl2 & 0xFE) >> 1) | ((ctrl1 & 0x01) << 7)
      apci.recvSeq = ((ctrl4 & 0xFE) >> 1) | ((ctrl3 & 0x01) << 7)
    } else if (frameType === 'S-Format') {
      apci.recvSeq = ((ctrl4 & 0xFE) >> 1) | ((ctrl3 & 0x01) << 7)
    } else if (frameType === 'U-Format') {
      if (ctrl1 === 0x07 && ctrl2 === 0x00) apci.cause = '启动帧 (STARTDT)'
      else if (ctrl1 === 0x0B && ctrl2 === 0x00) apci.cause = '停止帧 (STOPDT)'
      else if (ctrl1 === 0x03 && ctrl2 === 0x00) apci.cause = '测试帧 (TESTFR)'
      else apci.cause = '未知 U 帧命令'
    }
    parsed.value.apci = apci

    // 如果是 I 帧且有 ASDU
    if (frameType !== 'I-Format' || bytes.length <= 6) {
      computeHighlights(6, undefined, [])
      return
    }

    let pos = 6
    const typeId = bytes[pos++]
    const vsq = bytes[pos++]
    const causeTx = bytes[pos++]
    const originAddr = bytes[pos++]   // 原发地址低字节（通常为0）
    const commonAddrLow = bytes[pos++]
    const commonAddrHigh = bytes[pos++]
    const commonAddr = (commonAddrHigh << 8) | commonAddrLow

    const isSequence = (vsq & 0x80) !== 0
    const infoCount = vsq & 0x7F

    const asdu = {
      typeId: `0x${typeId.toString(16).toUpperCase().padStart(2, '0')}`,
      typeName: getTypeName(typeId),
      vsq: `0x${vsq.toString(16).toUpperCase().padStart(2, '0')}`,
      infoCount,
      isSequence,
      causeTx: `0x${causeTx.toString(16).toUpperCase().padStart(2, '0')}`,
      causeName: getCauseName(causeTx),
      originAddr,
      commonAddr,
      infos: []
    }

    // 解析信息体，同时记录位置用于高亮
    const infoPositions = []
    for (let i = 0; i < infoCount; i++) {
      if (pos + 2 >= bytes.length) break
      const addrStart = pos
      const addrLow = bytes[pos++]
      const addrMid = bytes[pos++]
      const addrHigh = bytes[pos++]
      const address = (addrHigh << 16) | (addrMid << 8) | addrLow
      const dataLen = getInfoElementLength(typeId)
      const dataStart = pos
      if (pos + dataLen > bytes.length) {
        errorMsg.value = `信息体 ${i + 1} 数据不足`
        break
      }
      const dataBytes = bytes.slice(pos, pos + dataLen)
      pos += dataLen
      const dataHex = dataBytes.map(b => b.toString(16).toUpperCase().padStart(2, '0')).join(' ')
      // 智能解析值
      let interpreted = interpretInfoElement(typeId, dataBytes)
      asdu.infos.push({ address, dataHex, interpreted })
      infoPositions.push({ addrStart, dataStart, dataLen })
    }

    parsed.value.asdu = asdu
    // 计算高亮：APCI 0-5，ASDU头部 6-11，信息体地址和数据
    computeHighlights(6, 6, infoPositions.map((p, idx) => ({
      addrStart: p.addrStart,
      dataStart: p.dataStart,
      dataLen: p.dataLen
    })))
  } catch (e) {
    errorMsg.value = e.message
    computeHighlights(6, undefined, [])
  }
}

// 帧类型判断
function getFrameType(ctrl1, ctrl2) {
  if ((ctrl1 & 0x01) === 0 && (ctrl2 & 0x01) === 0) return 'I-Format'
  if ((ctrl1 & 0x01) === 1 && (ctrl2 & 0x01) === 0 && (ctrl1 & 0x02) === 0) return 'S-Format'
  return 'U-Format'
}

// 类型标识名称映射（与之前相同，可扩展）
function getTypeName(typeId) {
  const map = {
    0x01: 'M_SP_NA_1 (单点信息)',
    0x02: 'M_SP_TA_1 (单点时标)',
    0x03: 'M_DP_NA_1 (双点信息)',
    0x04: 'M_DP_TA_1 (双点时标)',
    0x05: 'M_ST_NA_1 (步位置)',
    0x06: 'M_ST_TA_1 (步位置时标)',
    0x07: 'M_BO_NA_1 (32比特串)',
    0x08: 'M_BO_TA_1 (32比特串时标)',
    0x09: 'M_ME_NA_1 (测量值，归一化值)',
    0x0A: 'M_ME_TA_1 (测量值，归一化值时标)',
    0x0B: 'M_ME_NB_1 (测量值，标度化值)',
    0x0C: 'M_ME_TB_1 (测量值，标度化值时标)',
    0x0D: 'M_ME_NC_1 (测量值，短浮点数)',
    0x0E: 'M_ME_TC_1 (测量值，短浮点数时标)',
    0x0F: 'M_IT_NA_1 (累计量)',
    0x10: 'M_IT_TA_1 (累计量时标)',
    0x1E: 'C_SC_NA_1 (单命令)',
    0x1F: 'C_DC_NA_1 (双命令)',
    0x20: 'C_RC_NA_1 (步调节命令)',
    0x21: 'C_SE_NA_1 (设定值，归一化值)',
    0x22: 'C_SE_NB_1 (设定值，标度化值)',
    0x23: 'C_SE_NC_1 (设定值，短浮点数)'
  }
  return map[typeId] || `未知类型 (0x${typeId.toString(16)})`
}

// 传送原因名称
function getCauseName(cause) {
  const map = {
    1: '周期、循环',
    2: '背景扫描',
    3: '突发、自发',
    4: '初始化',
    5: '请求或被召唤',
    6: '激活',
    7: '激活确认',
    8: '停止激活',
    9: '停止激活确认',
    10: '激活终止',
    20: '响应站召唤',
    21: '响应组召唤'
  }
  return map[cause] || '保留/未知'
}

// 根据类型标识获取信息体元素长度（字节）
function getInfoElementLength(typeId) {
  // 单点、双点、步位置、比特串、命令等：1字节
  if ([0x01, 0x03, 0x05, 0x07, 0x1E, 0x1F, 0x20].includes(typeId)) return 1
  // 测量值归一化/标度化、设定值归一化/标度化：2字节
  if ([0x09, 0x0B, 0x21, 0x22].includes(typeId)) return 2
  // 短浮点数：4字节
  if ([0x0D, 0x23].includes(typeId)) return 4
  // 带时标的类型：数据部分4字节（实际标准中时标7字节，但用户简化处理为4字节）
  if ([0x02, 0x04, 0x06, 0x08, 0x0A, 0x0C, 0x0E, 0x0F, 0x10].includes(typeId)) return 4
  return 1
}

// 智能解析信息体元素（根据类型标识翻译）
function interpretInfoElement(typeId, dataBytes) {
  if (!dataBytes.length) return '无数据'
  const len = dataBytes.length
  const view = new DataView(new Uint8Array(dataBytes).buffer)
  switch (typeId) {
    case 0x01: // 单点信息
      return dataBytes[0] === 0 ? 'OFF/分' : 'ON/合'
    case 0x03: // 双点信息
      const val = dataBytes[0] & 0x03
      if (val === 0) return '不确定/中间状态'
      if (val === 1) return '确定/合'
      if (val === 2) return '确定/分'
      return '不确定'
    case 0x09: // 归一化测量值
      if (len >= 2) return (view.getInt16(0, true) / 32767).toFixed(4)
      break
    case 0x0B: // 标度化测量值
      if (len >= 2) return view.getInt16(0, true).toString()
      break
    case 0x0D: // 短浮点数
      if (len >= 4) return view.getFloat32(0, true).toFixed(6)
      break
    default:
      // 默认显示数值（小端）
      if (len === 1) return dataBytes[0].toString()
      if (len === 2) return (dataBytes[0] | (dataBytes[1] << 8)).toString()
      if (len === 4) {
        const intVal = dataBytes[0] | (dataBytes[1] << 8) | (dataBytes[2] << 16) | (dataBytes[3] << 24)
        return intVal.toString()
      }
  }
  return dataBytes.map(b => b.toString(16).padStart(2, '0')).join(' ')
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