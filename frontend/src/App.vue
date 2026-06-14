<script setup lang="ts">
import { onMounted } from "vue";
import { useStationsStore } from "./stores/stations";
import StationMap from "./components/StationMap.vue";
import AppCard from "./components/ui/AppCard.vue";

const store = useStationsStore();
onMounted(() => {
  store.fetchStations();
});
</script>

<template>
  <main style="padding: 1rem">
    <AppCard title="Sensor Map" subtitle="Live AQI by station">
      <StationMap />
    </AppCard>
  </main>
</template>
