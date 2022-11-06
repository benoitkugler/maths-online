<template>
  <v-dialog
    :model-value="exerciceToDelete != null"
    @update:model-value="exerciceToDelete = null"
    max-width="700"
  >
    <v-card title="Confirmer">
      <v-card-text
        >Etes-vous certain de vouloir supprimer la variante
        <i>{{ exerciceToDelete?.Id }} - {{ exerciceToDelete?.Subtitle }}</i> ?
        <br />
        Cette opération est irréversible.

        <div v-if="ownVariants.length == 1" class="mt-2">
          L'exercice associé sera aussi supprimé.
        </div>
      </v-card-text>
      <v-card-actions>
        <v-btn @click="exerciceToDelete = null">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn color="red" @click="deleteVariante" variant="outlined">
          Supprimer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

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

  <v-card class="mt-3 px-2">
    <v-row no-gutters>
      <v-col cols="auto" align-self="center" class="pr-2">
        <v-btn
          size="small"
          icon
          title="Retour aux exercices"
          @click="backToList"
        >
          <v-icon icon="mdi-arrow-left"></v-icon>
        </v-btn>
      </v-col>

      <v-col>
        <v-text-field
          class="my-2 input-small"
          variant="outlined"
          density="comfortable"
          label="Nom de l'exercice"
          v-model="group.Group.Title"
          :readonly="isReadonly"
          @blur="updateExercicegroup"
          hide-details
        ></v-text-field>
      </v-col>

      <v-col cols="3" class="px-1" align-self="center">
        <TagListField
          label="Catégories"
          :all-tags="props.allTags"
          :model-value="group.Tags || []"
          @update:model-value="saveTags"
          :readonly="isReadonly"
        ></TagListField
      ></v-col>

      <v-col cols="auto" align-self="center" class="px-1">
        <VariantsSelector
          :variants="ownVariants"
          :readonly="isReadonly"
          v-model="variantIndex"
          @delete="(qu) => (exerciceToDelete = qu)"
          @duplicate="duplicateVariante"
        ></VariantsSelector>
      </v-col>
    </v-row>

    <ExerciceVariantPannel
      :exercice-header="ownVariants[variantIndex]"
      :is-readonly="isReadonly"
      :session-id="props.session_id"
      :all-tags="props.allTags"
      :show-variant-meta="true"
      @update="(qu) => (ownVariants[variantIndex] = qu)"
    ></ExerciceVariantPannel>
  </v-card>
</template>

<script setup lang="ts">
import type {
  ExercicegroupExt,
  ExerciceHeader,
  QuestionExerciceUses,
  Sheet,
} from "@/controller/api_gen";
import { Visibility } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import type { VariantG } from "@/controller/editor";
import { copy } from "@/controller/utils";
import { useRouter } from "vue-router";
import { $computed, $ref } from "vue/macros";
import TagListField from "../TagListField.vue";
import UsesCard from "../UsesCard.vue";
import VariantsSelector from "../VariantsSelector.vue";
import ExerciceVariantPannel from "./ExerciceVariantPannel.vue";

interface Props {
  session_id: string;
  group: ExercicegroupExt;
  allTags: string[]; // to provide auto completion
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
}>();

const router = useRouter();

let group = $ref(copy(props.group));
let ownVariants = $ref(copy(props.group.Variants || []));

let variantIndex = $ref(0);

let isReadonly = $computed(
  () => props.group.Origin.Visibility != Visibility.Personnal
);

async function updateExercicegroup() {
  if (isReadonly) return;
  await controller.EditorUpdateExercicegroup(group.Group);
}

let deletedBlocked = $ref<QuestionExerciceUses>(null);
function goToSheet(sh: Sheet) {
  deletedBlocked = null;

  router.push({ name: "homework", query: { idSheet: sh.Id } });
}

let exerciceToDelete: ExerciceHeader | null = $ref(null);
async function deleteVariante() {
  const exe = exerciceToDelete;
  if (exe == null) return;
  const res = await controller.EditorDeleteExercice({
    id: exe.Id,
  });
  exerciceToDelete = null;

  if (res == undefined) return;
  // check if the question is used
  if (!res.Deleted) {
    deletedBlocked = res.BlockedBy;
    return;
  }

  ownVariants = ownVariants.filter((qu) => qu.Id != exe.Id);
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

async function saveTags(newTags: string[]) {
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
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
