<template>
  <v-dialog v-model="showSelector">
    <exercice-selector
      @close="showSelector = false"
      @select="addExercice"
      :exercices="allExercices"
    ></exercice-selector>
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
        <v-list-item v-if="!sheet.Exercices?.length" style="text-align: center">
          <i>Aucun exercice.</i>
        </v-list-item>
        <v-list-item v-for="(exercice, index) in sheet.Exercices" :key="index">
          <v-row>
            <v-col cols="3" align-self="center">
              <v-btn
                icon
                size="small"
                variant="flat"
                @click="removeExercice(index)"
              >
                <v-icon color="red" icon="mdi-close"></v-icon>
              </v-btn>
            </v-col>
            <v-col align-self="center">
              {{ exercice.Exercice.Title }}
            </v-col>
            <v-col align-self="center" style="text-align: right">
              / {{ exerciceBareme(exercice) }}
            </v-col>
          </v-row>
        </v-list-item>
        <v-list-item v-if="sheet.Exercices?.length">
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
import type { ExerciceHeader, SheetExt } from "@/controller/api_gen";
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
  (e: "update", exercices: ExerciceHeader[]): void;
}>();

onMounted(fetchExercices);

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
  const exes = (props.sheet.Exercices || []).concat(ex);
  emit("update", exes);
}

function removeExercice(index: number) {
  props.sheet.Exercices?.splice(index, 1);
  emit("update", props.sheet.Exercices || []);
}
</script>
