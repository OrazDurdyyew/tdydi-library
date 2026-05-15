// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2026-01-01',
  future: { compatibilityVersion: 4 },
  devtools: { enabled: true },
  devServer: {
    host: '0.0.0.0',
    port: 2026
  },

  modules: [
    '@nuxtjs/tailwindcss',
    '@nuxtjs/color-mode',
    '@nuxtjs/i18n',
    '@pinia/nuxt',
    '@vueuse/nuxt',
    'shadcn-nuxt'
  ],

  css: ['~/assets/css/main.css'],

  colorMode: {
    classSuffix: '',
    preference: 'system',
    fallback: 'light'
  },

  shadcn: {
    prefix: '',
    componentDir: './app/components/ui'
  },

  i18n: {
    defaultLocale: 'tm',
    // tm — основной язык, в URL без префикса. ru/en — с префиксом (/ru/..., /en/...).
    strategy: 'prefix_except_default',
    detectBrowserLanguage: false,
    locales: [
      { code: 'tm', file: 'tm.json', name: 'Türkmençe' },
      { code: 'ru', file: 'ru.json', name: 'Русский' },
      { code: 'en', file: 'en.json', name: 'English' }
    ]
  },

  runtimeConfig: {
    apiBase: process.env.API_BASE_INTERNAL || 'http://api:2020/api',
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:2020/api'
    }
  },

  typescript: {
    strict: true
  },

  app: {
    head: {
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Source+Serif+4:ital,opsz,wght@0,8..60,400;0,8..60,500;0,8..60,600;1,8..60,400&display=swap'
        }
      ]
    }
  }
})
