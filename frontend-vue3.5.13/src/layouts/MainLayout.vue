<template>
  <el-container style="height: 100vh">
    <el-aside width="200px">
      <el-menu :default-active="activeMenu" router>
        <el-menu-item v-for="item in menus" :key="item.path" :index="item.path">
          <el-icon><component :is="item.icon"/></el-icon>
          <span>{{item.name}}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header style="display:flex;justify-content:space-between;align-items:center;">
        <span style="font-weight:bold;font-size:20px;">JetJob 管理后台</span>
        <div>
          <el-button type="text" @click="logout">退出登录</el-button>
        </div>
      </el-header>
      <el-main>
        <router-view/>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { PieChart, Cpu, Monitor, SwitchButton } from '@element-plus/icons-vue'
const router = useRouter()
const route = useRoute()

const menus = [
  { name: '总览', path: '/', icon: PieChart },
  { name: '任务管理', path: '/tasks', icon: Cpu },
  { name: '节点管理', path: '/nodes', icon: Monitor },
  { name: '任务日志', path: '/logs', icon: SwitchButton }
]
const activeMenu = computed(() => route.path)

function logout() {
  localStorage.removeItem('jetjob_token')
  router.replace('/login')
}
</script>