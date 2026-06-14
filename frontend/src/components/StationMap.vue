<script setup lang="ts">
import { onMounted, watch } from "vue";
import L from "leaflet";
import { useStationsStore } from "../stores/stations";
import { colorFor, AQI_BANDS } from "../utils/aqi";
import { markerNoData, surface, aqiColors, space, radius } from "../styles/tokens";
import type { Station } from "../types/station";

const store = useStationsStore();
let map: L.Map;

const markersById = new Map<string, L.CircleMarker>();

function fillFor(station: Station): string {
  const value = station[store.pollutant];
  return value !== null ? colorFor(value) : markerNoData;
}

function styleBase(marker: L.CircleMarker, station: Station) {
  const color = fillFor(station);
  marker.setStyle({ color, fillColor: color, weight: 1, fillOpacity: 0.8 });
  marker.setRadius(7);
}

function styleSelected(marker: L.CircleMarker) {
  marker.setStyle({ color: surface.text, weight: 3 });
  marker.setRadius(11);
  marker.bringToFront();
}

onMounted(() => {
  map = L.map("map").setView([52.52, 13.405], 11);
  L.tileLayer("https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}.png", {
    attribution: "© OpenStreetMap contributors © CARTO",
    subdomains: "abcd",
    maxZoom: 19,
  }).addTo(map);
});

// Draw (or redraw) markers whenever the fleet data OR the selected pollutant changes.
watch(
  () => [store.stations, store.pollutant],
  () => {
    if (!map) return;

    markersById.forEach((marker) => marker.remove());
    markersById.clear();

    store.stations.forEach((station) => {
      const color = fillFor(station);
      const marker = L.circleMarker([station.lat, station.lng], {
        color,
        fillColor: color,
        weight: 1,
        fillOpacity: 0.8,
        radius: 7,
      })
        .addTo(map)
        .bindTooltip(
          `<strong>${station.district}</strong><br>${station.id}<br>` +
            `PM2.5 ${station.pm2_5 ?? "—"} · PM10 ${station.pm10 ?? "—"}`,
        )
        .on("click", () => store.selectStation(station.id));
      markersById.set(station.id, marker);
    });

    // Re-apply the highlight if a station was already selected before a redraw.
    if (store.selectedStationId) {
      const marker = markersById.get(store.selectedStationId);
      if (marker) styleSelected(marker);
    }
  },
);

// Move the highlight when the shared selection changes (from the map or elsewhere).
watch(
  () => store.selectedStationId,
  (newId, oldId) => {
    if (oldId) {
      const prev = markersById.get(oldId);
      const station = store.stations.find((s) => s.id === oldId);
      if (prev && station) styleBase(prev, station);
    }
    if (newId) {
      const next = markersById.get(newId);
      if (next) {
        styleSelected(next);
        map.panTo(next.getLatLng());
      }
    }
  },
);
</script>

<template>
  <div class="map-wrap">
    <div id="map" style="height: 500px"></div>
    <div
      class="legend"
      :style="{
        background: surface.card,
        border: `1px solid ${surface.border}`,
        borderRadius: `${radius.md}px`,
        padding: `${space.sm}px ${space.md}px`,
      }"
    >
      <div
        v-for="b in AQI_BANDS"
        :key="b.band"
        class="legend__row"
        :style="{ color: surface.text }"
      >
        <span class="legend__swatch" :style="{ background: aqiColors[b.band] }"></span>
        {{ b.label }} · {{ b.max === Infinity ? `${b.min}+` : `${b.min}–${b.max}` }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.map-wrap {
  position: relative;
}
.legend {
  position: absolute;
  bottom: 12px;
  left: 12px;
  z-index: 1000;
  font-size: 0.8rem;
}
.legend__row {
  display: flex;
  align-items: center;
  gap: 6px;
}
.legend__swatch {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}
</style>
