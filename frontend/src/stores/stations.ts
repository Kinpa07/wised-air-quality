import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { Station, Pollutant } from "../types/station";

export const useStationsStore = defineStore("stations", () => {
  const stations = ref<Station[]>([]);
  const selectedStationId = ref<string | null>(null);
  const pollutant = ref<Pollutant>("pm2_5");
  const search = ref("");

  async function fetchStations() {
    const response = await fetch("/api/sensor-readings-collector/v1/stations");
    const json = await response.json();
    stations.value = json.data;
  }

  function selectStation(id: string) {
    selectedStationId.value = id;
  }

  function setPollutant(p: Pollutant) {
    pollutant.value = p;
  }

  function setSearch(query: string) {
    search.value = query;
  }

  const selectedStation = computed(() => {
    return stations.value.find((station) => station.id === selectedStationId.value);
  });

  const filteredStations = computed(() => {
    const query = search.value.trim().toLowerCase();
    if (!query) return stations.value;
    return stations.value.filter(
      (station) =>
        station.id.toLowerCase().includes(query) || station.district.toLowerCase().includes(query),
    );
  });

  return {
    stations,
    selectedStationId,
    pollutant,
    search,
    fetchStations,
    selectStation,
    setPollutant,
    setSearch,
    filteredStations,
    selectedStation,
  };
});
