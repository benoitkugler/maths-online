<template>
  <v-card title="Modifier la variante">
    <v-card-text>
      <v-row>
        <v-col>
          <v-text-field
            class="my-2"
            variant="outlined"
            density="compact"
            label="Sous-titre"
            v-model="props.variant.Subtitle"
            :readonly="props.readonly"
            hint="Le sous-titre est affiché uniquement sur l'interface Professeur."
          ></v-text-field
        ></v-col>
      </v-row>
      <v-row>
        <v-col align-self="center">
          <DifficultyField
            class="px-1"
            v-model="props.variant.Difficulty"
            :readonly="props.readonly"
          ></DifficultyField>
        </v-col>
      </v-row>
    </v-card-text>
    <v-card-actions>
      <v-btn @click="emit('close')">Retour</v-btn>
      <v-spacer></v-spacer>
      <v-btn
        color="success"
        :disabled="props.readonly"
        @click="emit('save', props.variant)"
        >Enregistrer</v-btn
      >
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import DifficultyField from "./utils/DifficultyField.vue";
import type { VariantG } from "@/controller/editor";

interface Props {
  variant: VariantG;
  readonly: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "save", v: VariantG): void;
}>();
</script>

<style scoped></style>
