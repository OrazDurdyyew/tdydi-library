<script setup lang="ts">
const { state, cancel, confirm } = useConfirm()
</script>

<template>
  <Dialog v-if="state.open" :open="state.open" @update:open="(v) => !v && cancel()">
    <DialogHeader>
      <DialogTitle>{{ state.title }}</DialogTitle>
      <DialogDescription v-if="state.description">{{ state.description }}</DialogDescription>
    </DialogHeader>
    <DialogFooter>
      <Button variant="outline" @click="cancel">{{ state.cancelText }}</Button>
      <Button :variant="state.variant === 'destructive' ? 'destructive' : 'default'" @click="confirm">
        {{ state.confirmText }}
      </Button>
    </DialogFooter>
  </Dialog>
</template>
