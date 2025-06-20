<template>
  <el-card>
    <div style="margin-bottom:8px;">
      <el-button @click="connect">重新连接</el-button>
    </div>
    <div ref="logBox" style="background:#141414;color:#35ff55;padding:10px;font-family:monospace;height:300px;overflow:auto">
      <div v-for="(line, i) in logLines" :key="i">{{line}}</div>
    </div>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const logLines = ref([])
const logBox = ref(null)
let ws = null

function connect() {
  if (ws) ws.close()
  ws = new WebSocket('ws://localhost:8090/ws/logs')
  ws.onmessage = (evt) => {
    logLines.value.push(evt.data)
    if (logBox.value) {
      logBox.value.scrollTop = logBox.value.scrollHeight
    }
  }
}
onMounted(connect)
</script>