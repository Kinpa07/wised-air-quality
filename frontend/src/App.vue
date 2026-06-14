<script setup lang="ts">
import { onMounted } from "vue";
import { useStationsStore } from "./stores/stations";
import StationMap from "./components/StationMap.vue";
import AppCard from "./components/ui/AppCard.vue";
import AppToggle from "./components/ui/AppToggle.vue";
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
  <main style="padding: 1rem">
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
  </main>
</template>
