<template>
  <v-select
    label="Position de la légende"
    density="compact"
    variant="outlined"
    :items="items.map((i) => i.text)"
    :model-value="items.find((v) => v.value == props.modelValue)?.text"
    @update:model-value="onChange"
    hide-details
  ></v-select>
</template>

<script setup lang="ts">
import { LabelPosLabels, type LabelPos } from "@/controller/api_gen";

interface Props {
  modelValue: LabelPos;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: LabelPos): void;
}>();

const items = Object.entries(LabelPosLabels).map((k) => ({
  value: Number(k[0]) as LabelPos,
  text: k[1],
}));

function onChange(v: string) {
  const pos = items.find((item) => item.text == v)!.value;
  emit("update:model-value", pos);
}
</script>

<style scoped></style>
