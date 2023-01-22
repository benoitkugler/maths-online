<template>
  <v-dialog
    :model-value="deletedBlocked != null"
    @update:model-value="deletedBlocked = null"
    max-width="600"
  >
    <UsesCard
      :uses="deletedBlocked || []"
      @back="deletedBlocked = null"
      @go-to="goToSheet"
    ></UsesCard>
  </v-dialog>

  <ResourceScafold
    :resource="resource"
    :readonly="isReadonly"
    :all-tags="props.allTags"
    v-model="variantIndex"
    @back="backToList"
    @update-title="updateTitle"
    @update-tags="saveTags"
    @update-variant="updateVariant"
    @duplicate-variant="duplicateVariante"
    @delete-variant="deleteVariante"
  >
    <QuestionVariantPannel
      :question="ownVariants[variantIndex]"
      :readonly="isReadonly"
      @update="(qu) => (ownVariants[variantIndex] = qu)"
      @preview="(qu) => emit('preview', qu)"
    ></QuestionVariantPannel>
  </ResourceScafold>
</template>

<script setup lang="ts">
import type {
  LoopbackShowQuestion,
  Question,
  QuestionExerciceUses,
  QuestiongroupExt,
  Sheet,
  Tags,
  TagsDB,
} from "@/controller/api_gen";
import { Visibility } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import type { ResourceGroup, VariantG } from "@/controller/editor";
import { copy } from "@/controller/utils";
import { useRouter } from "vue-router";
import { $computed, $ref } from "vue/macros";
import UsesCard from "../UsesCard.vue";
import QuestionVariantPannel from "./QuestionVariantPannel.vue";
import ResourceScafold from "../ResourceScafold.vue";
import { computed } from "vue";

interface Props {
  group: QuestiongroupExt;
  variants: Question[];
  allTags: TagsDB; // to provide auto completion
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
  (e: "preview", question: LoopbackShowQuestion): void;
}>();

const router = useRouter();

let group = $ref(copy(props.group));
let ownVariants = $ref(copy(props.variants));
let resource = computed<ResourceGroup>(() => ({
  Id: group.Group.Id,
  Title: group.Group.Title,
  Tags: group.Tags,
  Variants: ownVariants,
}));

let variantIndex = $ref(0);

let isReadonly = $computed(
  () => props.group.Origin.Visibility != Visibility.Personnal
);

function updateTitle(t: string) {
  group.Group.Title = t;
  updateQuestiongroup();
}

async function updateQuestiongroup() {
  if (isReadonly) return;
  await controller.EditorUpdateQuestiongroup(group.Group);
}

let deletedBlocked = $ref<QuestionExerciceUses>(null);
function goToSheet(sh: Sheet) {
  deletedBlocked = null;

  router.push({ name: "homework", query: { idSheet: sh.Id } });
}

async function deleteVariante(que: VariantG) {
  const res = await controller.EditorDeleteQuestion({ id: que.Id });
  if (res == undefined) return;
  // check if the question is used
  if (!res.Deleted) {
    deletedBlocked = res.BlockedBy;
    return;
  }

  ownVariants = ownVariants.filter((qu) => qu.Id != que.Id);
  if (ownVariants.length && variantIndex >= ownVariants.length) {
    variantIndex = 0;
  }
  // if there is no more variant, that means the questiongroup is deleted:
  // go back
  if (!ownVariants.length) {
    emit("back");
  }
}

async function duplicateVariante(variant: VariantG) {
  const newQuestion = await controller.EditorDuplicateQuestion({
    id: variant.Id,
  });
  if (newQuestion == undefined) {
    return;
  }
  ownVariants.push(newQuestion);
  variantIndex = ownVariants.length - 1; // go to the new question
}

async function saveTags(newTags: Tags) {
  const rep = await controller.EditorUpdateQuestionTags({
    Id: group.Group.Id,
    Tags: newTags,
  });
  if (rep === undefined) {
    return;
  }
  group.Tags = newTags;
}

function backToList() {
  emit("back");
}

async function updateVariant(variant: VariantG) {
  if (isReadonly) {
    return;
  }
  ownVariants[variantIndex].Subtitle = variant.Subtitle;
  ownVariants[variantIndex].Difficulty = variant.Difficulty;
  await controller.EditorSaveQuestionMeta({
    Question: ownVariants[variantIndex],
  });
}
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
