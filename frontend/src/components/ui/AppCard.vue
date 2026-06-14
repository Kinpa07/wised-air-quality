<script setup lang="ts">
import Card from "primevue/card";

defineProps<{
  title: string;
  subtitle?: string;
}>();
</script>

<template>
  <Card>
    <template #title>{{ title }}</template>
    <template v-if="subtitle" #subtitle>{{ subtitle }}</template>
    <template #content>
      <slot />
    </template>
  </Card>
</template>
