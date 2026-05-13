// Глобальный middleware: защита маршрутов /admin/*
// Логика: если путь /admin/* и не /admin/login, проверяем наличие JWT в localStorage
// При отсутствии — редирект на /admin/login
//
// Бэкенд независимо валидирует JWT на каждом запросе — это только UX-защита.

export default defineNuxtRouteMiddleware((to) => {
  const isAdminRoute = to.path.includes('/admin')
  const isLoginPage = to.path.endsWith('/admin/login') || to.path.includes('/admin/login')

  if (!isAdminRoute || isLoginPage) {
    return
  }

  if (import.meta.client) {
    const token = localStorage.getItem('auth_token')
    if (!token) {
      const localePath = useLocalePath()
      return navigateTo(localePath('/admin/login'))
    }
  }
})
