import { computed, ref } from 'vue'
import type { User } from './types'

const TOKEN_KEY = 'it_token'
const USER_KEY = 'it_user'

const user = ref<User | null>(readUser())

function readUser(): User | null {
  const raw = localStorage.getItem(USER_KEY)
  if (!raw) {
    return null
  }
  try {
    return JSON.parse(raw) as User
  } catch {
    return null
  }
}

export function useAuth() {
  const isLoggedIn = computed(() => Boolean(user.value && localStorage.getItem(TOKEN_KEY)))
  const isAdmin = computed(() => user.value?.role === 'admin')

  function setSession(token: string, nextUser: User) {
    localStorage.setItem(TOKEN_KEY, token)
    localStorage.setItem(USER_KEY, JSON.stringify(nextUser))
    user.value = nextUser
  }

  function clearSession() {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
    user.value = null
  }

  return { user, isLoggedIn, isAdmin, setSession, clearSession }
}
