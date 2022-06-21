<template>
  <v-dialog v-model="showEditDescription">
    <description-pannel
      v-model="props.exercice.Exercice.Description"
      :readonly="isReadonly"
    ></description-pannel>
  </v-dialog>

  <v-dialog v-model="showImportQuestion" :retain-focus="false">
    <keep-alive>
      <question-selector
        :tags="props.allTags"
        :query="query"
        @closed="showImportQuestion = false"
        @selected="
          (q) => {
            showImportQuestion = false;
            addQuestion(q.Id);
          }
        "
      ></question-selector>
    </keep-alive>
  </v-dialog>

  <v-card class="mt-3 px-2">
    <v-row no-gutters class="mb-2">
      <v-col cols="auto" align-self="center" class="pr-2">
        <v-btn
          size="small"
          icon
          title="Retour aux exercices"
          @click="backToList"
        >
          <v-icon icon="mdi-arrow-left"></v-icon>
        </v-btn>
      </v-col>

      <v-col>
        <v-row no-gutters>
          <v-col>
            <v-text-field
              class="my-2 input-small"
              variant="outlined"
              density="compact"
              label="Nom de l'exercice"
              v-model="props.exercice.Exercice.Title"
              :readonly="isReadonly"
              hide-details
            ></v-text-field
          ></v-col>
          <v-col cols="auto" align-self="center">
            <v-btn
              class="mx-2"
              icon
              @click="save"
              :title="
                isReadonly ? 'Visualiser' : 'Enregistrer et prévisualiser'
              "
              size="small"
            >
              <v-icon
                :icon="isReadonly ? 'mdi-eye' : 'mdi-content-save'"
                size="small"
              ></v-icon>
            </v-btn>

            <v-menu offset-y close-on-content-click>
              <template v-slot:activator="{ isActive, props }">
                <v-btn
                  icon
                  title="Plus d'options"
                  v-on="{ isActive }"
                  v-bind="props"
                  size="x-small"
                >
                  <v-icon icon="mdi-dots-vertical"></v-icon>
                </v-btn>
              </template>
              <v-list>
                <v-list-item>
                  <v-btn
                    size="small"
                    @click="showEditDescription = true"
                    title="Editer le commentaire"
                  >
                    <v-icon
                      class="mr-2"
                      icon="mdi-message-reply-text"
                      size="small"
                    ></v-icon>
                    Commentaire
                  </v-btn>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-col>
        </v-row>

        <v-row no-gutters>
          <!-- <v-col class="pr-2">
            <tag-list-field
              label="Catégories"
              v-model="tags"
              :all-tags="props.allTags"
              @update:model-value="saveTags"
              :readonly="isReadonly"
            ></tag-list-field
          ></v-col> -->
          <v-col cols="auto">
            <v-btn
              title="Créer et ajouter une question"
              size="small"
              @click="createQuestion()"
            >
              <v-icon icon="mdi-plus" color="green"></v-icon>
              Ajouter une question
            </v-btn>
            <v-btn
              title="Importer une question existante"
              size="small"
              class="mx-2"
              @click="showImportQuestion = true"
            >
              <v-icon icon="mdi-plus" color="green"></v-icon>
              Importer une question
            </v-btn>
          </v-col>
        </v-row>
      </v-col>
    </v-row>

    <v-list @dragstart="onDragStart" @dragend="onDragEnd">
      <drop-zone
        v-if="showDropZone"
        @drop="(origin) => swapQuestions(origin, 0)"
      ></drop-zone>

      <div v-for="(question, index) in props.exercice.Questions" :key="index">
        <v-list-item
          draggable="true"
          @dragstart="(e) => onItemDragStart(e, index)"
        >
          <v-row>
            <v-col cols="auto" align-self="center">
              <v-btn
                v-if="!isReadonly"
                size="x-small"
                icon
                @click.stop="removeQuestion(index)"
                title="Supprimer"
              >
                <v-icon icon="mdi-delete" color="red" size="small"></v-icon>
              </v-btn>
            </v-col>
            <v-col> {{ question.Title }}</v-col>
            <v-col cols="3">
              <v-menu
                offset-y
                close-on-content-click
                :model-value="questionIndexToEdit == index"
                @update:model-value="
                  questionIndexToEdit = null;
                  questionToEdit = null;
                "
              >
                <template v-slot:activator="{ isActive, props }">
                  <v-chip
                    v-on="{ isActive }"
                    v-bind="props"
                    color="primary"
                    @click="
                      questionToEdit = copy(question.Question);
                      questionIndexToEdit = index;
                    "
                    :disabled="isReadonly"
                    >/ {{ question.Question.bareme }}
                  </v-chip>
                </template>
                <v-card subtitle="Modifier le barème">
                  <v-card-text>
                    <v-text-field
                      v-if="questionToEdit != null"
                      variant="outlined"
                      density="compact"
                      type="number"
                      label="Barème"
                      hide-details
                      v-model.number="questionToEdit.bareme"
                    ></v-text-field>
                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="success" @click="saveEditedQuestion"
                      >Enregistrer</v-btn
                    >
                  </v-card-actions>
                </v-card>
              </v-menu>
            </v-col>
          </v-row>
        </v-list-item>

        <drop-zone
          v-if="showDropZone"
          @drop="(origin) => swapQuestions(origin, index + 1)"
        ></drop-zone>
      </div>
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import {
  Visibility,
  type ExerciceExt,
  type ExerciceQuestion,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onDragListItemStart, swapItems } from "@/controller/editor";
