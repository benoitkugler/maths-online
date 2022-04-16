<template>
  <v-row>
    <v-col md="8">
      <v-row no-gutters>
        <v-col md="12">
          <interpolated-text
            :model-value="props.modelValue.Parts"
            @update:model-value="onTextChanged"
            class="px-2"
          ></interpolated-text>
        </v-col>
        <v-col md="12">
          <small class="text-grey mt-1 d-block">
            Insérer du code LaTeX avec : $\frac{a}{b}$. Insérer une expression
            avec : !2x + 1!
          </small>
        </v-col>
      </v-row>
    </v-col>
    <v-col md="2">
      <v-checkbox
        class="pr-2"
        hide-details
        label="Comme conseil"
        :model-value="props.modelValue.IsHint"
        @update:model-value="onIsHintChanged"
      ></v-checkbox>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { TextBlock } from "@/controller/exercice_gen";
import InterpolatedText from "../InterpolatedText.vue";

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

function onIsHintChanged(b: boolean) {
  props.modelValue.IsHint = b;
  emit("update:modelValue", props.modelValue);
}
</script>

<style></style>
