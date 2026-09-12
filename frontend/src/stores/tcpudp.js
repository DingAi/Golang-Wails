import { defineStore } from 'pinia'
export const useTcpUdpStore = defineStore('tcpudp', {
  state: () => ({
    isRunning: false,
    mode: 'tcp-client',
    // TCP Client
    tcpClientHost: '127.0.0.1',
    tcpClientPort: 8080,
    tcpClientTimeout: 3000,
    // TCP Server
    tcpServerHost: '0.0.0.0',
    tcpServerPort: 8080,
    tcpServerMaxConn: 8,
    tcpClientConnCount: 0,
    // UDP
    udpLocalPort: 6000,
    udpRemoteHost: '127.0.0.1',
    udpRemotePort: 7000,

    // 通用接收
    autoSplit: false,
    splitMs: 20,
    messages: [],
    recvMode: 'ascii',
    sendMode: 'ascii',
    autoSend: false,
    sendSec: 1
  }),
  actions: {
    start(mode) {},
    stop() {},
    send(payload) {},
    startAutoSend(payload) {},
    stopAutoSend() {},
    clearBuffer() {},
    setRecvMode(key) {},
    setSendMode(key) {},
    updateAutoSplit() {},
    initEventListeners() {},
    destroyEventListeners() {}
  }
})
