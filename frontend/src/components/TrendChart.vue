<script setup lang="ts">
import { ref, computed } from "vue";
import { storeToRefs } from "pinia";
import Chart from "primevue/chart";
import Skeleton from "primevue/skeleton";
import AppToggle from "./ui/AppToggle.vue";
import { useStationsStore } from "../stores/stations";
import { useStationReadings, type ChartRange } from "../composables/useStationReadings";
import { chart as chartColors, surface, sizes } from "../styles/tokens";

const store = useStationsStore();
const { selectedStationId } = storeToRefs(store);

const range = ref<ChartRange>("24h");
const rangeOptions: { value: ChartRange; label: string }[] = [
  { value: "24h", label: "24H" },
  { value: "7d", label: "7D" },
  { value: "30d", label: "30D" },
];

const { readings, loading, error } = useStationReadings(selectedStationId, range);

function formatLabel(ts: string): string {
  return new Date(ts).toLocaleString([], {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

const data = computed(() => ({
  labels: readings.value.map((r) => formatLabel(r.timestamp)),
  datasets: [
    {
      type: "bar",
      label: "PM10",
      data: readings.value.map((r) => r.pm10),
      backgroundColor: chartColors.bar,
      order: 2,
    },
    {
      type: "line",
      label: "PM2.5",
      data: readings.value.map((r) => r.pm2_5),
      borderColor: chartColors.line,
      borderWidth: 2,
      tension: 0.4,
      pointRadius: 0,
      order: 1,
    },
  ],
}));

const options = {
  responsive: true,
  maintainAspectRatio: false,
  interaction: { mode: "index", intersect: false },
  plugins: {
    legend: { position: "bottom", labels: { color: surface.muted } },
  },
  scales: {
    x: { grid: { display: false }, ticks: { color: surface.muted, maxTicksLimit: 8 } },
    y: { beginAtZero: true, grid: { color: surface.border }, ticks: { color: surface.muted } },
  },
};
</script>

<template>
  <div class="trend">
    <div class="trend__header">
      <AppToggle v-model="range" :options="rangeOptions" />
    </div>
    <div class="trend__body" :style="{ height: `${sizes.chart}px` }">
      <Skeleton v-if="loading" width="100%" height="100%" />
      <div v-else-if="error" class="trend__msg" :style="{ color: surface.muted }">
        Couldn’t load readings.
      </div>
      <div v-else-if="!selectedStationId" class="trend__msg" :style="{ color: surface.muted }">
        Select a station to see its particulate trend.
      </div>
      <div v-else-if="readings.length === 0" class="trend__msg" :style="{ color: surface.muted }">
        No readings in this window.
      </div>
      <Chart v-else type="bar" :data="data" :options="options" class="trend__chart" />
    </div>
  </div>
</template>

<style scoped>
.trend {
  display: flex;
  flex-direction: column;
}
.trend__header {
  display: flex;
  justify-content: flex-end;
}
.trend__msg {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}
.trend__chart {
  height: 100%;
}
</style>
