<script setup lang="ts">
import { computed } from "vue";
import Skeleton from "primevue/skeleton";
import AppStat from "./ui/AppStat.vue";
import { useStationsStore } from "../stores/stations";
import { colorFor, labelFor } from "../utils/aqi";
import { space, surface } from "../styles/tokens";

const store = useStationsStore();

interface Tile {
  icon: string;
  label: string;
  value: string;
  unit?: string;
  iconColor?: string;
}

const tiles = computed((): Tile[] => {
  const s = store.stats;
  if (!s) return [];
  const avg = s.avg_pm2_5;
  return [
    { icon: "pi pi-wifi", label: "Active Sensors", value: String(s.active_sensors) },
    avg === null
      ? { icon: "pi pi-cloud", label: "Avg PM2.5", value: "—" }
      : {
          icon: "pi pi-cloud",
          label: `Avg PM2.5 · ${labelFor(avg)}`,
          value: avg.toFixed(1),
          unit: "µg",
          iconColor: colorFor(avg),
        },
    {
      icon: "pi pi-exclamation-triangle",
      label: "Poor Connection",
      value: String(s.poor_connection),
    },
    {
      icon: "pi pi-chart-line",
      label: "Network Stability",
      value: `${s.network_stability.toFixed(1)}%`,
    },
  ];
});
</script>

<template>
  <div class="kpi-panel" :style="{ gap: `${space.md}px` }">
    <template v-if="store.statsLoading">
      <Skeleton v-for="n in 4" :key="n" height="6rem" />
    </template>
    <div v-else-if="store.statsError" class="kpi-panel__error" :style="{ color: surface.muted }">
      Couldn’t load fleet stats.
    </div>
    <template v-else>
      <AppStat
        v-for="tile in tiles"
        :key="tile.label"
        :icon="tile.icon"
        :label="tile.label"
        :value="tile.value"
        :unit="tile.unit"
        :icon-color="tile.iconColor"
      />
    </template>
  </div>
</template>

<style scoped>
.kpi-panel {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
}
.kpi-panel__error {
  grid-column: 1 / -1;
}
@media (max-width: 900px) {
  .kpi-panel {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
