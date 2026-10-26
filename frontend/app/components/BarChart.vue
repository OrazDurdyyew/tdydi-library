<script setup lang="ts">
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale
} from 'chart.js'
import type { ChartData, ChartOptions } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

withDefaults(
  defineProps<{
    data: ChartData<'bar'>
    options?: ChartOptions<'bar'>
  }>(),
  {
    options: () => ({
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { display: false } }
    })
  }
)
</script>

<template>
  <div class="relative h-64 w-full">
    <Bar :data="data" :options="options" />
  </div>
</template>
