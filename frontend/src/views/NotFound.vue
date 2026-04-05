<script setup>
import { useRouter } from 'vue-router'
import { NResult, NButton } from 'naive-ui'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const { t } = useI18n()

const goHome = () => {
  router.push('/')
}
</script>

<template>
  <div class="not-found-page">
    <NResult
      status="404"
      title="404 Not Found"
      description="您访问的页面不存在或已被移除"
      class="result-container"
    >
      <template #footer>
        <NButton
          type="primary"
          size="large"
          @click="goHome"
        >
          返回首页
        </NButton>
      </template>
    </NResult>
  </div>
</template>

<style scoped>
.not-found-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--n-color);
}

.result-container {
  transform: translateY(-10vh);
}

:deep(.n-result-header__title) {
  color: var(--color-text-primary);
  font-size: 32px;
  font-weight: 800;
}

:deep(.n-result-header__description) {
  color: var(--color-text-secondary);
  font-size: 16px;
  margin-top: 12px;
}
</style>
