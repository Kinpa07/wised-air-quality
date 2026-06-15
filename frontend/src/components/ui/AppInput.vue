<script setup lang="ts">
import IconField from "primevue/iconfield";
import InputIcon from "primevue/inputicon";
import InputText from "primevue/inputtext";

defineProps<{
  modelValue: string;
  placeholder?: string;
  icon?: string;
}>();

defineEmits<{
  "update:modelValue": [value: string];
}>();
</script>

<template>
  <IconField>
    <InputIcon v-if="icon" :class="icon" />
    <InputText
      :model-value="modelValue"
      :placeholder="placeholder"
      fluid
      @update:model-value="$emit('update:modelValue', $event ?? '')"
    />
  </IconField>
</template>
