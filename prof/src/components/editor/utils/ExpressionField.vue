<template>
  <v-text-field
    variant="outlined"
    density="compact"
    :model-value="props.modelValue"
    @update:model-value="onTextChange"
    :color="color"
    :hide-details="!hint"
    :hint="hint"
    class="expression-field-input"
    :label="label"
    :prefix="props.prefix"
  >
  </v-text-field>
</template>

<script setup lang="ts">
import { TextKind } from "@/controller/api_gen";
import { colorByKind } from "@/controller/editor";
import { computed } from "vue";

interface Props {
  modelValue: string;
  label?: string;
  hint?: string;
  prefix?: string;
  center?: boolean;
  width?: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: string): void;
}>();

const exprSeparator = "&";

const color = colorByKind[TextKind.Expression];
const align = computed(() => (props.center ? "center" : "left"));

function onTextChange(s: string) {
  s = s.trim();
  if (s.startsWith(exprSeparator) && s.endsWith(exprSeparator)) {
    // remove the unwanted &&
    s = s.substring(1, s.length - 1);
  }
  emit("update:model-value", s);
}
</script>

<style scoped>
.expression-field-input:deep(input) {
  font-size: 14px;
  text-align: v-bind(align);
  width: v-bind("props.width");
}
</style>
