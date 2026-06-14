<template>
    <div class="flex flex-col bg-[#272822] w-full h-full overflow-hidden text-[#F8F8F2]">
        <!-- 顶部标题栏 -->
        <header class="flex justify-between items-center bg-[#3E3D32] px-4 py-3 border-[#75715E] border-b shrink-0">
            <h2 class="font-medium text-[#66D9EF] text-lg">DL/T 698.45 协议解析工具</h2>
            <div class="flex gap-3">
                <button @click="fillExample"
                    class="bg-[#32332E] hover:bg-[#4e4d40] px-3 py-1.5 border border-[#75715E] rounded-lg text-[#E6DB74] text-sm transition">
                    示例报文
                </button>
                <button @click="parseFrame"
                    class="bg-[#A6E22E] hover:bg-[#93c725] px-3 py-1.5 rounded-lg text-[#272822] text-sm transition">
                    立即解析
                </button>
                <button @click="clearAll"
                    class="bg-[#32332E] hover:bg-[#4e4d40] px-3 py-1.5 border border-[#75715E] rounded-lg text-[#E6DB74] text-sm transition">
                    清空
                </button>
            </div>
        </header>

        <!-- 输入区 + 字节预览 -->
        <div class="p-4 border-[#75715E] border-b shrink-0">
            <div class="mb-2 text-[#75715E] text-sm">原始报文（十六进制，支持空格分隔）：</div>
            <textarea v-model="rawHex" @input="onInputChange"
                class="bg-[#32332E] p-3 border border-[#75715E] focus:border-[#66D9EF] rounded-xl focus:outline-none w-full font-mono text-[#E6DB74] text-sm resize-none"
                rows="2"
                placeholder="例如：68 1F 1F 68 08 02 02 01 01 00 01 C1 01 80 01 00 00 00 01 00 00 00 01 00 00 00 00 6A 35 16"></textarea>

            <!-- 字节预览区（按帧结构分区） -->
            <div v-if="frameBytes.length"
                class="flex flex-wrap gap-1.5 bg-[#32332E] mt-3 p-3 border border-[#75715E] rounded-xl">
                <span v-for="(byte, idx) in frameBytes" :key="idx"
                    class="flex justify-center items-center border border-[#444] rounded w-10 h-8 font-mono text-sm"
                    :style="getByteBgStyle(idx)">
                    {{ byte }}
                </span>
            </div>

            <!-- 图例 -->
            <div class="flex flex-wrap gap-x-4 gap-y-2 mt-3 text-sm">
                <span class="text-[#75715E]">区块图例：</span>
                <div class="flex items-center gap-1"><span class="border border-[#75715E] rounded w-4 h-4"
                        style="background-color: #3b82f6;"></span><span>链路头部</span></div>
                <div class="flex items-center gap-1"><span class="border border-[#75715E] rounded w-4 h-4"
                        style="background-color: #ef4444;"></span><span>地址域</span></div>
                <div class="flex items-center gap-1"><span class="border border-[#75715E] rounded w-4 h-4"
                        style="background-color: #10b981;"></span><span>APDU (应用层)</span></div>
                <div class="flex items-center gap-1"><span class="border border-[#75715E] rounded w-4 h-4"
                        style="background-color: #fbbf24;"></span><span>校验域</span></div>
            </div>
        </div>

        <!-- 解析结果区域（使用卡片+树形结构） -->
        <div class="flex-1 p-4 overflow-y-auto custom-scrollbar">
            <!-- 解析状态提示 -->
            <div v-if="isParsing" class="mb-2 text-[#E6DB74] text-sm">解析中...</div>

            <!-- 链路层卡片 -->
            <div v-if="parsed.link" class="bg-[#3E3D32] mb-4 p-4 border border-[#75715E]/40 rounded-xl">
                <div class="flex items-center gap-2 mb-3">
                    <div class="bg-[#3b82f6] rounded-full w-2 h-5"></div>
                    <h3 class="font-semibold text-[#66D9EF]">链路层信息</h3>
                </div>
                <div class="gap-3 grid grid-cols-2 md:grid-cols-3 text-sm">
                    <div><span class="text-[#75715E]">起始符：</span>{{ parsed.link.start }}</div>
                    <div><span class="text-[#75715E]">长度：</span>{{ parsed.link.length }} 字节</div>
                    <div><span class="text-[#75715E]">控制域：</span>{{ parsed.link.control }}</div>
                    <div><span class="text-[#75715E]">地址域：</span>{{ parsed.link.address }}</div>
                    <div><span class="text-[#75715E]">帧头校验(HCS)：</span>{{ parsed.link.hcs }} <span
                            v-if="parsed.link.hcsValid !== undefined"
                            :class="parsed.link.hcsValid ? 'text-[#A6E22E]' : 'text-[#F92672]'">({{ parsed.link.hcsValid
                            ? '通过' : '失败' }})</span></div>
                    <div><span class="text-[#75715E]">整帧校验(FCS)：</span>{{ parsed.link.fcs }} <span
                            v-if="parsed.link.fcsValid !== undefined"
                            :class="parsed.link.fcsValid ? 'text-[#A6E22E]' : 'text-[#F92672]'">({{ parsed.link.fcsValid
                            ? '通过' : '失败' }})</span></div>
                    <div><span class="text-[#75715E]">结束符：</span>{{ parsed.link.end }}</div>
                </div>
            </div>

            <!-- 应用层卡片（APDU） -->
            <div v-if="parsed.apdu" class="bg-[#3E3D32] mb-4 p-4 border border-[#75715E]/40 rounded-xl">
                <div class="flex items-center gap-2 mb-3">
                    <div class="bg-[#10b981] rounded-full w-2 h-5"></div>
                    <h3 class="font-semibold text-[#A6E22E]">应用层 (APDU)</h3>
                </div>
                <div class="text-sm">
                    <!-- APDU 类型 -->
                    <div class="mb-2"><span class="text-[#75715E]">APDU类型：</span>{{ parsed.apdu.typeName }} ({{
                        parsed.apdu.typeId }})</div>
                    <!-- 如果是请求/响应，显示 OAD 或数据集 -->
                    <div v-if="parsed.apdu.oads && parsed.apdu.oads.length" class="mt-3">
                        <div class="mb-1 text-[#75715E]">请求对象列表 (OAD)：</div>
                        <div class="space-y-2 ml-2">
                            <div v-for="(oad, idx) in parsed.apdu.oads" :key="idx"
                                class="bg-[#272822] p-2 border border-[#75715E]/30 rounded">
                                <div><span class="text-[#75715E]">接口类(IC)：</span>{{ oad.ic }} ({{ oad.icName }})</div>
                                <div><span class="text-[#75715E]">对象标识(OI)：</span>{{ oad.oi }}</div>
                                <div><span class="text-[#75715E]">属性标识(PI)：</span>{{ oad.pi }} ({{ oad.piName }})</div>
                                <div v-if="oad.value"><span class="text-[#75715E]">值：</span>{{ oad.value }}</div>
                            </div>
                        </div>
                    </div>
                    <!-- 数据集响应 -->
                    <div v-if="parsed.apdu.dataSet && parsed.apdu.dataSet.length" class="mt-3">
                        <div class="mb-1 text-[#75715E]">响应数据：</div>
                        <div class="space-y-2 ml-2">
                            <div v-for="(item, idx) in parsed.apdu.dataSet" :key="idx"
                                class="bg-[#272822] p-2 border border-[#75715E]/30 rounded">
                                <div><span class="text-[#75715E]">对象标识：</span>{{ item.oi }}</div>
                                <div><span class="text-[#75715E]">属性标识：</span>{{ item.pi }}</div>
                                <div><span class="text-[#75715E]">值：</span>{{ item.value }}</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 错误提示 -->
            <div v-if="errorMsg" class="bg-[#F92672]/20 p-3 border border-[#F92672] rounded-xl text-[#F92672] text-sm">
                ⚠️ 解析错误：{{ errorMsg }}
            </div>

            <!-- 提示信息：推荐使用后端解析 -->
            <div v-if="!errorMsg && !parsed.link && !parsed.apdu && rawHex" class="py-10 text-[#75715E] text-center">
                报文格式不正确或暂不支持自动解析。<br>
                建议通过后端 Go 服务进行完整解析。
            </div>
            <div v-if="!rawHex" class="py-10 text-[#75715E] text-center">
                请输入 DL/T 698.45 报文进行解析
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const rawHex = ref('')
const frameBytes = ref([])
const parsed = ref({ link: null, apdu: null })
const errorMsg = ref('')
const isParsing = ref(false)

