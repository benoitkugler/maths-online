<template>
  <v-card title="Paramètres de la question">
    <v-card-text class="mt-2">
      <v-row>
        <v-col>
          <v-text-field
            label="Nombre de répétitions"
            density="compact"
            variant="outlined"
            type="number"
            min="1"
            hide-details
            v-model.number="props.monoquestion.NbRepeat"
          ></v-text-field>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-text-field
            label="Barème"
            density="compact"
            variant="outlined"
            type="number"
            min="0"
            hint="Points attribués pour une question"
            v-model.number="props.monoquestion.Bareme"
          >
          </v-text-field>
        </v-col>
      </v-row>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn color="success" @click="emit('update', props.monoquestion)">
        Enregistrer
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import type { Monoquestion } from "@/controller/api_gen";

interface Props {
  monoquestion: Monoquestion;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", monoquestion: Monoquestion): void;
}>();
</script>

<style scoped></style>
