<template>
  <el-card>
    <div style="margin-bottom:16px;display:flex;justify-content:space-between">
      <el-button type="primary" @click="dialogVisible=true">新建任务</el-button>
      <el-button @click="fetchTasks">刷新</el-button>
    </div>
    <el-table :data="tasks" style="width:100%">
      <el-table-column prop="id" label="ID" width="50"/>
      <el-table-column prop="name" label="任务名"/>
      <el-table-column prop="command" label="命令"/>
      <el-table-column prop="schedule" label="定时"/>
      <el-table-column prop="status" label="状态"/>
      <el-table-column label="操作" width="140">
        <template #default="scope">
          <el-button type="danger" size="small" @click="removeTask(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="新建任务">
      <el-form :model="form">
        <el-form-item label="任务名"><el-input v-model="form.name"/></el-form-item>
        <el-form-item label="命令"><el-input v-model="form.command"/></el-form-item>
        <el-form-item label="定时"><el-input v-model="form.schedule"/></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible=false">取消</el-button>
        <el-button type="primary" @click="createTask">创建</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const tasks = ref([])
const dialogVisible = ref(false)
const form = ref({ name: '', command: '', schedule: '', status: 'active' })

async function fetchTasks() {
  tasks.value = (await api.get('/tasks')).data
}

async function createTask() {
  await api.post('/tasks', form.value)
  dialogVisible.value = false
  form.value = { name: '', command: '', schedule: '', status: 'active' }
  fetchTasks()
}

async function removeTask(id) {
  await api.delete(`/tasks/${id}`)
  fetchTasks()
}

onMounted(fetchTasks)
</script>