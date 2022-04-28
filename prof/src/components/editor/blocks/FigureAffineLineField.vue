<template>
  <v-card class="my-1">
    <v-card-subtitle class="bg-secondary py-3"
      >Définition de la réponse</v-card-subtitle
    >
    <v-card-text>
      <v-row class="fix-input-width">
        <v-col align-self="center" cols="4">
          <v-text-field
            density="compact"
            variant="outlined"
            label="Légende"
            v-model="props.modelValue.Label"
          ></v-text-field>
        </v-col>
        <v-col align-self="center">
          <v-text-field
            density="compact"
            variant="outlined"
            label="A"
            hint="Expression du coefficient directeur"
            v-model="props.modelValue.A"
            :color="expressionColor"
            class="no-hint-padding"
          ></v-text-field>
        </v-col>
        <v-col align-self="center">
          <v-text-field
            density="compact"
            variant="outlined"
            label="B"
            v-model="props.modelValue.B"
            hint="Expression de l'ordonnée à l'origine"
            :color="expressionColor"
            class="no-hint-padding"
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
import { colorByKind } from "@/controller/editor";
import type {
  FigureAffineLineFieldBlock,
  Variable
} from "@/controller/exercice_gen";
import { TextKind } from "@/controller/exercice_gen";
import FigureVue from "./Figure.vue";

interface Props {
  modelValue: FigureAffineLineFieldBlock;
  availableParameters: Variable[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FigureAffineLineFieldBlock): void;
}>();

const expressionColor = colorByKind[TextKind.Expression];
</script>

<style scoped>
.no-hint-padding:deep(.v-input__details) {
  padding-inline: 0px;
}

.fix-input-width:deep(input) {
  width: 100%;
}
</style>
