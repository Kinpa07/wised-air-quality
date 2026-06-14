<script setup lang="ts">
import { computed } from "vue";
import Chart from "primevue/chart";
import { chart as chartColors } from "../../styles/tokens";

const props = defineProps<{
  values: number[];
}>();

const data = computed(() => ({
  labels: props.values.map((_, i) => i),
  datasets: [
    {
      data: props.values,
      borderColor: chartColors.line,
      borderWidth: 1.5,
      tension: 0.4,
      pointRadius: 0,
      fill: false,
    },
  ],
}));

const options = {
  responsive: true,
  maintainAspectRatio: false,
  animation: false,
  plugins: {
    legend: { display: false },
    tooltip: { enabled: false },
  },
  scales: {
    x: { display: false },
    y: { display: false },
  },
};
</script>

<template>
  <Chart type="line" :data="data" :options="options" class="sparkline" />
</template>

<style scoped>
.sparkline {
  width: 120px;
  height: 36px;
}
</style>
