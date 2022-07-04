<template>
  <v-dialog
    :model-value="exerciceToDelete != null"
    @update:model-value="exerciceToDelete = null"
  >
    <v-card title="Confirmer">
      <v-card-text
        >Etes-vous certain de vouloir supprimer l'exercice
        <i>{{ exerciceToDelete?.Exercice.Title }}</i> ? <br />
        Cette opération est irréversible.
      </v-card-text>
      <v-card-actions>
        <v-btn @click="exerciceToDelete = null">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn
          color="red"
          @click="deleteExercice(true)"
          variant="contained"
          v-if="exerciceToDelete?.Questions?.length"
        >
          Supprimer aussi les questions
        </v-btn>
        <v-btn
          color="orange"
          @click="deleteExercice(false)"
          variant="contained"
        >
          Supprimer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="my-5 mx-auto" width="80%">
    <v-row>
      <v-col md="9" sm="6">
        <v-card-title>Liste des exercices</v-card-title>
      </v-col>

      <v-col align-self="center" style="text-align: right" md="3" sm="6">
        <v-btn class="mx-2" @click="createExercice" title="Créer un exercice">
          <v-icon icon="mdi-plus" color="success"></v-icon>
          Créer
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text>
      <v-list>
        <exercice-row
          v-for="(exercice, index) in exercices"
          :exercice="exercice"
          :key="index"
          @clicked="emit('clicked', exercice)"
          @delete="exerciceToDelete = exercice"
          @update-public="updatePublic"
        ></exercice-row>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { Exercice, ExerciceExt } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted } from "vue";
import { $ref } from "vue/macros";
import ExerciceRow from "./ExerciceRow.vue";

const emit = defineEmits<{
  (e: "clicked", exercice: ExerciceExt): void;
}>();

let exercices = $ref<ExerciceExt[]>([]);

onMounted(() => {
  fetchExercices();
});

async function fetchExercices() {
  const res = await controller.ExercicesGetList();
  exercices = res || [];
}

async function createExercice() {
  const newEx = await controller.ExerciceCreate();
  if (newEx == undefined) {
    return;
  }
  emit("clicked", newEx);
  await fetchExercices();
}

let exerciceToDelete = $ref<ExerciceExt | null>(null);
async function deleteExercice(deleteQuestions: boolean) {
  if (exerciceToDelete == null) {
    return;
  }
  await controller.ExerciceDelete({
    id: exerciceToDelete.Exercice.Id,
    delete_questions: deleteQuestions,
  });
  exerciceToDelete = null;
  await fetchExercices();
}

let exerciceToUpdate = $ref<Exercice | null>(null);
async function updateExercice() {
  if (exerciceToUpdate == null) {
    return;
  }
  const res = await controller.ExerciceUpdate(exerciceToUpdate);
  exerciceToUpdate = null;
  if (res == undefined) {
    return;
  }

  const index = exercices.findIndex((cl) => cl.Exercice.Id == res.Id);
  exercices[index].Exercice = res;
}

// let exerciceToShow = $ref<ExerciceExt | null>(null);

async function updatePublic(questionID: number, isPublic: boolean) {
  // TODO:
  //   const res = await controller.QuestionUpdateVisiblity({
  //     QuestionID: questionID,
  //     Public: isPublic,
  //   });
  //   if (res === undefined) {
  //     return;
  //   }
  //   questions.forEach((group) => {
  //     const index = group.Questions?.findIndex((qu) => qu.Id == questionID);
  //     if (index !== undefined) {
  //       group.Questions![index].Origin.IsPublic = isPublic;
  //     }
  //   });
}
</script>

<style>
:deep(.v-dialog .v-overlay__content) {
  max-width: 900px;
}
</style>
