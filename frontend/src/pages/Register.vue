<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NForm, NFormItem, NInput, NButton, NSpace, NA, useMessage } from 'naive-ui'
import type { FormInst } from 'naive-ui'
import { register } from '../api/user'

const router = useRouter()
const message = useMessage()
const formRef = ref<FormInst | null>(null)
const loading = ref(false)

const formValue = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  display_name: '',
})

const rules = {
  username: {
    required: true,
    message: '请输入用户名',
    trigger: 'blur',
  },
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
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string) => {
        if (value !== formValue.value.password) {
          return new Error('两次密码输入不一致')
        }
        return true
      },
      trigger: 'blur',
    },
  ],
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  loading.value = true
  try {
    const res = await register({
      username: formValue.value.username,
      email: formValue.value.email,
      password: formValue.value.password,
      display_name: formValue.value.display_name || undefined,
    })
    if (res.data.success) {
      message.success('注册成功，请登录')
      router.push('/login')
    } else {
      message.error(res.data.message || '注册失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '注册失败')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="register-container">
    <NCard title="注册" class="register-card">
      <NForm ref="formRef" :model="formValue" :rules="rules">
        <NFormItem label="用户名" path="username">
          <NInput v-model:value="formValue.username" placeholder="请输入用户名" />
        </NFormItem>
        <NFormItem label="邮箱" path="email">
          <NInput v-model:value="formValue.email" placeholder="请输入邮箱" />
        </NFormItem>
        <NFormItem label="显示名称" path="display_name">
          <NInput v-model:value="formValue.display_name" placeholder="可选，用于显示的名称" />
        </NFormItem>
        <NFormItem label="密码" path="password">
          <NInput
            v-model:value="formValue.password"
            type="password"
            show-password-on="click"
            placeholder="请输入密码"
          />
        </NFormItem>
        <NFormItem label="确认密码" path="confirmPassword">
          <NInput
            v-model:value="formValue.confirmPassword"
            type="password"
            show-password-on="click"
            placeholder="请再次输入密码"
          />
        </NFormItem>
        <NSpace vertical :size="16" style="width: 100%">
          <NButton type="primary" block :loading="loading" @click="handleSubmit">
            注册
          </NButton>
          <div style="text-align: center">
            已有账号？
            <NA @click="router.push('/login')">立即登录</NA>
          </div>
        </NSpace>
      </NForm>
    </NCard>
  </div>
</template>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-card {
  width: 420px;
}
</style>
