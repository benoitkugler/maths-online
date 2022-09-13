<template>
  <v-dialog
    :model-value="exerciceToDelete != null"
    @update:model-value="exerciceToDelete = null"
  >
    <v-card title="Confirmer">
      <v-card-text
        >Etes-vous certain de vouloir supprimer la variante
        <i>{{ exerciceToDelete?.Id }} - {{ exerciceToDelete?.Subtitle }}</i> ?
        <br />
        Cette opération est irréversible.

        <div v-if="variants.length == 1" class="mt-2">
          Le groupe d'exercices associé sera aussi supprimé.
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

  <v-card class="mt-3 px-2">
    <v-row no-gutters class="mb-2">
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
        <ExerciceVariantsSelector
          :variants="variants"
          :readonly="isReadonly"
          v-model="variantIndex"
          @delete="(qu) => (exerciceToDelete = qu)"
          @duplicate="duplicateVariante"
        ></ExerciceVariantsSelector>
      </v-col>
    </v-row>

    TODO
    <!-- <VariantPannel
      :exercice="variants[variantIndex]"
      :readonly="isReadonly"
      :session_id="props.session_id"
      :all-tags="props.allTags"
      @update="(qu) => (variants[variantIndex] = qu)"
    ></VariantPannel> -->
  </v-card>
</template>

<script setup lang="ts">
import type { ExercicegroupExt, ExerciceHeader } from "@/controller/api_gen";
import { Visibility } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { copy } from "@/controller/utils";
import { $computed, $ref } from "vue/macros";
import TagListField from "../TagListField.vue";
import ExerciceVariantsSelector from "./ExerciceVariantsSelector.vue";

interface Props {
  session_id: string;
  group: ExercicegroupExt;
  allTags: string[]; // to provide auto completion
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
}>();

let group = $ref(copy(props.group));
let variants = $ref(copy(props.group.Variants || []));

let variantIndex = $ref(0);

let isReadonly = $computed(
  () => props.group.Origin.Visibility != Visibility.Personnal
);

async function updateExercicegroup() {
  if (isReadonly) return;
  await controller.EditorUpdateExercicegroup(group.Group);
}

let exerciceToDelete: ExerciceHeader | null = $ref(null);
async function deleteVariante() {
  await controller.EditorDeleteExercice({ id: exerciceToDelete!.Id });

  variants = variants.filter((qu) => qu.Id != exerciceToDelete!.Id);
  exerciceToDelete = null;

  if (variants.length && variantIndex >= variants.length) {
    variantIndex = 0;
  }
  // if there is no more variant, that means the exercicegroup is deleted:
  // go back
  if (!variants.length) {
    emit("back");
  }
}

async function duplicateVariante(exercice: ExerciceHeader) {
  const newExercice = await controller.EditorDuplicateExercice({
    id: exercice.Id,
  });
  if (newExercice == undefined) {
    return;
  }
  variants.push(newExercice);
  variantIndex = variants.length - 1; // go to the new exercice
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
