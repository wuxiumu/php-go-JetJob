<template>
  <el-row justify="center" align="middle" style="height:100vh;">
    <el-col :span="8">
      <el-card>
        <el-form @submit.prevent="login">
          <el-form-item label="用户名">
            <el-input v-model="username" placeholder="admin" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="password" show-password placeholder="123456" />
          </el-form-item>
          <el-button type="primary" @click="login" style="width:100%">登录</el-button>
        </el-form>
      </el-card>
    </el-col>
  </el-row>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

const username = ref('')
const password = ref('')
const router = useRouter()

async function login() {
  try {
    const res = await api.post('/login', {
      username: username.value,
      password: password.value
    })
    if (res.data.token) {
      localStorage.setItem('jetjob_token', res.data.token)
      router.replace('/')
    }
  } catch (err) {
    alert('登录失败，账号或密码错误')
  }
}
</script>