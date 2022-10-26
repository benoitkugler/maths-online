<template>
  <v-dialog
    :model-value="questionToDelete != null"
    @update:model-value="questionToDelete = null"
    max-width="800"
  >
    <v-card title="Confirmer">
      <v-card-text
        >Etes-vous certain de vouloir supprimer la variante
        <i>{{ questionToDelete?.Id }} - {{ questionToDelete?.Subtitle }}</i> ?
        <br />
        Cette opération est irréversible.

        <div v-if="ownVariants.length == 1" class="mt-2">
          La question associée sera aussi supprimée.
        </div>
      </v-card-text>
      <v-card-actions>
        <v-btn @click="questionToDelete = null">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn color="red" @click="deleteVariante" variant="outlined">
          Supprimer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="mt-1 px-2">
    <v-row no-gutters>
      <v-col cols="auto" align-self="center" class="pr-2">
        <v-btn
          size="small"
          icon
          title="Retour aux questions"
          @click="backToList"
        >
          <v-icon icon="mdi-arrow-left"></v-icon>
        </v-btn>
      </v-col>

      <v-col>
        <v-text-field
          class="my-2 input-small"
          variant="outlined"
          density="compact"
          label="Nom de la question"
          v-model="group.Group.Title"
          :readonly="isReadonly"
          @blur="updateQuestiongroup"
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
        <QuestionVariantsSelector
          :variants="ownVariants"
          :readonly="isReadonly"
          v-model="variantIndex"
          @delete="(qu) => (questionToDelete = qu)"
          @duplicate="duplicateVariante"
        ></QuestionVariantsSelector>
      </v-col>
    </v-row>

    <QuestionVariantPannel
      :question="ownVariants[variantIndex]"
      :readonly="isReadonly"
      :session_id="props.session_id"
      :all-tags="props.allTags"
      :show-variant-meta="ownVariants.length >= 2"
      @update="(qu) => (ownVariants[variantIndex] = qu)"
    ></QuestionVariantPannel>
  </v-card>
</template>

<script setup lang="ts">
import type { Question, QuestiongroupExt } from "@/controller/api_gen";
import { Visibility } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { copy } from "@/controller/utils";
import { $computed, $ref } from "vue/macros";
import TagListField from "../TagListField.vue";
import QuestionVariantPannel from "./QuestionVariantPannel.vue";
import QuestionVariantsSelector from "./QuestionVariantsSelector.vue";

interface Props {
  session_id: string;
  group: QuestiongroupExt;
  variants: Question[];
  allTags: string[]; // to provide auto completion
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
}>();

let group = $ref(copy(props.group));
let ownVariants = $ref(copy(props.variants));

let variantIndex = $ref(0);

let isReadonly = $computed(
  () => props.group.Origin.Visibility != Visibility.Personnal
);

async function updateQuestiongroup() {
  if (isReadonly) return;
  await controller.EditorUpdateQuestiongroup(group.Group);
}

let questionToDelete: Question | null = $ref(null);
async function deleteVariante() {
  await controller.EditorDeleteQuestion({ id: questionToDelete!.Id });

  ownVariants = ownVariants.filter((qu) => qu.Id != questionToDelete!.Id);
  questionToDelete = null;

  if (ownVariants.length && variantIndex >= ownVariants.length) {
    variantIndex = 0;
  }
  // if there is no more variant, that means the questiongroup is deleted:
  // go back
  if (!ownVariants.length) {
    emit("back");
  }
}

async function duplicateVariante(question: Question) {
  const newQuestion = await controller.EditorDuplicateQuestion({
    id: question.Id,
  });
  if (newQuestion == undefined) {
    return;
  }
  ownVariants.push(newQuestion);
  variantIndex = ownVariants.length - 1; // go to the new question
}

async function saveTags(newTags: string[]) {
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
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
