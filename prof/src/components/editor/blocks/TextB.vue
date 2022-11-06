<template>
  <v-row>
    <v-col md="7">
      <v-row no-gutters>
        <v-col md="12">
          <interpolated-text
            :model-value="props.modelValue.Parts"
            @update:model-value="onTextChanged"
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
            density="compact"
            class="pr-2"
            hide-details
            label="Gras"
            :model-value="props.modelValue.Bold"
            @update:model-value="onBoldChanged"
          ></v-checkbox>
        </v-col>
        <v-col>
          <v-checkbox
            density="compact"
            class="pr-2"
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
import type { TextBlock } from "@/controller/api_gen";
import InterpolatedText from "../utils/InterpolatedText.vue";

interface Props {
  modelValue: TextBlock;
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
</script>

<style></style>
