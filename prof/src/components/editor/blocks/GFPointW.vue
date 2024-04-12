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
            label="X"
            hint="Expression, comparée à l'unité prés."
            :color="expressionColor"
            v-model="props.modelValue.Answer.X"
            @update:model-value="
              s => {
                completePoint(s, props.modelValue.Answer);
                emit('update:modelValue', props.modelValue);
              }
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
            @update:model-value="emit('update:modelValue', props.modelValue)"
          ></v-text-field>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { GFPoint } from "@/controller/api_gen";
import { TextKind } from "@/controller/api_gen";
import { colorByKind, completePoint } from "@/controller/editor";

interface Props {
  modelValue: GFPoint;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: GFPoint): void;
}>();

const expressionColor = colorByKind[TextKind.Expression];
</script>

<style scoped></style>
