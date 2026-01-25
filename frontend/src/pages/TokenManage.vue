<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import {
  NCard, NDataTable, NButton, NSpace, NPopconfirm, NModal,
  NForm, NFormItem, NInput, NSwitch, NInputNumber, NDatePicker,
  useMessage, NTag, NIcon, NDescriptions, NDescriptionsItem
} from 'naive-ui'
import { AddOutline, TrashOutline, CreateOutline, CopyOutline, AnalyticsOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst } from 'naive-ui'
import { getAllTokens, addToken, updateToken, deleteToken, deleteTokenBatch } from '../api/token'
import { getTokenUsageById } from '../api/usage'
import type { Token, TokenForm, TokenUsage } from '../types'

const message = useMessage()
const loading = ref(false)
const tokens = ref<Token[]>([])
const selectedRowKeys = ref<number[]>([])

const showModal = ref(false)
const modalTitle = ref('添加Token')
const formRef = ref<FormInst | null>(null)
const formLoading = ref(false)
const formValue = ref<TokenForm>({
  name: '',
  unlimited_quota: false,
  remain_quota: 1000,
  group: '',
})

const rules = {
  name: { required: true, message: '请输入Token名称', trigger: 'blur' },
}

const usageModal = ref(false)
const usageLoading = ref(false)
const usageData = ref<TokenUsage | null>(null)
const usageTitle = ref('Token 用量')

const columns: DataTableColumns<Token> = [
  { type: 'selection' },
  { title: '名称', key: 'name', width: 150 },
  { title: '分组', key: 'group', width: 120 },
  {
    title: 'Key',
    key: 'key',
    width: 280,
    render(row: Token) {
      return h(NSpace, { align: 'center' }, () => [
        h('code', { style: 'font-size: 12px' }, row.key?.slice(0, 24) + '...'),
        h(NButton, {
          size: 'tiny',
          quaternary: true,
          onClick: () => copyKey(row.key)
        }, () => h(NIcon, null, () => h(CopyOutline)))
      ])
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render(row: Token) {
      return h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small' },
        () => row.status === 1 ? '启用' : '禁用')
    }
  },
  {
    title: '配额',
    key: 'quota',
    width: 150,
    render(row: Token) {
      if (row.unlimited_quota) {
        return h(NTag, { type: 'info', size: 'small' }, () => '无限制')
      }
      return `${row.used_quota} / ${row.remain_quota + row.used_quota}`
    }
  },
  { title: '过期时间', key: 'expired_at', width: 180 },
  { title: '创建时间', key: 'created_at', width: 180 },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render(row: Token) {
      return h(NSpace, null, () => [
        h(NButton, { size: 'small', onClick: () => handleUsage(row) }, () => [
          h(NIcon, null, () => h(AnalyticsOutline)),
        ]),
        h(NButton, { size: 'small', onClick: () => handleEdit(row) }, () => [
          h(NIcon, null, () => h(CreateOutline)),
        ]),
        h(NPopconfirm, {
          onPositiveClick: () => handleDelete(row.id)
        }, {
          trigger: () => h(NButton, { size: 'small', type: 'error' }, () => [
            h(NIcon, null, () => h(TrashOutline)),
          ]),
          default: () => '确定删除此Token?'
        })
      ])
    }
  }
]

function copyKey(key: string) {
  navigator.clipboard.writeText(key)
  message.success('已复制到剪贴板')
}

async function fetchTokens() {
  loading.value = true
  try {
    const res = await getAllTokens()
    if (res.data.success && res.data.data) {
      tokens.value = res.data.data.items ?? []
    }
  } finally {
    loading.value = false
  }
}

function handleAdd() {
  modalTitle.value = '添加Token'
  formValue.value = {
    name: '',
    unlimited_quota: false,
    remain_quota: 1000,
    group: '',
  }
  showModal.value = true
}

function handleEdit(row: Token) {
  modalTitle.value = '编辑Token'
  formValue.value = {
    id: row.id,
    name: row.name,
    unlimited_quota: row.unlimited_quota,
    remain_quota: row.remain_quota,
    expired_at: row.expired_at,
    group: row.group || '',
  }
  showModal.value = true
}

