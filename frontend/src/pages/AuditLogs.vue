<script setup lang="ts">
import { ref, onMounted, h, computed } from 'vue'
import {
  NCard,
  NDataTable,
  NButton,
  NSpace,
  NInput,
  NInputNumber,
  NDatePicker,
  useMessage,
  NTag,
  NPagination,
  NSelect,
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getRequestLogs } from '../api/request_log'
import type { RequestLog } from '../types'
import { useUserStore } from '../stores/user'

const message = useMessage()
const userStore = useUserStore()

function formatDateTimeSeconds(value: string) {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

const loading = ref(false)
const logs = ref<RequestLog[]>([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const provider = ref('')
const endpoint = ref('')
const requestId = ref('')
const statusCode = ref<number | null>(null)
const tokenId = ref<number | null>(null)
const channelId = ref<number | null>(null)
const userId = ref<number | null>(null)
const action = ref('')

// Unix ms range
const timeRange = ref<[number, number] | null>(null)

const isAdmin = computed(() => userStore.isAdmin)

const actionOptions = [
  { label: '全部', value: '' },
  { label: '请求', value: 'proxy_request' },
  { label: '测试', value: 'channel_test' },
  { label: '充值', value: 'token_charge' },
]

const columns: DataTableColumns<RequestLog> = [
  {
    title: '时间',
    key: 'created_at',
    width: 170,
    render(row: RequestLog) {
      const text = formatDateTimeSeconds(row.created_at)
      return h('span', { title: row.created_at }, text)
    },
  },
  {
    title: '类型',
    key: 'action',
    width: 90,
    render(row: RequestLog) {
      const map: Record<string, { type: 'success' | 'warning' | 'info' | 'error'; text: string }> = {
        proxy_request: { type: 'success', text: '请求' },
        channel_test: { type: 'info', text: '测试' },
        token_charge: { type: 'warning', text: '充值' },
      }
      const v = row.action || 'proxy_request'
      const t = map[v] || { type: 'info', text: v }
      return h(NTag, { type: t.type as any, size: 'small' }, () => t.text)
    },
  },
  { title: 'Provider', key: 'provider', width: 140 },
  { title: 'Endpoint', key: 'endpoint', width: 180 },
  {
    title: '状态码',
    key: 'status_code',
    width: 90,
    render(row: RequestLog) {
      const t = row.status_code >= 200 && row.status_code < 300 ? 'success' : row.status_code >= 500 ? 'error' : 'warning'
      return h(NTag, { type: t, size: 'small' }, () => String(row.status_code))
    },
  },
  { title: '耗时(ms)', key: 'latency_ms', width: 100 },
  { title: 'Cost', key: 'cost', width: 70 },
  { title: 'User', key: 'user_id', width: 80 },
  { title: 'Token', key: 'token_id', width: 80 },
  { title: 'Channel', key: 'channel_id', width: 80 },
  { title: 'RequestID', key: 'request_id', width: 220 },
  {
    title: '错误',
    key: 'error',
    width: 240,
    render(row: RequestLog) {
      const text = row.error || ''
      const short = text.length > 120 ? text.slice(0, 120) + '...' : text
      return h('span', { title: text }, short)
    },
  },
]

async function fetchLogs() {
  loading.value = true
  try {
    const [start, end] = timeRange.value ?? [undefined, undefined]
    const res = await getRequestLogs({
      p: page.value,
      page_size: pageSize.value,
      provider: provider.value || undefined,
      action: action.value || undefined,
      endpoint: endpoint.value || undefined,
      request_id: requestId.value || undefined,
      status_code: statusCode.value ?? undefined,
      token_id: tokenId.value ?? undefined,
      channel_id: channelId.value ?? undefined,
      user_id: isAdmin.value ? userId.value ?? undefined : undefined,
      start_time: start,
      end_time: end,
    })
    if (res.data.success && res.data.data) {
      logs.value = res.data.data.items ?? []
      total.value = res.data.data.total ?? 0
    } else {
      message.error(res.data.message || '获取日志失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '获取日志失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  fetchLogs()
}

function handleReset() {
  provider.value = ''
  action.value = ''
  endpoint.value = ''
  requestId.value = ''
  statusCode.value = null
  tokenId.value = null
  channelId.value = null
  userId.value = null
  timeRange.value = null
  handleSearch()
}

function handlePageChange(p: number) {
  page.value = p
  fetchLogs()
}

function handlePageSizeChange(ps: number) {
  pageSize.value = ps
  page.value = 1
  fetchLogs()
}

onMounted(fetchLogs)
</script>

<template>
  <div>
    <NCard title="请求日志">
      <NSpace vertical :size="12" style="margin-bottom: 12px">
        <NSpace wrap>
          <NInput v-model:value="provider" placeholder="provider" style="width: 180px" />
          <NSelect v-model:value="action" :options="actionOptions" placeholder="type" style="width: 160px" />
          <NInput v-model:value="endpoint" placeholder="endpoint" style="width: 220px" />
          <NInput v-model:value="requestId" placeholder="request_id" style="width: 220px" />
          <NInputNumber v-model:value="statusCode" placeholder="status_code" :min="100" :max="599" style="width: 140px" />
          <NInputNumber v-model:value="tokenId" placeholder="token_id" :min="1" style="width: 140px" />
          <NInputNumber v-model:value="channelId" placeholder="channel_id" :min="1" style="width: 140px" />
          <NInputNumber
            v-if="isAdmin"
            v-model:value="userId"
            placeholder="user_id"
            :min="1"
            style="width: 140px"
          />
          <NDatePicker
            v-model:value="timeRange"
            type="datetimerange"
            clearable
            format="yyyy-MM-dd HH:mm:ss"
            :time-picker-props="{ format: 'HH:mm:ss' }"
            style="width: 320px"
          />
        </NSpace>
        <NSpace>
          <NButton type="primary" @click="handleSearch">查询</NButton>
          <NButton @click="handleReset">重置</NButton>
        </NSpace>
      </NSpace>

      <NDataTable
        :columns="columns"
        :data="logs"
        :loading="loading"
        :row-key="(row: RequestLog) => row.id"
        :scroll-x="1500"
      />

      <div style="display: flex; justify-content: flex-end; margin-top: 12px;">
        <NPagination
          :page="page"
          :page-size="pageSize"
          :item-count="total"
          show-size-picker
          :page-sizes="[10, 20, 50, 100]"
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </NCard>
  </div>
</template>
