<template>
  <div v-if="exercice == null">Loading {{ props.exerciceHeader }}</div>
  <div v-else>
    <!-- Display either the skeleton or the question editor -->
    <ExSkeleton
      v-if="viewMode == 'skeleton'"
      :session-id="props.sessionId"
      :all-tags="props.allTags"
      :exercice="exercice"
      :is-readonly="props.isReadonly"
      :show-variant-meta="props.showVariantMeta"
      @go-to-question="goToQuestion"
      @update="notifieUpdate"
    ></ExSkeleton>
    <ExEditor
      v-else
      :session-id="props.sessionId"
      :question-index="questionIndex"
      :exercice="exercice"
      :is-readonly="props.isReadonly"
      @update="notifieUpdate"
      @back="questionIndex = -1"
    ></ExEditor>
  </div>
</template>

<script setup lang="ts">
import {
  DifficultyTag,
  type ExerciceExt,
  type ExerciceHeader,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { refreshExercicePreview } from "@/controller/editor";
import { computed } from "@vue/reactivity";
import { onMounted, watch } from "vue";
import { $ref } from "vue/macros";
import ExEditor from "./ExEditor.vue";
import ExSkeleton from "./ExSkeleton.vue";

interface Props {
  exerciceHeader: ExerciceHeader;
  sessionId: string;
  isReadonly: boolean;
  allTags: string[]; // to provide auto completion
  showVariantMeta: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", ex: ExerciceHeader): void;
}>();

let questionIndex = $ref(-1);

const viewMode = computed(() => (questionIndex >= 0 ? "question" : "skeleton"));

let exercice = $ref<ExerciceExt | null>(null);

onMounted(() => {
  fetchExercice();
  refreshExercicePreview(props.sessionId, props.exerciceHeader.Id);
});

watch(props, (_) => {
  fetchExercice();
  refreshExercicePreview(props.sessionId, props.exerciceHeader.Id);
});

function goToQuestion(index: number) {
  questionIndex = index;
}

async function fetchExercice() {
  const res = await controller.EditorGetExerciceContent({
    id: props.exerciceHeader.Id,
  });
  if (res == undefined) return;
  exercice = res;
}

function notifieUpdate(ex: ExerciceExt) {
  exercice = ex;
  emit("update", {
    Id: ex.Exercice.Id,
    Difficulty: DifficultyTag.DiffEmpty, // for now, we dont support difficulty on exercice
    Subtitle: ex.Exercice.Subtitle,
  });
}
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
