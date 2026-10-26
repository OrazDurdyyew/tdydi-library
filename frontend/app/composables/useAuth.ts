export interface AdminInfo {
  id: number
  username: string
  email?: string
  role: 'superadmin' | 'content_manager' | 'publication_manager'
  all_sections: boolean
  allowed_section_ids: number[]
}

interface LoginResponse {
  access_token: string
  refresh_token: string
  admin: AdminInfo
}

const ACCESS_KEY = 'auth_token'
const REFRESH_KEY = 'refresh_token'

export const useAuth = () => {
  const config = useRuntimeConfig()
  const token = useState<string | null>('auth_token', () => null)
  const admin = useState<AdminInfo | null>('auth_admin', () => null)

  // Восстановить токен из localStorage при первом обращении на клиенте
  if (import.meta.client && token.value === null) {
    const stored = localStorage.getItem(ACCESS_KEY)
    if (stored) token.value = stored
  }

  async function login(username: string, password: string) {
    const res = await $fetch<LoginResponse>(`${config.public.apiBase}/auth/login`, {
      method: 'POST',
      body: { username, password }
    })
    token.value = res.access_token
    admin.value = res.admin
    if (import.meta.client) {
      localStorage.setItem(ACCESS_KEY, res.access_token)
      localStorage.setItem(REFRESH_KEY, res.refresh_token)
    }
  }

  async function fetchMe() {
    if (!token.value) return null
    try {
      const data = await $fetch<AdminInfo>(`${config.public.apiBase}/admin/me`, {
        headers: { Authorization: `Bearer ${token.value}` }
      })
      admin.value = data
      return data
    } catch {
      logout()
      return null
    }
  }

  function logout() {
    token.value = null
    admin.value = null
    if (import.meta.client) {
      localStorage.removeItem(ACCESS_KEY)
      localStorage.removeItem(REFRESH_KEY)
    }
  }

  return { token, admin, login, logout, fetchMe }
}
