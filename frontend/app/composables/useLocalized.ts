// Выбирает поле объекта по текущей локали: useLocalized()(section, 'name') → name_ru / name_tm / name_en.
// Фолбэк — туркменский (язык по умолчанию платформы).
export const useLocalized = () => {
  const { locale } = useI18n()
  return (obj: Record<string, any> | null | undefined, base: string): string => {
    if (!obj) return ''
    return obj[`${base}_${locale.value}`] || obj[`${base}_tm`] || ''
  }
}
