<template>
  <v-text-field
    variant="outlined"
    density="compact"
    :model-value="text"
    @update:model-value="onTextChange"
    :color="color"
    :hide-details="!hint"
    :hint="hint"
    class="input-small"
    :label="label"
  >
  </v-text-field>
</template>

<script setup lang="ts">
import { colorByKind } from "@/controller/editor";
import type { TextPart } from "@/controller/exercice_gen";
import { TextKind } from "@/controller/exercice_gen";
import { computed } from "@vue/runtime-core";

interface Props {
  modelValue: TextPart;
  label?: string;
  hint?: string;
  forceLatex?: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: TextPart): void;
}>();

const text = computed(() => {
  switch (props.modelValue.Kind) {
    case TextKind.Text:
      return props.modelValue.Content;
    case TextKind.StaticMath:
      if (props.forceLatex) {
        return props.modelValue.Content;
      }
      return "$" + props.modelValue.Content + "$";
    case TextKind.Expression:
      return "!" + props.modelValue.Content + "!";
  }
});

const color = computed(() => colorByKind[props.modelValue.Kind]);

function onTextChange(s: string) {
  s = s.trim();
  if (s.startsWith("$") && s.endsWith("$") && s.length >= 3) {
    emit("update:model-value", {
      Kind: TextKind.StaticMath,
      Content: s.substring(1, s.length - 1)
    });
  } else if (s.startsWith("!") && s.endsWith("!") && s.length >= 3) {
    emit("update:model-value", {
      Kind: TextKind.Expression,
      Content: s.substring(2, s.length - 1)
    });
  } else {
    emit("update:model-value", {
      Kind: props.forceLatex ? TextKind.StaticMath : TextKind.Text,
      Content: s
    });
  }
}
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 12px;
}
.input-small:deep(.v-field__input) {
  padding: 0px 6px;
}
</style>