import { copy } from "@/controller/utils";
import { computed } from "vue";
import { $ref } from "vue/macros";
import DescriptionPannel from "../editor/DescriptionPannel.vue";
import DropZone from "../editor/DropZone.vue";
import QuestionSelector, { type Query } from "./QuestionSelector.vue";

interface Props {
  exercice: ExerciceExt;
  allTags: string[];
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "back"): void;
  (e: "duplicate", exercice: ExerciceExt): void;
}>();

const isReadonly = computed(
  () => props.exercice.Origin.Visibility != Visibility.Personnal
);

function backToList() {
  emit("back");
}

let showEditDescription = $ref(false);

let query = $ref<Query>({ search: "", tags: [] });

async function save() {
  const res = await controller.ExerciceUpdate(props.exercice.Exercice);
  if (res == undefined) {
    return;
  }
  props.exercice.Exercice = res;
}

let showImportQuestion = $ref(false);

async function createQuestion() {
  const l = await controller.ExerciceCreateQuestion({
    IdExercice: props.exercice.Exercice.Id,
  });
  if (l == undefined) {
    return;
  }
  props.exercice.Questions = l;
}

async function addQuestion(idQuestion: number) {
  const current = (props.exercice.Questions || []).map((v) => v.Question);
  current.push({
    bareme: 1,
    id_question: idQuestion,
    id_exercice: -1,
  });
  const l = await controller.ExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: current,
  });
  if (l == undefined) {
    return;
  }
  props.exercice.Questions = l;
}

async function removeQuestion(index: number) {
  const l = (props.exercice.Questions || []).map((v) => v.Question);
  l.splice(index, 1);
  const res = await controller.ExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: l,
  });
  if (res == undefined) {
    return;
  }
  props.exercice.Questions = res;
}

let questionToEdit = $ref<ExerciceQuestion | null>(null);
let questionIndexToEdit = $ref<number | null>(null);
async function saveEditedQuestion() {
  if (questionIndexToEdit == null || questionToEdit == null) {
    return;
  }
  const current = (props.exercice.Questions || []).map((v) => v.Question);
  current[questionIndexToEdit] = questionToEdit;

  questionToEdit = null;
  questionIndexToEdit = null;
  const res = await controller.ExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: current,
  });
  if (res == undefined) {
    return;
  }
  props.exercice.Questions = res;
}

let showDropZone = $ref(false);

function onDragStart() {
  setTimeout(() => (showDropZone = true), 100); // workaround bug
}

function onDragEnd(ev: DragEvent) {
  showDropZone = false;
}

function onItemDragStart(paylod: DragEvent, index: number) {
  onDragListItemStart(paylod, index);
}

/** take the question at the index `origin` and insert it right before
the block at index `target` (which is between 0 and nbBlocks)
 */
async function swapQuestions(origin: number, target: number) {
  const l = swapItems(origin, target, props.exercice.Questions!).map(
    (v) => v.Question
  );
  const res = await controller.ExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: l,
  });
  if (res == undefined) {
    return;
  }
  props.exercice.Questions = res;
}
</script>

<style></style>
