<script setup lang="ts">
import { onMounted } from "vue";
import { useStationsStore } from "./stores/stations";
import DashboardHeader from "./components/DashboardHeader.vue";
import StationMap from "./components/StationMap.vue";
import StationWatchlist from "./components/StationWatchlist.vue";
import AppCard from "./components/ui/AppCard.vue";
import AppToggle from "./components/ui/AppToggle.vue";
import { space, sizes } from "./styles/tokens";
import type { Pollutant } from "./types/station";

const store = useStationsStore();

const pollutantOptions: { value: Pollutant; label: string }[] = [
  { value: "pm2_5", label: "PM2.5" },
  { value: "pm10", label: "PM10" },
];

onMounted(() => {
  store.fetchStations();
});
</script>

<template>
  <main
    class="dashboard"
    :style="{ padding: `${space.md}px`, gap: `${space.md}px`, maxWidth: `${sizes.maxWidth}px` }"
  >
    <DashboardHeader />
    <div class="dashboard__grid" :style="{ gap: `${space.md}px` }">
      <AppCard title="Sensor Map" subtitle="Live AQI by station">
        <template #actions>
          <AppToggle
            :model-value="store.pollutant"
            :options="pollutantOptions"
            @update:model-value="store.setPollutant"
          />
        </template>
        <StationMap />
      </AppCard>
      <AppCard title="Districts" subtitle="Tracked stations">
        <StationWatchlist />
      </AppCard>
    </div>
  </main>
</template>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  margin: 0 auto;
}
.dashboard__grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  align-items: start;
}
/* min-width: 0 overrides the grid item default (auto) so the wide table can't
   push past its track. */
.dashboard__grid > * {
  min-width: 0;
}
</style>
