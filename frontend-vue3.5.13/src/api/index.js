import axios from 'axios'

// 动态适配API地址
const API_BASE = import.meta.env.VITE_API_BASE

const instance = axios.create({
    baseURL: API_BASE,
})

// token注入+401拦截
instance.interceptors.request.use(config => {
    const token = localStorage.getItem('jetjob_token')
    if (token) config.headers['Authorization'] = `Bearer ${token}`
    return config
})
instance.interceptors.response.use(
    res => res,
    err => {
        if (err.response && err.response.status === 401) {
            localStorage.removeItem('jetjob_token')
            window.location = '/login'
        }
        return Promise.reject(err)
    }
)

export default instance