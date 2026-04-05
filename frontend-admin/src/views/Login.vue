<script setup>
import { h, ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { LockOutlined } from '@ant-design/icons-vue'
import { useAdmin } from '../composables/useAdmin'

const router = useRouter()
const { saveAdminToken, listAdminInspirations } = useAdmin()

const loginTokenInput = ref('')
const authLoading = ref(false)

const authenticate = async () => {
  const token = loginTokenInput.value.trim()
  if (!token) {
    message.warning('请输入管理端 Token')
    return
  }

  authLoading.value = true
  saveAdminToken(token)
  
  try {
    // Verify token by making a lightweight request
    await listAdminInspirations({ limit: 1, offset: 0, review_status: 'all' })
    message.success('登录成功')
    router.push('/')
  } catch (error) {
    saveAdminToken('') // clear incorrect token
    const errMsg = error?.response?.data?.error || 'Token 无效或服务不可用'
    message.error(errMsg)
  } finally {
    authLoading.value = false
  }
}
</script>

<template>
  <div class="login-shell">
    <div class="login-grid">
      <section class="intro-panel">
        <span class="intro-badge">ADMIN PANEL</span>
        <h1>O2AI 管理端</h1>
        <p>统一管理灵感审核、用户资产与生成风控审查。凭后台 ADMIN_TOKEN 安全介入系统运营。</p>
        <ul>
          <li>灵感内容审核</li>
          <li>用户列表与资产调配</li>
          <li>生成列表审计日志</li>
        </ul>
      </section>

      <a-card class="login-card seamless-card" :bordered="false">
        <div class="login-head">
          <span class="logo-mark">O2</span>
          <div>
            <h2>管理端登录</h2>
            <p>请输入后端配置的 ADMIN_TOKEN</p>
          </div>
        </div>
        <a-alert type="info" show-icon message="仅授权管理员可访问此后台" class="login-alert" />
        <a-input-password
          v-model:value="loginTokenInput"
          size="large"
          :prefix="h(LockOutlined)"
          placeholder="请输入 ADMIN_TOKEN"
          @press-enter="authenticate"
        />
        <a-button type="primary" size="large" block :loading="authLoading" @click="authenticate" style="margin-top: 16px;">
          验证并登录
        </a-button>
      </a-card>
    </div>
  </div>
</template>

<style scoped>
.login-shell {
  min-height: 100vh;
  padding: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background:
    radial-gradient(circle at 10% 15%, rgba(79, 70, 229, 0.12), transparent 36%),
    radial-gradient(circle at 84% 5%, rgba(124, 58, 237, 0.12), transparent 42%),
    #f8fafc;
}

.login-grid {
  width: 100%;
  max-width: 1080px;
  display: grid;
  grid-template-columns: 1.1fr 1fr;
  gap: 18px;
}

.intro-panel {
  border-radius: 16px;
  padding: 40px;
  color: #fff;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  box-shadow: 0 16px 40px rgba(79, 70, 229, 0.25);
}

.intro-badge {
  font-size: 11px;
  letter-spacing: 0.12em;
  opacity: 0.78;
}

.intro-panel h1 {
  margin: 14px 0 10px;
  font-size: 36px;
  line-height: 1.15;
}

.intro-panel p {
  margin: 0;
  font-size: 15px;
  color: rgba(255, 255, 255, 0.86);
}

.intro-panel ul {
  margin: 22px 0 0;
  padding-left: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.login-card {
  border-radius: 16px;
  padding: 12px 10px;
}

.login-head {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.login-head h2 {
  margin: 0;
  font-size: 22px;
}

.login-head p {
  margin: 2px 0 0;
  color: rgba(0, 0, 0, 0.45);
  font-size: 13px;
}

.login-alert {
  margin-bottom: 12px;
}

.logo-mark {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  color: #fff;
  background: linear-gradient(135deg, #6366f1, #a855f7);
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3);
}

@media (max-width: 1200px) {
  .login-grid {
    grid-template-columns: 1fr;
  }
}
</style>
