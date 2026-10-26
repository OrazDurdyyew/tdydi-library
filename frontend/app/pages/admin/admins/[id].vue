<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const route = useRoute()
const api = useApi()
const toast = useToast()
const { admin: me } = useAuth()
const id = Number(route.params.id)

const { data: adminData, refresh } = await useAsyncData(`admin-${id}`, () =>
  api.get<any>(`/admin/admins/${id}`)
)
const { data: sectionsResp } = await useAsyncData('sections-for-admin-edit', () =>
  api.get<{ data: any[] }>('/sections')
)
const sections = computed(() => sectionsResp.value?.data ?? [])

const form = reactive({
  email: adminData.value?.email ?? '',
  role: adminData.value?.role ?? 'content_manager',
  all_sections: adminData.value?.all_sections ?? true,
  password: '',
})

const permSectionIds = ref<number[]>(
  (adminData.value?.section_permissions ?? []).map((p: any) => p.section_id)
)

const saving = ref(false)
const savingPerms = ref(false)
const isSelf = computed(() => me.value?.id === id)

async function submit() {
  if (isSelf.value) {
    toast.error(t('admin.admins.cannot_edit_self'))
    return
  }
  saving.value = true
  try {
    await api.put(`/admin/admins/${id}`, form)
    await refresh()
    toast.success(t('admin.common.saved'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error'))
  } finally {
    saving.value = false
  }
}

async function savePermissions() {
  savingPerms.value = true
  try {
    await api.put(`/admin/admins/${id}/permissions`, {
      all_sections: form.all_sections,
      section_ids: permSectionIds.value,
    })
    await refresh()
    toast.success(t('admin.admins.permissions_saved'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error'))
  } finally {
    savingPerms.value = false
  }
}

useHead({ title: () => `${t('admin.admins.edit_title')}: ${adminData.value?.username}` })
</script>

<template>
  <div>
    <div class="flex items-center gap-3 mb-6">
      <NuxtLink to="/admin/admins" class="text-muted-foreground hover:text-foreground text-sm">← {{ t('admin.admins.title') }}</NuxtLink>
      <span class="text-muted-foreground">/</span>
      <h1 class="text-2xl font-bold font-serif font-mono">{{ adminData?.username }}</h1>
    </div>

    <Card
      v-if="isSelf"
      class="mb-6 max-w-xl bg-amber-50 dark:bg-amber-950/30 border-amber-200 dark:border-amber-900/50"
    >
      <CardContent class="pt-6 text-sm text-amber-800 dark:text-amber-200">
        {{ t('admin.admins.self_account_warning') }}
      </CardContent>
    </Card>

    <form @submit.prevent="submit" class="space-y-6 max-w-xl mb-10">
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">{{ t('admin.admins.account_data') }}</CardTitle>
        </CardHeader>
        <CardContent class="space-y-5">
          <div class="space-y-2">
            <Label for="email">{{ t('admin.settings.contact_email') }}</Label>
            <Input id="email" v-model="form.email" type="email" :disabled="isSelf" />
          </div>
          <div class="space-y-2">
            <Label for="password">{{ t('admin.admins.password_new') }} ({{ t('admin.admins.password_keep_empty') }})</Label>
            <Input id="password" v-model="form.password" type="password" autocomplete="new-password" :disabled="isSelf" />
          </div>
          <div class="space-y-2">
            <Label for="role">{{ t('admin.admins.role') }}</Label>
            <Select id="role" v-model="form.role" :disabled="isSelf">
              <option value="content_manager">{{ t('admin.admins.role_content_manager') }}</option>
              <option value="publication_manager">{{ t('admin.admins.role_publication_manager') }}</option>
              <option value="superadmin">{{ t('admin.admins.role_superadmin') }}</option>
            </Select>
          </div>
        </CardContent>
      </Card>

      <div class="flex gap-3">
        <Button type="submit" :disabled="saving || isSelf">
          {{ saving ? t('admin.common.saving') : t('admin.common.save') }}
        </Button>
        <NuxtLink to="/admin/admins">
          <Button type="button" variant="outline">{{ t('admin.common.back') }}</Button>
        </NuxtLink>
      </div>
    </form>

    <div v-if="form.role === 'content_manager'" class="max-w-xl">
      <h2 class="text-xl font-bold font-serif mb-4">{{ t('admin.admins.allowed_sections') }}</h2>
      <Card>
        <CardContent class="pt-6 space-y-4">
          <label class="flex items-center gap-3 cursor-pointer">
            <Checkbox v-model="form.all_sections" />
            <span class="text-sm font-medium">{{ t('admin.admins.all_sections') }}</span>
          </label>
          <div v-if="!form.all_sections" class="border border-border rounded-lg p-3 max-h-56 overflow-y-auto space-y-1.5">
            <label
              v-for="s in sections"
              :key="s.id"
              class="flex items-center gap-3 cursor-pointer hover:bg-muted px-2 py-1 rounded-md"
            >
              <Checkbox :value="s.id" v-model="permSectionIds" />
              <span :class="s.parent_id === null ? 'font-medium text-sm' : 'text-sm pl-4 text-muted-foreground'">
                {{ s.name_ru || s.name_tm }}
              </span>
            </label>
          </div>
          <Button @click="savePermissions" :disabled="savingPerms">
            {{ savingPerms ? t('admin.common.saving') : t('admin.common.save') }}
          </Button>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
