<template>
  <v-text-field
    variant="outlined"
    density="compact"
    :suffix="props.suffix"
    :prefix="props.prefix"
    :label="label"
    hide-details
    :model-value="variableString"
    @update:model-value="onVariableChange"
    class="centered-input"
  ></v-text-field>
</template>

<script setup lang="ts">
import type { Variable } from "@/controller/exercice_gen";
import { computed } from "@vue/runtime-core";

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

const variableString = computed(() => {
  let name = String.fromCodePoint(props.modelValue.Name);
  if (props.modelValue.Indice) {
    name += "_" + props.modelValue.Indice;
  }
  return name;
});

function onVariableChange(s: string) {
  if (s == "") {
    return;
  }
  const chunks = s.trim().split("_");
  if (chunks.length == 0) {
    return;
  }

  const variable = { Name: chunks[0].codePointAt(0)!, Indice: "" };
  if (chunks.length >= 2) {
    variable.Indice = chunks[1].trim();
  }

  emit("update:model-value", variable);
}
</script>

<style scoped>
.centered-input:deep(input) {
  text-align: center;
  font-size: 14px;
}
.centered-input:deep(.v-field__input) {
  padding-left: 4px;
  padding-right: 4px;
  width: 50px;
}
</style>
