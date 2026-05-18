<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const api = useApi()
const toast = useToast()
const globalSettings = useSiteSettings()

const { data: settings, refresh } = await useAsyncData('font-settings', () =>
  api.get<any>('/site-settings')
)

const showUpload = ref(false)
const activeTab = ref<'file' | 'url'>('file')
const fontFile = ref<File | null>(null)
const fontUrlInput = ref('')
const fontName = ref('')
const saving = ref(false)

const fontFamily = computed(() => settings.value?.font_family || 'Inter')
const hasCustomFont = computed(() => !!(settings.value?.font_url || settings.value?.font_file_path))

const currentFontInfo = computed(() => {
  if (settings.value?.font_url) return `Внешний шрифт: ${fontFamily.value}`
  if (settings.value?.font_file_path) return `Загруженный файл: ${fontFamily.value}`
  return `Стандартный шрифт: ${fontFamily.value}`
})

const canSave = computed(() => {
  if (activeTab.value === 'file') return !!(fontFile.value && fontName.value.trim())
  return !!(fontUrlInput.value.trim() && fontName.value.trim())
})

const onFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  fontFile.value = target.files?.[0] ?? null
}

const formatDate = (dateString: string) =>
  new Date(dateString).toLocaleDateString('ru-RU', {
    year: 'numeric', month: 'long', day: 'numeric',
    hour: '2-digit', minute: '2-digit'
  })

const resetForm = () => {
  fontFile.value = null
  fontUrlInput.value = ''
  fontName.value = ''
}

async function saveFont() {
  saving.value = true
  try {
    if (activeTab.value === 'file' && fontFile.value) {
      // 1. Шрифт файлы ýükle
      await api.upload('/admin/site-settings/font', fontFile.value)
      // 2. Шрифт adyny sakla
      await api.put('/admin/site-settings', {
        ...settings.value,
        font_family: fontName.value.trim(),
      })
      toast.success('Файл шрифта загружен и применён')
    } else if (activeTab.value === 'url') {
      await api.put('/admin/site-settings', {
        ...settings.value,
        font_family: fontName.value.trim(),
        font_url: fontUrlInput.value.trim(),
      })
      toast.success('Ссылка на шрифт сохранена')
    }

    await refresh()
    globalSettings.value = settings.value
    showUpload.value = false
    resetForm()
  } catch (e: any) {
    toast.error(e.data?.error || 'Не удалось сохранить шрифт')
  } finally {
    saving.value = false
  }
}

async function deleteFont() {
  saving.value = true
  try {
    await api.delete('/admin/site-settings/font')
    await api.put('/admin/site-settings', {
      ...settings.value,
      font_family: 'Inter',
      font_url: '',
    })
    await refresh()
    globalSettings.value = settings.value
    toast.success('Пользовательский шрифт удалён')
  } catch (e: any) {
    toast.error(e.data?.error || 'Не удалось удалить шрифт')
  } finally {
    saving.value = false
  }
}

useHead({ title: 'Шрифт сайта' })
</script>

<template>
  <div>
    <div class="flex items-center gap-4 mb-8">
      <NuxtLink to="/admin/settings">
        <Button variant="ghost" size="sm">← Настройки</Button>
      </NuxtLink>
      <h1 class="text-3xl font-bold font-serif">Шрифт сайта</h1>
    </div>

    <div class="grid gap-6 max-w-2xl">
      <!-- Текущий шрифт -->
      <Card>
        <CardHeader>
          <CardTitle>Текущий шрифт</CardTitle>
          <CardDescription>{{ currentFontInfo }}</CardDescription>
        </CardHeader>
        <CardContent class="space-y-3">
          <div class="flex items-center gap-3">
            <span class="text-2xl" :style="{ fontFamily: fontFamily }">
              Aa Bb Cc — Пример текста
            </span>
            <Badge variant="outline">
              {{ hasCustomFont ? 'Пользовательский' : 'Стандартный' }}
            </Badge>
          </div>
          <div class="text-sm text-muted-foreground space-y-1">
            <div>Название: <strong>{{ fontFamily }}</strong></div>
            <div v-if="settings?.font_url">URL: {{ settings.font_url }}</div>
            <div v-if="settings?.font_file_path">Файл: {{ settings.font_file_path }}</div>
            <div v-if="settings?.font_uploaded_at">
              Загружен: {{ formatDate(settings.font_uploaded_at) }}
            </div>
          </div>
        </CardContent>
        <CardFooter class="flex justify-between">
          <Button
            variant="outline"
            :disabled="!hasCustomFont || saving"
            @click="deleteFont"
          >
            Удалить шрифт
          </Button>
          <Button @click="showUpload = !showUpload">
            {{ showUpload ? 'Отмена' : 'Изменить шрифт' }}
          </Button>
        </CardFooter>
      </Card>

      <!-- Загрузка шрифта -->
      <Card v-if="showUpload">
        <CardHeader>
          <CardTitle>Загрузить шрифт</CardTitle>
          <CardDescription>
            Файл (TTF, OTF, WOFF, WOFF2) или ссылка на Google Fonts
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Tabs v-model="activeTab" class="w-full">
            <TabsList class="grid w-full grid-cols-2 mb-4">
              <TabsTrigger value="file">Файл</TabsTrigger>
              <TabsTrigger value="url">URL</TabsTrigger>
            </TabsList>

            <TabsContent value="file" class="space-y-4">
              <div class="space-y-2">
                <Label for="font-file">Файл шрифта</Label>
                <Input
                  id="font-file"
                  type="file"
                  accept=".ttf,.otf,.woff,.woff2,.eot"
                  @change="onFileChange"
                />
                <p class="text-xs text-muted-foreground">TTF, OTF, WOFF, WOFF2, EOT</p>
              </div>
              <div class="space-y-2">
                <Label for="font-name-file">Название шрифта (font-family)</Label>
                <Input
                  id="font-name-file"
                  v-model="fontName"
                  placeholder="Например: MyFont"
                />
              </div>
            </TabsContent>

            <TabsContent value="url" class="space-y-4">
              <div class="space-y-2">
                <Label for="font-url-input">URL шрифта</Label>
                <Input
                  id="font-url-input"
                  v-model="fontUrlInput"
                  placeholder="https://fonts.googleapis.com/css2?family=Roboto&display=swap"
                />
                <p class="text-xs text-muted-foreground">Google Fonts или другой внешний источник</p>
              </div>
              <div class="space-y-2">
                <Label for="font-name-url">Название шрифта (font-family)</Label>
                <Input
                  id="font-name-url"
                  v-model="fontName"
                  placeholder="Например: Roboto"
                />
              </div>
            </TabsContent>
          </Tabs>
        </CardContent>
        <CardFooter class="flex justify-end">
          <Button :disabled="!canSave || saving" @click="saveFont">
            {{ saving ? 'Сохранение…' : 'Сохранить' }}
          </Button>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>
