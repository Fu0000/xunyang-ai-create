<script setup>
import { computed, onMounted, ref, onBeforeUnmount } from 'vue'
import { message } from 'ant-design-vue'
import { useAdmin } from '../composables/useAdmin'

const { listUsers, updateUserCredits, updateUserStatus } = useAdmin()

const loading = ref(false)
const items = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

// Filters
const searchUid = ref('')
const searchEmail = ref('')
const searchStatus = ref('')

const statusOptions = [
  { label: '全部状态', value: '' },
  { label: '正常 (active)', value: 'active' },
  { label: '封禁 (banned)', value: 'banned' },
  { label: '停用 (disabled)', value: 'disabled' },
]

// Modal - Credits
const showCreditModal = ref(false)
const creditLoading = ref(false)
const targetUser = ref(null)
const creditForm = ref({ delta: 0, note: '' })

// Actions
const fetchList = async (resetPage = false) => {
  if (resetPage) page.value = 1
  loading.value = true
  try {
    const offset = (page.value - 1) * pageSize.value
    const data = await listUsers({
      limit: pageSize.value,
      offset,
      user_id: searchUid.value.trim(),
      email: searchEmail.value.trim(),
      status: searchStatus.value
    })
    items.value = data.items || []
    total.value = data.total || 0
  } catch (error) {
    message.error(error?.response?.data?.error || '加载用户列表失败')
  } finally {
    loading.value = false
  }
}

const onSearch = async () => await fetchList(true)
const onReset = async () => {
  searchUid.value = ''
  searchEmail.value = ''
  searchStatus.value = ''
  await fetchList(true)
}
const onPageChange = async (nextPage) => {
  page.value = nextPage
  await fetchList(false)
}

// Credits Management
const openCreditModal = (user) => {
  targetUser.value = user
  creditForm.value = { delta: 100, note: 'Admin grant' }
  showCreditModal.value = true
}

const submitCreditUpdate = async () => {
  if (!targetUser.value || creditLoading.value) return
  if (creditForm.value.delta === 0) {
    message.warning('变动额度不能为 0')
    return
  }

  creditLoading.value = true
  try {
    await updateUserCredits(targetUser.value.id, creditForm.value)
    message.success('资产更新成功')
    showCreditModal.value = false
    await fetchList(false)
  } catch (error) {
    message.error(error?.response?.data?.error || '变更资产失败')
  } finally {
    creditLoading.value = false
  }
}

// Status Management
const doChangeStatus = async (user, targetStatus) => {
  if (!confirm(`确认将 ${user.email} (${user.id}) 状态修改为 ${targetStatus}?`)) return
  
  try {
    await updateUserStatus(user.id, { status: targetStatus })
    message.success('状态变更成功')
    await fetchList(false)
  } catch (error) {
    message.error(error?.response?.data?.error || '状态变更失败')
  }
}

