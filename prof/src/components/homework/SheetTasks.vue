<template>
  <v-dialog v-model="showSelector">
    <exercice-selector
      @close="showSelector = false"
      @select="addExercice"
      :exercices="allExercices"
    ></exercice-selector>
  </v-dialog>

  <v-dialog
    :model-value="taskToRemove != null"
    @update:model-value="taskToRemove = null"
  >
    <v-card title="Confirmer">
      <v-card-text
        >Etes-vous certain de vouloir retirer l'exercice
        <i>{{ taskToRemove?.Exercice.Exercice.Title }}</i> ? <br />
        La progression des {{ taskToRemove?.NbProgressions }} élève(s) sera
        perdue, et cette opération est irréversible.
      </v-card-text>
      <v-card-actions>
        <v-btn @click="taskToRemove = null">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn
          color="red"
          @click="
            taskToRemove = null;
            emit('remove', taskToRemove!);
          "
          variant="elevated"
        >
          Supprimer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card>
    <v-row>
      <v-col align-self="center">
        <v-card-title> Liste des exercices </v-card-title>
      </v-col>
      <v-col align-self="center" style="text-align: right">
        <v-btn class="my-1 mr-2" @click="showSelector = true">
          <v-icon color="green" icon="mdi-plus"></v-icon>
          Ajouter
        </v-btn>
      </v-col>
    </v-row>
    <v-card-text>
      <v-list class="overflow-y-auto" style="max-height: 50vh">
        <v-list-item v-if="!sheet.Tasks?.length" style="text-align: center">
          <i>Aucun exercice.</i>
        </v-list-item>
        <v-list-item v-for="(task, index) in sheet.Tasks" :key="index">
          <v-row>
            <v-col cols="3" align-self="center">
              <v-btn
                icon
                size="small"
                variant="flat"
                @click="removeExercice(index)"
                title="Retirer l'exercice"
              >
                <v-icon color="red" icon="mdi-close"></v-icon>
              </v-btn>
            </v-col>
            <v-col align-self="center">
              {{ task.Exercice.Exercice.Title }}
            </v-col>
            <v-col align-self="center">
              <span v-if="!task.NbProgressions" class="text-grey"
                >Non démarré</span
              >
              <v-chip v-else color="secondary"
                >Démarré par {{ task.NbProgressions }}</v-chip
              >
            </v-col>
            <v-col align-self="center" style="text-align: right">
              / {{ exerciceBareme(task.Exercice) }}
            </v-col>
          </v-row>
        </v-list-item>
        <v-list-item v-if="sheet.Tasks?.length">
          <v-row>
            <v-col align-self="center">Total</v-col>
            <v-spacer></v-spacer>
            <v-col align-self="center" style="text-align: right">
              / {{ sheetBareme(sheet) }}
            </v-col>
          </v-row>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { ExerciceHeader, SheetExt, TaskExt } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { exerciceBareme, sheetBareme } from "@/controller/utils";
import { onMounted } from "vue";
import { $ref } from "vue/macros";
import ExerciceSelector from "./ExerciceSelector.vue";

interface Props {
  sheet: SheetExt;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "add", ex: ExerciceHeader): void;
  (e: "remove", task: TaskExt): void;
  (e: "reorder", tasks: TaskExt[]): void;
}>();

onMounted(fetchExercices);

let taskToRemove = $ref<TaskExt | null>(null);

let showSelector = $ref(false);
let allExercices = $ref<ExerciceHeader[]>([]);
async function fetchExercices() {
  const res = await controller.ExercicesGetList();
  if (res == undefined) {
    return;
  }
  allExercices = res;
}

function addExercice(ex: ExerciceHeader) {
  emit("add", ex);
}

function removeExercice(index: number) {
  const task = props.sheet.Tasks![index];
  // ask confirmation if progression has started
  if (task.NbProgressions) {
    taskToRemove = task;
  } else {
    emit("remove", task);
  }
}

// TODO: implement reorder
</script>
