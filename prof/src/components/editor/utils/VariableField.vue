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
import type { Variable } from "@/controller/api_gen";
import { variableToString } from "@/controller/editor";
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

const variableString = computed(() => variableToString(props.modelValue));

function onVariableChange(s: string) {
  if (s == "") {
    return;
  }
  const chunks = s.trim().split("_");
  if (chunks.length == 0) {
    return;
  }

  let indice = chunks.length >= 2 ? chunks[1].trim() : "";
  // for 'Vi', insert _ automatically
  if (chunks[0].length > 1) {
    indice = chunks[0].substr(1) + indice;
  }

  const variable = { Name: chunks[0].codePointAt(0)!, Indice: indice };

  emit("update:model-value", variable);
}
</script>

<style scoped>
.centered-input:deep(input) {
  text-align: center;
  font-size: 14px;
  width: 100%;
}
.centered-input:deep(.v-field__input) {
  padding-left: 4px;
  padding-right: 4px;
  width: 50px;
}
</style>
