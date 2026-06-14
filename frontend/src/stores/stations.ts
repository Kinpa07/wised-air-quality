import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { Station } from "../types/station";

export const useStationsStore = defineStore("stations", () => {
  const stations = ref<Station[]>([]);
  const selectedStationId = ref<string | null>(null);

  async function fetchStations() {
    const response = await fetch("/api/sensor-readings-collector/v1/stations");
    const json = await response.json();
    stations.value = json.data;
  }

  function selectStation(id: string) {
    selectedStationId.value = id;
  }

  const selectedStation = computed(() => {
    return stations.value.find((station) => station.id === selectedStationId.value);
  });

  return { stations, selectedStationId, fetchStations, selectStation, selectedStation };
});
