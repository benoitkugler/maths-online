<template>
  <v-select
    label="Notation"
    density="compact"
    variant="outlined"
    :items="items.map((i) => i.text)"
    :model-value="items.find((v) => v.value == props.modelValue)?.text"
    @update:model-value="onChange"
    :hint="hint"
    persistent-hint
  ></v-select>
</template>

<script setup lang="ts">
import { Notation } from "@/controller/api_gen";
import { computed } from "vue";

interface Props {
  modelValue: Notation;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: Notation): void;
}>();

const labels: { [key in Notation]: string } = {
  [Notation.NoNotation]: "Pas de notation",
  [Notation.SuccessNotation]: "Notation en fonction de la réussite",
};

const items = Object.entries(labels).map((k) => ({
  value: Number(k[0]) as Notation,
  text: k[1],
}));

function onChange(v: string) {
  const kind = items.find((item) => item.text == v)!.value;
  emit("update:model-value", kind);
}

const hint = computed(() => {
  return {
    [Notation.NoNotation]: "La fiche n'est pas notée.",
    [Notation.SuccessNotation]:
      "Les points d'une question sont attribués si l'élève a réussi une fois la question (peu importe le nombre de tentatives).",
  }[props.modelValue];
});
</script>

<style scoped></style>
