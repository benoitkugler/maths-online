<template>
  <InterpolatedText
    :model-value="props.modelValue.Terms"
    @update:model-value="(s) => emit('update:modelValue', { Terms: s })"
    force-latex
    center
    :transform="textTransform"
  ></InterpolatedText>
</template>

<script setup lang="ts">
import type { ProofEquality } from "@/controller/api_gen";
import type Quill from "quill";
import InterpolatedText from "../../utils/InterpolatedText.vue";

interface Props {
  modelValue: ProofEquality;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: ProofEquality): void;
}>();

function textTransform(quill: Quill) {
  const text = quill.getText();

  const avecSep = "avec";
  const avecIndex = text.indexOf(avecSep);
  if (avecIndex != -1) {
    quill.formatText(avecIndex, avecSep.length, {
      bold: true,
      color: "blue",
    });
  }

  for (const match of text.matchAll(/=/g)) {
    // do not highlight = after avec sep
    if (avecIndex != -1 && match.index! >= avecIndex) {
      continue;
    }
    quill.formatText(match.index!, match.length, {
      bold: true,
      color: "blue",
    });
  }
}
</script>

<style></style>
