<template>
  <!-- <v-dialog
    :model-value="questionToDelete != null"
    @update:model-value="questionToDelete = null"
  >
    <v-card title="Confirmer">
      <v-card-text
        >Etes-vous certain de vouloir supprimer la question
        <i>{{ questionToDelete?.Title }}</i> ? <br />
        Cette opération est irréversible.
      </v-card-text>
      <v-card-actions>
        <v-btn @click="questionToDelete = null">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn color="red" @click="deleteQuestion" variant="outlined">
          Supprimer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog> -->

  <v-card class="mt-3 px-2">
    <v-row no-gutters class="mb-2">
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
              hide-details
            ></v-text-field
          >
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
              :variants="variants"
              :readonly="isReadonly"
              v-model="variantIndex"
              ></VariantsSelector>
      </v-col>
    </v-row>

    TODO
  </v-card>
</template>

<script setup lang="ts">
import type { Question, QuestiongroupExt } from "@/controller/api_gen";
import { Visibility } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { copy } from "@/controller/utils";
import { computed } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import TagListField from "../TagListField.vue";
import VariantsSelector from "./VariantsSelector.vue";

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
let variants = $ref(copy(props.variants));

let variantIndex = $ref(0);

let questionToDelete: Question | null = $ref(null);
async function deleteQuestion() {
  await controller.EditorDeleteQuestion({ id: questionToDelete!.Id });

  questionToDelete = null;
}

const isReadonly = computed(
  () => props.group.Origin.Visibility != Visibility.Personnal
);

async function saveTags(newTags: string[]) {
  const rep = await controller.EditorUpdateTags({ Id: group.Group.Id, Tags: newTags });
  if (rep === undefined) {
    return;
  }
  group.Tags = newTags;
}

async function duplicate(question: Question) {
  const newQuestion = await controller.EditorDuplicateQuestion({
    id: question.Id,
  });
  if (newQuestion == undefined) {
    return;
  }
  variants.push(newQuestion);
  variantIndex += 1; // go to the new question
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
