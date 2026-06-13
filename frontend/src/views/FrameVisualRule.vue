<template>
  <div class="flex flex-col bg-[#272822] w-full h-full overflow-hidden text-[#F8F8F2]">
    <!-- 顶部标题栏 -->
    <header class="flex justify-between items-center bg-[#3E3D32] px-4 py-3 border-[#75715E] border-b shrink-0">
      <h2 class="font-medium text-[#66D9EF] text-lg">数据帧翻译规则设计</h2>
      <div class="flex gap-3">
        <button
          @click="addField"
          class="bg-[#A6E22E] hover:bg-[#93c725] px-3 py-1.5 rounded-lg text-[#272822] text-sm transition"
        >
          + 添加字段分区
        </button>
        <button
          @click="clearAll"
          class="bg-[#32332E] hover:bg-[#4e4d40] px-3 py-1.5 border border-[#75715E] rounded-lg text-[#E6DB74] text-sm transition"
        >
          清空所有分区
        </button>
      </div>
    </header>

    <!-- 1. 上方：数据帧输入 + 预览高亮 -->
    <div class="p-4 border-[#75715E] border-b shrink-0">
      <div class="mb-2 text-[#75715E] text-sm">
        数据帧（十六进制，支持空格分隔）：
      </div>

      <!-- 可编辑输入框 -->
      <textarea
        v-model="frameInput"
        @input="parseFrame"
        class="bg-[#32332E] p-3 border border-[#75715E] focus:border-[#66D9EF] rounded-xl focus:outline-none w-full font-mono text-[#E6DB74] text-sm resize-none"
        rows="2"
        placeholder="例如：01 03 00 00 00 01 94 0B"
      ></textarea>

      <!-- 帧字节预览区，每个字节独立span，背景色分区 -->
      <div class="flex flex-wrap gap-1.5 bg-[#32332E] mt-3 p-3 border border-[#75715E] rounded-xl">
        <span
          v-for="(byte, idx) in frameBytes"
          :key="idx"
          class="flex justify-center items-center border border-[#444] rounded w-10 h-8 font-mono text-sm"
          :style="getByteBgStyle(idx)"
        >
          {{ byte }}
        </span>
      </div>

      <!-- 图例 -->
      <div class="flex flex-wrap gap-x-4 gap-y-2 mt-3 text-sm">
        <span class="text-[#75715E]">分区图例：</span>
        <template v-for="(field, idx) in fieldList" :key="idx">
          <div class="flex items-center gap-1">
            <span
              class="border border-[#75715E] rounded w-4 h-4"
              :style="{ backgroundColor: field.color }"
            ></span>
            <span>{{ field.name || `未命名分区${idx+1}` }}</span>
          </div>
        </template>
      </div>
    </div>

    <!-- 2. 下方：字段分区配置列表 -->
    <div class="flex-1 p-4 overflow-y-auto custom-scrollbar">
      <div class="gap-3 grid">
        <!-- 单条分区配置项 -->
        <div
          v-for="(field, idx) in fieldList"
          :key="idx"
          class="bg-[#3E3D32] p-3 border border-[#75715E]/40 rounded-xl"
        >
          <div class="flex justify-between items-center mb-3">
            <span class="font-medium text-[#66D9EF] text-sm">分区 #{{ idx + 1 }}</span>
            <div class="flex gap-2">
              <button
                @click="moveUp(idx)"
                :disabled="idx === 0"
                class="flex justify-center items-center bg-[#32332E] hover:bg-[#4e4d40] disabled:opacity-40 rounded w-6 h-6 text-xs"
              >↑</button>
              <button
                @click="moveDown(idx)"
                :disabled="idx === fieldList.length - 1"
                class="flex justify-center items-center bg-[#32332E] hover:bg-[#4e4d40] disabled:opacity-40 rounded w-6 h-6 text-xs"
              >↓</button>
              <button
                @click="delField(idx)"
                class="flex justify-center items-center bg-[#F92672]/20 hover:bg-[#F92672]/40 rounded w-6 h-6 text-[#F92672] text-xs"
              >×</button>
            </div>
          </div>

          <!-- 行1：分区名称 + 颜色选择 -->
          <div class="flex items-center gap-3 mb-3">
            <label class="w-20 text-[#75715E] text-sm">分区名称</label>
            <input
              v-model="field.name"
              placeholder="例：设备地址、功能码、数据、CRC校验"
              class="flex-1 bg-[#32332E] px-3 py-1.5 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none text-sm"
            />
            <label class="text-[#75715E] text-sm">标记色</label>
            <input
              v-model="field.color"
              type="color"
              class="bg-transparent p-0 border-0 rounded w-8 h-8 cursor-pointer"
            />
          </div>

          <!-- 行2：起始偏移 + 字节长度 -->
          <div class="flex items-center gap-3">
            <label class="w-20 text-[#75715E] text-sm">起始偏移(字节)</label>
            <input
              v-model.number="field.start"
              type="number"
              min="0"
              class="bg-[#32332E] px-3 py-1.5 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none w-24 text-sm"
            />

            <label class="text-[#75715E] text-sm">截取长度(字节)</label>
            <input
              v-model.number="field.len"
              type="number"
              min="1"
              class="bg-[#32332E] px-3 py-1.5 border border-[#75715E] focus:border-[#66D9EF] rounded-lg focus:outline-none w-24 text-sm"
            />
            <span class="text-[#75715E] text-xs">字节从 0 开始计数</span>
          </div>
        </div>
      </div>

      <!-- 空提示 -->
      <div v-if="fieldList.length === 0" class="py-10 text-[#75715E] text-center">
        暂无分区规则，点击「添加字段分区」开始设计数据帧规则
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

