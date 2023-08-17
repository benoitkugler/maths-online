<template>
  <v-card density="compact">
    <v-card-text>
      <v-row justify="space-between">
        <v-col align-self="center" cols="auto">
          <v-menu v-if="sheet.Origin.Visibility == Visibility.Admin">
            <template v-slot:activator="{ isActive, props }">
              <v-card
                v-on="{ isActive }"
                v-bind="props"
                variant="elevated"
                :title="sheet.Sheet.Title"
                :subtitle="subtitle"
                color="yellow"
              >
              </v-card>
            </template>
            <PreviewSheet :sheet="sheet"></PreviewSheet>
          </v-menu>
          <v-card
            v-else
            variant="elevated"
            @click="emit('editSheet', props.sheet)"
            :title="sheet.Sheet.Title"
            :subtitle="subtitle"
            color="grey-lighten-3"
          >
          </v-card>
        </v-col>
        <v-col cols="auto" align-self="center">
          <v-tooltip
            text="Enregistrer cette feuille dans les favoris"
            location="top"
            v-if="sheet.Sheet.Anonymous.Valid"
          >
            <template v-slot:activator="{ isActive, props }">
              <v-btn
                density="comfortable"
                icon
                size="small"
                class="mr-1"
                v-bind="props"
                v-on="{ isActive }"
                @click="emit('setFavorite', sheet.Sheet)"
              >
                <v-icon color="secondary"> mdi-heart </v-icon>
              </v-btn>
            </template>
          </v-tooltip>
          <v-menu offset-y close-on-content-click>
            <template v-slot:activator="{ isActive, props }">
              <v-btn
                v-on="{ isActive }"
                v-bind="props"
                density="comfortable"
                icon
                size="small"
                title="Copier vers une autre classe..."
                class="mr-1"
              >
                <v-icon icon="mdi-content-copy" size="small"></v-icon>
              </v-btn>
            </template>
            <v-list>
              <v-list-subheader>Copier vers...</v-list-subheader>
              <v-list-item
                v-for="(classroom, index) in classrooms"
                :key="index"
                link
                @click="emit('copy', classroom.id)"
              >
                {{ classroom.name }}
              </v-list-item>
            </v-list>
          </v-menu>

          <v-btn
            density="comfortable"
            @click="emit('delete')"
            icon
            title="Supprimer le travail"
            size="small"
          >
            <v-icon color="red" icon="mdi-delete" size="small"></v-icon>
          </v-btn>
        </v-col>
      </v-row>
      <v-row justify="space-between" class="mt-0">
        <v-col cols="auto" align-self="center">
          <v-switch
            label="Accès à durée limité"
            v-model="travail.Noted"
            @update:model-value="emit('update', travail)"
            hide-details
            color="primary"
          >
          </v-switch>
        </v-col>

        <v-col v-if="travail.Noted" cols="auto" align-self="center">
          Clôture :
          <v-menu
            offset-y
            :close-on-content-click="false"
            :model-value="deadlineToEdit != null"
            @update:model-value="deadlineToEdit = null"
          >
            <template v-slot:activator="{ isActive, props }">
              <v-chip
                v-on="{ isActive }"
                v-bind="props"
                style="text-align: right"
                class="ml-1"
                color="primary"
                variant="outlined"
                @click="deadlineToEdit = travail.Deadline"
              >
                {{ deadline }}
              </v-chip>
            </template>
            <v-card width="300px">
              <v-card-text class="pb-0" v-if="deadlineToEdit != null">
                <TimeField v-model="deadlineToEdit"></TimeField>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  color="success"
                  @click="
                    travail.Deadline = deadlineToEdit!;
                    deadlineToEdit = null;
                    emit('update', travail);
                  "
                  >Enregistrer</v-btn
                >
              </v-card-actions>
            </v-card>
          </v-menu>
        </v-col>
        <v-col cols="auto" align-self="center" v-else>
          <v-chip> Feuille en accès libre </v-chip>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import {
  Visibility,
  type Classroom,
  type Sheet,
  type SheetExt,
  type Time,
  type Travail
} from "@/controller/api_gen";
import { formatTime } from "@/controller/utils";
import { computed } from "vue";
import TimeField from "./TimeField.vue";
import { $ref } from "vue/macros";
import PreviewSheet from "./PreviewSheet.vue";

interface Props {
  travail: Travail;
  sheet: SheetExt;
  classrooms: Classroom[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "delete"): void;
  (e: "update", tr: Travail): void;
  (e: "copy", idTarget: number): void;
  (e: "setFavorite", sheet: Sheet): void;
  (e: "editSheet", sheet: SheetExt): void;
}>();

const nbTasks = computed(() => props.sheet.Tasks?.length || 0);

const subtitle = computed(() => {
  if (nbTasks.value == 0) {
    return "Aucune tâche";
  } else if (nbTasks.value == 1) {
    return "1 tâche";
  } else {
    return `${nbTasks.value} tâches`;
  }
});

const deadline = computed(() => formatTime(props.travail.Deadline));

let deadlineToEdit = $ref<Time | null>(null);
</script>
