import axios from 'axios'

const baseURL = import.meta.env.VITE_API_BASE_URL || '/api'

const request = axios.create({
    baseURL,
    timeout: 15000,
})

export default request