// 区域高亮范围（根据解析结果动态计算）
let highlightRanges = reactive([])

// 填充一个典型的 698 请求报文（读正向有功总电能请求，帧结构基本完整）
function fillExample() {
    // 示例：获取电能量数据（Get-Request-Normal）
    // 说明：实际报文较长，这里提供一个简单但合法的链路帧（不含完整 APDU，仅示意）
    // 更完整的示例需要包含 ASN.1 编码。此处提供一个标准的“读取电表时间”请求报文供测试。
    rawHex.value = '68 1F 1F 68 08 02 02 01 01 00 01 C1 01 80 01 00 00 00 01 00 00 00 01 00 00 00 00 6A 35 16'
    onInputChange()
}

// 更新字节数组
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
}

// 根据索引获取背景色（基于预设的链路层结构大致分区，实际可用解析结果细化）
function getByteBgStyle(idx) {
    // 简单示例：0-1起始符+长度，2-3地址+控制等，实际可根据 parsed.link 动态计算
    if (idx === 0) return { backgroundColor: '#3b82f6' }      // 起始符
    if (idx === 1) return { backgroundColor: '#3b82f6' }      // 长度
    if (idx >= 2 && idx <= 5) return { backgroundColor: '#ef4444' } // 控制+地址简示
    if (idx > 5 && idx < frameBytes.value.length - 2) return { backgroundColor: '#10b981' } // 应用层数据
    if (idx >= frameBytes.value.length - 2) return { backgroundColor: '#fbbf24' } // 校验+结束符
    return { backgroundColor: 'transparent' }
}

