<script setup lang="ts">
import { computed } from "vue";
import { chart as chartColors } from "../../styles/tokens";

const props = defineProps<{
  values: number[];
}>();

const WIDTH = 80;
const HEIGHT = 32;
const PAD = 2; // keep the stroke off the edges

// Map the values to a polyline point string, scaled to the box. Flat input
// (all equal) avoids a divide-by-zero via the `|| 1` span fallback.
const points = computed(() => {
  const vals = props.values;
  if (vals.length === 0) return "";
  const min = Math.min(...vals);
  const max = Math.max(...vals);
  const span = max - min || 1;
  const stepX = vals.length > 1 ? (WIDTH - PAD * 2) / (vals.length - 1) : 0;
  return vals
    .map((v, i) => {
      const x = PAD + i * stepX;
      const y = PAD + (1 - (v - min) / span) * (HEIGHT - PAD * 2);
      return `${x.toFixed(1)},${y.toFixed(1)}`;
    })
    .join(" ");
});
</script>

<template>
  <svg class="sparkline" :viewBox="`0 0 ${WIDTH} ${HEIGHT}`">
    <polyline
      v-if="points"
      :points="points"
      fill="none"
      :stroke="chartColors.line"
      stroke-width="1.5"
      stroke-linejoin="round"
      stroke-linecap="round"
    />
  </svg>
</template>

<style scoped>
.sparkline {
  display: block;
  width: 80px;
  height: 32px;
}
</style>
