// Применение брендинга (палитра + шрифт + фавикон) из site_settings.
// Используется и в plugins/branding.ts (первичная загрузка), и в админке настроек
// (мгновенное применение после сохранения, без перезагрузки страницы).

export function brandingCss(s: Record<string, any>): string {
  const font = (s.font_family || 'Inter').trim()
  const fontVars = `  --font-sans: '${font}', system-ui, -apple-system, sans-serif;
  --font-family: var(--font-sans);`

  // Светлая и тёмная палитры — отдельные наборы полей в site_settings.
  // :root:not(.dark) и :root.dark имеют специфичность (0,2,0) — перебивают
  // :root и .dark из main.css независимо от порядка подключения стилей.
  return `:root:not(.dark) {
  --primary: ${s.primary_color};
  --secondary: ${s.secondary_color};
  --accent: ${s.accent_color};
  --ring: ${s.ring_color};
  --background: ${s.background_color};
  --foreground: ${s.foreground_color};
  --muted: ${s.muted_color};
  --border: ${s.border_color};
${fontVars}
}
:root.dark {
  --primary: ${s.dark_primary_color};
  --secondary: ${s.dark_secondary_color};
  --accent: ${s.dark_accent_color};
  --ring: ${s.dark_ring_color};
  --background: ${s.dark_background_color};
  --foreground: ${s.dark_foreground_color};
  --muted: ${s.dark_muted_color};
  --border: ${s.dark_border_color};
${fontVars}
}`
}

export function fontHref(s: Record<string, any>): string {
  const font = (s.font_family || 'Inter').trim()
  return `https://fonts.googleapis.com/css2?family=${font.replace(/ /g, '+')}:wght@400;500;600;700&display=swap`
}

// Применяет брендинг напрямую через DOM — для мгновенного эффекта на клиенте.
export function applyBrandingClient(s: Record<string, any>, apiBase: string): void {
  if (!import.meta.client || !s) return

  let style = document.getElementById('site-branding') as HTMLStyleElement | null
  if (!style) {
    style = document.createElement('style')
    style.id = 'site-branding'
    document.head.appendChild(style)
  }
  style.textContent = brandingCss(s)

  let fontLink = document.getElementById('site-font') as HTMLLinkElement | null
  if (!fontLink) {
    fontLink = document.createElement('link')
    fontLink.id = 'site-font'
    fontLink.rel = 'stylesheet'
    document.head.appendChild(fontLink)
  }
  fontLink.href = fontHref(s)

  if (s.favicon_path) {
    let fav = document.querySelector("link[rel='icon']") as HTMLLinkElement | null
    if (!fav) {
      fav = document.createElement('link')
      fav.rel = 'icon'
      document.head.appendChild(fav)
    }
    fav.href = `${apiBase}/files/${s.favicon_path}`
  }
}
