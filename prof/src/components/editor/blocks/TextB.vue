<template>
  <v-row>
    <v-col md="7">
      <v-row no-gutters>
        <v-col md="12">
          <interpolated-text
            :model-value="props.modelValue.Parts"
            @update:model-value="onTextChanged"
            :custom-tokenize="tokenizeText"
          ></interpolated-text>
        </v-col>
        <v-col md="12">
          <small class="text-grey mt-1 d-block">
            Insérer du code LaTeX avec : $\frac{a}{b}$. Insérer une expression
            avec : &2x + 1&
          </small>
        </v-col>
      </v-row>
    </v-col>
    <v-col md="5">
      <v-row no-gutters>
        <v-col>
          <v-checkbox
            class="pr-1"
            density="compact"
            hide-details
            label="Gras"
            :model-value="props.modelValue.Bold"
            @update:model-value="onBoldChanged"
          ></v-checkbox>
        </v-col>
        <v-col>
          <v-checkbox
            class="pr-0"
            density="compact"
            hide-details
            label="Italique"
            :model-value="props.modelValue.Italic"
            @update:model-value="onItalicChanged"
          ></v-checkbox>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-checkbox
            density="compact"
            class="pr-2"
            hide-details
            label="Taille réduite"
            :model-value="props.modelValue.Smaller"
            @update:model-value="onSmallerChanged"
          ></v-checkbox>
        </v-col>
      </v-row>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { TextBlock, Variable } from "@/controller/api_gen";
import InterpolatedText from "../utils/InterpolatedText.vue";
import {
  itemize,
  defautTokenize,
  partToToken,
  type Token
} from "../utils/interpolated_text";
import { TextKind } from "@/controller/loopback_gen";

interface Props {
  modelValue: TextBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: TextBlock): void;
}>();

function onTextChanged(s: string) {
  props.modelValue.Parts = s;
  emit("update:modelValue", props.modelValue);
}

function onBoldChanged(b: boolean) {
  props.modelValue.Bold = b;
  emit("update:modelValue", props.modelValue);
}
function onItalicChanged(b: boolean) {
  props.modelValue.Italic = b;
  emit("update:modelValue", props.modelValue);
}
function onSmallerChanged(b: boolean) {
  props.modelValue.Smaller = b;
  emit("update:modelValue", props.modelValue);
}

// support for formulas
function tokenizeText(input: string) {
  const out: Token[] = [];
  const lines = input.split("\n");
  let currentLines: string[] = [];

  lines.forEach((line, index) => {
    const lineT = line.trim();
    if (lineT != "$$" && lineT.startsWith("$$") && lineT.endsWith("$$")) {
      // found a formula
      if (currentLines.length) {
        // flush the previous block
        // apply the regular tokenization
        const regularBlock = currentLines.join("\n") + "\n";

        out.push(...defautTokenize(regularBlock));
        currentLines = [];
      }

      const start = line.indexOf("$$");
      const end = line.lastIndexOf("$$");
      const style = "color: blue;";
      const innerFormula = itemize(lineT.substring(2, lineT.length - 2));
      const tokens = innerFormula.map(tp =>
        tp.Kind == TextKind.Expression
          ? partToToken(tp)
          : { Content: tp.Content, Kind: style }
      );
      out.push(
        { Content: line.substring(0, start + 2), Kind: style },
        ...tokens,
        {
          Content:
            index == lines.length - 1
              ? line.substring(end)
              : line.substring(end) + "\n",
          Kind: style
        }
      );
    } else {
      currentLines.push(line);
    }
  });
  if (currentLines.length) {
    out.push(...defautTokenize(currentLines.join("\n")));
  }

  return out;
}
</script>

<style></style>
