<template>
  <v-card class="px-2" title="Questions et barême">
    <v-list @dragstart="onDragStart" @dragend="onDragEnd">
      <drop-zone
        v-if="showDropZone"
        @drop="(origin) => swapQuestions(origin, 0)"
      ></drop-zone>

      <div v-for="(question, index) in props.exercice.Questions" :key="index">
        <v-list-item>
          <v-row>
            <v-col cols="auto" align-self="center" v-if="!props.isReadonly">
              <drag-icon
                color="black"
                @start="(e) => onItemDragStart(e, index)"
              ></drag-icon>
            </v-col>
            <v-col cols="auto" align-self="center" class="my-1">
              <v-btn
                v-if="!props.isReadonly"
                class="mx-2 my-1"
                size="x-small"
                icon
                @click.stop="duplicateQuestion(index as Int)"
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
                    class="my-1"
                    elevation="3"
                    v-on="{ isActive }"
                    v-bind="props"
                    color="primary-darken"
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
import type {
  ExerciceExt,
  ExerciceQuestionExt,
  IdExercice,
  Int,
  LoopbackShowExercice,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { copy, onDragListItemStart, swapItems } from "@/controller/utils";
import DragIcon from "../../DragIcon.vue";
import { ref } from "vue";

interface Props {
  exercice: ExerciceExt;
  isReadonly: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "update", exercice: ExerciceExt): void;
  (e: "preview", exercice: LoopbackShowExercice): void;
}>();

function toExerciceQuestions(questions: ExerciceQuestionExt[]) {
  return questions.map((qu) => ({
    id_exercice: -1 as IdExercice,
    id_question: qu.Question.Id,
    bareme: qu.Bareme,
  }));
}

async function removeQuestion(index: number) {
  const l = toExerciceQuestions(props.exercice.Questions || []);
  l.splice(index, 1);
  const res = await controller.EditorExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: l,
  });
  if (res == undefined) return;
  controller.showMessage("Question supprimée avec succès.");

  emit("update", res.Ex);
  emit("preview", res.Preview);
}

async function duplicateQuestion(index: Int) {
  const res = await controller.EditorExerciceDuplicateQuestion({
    IdExercice: props.exercice.Exercice.Id,
    QuestionIndex: index,
  });
  if (res == undefined) return;
  controller.showMessage("Question dupliquée avec succès.");

  emit("update", res.Ex);
  emit("preview", res.Preview);
}

const questionToEdit = ref<ExerciceQuestionExt | null>(null);
const questionIndexToEdit = ref<number | null>(null);
async function saveEditedQuestion() {
  if (questionIndexToEdit.value == null || questionToEdit.value == null) {
    return;
  }
  const current = copy(props.exercice.Questions || []);
  current[questionIndexToEdit.value] = questionToEdit.value;

  questionToEdit.value = null;
  questionIndexToEdit.value = null;
  const res = await controller.EditorExerciceUpdateQuestions({
    IdExercice: props.exercice.Exercice.Id,
    Questions: toExerciceQuestions(current),
  });
  if (res == undefined) {
    return;
  }
  controller.showMessage("Questions modifiées avec succès.");

  emit("update", res.Ex);
  emit("preview", res.Preview);
}

const showDropZone = ref(false);

function onDragStart() {
  setTimeout(() => (showDropZone.value = true), 100); // workaround bug
}

function onDragEnd() {
  showDropZone.value = false;
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
  });
  if (res == undefined) {
    return;
  }
  controller.showMessage("Questions modifiées avec succès.");

  emit("update", res.Ex);
  emit("preview", res.Preview);
}
</script>

<style scoped>
:deep(.v-input__append) {
  padding-top: 0;
  align-self: center;
  margin-inline-start: 4px;
}
</style>
