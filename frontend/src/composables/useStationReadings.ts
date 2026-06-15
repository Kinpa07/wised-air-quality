import { ref, watch, type Ref } from "vue";
import type { Reading } from "../types/reading";

export type ChartRange = "24h" | "7d" | "30d";

const RANGE_HOURS: Record<ChartRange, number> = {
  "24h": 24,
  "7d": 24 * 7,
  "30d": 24 * 30,
};

export function useStationReadings(stationId: Ref<string | null>, range: Ref<ChartRange>) {
  const readings = ref<Reading[]>([]);
  const loading = ref(false);
  const error = ref(false);

  async function load() {
    if (!stationId.value) {
      readings.value = [];
      return;
    }
    loading.value = true;
    error.value = false;
    try {
      const since = new Date(Date.now() - RANGE_HOURS[range.value] * 3600 * 1000).toISOString();
      const res = await fetch(
        `/api/sensor-readings-collector/v1/clients/${encodeURIComponent(stationId.value)}/readings?since=${since}`,
      );
      if (!res.ok) throw new Error(`readings: ${res.status}`);
      const json = await res.json();
      // Endpoint returns newest-first; the chart plots oldest -> newest.
      readings.value = (json.data as Reading[])
        .slice()
        .sort((a, b) => a.timestamp.localeCompare(b.timestamp));
    } catch {
      error.value = true;
      readings.value = [];
    } finally {
      loading.value = false;
    }
  }

  watch([stationId, range], load, { immediate: true });

  return { readings, loading, error };
}
