<script setup>
import { onMounted, ref, onBeforeUnmount } from 'vue'
import { message } from 'ant-design-vue'
import { useAdmin } from '../composables/useAdmin'

const { listGenerations } = useAdmin()

const loading = ref(false)
const items = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

// Filters
const searchUid = ref('')
const searchType = ref('')
const searchStatus = ref('')
const searchTaskId = ref('')

const typeOptions = [
  { label: '全部模型/类型', value: '' },
  { label: '图片 (image)', value: 'image' },
  { label: '视频 (video)', value: 'video' },
  { label: 'Kling', value: 'kling' },
  { label: 'Luma', value: 'luma' },
  { label: 'Midjourney', value: 'midjourney' }
]

const statusOptions = [
  { label: '全部状态', value: '' },
  { label: '处理中 (processing)', value: 'processing' },
  { label: '成功 (success)', value: 'success' },
  { label: '失败/退款 (failed)', value: 'failed' }
]

const fetchList = async (resetPage = false) => {
  if (resetPage) page.value = 1
  loading.value = true
  try {
    const offset = (page.value - 1) * pageSize.value
    const data = await listGenerations({
      limit: pageSize.value,
      offset,
      user_id: searchUid.value.trim(),
      type: searchType.value.trim(),
      status: searchStatus.value,
      task_id: searchTaskId.value.trim()
    })
    items.value = data.items || []
    total.value = data.total || 0
  } catch (error) {
    message.error(error?.response?.data?.error || '加载生成记录失败')
  } finally {
    loading.value = false
  }
}

const onSearch = async () => await fetchList(true)
const onReset = async () => {
  searchUid.value = ''
  searchType.value = ''
  searchStatus.value = ''
  searchTaskId.value = ''
  await fetchList(true)
}
const onPageChange = async (nextPage) => {
  page.value = nextPage
  await fetchList(false)
}

const formatTime = (ts) => {
  if (!ts) return '-'
  return new Date(ts).toLocaleString('zh-CN', { hour12: false })
}

const columns = [
  { title: 'Task ID / 追溯码', dataIndex: 'task_id', width: 140 },
  { title: '用户归属', dataIndex: 'user_meta', width: 160 },
  { title: '调用类型', dataIndex: 'type', width: 110 },
  { title: '输入内容 (Prompt/Refs)', dataIndex: 'prompt', width: 300 },
  { title: '状态', dataIndex: 'status', width: 110 },
  { title: '结果产物', dataIndex: 'media', width: 100 },
  { title: '消费/退款', dataIndex: 'credits_cost', width: 100 },
  { title: '调用时间', dataIndex: 'created_at', width: 140 }
]

const handleRefresh = () => fetchList(false)

onMounted(() => {
  fetchList(true)
  window.addEventListener('admin-refresh-list', handleRefresh)
})

onBeforeUnmount(() => {
  window.removeEventListener('admin-refresh-list', handleRefresh)
})
</script>

