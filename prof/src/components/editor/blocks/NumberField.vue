<template>
  <v-row class="mt-2">
    <v-col>
      <v-text-field
        variant="outlined"
        density="compact"
        :model-value="props.modelValue.Expression"
        @update:model-value="
          (v) => emit('update:modelValue', { Expression: v })
        "
        label="Réponse"
        hint="Expression s'évaluant en un nombre (à virgule flottante). Ne crée pas de nouvelle ligne."
        persistent-hint
        :color="color"
      >
      </v-text-field>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { NumberFieldBlock, Variable } from "@/controller/api_gen";
import { TextKind } from "@/controller/api_gen";
import { colorByKind } from "@/controller/editor";

interface Props {
  modelValue: NumberFieldBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();
const color = colorByKind[TextKind.Expression];

const emit = defineEmits<{
  (event: "update:modelValue", value: NumberFieldBlock): void;
}>();
</script>

<style></style>