const formatTime = (ts) => {
  if (!ts) return '-'
  return new Date(ts).toLocaleString('zh-CN', { hour12: false })
}

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '账号特征', dataIndex: 'user_meta', width: 220 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '剩余钻石', dataIndex: 'credits', width: 100, align: 'right' },
  { title: '使用/推广', dataIndex: 'activity', width: 130 },
  { title: '注册时间', dataIndex: 'created_at', width: 160 },
  { title: '最后登录', dataIndex: 'last_login_at', width: 160 },
  { title: '快速操作', dataIndex: 'actions', width: 200 }
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
  <div class="user-list-view">
    <a-card class="mb16 seamless-card" :bordered="false" style="padding-bottom: 0;">
      <a-form layout="vertical">
        <a-row :gutter="[16, 16]">
          <a-col :xs="24" :sm="12" :md="6" :lg="4">
            <a-form-item label="UID" style="margin-bottom: 0;">
              <a-input v-model:value="searchUid" allow-clear placeholder="精确匹配数字 ID" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="10" :lg="6">
            <a-form-item label="邮箱" style="margin-bottom: 0;">
              <a-input v-model:value="searchEmail" allow-clear placeholder="模糊匹配" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8" :lg="4">
            <a-form-item label="状态" style="margin-bottom: 0;">
              <a-select v-model:value="searchStatus" :options="statusOptions" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="24" :lg="10">
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

    <a-card class="table-card seamless-card" :bordered="false" :body-style="{ padding: '0' }">
      <a-table
        :columns="columns"
        :data-source="items"
        :loading="loading"
        :pagination="false"
        :row-key="record => record.id"
        size="middle"
        :scroll="{ x: 1000 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'user_meta'">
            <div class="cell-stack">
              <strong>{{ record.email }}</strong>
              <div class="tags">
                <a-tag size="small" :color="record.is_linuxdo ? 'blue' : 'default'">
                  {{ record.is_linuxdo ? 'Lडू' : '邮箱注册' }}
                </a-tag>
                <span class="muted">{{ record.nickname }}</span>
              </div>
            </div>
          </template>

          <template v-else-if="column.dataIndex === 'status'">
            <a-tag :color="record.status === 'active' ? 'success' : 'error'">
              {{ record.status }}
            </a-tag>
          </template>

          <template v-else-if="column.dataIndex === 'credits'">
            <span style="color: #faad14; font-weight: 600;">💎 {{ record.credits }}</span>
          </template>

          <template v-else-if="column.dataIndex === 'activity'">
            <div class="cell-stack">
              <span>调用: <b>{{ record.usage_count }}</b> 次</span>
              <span class="muted">邀新: {{ record.invite_count }} 人</span>
            </div>
          </template>

          <template v-else-if="column.dataIndex === 'created_at'">
            {{ formatTime(record.created_at) }}
          </template>

          <template v-else-if="column.dataIndex === 'last_login_at'">
            <span class="muted">{{ formatTime(record.last_login_at) }}</span>
          </template>

          <template v-else-if="column.dataIndex === 'actions'">
            <div class="action-cell">
              <a-button size="small" type="primary" ghost @click="openCreditModal(record)">资产调整</a-button>
              
              <a-button 
                v-if="record.status === 'active'"
                size="small" danger @click="doChangeStatus(record, 'banned')"
              >
                封禁账号
              </a-button>
              <a-button 
                v-else
                size="small" type="primary" success @click="doChangeStatus(record, 'active')"
              >
                解冻恢复
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

    <a-modal
      v-model:open="showCreditModal"
      title="手工调整用户钻石资产"
      @ok="submitCreditUpdate"
      :confirmLoading="creditLoading"
    >
      <div v-if="targetUser" style="margin-bottom: 24px;">
        <p>正为用户 <strong>{{ targetUser.email }}</strong> 进行资产调整</p>
        <p>当前余额：<span style="color: #faad14; font-weight: bold;">💎 {{ targetUser.credits }}</span></p>
        <a-form layout="vertical">
          <a-form-item label="变动额度 (Delta)">
            <a-input-number 
              v-model:value="creditForm.delta" 
              style="width: 100%" 
              :formatter="value => `${value > 0 ? '+' : ''}${value}`"
            />
            <div class="muted" style="margin-top: 4px;">负数表示扣除。计算结果不能低于 0 余额。</div>
          </a-form-item>
          <a-form-item label="备注说明 (记录入账流水源)">
            <a-input v-model:value="creditForm.note" placeholder="例如: 参与线下社群抽奖赠送 100" />
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<style scoped>
.mb16 { margin-bottom: 24px; }
.cell-stack { display: flex; flex-direction: column; gap: 4px; }
.tags { display: flex; align-items: center; gap: 6px; }
.muted { color: rgba(0, 0, 0, 0.45); font-size: 12px; }
.action-cell { display: flex; gap: 6px; flex-wrap: wrap; }
.table-pagination { display: flex; justify-content: flex-end; padding: 16px; border-top: 1px solid #f0f0f0; }
</style>
