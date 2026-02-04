<template>
  <v-dialog v-model="showExerciceSelector" max-width="1200">
    <resource-selector
      :query="exerciceQuery"
      :tags="props.allTags"
      :mode="'exercices'"
      @closed="showExerciceSelector = false"
      @selected-variant="(v) => emit('addExercice', v)"
      @update:query="(q) => (exerciceQuery = q)"
    ></resource-selector>
  </v-dialog>

  <v-dialog v-model="showMonoquestionSelector" max-width="1200">
    <resource-selector
      :query="exerciceQuery"
      :tags="props.allTags"
      :mode="'questions'"
      @closed="showMonoquestionSelector = false"
      @selected-variant="(v) => emit('addMonoquestion', v)"
      @selected-group="(gr) => emit('addRandomMonoquestion', gr)"
      @update:query="(q) => (questionQuery = q)"
    ></resource-selector>
  </v-dialog>

  <v-dialog
    :model-value="taskToRemove != null"
    @update:model-value="taskToRemove = null"
    max-width="600px"
  >
    <v-card title="Confirmer">
      <v-card-text
        >Etes-vous certain de vouloir retirer la tâche
        <i>{{ taskToRemove?.Title }}</i> ? <br />
        La progression des {{ taskToRemove?.NbProgressions }} élève(s) sera
        perdue, et cette opération est irréversible.
      </v-card-text>
      <v-card-actions>
        <v-btn @click="taskToRemove = null" color="warning">Retour</v-btn>
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
        <v-card-title> Liste des tâches </v-card-title>
      </v-col>
      <v-col align-self="center" style="text-align: right">
        <v-btn
          size="small"
          class="my-1 mr-2"
          @click="showExerciceSelector = true"
          title="Ajouter un exercice..."
        >
          <v-icon color="green" icon="mdi-plus"></v-icon>
          Exercice
        </v-btn>
        <v-btn
          size="small"
          class="my-1 mr-2"
          title="Ajouter une question (avec répétitions)..."
          @click="showMonoquestionSelector = true"
        >
          <v-icon color="green" icon="mdi-plus"></v-icon>
          Question
        </v-btn>
      </v-col>
    </v-row>
    <v-card-text>
      <MissingResourcesHint
        :pattern="missingHint.Pattern || []"
        :missing-exercices="missingHint.MissingExercices || []"
        :missing-questions="missingHint.MissingQuestions || []"
      ></MissingResourcesHint>
      <v-list
        class="overflow-y-auto"
        style="max-height: 50vh"
        @dragend="onDragend"
      >
        <v-list-item v-if="!sheet.Tasks?.length" style="text-align: center">
          <i>Aucun exercice.</i>
        </v-list-item>

        <drop-zone
          v-if="showDropZone"
          @drop="(origin) => swap(origin, 0)"
        ></drop-zone>

        <div v-for="(task, index) in sheet.Tasks" :key="index">
          <v-list-item>
            <v-row no-gutters>
              <v-col cols="auto" align-self="center">
                <drag-icon
                  color="black"
                  @start="(e) => onItemDragStart(e, index)"
                ></drag-icon>
              </v-col>
              <v-col cols="auto" align-self="center">
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
                <v-tooltip>
                  <template v-slot:activator="{ isActive, props }">
                    <span v-on="{ isActive }" v-bind="props">
                      {{ task.Title }}
                    </span>
                  </template>
                  {{ taskTooltip(task) }}
                </v-tooltip>
              </v-col>
              <v-col cols="3" align-self="center" class="text-center">
                <TagIndex :tags="task.Tags || []"></TagIndex>
              </v-col>
              <v-col cols="auto" align-self="center" class="px-3">
                <span v-if="!task.NbProgressions" class="text-grey"
                  >Non démarrée</span
                >
                <v-chip v-else color="secondary" variant="elevated"
                  >Démarrée par {{ task.NbProgressions }}</v-chip
                >
              </v-col>
              <v-col
                cols="auto"
                align-self="center"
                style="text-align: right"
                class="pl-2"
              >
                <task-details-chip
                  :task="task"
                  @update-monoquestion="(qu) => emit('updateMonoquestion', qu)"
                  @update-random-monoquestion="
                    (qu) => emit('updateRandomMonoquestion', qu)
                  "
                >
                </task-details-chip>
              </v-col>
            </v-row>
          </v-list-item>
          <drop-zone
            v-if="showDropZone"
            @drop="(origin) => swap(origin, index + 1)"
          ></drop-zone>
        </div>

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
import {
  OriginKind,
  WorkKind,
  type Monoquestion,
  type Query,
  type SheetExt,
  type TagsDB,
  type TaskExt,
  type RandomMonoquestion,
  MissingTasksHint,
} from "@/controller/api_gen";
import {
  onDragListItemStart,
  sheetBareme,
  swapItems,
} from "@/controller/utils";
import DragIcon from "../DragIcon.vue";
import DropZone from "../DropZone.vue";
import { ref, watch } from "vue";
import ResourceSelector from "../ResourceSelector.vue";
import type { ResourceGroup, VariantG } from "@/controller/editor";
import TaskDetailsChip from "./TaskDetailsChip.vue";
import { controller } from "@/controller/controller";
import TagIndex from "../TagIndex.vue";
import MissingResourcesHint from "../MissingResourcesHint.vue";
import { onMounted } from "vue";

