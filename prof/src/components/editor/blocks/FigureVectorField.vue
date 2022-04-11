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
            label="X (Vecteur)"
            hint="Expression, comparée à l'unité prés."
            :color="expressionColor"
            v-model="props.modelValue.Answer.X"
            @update:model-value="s => completePoint(s, props.modelValue.Answer)"
          ></v-text-field>
        </v-col>
        <v-col cols="6">
          <v-text-field
            variant="outlined"
            density="compact"
            label="Y (Vecteur)"
            hint="Expression, comparée à l'unité prés."
            :color="expressionColor"
            v-model="props.modelValue.Answer.Y"
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row class="py-0">
        <v-col class="py-0">
          <v-switch
            v-model="props.modelValue.MustHaveOrigin"
            color="secondary"
            label="Evaluer le point d'origine du vecteur."
            hide-details
          >
          </v-switch>
        </v-col>
      </v-row>
      <v-row v-if="props.modelValue.MustHaveOrigin">
        <v-col cols="6">
          <v-text-field
            variant="outlined"
            density="compact"
            label="X (Origine)"
            hint="Expression, comparée à l'unité prés."
            :color="expressionColor"
            v-model="props.modelValue.AnswerOrigin.X"
            @update:model-value="
              s => completePoint(s, props.modelValue.AnswerOrigin)
            "
          ></v-text-field>
        </v-col>
        <v-col cols="6">
          <v-text-field
            variant="outlined"
            density="compact"
            label="Y (Origine)"
            hint="Expression, comparée à l'unité prés."
            :color="expressionColor"
            v-model="props.modelValue.AnswerOrigin.Y"
          ></v-text-field>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
  <figure-vue
    v-model="props.modelValue.Figure"
    :available-parameters="props.availableParameters"
  ></figure-vue>
</template>

<script setup lang="ts">
import { colorByKind, completePoint } from "@/controller/editor";
import type {
  FigureVectorFieldBlock,
  Variable
} from "@/controller/exercice_gen";
import { TextKind } from "@/controller/exercice_gen";
import FigureVue from "./Figure.vue";

interface Props {
  modelValue: FigureVectorFieldBlock;
  availableParameters: Variable[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FigureVectorFieldBlock): void;
}>();

const expressionColor = colorByKind[TextKind.Expression];
</script>

<style scoped></style>
