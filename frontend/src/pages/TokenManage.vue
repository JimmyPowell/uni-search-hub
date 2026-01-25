<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import {
  NCard, NDataTable, NButton, NSpace, NPopconfirm, NModal,
  NForm, NFormItem, NInput, NSwitch, NInputNumber, NDatePicker,
  useMessage, NTag, NIcon
} from 'naive-ui'
import { AddOutline, TrashOutline, CreateOutline, CopyOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst } from 'naive-ui'
import { getAllTokens, addToken, updateToken, deleteToken, deleteTokenBatch } from '../api/token'
import type { Token, TokenForm } from '../types'

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
})

const rules = {
  name: { required: true, message: '请输入Token名称', trigger: 'blur' },
}

const columns: DataTableColumns<Token> = [
  { type: 'selection' },
  { title: '名称', key: 'name', width: 150 },
  {
    title: 'Key',
    key: 'key',
    width: 280,
    render(row) {
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
    render(row) {
      return h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small' },
        () => row.status === 1 ? '启用' : '禁用')
    }
  },
  {
    title: '配额',
    key: 'quota',
    width: 150,
    render(row) {
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
    render(row) {
      return h(NSpace, null, () => [
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
  }
  showModal.value = true
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
    <NCard title="Token管理">
      <template #header-extra>
        <NSpace>
          <NPopconfirm @positive-click="handleBatchDelete">
            <template #trigger>
              <NButton type="error" :disabled="selectedRowKeys.length === 0">
                <template #icon><NIcon><TrashOutline /></NIcon></template>
                批量删除
              </NButton>
            </template>
            确定删除选中的 {{ selectedRowKeys.length }} 个Token?
          </NPopconfirm>
          <NButton type="primary" @click="handleAdd">
            <template #icon><NIcon><AddOutline /></NIcon></template>
            添加Token
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
  </div>
</template>
