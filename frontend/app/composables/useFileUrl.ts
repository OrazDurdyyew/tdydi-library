export const useFileUrl = () => {
  const config = useRuntimeConfig()
  // Для файлов всегда используем публичный URL (браузер делает запрос).
  // Server-side рендеринг тоже отдаёт ссылки клиенту, не загружает файлы сам.
  const apiBase = config.public.apiBase

  return (relativePath: string | null | undefined): string => {
    if (!relativePath) return ''
    return `${apiBase}/files/${relativePath}`
  }
}
