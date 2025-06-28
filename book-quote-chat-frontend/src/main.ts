import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import './assets/tailwind.css';
import 'element-plus/dist/index.css'
import VueMatomo from 'vue-matomo';
import ElementPlus from 'element-plus'

createApp(App)
  .use(router)
  .use(ElementPlus)
    .use(VueMatomo, {
    // Configure your matomo server and site by providing
    host: 'https://tj.97gaoqian.com',
    siteId: 20,
  })
  .mount('#app');

window._paq = window._paq || [];
window._paq.push(['trackPageView']);