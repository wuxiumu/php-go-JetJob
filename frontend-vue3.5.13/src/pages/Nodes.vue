<template>
  <el-card>
    <div style="margin-bottom:16px;">
      <el-button @click="fetchNodes">刷新</el-button>
    </div>
    <el-table :data="nodes" style="width:100%">
      <el-table-column prop="id" label="ID" width="50"/>
      <el-table-column prop="name" label="节点名"/>
      <el-table-column prop="host" label="主机"/>
      <el-table-column prop="status" label="状态">
        <template #default="scope">
          <el-tag :type="scope.row.status==='active'?'success':'info'">
            {{scope.row.status}}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="last_heartbeat" label="心跳时间"/>
    </el-table>
  </el-card>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'
const nodes = ref([])
async function fetchNodes() {
  nodes.value = (await api.get('/nodes')).data
}
onMounted(fetchNodes)
</script>