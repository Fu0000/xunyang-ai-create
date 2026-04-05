<script setup>
import { ref, computed, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  UserOutlined,
  FileSearchOutlined,
  CameraOutlined,
  ReloadOutlined
} from '@ant-design/icons-vue'
import { useAdmin } from '../composables/useAdmin'

const router = useRouter()
const route = useRoute()
const { saveAdminToken } = useAdmin()

const collapsed = ref(false)

const moduleMenuItems = [
  { key: 'inspirations', icon: () => h(FileSearchOutlined), label: '灵感内容审核' },
  { key: 'users', icon: () => h(UserOutlined), label: '用户列表' },
  { key: 'generations', icon: () => h(CameraOutlined), label: '生成列表' }
]

const activeMenuKey = computed(() => route.name || '')
const moduleTitle = computed(() => route.meta.title || '工作台')
const moduleSubTitle = computed(() => route.meta.subTitle || '')

const onModuleSelect = ({ key }) => {
  router.push({ name: key })
}

const logout = () => {
  saveAdminToken('')
  message.info('已退出登录')
  router.push('/login')
}

// Support a way to refresh current view if needed, but standard SPA router view re-keying or explicit event bus is better
// We'll just provide an event through a global or window
const refreshCurrent = () => {
  window.dispatchEvent(new Event('admin-refresh-list'))
}
</script>

<template>
  <a-layout class="admin-layout">
    <a-layout-sider 
      v-model:collapsed="collapsed"
      :width="232" 
      theme="light" 
      class="admin-sider"
      breakpoint="lg"
      collapsed-width="0"
    >
      <div class="sider-logo">
        <span class="logo-mark">O2</span>
        <div class="logo-text">
          <strong>O2AI</strong>
          <span>Admin Workspace</span>
        </div>
      </div>
      <a-menu
        theme="light"
        mode="inline"
        :items="moduleMenuItems"
        :selected-keys="[activeMenuKey]"
        @select="onModuleSelect"
      />
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="admin-header">
        <div class="header-left">
          <h2>{{ moduleTitle }}</h2>
          <span>{{ moduleSubTitle }}</span>
        </div>
        <div class="header-right">
          <a-button :icon="h(ReloadOutlined)" @click="refreshCurrent" />
          <a-tag color="blue">已登录</a-tag>
          <a-button @click="logout">退出登录</a-button>
        </div>
      </a-layout-header>

      <a-layout-content class="admin-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<style scoped>
.admin-layout {
  min-height: 100vh;
  background: #f3f5f9;
}

.admin-sider {
  border-right: 1px solid rgba(0, 0, 0, 0.04);
}

.sider-logo {
  height: 64px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px;
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

.logo-text {
  display: flex;
  flex-direction: column;
  line-height: 1.1;
  color: #0f172a;
}

.logo-text strong {
  font-size: 13px;
}

.logo-text span {
  font-size: 12px;
  color: #64748b;
}

.admin-header {
  height: 64px;
  padding: 0 32px;
  background: rgba(255, 255, 255, 0.8) !important;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.02);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.header-left h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  line-height: 1;
  color: rgba(0, 0, 0, 0.88);
}

.header-left span {
  font-size: 13px;
  line-height: 1;
  color: #64748b;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.admin-content {
  padding: 24px;
}

@media (max-width: 900px) {
  .admin-header {
    padding: 0 16px;
  }
  .mobile-toggle {
    display: inline-flex;
  }
  .header-left {
    gap: 10px;
  }
  .header-title {
    flex-direction: column;
    align-items: flex-start;
    gap: 2px;
  }
  .header-title h2 {
    font-size: 16px;
  }
  .header-title span {
    font-size: 11px;
  }
  .admin-content {
    padding: 12px;
  }
}
</style>
