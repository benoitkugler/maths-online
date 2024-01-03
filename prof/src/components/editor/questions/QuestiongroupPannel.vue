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
    ref="scafold"
  >
    <QuestionVariantPannel
      :question="ownVariants[variantIndex]"
      :readonly="isReadonly"
      @update="(qu: Question) => (ownVariants[variantIndex] = qu)"
      @preview="(qu: LoopbackShowQuestion) => emit('preview', qu)"
    ></QuestionVariantPannel>
  </ResourceScafold>
</template>

<script setup lang="ts">
import type {
  LoopbackShowQuestion,
  Question,
  QuestiongroupExt,
  Sheet,
  Tags,
  TagsDB,
  TaskUses,
} from "@/controller/api_gen";
import { Visibility } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import type { ResourceGroup, VariantG } from "@/controller/editor";
import { copy } from "@/controller/utils";
import { useRouter } from "vue-router";
import UsesCard from "../UsesCard.vue";
import QuestionVariantPannel from "./QuestionVariantPannel.vue";
import ResourceScafold from "../ResourceScafold.vue";
import { computed, ref } from "vue";

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

const group = ref(copy(props.group));
const ownVariants = ref(copy(props.variants));
const resource = computed<ResourceGroup>(() => ({
  Id: group.value.Group.Id,
  Title: group.value.Group.Title,
  Tags: group.value.Tags,
  Variants: ownVariants.value,
  Origin: group.value.Origin,
}));

const variantIndex = ref(0);

const isReadonly = computed(
  () => props.group.Origin.Visibility != Visibility.Personnal
);

function updateTitle(t: string) {
  group.value.Group.Title = t;
  updateQuestiongroup();
}

async function updateQuestiongroup() {
  if (isReadonly.value) return;
  await controller.EditorUpdateQuestiongroup(group.value.Group);
}

const deletedBlocked = ref<TaskUses>(null);
function goToSheet(sh: Sheet) {
  deletedBlocked.value = null;

  router.push({ name: "homework", query: { idSheet: sh.Id } });
}

async function deleteVariante(que: VariantG) {
  const res = await controller.EditorDeleteQuestion({ id: que.Id });
  if (res == undefined) return;
  // check if the variant is used
  if (!res.Deleted) {
    deletedBlocked.value = res.BlockedBy;
    return;
  }

  ownVariants.value = ownVariants.value.filter((qu) => qu.Id != que.Id);
  if (
    ownVariants.value.length &&
    variantIndex.value >= ownVariants.value.length
  ) {
    variantIndex.value = 0;
  }
  // if there is no more variant, that means the questiongroup is deleted:
  // go back
  if (!ownVariants.value.length) {
    emit("back");
  }
}

const scafold = ref<InstanceType<typeof ResourceScafold> | null>(null);

async function duplicateVariante(variant: VariantG) {
  const newQuestion = await controller.EditorDuplicateQuestion({
    id: variant.Id,
  });
  if (newQuestion == undefined) {
    return;
  }
  ownVariants.value.push(newQuestion);
  variantIndex.value = ownVariants.value.length - 1; // go to the new question

  if (scafold.value) scafold.value.showEditVariant(newQuestion);
}

async function saveTags(newTags: Tags) {
  const rep = await controller.EditorUpdateQuestionTags({
    Id: group.value.Group.Id,
    Tags: newTags,
  });
  if (rep === undefined) {
    return;
  }
  group.value.Tags = newTags;
}

function backToList() {
  emit("back");
}

async function updateVariant(variant: VariantG) {
  if (isReadonly.value) {
    return;
  }
  ownVariants.value[variantIndex.value].Subtitle = variant.Subtitle;
  ownVariants.value[variantIndex.value].Difficulty = variant.Difficulty;
  await controller.EditorSaveQuestionMeta({
    Question: ownVariants.value[variantIndex.value],
  });
}
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
