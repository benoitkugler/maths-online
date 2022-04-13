<template>
  <v-row class="mt-2">
    <v-col>
      <v-text-field
        variant="outlined"
        density="compact"
        v-model="props.modelValue.Expression"
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
import { colorByKind } from "@/controller/editor";
import type { NumberFieldBlock } from "@/controller/exercice_gen";
import { TextKind } from "@/controller/exercice_gen";

interface Props {
  modelValue: NumberFieldBlock;
}
const props = defineProps<Props>();
const color = colorByKind[TextKind.Expression];

const emit = defineEmits<{
  (event: "update:modelValue", value: NumberFieldBlock): void;
}>();
</script>

<style></style>
