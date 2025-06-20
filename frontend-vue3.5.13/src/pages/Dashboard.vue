<template>
  <el-row :gutter="24">
    <el-col :span="8">
      <el-card><div>任务数</div><div style="font-size:32px">{{taskCount}}</div></el-card>
    </el-col>
    <el-col :span="8">
      <el-card><div>在线节点</div><div style="font-size:32px">{{activeNodes}}</div></el-card>
    </el-col>
    <el-col :span="8">
      <el-card><div>离线节点</div><div style="font-size:32px">{{offlineNodes}}</div></el-card>
    </el-col>
  </el-row>
  <el-row style="margin-top:32px">
    <el-col :span="24">
      <el-card>
        <v-chart :option="taskChartOption" style="height:300px"/>
      </el-card>
    </el-col>
  </el-row>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'
import VChart from 'vue-echarts'

const taskCount = ref(0), activeNodes = ref(0), offlineNodes = ref(0)
const taskChartOption = ref({})

async function fetchData() {
  const tasks = (await api.get('/tasks')).data
  taskCount.value = tasks.length
  const nodes = (await api.get('/nodes')).data
  activeNodes.value = nodes.filter(n=>n.status==='active').length
  offlineNodes.value = nodes.filter(n=>n.status==='offline').length

  // 简单任务类型分布
  const typeMap = {}
  tasks.forEach(t => typeMap[t.type] = (typeMap[t.type]||0)+1)
  taskChartOption.value = {
    title: { text: '任务类型分布', left: 'center' },
    tooltip: { trigger: 'item' },
    legend: { bottom: 0 },
    series: [{
      type: 'pie',
      radius: '50%',
      data: Object.keys(typeMap).map(k => ({ name: k, value: typeMap[k] })),
      label: { formatter: '{b}: {c} ({d}%)' }
    }]
  }
}
onMounted(fetchData)
</script>