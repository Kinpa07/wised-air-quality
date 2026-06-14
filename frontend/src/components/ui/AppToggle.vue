<script setup lang="ts" generic="T extends string">
import SelectButton from "primevue/selectbutton";

defineProps<{
  modelValue: T;
  options: { label: string; value: T }[];
}>();

defineEmits<{
  "update:modelValue": [value: T];
}>();
</script>

<template>
  <SelectButton
    :model-value="modelValue"
    :options="options"
    option-label="label"
    option-value="value"
    :allow-empty="false"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>
