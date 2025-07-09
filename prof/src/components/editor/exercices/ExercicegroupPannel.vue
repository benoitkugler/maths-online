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
    <ExerciceVariantPannel
      :exercice-header="ownVariants[variantIndex]"
      :readonly="isReadonly"
      :all-tags="props.allTags"
      :show-variant-meta="true"
      @update="(ex) => (ownVariants[variantIndex] = ex)"
      ref="editor"
    ></ExerciceVariantPannel>
  </ResourceScafold>
</template>

<script setup lang="ts">
import type {
  ExercicegroupExt,
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
import ExerciceVariantPannel from "./ExerciceVariantPannel.vue";
import ResourceScafold from "../ResourceScafold.vue";
import { ref } from "vue";
import { computed } from "vue";

interface Props {
  group: ExercicegroupExt;
  allTags: TagsDB; // to provide auto completion
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
}>();

const router = useRouter();

const group = ref(copy(props.group));
const ownVariants = ref(copy(props.group.Variants || []));

const variantIndex = ref(0);

const isReadonly = computed(
  () => props.group.Origin.Visibility != Visibility.Personnal
);

const resource = computed<ResourceGroup>(() => ({
  Id: group.value.Group.Id,
  Title: group.value.Group.Title,
  Tags: group.value.Tags,
  Variants: ownVariants.value,
  Origin: group.value.Origin,
}));

function updateTitle(t: string) {
  group.value.Group.Title = t;
  updateExercicegroup();
}

const editor = ref<InstanceType<typeof ExerciceVariantPannel> | null>(null);

async function updateExercicegroup() {
  if (isReadonly.value) return;
  const res = await controller.EditorUpdateExercicegroup(group.value.Group);
  if (res === undefined) return;
  controller.showMessage("Exercice modifié avec succès.");

  // refresh the preview
  editor.value?.refreshExercicePreview();
}

const deletedBlocked = ref<TaskUses>(null);
function goToSheet(sh: Sheet) {
  deletedBlocked.value = null;

  router.push({ name: "homework", query: { idSheet: sh.Id } });
}

async function deleteVariante(variant: VariantG) {
  const res = await controller.EditorDeleteExercice({
    id: variant.Id,
  });
  if (res == undefined) return;
  // check if the variant is used
  if (!res.Deleted) {
    deletedBlocked.value = res.BlockedBy;
    return;
  } else {
    controller.showMessage("Exercice supprimé avec succès.");
  }

  ownVariants.value = ownVariants.value.filter((qu) => qu.Id != variant.Id);
  if (
    ownVariants.value.length &&
    variantIndex.value >= ownVariants.value.length
  ) {
    variantIndex.value = 0;
  }
  // if there is no more variant, that means the exercicegroup is deleted:
  // go back
  if (!ownVariants.value.length) {
    emit("back");
  }
}

const scafold = ref<InstanceType<typeof ResourceScafold> | null>(null);

async function duplicateVariante(exercice: VariantG) {
  const newExercice = await controller.EditorDuplicateExercice({
    id: exercice.Id,
  });
  if (newExercice == undefined) {
    return;
  }
  controller.showMessage("Variante dupliquée avec succès.");

  ownVariants.value.push(newExercice);
  variantIndex.value = ownVariants.value.length - 1; // go to the new exercice

  if (scafold.value) scafold.value.showEditVariant(newExercice);
}

async function saveTags(newTags: Tags) {
  const rep = await controller.EditorUpdateExerciceTags({
    Id: group.value.Group.Id,
    Tags: newTags,
  });
  if (rep === undefined) {
    return;
  }
  controller.showMessage("Etiquettes modifiées avec succès.");

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
  const res = await controller.EditorSaveExerciceMeta(
    ownVariants.value[variantIndex.value]
  );
  if (res === undefined) return;
  controller.showMessage("Exercice modifié avec succès.");
}
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
