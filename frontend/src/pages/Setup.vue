<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NForm, NFormItem, NInput, NButton, NSpace, useMessage } from 'naive-ui'
import type { FormInst } from 'naive-ui'
import { getSetupStatus, setupRootUser } from '../api/setup'

const router = useRouter()
const message = useMessage()
const formRef = ref<FormInst | null>(null)
const loading = ref(false)

const formValue = ref({
  username: '',
  password: '',
  confirmPassword: '',
})

const rules = {
  username: {
    required: true,
    message: '请输入用户名',
    trigger: 'blur',
  },
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 8, message: '密码至少8个字符', trigger: 'blur' },
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
    const res = await setupRootUser({
      username: formValue.value.username,
      password: formValue.value.password,
      confirmPassword: formValue.value.confirmPassword,
    })
    if (res.data.success) {
      message.success('初始化成功，请登录')
      router.push('/login')
    } else {
      message.error(res.data.message || '初始化失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '初始化失败')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const res = await getSetupStatus()
    if (res.data.success && res.data.data?.setupRequired === false) {
      router.replace('/login')
    }
  } catch {
    // ignore: allow user to try setup if status endpoint fails
  }
})
</script>

<template>
  <div class="setup-container">
    <NCard title="初始化 Root 管理员" class="setup-card">
      <NForm ref="formRef" :model="formValue" :rules="rules">
        <NFormItem label="用户名" path="username">
          <NInput v-model:value="formValue.username" placeholder="请输入用户名" />
        </NFormItem>
        <NFormItem label="密码" path="password">
          <NInput
            v-model:value="formValue.password"
            type="password"
            show-password-on="click"
            placeholder="请输入密码（至少8位）"
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
            初始化
          </NButton>
        </NSpace>
      </NForm>
    </NCard>
  </div>
</template>

<style scoped>
.setup-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.setup-card {
  width: 420px;
}
</style>

