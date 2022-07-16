<template>
  <v-select
    label="Type de tracé"
    density="compact"
    variant="outlined"
    :items="items.map((i) => i.text)"
    :model-value="items.find((v) => v.value == props.modelValue)?.text"
    @update:model-value="onChange"
    hide-details
  ></v-select>
</template>

<script setup lang="ts">
import { SegmentKindLabels, type SegmentKind } from "@/controller/api_gen";

interface Props {
  modelValue: SegmentKind;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: SegmentKind): void;
}>();

const items = Object.entries(SegmentKindLabels).map((k) => ({
  value: Number(k[0]) as SegmentKind,
  text: k[1],
}));

function onChange(v: string) {
  const kind = items.find((item) => item.text == v)!.value;
  emit("update:model-value", kind);
}
</script>

<style scoped></style>
