import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('node_modules')) {
            if (id.includes('element-plus')) return 'element-plus' // 不要拆分 icons/locale
            if (id.includes('ant-design-vue')) return 'ant-design-vue'
            if (id.includes('agora')) return 'agora'
            if (id.includes('axios')) return 'axios'
            if (id.includes('dayjs')) return 'dayjs'
            if (id.includes('lodash')) return 'lodash'
            if (id.includes('vue')) return 'vue-vendor'
            return 'vendor'
          }
          // 可按业务大页面分包，假设 chat 相关很大：
          if (id.includes('/src/views/ChatRoom.vue')) return 'chatroom'
          if (id.includes('/src/views/Friends.vue')) return 'friends'
        }
      }
    },
    chunkSizeWarningLimit: 12000 // 可选，避免过多无用警告
  }
})