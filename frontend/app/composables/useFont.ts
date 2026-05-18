import type { SiteSettings } from './useSiteSettings'

export function useFont() {
  // useState — Nuxt auto-import, composable içinde işleýär
  const settings = useState<SiteSettings | null>('site_settings')

  const loadFont = () => {
    if (!settings.value) return

    const fontFamily = settings.value.font_family || 'Inter'
    const fontUrl = settings.value.font_url
    const fontFilePath = settings.value.font_file_path

    // Удаляем старый элемент шрифта, если он существует
    const oldLink = document.getElementById('custom-font-link')
    if (oldLink) {
      oldLink.remove()
    }

    // Если есть URL шрифта (Google Fonts или другой внешний источник)
    if (fontUrl) {
      const link = document.createElement('link')
      link.id = 'custom-font-link'
      link.rel = 'stylesheet'
      link.href = fontUrl
      document.head.appendChild(link)
    }
    // Если есть загруженный файл шрифта
    else if (fontFilePath) {
      const safeFontName = fontFamily.replace(/\s+/g, '-')
      const fileUrl = `/api/files/${fontFilePath}`

      const style = document.createElement('style')
      style.id = 'custom-font-link'
      style.textContent = `
        @font-face {
          font-family: '${safeFontName}';
          src: url('${fileUrl}') format('woff2'),
               url('${fileUrl}') format('woff'),
               url('${fileUrl}') format('truetype');
          font-weight: normal;
          font-style: normal;
          font-display: swap;
        }
      `
      document.head.appendChild(style)
    }

    // Применяем шрифт к документу
    if (fontFamily) {
      document.documentElement.style.setProperty('--font-family', fontFamily)
      document.body.style.fontFamily = `'${fontFamily}', system-ui, -apple-system, sans-serif`
    }
  }

  const getFontFamily = (): string => {
    if (!settings.value) return 'Inter'
    return settings.value.font_family || 'Inter'
  }

  const hasCustomFont = (): boolean => {
    if (!settings.value) return false
    return !!(settings.value.font_url || settings.value.font_file_path)
  }

  return {
    loadFont,
    getFontFamily,
    hasCustomFont,
  }
}
