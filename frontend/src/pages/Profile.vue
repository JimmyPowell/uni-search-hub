<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NCard, NForm, NFormItem, NInput, NButton, NPopconfirm, NInputNumber, useMessage, NAlert } from 'naive-ui'
import type { FormInst } from 'naive-ui'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { updateSelf, deleteSelf, topUp } from '../api/user'

const router = useRouter()
const message = useMessage()
const userStore = useUserStore()
const formRef = ref<FormInst | null>(null)
const passwordFormRef = ref<FormInst | null>(null)
const loading = ref(false)
const passwordLoading = ref(false)

const formValue = ref({
  username: '',
  email: '',
  display_name: '',
})

const passwordForm = ref({
  password: '',
  confirmPassword: '',
})

const redemptionForm = ref({
  key: '',
})
const redemptionFormRef = ref<FormInst | null>(null)
const redemptionLoading = ref(false)

const rules = {
  email: [
    {
      validator: (_rule: any, value: string) => {
        if (!value) return true
        const ok = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)
        return ok ? true : new Error('请输入有效的邮箱地址')
      },
      trigger: 'blur',
    },
  ],
}

const passwordRules = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string) => {
        if (value !== passwordForm.value.password) {
          return new Error('两次密码输入不一致')
        }
        return true
      },
      trigger: 'blur',
    },
  ],
}

onMounted(() => {
  if (userStore.user) {
    formValue.value = {
      username: userStore.user.username,
      email: userStore.user.email,
      display_name: userStore.user.display_name || '',
    }
  }
})

async function handleUpdateProfile() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  loading.value = true
  try {
    const res = await updateSelf({
      email: formValue.value.email,
      display_name: formValue.value.display_name,
    })
    if (res.data.success) {
      message.success('更新成功')
      await userStore.fetchUser()
    } else {
      message.error(res.data.message || '更新失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '更新失败')
  } finally {
    loading.value = false
  }
}

async function handleChangePassword() {
  try {
    await passwordFormRef.value?.validate()
  } catch {
    return
  }

  passwordLoading.value = true
  try {
    const res = await updateSelf({
      password: passwordForm.value.password,
    })
    if (res.data.success) {
      message.success('密码修改成功')
      passwordForm.value = { password: '', confirmPassword: '' }
    } else {
      message.error(res.data.message || '密码修改失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '密码修改失败')
  } finally {
    passwordLoading.value = false
  }
}

async function handleDeleteAccount() {
  try {
    const res = await deleteSelf()
    if (res.data.success) {
      message.success('账户已删除')
      userStore.clearUser()
      router.push('/login')
    } else {
      message.error(res.data.message || '删除失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '删除失败')
  }
}

async function handleRedemption() {
  try {
    await redemptionFormRef.value?.validate()
  } catch {
    return
  }

  redemptionLoading.value = true
  try {
    const res = await topUp(redemptionForm.value.key)
    if (res.data.success) {
      message.success(`充值成功，获得 ${res.data.data ?? 0} 配额`)
      redemptionForm.value.key = ''
      await userStore.fetchUser()
    } else {
      message.error(res.data.message || '充值失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '充值失败')
  } finally {
    redemptionLoading.value = false
  }
}

const redemptionRules = {
  key: { required: true, message: '请输入兑换码', trigger: 'blur' },
}
</script>

<template>
  <div class="profile-page">
    <NCard title="个人资料">
      <NForm ref="formRef" :model="formValue" :rules="rules" label-placement="left" label-width="100">
        <NFormItem label="用户名">
          <NInput v-model:value="formValue.username" disabled />
        </NFormItem>
        <NFormItem label="邮箱" path="email">
          <NInput v-model:value="formValue.email" placeholder="请输入邮箱" />
        </NFormItem>
        <NFormItem label="显示名称" path="display_name">
          <NInput v-model:value="formValue.display_name" placeholder="请输入显示名称" />
        </NFormItem>
        <NFormItem>
          <NButton type="primary" :loading="loading" @click="handleUpdateProfile">保存修改</NButton>
        </NFormItem>
      </NForm>
    </NCard>

    <NCard title="修改密码" style="margin-top: 16px">
      <NForm ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-placement="left" label-width="100">
        <NFormItem label="新密码" path="password">
          <NInput v-model:value="passwordForm.password" type="password" show-password-on="click" placeholder="请输入新密码" />
        </NFormItem>
        <NFormItem label="确认密码" path="confirmPassword">
          <NInput v-model:value="passwordForm.confirmPassword" type="password" show-password-on="click" placeholder="请再次输入新密码" />
        </NFormItem>
        <NFormItem>
          <NButton type="primary" :loading="passwordLoading" @click="handleChangePassword">修改密码</NButton>
        </NFormItem>
      </NForm>
    </NCard>

    <NCard title="兑换码充值" style="margin-top: 16px">
      <NAlert type="info" style="margin-bottom: 16px">
        如果您有兑换码，请在此输入以获得配额。
      </NAlert>
      <NForm ref="redemptionFormRef" :model="redemptionForm" :rules="redemptionRules" label-placement="left" label-width="100">
        <NFormItem label="兑换码" path="key">
          <NInput v-model:value="redemptionForm.key" placeholder="请输入兑换码" />
        </NFormItem>
        <NFormItem>
          <NButton type="primary" :loading="redemptionLoading" @click="handleRedemption">兑换</NButton>
        </NFormItem>
      </NForm>
    </NCard>

    <NCard title="危险操作" style="margin-top: 16px">
      <p style="color: #999; margin-bottom: 16px">删除账户后，所有数据将被永久清除且无法恢复。</p>
      <NPopconfirm @positive-click="handleDeleteAccount">
        <template #trigger>
          <NButton type="error">删除账户</NButton>
        </template>
        确定要删除您的账户吗？此操作不可逆！
      </NPopconfirm>
    </NCard>
  </div>
</template>

<style scoped>
.profile-page {
  max-width: 600px;
}
</style>
