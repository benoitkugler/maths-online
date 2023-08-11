<template>
  <v-card class="my-2">
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
            @update:model-value="
              s => {
                completePoint(s, props.modelValue.Answer);
                emitUpdate();
              }
            "
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
            @update:model-value="emitUpdate()"
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row class="py-0">
        <v-col class="py-0">
          <v-switch
            v-model="props.modelValue.MustHaveOrigin"
            @update:model-value="emitUpdate()"
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
              s => {
                completePoint(s, props.modelValue.AnswerOrigin);
                emitUpdate();
              }
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
            @update:model-value="emitUpdate()"
          ></v-text-field>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { GFVector } from "@/controller/api_gen";
import { TextKind } from "@/controller/api_gen";
import { colorByKind, completePoint } from "@/controller/editor";

interface Props {
  modelValue: GFVector;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: GFVector): void;
}>();

const expressionColor = colorByKind[TextKind.Expression];

function emitUpdate() {
  emit("update:modelValue", props.modelValue);
}
</script>

<style scoped></style>
