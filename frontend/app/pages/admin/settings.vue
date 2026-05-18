<script setup lang="ts">
import { Upload, Plus, Trash2 } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const api = useApi()
const config = useRuntimeConfig()
const toast = useToast()
const globalSettings = useSiteSettings()

const { data: settings, refresh } = await useAsyncData('site-settings-admin', () =>
  api.get<any>('/site-settings')
)

const COLOR_FIELDS = [
  { key: 'primary_color', label: 'Primary', fallback: '30 58 95' },
  { key: 'secondary_color', label: 'Secondary', fallback: '44 95 141' },
  { key: 'accent_color', label: 'Accent', fallback: '200 169 107' },
  { key: 'background_color', label: 'Background', fallback: '248 250 252' },
  { key: 'foreground_color', label: 'Foreground', fallback: '17 24 39' },
  { key: 'muted_color', label: 'Muted', fallback: '241 245 249' },
  { key: 'border_color', label: 'Border', fallback: '229 231 235' },
  { key: 'ring_color', label: 'Ring', fallback: '44 95 141' }
] as const

const DARK_COLOR_FIELDS = [
  { key: 'dark_primary_color', label: 'Primary', fallback: '122 162 201' },
  { key: 'dark_secondary_color', label: 'Secondary', fallback: '96 141 184' },
  { key: 'dark_accent_color', label: 'Accent', fallback: '200 169 107' },
  { key: 'dark_background_color', label: 'Background', fallback: '15 23 42' },
  { key: 'dark_foreground_color', label: 'Foreground', fallback: '249 250 251' },
  { key: 'dark_muted_color', label: 'Muted', fallback: '30 41 59' },
  { key: 'dark_border_color', label: 'Border', fallback: '31 41 55' },
  { key: 'dark_ring_color', label: 'Ring', fallback: '122 162 201' }
] as const

const form = reactive<Record<string, string>>({
  name_tm: settings.value?.name_tm ?? '',
  name_ru: settings.value?.name_ru ?? '',
  name_en: settings.value?.name_en ?? '',
  font_family: settings.value?.font_family ?? 'Inter',
  contact_address: settings.value?.contact_address ?? '',
  contact_phone: settings.value?.contact_phone ?? '',
  contact_email: settings.value?.contact_email ?? '',
  footer_tm: settings.value?.footer_tm ?? '',
  footer_ru: settings.value?.footer_ru ?? '',
  footer_en: settings.value?.footer_en ?? ''
})
for (const f of [...COLOR_FIELDS, ...DARK_COLOR_FIELDS]) {
  form[f.key] = settings.value?.[f.key] || f.fallback
}

// Соц. ссылки: объект { telegram: url } ↔ массив строк для редактирования
const socialLinks = reactive<{ key: string; value: string }[]>(
  Object.entries(settings.value?.social_links ?? {}).map(([key, value]) => ({
    key,
    value: String(value)
  }))
)
function addSocial() {
  socialLinks.push({ key: '', value: '' })
}
function removeSocial(i: number) {
  socialLinks.splice(i, 1)
}

// ─── Конвертация "R G B" ↔ "#rrggbb" ───
function rgbToHex(rgb: string | undefined): string {
  const parts = (rgb || '').trim().split(/\s+/).map(Number)
  if (parts.length !== 3 || parts.some((n) => Number.isNaN(n))) return '#000000'
  return (
    '#' +
    parts
      .map((n) => Math.max(0, Math.min(255, n)).toString(16).padStart(2, '0'))
      .join('')
  )
}
function hexToRgb(hex: string): string {
  const m = hex.replace('#', '')
  const r = parseInt(m.slice(0, 2), 16)
  const g = parseInt(m.slice(2, 4), 16)
  const b = parseInt(m.slice(4, 6), 16)
  return `${r} ${g} ${b}`
}

const previewStyle = computed(() => ({
  '--p-bg': `rgb(${form.background_color})`,
  '--p-fg': `rgb(${form.foreground_color})`,
  '--p-primary': `rgb(${form.primary_color})`,
  '--p-muted': `rgb(${form.muted_color})`,
  '--p-border': `rgb(${form.border_color})`
}))

const darkPreviewStyle = computed(() => ({
  '--p-bg': `rgb(${form.dark_background_color})`,
  '--p-fg': `rgb(${form.dark_foreground_color})`,
  '--p-primary': `rgb(${form.dark_primary_color})`,
  '--p-muted': `rgb(${form.dark_muted_color})`,
  '--p-border': `rgb(${form.dark_border_color})`
}))

// ─── Загрузка логотипа / фавикона ───
const fileUrl = useFileUrl()
const uploading = ref<string | null>(null)
const assetPaths = reactive<Record<string, string>>({
  logo_light: settings.value?.logo_light_path ?? '',
  logo_dark: settings.value?.logo_dark_path ?? '',
  favicon: settings.value?.favicon_path ?? ''
})

