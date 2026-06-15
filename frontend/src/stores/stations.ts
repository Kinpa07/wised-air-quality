import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { Station, Pollutant } from "../types/station";
import type { Stats } from "../types/stats";

export const useStationsStore = defineStore("stations", () => {
  const stations = ref<Station[]>([]);
  const stationsLoading = ref(false);
  const stationsError = ref(false);

  const stats = ref<Stats | null>(null);
  const statsLoading = ref(false);
  const statsError = ref(false);

  const selectedStationId = ref<string | null>(null);
  const pollutant = ref<Pollutant>("pm2_5");
  const search = ref("");

  async function fetchStations() {
    stationsLoading.value = true;
    stationsError.value = false;
    try {
      const response = await fetch("/api/sensor-readings-collector/v1/stations");
      if (!response.ok) throw new Error(`stations: ${response.status}`);
      const json = await response.json();
      stations.value = json.data;
    } catch {
      stationsError.value = true;
    } finally {
      stationsLoading.value = false;
    }
  }

  async function fetchStats() {
    statsLoading.value = true;
    statsError.value = false;
    try {
      const response = await fetch("/api/sensor-readings-collector/v1/stats");
      if (!response.ok) throw new Error(`stats: ${response.status}`);
      const json = await response.json();
      stats.value = json.data;
    } catch {
      statsError.value = true;
    } finally {
      statsLoading.value = false;
    }
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
    stationsLoading,
    stationsError,
    stats,
    statsLoading,
    statsError,
    selectedStationId,
    pollutant,
    search,
    fetchStations,
    fetchStats,
    selectStation,
    setPollutant,
    setSearch,
    filteredStations,
    selectedStation,
  };
});
