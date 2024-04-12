<template>
  <v-card :subtitle="props.sheet.Sheet.Title">
    <template v-slot:append>
      <v-btn
        :icon="showQuestions ? 'mdi-chevron-up' : 'mdi-chevron-down'"
        size="small"
        @click="showQuestions = !showQuestions"
      ></v-btn>
    </template>
    <v-card-text>
      <v-card v-for="(task, j) in props.data.TaskStats" :key="j" class="my-2">
        <v-card-text>
          <div class="text-subtitle-2">{{ task.Title }}</div>
          <PercentageBar
            class="mt-4"
            :success="task.NbSuccess"
            :failure="task.NbFailure"
            :height="30"
          ></PercentageBar>
          <div v-if="showQuestions" class="mt-4">
            <v-card-subtitle class="text-center mb-2">
              {{ workSubtitle(task.IdWork) }}</v-card-subtitle
            >

            <v-row v-for="(question, index) in task.QuestionStats" :key="index">
              <v-col cols="2" class="text-center">
                <TagChip :tag="{ Tag: question.Difficulty }"></TagChip>
              </v-col>
              <v-col>
                <small>({{ question.Id }})</small> {{ question.Description }}
              </v-col>
              <v-col cols="4">
                <PercentageBar
                  :success="question.NbSuccess"
                  :failure="question.NbFailure"
                  :height="20"
                ></PercentageBar>
              </v-col>
            </v-row>
          </div>
        </v-card-text>
      </v-card>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { SheetExt, TravailMarks, WorkID } from "@/controller/api_gen";
import { WorkKind } from "@/controller/api_gen";
import PercentageBar from "./PercentageBar.vue";
import { ref } from "vue";
import TagChip from "../editor/utils/TagChip.vue";

interface Props {
  data: TravailMarks;
  sheet: SheetExt;
}

const props = defineProps<Props>();

const showQuestions = ref(false);

function workSubtitle(work: WorkID) {
  switch (work.Kind) {
    case WorkKind.WorkExercice:
      return "Détails de l'exercice";
    case WorkKind.WorkMonoquestion:
      return "Question unique";
    case WorkKind.WorkRandomMonoquestion:
      return "Détails des variantes";
  }
}
</script>
