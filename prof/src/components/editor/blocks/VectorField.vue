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
            hint="Expression, comparée comme nombre à virgule."
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
            label="Y (Vecteur)"
            hint="Expression, comparée comme nombre à virgule."
            :color="expressionColor"
            v-model="props.modelValue.Answer.Y"
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row class="py-0">
        <v-col class="py-0" cols="12">
          <v-switch
            density="compact"
            v-model="props.modelValue.AcceptColinear"
            color="secondary"
            label="Accepter un vecteur colinéaire (non nul)"
            hide-details
          >
          </v-switch>
        </v-col>
        <v-col class="py-0" cols="12">
          <v-switch
            density="compact"
            v-model="props.modelValue.DisplayColumn"
            color="secondary"
            label="Afficher en colonne"
            hide-details
          >
          </v-switch>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { colorByKind, completePoint } from "@/controller/editor";
import type { Variable, VectorFieldBlock } from "@/controller/exercice_gen";
import { TextKind } from "@/controller/exercice_gen";

interface Props {
  modelValue: VectorFieldBlock;
  availableParameters: Variable[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: VectorFieldBlock): void;
}>();

const expressionColor = colorByKind[TextKind.Expression];
</script>

<style scoped></style>
