import axios from 'axios'
import type { AxiosRequestConfig, AxiosResponse } from 'axios'

// 可根据环境自动切换 API 基础地址
const baseURL = import.meta.env.VITE_API_BASE_URL || '/'

const instance = axios.create({
    baseURL,
    timeout: 10000,
    withCredentials: true
})

// 请求拦截器（如需要带 token）
instance.interceptors.request.use(
    config => {
        const token = localStorage.getItem('token')
        if (token) config.headers.Authorization = `Bearer ${token}`
        return config
    },
    error => Promise.reject(error)
)

// 响应拦截器（统一处理错误）
instance.interceptors.response.use(
    (response: AxiosResponse) => response,
    error => {
        // 可以在这里统一弹窗或输出
        if (error.response) {
            console.error(`[API Error] ${error.response.status}: ${error.response.data?.message || error.message}`)
        } else {
            console.error('[API Error] 网络异常或超时', error)
        }
        return Promise.reject(error)
    }
)

const request = {
    get<T = any>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
        return instance.get<T>(url, config)
    },
    post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
        return instance.post<T>(url, data, config)
    },
    // 新增 delete 方法
    delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
        return instance.delete<T>(url, config)
    }
}

export default request