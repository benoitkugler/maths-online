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
      @addExercice="addExerciceToSheet"
      @addMonoquestion="addMonoquestionToSheet"
      @udpate-monoquestion="updateMonoquestion"
      @removeTask="removeTaskFromSheet"
      @reorderTasks="reorderSheetTasks"
      @close="sheetToUpdate = null"
    >
    </sheet-details>
  </v-dialog>

  <v-dialog
    :model-value="sheetToDelete != null"
    @update:model-value="sheetToDelete = null"
    max-width="800px"
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
          @update="(s) => (sheetToUpdate = s)"
        ></classroom-sheets-row>
      </v-col>
    </v-row>

    <v-row v-if="!classrooms.length" justify="center">
      <v-col style="text-align: center" class="my-3">
        <i>Aucune classe n'est enregistrée.</i>
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import type {
  ClassroomSheets,
  ExerciceHeader,
  Monoquestion,
  QuestionHeader,
  Sheet,
  SheetExt,
  TaskExt,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { computed, onActivated, onMounted } from "vue";
import { useRoute } from "vue-router";
import { $ref } from "vue/macros";
import ClassroomSheetsRow from "../components/homework/ClassroomSheetsRow.vue";
import SheetDetails from "../components/homework/SheetDetails.vue";

let classrooms = $ref<ClassroomSheets[]>([]);

const route = useRoute();

onMounted(() => {
  fetchClassrooms();
});

onActivated(async () => {
  await fetchClassrooms();
  showDetailsFromQuery();
});

const classroomList = computed(() => classrooms.map((cl) => cl.Classroom));

let sheetToUpdate = $ref<SheetExt | null>(null);
let sheetToDelete = $ref<SheetExt | null>(null);

function showDetailsFromQuery() {
  const idSheet = Number(route.query["idSheet"]);
  if (isNaN(idSheet)) return;
  classrooms.forEach((cl) =>
    cl.Sheets?.forEach((sh) => {
      if (sh.Sheet.Id == idSheet) {
        sheetToUpdate = sh;
        return;
      }
    })
  );
}

async function fetchClassrooms() {
  const res = await controller.HomeworkGetSheets();
  if (res == undefined) {
    return;
  }
  classrooms = res;
}

function indexes(idSheet: number, idClassroom: number) {
  const clIndex = classrooms.findIndex((cl) => cl.Classroom.id == idClassroom)!;
  const sheetIndex = (classrooms[clIndex].Sheets || []).findIndex(
    (s) => s.Sheet.Id == idSheet
  );
  return { clIndex, sheetIndex };
}

async function createSheet(idClassroom: number) {
  const res = await controller.HomeworkCreateSheet({
    IdClassroom: idClassroom,
  });
  if (res == undefined) {
    return;
  }

  const { clIndex } = indexes(-1, idClassroom);
  const cl = classrooms[clIndex];
  console.log(clIndex, cl);

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

  const { clIndex } = indexes(-1, sheetToDelete?.Sheet.IdClassroom);
  const cl = classrooms[clIndex];
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

  const { clIndex } = indexes(idSheet, idClassroom);
  const cl = classrooms[clIndex];
  cl.Sheets = (cl.Sheets || []).concat(res);
}

async function updateSheet(sheet: Sheet) {
  await controller.HomeworkUpdateSheet(sheet);

  const { clIndex, sheetIndex } = indexes(sheet.Id, sheet.IdClassroom);
  const cl = classrooms[clIndex];
  cl.Sheets![sheetIndex].Sheet = sheet;
}

async function addExerciceToSheet(sheet: Sheet, exercice: ExerciceHeader) {
  const newTask = await controller.HomeworkAddExercice({
    IdExercice: exercice.Id,
    IdSheet: sheet.Id,
  });
  if (newTask == undefined) {
    return;
  }

  const { clIndex, sheetIndex } = indexes(sheet.Id, sheet.IdClassroom);
  const cl = classrooms[clIndex];
  cl.Sheets![sheetIndex].Tasks = (cl.Sheets![sheetIndex].Tasks || []).concat(
    newTask
  );
}

async function addMonoquestionToSheet(sheet: Sheet, question: QuestionHeader) {
  const newTask = await controller.HomeworkAddMonoquestion({
    IdQuestion: question.Id,
    IdSheet: sheet.Id,
  });
  if (newTask == undefined) {
    return;
  }

  const { clIndex, sheetIndex } = indexes(sheet.Id, sheet.IdClassroom);
  const cl = classrooms[clIndex];
  cl.Sheets![sheetIndex].Tasks = (cl.Sheets![sheetIndex].Tasks || []).concat(
    newTask
  );
}

async function updateMonoquestion(sheet: Sheet, qu: Monoquestion) {
  const task = await controller.HomeworkUpdateMonoquestion(qu);
  if (task == undefined) return;
  const { clIndex, sheetIndex } = indexes(sheet.Id, sheet.IdClassroom);
  const cl = classrooms[clIndex];
  const tasks = cl.Sheets![sheetIndex].Tasks || [];
  const index = tasks.findIndex((v) => v.Id == task.Id);
  tasks[index] = task;
}

async function removeTaskFromSheet(sheet: Sheet, task: TaskExt) {
  await controller.HomeworkRemoveTask({ "id-task": task.Id });

  const { clIndex, sheetIndex } = indexes(sheet.Id, sheet.IdClassroom);
  const cl = classrooms[clIndex];
  cl.Sheets![sheetIndex].Tasks = (cl.Sheets![sheetIndex].Tasks || []).filter(
    (ta) => ta.Id != task.Id
  );
}

async function reorderSheetTasks(sheet: Sheet, tasks: TaskExt[]) {
  await controller.HomeworkReorderSheetTasks({
    IdSheet: sheet.Id,
    Tasks: tasks.map((t) => t.Id),
  });

  const { clIndex, sheetIndex } = indexes(sheet.Id, sheet.IdClassroom);
  const cl = classrooms[clIndex];
  cl.Sheets![sheetIndex].Tasks = tasks;
}
</script>

<style>
:deep(.v-dialog .v-overlay__content) {
  max-width: 900px;
}
</style>
