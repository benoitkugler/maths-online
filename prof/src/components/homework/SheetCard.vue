<template>
  <v-card :title="sheet.Sheet.Title" :color="color">
    <v-card-text>
      <v-row>
        <v-col>
          <v-chip>{{ sheet.Tasks?.length || 0 }} tâche(s) </v-chip>
        </v-col>
        <v-col style="text-align: right">
          <v-chip :color="sheet.Sheet.Notation == 0 ? 'black' : 'primary'">
            <span v-if="sheet.Sheet.Notation != 0">Notée</span>
            <span v-else>Non notée</span>
          </v-chip>
        </v-col>
      </v-row>

      <v-row no-gutters v-if="deadline" class="mt-2">
        <v-col align-self="center"> Clôture :</v-col>
        <v-col align-self="center">
          <v-chip style="text-align: right">
            {{ deadline }}
          </v-chip>
        </v-col>
      </v-row>
    </v-card-text>
    <v-card-actions>
      <v-btn
        density="comfortable"
        :color="sheet.Sheet.Activated ? 'grey' : 'red'"
        @click="emit('delete')"
        icon
        title="Supprimer la fiche"
        :disabled="sheet.Sheet.Activated"
        size="small"
      >
        <v-icon icon="mdi-delete"></v-icon>
      </v-btn>

      <v-menu offset-y close-on-content-click>
        <template v-slot:activator="{ isActive, props }">
          <v-btn
            v-on="{ isActive }"
            v-bind="props"
            density="comfortable"
            icon
            title="Copier vers..."
            size="small"
          >
            <v-icon icon="mdi-content-copy" color="secondary"></v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item
            v-for="(classroom, index) in classrooms"
            :key="index"
            @click="emit('copy', classroom.id)"
          >
            {{ classroom.name }}
          </v-list-item>
        </v-list>
      </v-menu>

      <v-spacer></v-spacer>
      <v-btn title="Modifier le contenu de la fiche" @click="emit('update')"
        >Modifier</v-btn
      >
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import type { Classroom, SheetExt } from "@/controller/api_gen";
import { formatTime } from "@/controller/utils";
import { computed } from "vue";

interface Props {
  sheet: SheetExt;
  classrooms: Classroom[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "delete"): void;
  (e: "update"): void;
  (e: "copy", target: number): void;
}>();

const color = computed(() =>
  props.sheet.Sheet.Activated ? "blue-lighten-4" : "grey-lighten-4"
);

const deadline = computed(() => formatTime(props.sheet.Sheet.Deadline));
</script>
