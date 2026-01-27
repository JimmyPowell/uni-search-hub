<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import {
  NCard, NDataTable, NButton, NSpace, NPopconfirm, NModal,
  NForm, NFormItem, NInput, NSwitch, NInputNumber, NDatePicker,
  useMessage, NTag, NIcon, NAlert
} from 'naive-ui'
import { AddOutline, TrashOutline, CreateOutline, CopyOutline, ChevronForwardOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst } from 'naive-ui'
import { getAllRedemptions, addRedemption, updateRedemption, deleteRedemption, deleteInvalidRedemptions } from '../api/redemption'
import type { Redemption, RedemptionForm } from '../types'

const message = useMessage()
const loading = ref(false)
const redemptions = ref<Redemption[]>([])

const showModal = ref(false)
const modalTitle = ref('添加兑换码')
const formRef = ref<FormInst | null>(null)
const formLoading = ref(false)
const createdKeys = ref<string[]>([])

const showKeysModal = ref(false)
const formValue = ref<RedemptionForm>({
  name: '',
  quota: 100,
  count: 1,
  expired_time: undefined,
})

const rules = {
  name: { required: true, message: '请输入兑换码名称', trigger: 'blur' },
  quota: { required: true, type: 'number', message: '请输入配额', trigger: 'blur' },
  count: { required: true, type: 'number', message: '请输入生成数量', trigger: 'blur' },
}

