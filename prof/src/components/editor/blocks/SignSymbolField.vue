<template>
  <v-btn rounded @click="nextSymbol" size="small" class="my-2" block>
    {{ symbolItems[props.modelValue] }}
  </v-btn>
</template>

<script setup lang="ts">
import type { SignSymbol } from "@/controller/api_gen";
import { SignSymbolLabels } from "@/controller/api_gen";

interface Props {
  modelValue: SignSymbol;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: SignSymbol): void;
}>();

const symbolItems = SignSymbolLabels;

function nextSymbol() {
  const symbol = (props.modelValue + 1) % Object.keys(symbolItems).length;
  emit("update:model-value", symbol as SignSymbol);
}
</script>

<style scoped>
.centered-input:deep(input) {
  text-align: center;
}
</style>
