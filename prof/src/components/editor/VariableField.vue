<template>
  <v-text-field
    variant="outlined"
    density="compact"
    :suffix="props.suffix"
    :prefix="props.prefix"
    :label="label"
    hide-details
    :model-value="String.fromCodePoint(props.modelValue)"
    @update:model-value="onVariableChange"
    class="centered-input"
  ></v-text-field>
</template>

<script setup lang="ts">
import type { Variable } from "@/controller/exercice_gen";

interface Props {
  modelValue: Variable;
  suffix?: string;
  prefix?: string;
  label?: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: Variable): void;
}>();

function onVariableChange(s: string) {
  if (s != "") {
    const variable = s.codePointAt(0)!;
    emit("update:model-value", variable);
  }
}
</script>

<style scoped>
.centered-input:deep(input) {
  text-align: center;
}
.centered-input:deep(.v-field__input) {
  padding-left: 4px;
}
</style>
