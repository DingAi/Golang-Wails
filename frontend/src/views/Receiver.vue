<template>
<div>
  <div class="bg-[#272822] shadow-lg p-7 border border-[#3E3D32] rounded-2xl">
    <div class="mb-5">
      <span class="text-[#A6E22E]">运行状态：</span>
      <span class="text-[#A6E22E]">● 正在监听，等待设备接入</span>
    </div>

    <label class="block mb-2 text-[#A6E22E] text-sm">实时接收日志</label>
    <div class="bg-[#3E3D32] p-4 border border-[#75715E] rounded-xl h-80 overflow-y-auto font-mono text-[#E6DB74] whitespace-pre-wrap">
[2026-06-02 14:30:01] 接收服务启动成功
[2026-06-02 14:30:03] 等待串口/TCP客户端连接...
    </div>

    <div class="flex gap-4 mt-6">
      <button class="bg-[#F92672] hover:bg-[#dd1e63] px-6 py-3 rounded-lg text-white transition">清空日志</button>
      <button class="bg-[#66D9EF] hover:bg-[#52b9d4] px-6 py-3 rounded-lg text-[#272822] transition">导出日志</button>
    </div>
  </div>
</div>
</template>