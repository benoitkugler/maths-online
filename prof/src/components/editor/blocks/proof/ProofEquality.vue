<template>
  <InterpolatedText
    :model-value="props.modelValue.Terms"
    @update:model-value="(s) => emit('update:modelValue', { Terms: s })"
    force-latex
    center
    :custom-tokenize="tokenizeText"
  ></InterpolatedText>
</template>

<script setup lang="ts">
import { TextKind, type ProofEquality } from "@/controller/api_gen";
import InterpolatedText from "../../utils/InterpolatedText.vue";
import {
  itemize,
  partToToken,
  splitByRegexp,
  type Token,
} from "../../utils/interpolated_text";

interface Props {
  modelValue: ProofEquality;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: ProofEquality): void;
}>();

function tokenizeText(input: string) {
  const parts = itemize(input);
  const out: Token[] = [];

  parts.forEach((part) => {
    if (part.Kind == TextKind.Text) {
      const matchSep = /=/g;
      const matchStyle = "color: blue; font-weight: bold";
      // do not split after avec
      const index = part.Content.indexOf("avec");
      if (index != -1) {
        const toColorize = part.Content.substring(0, index);
        out.push(...splitByRegexp(matchSep, toColorize, matchStyle, ""));
        out.push({
          Content: part.Content.substring(index),
          Kind: "font-style: italic",
        });
        return;
      }
      out.push(...splitByRegexp(matchSep, part.Content, matchStyle, ""));
    } else {
      // defaut to regular color
      out.push(partToToken(part));
    }
  });

  return out;
}
</script>

<style></style>
