<script setup lang="ts">
import { onMounted } from "vue";
import { useStationsStore } from "./stores/stations";
import StationMap from "./components/StationMap.vue";

const store = useStationsStore();
onMounted(() => {
  store.fetchStations();
});
</script>

<template>
  <main style="padding: 1rem">
    <h1>Wised Berlin Air-Quality Monitor</h1>
    <StationMap />
  </main>
</template>
