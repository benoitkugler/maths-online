<template>
  <v-dialog
    :model-value="sheetToUpdate != null"
    @update:model-value="sheetToUpdate = null"
    :retain-focus="false"
  >
    <sheet-details
      v-if="sheetToUpdate != null"
      :sheet="sheetToUpdate"
      @update="updateSheet"
      @close="sheetToUpdate = null"
    >
    </sheet-details>
  </v-dialog>

  <v-dialog
    :model-value="sheetToDelete != null"
    @update:model-value="sheetToDelete = null"
  >
    <v-card title="Confirmer">
      <v-card-text
        >Etes-vous certain de vouloir supprimer la fiche
        <i>{{ sheetToDelete?.Sheet.Title }}</i> ? <br />
        La progression éventuelle des élèves sera perdue, et cette opération est
        irréversible.
      </v-card-text>
      <v-card-actions>
        <v-btn @click="sheetToDelete = null">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn color="red" @click="deleteSheet" variant="elevated">
          Supprimer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- <v-dialog
    fullscreen
    :model-value="showMonitor"
    @update:model-value="showMonitorChanged"
    :retain-focus="false"
  >
    <session-monitor @closed="closeMonitor"></session-monitor>
  </v-dialog> -->

  <v-card class="my-5 mx-auto" width="80%">
    <v-row>
      <v-col>
        <v-card-title>Travail à la maison</v-card-title>
        <v-card-subtitle
          >Configurer les exercices à faire à la maison</v-card-subtitle
        >
      </v-col>
    </v-row>

    <v-row v-for="(classroom, index) in classrooms" :key="index" class="my-1">
      <v-col>
        <classroom-sheets-row
          :classroom="classroom"
          :classrooms="classroomList"
          @add="createSheet(classroom.Classroom.id)"
          @delete="(sheet) => (sheetToDelete = sheet)"
          @copy="copySheet"
          @update="(s) => (sheetToUpdate = copy(s))"
        ></classroom-sheets-row>
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import type { ClassroomSheets, SheetExt } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { copy } from "@/controller/utils";
import { computed, onActivated, onMounted } from "vue";
import { $ref } from "vue/macros";
import ClassroomSheetsRow from "../components/homework/ClassroomSheetsRow.vue";
import SheetDetails from "../components/homework/SheetDetails.vue";

let classrooms = $ref<ClassroomSheets[]>([]);

onMounted(() => {
  fetchClassrooms();
});

onActivated(fetchClassrooms);

const classroomList = computed(() => classrooms.map((cl) => cl.Classroom));

let sheetToUpdate = $ref<SheetExt | null>(null);
let sheetToDelete = $ref<SheetExt | null>(null);

async function fetchClassrooms() {
  const res = await controller.HomeworkGetSheets();
  if (res == undefined) {
    return;
  }
  classrooms = res;
}

async function createSheet(idClassroom: number) {
  const res = await controller.HomeworkCreateSheet({
    IdClassroom: idClassroom,
  });
  if (res == undefined) {
    return;
  }
  const cl = classrooms.find((cl) => cl.Classroom.id == idClassroom)!;
  cl.Sheets = (cl.Sheets || []).concat(res);
}

async function deleteSheet() {
  if (sheetToDelete == null) {
    return;
  }
  const id = sheetToDelete.Sheet.Id;
  const res = await controller.HomeworkDeleteSheet({
    id: id,
  });
  if (res == undefined) {
    return;
  }
  const cl = classrooms.find(
    (cl) => cl.Classroom.id == sheetToDelete?.Sheet.IdClassroom
  )!;
  cl.Sheets = (cl.Sheets || []).filter((sh) => sh.Sheet.Id != id);
  sheetToDelete = null;
}

async function copySheet(idSheet: number, idClassroom: number) {
  const res = await controller.HomeworkCopySheet({
    IdSheet: idSheet,
    IdClassroom: idClassroom,
  });
  if (res == undefined) {
    return;
  }
  const cl = classrooms.find((cl) => cl.Classroom.id == idClassroom)!;
  cl.Sheets = (cl.Sheets || []).concat(res);
}

async function updateSheet(sheet: SheetExt) {
  await controller.HomeworkUpdateSheet({
    Sheet: sheet.Sheet,
    Exercices: (sheet.Exercices || []).map((ex) => ex.Exercice.Id),
  });

  const cl = classrooms.find(
    (cl) => cl.Classroom.id == sheet.Sheet.IdClassroom
  )!;
  const index = cl.Sheets!.findIndex((s) => s.Sheet.Id == sheet.Sheet.Id);
  cl.Sheets![index] = sheet;
}
</script>

<style>
:deep(.v-dialog .v-overlay__content) {
  max-width: 900px;
}
</style>
