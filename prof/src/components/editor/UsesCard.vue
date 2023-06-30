<template>
  <v-card>
    <v-card-title class="bg-warning"> Tâches concernées </v-card-title>
    <v-card-text>
      Le contenu est utilisé dans les tâches suivantes :
      <v-list class="rounded my-2">
        <v-list-item
          rounded
          v-for="(use, index) in props.uses"
          :key="index"
          link
          @click="emit('goTo', use.Sheet)"
        >
          <v-list-item-title>
            Feuille : <b>{{ use.Sheet.Title }}</b>
          </v-list-item-title>
          <v-list-item-subtitle>
            Niveau {{ use.Sheet.Level }}
          </v-list-item-subtitle>
        </v-list-item>
      </v-list>

      Si vous souhaitez réellement supprimer le contenu, veuillez d'abord le
      retirer de cette liste.
    </v-card-text>
    <v-card-actions>
      <v-btn @click="emit('back')" color="warning">Retour</v-btn>
      <!-- <v-spacer></v-spacer>
        <v-btn color="red" @click="deleteVariante" variant="outlined">
          OK
        </v-btn> -->
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import type { QuestionExerciceUses, Sheet } from "@/controller/api_gen";

interface Props {
  uses: QuestionExerciceUses;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
  (e: "goTo", sh: Sheet): void;
}>();
</script>

<style scoped></style>