// 默认示例 Modbus RTU 数据帧
const defaultFrameStr = "01 03 00 00 00 01 94 0B"
const frameInput = ref(defaultFrameStr)
const frameBytes = ref([])

// 字段分区配置项结构
const fieldList = ref([
  { name: "设备地址", start: 0, len: 1, color: "#ff6b6b" },
  { name: "功能码", start: 1, len: 1, color: "#4ecdc4" },
  { name: "数据域", start: 2, len: 4, color: "#ffe066" },
  { name: "CRC校验", start: 6, len: 2, color: "#a8e6cf" }
])

// 解析输入的十六进制帧为字节数组
const parseFrame = () => {
  // 去掉所有空格，按2个字符一组拆分
  const str = frameInput.value.replace(/\s/g, '')
  const bytes = []
  for (let i = 0; i < str.length; i += 2) {
    const byte = str.substr(i, 2).toUpperCase()
    if (/^[0-9A-F]{2}$/.test(byte)) {
      bytes.push(byte)
    }
  }
  frameBytes.value = bytes
}

// 初始化解析
parseFrame()

// 根据字节下标，匹配对应分区背景色
const getByteBgStyle = (byteIndex) => {
  let bg = "transparent"
  // 倒序匹配：后添加的分区覆盖前面（防止区间重叠混乱）
  for (let i = fieldList.value.length - 1; i >= 0; i--) {
    const f = fieldList.value[i]
    const s = f.start
    const e = f.start + f.len - 1
    if (byteIndex >= s && byteIndex <= e) {
      bg = f.color
      break
    }
  }
  return { backgroundColor: bg }
}

// 新增分区
const addField = () => {
  fieldList.value.push({
    name: "",
    start: 0,
    len: 1,
    color: "#9d65c9"
  })
}

// 删除分区
const delField = (idx) => {
  fieldList.value.splice(idx, 1)
}

// 上移
const moveUp = (idx) => {
  if (idx <= 0) return
  const temp = fieldList.value[idx]
  fieldList.value[idx] = fieldList.value[idx - 1]
  fieldList.value[idx - 1] = temp
}

// 下移
const moveDown = (idx) => {
  if (idx >= fieldList.value.length - 1) return
  const temp = fieldList.value[idx]
  fieldList.value[idx] = fieldList.value[idx + 1]
  fieldList.value[idx + 1] = temp
}

// 清空全部
const clearAll = () => {
  fieldList.value = []
}
</script>

<style scoped>
/* 自定义滚动条 同项目风格 */
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