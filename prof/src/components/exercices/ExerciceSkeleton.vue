<template>
  <v-dialog v-model="showEditDescription">
    <description-pannel
      v-model="props.exercice.Exercice.Description"
      :readonly="isReadonly"
    ></description-pannel>
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
        <v-row no-gutters>
          <v-col>
            <v-text-field
              class="my-2 input-small"
              variant="outlined"
              density="compact"
              label="Nom de l'exercice"
              v-model="props.exercice.Exercice.Title"
              :readonly="isReadonly"
              hide-details
            ></v-text-field
          ></v-col>
          <v-col cols="auto" align-self="center">
            <v-btn
              class="mx-2"
              icon
              @click="save"
              :title="
                isReadonly ? 'Visualiser' : 'Enregistrer et prévisualiser'
              "
              size="small"
            >
              <v-icon
                :icon="isReadonly ? 'mdi-eye' : 'mdi-content-save'"
                size="small"
              ></v-icon>
            </v-btn>

            <v-menu offset-y close-on-content-click>
              <template v-slot:activator="{ isActive, props }">
                <v-btn
                  icon
                  title="Plus d'options"
                  v-on="{ isActive }"
                  v-bind="props"
                  size="x-small"
                >
                  <v-icon icon="mdi-dots-vertical"></v-icon>
                </v-btn>
              </template>
              <v-list>
                <v-list-item>
                  <v-btn
                    size="small"
                    @click="showEditDescription = true"
                    title="Editer le commentaire"
                  >
                    <v-icon
                      class="mr-2"
                      icon="mdi-message-reply-text"
                      size="small"
                    ></v-icon>
                    Commentaire
                  </v-btn>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-col>
        </v-row>

        <v-row no-gutters>
          <!-- <v-col class="pr-2">
            <tag-list-field
              label="Catégories"
              v-model="tags"
              :all-tags="props.allTags"
              @update:model-value="saveTags"
              :readonly="isReadonly"
            ></tag-list-field
          ></v-col> -->
          <v-col cols="auto">
            <v-btn
              title="Créer et ajouter une question"
              size="small"
              @click="addQuestion"
            >
              <v-icon icon="mdi-plus" color="green"></v-icon>
              Ajouter une question
            </v-btn>
            <v-btn
              title="Importer une question existante"
              size="small"
              class="mx-2"
            >
              <v-icon icon="mdi-plus" color="green"></v-icon>
              Importer une question
            </v-btn>
          </v-col>
        </v-row>
      </v-col>
    </v-row>

    <v-list>
      <v-list-item
        v-for="(question, index) in props.exercice.Questions"
        :key="index"
      >
        {{ question }}
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import { Visibility, type ExerciceExt } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { computed } from "vue";
import { $ref } from "vue/macros";
import DescriptionPannel from "../editor/DescriptionPannel.vue";

interface Props {
  exercice: ExerciceExt;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "back"): void;
  (e: "delete", exercice: ExerciceExt): void;
  (e: "clicked", exercice: ExerciceExt): void;
  (e: "duplicate", exercice: ExerciceExt): void;
  (e: "updatePublic", exerciceID: number, isPublic: boolean): void;
}>();

const isReadonly = computed(
  () => props.exercice.Origin.Visibility != Visibility.Personnal
);

function backToList() {
  emit("back");
}

let showEditDescription = $ref(false);

async function save() {
  console.log("TODO");
}

async function addQuestion() {
  const l = await controller.ExerciceAddQuestion({
    IdExercice: props.exercice.Exercice.Id,
    IdQuestion: -1,
  });
  if (l == undefined) {
    return;
  }
  props.exercice.Questions = l;
}
</script>

<style></style>
