<script setup lang="ts">
import { MapPin, Phone, Mail, ExternalLink, BookMarked } from 'lucide-vue-next'

const settings = useSiteSettings()
const { t } = useI18n()
const loc = useLocalized()
const fileUrl = useFileUrl()
const colorMode = useColorMode()
const { data: menu } = await useMenu()

const siteName = computed(() => loc(settings.value, 'name') || t('site.name'))
const footerText = computed(() => loc(settings.value, 'footer'))

const logoUrl = computed(() => {
  if (!settings.value) return ''
  const path =
    colorMode.value === 'dark' && settings.value.logo_dark_path
      ? settings.value.logo_dark_path
      : settings.value.logo_light_path
  return fileUrl(path)
})

const socialLinks = computed(() => {
  const s = settings.value?.social_links
  if (!s || typeof s !== 'object') return []
  return Object.entries(s)
    .filter(([, v]) => !!v)
    .map(([key, value]) => ({ key, value: String(value) }))
})

const year = new Date().getFullYear()
</script>

<template>
  <footer class="mt-20 border-t border-border bg-muted/40">
    <div class="page-container py-12">
      <div class="grid grid-cols-1 md:grid-cols-12 gap-10">
        <!-- О библиотеке -->
        <div class="md:col-span-5">
          <div class="flex items-center gap-2.5 mb-3">
            <img
              v-if="logoUrl"
              :src="logoUrl"
              :alt="siteName"
              class="h-9 w-auto object-contain"
            />
            <span
              v-else
              class="inline-flex items-center justify-center w-9 h-9 rounded-lg bg-primary text-primary-foreground"
            >
              <BookMarked class="w-4 h-4" />
            </span>
            <span class="font-semibold text-[0.95rem]">{{ siteName }}</span>
          </div>
          <div
            v-if="footerText"
            class="text-sm text-muted-foreground leading-relaxed max-w-sm"
            v-html="footerText"
          />
          <p v-else class="text-sm text-muted-foreground max-w-sm">
            {{ t('home.subtitle') }}
          </p>

          <div v-if="socialLinks.length" class="flex flex-wrap gap-2 mt-5">
            <a
              v-for="link in socialLinks"
              :key="link.key"
              :href="link.value"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg border border-border bg-card text-xs font-medium capitalize text-muted-foreground hover:text-secondary hover:border-secondary/40 transition-colors"
            >
              {{ link.key }}
              <ExternalLink class="w-3 h-3" />
            </a>
          </div>
        </div>

        <!-- Навигация -->
        <nav class="md:col-span-4">
          <h3 class="footer-heading">{{ t('footer.navigation') }}</h3>
          <ul class="mt-4 space-y-2.5 text-sm">
            <li><NuxtLinkLocale to="/" class="footer-link">{{ t('nav.home') }}</NuxtLinkLocale></li>
            <li><NuxtLinkLocale to="/books" class="footer-link">{{ t('nav.books') }}</NuxtLinkLocale></li>
            <li>
              <NuxtLinkLocale to="/periodicals" class="footer-link">
                {{ t('nav.periodicals') }}
              </NuxtLinkLocale>
            </li>
            <li><NuxtLinkLocale to="/news" class="footer-link">{{ t('nav.news') }}</NuxtLinkLocale></li>
            <li>
              <NuxtLinkLocale to="/collections" class="footer-link">
                {{ t('nav.collections') }}
              </NuxtLinkLocale>
            </li>
            <li><NuxtLinkLocale to="/faq" class="footer-link">{{ t('nav.faq') }}</NuxtLinkLocale></li>
            <li v-for="page in menu?.pages" :key="page.slug">
              <NuxtLinkLocale :to="`/${page.slug}`" class="footer-link">
                {{ loc(page, 'title') }}
              </NuxtLinkLocale>
            </li>
          </ul>
        </nav>

        <!-- Контакты -->
        <div class="md:col-span-3">
          <h3 class="footer-heading">{{ t('footer.contacts') }}</h3>
          <div class="mt-4 space-y-3 text-sm">
            <div
              v-if="settings?.contact_address"
              class="flex items-start gap-2 text-muted-foreground"
            >
              <MapPin class="w-4 h-4 mt-0.5 icon-primary shrink-0" />
              <span>{{ settings.contact_address }}</span>
            </div>
            <div v-if="settings?.contact_phone" class="flex items-center gap-2">
              <Phone class="w-4 h-4 icon-primary shrink-0" />
              <a :href="`tel:${settings.contact_phone}`" class="footer-link">
                {{ settings.contact_phone }}
              </a>
            </div>
            <div v-if="settings?.contact_email" class="flex items-center gap-2">
              <Mail class="w-4 h-4 icon-primary shrink-0" />
              <a :href="`mailto:${settings.contact_email}`" class="footer-link">
                {{ settings.contact_email }}
              </a>
            </div>
          </div>
        </div>
      </div>

      <div class="divider-double mt-10 mb-5" />
      <p class="text-xs text-muted-foreground">
        © {{ year }} {{ siteName }}. {{ t('footer.rights') }}.
      </p>
    </div>
  </footer>
</template>

<style scoped>
.footer-heading {
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-size: 0.72rem;
  font-weight: 600;
  color: rgb(var(--muted-foreground));
}
.footer-link {
  color: rgb(var(--muted-foreground));
  transition: color 140ms ease;
}
.footer-link:hover {
  color: rgb(var(--secondary));
}
</style>
