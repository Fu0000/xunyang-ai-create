import { createApp } from 'vue'
import { createPinia } from 'pinia'
import naive from 'naive-ui'
import router from './router'
import i18n from './i18n'
import './style.css'
import App from './App.vue'
import { usePricingStore } from './stores/pricing'

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)
app.use(i18n)
app.use(naive)

// TASK-18: 全局 Vue 错误处理，防止未捕获异常导致白屏
app.config.errorHandler = (err, _vm, info) => {
  console.error('[Vue Error]', err, '\nComponent Info:', info)
  // 如需接入 Sentry 等监控，在此添加：
  // Sentry.captureException(err, { extra: { info } })
}

// TASK-18: 开发模式下显示 Vue 警告
app.config.warnHandler = (msg, _vm, trace) => {
  if (import.meta.env.DEV) {
    console.warn('[Vue Warn]', msg, trace)
  }
}

// TASK-17: 监听 auth:logout 事件（由 request.js 在 401 时派发）
window.addEventListener('auth:logout', () => {
  // 使用 pinia 执行登出逻辑，需要在 app.use(pinia) 之后才能调用
  import('./stores/user').then(({ useUserStore }) => {
    const userStore = useUserStore()
    userStore.logout()
    router.push('/')
  })
})

// Fetch pricing data before mounting
usePricingStore().fetchPricing().then(() => {
  app.mount('#app')
})

