<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()

interface Admin {
  id: number
  username: string
  email: string
  role: string
  all_sections: boolean
  created_at: string
}

const api = useApi()
const toast = useToast()
const { ask } = useConfirm()
const { admin: me } = useAuth()

const { data, refresh, pending } = await useAsyncData('admin-admins', () =>
  api.get<Admin[]>('/admin/admins')
)

const admins = computed(() => data.value ?? [])

async function remove(id: number, username: string) {
  const ok = await ask({
    title: t('admin.admins.confirm_delete_title'),
    description: t('admin.admins.confirm_delete_description', { username }),
    confirmText: t('admin.common.delete'),
    variant: 'destructive'
  })
  if (!ok) return
  try {
    await api.delete(`/admin/admins/${id}`)
    await refresh()
    toast.success(t('admin.admins.deleted'))
  } catch (e: any) {
    toast.error(e.data?.error || t('admin.common.error_delete'))
  }
}

const roleLabel = (r: string) => (r === 'superadmin' ? t('admin.admins.role_superadmin') : t('admin.admins.role_content_manager'))

useHead({ title: () => t('admin.admins.title') })
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold font-serif">{{ t('admin.admins.title') }}</h1>
      <NuxtLink to="/admin/admins/new">
        <Button><Plus class="w-4 h-4" /> {{ t('admin.common.add') }}</Button>
      </NuxtLink>
    </div>

    <div v-if="pending" class="text-muted-foreground">{{ t('admin.common.loading') }}</div>

    <Card v-else class="p-0 overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{{ t('admin.login.username') }}</TableHead>
            <TableHead>{{ t('admin.settings.contact_email') }}</TableHead>
            <TableHead>{{ t('admin.admins.role') }}</TableHead>
            <TableHead>{{ t('admin.nav.sections') }}</TableHead>
            <TableHead class="text-right">{{ t('admin.common.actions') }}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="a in admins" :key="a.id">
            <TableCell class="font-mono">
              <div class="flex items-center gap-2">
                <span class="font-medium">{{ a.username }}</span>
                <Badge v-if="me?.id === a.id" variant="success">{{ t('admin.admins.you') }}</Badge>
              </div>
            </TableCell>
            <TableCell class="text-muted-foreground">{{ a.email || '—' }}</TableCell>
            <TableCell>
              <Badge :variant="a.role === 'superadmin' ? 'warning' : 'secondary'">
                {{ roleLabel(a.role) }}
              </Badge>
            </TableCell>
            <TableCell class="text-muted-foreground">
              {{ a.all_sections ? t('admin.admins.all_sections') : t('admin.admins.sections_selective') }}
            </TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <NuxtLink :to="`/admin/admins/${a.id}`">
                  <Button variant="ghost" size="sm">{{ t('admin.common.edit') }}</Button>
                </NuxtLink>
                <Button
                  v-if="me?.id !== a.id"
                  variant="ghost"
                  size="sm"
                  class="text-destructive hover:text-destructive"
                  @click="remove(a.id, a.username)"
                >
                  {{ t('admin.common.delete') }}
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableEmpty v-if="admins.length === 0" :colspan="5">{{ t('admin.common.no_data') }}</TableEmpty>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>
