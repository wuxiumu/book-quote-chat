import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import './assets/tailwind.css';
import 'element-plus/dist/index.css'
import VueMatomo from 'vue-matomo';
import ElementPlus from 'element-plus'
import VChart from 'vue-echarts'

// echarts 渲染器必须注册
import * as echarts from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
echarts.use([CanvasRenderer])

const app = createApp(App)
app.use(router)
app.use(ElementPlus)
app.use(VueMatomo, {
    host: 'https://tj.97gaoqian.com',
    siteId: 20,
})
app.component('v-chart', VChart) // 注册全局 v-chart
app.mount('#app')

// Matomo 追踪代码（可选）
window._paq = window._paq || [];
window._paq.push(['trackPageView']);