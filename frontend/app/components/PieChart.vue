<script setup lang="ts">
import { Pie } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, ArcElement } from 'chart.js'
import type { ChartData, ChartOptions } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, ArcElement)

withDefaults(
  defineProps<{
    data: ChartData<'pie'>
    options?: ChartOptions<'pie'>
  }>(),
  {
    options: () => ({
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { position: 'bottom' } }
    })
  }
)
</script>

<template>
  <div class="relative h-64 w-full">
    <Pie :data="data" :options="options" />
  </div>
</template>
