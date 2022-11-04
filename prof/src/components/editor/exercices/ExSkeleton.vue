<template>
  <!-- description for the whole exercice, not a particular question -->
  <v-dialog v-model="showEditDescription" max-width="600px">
    <description-pannel
      :description="props.exercice.Exercice.Description"
      :readonly="props.isReadonly"
      @save="saveDescription"
    ></description-pannel>
  </v-dialog>

  <v-dialog v-model="showImportQuestion" :retain-focus="false">
    <keep-alive>
      <question-selector
        :tags="props.allTags"
        :query="questionQuery"
        @closed="showImportQuestion = false"
        @selected="
          (q) => {
            showImportQuestion = false;
            importQuestion(q.Id);
          }
        "
      ></question-selector>
    </keep-alive>
  </v-dialog>

  <v-card class="px-2">
    <v-row no-gutters class="mb-2">
      <template v-if="props.showVariantMeta">
        <v-col>
          <v-text-field
            class="my-2 input-small"
            variant="outlined"
            density="compact"
            label="Sous-titre de la variante (optionnel)"
            v-model="props.exercice.Exercice.Subtitle"
            :readonly="props.isReadonly"
            hide-details
            @blur="saveMeta"
          ></v-text-field
        ></v-col>
      </template>
      <v-spacer v-else></v-spacer>

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

      <v-col cols="auto" align-self="center">
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
                class="my-1"
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

    <v-list @dragstart="onDragStart" @dragend="onDragEnd" style="height: 67vh">
      <drop-zone
        v-if="showDropZone"
        @drop="(origin) => swapQuestions(origin, 0)"
      ></drop-zone>

      <div v-for="(question, index) in props.exercice.Questions" :key="index">
        <v-list-item link @click="emit('goToQuestion', index)">
          <v-row>
            <v-col cols="auto" align-self="center">
              <drag-icon
                color="black"
                @start="(e) => onItemDragStart(e, index)"
              ></drag-icon>
            </v-col>
            <v-col cols="auto" align-self="center" class="my-1">
              <v-btn
                class="mx-2 my-1"
                size="x-small"
                icon
                @click.stop="duplicateQuestion(index)"
                title="Dupliquer cette question"
              >
                <v-icon icon="mdi-content-copy" color="secondary"></v-icon>
              </v-btn>
              <v-btn
                v-if="!isReadonly"
                size="x-small"
                icon
                @click.stop="removeQuestion(index)"
                title="Supprimer cette question"
              >
                <v-icon icon="mdi-delete" color="red"></v-icon>
              </v-btn>
            </v-col>
            <v-col align-self="center">
              Question {{ index + 1 }}
              <small>(Id : {{ question.Question.Id }})</small>
            </v-col>
            <v-col cols="2" align-self="center">
              <v-menu
                offset-y
                :close-on-content-click="false"
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
                    >/ {{ question.Bareme }}
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
                      v-model.number="questionToEdit.Bareme"
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
import DropZone from "@/components/DropZone.vue";
import type { ExerciceExt, ExerciceQuestionExt } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { copy, onDragListItemStart, swapItems } from "@/controller/utils";
import { $ref } from "vue/macros";
import DragIcon from "../../DragIcon.vue";
import QuestionSelector, {
  type QuestionQuery,
} from "../../QuestionSelector.vue";
import DescriptionPannel from "../DescriptionPannel.vue";

interface Props {
  sessionId: string;
  exercice: ExerciceExt;
  isReadonly: boolean;
  allTags: string[]; // used to select the question to import
  showVariantMeta: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "update", exercice: ExerciceExt): void;
  (e: "goToQuestion", questionIndex: number): void;
}>();

let questionQuery = $ref<QuestionQuery>({ search: "", tags: [] });
let showImportQuestion = $ref(false);

let showEditDescription = $ref(false);
async function saveDescription(description: string) {
  showEditDescription = false;
  props.exercice.Exercice.Description = description;
  saveMeta();
}

async function saveMeta() {
  if (props.isReadonly) {
    return;
  }
  await controller.EditorSaveExerciceMeta(props.exercice.Exercice);
  emit("update", props.exercice);
}

async function createQuestion() {
  const res = await controller.EditorExerciceCreateQuestion({
    IdExercice: props.exercice.Exercice.Id,
    SessionID: props.sessionId,
  });
  if (res == undefined) {
    return;
  }
  emit("update", res);
}

function toExerciceQuestions(questions: ExerciceQuestionExt[]) {
  return questions.map((qu) => ({
    id_exercice: -1,
    id_question: qu.Question.Id,
    bareme: qu.Bareme,
  }));
}

async function importQuestion(idQuestion: number) {
  const res = await controller.EditorExerciceImportQuestion({
    IdExercice: props.exercice.Exercice.Id,
    IdQuestion: idQuestion,
    SessionID: props.sessionId,
  });
  if (res == undefined) {
    return;
  }
  emit("update", res);
}

async function removeQuestion(index: number) {
  const l = toExerciceQuestions(props.exercice.Questions || []);
  l.splice(index, 1);
  const res = await controller.EditorExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: l,
    SessionID: props.sessionId,
  });
  if (res == undefined) return;
  emit("update", res);
}

async function duplicateQuestion(index: number) {
  const res = await controller.EditorExerciceDuplicateQuestion({
    IdExercice: props.exercice.Exercice.Id,
    QuestionIndex: index,
    SessionID: props.sessionId,
  });
  if (res == undefined) return;
  emit("update", res);
}

let questionToEdit = $ref<ExerciceQuestionExt | null>(null);
let questionIndexToEdit = $ref<number | null>(null);
async function saveEditedQuestion() {
  if (questionIndexToEdit == null || questionToEdit == null) {
    return;
  }
  const current = copy(props.exercice.Questions || []);
  current[questionIndexToEdit] = questionToEdit;

  questionToEdit = null;
  questionIndexToEdit = null;
  const res = await controller.EditorExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: toExerciceQuestions(current),
    SessionID: props.sessionId,
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
  const res = await controller.EditorExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: toExerciceQuestions(l),
    SessionID: props.sessionId,
  });
  if (res == undefined) {
    return;
  }
  emit("update", res);
}
</script>

<style scoped>
:deep(.v-input__append) {
  padding-top: 0;
  align-self: center;
  margin-inline-start: 4px;
}
</style>
