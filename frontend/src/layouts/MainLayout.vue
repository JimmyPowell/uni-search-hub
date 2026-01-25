<script setup lang="ts">
import { h, computed } from 'vue'
import { useRouter } from 'vue-router'
import { NLayout, NLayoutSider, NLayoutContent, NMenu, NIcon, NButton, NAvatar, NDropdown, NSpace, NFlex } from 'naive-ui'
import { HomeOutline, KeyOutline, PersonOutline, PeopleOutline, SettingsOutline, DocumentTextOutline, ServerOutline } from '@vicons/ionicons5'
import type { MenuOption } from 'naive-ui'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

function renderIcon(icon: any) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions = computed<MenuOption[]>(() => {
  const options: MenuOption[] = [
    {
      label: '仪表盘',
      key: 'dashboard',
      icon: renderIcon(HomeOutline),
    },
    {
      label: 'Token管理',
      key: 'tokens',
      icon: renderIcon(KeyOutline),
    },
    {
      label: '个人资料',
      key: 'profile',
      icon: renderIcon(PersonOutline),
    },
    {
      label: '请求日志',
      key: 'logs',
      icon: renderIcon(DocumentTextOutline),
    },
  ]

  if (userStore.isAdmin) {
    options.push({
      label: '用户管理',
      key: 'users',
      icon: renderIcon(PeopleOutline),
    })
  }

  if (userStore.isRoot) {
    options.push({
      label: '渠道管理',
      key: 'channels',
      icon: renderIcon(ServerOutline),
    })
    options.push({
      label: '系统设置',
      key: 'settings',
      icon: renderIcon(SettingsOutline),
    })
  }

  return options
})

const activeKey = computed(() => {
  const path = router.currentRoute.value.path
  if (path === '/') return 'dashboard'
  return path.slice(1)
})

function handleMenuUpdate(key: string) {
  if (key === 'dashboard') {
    router.push('/')
  } else {
    router.push(`/${key}`)
  }
}

const dropdownOptions = [
  {
    label: '个人资料',
    key: 'profile',
  },
  {
    type: 'divider',
    key: 'd1',
  },
  {
    label: '退出登录',
    key: 'logout',
  },
]

async function handleDropdown(key: string) {
  if (key === 'profile') {
    router.push('/profile')
  } else if (key === 'logout') {
    await userStore.logout()
    router.push('/login')
  }
}
</script>

<template>
  <NLayout has-sider style="height: 100vh">
    <NLayoutSider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="220"
      show-trigger
      style="height: 100vh"
    >
      <div class="logo">
        <h2>UniSearch</h2>
      </div>
      <NMenu
        :options="menuOptions"
        :value="activeKey"
        @update:value="handleMenuUpdate"
      />
    </NLayoutSider>
    <NLayout>
      <div class="header">
        <NFlex justify="end" align="center" style="height: 100%; padding-right: 24px;">
          <NDropdown :options="dropdownOptions" @select="handleDropdown">
            <NButton quaternary>
              <NSpace align="center">
                <NAvatar round size="small">{{ userStore.displayName?.charAt(0)?.toUpperCase() }}</NAvatar>
                <span>{{ userStore.displayName }}</span>
              </NSpace>
            </NButton>
          </NDropdown>
        </NFlex>
      </div>
      <NLayoutContent content-style="padding: 24px;" :native-scrollbar="false">
        <router-view />
      </NLayoutContent>
    </NLayout>
  </NLayout>
</template>

<style scoped>
.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid var(--n-border-color);
}

.logo h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--n-text-color);
}

.header {
  height: 64px;
  border-bottom: 1px solid var(--n-border-color);
  background: var(--n-color);
}
</style>
