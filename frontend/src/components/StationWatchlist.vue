<script setup lang="ts">
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import AqiBadge from "./ui/AqiBadge.vue";
import StationSparkline from "./StationSparkline.vue";
import { useStationsStore } from "../stores/stations";
import { space, sizes, surface } from "../styles/tokens";
import type { Station } from "../types/station";

const store = useStationsStore();

function selectRow(station: Station | null) {
  if (station) store.selectStation(station.id);
}
</script>

<template>
  <!-- Fixed height = the map's, so the cards align; scroll-height="flex" makes
       the table body fill what's left under the header. -->
  <div class="watchlist-wrap" :style="{ height: `${sizes.map}px` }">
    <DataTable
      :value="store.filteredStations"
      data-key="id"
      selection-mode="single"
      :selection="store.selectedStation"
      scrollable
      scroll-height="flex"
      size="small"
      @update:selection="selectRow"
    >
      <Column header="District">
        <template #body="{ data }: { data: Station }">{{ data.district || "Unknown" }}</template>
      </Column>
      <Column header="Sensor">
        <template #body="{ data }: { data: Station }">
          <div :title="data.id">{{ data.id.slice(0, 8) }}…</div>
          <small :style="{ color: surface.muted }"
            >{{ data.lat.toFixed(3) }}, {{ data.lng.toFixed(3) }}</small
          >
        </template>
      </Column>
      <Column header="PM2.5" header-class="pm-col">
        <template #body="{ data }: { data: Station }">
          <div class="pm-cell" :style="{ gap: `${space.sm}px`, color: surface.text }">
            <span class="pm-cell__value">{{ data.pm2_5 ?? "—" }}</span>
            <AqiBadge :value="data.pm2_5" />
          </div>
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
  </div>
</template>

<style scoped>
/* Centre value+badge (and the header, below) in the wide PM2.5 column. */
.pm-cell {
  display: flex;
  align-items: center;
  justify-content: center;
}
.watchlist-wrap :deep(.pm-col .p-datatable-column-header-content) {
  justify-content: center;
}
.pm-cell__value {
  font-weight: 600;
  min-width: 40px; /* fixed so the badge after it always starts at the same x */
  text-align: right;
}
</style>
