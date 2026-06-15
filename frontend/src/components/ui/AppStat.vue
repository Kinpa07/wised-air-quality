<script setup lang="ts">
import Card from "primevue/card";
import { space, surface } from "../../styles/tokens";

defineProps<{
  icon: string;
  label: string;
  value: string;
  unit?: string;
  iconColor?: string;
}>();
</script>

<template>
  <Card>
    <template #content>
      <div class="app-stat" :style="{ gap: `${space.sm}px` }">
        <i :class="icon" :style="{ color: iconColor ?? surface.muted }" />
        <div class="app-stat__value" :style="{ color: surface.text }">
          {{ value
          }}<span
            v-if="unit"
            class="app-stat__unit"
            :style="{ color: surface.muted, marginLeft: `${space.xs}px` }"
            >{{ unit }}</span
          >
        </div>
        <div class="app-stat__label" :style="{ color: surface.muted }">{{ label }}</div>
      </div>
    </template>
  </Card>
</template>

<style scoped>
.app-stat {
  display: flex;
  flex-direction: column;
}
.app-stat__value {
  font-size: 1.75rem;
  font-weight: 600;
  line-height: 1.2;
}
.app-stat__unit {
  font-size: 0.9rem;
  font-weight: 400;
}
.app-stat__label {
  font-size: 0.8rem;
}
</style>
