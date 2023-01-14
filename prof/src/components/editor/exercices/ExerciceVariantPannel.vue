<template>
  <div v-if="exercice == null">Chargement de l'exercice...</div>
  <div v-else>
    <!-- Display either the skeleton or the question editor -->
    <ExSkeleton
      v-if="viewMode == 'skeleton'"
      :all-tags="props.allTags"
      :exercice="exercice"
      :is-readonly="props.isReadonly"
      :show-variant-meta="props.showVariantMeta"
      @go-to-question="goToQuestion"
      @update="notifieUpdate"
      @preview="(qu) => emit('preview', qu)"
    ></ExSkeleton>
    <ExEditor
      v-else
      :question-index="questionIndex"
      :exercice="exercice"
      :is-readonly="props.isReadonly"
      @update="notifieUpdate"
      @preview="(qu) => emit('preview', qu)"
      @back="questionIndex = -1"
    ></ExEditor>
  </div>
</template>

<script setup lang="ts">
import {
  DifficultyTag,
  type ExerciceExt,
  type ExerciceHeader,
  type IdExercice,
  type LoopbackShowExercice,
  type TagsDB,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { computed } from "@vue/reactivity";
import { onMounted, watch } from "vue";
import { $ref } from "vue/macros";
import ExEditor from "./ExEditor.vue";
import ExSkeleton from "./ExSkeleton.vue";

interface Props {
  exerciceHeader: ExerciceHeader;
  isReadonly: boolean;
  allTags: TagsDB; // to provide auto completion
  showVariantMeta: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", ex: ExerciceHeader): void;
  (e: "preview", ex: LoopbackShowExercice): void;
}>();

let questionIndex = $ref(-1);

const viewMode = computed(() => (questionIndex >= 0 ? "question" : "skeleton"));

let exercice = $ref<ExerciceExt | null>(null);

onMounted(() => {
  fetchExercice();
  refreshExercicePreview(props.exerciceHeader.Id);
});

watch(props, (_) => {
  fetchExercice();
  refreshExercicePreview(props.exerciceHeader.Id);
});

async function refreshExercicePreview(id: IdExercice) {
  const res = await controller.EditorSaveExerciceAndPreview({
    OnlyPreview: true,
    IdExercice: id,
    Parameters: { Intrinsics: [], Variables: [] }, // ignored
    Questions: [], // ignored
  });
  if (res == undefined) return;
  emit("preview", res.Preview);
}

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
