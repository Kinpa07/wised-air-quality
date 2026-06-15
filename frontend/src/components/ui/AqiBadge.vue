<script setup lang="ts">
import { computed } from "vue";
import Tag from "primevue/tag";
import { colorFor, labelFor } from "../../utils/aqi";
import { surface, sizes } from "../../styles/tokens";

const props = defineProps<{
  value: number | null;
}>();

const label = computed(() => (props.value === null ? "No data" : labelFor(props.value)));
const background = computed(() => (props.value === null ? surface.muted : colorFor(props.value)));
</script>

<template>
  <Tag
    :value="label"
    :style="{
      background,
      color: surface.card,
      minWidth: `${sizes.badge}px`,
      justifyContent: 'center',
    }"
  />
</template>
