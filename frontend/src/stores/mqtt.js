// stores/mqtt.js 伪接口定义
import { defineStore } from 'pinia'
export const useMqttStore = defineStore('mqtt', {
  state: () => ({
    isConnected: false,
    brokerUrl: 'mqtt://127.0.0.1:1883',
    clientId: '',
    username: '',
    password: '',
    keepalive: 60,
    subList: [], // {topic, qos}
    messages: [],
    recvMode: 'ascii',
    sendMode: 'ascii',
    autoPub: false,
    pubSec: 1
  }),
  actions: {
    connect(){},
    disconnect(){},
    publish(topic,payload,qos){},
    startAutoPub(payload,topic,qos){},
    stopAutoPub(){},
    clearMessages(){},
    setRecvMode(key){},
    setSendMode(key){},
    initEventListeners(){},
    destroyEventListeners(){}
  }
})