async function uploadAsset(kind: 'logo_light' | 'logo_dark' | 'favicon', e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  uploading.value = kind
  try {
    const res = await api.upload<{ path: string }>(
      `/admin/site-settings/asset?kind=${kind}`,
      file
    )
    assetPaths[kind] = res.path
    toast.success('Файл загружен')
  } catch (err: any) {
    toast.error(err.data?.error || 'Ошибка загрузки')
  } finally {
    uploading.value = null
    input.value = ''
  }
}

// ─── Сохранение ───
const saving = ref(false)
async function submit() {
  saving.value = true
  try {
    const social: Record<string, string> = {}
    for (const s of socialLinks) {
      if (s.key.trim() && s.value.trim()) social[s.key.trim()] = s.value.trim()
    }
    await api.put<any>('/admin/site-settings', { ...form, social_links: social })
    await refresh()
    globalSettings.value = settings.value
    // Применяем палитру/шрифт/фавикон сразу, без перезагрузки страницы.
    applyBrandingClient(settings.value, config.public.apiBase)
    toast.success('Настройки сохранены')
  } catch (e: any) {
    toast.error(e.data?.error || 'Ошибка сохранения')
  } finally {
    saving.value = false
  }
}

useHead({ title: 'Настройки сайта' })
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold font-serif mb-8">Настройки сайта</h1>

    <form @submit.prevent="submit" class="space-y-6 max-w-3xl">
      <!-- Название -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Название института</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="name-tm">Туркменский</Label>
            <Input id="name-tm" v-model="form.name_tm" />
          </div>
          <div class="space-y-2">
            <Label for="name-ru">Русский</Label>
            <Input id="name-ru" v-model="form.name_ru" />
          </div>
          <div class="space-y-2">
            <Label for="name-en">English</Label>
            <Input id="name-en" v-model="form.name_en" />
          </div>
        </CardContent>
      </Card>

      <!-- Логотип и фавикон -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Логотип и фавикон</CardTitle>
          <CardDescription>PNG / JPG / WebP. Фавикон лучше квадратный (например 64×64).</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-5">
            <div
              v-for="asset in [
                { kind: 'logo_light', label: 'Логотип (светлый)' },
                { kind: 'logo_dark', label: 'Логотип (тёмный)' },
                { kind: 'favicon', label: 'Фавикон' }
              ]"
              :key="asset.kind"
              class="space-y-2"
            >
              <Label>{{ asset.label }}</Label>
              <div
                class="aspect-square rounded-md border border-border bg-muted flex items-center justify-center overflow-hidden"
                :class="asset.kind === 'logo_dark' ? 'bg-secondary' : ''"
              >
                <img
                  v-if="assetPaths[asset.kind]"
                  :src="fileUrl(assetPaths[asset.kind])"
                  :alt="asset.label"
                  class="max-w-[80%] max-h-[80%] object-contain"
                />
                <Upload v-else class="w-7 h-7 text-muted-foreground" />
              </div>
              <label class="btn-outline w-full cursor-pointer text-sm">
                <Upload class="w-4 h-4" />
                {{ uploading === asset.kind ? 'Загрузка…' : 'Выбрать файл' }}
                <input
                  type="file"
                  accept="image/*"
                  class="hidden"
                  :disabled="uploading === asset.kind"
                  @change="uploadAsset(asset.kind as any, $event)"
                />
              </label>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Цветовые токены — светлая тема -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Палитра — светлая тема</CardTitle>
          <CardDescription>Цвета применяются на сайте без передеплоя.</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
            <div v-for="f in COLOR_FIELDS" :key="f.key" class="space-y-1.5">
              <Label :for="`c-${f.key}`">{{ f.label }}</Label>
              <div class="flex items-center gap-2">
                <input
                  :id="`c-${f.key}`"
                  type="color"
                  class="h-9 w-12 rounded border border-border cursor-pointer bg-transparent p-0.5"
                  :value="rgbToHex(form[f.key])"
                  @input="form[f.key] = hexToRgb(($event.target as HTMLInputElement).value)"
                />
                <input
                  v-model="form[f.key]"
                  class="input font-mono text-xs flex-1"
                  placeholder="R G B"
                />
              </div>
            </div>
          </div>

          <!-- Live preview -->
          <div
            class="mt-6 rounded-lg border p-5"
            :style="{ ...previewStyle, backgroundColor: 'var(--p-bg)', borderColor: 'var(--p-border)' }"
          >
            <div class="text-xs uppercase tracking-widest mb-2" :style="{ color: 'var(--p-primary)' }">
              Предпросмотр
            </div>
            <div class="font-serif font-bold text-xl" :style="{ color: 'var(--p-fg)' }">
              {{ form.name_ru || 'Электронная библиотека' }}
            </div>
            <div class="flex gap-2 mt-3">
              <span
                class="px-3 py-1.5 rounded-md text-sm font-medium"
                :style="{ backgroundColor: 'var(--p-primary)', color: 'var(--p-bg)' }"
              >
                Кнопка
              </span>
              <span
                class="px-3 py-1.5 rounded-md text-sm"
                :style="{ backgroundColor: 'var(--p-muted)', color: 'var(--p-fg)' }"
              >
                Вторичная
              </span>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Цветовые токены — тёмная тема -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Палитра — тёмная тема</CardTitle>
          <CardDescription>Применяется, когда посетитель выбрал тёмную тему.</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
            <div v-for="f in DARK_COLOR_FIELDS" :key="f.key" class="space-y-1.5">
              <Label :for="`c-${f.key}`">{{ f.label }}</Label>
              <div class="flex items-center gap-2">
                <input
                  :id="`c-${f.key}`"
                  type="color"
                  class="h-9 w-12 rounded border border-border cursor-pointer bg-transparent p-0.5"
                  :value="rgbToHex(form[f.key])"
                  @input="form[f.key] = hexToRgb(($event.target as HTMLInputElement).value)"
                />
                <input
                  v-model="form[f.key]"
                  class="input font-mono text-xs flex-1"
                  placeholder="R G B"
                />
              </div>
            </div>
          </div>

          <!-- Live preview (тёмная) -->
          <div
            class="mt-6 rounded-lg border p-5"
            :style="{ ...darkPreviewStyle, backgroundColor: 'var(--p-bg)', borderColor: 'var(--p-border)' }"
          >
            <div class="text-xs uppercase tracking-widest mb-2" :style="{ color: 'var(--p-primary)' }">
              Предпросмотр
            </div>
            <div class="font-serif font-bold text-xl" :style="{ color: 'var(--p-fg)' }">
              {{ form.name_ru || 'Электронная библиотека' }}
            </div>
            <div class="flex gap-2 mt-3">
              <span
                class="px-3 py-1.5 rounded-md text-sm font-medium"
                :style="{ backgroundColor: 'var(--p-primary)', color: 'var(--p-bg)' }"
              >
                Кнопка
              </span>
              <span
                class="px-3 py-1.5 rounded-md text-sm"
                :style="{ backgroundColor: 'var(--p-muted)', color: 'var(--p-fg)' }"
              >
                Вторичная
              </span>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Типографика -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Шрифт</CardTitle>
          <CardDescription>
            Выберите стандартный шрифт или
            <NuxtLink to="/admin/settings/font" class="underline text-primary">загрузите собственный</NuxtLink>.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-2 max-w-xs">
            <Label for="font">Семейство шрифтов</Label>
            <Select id="font" v-model="form.font_family">
              <option value="Inter">Inter</option>
              <option value="Roboto">Roboto</option>
              <option value="Open Sans">Open Sans</option>
              <option value="Noto Sans">Noto Sans</option>
              <option value="Montserrat">Montserrat</option>
              <option value="PT Sans">PT Sans</option>
            </Select>
          </div>
          <div v-if="settings?.font_file_path || settings?.font_url" class="mt-3 text-sm text-muted-foreground flex items-center gap-2">
            <span class="inline-block w-2 h-2 rounded-full bg-green-500"></span>
            Пользовательский шрифт загружен: <strong>{{ settings?.font_family }}</strong>
          </div>
        </CardContent>
      </Card>

      <!-- Контакты -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Контакты</CardTitle>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="addr">Адрес</Label>
            <Textarea id="addr" v-model="form.contact_address" rows="2" />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <Label for="phone">Телефон</Label>
              <Input id="phone" v-model="form.contact_phone" />
            </div>
            <div class="space-y-2">
              <Label for="email">Email</Label>
              <Input id="email" v-model="form.contact_email" type="email" />
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Соц. ссылки -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Социальные сети</CardTitle>
          <CardDescription>Название (например <code class="font-mono">telegram</code>) и ссылка.</CardDescription>
        </CardHeader>
        <CardContent class="space-y-3">
          <div v-for="(s, i) in socialLinks" :key="i" class="flex items-center gap-2">
            <Input v-model="s.key" placeholder="telegram" class="max-w-[160px]" />
            <Input v-model="s.value" placeholder="https://t.me/…" class="flex-1" />
            <Button type="button" variant="ghost" size="sm" class="text-destructive" @click="removeSocial(i)">
              <Trash2 class="w-4 h-4" />
            </Button>
          </div>
          <Button type="button" variant="outline" size="sm" @click="addSocial">
            <Plus class="w-4 h-4" /> Добавить ссылку
          </Button>
        </CardContent>
      </Card>

      <!-- Подвал -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Текст подвала</CardTitle>
          <CardDescription>Можно использовать HTML.</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="foot-tm">Туркменский</Label>
            <Textarea id="foot-tm" v-model="form.footer_tm" rows="2" />
          </div>
          <div class="space-y-2">
            <Label for="foot-ru">Русский</Label>
            <Textarea id="foot-ru" v-model="form.footer_ru" rows="2" />
          </div>
          <div class="space-y-2">
            <Label for="foot-en">English</Label>
            <Textarea id="foot-en" v-model="form.footer_en" rows="2" />
          </div>
        </CardContent>
      </Card>

      <Button type="submit" :disabled="saving">
        {{ saving ? 'Сохранение…' : 'Сохранить настройки' }}
      </Button>
    </form>
  </div>
</template>
