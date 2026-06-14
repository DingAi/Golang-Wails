<template>
  <div class="flex flex-col bg-[#272822] w-full h-full overflow-hidden text-[#F8F8F2]">
    <header class="flex justify-between items-center bg-[#3E3D32] px-4 py-3 border-[#75715E] border-b shrink-0">
      <h2 class="font-medium text-[#66D9EF] text-lg">Modbus规则设计</h2>
      <div class="flex gap-3">
        <button @click="addField"
          class="bg-[#A6E22E] hover:bg-[#93c725] px-3 py-1.5 rounded-lg text-[#272822] text-sm transition">+
          添加字段分区</button>
        <button @click="toggleLayout"
          class="bg-[#32332E] hover:bg-[#4e4d40] px-3 py-1.5 border border-[#75715E] rounded-lg text-[#E6DB74] text-sm transition">{{
            layoutMode === 'horizontal' ? '切换为竖向' : '切换为横向' }}</button>
        <button @click="clearAll"
          class="bg-[#32332E] hover:bg-[#4e4d40] px-3 py-1.5 border border-[#75715E] rounded-lg text-[#E6DB74] text-sm transition">清空所有分区</button>
      </div>
    </header>

    <!-- 数据帧输入区 -->
    <div class="p-4 border-[#75715E] border-b shrink-0">
      <div class="mb-2 text-[#75715E] text-sm">数据帧（十六进制，支持空格分隔）：</div>
      <textarea v-model="frameInput" @input="parseFrame"
        class="bg-[#32332E] p-3 border border-[#75715E] focus:border-[#66D9EF] rounded-xl focus:outline-none w-full font-mono text-[#E6DB74] text-sm resize-none"
        rows="2" placeholder="例如：01 03 00 00 00 01 94 0B"></textarea>
      <div class="flex flex-wrap gap-1.5 bg-[#32332E] mt-3 p-3 border border-[#75715E] rounded-xl">
        <span v-for="(byte, idx) in frameBytes" :key="idx"
          class="flex justify-center items-center border border-[#444] rounded w-10 h-8 font-mono text-sm"
          :style="getByteBgStyle(idx)">{{ byte }}</span>
      </div>
      <div class="flex flex-wrap gap-x-4 gap-y-2 mt-3 text-sm">
        <span class="text-[#75715E]">分区图例：</span>
        <template v-for="(field, idx) in fieldList" :key="idx">
          <div class="flex items-center gap-1"><span class="border border-[#75715E] rounded w-4 h-4"
              :style="{ backgroundColor: field.color }"></span><span>{{ field.name || `未命名分区${idx + 1}` }}</span></div>
        </template>
      </div>
    </div>

    <!-- 横向布局（左右滑动） -->
    <div v-if="layoutMode === 'horizontal'" class="flex-1 p-4 overflow-x-auto overflow-y-hidden custom-scrollbar">
      <div class="flex flex-row items-stretch gap-3" style="min-width: min-content;">
        <div v-for="(field, idx) in fieldList" :key="idx"
          class="flex-shrink-0 bg-[#3E3D32] p-3 border border-[#75715E]/40 rounded-xl" style="width: 360px;">
          <!-- 卡片内容 -->
          <div>
            <div class="flex justify-between items-center mb-3">
              <span class="font-medium text-[#66D9EF] text-sm">分区 #{{ idx + 1 }}</span>
              <div class="flex gap-2">
                <button @click="moveUp(idx)" :disabled="idx === 0"
                  class="flex justify-center items-center bg-[#32332E] hover:bg-[#4e4d40] disabled:opacity-40 rounded w-6 h-6 text-xs">↑</button>
                <button @click="moveDown(idx)" :disabled="idx === fieldList.length - 1"
                  class="flex justify-center items-center bg-[#32332E] hover:bg-[#4e4d40] disabled:opacity-40 rounded w-6 h-6 text-xs">↓</button>
                <button @click="delField(idx)"
                  class="flex justify-center items-center bg-[#F92672]/20 hover:bg-[#F92672]/40 rounded w-6 h-6 text-[#F92672] text-xs">×</button>
              </div>
            </div>
            <div class="flex items-center gap-3 mb-3">
              <label class="w-20 text-[#75715E] text-sm">分区名称</label>
              <input v-model="field.name" placeholder="例：设备地址、功能码、数据、CRC校验"
                class="flex-1 bg-[#32332E] px-3 py-1.5 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none text-sm" />
              <label class="text-[#75715E] text-sm">标记色</label>
              <input v-model="field.color" type="color"
                class="bg-transparent p-0 border-0 rounded w-8 h-8 cursor-pointer" />
            </div>
            <div class="flex items-center gap-3 mb-3">
              <label class="w-20 text-[#75715E] text-sm">起始偏移(字节)</label>
              <input v-model.number="field.start" type="number" min="0"
                class="bg-[#32332E] px-3 py-1.5 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none w-24 text-sm" />
              <label class="text-[#75715E] text-sm">截取长度(字节)</label>
              <input v-model.number="field.len" type="number" min="1"
                class="bg-[#32332E] px-3 py-1.5 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none w-24 text-sm" />
              <span class="text-[#75715E] text-xs">字节从 0 开始计数</span>
            </div>
            <div class="bg-[#272822] mt-2 p-2 border border-[#75715E]/30 rounded">
              <div class="flex items-center gap-3 mb-2"><label
                  class="flex items-center gap-1 text-sm cursor-pointer"><input type="checkbox"
                    v-model="field.enableParse" class="w-3.5 h-3.5" /><span class="text-[#75715E]">启用数值解析</span></label>
              </div>
              <div v-if="field.enableParse" class="gap-3 grid grid-cols-2 text-sm">
                <div><label class="block mb-1 text-[#75715E]">数值类型</label><select v-model="field.parseType"
                    class="custom-select">
                    <option value="uint16">无符号整型16位</option>
                    <option value="int16">有符号整型16位</option>
                    <option value="uint32">无符号整型32位</option>
                    <option value="int32">有符号整型32位</option>
                    <option value="float32">单精度浮点数32位</option>
                  </select></div>
                <div><label class="block mb-1 text-[#75715E]">字节序</label><select v-model="field.endian"
                    class="custom-select">
                    <option value="big">大端 (Big-Endian)</option>
                    <option value="little">小端 (Little-Endian)</option>
                  </select></div>
              </div>
              <div v-if="field.enableParse" class="mt-2 pt-2 border-[#75715E]/30 border-t text-sm">
                <div class="flex justify-between items-center"><span class="text-[#75715E]">原始数据：</span><span
                    class="font-mono text-[#E6DB74]">{{ getRawHex(field) }}</span></div>
                <div class="flex justify-between items-center mt-1"><span class="text-[#75715E]">解析值：</span><span
                    :class="getParseStatus(field).class" class="font-bold">{{ getParseStatus(field).text }}</span></div>
                <div v-if="getParseWarning(field)" class="mt-1 text-[#F92672] text-xs">⚠️ {{ getParseWarning(field) }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-if="fieldList.length === 0" class="py-10 text-[#75715E] text-center">暂无分区规则</div>
    </div>

    <!-- 竖向布局（上下滚动） -->
    <div v-else class="flex-1 p-4 overflow-y-auto custom-scrollbar">
      <div class="gap-3 grid">
        <div v-for="(field, idx) in fieldList" :key="idx"
          class="bg-[#3E3D32] p-3 border border-[#75715E]/40 rounded-xl">
          <!-- 卡片内容与上面完全一致（复制粘贴） -->
          <div>
            <div class="flex justify-between items-center mb-3">
              <span class="font-medium text-[#66D9EF] text-sm">分区 #{{ idx + 1 }}</span>
              <div class="flex gap-2">
                <button @click="moveUp(idx)" :disabled="idx === 0"
                  class="flex justify-center items-center bg-[#32332E] hover:bg-[#4e4d40] disabled:opacity-40 rounded w-6 h-6 text-xs">↑</button>
                <button @click="moveDown(idx)" :disabled="idx === fieldList.length - 1"
                  class="flex justify-center items-center bg-[#32332E] hover:bg-[#4e4d40] disabled:opacity-40 rounded w-6 h-6 text-xs">↓</button>
                <button @click="delField(idx)"
                  class="flex justify-center items-center bg-[#F92672]/20 hover:bg-[#F92672]/40 rounded w-6 h-6 text-[#F92672] text-xs">×</button>
              </div>
            </div>
            <div class="flex items-center gap-3 mb-3">
              <label class="w-20 text-[#75715E] text-sm">分区名称</label>
              <input v-model="field.name" placeholder="例：设备地址、功能码、数据、CRC校验"
                class="flex-1 bg-[#32332E] px-3 py-1.5 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none text-sm" />
              <label class="text-[#75715E] text-sm">标记色</label>
              <input v-model="field.color" type="color"
                class="bg-transparent p-0 border-0 rounded w-8 h-8 cursor-pointer" />
            </div>
            <div class="flex items-center gap-3 mb-3">
              <label class="w-20 text-[#75715E] text-sm">起始偏移(字节)</label>
              <input v-model.number="field.start" type="number" min="0"
                class="bg-[#32332E] px-3 py-1.5 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none w-24 text-sm" />
              <label class="text-[#75715E] text-sm">截取长度(字节)</label>
              <input v-model.number="field.len" type="number" min="1"
                class="bg-[#32332E] px-3 py-1.5 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none w-24 text-sm" />
              <span class="text-[#75715E] text-xs">字节从 0 开始计数</span>
            </div>
            <div class="bg-[#272822] mt-2 p-2 border border-[#75715E]/30 rounded">
              <div class="flex items-center gap-3 mb-2"><label
                  class="flex items-center gap-1 text-sm cursor-pointer"><input type="checkbox"
                    v-model="field.enableParse" class="w-3.5 h-3.5" /><span class="text-[#75715E]">启用数值解析</span></label>
              </div>
              <div v-if="field.enableParse" class="gap-3 grid grid-cols-2 text-sm">
                <div><label class="block mb-1 text-[#75715E]">数值类型</label><select v-model="field.parseType"
                    class="custom-select">
                    <option value="uint16">无符号整型16位</option>
                    <option value="int16">有符号整型16位</option>
                    <option value="uint32">无符号整型32位</option>
                    <option value="int32">有符号整型32位</option>
                    <option value="float32">单精度浮点数32位</option>
                  </select></div>
                <div><label class="block mb-1 text-[#75715E]">字节序</label><select v-model="field.endian"
                    class="custom-select">
                    <option value="big">大端 (Big-Endian)</option>
                    <option value="little">小端 (Little-Endian)</option>
                  </select></div>
              </div>
              <div v-if="field.enableParse" class="mt-2 pt-2 border-[#75715E]/30 border-t text-sm">
                <div class="flex justify-between items-center"><span class="text-[#75715E]">原始数据：</span><span
                    class="font-mono text-[#E6DB74]">{{ getRawHex(field) }}</span></div>
                <div class="flex justify-between items-center mt-1"><span class="text-[#75715E]">解析值：</span><span
                    :class="getParseStatus(field).class" class="font-bold">{{ getParseStatus(field).text }}</span></div>
                <div v-if="getParseWarning(field)" class="mt-1 text-[#F92672] text-xs">⚠️ {{ getParseWarning(field) }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-if="fieldList.length === 0" class="py-10 text-[#75715E] text-center">暂无分区规则</div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const defaultFrameStr = "01 03 00 00 00 01 94 0B"
const frameInput = ref(defaultFrameStr)
const frameBytes = ref([])

const fieldList = ref([
  { name: "设备地址", start: 0, len: 1, color: "#ff6b6b", enableParse: false, parseType: "uint16", endian: "big" },
  { name: "功能码", start: 1, len: 1, color: "#4ecdc4", enableParse: false, parseType: "uint16", endian: "big" },
  { name: "数据域", start: 2, len: 4, color: "#ffe066", enableParse: true, parseType: "uint32", endian: "big" },
  { name: "CRC校验", start: 6, len: 2, color: "#a8e6cf", enableParse: true, parseType: "uint16", endian: "big" }
])

const layoutMode = ref('horizontal')  // 默认横向
const toggleLayout = () => { layoutMode.value = layoutMode.value === 'horizontal' ? 'vertical' : 'horizontal' }

const parseFrame = () => {
  const str = frameInput.value.replace(/\s/g, '')
  const bytes = []
  for (let i = 0; i < str.length; i += 2) {
    const byte = str.substr(i, 2).toUpperCase()
    if (/^[0-9A-F]{2}$/.test(byte)) bytes.push(byte)
  }
  frameBytes.value = bytes
}
parseFrame()

const getRawHex = (field) => {
  const start = field.start
  const end = start + field.len
  if (start >= frameBytes.value.length) return ''
  const slice = frameBytes.value.slice(start, Math.min(end, frameBytes.value.length))
  return slice.join(' ')
}

const parseValue = (field) => {
  if (!field.enableParse) return null
  const start = field.start
  const end = start + field.len
  if (start >= frameBytes.value.length) return null
  const rawBytes = frameBytes.value.slice(start, Math.min(end, frameBytes.value.length))
  if (rawBytes.length !== field.len) return null
  const buffer = new ArrayBuffer(field.len)
  const view = new DataView(buffer)
  for (let i = 0; i < rawBytes.length; i++) view.setUint8(i, parseInt(rawBytes[i], 16))
  const isLittle = field.endian === 'little'
  try {
    switch (field.parseType) {
      case 'uint16': return field.len === 2 ? view.getUint16(0, isLittle) : null
      case 'int16': return field.len === 2 ? view.getInt16(0, isLittle) : null
      case 'uint32': return field.len === 4 ? view.getUint32(0, isLittle) : null
      case 'int32': return field.len === 4 ? view.getInt32(0, isLittle) : null
      case 'float32': return field.len === 4 ? view.getFloat32(0, isLittle) : null
      default: return null
    }
  } catch { return null }
}

const getParseStatus = (field) => {
  if (!field.enableParse) return { text: '', class: '' }
  const requiredLen = (type) => {
    if (type === 'uint16' || type === 'int16') return 2
    if (type === 'uint32' || type === 'int32' || type === 'float32') return 4
    return 0
  }
  const req = requiredLen(field.parseType)
  if (field.len !== req) return { text: `长度不匹配（需要 ${req} 字节）`, class: 'text-[#F92672]' }
  const val = parseValue(field)
  if (val === null) return { text: '解析失败', class: 'text-[#F92672]' }
  return { text: val, class: 'text-[#A6E22E]' }
}

const getParseWarning = (field) => {
  if (!field.enableParse) return ''
  const requiredLen = (type) => {
    if (type === 'uint16' || type === 'int16') return 2
    if (type === 'uint32' || type === 'int32' || type === 'float32') return 4
    return 0
  }
  const req = requiredLen(field.parseType)
  if (field.len !== req) {
    const map = { uint16: '无符号16位整型', int16: '有符号16位整型', uint32: '无符号32位整型', int32: '有符号32位整型', float32: '单精度浮点数' }
    return `当前分区长度为 ${field.len} 字节，但 ${map[field.parseType]} 需要 ${req} 字节，无法解析`
  }
  if (field.start + field.len > frameBytes.value.length) return '数据帧长度不足，无法截取该分区数据'
  return ''
}

const getByteBgStyle = (byteIndex) => {
  let bg = "transparent"
  for (let i = fieldList.value.length - 1; i >= 0; i--) {
    const f = fieldList.value[i]
    const s = f.start
    const e = f.start + f.len - 1
    if (byteIndex >= s && byteIndex <= e) { bg = f.color; break }
  }
  return { backgroundColor: bg }
}

const addField = () => { fieldList.value.push({ name: "", start: 0, len: 1, color: "#9d65c9", enableParse: false, parseType: "uint16", endian: "big" }) }
const delField = (idx) => fieldList.value.splice(idx, 1)
const moveUp = (idx) => { if (idx <= 0) return; const temp = fieldList.value[idx]; fieldList.value[idx] = fieldList.value[idx - 1]; fieldList.value[idx - 1] = temp }
const moveDown = (idx) => { if (idx >= fieldList.value.length - 1) return; const temp = fieldList.value[idx]; fieldList.value[idx] = fieldList.value[idx + 1]; fieldList.value[idx + 1] = temp }
const clearAll = () => { fieldList.value = [] }
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

.custom-select {
  background-color: #32332E;
  border: 1px solid #75715E;
  border-radius: 0.5rem;
  padding: 0.25rem 0.5rem;
  color: #F8F8F2;
  font-size: 0.875rem;
  outline: none;
  cursor: pointer;
  width: 100%;
}

.custom-select:focus {
  border-color: #66D9EF;
}

.custom-select option {
  background-color: #3E3D32;
  color: #F8F8F2;
}
</style>