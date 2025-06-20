import { createRouter, createWebHistory } from 'vue-router'
import Login from '../pages/Login.vue'
import Dashboard from '../pages/Dashboard.vue'
import Tasks from '../pages/Tasks.vue'
import Nodes from '../pages/Nodes.vue'
import Logs from '../pages/Logs.vue'
import MainLayout from '../layouts/MainLayout.vue'

const routes = [
    { path: '/login', component: Login },
    {
        path: '/',
        component: MainLayout,
        children: [
            { path: '', component: Dashboard },
            { path: 'tasks', component: Tasks },
            { path: 'nodes', component: Nodes },
            { path: 'logs', component: Logs },
        ]
    }
]

const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('jetjob_token')
    if (to.path !== '/login' && !token) next('/login')
    else next()
})

export default router