const columns: DataTableColumns<Redemption> = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '名称', key: 'name', width: 150 },
  {
    title: '兑换码',
    key: 'key',
    width: 280,
    render(row: Redemption) {
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
    width: 100,
    render(row: Redemption) {
      const statusMap: Record<number, { type: 'success' | 'info' | 'error', text: string }> = {
        1: { type: 'success', text: '未使用' },
        2: { type: 'info', text: '已使用' },
        3: { type: 'error', text: '已禁用' },
      }
      const status = statusMap[row.status] || { type: 'error', text: '未知' }
      return h(NTag, { type: status.type, size: 'small' }, () => status.text)
    }
  },
  { title: '配额', key: 'quota', width: 100 },
  {
    title: '过期时间',
    key: 'expired_time',
    width: 180,
    render(row: Redemption) {
      return row.expired_time === 0 ? '永久有效' : new Date(row.expired_time * 1000).toLocaleString()
    }
  },
  {
    title: '创建时间',
    key: 'created_time',
    width: 180,
    render(row: Redemption) {
      return new Date(row.created_time * 1000).toLocaleString()
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render(row: Redemption) {
      return h(NSpace, null, () => [
        h(NButton, {
          size: 'small',
          disabled: row.status !== 1,
          onClick: () => handleEdit(row)
        }, () => [
          h(NIcon, null, () => h(CreateOutline)),
        ]),
        h(NPopconfirm, {
          onPositiveClick: () => handleDelete(row.id)
        }, {
          trigger: () => h(NButton, { size: 'small', type: 'error' }, () => [
            h(NIcon, null, () => h(TrashOutline)),
          ]),
          default: () => '确定删除此兑换码?'
        })
      ])
    }
  }
]

function copyKey(key: string) {
  navigator.clipboard.writeText(key)
  message.success('已复制到剪贴板')
}

async function fetchRedemptions() {
  loading.value = true
  try {
    const res = await getAllRedemptions()
    if (res.data.success && res.data.data) {
      redemptions.value = res.data.data.items ?? []
    }
  } finally {
    loading.value = false
  }
}

function handleAdd() {
  modalTitle.value = '批量生成兑换码'
  formValue.value = {
    name: '',
    quota: 100,
    count: 1,
    expired_time: undefined,
  }
  showModal.value = true
}

function handleEdit(row: Redemption) {
  modalTitle.value = '编辑兑换码'
  formValue.value = {
    id: row.id,
    name: row.name,
    quota: row.quota,
    count: 1,
    expired_time: row.expired_time ? row.expired_time * 1000 : undefined,
    status: row.status,
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
    // 转换时间戳：毫秒 -> 秒
    const submitData = {
      ...formValue.value,
      expired_time: formValue.value.expired_time
        ? Math.floor(formValue.value.expired_time / 1000)
        : undefined
    }

    if (formValue.value.id) {
      const res = await updateRedemption(submitData)
      if (res.data.success) {
        message.success('更新成功')
        showModal.value = false
        fetchRedemptions()
      } else {
        message.error(res.data.message || '更新失败')
      }
    } else {
      const res = await addRedemption(submitData)
      if (res.data.success) {
        message.success('兑换码创建成功')
        createdKeys.value = res.data.data ?? []
        showKeysModal.value = true
        showModal.value = false
        fetchRedemptions()
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
    const res = await deleteRedemption(id)
    if (res.data.success) {
      message.success('删除成功')
      fetchRedemptions()
    } else {
      message.error(res.data.message || '删除失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '删除失败')
  }
}

async function handleDeleteInvalid() {
  try {
    const res = await deleteInvalidRedemptions()
    if (res.data.success) {
      message.success(`已删除 ${res.data.data ?? 0} 个无效兑换码`)
      fetchRedemptions()
    } else {
      message.error(res.data.message || '删除失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '删除失败')
  }
}

onMounted(fetchRedemptions)
</script>

<template>
  <div>
    <NCard title="兑换码管理">
      <template #header-extra>
        <NSpace>
          <NPopconfirm @positive-click="handleDeleteInvalid">
            <template #trigger>
              <NButton>
                <template #icon><NIcon><TrashOutline /></NIcon></template>
                清理无效兑换码
              </NButton>
            </template>
            确定删除所有已使用和已过期的兑换码?
          </NPopconfirm>
          <NButton type="primary" @click="handleAdd">
            <template #icon><NIcon><AddOutline /></NIcon></template>
            批量生成兑换码
          </NButton>
        </NSpace>
      </template>

      <NDataTable
        :columns="columns"
        :data="redemptions"
        :loading="loading"
        :row-key="(row: Redemption) => row.id"
      />
    </NCard>

    <NModal v-model:show="showModal" preset="card" :title="modalTitle" style="width: 500px">
      <NForm ref="formRef" :model="formValue" :rules="rules" label-placement="left" label-width="100">
        <NFormItem label="名称" path="name">
          <NInput v-model:value="formValue.name" placeholder="请输入兑换码名称" />
        </NFormItem>
        <NFormItem label="配额" path="quota">
          <NInputNumber v-model:value="formValue.quota" :min="1" style="width: 100%" />
        </NFormItem>
        <NFormItem v-if="!formValue.id" label="生成数量" path="count">
          <NInputNumber v-model:value="formValue.count" :min="1" :max="100" style="width: 100%" />
        </NFormItem>
        <NFormItem label="过期时间" path="expired_time">
          <NDatePicker
            v-model:formatted-value="formValue.expired_time"
            type="datetime"
            clearable
            style="width: 100%"
            value-type="timestamp"
            :is-date-disabled="(timestamp: number) => timestamp < Date.now()"
          />
        </NFormItem>
        <NFormItem v-if="formValue.id" label="状态" path="status">
          <NSwitch
            v-model:value="formValue.status"
            :checked-value="1"
            :unchecked-value="3"
          >
            <template #checked>启用</template>
            <template #unchecked>禁用</template>
          </NSwitch>
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace justify="end">
          <NButton @click="showModal = false">取消</NButton>
          <NButton type="primary" :loading="formLoading" @click="handleSubmit">确定</NButton>
        </NSpace>
      </template>
    </NModal>

    <NModal v-model:show="showKeysModal" preset="card" title="兑换码创建成功" style="width: 600px">
      <NAlert type="success" style="margin-bottom: 16px">
        已成功生成 {{ createdKeys.length }} 个兑换码，请及时保存！
      </NAlert>
      <div style="max-height: 400px; overflow-y: auto;">
        <div v-for="(key, index) in createdKeys" :key="key" style="margin-bottom: 8px;">
          <NSpace align="center">
            <span style="color: #999; width: 30px;">{{ index + 1 }}.</span>
            <code>{{ key }}</code>
            <NButton size="tiny" quaternary @click="copyKey(key)">
              <template #icon><NIcon><CopyOutline /></NIcon></template>
            </NButton>
          </NSpace>
        </div>
      </div>
      <template #footer>
        <NSpace justify="end">
          <NButton type="primary" @click="showKeysModal = false">关闭</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>