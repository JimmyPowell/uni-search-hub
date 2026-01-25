<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NCard, NGrid, NGridItem, NStatistic, NSpin, NSpace } from 'naive-ui'
import { useUserStore } from '../stores/user'
import { getAllTokens } from '../api/token'
import type { Token } from '../types'

const userStore = useUserStore()
const loading = ref(true)
const tokens = ref<Token[]>([])

const stats = ref({
  totalTokens: 0,
  activeTokens: 0,
  totalQuota: 0,
  usedQuota: 0,
})

onMounted(async () => {
  try {
    const res = await getAllTokens()
    if (res.data.success && res.data.data) {
      tokens.value = res.data.data.items ?? []
      stats.value.totalTokens = tokens.value.length
      stats.value.activeTokens = tokens.value.filter(t => t.status === 1).length
      stats.value.totalQuota = tokens.value.reduce((sum, t) => sum + (t.unlimited_quota ? 0 : t.remain_quota), 0)
      stats.value.usedQuota = tokens.value.reduce((sum, t) => sum + t.used_quota, 0)
    }
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="dashboard">
    <h2 style="margin-bottom: 24px">欢迎回来，{{ userStore.displayName }}</h2>
    
    <NSpin :show="loading">
      <NGrid :cols="4" :x-gap="16" :y-gap="16">
        <NGridItem>
          <NCard>
            <NStatistic label="Token总数" :value="stats.totalTokens" />
          </NCard>
        </NGridItem>
        <NGridItem>
          <NCard>
            <NStatistic label="活跃Token" :value="stats.activeTokens" />
          </NCard>
        </NGridItem>
        <NGridItem>
          <NCard>
            <NStatistic label="剩余配额" :value="stats.totalQuota" />
          </NCard>
        </NGridItem>
        <NGridItem>
          <NCard>
            <NStatistic label="已使用配额" :value="stats.usedQuota" />
          </NCard>
        </NGridItem>
      </NGrid>

      <NSpace vertical :size="16" style="margin-top: 24px">
        <NCard title="快速开始">
          <p>UniSearch Hub 是一个统一搜索API聚合服务，支持多种搜索引擎（Tavily、Brave、Jina等）。</p>
          <ul>
            <li>在 <strong>Token管理</strong> 中创建API Token来访问搜索服务</li>
            <li>在 <strong>个人资料</strong> 中修改您的账户信息</li>
          </ul>
        </NCard>
      </NSpace>
    </NSpin>
  </div>
</template>

<style scoped>
.dashboard {
  max-width: 1200px;
}
</style>
