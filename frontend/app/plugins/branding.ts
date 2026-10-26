// Загружает настройки сайта и применяет CSS-токены (палитра), шрифт и фавикон.
// ВАЖНО: useHead вызывается после await — контекст Nuxt теряется, поэтому
// оборачиваем в nuxtApp.runWithContext(), иначе <style> не попадает в <head>.

export default defineNuxtPlugin(async (nuxtApp) => {
  const config = useRuntimeConfig()
  const settings = useSiteSettings()

  const apiBase = import.meta.server
    ? config.apiBase || config.public.apiBase
    : config.public.apiBase

  let data: any
  try {
    data = await $fetch<any>(`${apiBase}/site-settings`)
  } catch {
    return // graceful fallback на дефолтные CSS-переменные из main.css
  }

  settings.value = data

  const link: Record<string, any>[] = [
    { rel: 'stylesheet', href: fontHref(data), id: 'site-font' }
  ]
  if (data.favicon_path) {
    link.push({ rel: 'icon', href: `${config.public.apiBase}/files/${data.favicon_path}` })
  }

  nuxtApp.runWithContext(() => {
    useHead({
      style: [{ innerHTML: brandingCss(data), id: 'site-branding', tagPriority: 'critical' } as any],
      link
    })
  })
})
