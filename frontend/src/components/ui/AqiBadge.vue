<script setup lang="ts">
import { computed } from "vue";
import { colorFor, labelFor } from "../../utils/aqi";
import { surface, radius, space } from "../../styles/tokens";

const props = defineProps<{
  value: number | null;
}>();

const label = computed(() => (props.value === null ? "No data" : labelFor(props.value)));
const background = computed(() => (props.value === null ? surface.muted : colorFor(props.value)));
</script>

<template>
  <span
    class="aqi-badge"
    :style="{
      background,
      color: surface.card,
      borderRadius: `${radius.sm}px`,
      padding: `${space.xs}px ${space.sm}px`,
    }"
  >
    {{ label }}
  </span>
</template>

<style scoped>
.aqi-badge {
  display: inline-block;
  font-size: 0.75rem;
  font-weight: 600;
  white-space: nowrap;
}
</style>
