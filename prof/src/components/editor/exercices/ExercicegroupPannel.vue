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
    <ExerciceVariantPannel
      :exercice-header="ownVariants[variantIndex]"
      :is-readonly="isReadonly"
      :all-tags="props.allTags"
      :show-variant-meta="true"
      @update="(ex) => (ownVariants[variantIndex] = ex)"
      @preview="(ex) => emit('preview', ex)"
    ></ExerciceVariantPannel>
  </ResourceScafold>
</template>

<script setup lang="ts">
import type {
  ExercicegroupExt,
  ExerciceHeader,
  LoopbackShowExercice,
  QuestionExerciceUses,
  Sheet,
  Tags,
  TagsDB,
} from "@/controller/api_gen";
import { Visibility } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import type { ResourceGroup, VariantG } from "@/controller/editor";
import { copy } from "@/controller/utils";
import { computed } from "vue";
import { useRouter } from "vue-router";
import { $computed, $ref } from "vue/macros";
import UsesCard from "../UsesCard.vue";
import ExerciceVariantPannel from "./ExerciceVariantPannel.vue";
import ResourceScafold from "../ResourceScafold.vue";

interface Props {
  group: ExercicegroupExt;
  allTags: TagsDB; // to provide auto completion
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
  (e: "preview", ex: LoopbackShowExercice): void;
}>();

const router = useRouter();

let group = $ref(copy(props.group));
let ownVariants = $ref(copy(props.group.Variants || []));

let variantIndex = $ref(0);

let isReadonly = $computed(
  () => props.group.Origin.Visibility != Visibility.Personnal
);

const resource = computed<ResourceGroup>(() => ({
  Id: group.Group.Id,
  Title: group.Group.Title,
  Tags: group.Tags,
  Variants: ownVariants,
}));

function updateTitle(t: string) {
  group.Group.Title = t;
  updateExercicegroup();
}

async function updateExercicegroup() {
  if (isReadonly) return;
  await controller.EditorUpdateExercicegroup(group.Group);

  // refresh the preview
  const res = await controller.EditorSaveExerciceAndPreview({
    OnlyPreview: true,
    IdExercice: ownVariants[variantIndex].Id,
    Parameters: [], // ignored
    Questions: [], // ignored
    CurrentQuestion: -1,
    ShowCorrection: false,
  });
  if (res == undefined) return;
  emit("preview", res.Preview);
}

let deletedBlocked = $ref<QuestionExerciceUses>(null);
function goToSheet(sh: Sheet) {
  deletedBlocked = null;

  router.push({ name: "homework", query: { idSheet: sh.Id } });
}

async function deleteVariante(variant: VariantG) {
  const res = await controller.EditorDeleteExercice({
    id: variant.Id,
  });
  if (res == undefined) return;
  // check if the variant is used
  if (!res.Deleted) {
    deletedBlocked = res.BlockedBy;
    return;
  }

  ownVariants = ownVariants.filter((qu) => qu.Id != variant.Id);
  if (ownVariants.length && variantIndex >= ownVariants.length) {
    variantIndex = 0;
  }
  // if there is no more variant, that means the exercicegroup is deleted:
  // go back
  if (!ownVariants.length) {
    emit("back");
  }
}

async function duplicateVariante(exercice: VariantG) {
  const newExercice = await controller.EditorDuplicateExercice({
    id: exercice.Id,
  });
  if (newExercice == undefined) {
    return;
  }
  ownVariants.push(newExercice);
  variantIndex = ownVariants.length - 1; // go to the new exercice
}

async function saveTags(newTags: Tags) {
  const rep = await controller.EditorUpdateExerciceTags({
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
  await controller.EditorSaveExerciceMeta(ownVariants[variantIndex]);
}
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
