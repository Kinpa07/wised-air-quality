<script setup lang="ts">
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import AqiBadge from "./ui/AqiBadge.vue";
import StationSparkline from "./StationSparkline.vue";
import { useStationsStore } from "../stores/stations";
import type { Station } from "../types/station";

const store = useStationsStore();

function selectRow(station: Station | null) {
  if (station) store.selectStation(station.id);
}
</script>

<template>
  <DataTable
    :value="store.stations"
    data-key="id"
    selection-mode="single"
    :selection="store.selectedStation"
    scrollable
    scroll-height="400px"
    @update:selection="selectRow"
  >
    <Column header="District">
      <template #body="{ data }: { data: Station }">{{ data.district || "Unknown" }}</template>
    </Column>
    <Column header="Sensor">
      <template #body="{ data }: { data: Station }">
        <div>{{ data.id }}</div>
        <small>{{ data.lat }}, {{ data.lng }}</small>
      </template>
    </Column>
    <Column header="PM2.5">
      <template #body="{ data }: { data: Station }">
        {{ data.pm2_5 ?? "—" }} <AqiBadge :value="data.pm2_5" />
      </template>
    </Column>
    <Column header="24h">
      <template #body="{ data }: { data: Station }">
        <StationSparkline :station-id="data.id" />
      </template>
    </Column>
    <Column field="pm10" header="PM10" />
    <Column field="connection" header="Connection" />
    <Column header="Stability">
      <template #body="{ data }: { data: Station }">{{ data.stability.toFixed(0) }}%</template>
    </Column>
  </DataTable>
</template>
