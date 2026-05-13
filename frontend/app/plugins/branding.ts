// Загружает настройки сайта и применяет CSS-токены (цвета палитры) + фавикон.
// Работает универсально (SSR + клиент) — на SSR токены попадают в HTML <head>,
// что предотвращает FOUC при первой загрузке.

export default defineNuxtPlugin(async () => {
  const config = useRuntimeConfig()
  const settings = useSiteSettings()

  const apiBase = import.meta.server
    ? (config.apiBase || config.public.apiBase)
    : config.public.apiBase

  try {
    const data = await $fetch<any>(`${apiBase}/site-settings`)
    settings.value = data

    const css = `:root {
  --primary: ${data.primary_color};
  --secondary: ${data.secondary_color};
  --accent: ${data.accent_color};
  --background: ${data.background_color};
  --foreground: ${data.foreground_color};
  --muted: ${data.muted_color};
  --border: ${data.border_color};
  --ring: ${data.ring_color};
  --font-family: '${data.font_family}', system-ui, sans-serif;
}`

    const head: Record<string, any> = {
      style: [{ children: css, id: 'site-branding' }]
    }

    if (data.favicon_path) {
      head.link = [{
        rel: 'icon',
        href: `${config.public.apiBase}/files/${data.favicon_path}`
      }]
    }

    useHead(head)
  } catch {
    // Graceful fallback на дефолтные CSS-переменные из main.css
  }
})
