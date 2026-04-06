<script setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { useAdmin } from '../composables/useAdmin'

const { getSettings, updateSetting } = useAdmin()

const loading = ref(false)
const saving = ref(false)

const settings = ref({
  inspiration_auto_approve: true
})

const loadSettings = async () => {
  loading.value = true
  try {
    const data = await getSettings()
    if (data && data.settings) {
      const autoApproveItem = data.settings.find(s => s.key === 'inspiration_auto_approve')
      if (autoApproveItem) {
        settings.value.inspiration_auto_approve = autoApproveItem.value === 'true'
      }
    }
  } catch (error) {
    message.error('加载系统设置失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const toggleAutoApprove = async (checked) => {
  saving.value = true
  try {
    await updateSetting({
      key: 'inspiration_auto_approve',
      value: checked ? 'true' : 'false'
    })
    message.success('配置已保存生效')
    settings.value.inspiration_auto_approve = checked
  } catch (error) {
    message.error('保存设置失败: ' + error.message)
    // Revert visually on error
    settings.value.inspiration_auto_approve = !checked
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadSettings()
})
</script>

<template>
  <div class="settings-container">
    <a-card title="全局参数控制" :bordered="false" class="settings-card">
      <a-spin :spinning="loading">
        <div class="setting-item">
          <div class="setting-info">
            <h3 class="setting-title">灵感分享免审核自动通过</h3>
            <p class="setting-desc">开启后，用户在前台分享的作品将不经过「灵感内容审核」模块，直接进入公开瀑布流大厅。关闭后需要管理员手动审核通过才能展示。</p>
          </div>
          <div class="setting-action">
            <a-switch
              :checked="settings.inspiration_auto_approve"
              :loading="saving"
              checked-children="已开启"
              un-checked-children="需审核"
              @change="toggleAutoApprove"
            />
          </div>
        </div>
      </a-spin>
    </a-card>
  </div>
</template>

<style scoped>
.settings-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
  max-width: 800px;
}

.settings-card {
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-info {
  flex: 1;
  padding-right: 24px;
}

.setting-title {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 500;
  color: #1f2937;
}

.setting-desc {
  margin: 0;
  font-size: 13px;
  color: #6b7280;
  line-height: 1.5;
}

.setting-action {
  flex-shrink: 0;
}
</style>
