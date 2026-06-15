<script setup lang="ts">
import { useStationsStore } from "../stores/stations";
import AppInput from "./ui/AppInput.vue";
import { space, surface, sizes } from "../styles/tokens";

const store = useStationsStore();
</script>

<template>
  <header class="dashboard-header" :style="{ gap: `${space.md}px` }">
    <h1 class="dashboard-header__title" :style="{ color: surface.text }">
      Berlin · Particulate Monitoring
    </h1>
    <AppInput
      :model-value="store.search"
      icon="pi pi-search"
      placeholder="Search district or sensor ID"
      :style="{ width: `${sizes.search}px` }"
      @update:model-value="store.setSearch"
    />
  </header>
</template>

<style scoped>
.dashboard-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.dashboard-header__title {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
}
</style>
