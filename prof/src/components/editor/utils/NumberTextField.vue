<template>
  <v-text-field
    class="number"
    variant="outlined"
    density="compact"
    :model-value="props.modelValue"
    @update:model-value="s => emit('update:model-value', Number(s))"
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
</script>

<style scoped>
.number:deep(input) {
  width: 70%;
}
</style>
