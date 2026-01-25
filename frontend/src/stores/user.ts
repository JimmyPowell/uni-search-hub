import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '../types'
import { UserRole } from '../types'
import { getSelf, logout as apiLogout } from '../api/user'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const loading = ref(false)

  const isLoggedIn = computed(() => !!user.value)
  const isAdmin = computed(() => user.value && user.value.role >= UserRole.Admin)
  const isRoot = computed(() => user.value && user.value.role >= UserRole.Root)
  const username = computed(() => user.value?.username || '')
  const displayName = computed(() => user.value?.display_name || user.value?.username || '')

  async function fetchUser() {
    loading.value = true
    try {
      const res = await getSelf()
      if (res.data.success && res.data.data) {
        user.value = res.data.data
      }
    } catch (error) {
      user.value = null
    } finally {
      loading.value = false
    }
  }

  function setUser(u: User) {
    user.value = u
  }

  async function logout() {
    try {
      await apiLogout()
    } finally {
      user.value = null
    }
  }

  function clearUser() {
    user.value = null
  }

  return {
    user,
    loading,
    isLoggedIn,
    isAdmin,
    isRoot,
    username,
    displayName,
    fetchUser,
    setUser,
    logout,
    clearUser,
  }
})
