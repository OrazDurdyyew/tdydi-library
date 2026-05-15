export const useApi = () => {
  const config = useRuntimeConfig()
  const { token, logout } = useAuth()

  // На SSR используем внутренний URL (api:2020 через docker DNS).
  // На клиенте — публичный URL (localhost:2020 или продакшен-домен).
  const apiBase = import.meta.server
    ? (config.apiBase || config.public.apiBase)
    : config.public.apiBase

  function headers(): Record<string, string> {
    return token.value ? { Authorization: `Bearer ${token.value}` } : {}
  }

  async function request<T>(path: string, opts: any = {}): Promise<T> {
    try {
      return await $fetch<T>(`${apiBase}${path}`, {
        ...opts,
        headers: { ...headers(), ...(opts.headers || {}) }
      })
    } catch (err: any) {
      if (err?.status === 401 || err?.statusCode === 401) {
        logout()
      }
      throw err
    }
  }

  return {
    get: <T>(path: string) => request<T>(path),
    post: <T>(path: string, body: any) => request<T>(path, { method: 'POST', body }),
    put: <T>(path: string, body: any) => request<T>(path, { method: 'PUT', body }),
    delete: <T>(path: string) => request<T>(path, { method: 'DELETE' }),
    upload: <T>(path: string, file: File) => {
      const formData = new FormData()
      formData.append('file', file)
      return request<T>(path, { method: 'POST', body: formData })
    },
    // uploadField — multipart-загрузка с произвольным именем поля (image, cover, file…).
    uploadField: <T>(path: string, field: string, file: File) => {
      const formData = new FormData()
      formData.append(field, file)
      return request<T>(path, { method: 'POST', body: formData })
    }
  }
}
