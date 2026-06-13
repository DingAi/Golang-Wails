import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { Events } from '@wailsio/runtime'
import * as SerialService from '../../bindings/project_service/internal/service/serialservice.js'

export const useSerialStore = defineStore('serial', () => {
  // ---------- State ----------
  const isOpen = ref(false)
  const portName = ref('COM5')
  const baudRate = ref(115200)
  const dataBits = ref(8)
  const parity = ref('无')
  const stopBits = ref(1)
  const rxBytes = ref(0)
  const txBytes = ref(0)
  const ports = ref([])

  // Chat-style message log: { type: 'sent'|'recv', text: string, time: string }
  const messages = ref([])

  const autoSplit = ref(true)
  const splitMs = ref(20)
  const saveToFile = ref(false)
  const autoSend = ref(false)
  const sendSec = ref(1.0)
  const recvMode = ref('hex')
  const sendMode = ref('hex')

  const statusText = computed(() => {
    return isOpen.value ? '打开' : '关闭'
  })

  const now = () => {
    const d = new Date()
    return d.toLocaleTimeString('zh-CN', { hour12: false })
  }

  // ---------- Init: listen to backend events ----------
  let statusUnlisten = null
  let dataUnlisten = null

  function initEventListeners() {
    statusUnlisten = Events.On('serial:status', (event) => {
      const s = event.data
      isOpen.value = s.is_open
      portName.value = s.port_name
      baudRate.value = s.baud_rate
      rxBytes.value = s.rx_bytes
      txBytes.value = s.tx_bytes
    })

    dataUnlisten = Events.On('serial:data', (event) => {
      const d = event.data
      const text = formatData(d, recvMode.value)
      messages.value.push({
        type: 'recv',
        text: text,
        time: now(),
      })
    })
  }

  function destroyEventListeners() {
    if (statusUnlisten) statusUnlisten()
    if (dataUnlisten) dataUnlisten()
  }

  // ---------- Actions ----------
  async function fetchPorts() {
    try {
      ports.value = await SerialService.GetSerialPorts()
    } catch (err) {
      console.error('获取串口列表失败:', err)
    }
  }

  async function open() {
    try {
      await SerialService.OpenSerialPort({
        port_name: portName.value,
        baud_rate: baudRate.value,
        data_bits: dataBits.value,
        parity: parity.value,
        stop_bits: stopBits.value,
      })
      await fetchPorts()
    } catch (err) {
      console.error('打开串口失败:', err)
      throw err
    }
  }

  async function close() {
    try {
      await SerialService.CloseSerialPort()
    } catch (err) {
      console.error('关闭串口失败:', err)
      throw err
    }
  }

  /**
   * Send data and add a sent bubble to the message log
   * @param {string} data - raw string to send
   * @param {string} [format] - format override, uses sendMode if omitted
   * @returns {Promise<number>} bytes sent
   */
  async function send(data, format) {
    if (!isOpen.value) {
      throw new Error('串口未打开')
    }
    const fmt = format || sendMode.value
    // Show what we're sending in the current send mode format
    const displayText = formatDisplayFromInput(data, fmt)
    messages.value.push({
      type: 'sent',
      text: displayText,
      time: now(),
    })
    return await SerialService.SendData(data, fmt)
  }

  async function startAutoSend(data, format) {
    if (!isOpen.value) return
    await SerialService.StartAutoSend(data, format || sendMode.value, sendSec.value)
  }

  async function stopAutoSend() {
    await SerialService.StopAutoSend()
  }

  async function clearBuffer() {
    messages.value = []
    await SerialService.ClearReceiveBuffer()
  }

  async function updateAutoSplit() {
    await SerialService.SetAutoSplit(autoSplit.value, splitMs.value)
  }

  async function refreshStatus() {
    try {
      const status = await SerialService.GetStatus()
      isOpen.value = status.is_open
      portName.value = status.port_name
      baudRate.value = status.baud_rate
      rxBytes.value = status.rx_bytes
      txBytes.value = status.tx_bytes
    } catch (err) {
      console.error('获取状态失败:', err)
    }
  }

  function setRecvMode(mode) {
    recvMode.value = mode
  }

  function setSendMode(mode) {
    sendMode.value = mode
  }

  return {
    // state
    isOpen, portName, baudRate, dataBits, parity, stopBits,
    rxBytes, txBytes, ports, messages,
    autoSplit, splitMs, saveToFile, autoSend, sendSec,
    recvMode, sendMode, statusText,
    // lifecycle
    initEventListeners, destroyEventListeners,
    // actions
    fetchPorts, open, close, send,
    startAutoSend, stopAutoSend, clearBuffer,
    updateAutoSplit, refreshStatus,
    setRecvMode, setSendMode,
  }
})

// ---------- Helpers ----------

// Format raw received data for display bubbles
function formatData(data, mode) {
  const { hex, ascii, raw } = data
  switch (mode) {
    case 'ascii':
      return ascii || ''
    case 'hex': {
      const h = hex || ''
      const pairs = []
      for (let i = 0; i < h.length; i += 2) {
        pairs.push(h.substring(i, i + 2))
      }
      return pairs.join(' ')
    }
    case 'bin':
      if (raw && raw.length) {
        return Array.from(raw).map(b => b.toString(2).padStart(8, '0')).join(' ')
      }
      return ''
    case 'dec':
      if (raw && raw.length) {
        return Array.from(raw).map(b => b.toString(10)).join(' ')
      }
      return ''
    case 'bcd':
      return hex || ''
    default:
      return hex || ''
  }
}

// Format input text for sent bubbles (no real conversion, just display as-is based on mode)
function formatDisplayFromInput(input, mode) {
  if (!input) return ''
  switch (mode) {
    case 'ascii':
      return input
    case 'hex':
      // Remove spaces and uppercase for clean display
      return input.replace(/\s+/g, ' ').trim().toUpperCase()
    case 'bin':
    case 'dec':
      return input.replace(/\s+/g, ' ').trim()
    case 'bcd':
      return input.replace(/\s+/g, ' ').trim().toUpperCase()
    default:
      return input
  }
}