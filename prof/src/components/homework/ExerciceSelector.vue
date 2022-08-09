<template>
  <v-card min-width="600">
    <v-row class="ma-1">
      <v-col>
        <v-card-title> Choisir un exercice </v-card-title>
      </v-col>
      <v-spacer></v-spacer>
      <v-col style="text-align: right">
        <v-btn icon @click="emit('close')" variant="flat">
          <v-icon icon="mdi-close"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-card-text>
      <v-list>
        <v-list-item
          v-for="exercice in exercices"
          :key="exercice.Exercice.Id"
          @click="emit('select', exercice)"
        >
          <v-row no-gutters>
            <v-col>
              <i>({{ exercice.Exercice.Id }}) </i>{{ exercice.Exercice.Title }}
            </v-col>
            <v-col style="text-align: right">
              <small>
                ({{ exercice.Questions?.length || 0 }}
                question(s))
              </small>
            </v-col>
          </v-row>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { ExerciceHeader } from "@/controller/api_gen";

interface Props {
  exercices: ExerciceHeader[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "select", ex: ExerciceHeader): void;
}>();
</script>
