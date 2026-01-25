<script setup lang="ts">
import { ref, onMounted, h, computed } from 'vue'
import {
  NCard,
  NDataTable,
  NButton,
  NSpace,
  NPopconfirm,
  NModal,
  NForm,
  NFormItem,
  NInput,
  NSelect,
  NSwitch,
  NInputNumber,
  useMessage,
  NTag,
  NIcon,
  NPagination,
} from 'naive-ui'
import { AddOutline, CreateOutline, TrashOutline, PlayOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst, SelectOption } from 'naive-ui'
import { createChannel, deleteChannel, getChannels, updateChannel, testChannel } from '../api/channel'
import type { Channel } from '../types'

const message = useMessage()
const loading = ref(false)

const channels = ref<Channel[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const showModal = ref(false)
const modalTitle = ref('添加渠道')
const formRef = ref<FormInst | null>(null)
const formLoading = ref(false)
const testingId = ref<number | null>(null)

type ChannelForm = {
  id?: number
  name: string
  provider: string
  api_key?: string
  enabled: boolean
  weight: number
}

const formValue = ref<ChannelForm>({
  name: '',
  provider: '',
  api_key: '',
  enabled: true,
  weight: 1,
})

const rules = {
  name: { required: true, message: '请输入渠道名称', trigger: 'blur' },
  provider: { required: true, message: '请选择/输入 provider', trigger: 'blur' },
}

const providerOptions = computed<SelectOption[]>(() => {
  return [
    { label: 'tavily', value: 'tavily' },
    { label: 'metaso', value: 'metaso' },
    { label: 'zhipu_web_search', value: 'zhipu_web_search' },
  ]
})

const columns: DataTableColumns<Channel> = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', key: 'name', width: 160 },
  { title: 'Provider', key: 'provider', width: 160 },
  {
    title: '启用',
    key: 'enabled',
    width: 80,
    render(row: Channel) {
      return h(NTag, { type: row.enabled ? 'success' : 'error', size: 'small' }, () => (row.enabled ? '启用' : '禁用'))
    },
  },
  { title: '权重', key: 'weight', width: 80 },
  { title: 'Fail', key: 'fail_count', width: 80 },
  { title: 'LastUsed', key: 'last_used_at', width: 120 },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render(row: Channel) {
      return h(NSpace, null, () => [
        h(
          NButton,
          {
            size: 'small',
            loading: testingId.value === row.id,
            onClick: () => handleTest(row.id),
          },
          {
            icon: () => h(NIcon, null, () => h(PlayOutline)),
            default: () => '测试',
          }
        ),
        h(
          NButton,
          { size: 'small', onClick: () => handleEdit(row) },
          () => h(NIcon, null, () => h(CreateOutline))
        ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete(row.id) },
          {
            trigger: () =>
              h(
                NButton,
                { size: 'small', type: 'error' },
                () => h(NIcon, null, () => h(TrashOutline))
              ),
            default: () => '确定删除此渠道？',
          }
        ),
      ])
    },
  },
]

async function handleTest(id: number) {
  if (testingId.value) return
  testingId.value = id
  try {
    const res = await testChannel(id)
    if (res.data.success) {
      const d: any = res.data.data || {}
      message.success(`测试成功（${d.provider || ''} ${d.latency_ms ?? ''}ms）`)
    } else {
      message.error(res.data.message || '测试失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '测试失败')
  } finally {
    testingId.value = null
  }
}

async function fetchChannels() {
  loading.value = true
  try {
    const res = await getChannels({ p: page.value, page_size: pageSize.value })
    if (res.data.success && res.data.data) {
      channels.value = res.data.data.items ?? []
      total.value = res.data.data.total ?? 0
    } else {
      message.error(res.data.message || '获取渠道失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '获取渠道失败')
  } finally {
    loading.value = false
  }
}

function handleAdd() {
  modalTitle.value = '添加渠道'
  formValue.value = {
    name: '',
    provider: '',
    api_key: '',
    enabled: true,
    weight: 1,
  }
  showModal.value = true
}

function handleEdit(row: Channel) {
  modalTitle.value = '编辑渠道'
  formValue.value = {
    id: row.id,
    name: row.name,
    provider: row.provider,
    api_key: '',
    enabled: row.enabled,
    weight: row.weight ?? 1,
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
      const payload: Partial<Channel> & { id: number } = {
        id: formValue.value.id,
        name: formValue.value.name,
        provider: formValue.value.provider,
        enabled: formValue.value.enabled,
        weight: formValue.value.weight,
      }
      if (formValue.value.api_key) {
        payload.api_key = formValue.value.api_key
      }
      const res = await updateChannel(payload)
      if (res.data.success) {
        message.success('更新成功')
        showModal.value = false
        fetchChannels()
      } else {
        message.error(res.data.message || '更新失败')
      }
    } else {
      const res = await createChannel({
        name: formValue.value.name,
        provider: formValue.value.provider,
        api_key: formValue.value.api_key || '',
        enabled: formValue.value.enabled,
        weight: formValue.value.weight,
      })
      if (res.data.success) {
        message.success('创建成功')
        showModal.value = false
        fetchChannels()
      } else {
        message.error(res.data.message || '创建失败')
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
    const res = await deleteChannel(id)
    if (res.data.success) {
      message.success('删除成功')
      fetchChannels()
    } else {
      message.error(res.data.message || '删除失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '删除失败')
  }
}

function handlePageChange(p: number) {
  page.value = p
  fetchChannels()
}

function handlePageSizeChange(ps: number) {
  pageSize.value = ps
  page.value = 1
  fetchChannels()
}

onMounted(fetchChannels)
</script>

<template>
  <div>
    <NCard title="渠道管理">
      <template #header-extra>
        <NButton type="primary" @click="handleAdd">
          <template #icon><NIcon><AddOutline /></NIcon></template>
          添加渠道
        </NButton>
      </template>

      <NDataTable
        :columns="columns"
        :data="channels"
        :loading="loading"
        :row-key="(row: Channel) => row.id"
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

    <NModal v-model:show="showModal" preset="card" :title="modalTitle" style="width: 520px">
      <NForm ref="formRef" :model="formValue" :rules="rules" label-placement="left" label-width="110">
        <NFormItem label="名称" path="name">
          <NInput v-model:value="formValue.name" placeholder="请输入渠道名称" />
        </NFormItem>

        <NFormItem label="Provider" path="provider">
          <NSelect
            v-model:value="formValue.provider"
            :options="providerOptions"
            filterable
            tag
            placeholder="选择或输入 provider"
          />
        </NFormItem>

        <NFormItem label="API Key" path="api_key">
          <NInput
            v-model:value="formValue.api_key"
            type="password"
            show-password-on="click"
            :placeholder="formValue.id ? '留空表示不修改' : '请输入 API Key'"
          />
        </NFormItem>

        <NFormItem label="启用" path="enabled">
          <NSwitch v-model:value="formValue.enabled" />
        </NFormItem>

        <NFormItem label="权重" path="weight">
          <NInputNumber v-model:value="formValue.weight" :min="1" style="width: 100%" />
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
