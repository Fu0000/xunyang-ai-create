<script setup>
import { computed, h, onMounted, ref, onBeforeUnmount } from 'vue'
import { message } from 'ant-design-vue'
import { useAdmin } from '../composables/useAdmin'

const { listAdminInspirations, reviewInspiration } = useAdmin()

const loading = ref(false)
const activePostID = ref(0)
const items = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const reviewStatus = ref('pending')
const userID = ref('')
const keyword = ref('')
const startDate = ref('')
const endDate = ref('')

const reviewStatusOptions = [
  { label: '待审核', value: 'pending' },
  { label: '全部', value: 'all' },
  { label: '已通过', value: 'approved' },
  { label: '已驳回', value: 'rejected' }
]

const pendingCountInPage = computed(() => items.value.filter((i) => i.review_status === 'pending').length)
const approvedCountInPage = computed(() => items.value.filter((i) => i.review_status === 'approved').length)
const rejectedCountInPage = computed(() => items.value.filter((i) => i.review_status === 'rejected').length)

const isVideoPost = (item) => item?.type === 'video' || !!item?.video_url

const statusText = (status) => {
  if (status === 'approved') return '已通过'
  if (status === 'rejected') return '已驳回'
  return '待审核'
}

const statusColor = (status) => {
  if (status === 'approved') return 'success'
  if (status === 'rejected') return 'error'
  return 'processing'
}

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  return new Date(timestamp).toLocaleString('zh-CN', { hour12: false })
}

const extractError = (error, fallback) => error?.response?.data?.error || fallback

const fetchList = async (resetPage = false) => {
  if (resetPage) page.value = 1
  loading.value = true
  try {
    const offset = (page.value - 1) * pageSize.value
    const data = await listAdminInspirations({
      limit: pageSize.value,
      offset,
      review_status: reviewStatus.value,
      user_id: userID.value.trim(),
      q: keyword.value.trim(),
      start_date: startDate.value,
      end_date: endDate.value
    })
    items.value = data.items || []
    total.value = data.total || 0
  } catch (error) {
    if ([401, 403].includes(error?.response?.status)) {
      message.error('登录失效')
      // Handled by interceptor theoretically, but safe to notice here
    } else {
      message.error(extractError(error, '加载审核列表失败'))
    }
  } finally {
    loading.value = false
  }
}

const onSearch = async () => {
  await fetchList(true)
}

const onReset = async () => {
  reviewStatus.value = 'pending'
  userID.value = ''
  keyword.value = ''
  startDate.value = ''
  endDate.value = ''
  await fetchList(true)
}

const onPageChange = async (nextPage) => {
  page.value = nextPage
  await fetchList(false)
}

const doReview = async (item, action) => {
  if (!item?.id) return

  let note = ''
  if (action === 'reject') {
    const value = window.prompt('请输入驳回原因（可选）', '')
    if (value === null) return
    note = value.trim()
  }

  activePostID.value = item.id
  try {
    const data = await reviewInspiration(item.id, { action, note })
    const next = data?.item
    if (next) {
      const index = items.value.findIndex((row) => row.id === next.id)
      if (index >= 0) items.value[index] = next
    }
    message.success(action === 'approve' ? '审核通过' : '已驳回')
  } catch (error) {
    message.error(extractError(error, '提交审核失败'))
  } finally {
    activePostID.value = 0
  }
}

