<script setup lang="ts">
import { surface, radius, space } from "../../styles/tokens";

defineProps<{
  title: string;
  subtitle?: string;
}>();
</script>

<template>
  <section
    class="app-card"
    :style="{
      background: surface.card,
      border: `1px solid ${surface.border}`,
      borderRadius: `${radius.lg}px`,
      padding: `${space.md}px`,
    }"
  >
    <header class="app-card__head" :style="{ marginBottom: `${space.sm}px` }">
      <h2 class="app-card__title" :style="{ color: surface.text }">{{ title }}</h2>
      <p v-if="subtitle" class="app-card__subtitle" :style="{ color: surface.muted }">
        {{ subtitle }}
      </p>
    </header>
    <slot />
  </section>
</template>

<style scoped>
.app-card__title {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
}
.app-card__subtitle {
  margin: 0;
  font-size: 0.8rem;
}
</style>
