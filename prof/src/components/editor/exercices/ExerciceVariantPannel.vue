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
      @go-to-question="goToQuestion"
      @update="(ex) => (exercice = ex)"
    ></ExSkeleton>
    <ExEditor
      v-else
      :session-id="props.sessionId"
      :question-index="questionIndex"
      :exercice="exercice"
      :is-readonly="props.isReadonly"
      @update="(ex) => (exercice = ex)"
    ></ExEditor>
  </div>
</template>

<script setup lang="ts">
import type { ExerciceExt, ExerciceHeader } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted, watch } from "vue";
import { $ref } from "vue/macros";
import ExEditor from "./ExEditor.vue";
import ExSkeleton from "./ExSkeleton.vue";

interface Props {
  exerciceHeader: ExerciceHeader;
  sessionId: string;
  isReadonly: boolean;
  allTags: string[]; // to provide auto completion
}

const props = defineProps<Props>();

let questionIndex = $ref(0);

let viewMode = $ref<"skeleton" | "question">("skeleton");

let exercice = $ref<ExerciceExt | null>(null);

onMounted(fetchExercice);

watch(props, (_) => {
  fetchExercice();
});

function goToQuestion(index: number) {
  questionIndex = index;
  viewMode = "question";
}

async function fetchExercice() {
  const res = await controller.EditorGetExerciceContent({
    id: props.exerciceHeader.Id,
  });
  if (res == undefined) return;
  exercice = res;
}
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