const columns = [
  { title: 'ID', dataIndex: 'id', width: 90 },
  { title: '封面', dataIndex: 'cover', width: 120 },
  { title: '标题', dataIndex: 'title' },
  { title: '作者', dataIndex: 'author', width: 180 },
  { title: '状态', dataIndex: 'review_status', width: 120 },
  { title: '发布时间', dataIndex: 'published_at', width: 190 },
  { title: '操作', dataIndex: 'actions', width: 180 }
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
  <div class="inspiration-review">
    <a-row :gutter="16" class="mb16">
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card seamless-card" :bordered="false" size="small">
          <span class="stat-label">总条数</span>
          <a-statistic :value="total" />
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card seamless-card" :bordered="false" size="small">
          <span class="stat-label">当前页待审核</span>
          <a-statistic :value="pendingCountInPage" />
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card seamless-card" :bordered="false" size="small">
          <span class="stat-label">当前页已通过</span>
          <a-statistic :value="approvedCountInPage" />
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card seamless-card" :bordered="false" size="small">
          <span class="stat-label">当前页已驳回</span>
          <a-statistic :value="rejectedCountInPage" />
        </a-card>
      </a-col>
    </a-row>

    <a-card class="mb16 filter-card seamless-card" :bordered="false">
      <a-form layout="vertical">
        <a-row :gutter="[16, 16]">
          <a-col :xs="24" :sm="12" :md="8" :lg="4">
            <a-form-item label="审核状态" style="margin-bottom: 0;">
              <a-select v-model:value="reviewStatus" :options="reviewStatusOptions" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8" :lg="4">
            <a-form-item label="用户 ID" style="margin-bottom: 0;">
              <a-input v-model:value="userID" allow-clear placeholder="数字 UID" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="12" :lg="4">
            <a-form-item label="关键词" style="margin-bottom: 0;">
              <a-input v-model:value="keyword" allow-clear placeholder="标题/提示词/所属 share id" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="12" :lg="4">
            <a-form-item label="开始日期" style="margin-bottom: 0;">
              <input v-model="startDate" type="date" class="native-date" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="12" :lg="4">
            <a-form-item label="结束日期" style="margin-bottom: 0;">
              <input v-model="endDate" type="date" class="native-date" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="24" :md="24" :lg="2">
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

    <a-card class="table-card seamless-card" :bordered="false">
      <a-table
        :columns="columns"
        :data-source="items"
        :loading="loading"
        :pagination="false"
        :row-key="(record) => record.id"
        size="middle"
        :scroll="{ x: 1000 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'cover'">
            <div class="cover-box">
              <video
                v-if="isVideoPost(record) && record.video_url"
                :src="record.video_url"
                :poster="record.cover_url"
                muted
                playsinline
                preload="metadata"
              />
              <img v-else :src="record.cover_url || record.images?.[0]" alt="cover" loading="lazy" />
            </div>
          </template>
          <template v-else-if="column.dataIndex === 'title'">
            <div class="title-cell">
              <strong>{{ record.title || '未命名内容' }}</strong>
              <span class="muted">share: {{ record.share_id }}</span>
            </div>
          </template>
          <template v-else-if="column.dataIndex === 'author'">
            <div class="title-cell">
              <span>{{ record.author?.nickname || '-' }}</span>
              <span class="muted">UID: {{ record.author?.user_id || '-' }}</span>
            </div>
          </template>
          <template v-else-if="column.dataIndex === 'review_status'">
            <a-tag :color="statusColor(record.review_status)">
              {{ statusText(record.review_status) }}
            </a-tag>
          </template>
          <template v-else-if="column.dataIndex === 'published_at'">
            {{ formatTime(record.published_at) }}
          </template>
          <template v-else-if="column.dataIndex === 'actions'">
            <div class="action-cell">
              <a-button
                type="primary"
                size="small"
                shape="round"
                :disabled="record.review_status === 'approved'"
                :loading="activePostID === record.id"
                @click="doReview(record, 'approve')"
                style="background-color: #10b981; border-color: #10b981;"
              >
                予以通过
              </a-button>
              <a-button
                danger
                size="small"
                shape="round"
                :disabled="record.review_status === 'rejected'"
                :loading="activePostID === record.id"
                @click="doReview(record, 'reject')"
              >
                驳回下架
              </a-button>
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
.stat-label {
  display: block;
  color: rgba(0, 0, 0, 0.45);
  font-size: 12px;
  margin-bottom: 4px;
}
.cover-box {
  width: 72px; height: 72px; border-radius: 8px; overflow: hidden; background: #f5f5f5; border: 1px solid #f0f0f0;
}
.cover-box img, .cover-box video {
  width: 100%; height: 100%; object-fit: cover;
}
.title-cell { display: flex; flex-direction: column; gap: 3px; }
.title-cell strong { color: rgba(0, 0, 0, 0.88); }
.muted { color: rgba(0, 0, 0, 0.45); font-size: 12px; }
.action-cell { display: flex; gap: 8px; }
.table-pagination { display: flex; justify-content: flex-end; margin-top: 16px; }
.native-date {
  width: 138px; height: 32px; border: 1px solid #d9d9d9; border-radius: 8px; padding: 0 10px; background: #fff; box-sizing: border-box; outline: none; transition: border-color 0.2s;
}
.native-date:focus {
  border-color: #4f46e5;
}
:deep(.ant-table-tbody > tr:hover > td) {
  background: #f7faff !important;
}
</style>
