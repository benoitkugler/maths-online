<template>
  <v-card class="my-1">
    <v-card-subtitle class="bg-secondary py-3"
      >Coordonnées de la réponse</v-card-subtitle
    >
    <v-card-text>
      <v-row>
        <v-col cols="6">
          <v-text-field
            variant="outlined"
            density="compact"
            label="X"
            hint="Expression, comparée à l'unité prés."
            :color="expressionColor"
            v-model="props.modelValue.Answer.X"
            @update:model-value="
              (s) => completePoint(s, props.modelValue.Answer)
            "
          ></v-text-field>
        </v-col>
        <v-col cols="6">
          <v-text-field
            variant="outlined"
            density="compact"
            label="Y"
            hint="Expression, comparée à l'unité prés."
            :color="expressionColor"
            v-model="props.modelValue.Answer.Y"
          ></v-text-field>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
  <figure-block-vue
    v-model="props.modelValue.Figure"
    @update:model-value="emit('update:modelValue', props.modelValue)"
    :available-parameters="props.availableParameters"
  ></figure-block-vue>
</template>

<script setup lang="ts">
import type { FigurePointFieldBlock, Variable } from "@/controller/api_gen";
import { TextKind } from "@/controller/api_gen";
import { colorByKind, completePoint } from "@/controller/editor";
import FigureBlockVue from "./FigureBlock.vue";

interface Props {
  modelValue: FigurePointFieldBlock;
  availableParameters: Variable[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FigurePointFieldBlock): void;
}>();

const expressionColor = colorByKind[TextKind.Expression];
</script>

<style scoped></style>
