<template>
  <div class="mx-0">
    <small v-if="props.label" class="ml-2 text-grey">{{ props.label }} </small>
    <HiglightedText
      :model-value="props.modelValue"
      @update:model-value="s => emit('update:modelValue', s)"
      :tokenizer="props.customTokenize ? props.customTokenize : defautTokenize"
      :focus-color="activeColor"
      :center="!!props.center"
      :class="props.hint ? 'mb-3' : ''"
    ></HiglightedText>
    <small v-if="props.hint" class="text-grey">{{ props.hint }} </small>
  </div>
</template>

<script setup lang="ts">
import { colorByKind } from "@/controller/editor";
import { TextKind } from "@/controller/api_gen";
import { computed } from "vue";
import { defautTokenize, type Token } from "./interpolated_text";
import HiglightedText from "./HiglightedText.vue";

type Props = {
  modelValue: string;
  label?: string;
  hint?: string;
  forceLatex?: boolean;
  center?: boolean;
  customTokenize?: (input: string) => Token[];
};

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:modelValue", modelValue: string): void;
}>();

const colorLatex = colorByKind[TextKind.StaticMath];
const activeColor = computed(() => (props.forceLatex ? colorLatex : "#444444"));
</script>

<style></style>
