<template>
  <v-text-field
    class="number"
    variant="outlined"
    density="compact"
    :model-value="props.modelValue"
    @update:model-value="onChange"
    type="number"
    :label="label"
    hide-details
  ></v-text-field>
</template>

<script setup lang="ts">
interface Props {
  modelValue: number;
  label?: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: number): void;
}>();

function isNumber(value: string | number): boolean {
  return value != null && value !== "" && !isNaN(Number(value.toString()));
}

function onChange(s: string) {
  if (!isNumber(s)) {
    return;
  }
  const out = Number(s);
  emit("update:model-value", out);
}
</script>

<style scoped>
.number:deep(input) {
  width: 70%;
}
</style>
