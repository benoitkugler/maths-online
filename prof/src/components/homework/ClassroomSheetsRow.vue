<template>
  <v-dialog v-model="showNotes">
    <classroom-notes
      v-if="showNotes"
      :classroom="props.classroom.Classroom"
      :sheets="
        props.classroom.Sheets?.filter((sh) =>
          selectedSheets.has(sh.Sheet.Id)
        ) || []
      "
    ></classroom-notes>
  </v-dialog>

  <v-card class="mx-2 bg-grey-lighten-4">
    <v-row no-gutters>
      <v-col>
        <v-card-title>{{ classroom.Classroom.name }}</v-card-title>
      </v-col>
      <v-col align-self="center" cols="auto">
        <v-btn
          @click="onShowNotes"
          title="Afficher les notes pour les feuilles sélectionnées"
          :disabled="inSelect && Array.from(selectedSheets).length == 0"
        >
          {{ inSelect ? "Afficher" : "Notes" }}
        </v-btn>
        <v-btn
          class="mx-2"
          @click="emit('add')"
          title="Créer une nouvelle fiche"
        >
          <v-icon icon="mdi-plus" color="success"></v-icon>
          Ajouter
        </v-btn>
      </v-col>
    </v-row>
    <v-card-text>
      <v-row>
        <v-col
          v-if="!classroom.Sheets?.length"
          align-self="center"
          style="text-align: center"
        >
          <i>Aucune fiche n'est encore définie.</i>
        </v-col>

        <v-col
          lg="3"
          md="4"
          sm="6"
          xs="12"
          v-for="(sheet, index) in classroom.Sheets"
          :key="index"
        >
          <sheet-card
            :sheet="sheet"
            :classrooms="props.classrooms"
            :status="sheetStatus(sheet)"
            @delete="emit('delete', sheet)"
            @copy="(idClassroom) => emit('copy', sheet.Sheet.Id, idClassroom)"
            @update="emit('update', sheet)"
            @[mayClick]="() => onToggle(sheet)"
          ></sheet-card>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type {
  Classroom,
  ClassroomSheets,
  SheetExt,
} from "@/controller/api_gen";
import { $computed, $ref } from "vue/macros";
import ClassroomNotes from "./ClassroomNotes.vue";
import SheetCard from "./SheetCard.vue";
import { SheetStatus } from "./utils";

interface Props {
  classroom: ClassroomSheets;
  classrooms: Classroom[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "add"): void;
  (e: "update", sheet: SheetExt): void;
  (e: "delete", sheet: SheetExt): void;
  (e: "copy", idSheet: number, idClassroom: number): void;
}>();

const mayClick = $computed(() => (inSelect ? "click" : ""));

let inSelect = $ref(false);
let selectedSheets = $ref(
  new Set(props.classroom.Sheets?.map((sh) => sh.Sheet.Id))
);

let showNotes = $ref(false);

function onToggle(sh: SheetExt) {
  if (selectedSheets.has(sh.Sheet.Id)) {
    selectedSheets.delete(sh.Sheet.Id);
  } else {
    selectedSheets.add(sh.Sheet.Id);
  }
}
function onShowNotes() {
  if (!inSelect) {
    inSelect = true;
  } else {
    inSelect = false;
    showNotes = true;
  }
}

function sheetStatus(sh: SheetExt) {
  if (!inSelect) return SheetStatus.normal;
  return selectedSheets.has(sh.Sheet.Id)
    ? SheetStatus.selected
    : SheetStatus.notSelected;
}
</script>
