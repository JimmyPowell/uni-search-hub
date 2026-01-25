<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NCard, NForm, NFormItem, NInput, NButton, NSpace, NSpin, useMessage } from 'naive-ui'
import { getOptions, updateOption } from '../api/option'
import type { Option } from '../types'

const message = useMessage()
const loading = ref(false)
const saving = ref(false)
const options = ref<Option[]>([])
const editedOptions = ref<Record<string, string>>({})

async function fetchOptions() {
  loading.value = true
  try {
    const res = await getOptions()
    if (res.data.success && res.data.data) {
      options.value = res.data.data
      editedOptions.value = {}
      options.value.forEach(opt => {
        editedOptions.value[opt.key] = opt.value
      })
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '获取配置失败')
  } finally {
    loading.value = false
  }
}

async function handleSave(key: string) {
  saving.value = true
  try {
    const value = editedOptions.value[key] ?? ''
    const res = await updateOption({ key, value })
    if (res.data.success) {
      message.success('保存成功')
      const opt = options.value.find(o => o.key === key)
      if (opt) {
        opt.value = value
      }
    } else {
      message.error(res.data.message || '保存失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

function getOptionLabel(key: string): string {
  const labels: Record<string, string> = {
    'SystemName': '系统名称',
    'SystemDescription': '系统描述',
    'FooterText': '页脚文本',
    'RegisterEnabled': '允许注册',
    'TurnstileEnabled': 'Turnstile验证',
    'TurnstileSiteKey': 'Turnstile Site Key',
    'TurnstileSecretKey': 'Turnstile Secret Key',
  }
  return labels[key] || key
}

function getOptionType(key: string): string {
  if (key.includes('Enabled')) return 'boolean'
  if (key.includes('Key') || key.includes('Secret')) return 'password'
  return 'text'
}

onMounted(fetchOptions)
</script>

<template>
  <div class="settings-page">
    <NCard title="系统设置">
      <NSpin :show="loading">
        <NForm label-placement="left" label-width="160">
          <template v-for="option in options" :key="option.key">
            <NFormItem :label="getOptionLabel(option.key)">
              <NSpace align="center" style="width: 100%">
                <NInput
                  v-if="getOptionType(option.key) === 'password'"
                  :value="editedOptions[option.key] ?? ''"
                  @update:value="(v) => (editedOptions[option.key] = v)"
                  type="password"
                  show-password-on="click"
                  style="width: 400px"
                />
                <NInput
                  v-else
                  :value="editedOptions[option.key] ?? ''"
                  @update:value="(v) => (editedOptions[option.key] = v)"
                  style="width: 400px"
                />
                <NButton
                  type="primary"
                  size="small"
                  :loading="saving"
                  :disabled="(editedOptions[option.key] ?? '') === option.value"
                  @click="handleSave(option.key)"
                >
                  保存
                </NButton>
              </NSpace>
            </NFormItem>
          </template>
        </NForm>
        
        <template v-if="options.length === 0 && !loading">
          <p style="text-align: center; color: #999">暂无配置项</p>
        </template>
      </NSpin>
    </NCard>

    <NCard title="说明" style="margin-top: 16px">
      <ul>
        <li><strong>RegisterEnabled</strong>: 设置为 true 允许用户注册，false 则关闭注册</li>
        <li><strong>TurnstileEnabled</strong>: 设置为 true 开启 Cloudflare Turnstile 人机验证</li>
        <li><strong>TurnstileSiteKey</strong> / <strong>TurnstileSecretKey</strong>: Cloudflare Turnstile 密钥</li>
      </ul>
    </NCard>
  </div>
</template>

<style scoped>
.settings-page {
  max-width: 800px;
}
</style>
