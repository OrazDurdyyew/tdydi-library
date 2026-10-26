import { watch } from 'vue'

export default defineNuxtPlugin((nuxtApp) => {
  const { loadFont } = useFont()
  const settings = useSiteSettings()

  // Загружаем шрифт после монтирования приложения
  nuxtApp.hook('app:mounted', () => {
    loadFont()
  })

  // Перезагружаем шрифт при изменении настроек
  watch(settings, () => {
    loadFont()
  }, { deep: true })
})
