// Форматирование даты под локаль сайта. tm → используем ru-локаль Intl (туркменской в Intl нет).
const INTL_LOCALE: Record<string, string> = {
  tm: 'ru-RU',
  ru: 'ru-RU',
  en: 'en-US'
}

export function formatDate(
  value: string | Date | null | undefined,
  locale = 'ru'
): string {
  if (!value) return ''
  const d = typeof value === 'string' ? new Date(value) : value
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleDateString(INTL_LOCALE[locale] || 'ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}
