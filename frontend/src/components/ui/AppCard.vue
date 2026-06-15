<script setup lang="ts">
import Card from "primevue/card";
import { sizes } from "../../styles/tokens";

defineProps<{
  title: string;
  subtitle?: string;
}>();
</script>

<template>
  <Card>
    <template #title>
      <div class="app-card__title-row" :style="{ minHeight: `${sizes.titleRow}px` }">
        <span>{{ title }}</span>
        <slot name="actions" />
      </div>
    </template>
    <template v-if="subtitle" #subtitle>{{ subtitle }}</template>
    <template #content>
      <slot />
    </template>
  </Card>
</template>

<style scoped>
.app-card__title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
