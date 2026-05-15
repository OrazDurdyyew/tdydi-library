<script setup lang="ts">
import { Swiper, SwiperSlide } from 'swiper/vue'
import { Navigation } from 'swiper/modules'
import 'swiper/css'
import 'swiper/css/navigation'

defineProps<{ books: any[] }>()
</script>

<template>
  <ClientOnly>
    <Swiper
      :modules="[Navigation]"
      navigation
      :space-between="20"
      :slides-per-view="2.2"
      :breakpoints="{
        480: { slidesPerView: 2.6 },
        640: { slidesPerView: 3.4 },
        768: { slidesPerView: 4.4 },
        1024: { slidesPerView: 6 }
      }"
      class="book-swiper"
    >
      <SwiperSlide v-for="book in books" :key="book.id">
        <BookCard :book="book" />
      </SwiperSlide>
    </Swiper>

    <template #fallback>
      <div class="grid grid-cols-3 md:grid-cols-6 gap-5">
        <BookCard v-for="book in books.slice(0, 6)" :key="book.id" :book="book" />
      </div>
    </template>
  </ClientOnly>
</template>

<style scoped>
.book-swiper {
  padding: 0.25rem 0 0.5rem;
}
.book-swiper :deep(.swiper-button-next),
.book-swiper :deep(.swiper-button-prev) {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 9999px;
  background-color: rgb(var(--card));
  border: 1px solid rgb(var(--border));
  color: rgb(var(--foreground));
  box-shadow: 0 4px 12px rgb(15 23 42 / 0.12);
  transition: background-color 140ms ease, color 140ms ease;
}
.book-swiper :deep(.swiper-button-next:hover),
.book-swiper :deep(.swiper-button-prev:hover) {
  background-color: rgb(var(--primary));
  color: rgb(var(--primary-foreground));
}
.book-swiper :deep(.swiper-button-next::after),
.book-swiper :deep(.swiper-button-prev::after) {
  font-size: 0.85rem;
  font-weight: 700;
}
.book-swiper :deep(.swiper-button-disabled) {
  opacity: 0;
}
</style>
