<template>
  <v-card min-width="1200">
    <v-row class="mx-2 mt-1">
      <v-col>
        <v-card-title>Modifier la fiche</v-card-title>
      </v-col>

      <v-col style="text-align: right">
        <v-btn icon @click="emit('close')" variant="text">
          <v-icon icon="mdi-close"></v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text>
      <v-row>
        <!-- fields -->
        <v-col cols="6">
          <v-row class="mb-2"
            ><v-col>
              <v-text-field
                variant="outlined"
                density="compact"
                label="Titre de la fiche"
                v-model="props.sheet.Sheet.Title"
                @blur="update"
                hide-details
              >
              </v-text-field> </v-col
          ></v-row>

          <v-row class="mt-2"
            ><v-col>
              <NotationField
                v-model="props.sheet.Sheet.Notation"
                @update:model-value="update"
              ></NotationField> </v-col
          ></v-row>

          <v-row>
            <v-col>
              <TimeField
                v-model="props.sheet.Sheet.Deadline"
                @update:model-value="update"
              ></TimeField>
            </v-col>
          </v-row>

          <v-row>
            <v-col>
              <v-checkbox
                :color="color"
                class="mt-3"
                density="compact"
                v-model="props.sheet.Sheet.Activated"
                label="Fiche active"
                :messages="
                  props.sheet.Sheet.Activated
                    ? `Désactiver la feuille, la rendant invisible aux élèves. 
        Pour empêcher la modification des notes, modifier plutôt la date de cloture.`
                    : `Activer la feuille, la rendant visible aux élèves. Pour une feuille notée, la date de cloture décide si les notes sont modifiables ou non.`
                "
                @update:model-value="update"
              >
              </v-checkbox>
            </v-col>
          </v-row>
        </v-col>

        <!-- exercice list -->
        <v-col cols="6">
          <SheetTasks
            :sheet="props.sheet"
            @add="(v) => emit('addTask', props.sheet.Sheet, v)"
            @remove="(v) => emit('removeTask', props.sheet.Sheet, v)"
            @reorder="(v) => emit('reorderTasks', props.sheet.Sheet, v)"
          ></SheetTasks>
        </v-col>
      </v-row>
    </v-card-text>
    <v-card-actions>
      <v-btn @click="emit('close')">Retour</v-btn>
      <v-spacer></v-spacer>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import type {
  ExerciceHeader,
  Sheet,
  SheetExt,
  TaskExt,
} from "@/controller/api_gen";
import { computed } from "vue";
import NotationField from "./NotationField.vue";
import SheetTasks from "./SheetTasks.vue";
import TimeField from "./TimeField.vue";

interface Props {
  sheet: SheetExt;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "update", sheet: Sheet): void;
  (e: "addTask", sheet: Sheet, ex: ExerciceHeader): void;
  (e: "removeTask", sheet: Sheet, task: TaskExt): void;
  (e: "reorderTasks", sheet: Sheet, tasks: TaskExt[]): void;
}>();

const color = computed(() =>
  props.sheet.Sheet.Activated ? "blue-lighten-4" : "grey-lighten-4"
);

function update() {
  emit("update", props.sheet.Sheet);
}
</script>
