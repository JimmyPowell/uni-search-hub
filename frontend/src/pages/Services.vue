<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { NCard, NGrid, NGridItem, NTag, NTable, NCode, NButton, useMessage } from 'naive-ui'
import { getServiceCatalog } from '../api/meta'
import type { ServiceItem, ServiceEndpoint } from '../types'

const message = useMessage()
const loading = ref(false)
const services = ref<ServiceItem[]>([])

const grouped = computed(() => {
  const m: Record<string, ServiceItem[]> = {}
  for (const s of services.value) {
    if (!m[s.provider]) m[s.provider] = []
    m[s.provider]!.push(s)
  }
  return Object.entries(m).sort((a, b) => a[0].localeCompare(b[0]))
})

async function fetchCatalog() {
  loading.value = true
  try {
    const res = await getServiceCatalog()
    if (res.data.success && res.data.data) {
      services.value = res.data.data.services ?? []
    } else {
      message.error(res.data.message || '获取服务目录失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '获取服务目录失败')
  } finally {
    loading.value = false
  }
}

function formatCurl(ep: ServiceEndpoint) {
  const headers = ep.example_headers || {}
  const headerArgs = Object.entries(headers)
    .map(([k, v]) => `-H '${k}: ${v}'`)
    .join(' ')
  const body = ep.example_body ? JSON.stringify(ep.example_body) : '{}'
  return `curl -sS ${headerArgs} -X ${ep.method} 'http://<host>${ep.path}' -d '${body}'`
}

function copy(text: string) {
  navigator.clipboard.writeText(text)
  message.success('已复制')
}

onMounted(fetchCatalog)
</script>

<template>
  <div>
    <NCard title="服务目录" :segmented="{ content: true }">
      <template #header-extra>
        <NTag v-if="loading" type="info">加载中</NTag>
      </template>

      <NGrid :cols="1" :y-gap="12">
        <NGridItem v-for="[provider, items] in grouped" :key="provider">
          <NCard :title="provider" size="small">
            <div style="margin-bottom: 8px;">
              <NTag :type="(items?.[0]?.enabled_channels_count ?? 0) > 0 ? 'success' : 'warning'">
                enabled channels: {{ items?.[0]?.enabled_channels_count ?? 0 }}
              </NTag>
            </div>

            <div v-for="svc in items" :key="svc.id" style="margin-bottom: 12px;">
              <div style="font-weight: 600; margin-bottom: 6px;">{{ svc.name }}</div>
              <NTable size="small" :single-line="false">
                <thead>
                  <tr>
                    <th style="width: 90px;">Method</th>
                    <th>Path</th>
                    <th style="width: 80px;">Auth</th>
                    <th>Notes</th>
                    <th style="width: 90px;">Action</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="ep in svc.endpoints" :key="ep.method + ep.path">
                    <td><code>{{ ep.method }}</code></td>
                    <td><code>{{ ep.path }}</code></td>
                    <td>
                      <NTag size="small" :type="ep.auth === 'token' ? 'info' : 'default'">{{ ep.auth }}</NTag>
                    </td>
                    <td>{{ ep.notes }}</td>
                    <td>
                      <NButton size="tiny" @click="copy(formatCurl(ep))">复制 curl</NButton>
                    </td>
                  </tr>
                </tbody>
              </NTable>

              <div style="margin-top: 8px;">
                <div style="color: #666; margin-bottom: 4px;">curl 示例（可复制）：</div>
                <NCode
                  :code="svc.endpoints?.[0] ? formatCurl(svc.endpoints[0]) : ''"
                  language="bash"
                  word-wrap
                />
              </div>
            </div>
          </NCard>
        </NGridItem>
      </NGrid>
    </NCard>
  </div>
</template>
