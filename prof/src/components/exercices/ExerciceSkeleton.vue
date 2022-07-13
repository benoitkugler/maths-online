<template>
  <v-dialog v-model="showEditDescription">
    <description-pannel
      v-model="props.exercice.Exercice.Description"
      :readonly="isReadonly"
    ></description-pannel>
  </v-dialog>

  <v-dialog v-model="showFlowDocumentation">
    <v-card title="Déroulement d'un exercice">
      <v-card-text>
        Un exercice peut se dérouler de deux manières différentes :
        séquentiellement ou en parallèle.
        <v-list>
          <v-list-item>
            <v-row>
              <v-col cols="3"> <b>Déroulé séquentiel</b></v-col>
              <v-col>
                L'élève répond aux questions l'une après l'autre, et ne peux pas
                accéder à la question suivante tant qu'il n'a pas réussi la
                question courante.</v-col
              >
            </v-row>
          </v-list-item>
          <v-list-item>
            <v-row>
              <v-col cols="3">
                <b>Déroulé parallèle</b>
              </v-col>
              <v-col>
                L'élève répond à toutes les questions avant de valider, dans
                l'ordre de son choix, puis peut reprendre les questions
                fausses.</v-col
              >
            </v-row>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
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
          @click="emit('back')"
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
            <v-select
              density="compact"
              variant="outlined"
              :items="flowItems.map((v) => v.text)"
              label="Déroulement"
              hide-details
              :model-value="
                flowItems.find((v) => v.value == props.exercice.Exercice.Flow)
                  ?.text
              "
              @update:model-value="v => props.exercice.Exercice.Flow = flowItems.find((item) => item.text == v)!.value"
            >
              <template v-slot:append>
                <v-btn
                  icon
                  size="x-small"
                  @click="showFlowDocumentation = true"
                >
                  <v-icon icon="mdi-help" color="info" size="x-small"></v-icon>
                </v-btn>
              </template>
            </v-select>
          </v-col>
          <v-spacer></v-spacer>
          <v-col cols="auto" align-self="center" class="pl-2">
            <v-btn
              title="Créer et ajouter une question"
              size="small"
              @click="createQuestion()"
            >
              <v-icon icon="mdi-plus" color="green"></v-icon>
              Créer une question
            </v-btn>
            <v-btn
              title="Importer une question existante"
              size="small"
              class="mx-2"
              @click="showImportQuestion = true"
            >
              <v-icon icon="mdi-plus" color="green"></v-icon>
              Importer
            </v-btn>
          </v-col>
        </v-row>
      </v-col>

      <v-col cols="auto" align-self="center" class="pl-2">
        <v-btn
          size="small"
          icon
          title="Editer le contenu des questions"
          @click="emit('next')"
        >
          <v-icon icon="mdi-arrow-right"></v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-list @dragstart="onDragStart" @dragend="onDragEnd" style="height: 66vh">
      <drop-zone
        v-if="showDropZone"
        @drop="(origin) => swapQuestions(origin, 0)"
      ></drop-zone>

      <div v-for="(question, index) in props.exercice.Questions" :key="index">
        <v-list-item>
          <v-row>
            <v-col cols="auto" align-self="center">
              <v-icon
                style="cursor: grab"
                draggable="true"
                @dragstart="(e) => onItemDragStart(e, index)"
                size="large"
                icon="mdi-drag-vertical"
              ></v-icon>
            </v-col>
            <v-col cols="auto" align-self="center">
              <v-btn
                v-if="!isReadonly"
                size="x-small"
                icon
                @click.stop="removeQuestion(index)"
                title="Retirer la question"
              >
                <v-icon icon="mdi-delete" color="red" size="small"></v-icon>
              </v-btn>
              <v-btn
                v-if="!isReadonly"
                class="mx-1"
                size="x-small"
                icon
                @click.stop="duplicateQuestion(index)"
                title="Dupliquer la question"
              >
                <v-icon
                  icon="mdi-content-copy"
                  color="info"
                  size="small"
                ></v-icon>
              </v-btn>
            </v-col>
            <v-col align-self="center">
              <small>({{ question.id_question }})</small>
              {{ getQuestion(question.id_question).Question.page.title }}</v-col
            >
            <v-col cols="auto">
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
                    elevation="3"
                    v-on="{ isActive }"
                    v-bind="props"
                    color="primary"
                    @click="
                      questionToEdit = copy(question);
                      questionIndexToEdit = index;
                    "
                    :disabled="isReadonly"
                    >/ {{ question.bareme }}
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
import type { Flow } from "@/controller/api_gen";
import {
  FlowLabels,
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
  session_id: string;
  allTags: string[];
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "back"): void;
  (e: "next"): void;
  (e: "update", exercice: ExerciceExt): void;
  (e: "duplicate", exercice: ExerciceExt): void;
}>();