<template>
  <div class="generation-list-view">
    <a-card class="mb16 seamless-card" :bordered="false" style="padding-bottom: 0;">
      <a-form layout="vertical">
        <a-row :gutter="[16, 16]">
          <a-col :xs="12" :sm="12" :md="6" :lg="4">
            <a-form-item label="类型" style="margin-bottom: 0;">
              <a-select v-model:value="searchType" :options="typeOptions" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="12" :sm="12" :md="6" :lg="4">
            <a-form-item label="状态" style="margin-bottom: 0;">
              <a-select v-model:value="searchStatus" :options="statusOptions" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="4">
            <a-form-item label="用户 UID" style="margin-bottom: 0;">
              <a-input v-model:value="searchUid" allow-clear placeholder="数字 ID" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="6">
            <a-form-item label="Task ID" style="margin-bottom: 0;">
              <a-input v-model:value="searchTaskId" allow-clear placeholder="三方渠道 Task ID" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="24" :md="24" :lg="6">
            <a-form-item label=" " style="margin-bottom: 0;">
              <div style="display: flex; gap: 8px;">
                <a-button type="primary" :disabled="loading" @click="onSearch">查询</a-button>
                <a-button :disabled="loading" @click="onReset">重置</a-button>
              </div>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-card>

    <a-alert
      v-if="searchStatus === 'failed'"
      message="目前可能正查看由于错误引起的退款请求，可通过 Task ID 在三方提供商后台查阅失败详情。"
      type="error"
      show-icon
      class="mb16"
    />

    <a-card class="table-card seamless-card" :bordered="false" :body-style="{ padding: '0' }">
      <a-table
        :columns="columns"
        :data-source="items"
        :loading="loading"
        :pagination="false"
        :row-key="record => record.id"
        size="small"
        :scroll="{ x: 1200 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'task_id'">
            <span class="code-font">{{ record.task_id || `INT-${record.id}` }}</span>
          </template>

          <template v-else-if="column.dataIndex === 'user_meta'">
            <div class="cell-stack">
              <strong>{{ record.user_email }}</strong>
              <span class="muted">UID: {{ record.user_id }}</span>
            </div>
          </template>

          <template v-else-if="column.dataIndex === 'type'">
            <a-tag color="blue">{{ record.type }}</a-tag>
          </template>

          <template v-else-if="column.dataIndex === 'prompt'">
            <div class="prompt-text" :title="record.prompt">{{ record.prompt }}</div>
            <div class="refs" v-if="record.reference_images && record.reference_images.length > 0">
              <a-badge :count="record.reference_images.length" color="#52c41a" title="垫图数量" />
            </div>
          </template>

          <template v-else-if="column.dataIndex === 'status'">
            <a-tag v-if="record.status === 'success'" color="success">成功</a-tag>
            <a-tooltip v-else-if="record.status === 'failed'" :title="record.error_msg">
              <a-tag color="error">失败 (退款)</a-tag>
            </a-tooltip>
            <a-tag v-else color="processing">{{ record.status }}</a-tag>
          </template>

          <template v-else-if="column.dataIndex === 'media'">
            <div class="media-box" v-if="record.status === 'success'">
              <video 
                v-if="record.video_url" 
                :src="record.video_url" 
                muted playsinline preload="none" 
              />
              <img 
                v-else-if="record.images && record.images.length > 0" 
                :src="record.images[0]" 
                alt="thumb" loading="lazy" 
              />
            </div>
            <span class="muted" v-else>-</span>
          </template>

          <template v-else-if="column.dataIndex === 'credits_cost'">
            <span :style="{ color: record.status === 'failed' ? '#ff4d4f' : '#faad14', fontWeight: 600 }">
              {{ record.status === 'failed' ? '已退回' : `💎 -${record.credits_cost}` }}
            </span>
          </template>

          <template v-else-if="column.dataIndex === 'created_at'">
            <div class="cell-stack">
              <span>{{ formatTime(record.created_at) }}</span>
              <span class="muted">{{ formatTime(record.updated_at) }}</span>
            </div>
          </template>
        </template>
      </a-table>

      <div class="table-pagination">
        <a-pagination
          :current="page"
          :page-size="pageSize"
          :total="total"
          :show-size-changer="false"
          @change="onPageChange"
        />
      </div>
    </a-card>
  </div>
</template>

<style scoped>
.mb16 { margin-bottom: 24px; }
.cell-stack { display: flex; flex-direction: column; gap: 2px; }
.code-font { font-family: ui-monospace, SFMono-Regular, Consolas, "Liberation Mono", Menlo, monospace; font-size: 12px; color: #1677ff; }

.prompt-text {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  font-size: 13px;
  line-height: 1.4;
  color: rgba(0,0,0,0.7);
}
.refs { margin-top: 4px; }

.media-box {
  width: 48px;
  height: 48px;
  border-radius: 6px;
  overflow: hidden;
  background: #f0f0f0;
  border: 1px solid #e8e8e8;
}
.media-box img, .media-box video {
  width: 100%; height: 100%; object-fit: cover;
}

.muted { color: rgba(0, 0, 0, 0.45); font-size: 12px; }
.table-pagination { display: flex; justify-content: flex-end; padding: 12px 16px; border-top: 1px solid #f0f0f0; }

:deep(.ant-table-tbody > tr:hover > td) {
  background: #fdfefe !important;
}
</style>