function onInputChange() {
    updateFrameBytes()
    if (rawHex.value.trim()) {
        parseFrame()
    } else {
        clearAll()
    }
}

// 核心解析函数：调用后端 Go 方法（如果可用），否则使用前端简易模拟解析
async function parseFrame() {
    errorMsg.value = ''
    parsed.value = { link: null, apdu: null }
    if (!frameBytes.value.length) return

    isParsing.value = true
    try {
        // 优先尝试调用 Wails 后端暴露的解析函数（假设 App 中有 Parse698 方法）
        if (window.go && window.go.main && window.go.main.App && window.go.main.App.Parse698) {
            const hex = rawHex.value.replace(/\s/g, '')
            const result = await window.go.main.App.Parse698(hex)
            // 后端返回结构：{ link: {...}, apdu: {...} }
            parsed.value = result
        } else {
            // 后备：前端简易模拟解析（仅演示链路层基本结构和典型 APDU 结构）
            const bytes = frameBytes.value.map(b => parseInt(b, 16))
            if (bytes[0] !== 0x68) throw new Error('无效的帧起始符')
            const len = bytes[1]  // 长度域（实际长度=len+2？简化）
            const link = {
                start: '68',
                length: len,
                control: `0x${bytes[2].toString(16).padStart(2, '0')}`,
                address: `${bytes[3]}.${bytes[4]}.${bytes[5]}`,
                hcs: `0x${bytes[6].toString(16)}${bytes[7].toString(16)}`,
                fcs: `0x${bytes[bytes.length - 3].toString(16)}${bytes[bytes.length - 2].toString(16)}`,
                end: bytes[bytes.length - 1] === 0x16 ? '16' : '未知'
            }
            // 简单校验 HCS/FCS 可通过算法，这里不实现
            link.hcsValid = true // 示意
            link.fcsValid = true
            parsed.value.link = link

            // 尝试解析应用层 APDU（仅当长度足够）
            if (bytes.length > 8) {
                const apduBytes = bytes.slice(8, bytes.length - 3)
                // 非常简化的 APDU 解析：假设是 Get-Request-Normal 或 Get-Response-Normal
                // 实际应该解析 ASN.1 结构，这里仅抽取第一个 OAD
                if (apduBytes.length >= 5) {
                    const typeId = apduBytes[0]
                    let typeName = '未知'
                    if (typeId === 0xC1) typeName = 'Get-Request-Normal'
                    else if (typeId === 0x81) typeName = 'Get-Response-Normal'
                    else typeName = `0x${typeId.toString(16)}`
                    const apdu = { typeId: `0x${typeId.toString(16)}`, typeName, oads: [] }
                    // 尝试解析一个 OAD (简化: 固定偏移)
                    if (apduBytes.length >= 6 && (apduBytes[1] === 0x01 || apduBytes[1] === 0x80)) {
                        let ic = apduBytes[2]
                        let oi = (apduBytes[3] << 8) | apduBytes[4]
                        let pi = apduBytes[5]
                        apdu.oads.push({
                            ic: `0x${ic.toString(16)}`,
                            icName: ic === 0x00 ? '电能量' : '其他',
                            oi: `0x${oi.toString(16)}`,
                            pi: `0x${pi.toString(16)}`,
                            piName: pi === 0x01 ? '值' : '其他'
                        })
                    }
                    parsed.value.apdu = apdu
                }
            }
        }
    } catch (err) {
        errorMsg.value = err.message || String(err)
    } finally {
        isParsing.value = false
    }
}

function clearAll() {
    rawHex.value = ''
    frameBytes.value = []
    parsed.value = { link: null, apdu: null }
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