const isReadonly = computed(
  () => props.exercice.Origin.Visibility != Visibility.Personnal
);

const flowItems = Object.entries(FlowLabels).map((k) => ({
  value: Number(k[0]) as Flow,
  text: k[1],
}));

let showEditDescription = $ref(false);

function getQuestion(questionID: number) {
  console.log(
    questionID,
    props.exercice.QuestionsSource,
    props.exercice.QuestionsSource![questionID]
  );

  return props.exercice.QuestionsSource![questionID];
}

let query = $ref<Query>({ search: "", tags: [] });

async function save() {
  const res = await controller.ExerciceUpdate({
    Exercice: props.exercice.Exercice,
    SessionID: props.session_id,
  });
  if (res == undefined) {
    return;
  }
  props.exercice.Exercice = res;
}

let showImportQuestion = $ref(false);

async function createQuestion() {
  const res = await controller.ExerciceCreateQuestion({
    IdExercice: props.exercice.Exercice.Id,
    SessionID: props.session_id,
  });
  if (res == undefined) {
    return;
  }
  emit("update", res);
}

async function addQuestion(idQuestion: number) {
  const current = copy(props.exercice.Questions || []);
  current.push({
    bareme: 1,
    id_question: idQuestion,
    id_exercice: -1,
  });
  const res = await controller.ExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: current,
    SessionID: props.session_id,
  });
  if (res == undefined) {
    return;
  }
  emit("update", res);
}

async function removeQuestion(index: number) {
  const l = copy(props.exercice.Questions || []);
  l.splice(index, 1);
  const res = await controller.ExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: l,
    SessionID: props.session_id,
  });
  if (res == undefined) {
    return;
  }
  emit("update", res);
}

async function duplicateQuestion(index: number) {
  const l = copy(props.exercice.Questions || []);
  const added = l.slice(0, index).concat(l[index]).concat(l.slice(index));
  const res = await controller.ExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: added,
    SessionID: props.session_id,
  });
  if (res == undefined) {
    return;
  }
  emit("update", res);
}

let questionToEdit = $ref<ExerciceQuestion | null>(null);
let questionIndexToEdit = $ref<number | null>(null);
async function saveEditedQuestion() {
  if (questionIndexToEdit == null || questionToEdit == null) {
    return;
  }
  const current = copy(props.exercice.Questions || []);
  current[questionIndexToEdit] = questionToEdit;

  questionToEdit = null;
  questionIndexToEdit = null;
  const res = await controller.ExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: current,
    SessionID: props.session_id,
  });
  if (res == undefined) {
    return;
  }
  emit("update", res);
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
  const l = swapItems(origin, target, props.exercice.Questions!);
  const res = await controller.ExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: l,
    SessionID: props.session_id,
  });
  if (res == undefined) {
    return;
  }
  emit("update", res);
}

let showFlowDocumentation = $ref(false);
</script>

<style scoped>
:deep(.v-input__append) {
  padding-top: 0;
  align-self: center;
  margin-inline-start: 4px;
}
</style>
