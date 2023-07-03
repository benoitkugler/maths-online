<template>
  <v-card>
    <v-row class="mt-0 mb-0 mx-1">
      <v-col>
        <v-card-title>Paramètres de la feuille</v-card-title>
      </v-col>

      <v-col style="text-align: right" align-self="center">
        <v-btn icon @click="emit('close')" variant="text">
          <v-icon icon="mdi-close"></v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text>
      <v-row>
        <!-- fields -->
        <v-col cols="4" align-self="center">
          <v-row>
            <v-col>
              <v-text-field
                variant="outlined"
                density="compact"
                label="Titre de la feuille"
                v-model="props.sheet.Sheet.Title"
                @blur="update"
                hide-details
              >
              </v-text-field> </v-col
          ></v-row>

          <v-row>
            <v-col>
              <v-combobox
                variant="outlined"
                density="compact"
                label="Niveau (classe)"
                v-model="props.sheet.Sheet.Level"
                :items="allTags.Levels || []"
                :color="LevelColor"
                @blur="update"
              ></v-combobox>
            </v-col>
          </v-row>
        </v-col>

        <!-- exercice list -->
        <v-col cols="8">
          <SheetTasks
            :sheet="props.sheet"
            :all-tags="allTags"
            @add-exercice="(v) => emit('addExercice', props.sheet.Sheet, v)"
            @add-monoquestion="
              (v) => emit('addMonoquestion', props.sheet.Sheet, v)
            "
            @add-random-monoquestion="
              (v) => emit('addRandomMonoquestion', props.sheet.Sheet, v)
            "
            @update-monoquestion="
              (v) => emit('udpateMonoquestion', props.sheet.Sheet, v)
            "
            @update-random-monoquestion="
              (v) => emit('udpateRandomMonoquestion', props.sheet.Sheet, v)
            "
            @remove="(v) => emit('removeTask', props.sheet.Sheet, v)"
            @reorder="(v) => emit('reorderTasks', props.sheet.Sheet, v)"
          ></SheetTasks>
        </v-col>
      </v-row>
    </v-card-text>
    <v-card-actions class="py-0 my-0"></v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import type {
  Monoquestion,
  RandomMonoquestion,
  Sheet,
  SheetExt,
  TagsDB,
  TaskExt,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import {
  LevelColor,
  emptyTagsDB,
  type VariantG,
  type ResourceGroup,
} from "@/controller/editor";
import { onMounted } from "vue";
import { $ref } from "vue/macros";
import SheetTasks from "./SheetTasks.vue";

interface Props {
  sheet: SheetExt;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "update", sheet: Sheet): void;
  (e: "addExercice", sheet: Sheet, ex: VariantG): void;
  (e: "addMonoquestion", sheet: Sheet, ex: VariantG): void;
  (e: "addRandomMonoquestion", sheet: Sheet, ex: ResourceGroup): void;
  (e: "udpateMonoquestion", sheet: Sheet, qu: Monoquestion): void;
  (e: "udpateRandomMonoquestion", sheet: Sheet, qu: RandomMonoquestion): void;
  (e: "removeTask", sheet: Sheet, task: TaskExt): void;
  (e: "reorderTasks", sheet: Sheet, tasks: TaskExt[]): void;
}>();

function update() {
  emit("update", props.sheet.Sheet);
}

onMounted(fetchTags);
let allTags = $ref<TagsDB>(emptyTagsDB());
async function fetchTags() {
  const tags = await controller.EditorGetTags();
  if (tags) allTags = tags;
}
</script>