async function handleUsage(row: Token) {
  usageTitle.value = `Token 用量：${row.name}`
  usageModal.value = true
  usageLoading.value = true
  usageData.value = null
  try {
    const res = await getTokenUsageById(row.id)
    if (res.data.success) {
      usageData.value = res.data.data ?? null
    } else {
      message.error(res.data.message || '获取用量失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '获取用量失败')
  } finally {
    usageLoading.value = false
  }
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  formLoading.value = true
  try {
    if (formValue.value.id) {
      const res = await updateToken(formValue.value)
      if (res.data.success) {
        message.success('更新成功')
        showModal.value = false
        fetchTokens()
      } else {
        message.error(res.data.message || '更新失败')
      }
    } else {
      const res = await addToken(formValue.value)
      if (res.data.success) {
        message.success('添加成功')
        showModal.value = false
        fetchTokens()
      } else {
        message.error(res.data.message || '添加失败')
      }
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '操作失败')
  } finally {
    formLoading.value = false
  }
}

async function handleDelete(id: number) {
  try {
    const res = await deleteToken(id)
    if (res.data.success) {
      message.success('删除成功')
      fetchTokens()
    } else {
      message.error(res.data.message || '删除失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '删除失败')
  }
}

async function handleBatchDelete() {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要删除的Token')
    return
  }
  try {
    const res = await deleteTokenBatch(selectedRowKeys.value)
    if (res.data.success) {
      message.success('批量删除成功')
      selectedRowKeys.value = []
      fetchTokens()
    } else {
      message.error(res.data.message || '批量删除失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '批量删除失败')
  }
}

onMounted(fetchTokens)
</script>

<template>
  <div>
    <NCard title="令牌管理">
      <template #header-extra>
        <NSpace>
          <NPopconfirm @positive-click="handleBatchDelete">
            <template #trigger>
              <NButton type="error" :disabled="selectedRowKeys.length === 0">
                <template #icon><NIcon><TrashOutline /></NIcon></template>
                批量删除
              </NButton>
            </template>
            确定删除选中的 {{ selectedRowKeys.length }} 个令牌?
          </NPopconfirm>
          <NButton type="primary" @click="handleAdd">
            <template #icon><NIcon><AddOutline /></NIcon></template>
            添加令牌
          </NButton>
        </NSpace>
      </template>

      <NDataTable
        :columns="columns"
        :data="tokens"
        :loading="loading"
        :row-key="(row: Token) => row.id"
        v-model:checked-row-keys="selectedRowKeys"
      />
    </NCard>

    <NModal v-model:show="showModal" preset="card" :title="modalTitle" style="width: 500px">
      <NForm ref="formRef" :model="formValue" :rules="rules" label-placement="left" label-width="100">
        <NFormItem label="名称" path="name">
          <NInput v-model:value="formValue.name" placeholder="请输入Token名称" />
        </NFormItem>
        <NFormItem label="分组" path="group">
          <NInput v-model:value="formValue.group" placeholder="可选：用于归类 Token" />
        </NFormItem>
        <NFormItem label="无限配额" path="unlimited_quota">
          <NSwitch v-model:value="formValue.unlimited_quota" />
        </NFormItem>
        <NFormItem v-if="!formValue.unlimited_quota" label="配额" path="remain_quota">
          <NInputNumber v-model:value="formValue.remain_quota" :min="0" style="width: 100%" />
        </NFormItem>
        <NFormItem label="过期时间" path="expired_at">
          <NDatePicker v-model:formatted-value="formValue.expired_at" type="datetime" clearable style="width: 100%" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="showModal = false">取消</NButton>
          <NButton type="primary" :loading="formLoading" @click="handleSubmit">确定</NButton>
        </NSpace>
      </template>
    </NModal>

    <NModal v-model:show="usageModal" preset="card" :title="usageTitle" style="width: 520px">
      <div v-if="usageLoading">加载中...</div>
      <template v-else>
        <NDescriptions v-if="usageData" :column="1" label-placement="left" bordered>
          <NDescriptionsItem label="名称">{{ usageData.name }}</NDescriptionsItem>
          <NDescriptionsItem label="Total Granted">{{ usageData.total_granted }}</NDescriptionsItem>
          <NDescriptionsItem label="Total Used">{{ usageData.total_used }}</NDescriptionsItem>
          <NDescriptionsItem label="Total Available">{{ usageData.total_available }}</NDescriptionsItem>
          <NDescriptionsItem label="Unlimited">{{ usageData.unlimited_quota }}</NDescriptionsItem>
          <NDescriptionsItem label="Expires At">{{ usageData.expires_at }}</NDescriptionsItem>
        </NDescriptions>
        <div v-else>暂无数据</div>
      </template>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="usageModal = false">关闭</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>
