<script setup lang="ts">
import { ref, onMounted } from "vue";
import AppSparkline from "./ui/AppSparkline.vue";

const props = defineProps<{
  stationId: string;
}>();

const values = ref<number[]>([]);

onMounted(async () => {
  const since = new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString();
  const res = await fetch(
    `/api/sensor-readings-collector/v1/clients/${encodeURIComponent(props.stationId)}/readings?since=${since}`,
  );
  const json = await res.json();
  const readings = json.data as { pm2_5: number; timestamp: string }[];
  readings.sort((a, b) => a.timestamp.localeCompare(b.timestamp));
  values.value = readings.map((r) => r.pm2_5);
});
</script>

<template>
  <AppSparkline :values="values" />
</template>
