import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // 关键配置：遇到 /api 就转发给 localhost:8080
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  },
  build: {
    // 打包输出目录配置，方便 Go 读取
    outDir: '../dist',
    emptyOutDir: true
  }
})