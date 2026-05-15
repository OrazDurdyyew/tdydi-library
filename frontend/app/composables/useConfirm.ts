interface ConfirmOptions {
  title: string
  description?: string
  confirmText?: string
  cancelText?: string
  variant?: 'default' | 'destructive'
}

interface ConfirmState extends Required<Omit<ConfirmOptions, 'description'>> {
  open: boolean
  description: string
}

const state = reactive<ConfirmState>({
  open: false,
  title: '',
  description: '',
  confirmText: 'Подтвердить',
  cancelText: 'Отмена',
  variant: 'default'
})

let resolver: ((value: boolean) => void) | null = null

export function useConfirm() {
  function ask(opts: ConfirmOptions): Promise<boolean> {
    state.title = opts.title
    state.description = opts.description ?? ''
    state.confirmText = opts.confirmText ?? 'Подтвердить'
    state.cancelText = opts.cancelText ?? 'Отмена'
    state.variant = opts.variant ?? 'default'
    state.open = true
    return new Promise<boolean>((resolve) => {
      resolver = resolve
    })
  }

  function confirm() {
    state.open = false
    resolver?.(true)
    resolver = null
  }

  function cancel() {
    state.open = false
    resolver?.(false)
    resolver = null
  }

  return { state, ask, confirm, cancel }
}
