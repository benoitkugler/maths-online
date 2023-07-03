<template>
  <v-card>
    <v-row class="pa-2">
      <v-col align-self="center" cols="8">
        <v-tooltip location="top">
          <template v-slot:activator="{ props: innerProps }">
            <v-card-subtitle v-bind="innerProps" class="pl-1"
              >{{ sheet.Sheet.Title }}
            </v-card-subtitle>
          </template>
          {{ sheet.Sheet.Title }} ({{ nbTasks }} tâche{{
            nbTasks > 1 ? "s" : ""
          }})
        </v-tooltip>
      </v-col>
      <v-col cols="auto">
        <v-menu offset-y close-on-content-click>
          <template v-slot:activator="{ isActive, props }">
            <v-btn
              v-on="{ isActive }"
              v-bind="props"
              density="comfortable"
              icon
              title="Copier vers une autre classe..."
              size="small"
              class="mr-1"
            >
              <v-icon
                icon="mdi-content-copy"
                color="secondary"
                size="small"
              ></v-icon>
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
    <v-card-text class="pt-1">
      <v-row justify="center">
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
      </v-row>

      <div class="mt-2">
        <v-row v-if="travail.Noted" no-gutters>
          <v-col align-self="center"> Clôture :</v-col>
          <v-col align-self="center">
            <v-menu
              offset-y
              :close-on-content-click="false"
              :model-value="deadlineToEdit != null"
              @update:model-value="deadlineToEdit = null"
            >
              <template v-slot:activator="{ isActive, props }">
                <v-chip
                  v-on="isActive"
                  v-bind="props"
                  style="text-align: right"
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
        </v-row>
        <v-row no-gutters justify="center" v-else>
          <v-col cols="auto"><v-chip> Feuille en accès libre </v-chip> </v-col>
        </v-row>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { Classroom, SheetExt, Time, Travail } from "@/controller/api_gen";
import { formatTime } from "@/controller/utils";
import { computed } from "vue";
import TimeField from "./TimeField.vue";
import { $ref } from "vue/macros";

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
}>();

const nbTasks = computed(() => props.sheet.Tasks?.length || 0);

const deadline = computed(() => formatTime(props.travail.Deadline));

let deadlineToEdit = $ref<Time | null>(null);
</script>
