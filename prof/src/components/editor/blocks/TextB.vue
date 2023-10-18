<template>
  <v-row>
    <v-col>
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
            Insérer une expression avec & &, insérer du code LaTeX avec $ $ ou
            $$ $$.
          </small>
        </v-col>
      </v-row>
    </v-col>
    <v-col cols="auto">
      <div class="mt-2">
        <v-tooltip text="Gras">
          <template v-slot:activator="{ isActive, props: inner }">
            <v-btn-toggle
              v-on="{ isActive }"
              v-bind="inner"
              density="comfortable"
              :model-value="props.modelValue.Bold ? 0 : -1"
              @update:model-value="i => onBoldChanged(i == 0)"
            >
              <v-btn icon="mdi-format-bold" class="py-1"></v-btn>
            </v-btn-toggle>
          </template>
        </v-tooltip>
      </div>

      <div>
        <v-tooltip text="Italique">
          <template v-slot:activator="{ isActive, props: inner }">
            <v-btn-toggle
              v-on="{ isActive }"
              v-bind="inner"
              density="comfortable"
              :model-value="props.modelValue.Italic ? 0 : -1"
              @update:model-value="i => onItalicChanged(i == 0)"
            >
              <v-btn icon="mdi-format-italic" class="py-1"></v-btn>
            </v-btn-toggle>
          </template>
        </v-tooltip>
      </div>

      <div>
        <v-tooltip text="Taille réduite">
          <template v-slot:activator="{ isActive, props: inner }">
            <v-btn-toggle
              v-on="{ isActive }"
              v-bind="inner"
              density="comfortable"
              :model-value="props.modelValue.Smaller ? 0 : -1"
              @update:model-value="i => onSmallerChanged(i == 0)"
            >
              <v-btn icon="mdi-format-font-size-decrease" class="py-1"></v-btn>
            </v-btn-toggle>
          </template>
        </v-tooltip>
      </div>
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

<style scoped>
.v-btn-toggle {
  flex-direction: column;
}
</style>
