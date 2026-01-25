<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import {
  NCard, NDataTable, NButton, NSpace, NPopconfirm, NModal,
  NForm, NFormItem, NInput, NSelect,
  useMessage, NTag, NIcon
} from 'naive-ui'
import { AddOutline, TrashOutline, CreateOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst } from 'naive-ui'
import { getAllUsers, createUser, updateUser, deleteUser } from '../api/user'
import type { User } from '../types'
import { UserRole } from '../types'

const message = useMessage()
const loading = ref(false)
const users = ref<User[]>([])

const showModal = ref(false)
const modalTitle = ref('添加用户')
const formRef = ref<FormInst | null>(null)
const formLoading = ref(false)
const isEdit = ref(false)

type UserForm = {
  id: number
  username: string
  email: string
  display_name: string
  password: string
  role: number
  status: number
}

const formValue = ref<UserForm>({
  id: 0,
  username: '',
  email: '',
  display_name: '',
  password: '',
  role: UserRole.User,
  status: 1,
})

const rules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    {
      validator: (_rule: any, value: string) => {
        const ok = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)
        return ok ? true : new Error('请输入有效的邮箱地址')
      },
      trigger: 'blur',
    },
  ],
  password: {
    required: true,
    message: '请输入密码',
    trigger: 'blur',
    validator: (_rule: any, value: string) => {
      if (!isEdit.value && !value) {
        return new Error('请输入密码')
      }
      return true
    }
  },
}

const roleOptions = [
  { label: '普通用户', value: UserRole.User },
  { label: '管理员', value: UserRole.Admin },
  { label: 'Root', value: UserRole.Root },
]

const statusOptions = [
  { label: '正常', value: 1 },
  { label: '禁用', value: 2 },
]

const columns: DataTableColumns<User> = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '用户名', key: 'username', width: 120 },
  { title: '显示名称', key: 'display_name', width: 120 },
  { title: '邮箱', key: 'email', width: 200 },
  {
    title: '角色',
    key: 'role',
    width: 100,
    render(row: User) {
      const roleMap: Record<number, { type: 'success' | 'warning' | 'error', text: string }> = {
        [UserRole.User]: { type: 'success', text: '用户' },
        [UserRole.Admin]: { type: 'warning', text: '管理员' },
        [UserRole.Root]: { type: 'error', text: 'Root' },
      }
      const r = roleMap[row.role] || { type: 'success', text: '用户' }
      return h(NTag, { type: r.type, size: 'small' }, () => r.text)
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render(row: User) {
      return h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small' },
        () => row.status === 1 ? '正常' : '禁用')
    }
  },
  { title: '创建时间', key: 'created_at', width: 180 },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render(row: User) {
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
          default: () => '确定删除此用户?'
        })
      ])
    }
  }
]

async function fetchUsers() {
  loading.value = true
  try {
    const res = await getAllUsers()
    if (res.data.success && res.data.data) {
      users.value = res.data.data.items ?? []
    }
  } finally {
    loading.value = false
  }
}

function handleAdd() {
  isEdit.value = false
  modalTitle.value = '添加用户'
  formValue.value = {
    id: 0,
    username: '',
    email: '',
    display_name: '',
    password: '',
    role: UserRole.User,
    status: 1,
  }
  showModal.value = true
}

function handleEdit(row: User) {
  isEdit.value = true
  modalTitle.value = '编辑用户'
  formValue.value = {
    id: row.id,
    username: row.username,
    email: row.email,
    display_name: row.display_name || '',
    password: '',
    role: row.role,
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
    if (isEdit.value) {
      const data: any = {
        id: formValue.value.id,
        email: formValue.value.email,
        display_name: formValue.value.display_name,
        role: formValue.value.role,
        status: formValue.value.status,
      }
      if (formValue.value.password) {
        data.password = formValue.value.password
      }
      const res = await updateUser(data)
      if (res.data.success) {
        message.success('更新成功')
        showModal.value = false
        fetchUsers()
      } else {
        message.error(res.data.message || '更新失败')
      }
    } else {
      const res = await createUser({
        username: formValue.value.username,
        email: formValue.value.email,
        display_name: formValue.value.display_name,
        password: formValue.value.password,
        role: formValue.value.role,
        status: formValue.value.status,
      })
      if (res.data.success) {
        message.success('添加成功')
        showModal.value = false
        fetchUsers()
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
    const res = await deleteUser(id)
    if (res.data.success) {
      message.success('删除成功')
      fetchUsers()
    } else {
      message.error(res.data.message || '删除失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '删除失败')
  }
}

onMounted(fetchUsers)
</script>

<template>
  <div>
    <NCard title="用户管理">
      <template #header-extra>
        <NButton type="primary" @click="handleAdd">
          <template #icon><NIcon><AddOutline /></NIcon></template>
          添加用户
        </NButton>
      </template>

      <NDataTable
        :columns="columns"
        :data="users"
        :loading="loading"
        :row-key="(row: User) => row.id"
      />
    </NCard>

    <NModal v-model:show="showModal" preset="card" :title="modalTitle" style="width: 500px">
      <NForm ref="formRef" :model="formValue" :rules="rules" label-placement="left" label-width="100">
        <NFormItem label="用户名" path="username">
          <NInput v-model:value="formValue.username" :disabled="isEdit" placeholder="请输入用户名" />
        </NFormItem>
        <NFormItem label="邮箱" path="email">
          <NInput v-model:value="formValue.email" placeholder="请输入邮箱" />
        </NFormItem>
        <NFormItem label="显示名称" path="display_name">
          <NInput v-model:value="formValue.display_name" placeholder="请输入显示名称" />
        </NFormItem>
        <NFormItem :label="isEdit ? '新密码' : '密码'" path="password">
          <NInput v-model:value="formValue.password" type="password" show-password-on="click" :placeholder="isEdit ? '留空则不修改' : '请输入密码'" />
        </NFormItem>
        <NFormItem label="角色" path="role">
          <NSelect v-model:value="formValue.role" :options="roleOptions" />
        </NFormItem>
        <NFormItem label="状态" path="status">
          <NSelect v-model:value="formValue.status" :options="statusOptions" />
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
