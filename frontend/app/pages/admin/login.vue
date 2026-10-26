<script setup lang="ts">
import { Button } from '~/components/ui/button'
import { Input } from '~/components/ui/input'
import { Label } from '~/components/ui/label'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '~/components/ui/card'

definePageMeta({ layout: false })

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const { login } = useAuth()
const localePath = useLocalePath()
const { t } = useI18n()

async function submit() {
  loading.value = true
  error.value = ''
  try {
    await login(username.value, password.value)
    await navigateTo(localePath('/admin'))
  } catch (e: any) {
    error.value = e.data?.error || t('admin.login.error_invalid')
  } finally {
    loading.value = false
  }
}

useHead({ title: () => t('admin.login.title') })
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-background p-4">
    <Card class="w-full max-w-sm">
      <CardHeader>
        <CardTitle>{{ t('admin.login.title') }}</CardTitle>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="submit" class="space-y-5">
          <div class="space-y-2">
            <Label for="username">{{ t('admin.login.username') }}</Label>
            <Input
              id="username"
              v-model="username"
              type="text"
              required
              autocomplete="username"
              placeholder="admin"
            />
          </div>
          <div class="space-y-2">
            <Label for="password">{{ t('admin.login.password') }}</Label>
            <Input
              id="password"
              v-model="password"
              type="password"
              required
              autocomplete="current-password"
              placeholder="••••••••"
            />
          </div>
          <div
            v-if="error"
            class="text-sm text-destructive px-3 py-2 bg-destructive/10 rounded-lg"
          >
            {{ error }}
          </div>
          <Button type="submit" :disabled="loading" class="w-full">
            {{ loading ? t('admin.login.submitting') : t('admin.login.submit') }}
          </Button>
        </form>
      </CardContent>
    </Card>
  </div>
</template>