interface Props {
  sheet: SheetExt;
  allTags: TagsDB;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "addExercice", ex: VariantG): void;
  (e: "addMonoquestion", qu: VariantG): void;
  (e: "addRandomMonoquestion", qu: ResourceGroup): void;
  (e: "updateMonoquestion", qu: Monoquestion): void;
  (e: "updateRandomMonoquestion", qu: RandomMonoquestion): void;
  (e: "remove", task: TaskExt): void;
  (e: "reorder", tasks: TaskExt[]): void;
}>();

const taskToRemove = ref<TaskExt | null>(null);

const showMonoquestionSelector = ref(false);

const showExerciceSelector = ref(false);
const exerciceQuery = ref<Query>({
  TitleQuery: "",
  LevelTags: props.sheet.Sheet.Level ? [props.sheet.Sheet.Level] : [],
  ChapterTags: [],
  SubLevelTags: [],
  Origin: OriginKind.All,
  Matiere: controller.settings.FavoriteMatiere,
});
const questionQuery = ref<Query>({
  TitleQuery: "",
  LevelTags: props.sheet.Sheet.Level ? [props.sheet.Sheet.Level] : [],
  ChapterTags: [],
  SubLevelTags: [],
  Origin: OriginKind.All,
  Matiere: controller.settings.FavoriteMatiere,
});

onMounted(fetchMissingHint);

watch(props.sheet, () => {
  exerciceQuery.value.LevelTags = props.sheet.Sheet.Level
    ? [props.sheet.Sheet.Level]
    : [];
  questionQuery.value.LevelTags = props.sheet.Sheet.Level
    ? [props.sheet.Sheet.Level]
    : [];
  fetchMissingHint();
});

function removeExercice(index: number) {
  const task = props.sheet.Tasks![index];
  // ask confirmation if progression has started
  if (task.NbProgressions) {
    taskToRemove.value = task;
  } else {
    emit("remove", task);
  }
}

const showDropZone = ref(false);
function onItemDragStart(payload: DragEvent, index: number) {
  onDragListItemStart(payload, index);
  setTimeout(() => (showDropZone.value = true), 100); // workaround bug
}

function onDragend() {
  showDropZone.value = false;
}

function swap(origin: number, target: number) {
  showDropZone.value = false;
  const l = swapItems(origin, target, props.sheet.Tasks!);
  emit("reorder", l);
}

function taskTooltip(task: TaskExt) {
  switch (task.IdWork.Kind) {
    case WorkKind.WorkExercice:
      return `Exercice ${task.IdWork.ID}`;
    case WorkKind.WorkMonoquestion:
      return `Variante d'une question (répétée ${task.Bareme?.length} fois)`;
    case WorkKind.WorkRandomMonoquestion:
      return `Groupe de questions (${task.Bareme?.length} variantes aléatoires)`;
  }
}

const missingHint = ref<MissingTasksHint>({
  Pattern: [],
  MissingExercices: [],
  MissingQuestions: [],
});
async function fetchMissingHint() {
  const res = await controller.HomeworkMissingTasksHint({
    "id-sheet": props.sheet.Sheet.Id,
  });
  if (res === undefined) return;
  missingHint.value = res;
}
</script>
