<script setup lang="ts">
import { MapPin, Phone, Mail } from 'lucide-vue-next'

const settings = useSiteSettings()
const { locale } = useI18n()

const siteName = computed(() => {
  if (!settings.value) return ''
  if (locale.value === 'tm') return settings.value.name_tm
  if (locale.value === 'en') return settings.value.name_en
  return settings.value.name_ru
})

const footerText = computed(() => {
  if (!settings.value) return ''
  if (locale.value === 'tm') return settings.value.footer_tm
  if (locale.value === 'en') return settings.value.footer_en
  return settings.value.footer_ru
})
</script>

<template>
  <footer class="mt-20 pb-12">
    <div class="page-container">
      <!-- Двойная орнаментальная линия в духе титульного листа -->
      <div class="divider-double max-w-md mx-auto mb-2" />
      <div class="text-center font-serif italic text-muted-foreground text-sm">❦</div>
      <div class="divider-double max-w-md mx-auto mt-2 mb-10" />

      <div class="flex flex-col md:flex-row md:justify-between gap-10 text-sm">
        <div class="max-w-md">
          <div class="font-serif font-bold text-lg text-foreground">{{ siteName }}</div>
          <div
            v-if="footerText"
            class="text-muted-foreground mt-3 leading-relaxed serif"
            v-html="footerText"
          />
          <div v-else class="text-muted-foreground mt-3 serif italic">
            anno domini {{ new Date().getFullYear() }}
          </div>
        </div>
        <div v-if="settings" class="space-y-2.5">
          <div
            v-if="settings.contact_address"
            class="flex items-start gap-2 text-muted-foreground"
          >
            <MapPin class="w-4 h-4 mt-0.5 icon-primary shrink-0" />
            <span class="serif italic">{{ settings.contact_address }}</span>
          </div>
          <div v-if="settings.contact_phone" class="flex items-center gap-2">
            <Phone class="w-4 h-4 icon-primary shrink-0" />
            <a
              :href="`tel:${settings.contact_phone}`"
              class="text-muted-foreground hover:text-foreground transition-colors"
            >
              {{ settings.contact_phone }}
            </a>
          </div>
          <div v-if="settings.contact_email" class="flex items-center gap-2">
            <Mail class="w-4 h-4 icon-primary shrink-0" />
            <a
              :href="`mailto:${settings.contact_email}`"
              class="text-muted-foreground hover:text-foreground transition-colors"
            >
              {{ settings.contact_email }}
            </a>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>
