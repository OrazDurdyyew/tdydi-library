export type ToastType = 'success' | 'error' | 'info'

export interface Toast {
  id: number
  type: ToastType
  message: string
}

const toasts = ref<Toast[]>([])
let nextId = 0

export const useToast = () => {
  function show(message: string, type: ToastType = 'info', duration = 3500) {
    const id = ++nextId
    toasts.value.push({ id, type, message })
    setTimeout(() => remove(id), duration)
  }

  function remove(id: number) {
    const i = toasts.value.findIndex((t) => t.id === id)
    if (i !== -1) toasts.value.splice(i, 1)
  }

  return {
    toasts: readonly(toasts),
    success: (msg: string) => show(msg, 'success'),
    error: (msg: string) => show(msg, 'error', 5000),
    info: (msg: string) => show(msg, 'info'),
    remove
  }
}
