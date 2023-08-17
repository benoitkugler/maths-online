<template>
  <v-dialog
    :model-value="sheetToUpdate != null"
    @update:model-value="sheetToUpdate = null"
  >
    <sheet-details
      v-if="sheetToUpdate != null"
      :sheet="sheetToUpdate"
      @update="updateSheet"
      @addExercice="addExerciceToSheet"
      @addMonoquestion="addMonoquestionToSheet"
      @add-random-monoquestion="addRandomMonoquestionToSheet"
      @udpate-monoquestion="updateMonoquestion"
      @udpate-random-monoquestion="updateRandomMonoquestion"
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
        <v-btn @click="sheetToDelete = null" color="warning">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn color="red" @click="deleteSheet" variant="elevated">
          Supprimer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="showFavorites">
    <sheet-folder
      :sheets="homeworks.Sheets"
      :classrooms="homeworks.Travaux.map(t => t.Classroom)"
      @create="createSheet"
      @duplicate="sh => duplicateSheet(sh.Sheet.Id)"
      @assign="createTravailWith"
      @delete="sh => (sheetToDelete = sh)"
      @edit="sh => (sheetToUpdate = sh)"
      @update-public="updatePublic"
      @create-review="createReview"
    ></sheet-folder>
  </v-dialog>

  <v-card class="my-5 mx-auto" width="90%">
    <v-row>
      <v-col>
        <v-card-title>Travail à la maison</v-card-title>
        <v-card-subtitle
          >Configurer les exercices à faire à la maison, en assignant à chaque
          classe une ou plusieurs feuilles.</v-card-subtitle
        >
      </v-col>
      <v-col cols="auto" align-self="center">
        <v-btn
          class="mr-3"
          title="Afficher les feuilles favorites..."
          @click="
            fetchHomeworks();
            showFavorites = true;
          "
        >
          <v-icon class="mr-2" color="secondary">mdi-heart</v-icon>
          Favoris
        </v-btn>
      </v-col>
    </v-row>

    <travaux-pannel
      :homeworks="homeworks"
      @create="createTravail"
      @create-with="createTravailWith"
      @set-favorite="setSheetFavorite"
      @update="updateTravail"
      @copy="copyTravailTo"
      @delete="deleteTravail"
      @edit-sheet="sh => (sheetToUpdate = sh)"
      ref="travauxPannel"
    ></travaux-pannel>
  </v-card>
</template>

<script setup lang="ts">
import {
  ReviewKind,
  type Monoquestion,
  type RandomMonoquestion,
  type Sheet,
  type SheetExt,
  type TaskExt,
  type Travail
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onActivated, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { $ref } from "vue/macros";
import SheetDetails from "../components/homework/SheetDetails.vue";
import SheetFolder from "@/components/homework/SheetFolder.vue";
import type { ResourceGroup, VariantG } from "@/controller/editor";
import TravauxPannel from "@/components/homework/TravauxPannel.vue";
import type { HomeworksT } from "@/controller/utils";
import { ref } from "vue";

let homeworks = $ref<HomeworksT>({ Sheets: new Map(), Travaux: [] });

const route = useRoute();
const router = useRouter();

onMounted(() => {
  fetchHomeworks();
});

onActivated(async () => {
  await fetchHomeworks();
  showDetailsFromRouteQuery();
});

let travauxPannel = $ref<InstanceType<typeof TravauxPannel> | null>(null);

let showFavorites = $ref(false);

let sheetToUpdate = $ref<SheetExt | null>(null);
let sheetToDelete = $ref<SheetExt | null>(null);

function showDetailsFromRouteQuery() {
  const idSheet = Number(route.query["idSheet"]);
  if (isNaN(idSheet)) return;

  const sh = homeworks.Sheets.get(idSheet);
  if (sh === undefined) return;
  sheetToUpdate = sh;
}

async function fetchHomeworks() {
  const res = await controller.HomeworkGetSheets();
  if (res == undefined) return;
  homeworks = {
    Sheets: new Map(
      Object.entries(res.Sheets || {}).map(v => [Number(v[0]), v[1]])
    ),
    Travaux: res.Travaux || []
  };
}

/** return the place of the given classroom in the Travaux list */
function travauxByClassroom(idClassroom: number) {
  return (homeworks.Travaux || []).find(cl => cl.Classroom.id == idClassroom);
}

async function createSheet() {
  const res = await controller.HomeworkCreateSheet();
  if (res == undefined) {
    return;
  }

  homeworks.Sheets.set(res.Sheet.Id, res);

  sheetToUpdate = res;
}

async function createTravail(idClassroom: number) {
  const res = await controller.HomeworkCreateTravail({
    "id-classroom": idClassroom
  });
  if (res == undefined) {
    return;
  }

  homeworks.Sheets.set(res.Sheet.Sheet.Id, res.Sheet);
  const cl = travauxByClassroom(idClassroom)!;
  cl.Travaux = (cl.Travaux || []).concat(res.Travail);

  sheetToUpdate = res.Sheet;
}

async function createTravailWith(idSheet: number, idClassroom: number) {
  const res = await controller.HomeworkCreateTravailWith({
    IdClassroom: idClassroom,
    IdSheet: idSheet
  });
  if (res == undefined) {
    return;
  }

  const cl = travauxByClassroom(idClassroom)!;
  cl.Travaux = (cl.Travaux || []).concat(res);

  showFavorites = false;
  if (travauxPannel != null) travauxPannel.showClassroom(idClassroom);
}

async function deleteSheet() {
  if (sheetToDelete == null) {
    return;
  }
  const id = sheetToDelete.Sheet.Id;
  const res = await controller.HomeworkDeleteSheet({
    id: id
  });
  if (res == undefined) {
    return;
  }

  homeworks.Sheets.delete(id);
  homeworks.Travaux?.forEach(cl => {
    cl.Travaux = (cl.Travaux || []).filter(tr => tr.IdSheet != id);
  });
  sheetToDelete = null;
}

async function duplicateSheet(idSheet: number) {
  const res = await controller.HomeworkCopySheet({
    IdSheet: idSheet
  });
  if (res == undefined) {
    return;
  }

  homeworks.Sheets.set(res.Sheet.Id, res);
}

async function updateSheet(sheet: Sheet) {
  const ok = await controller.HomeworkUpdateSheet(sheet);
  if (ok == undefined) return;

  homeworks.Sheets.get(sheet.Id)!.Sheet = sheet;
}

async function updatePublic(sheet: Sheet, pub: boolean) {
  sheet.Public = pub;
  const ok = await controller.HomeworkUpdateSheet(sheet);
  if (ok == undefined) return;

  const sheetExt = homeworks.Sheets.get(sheet.Id)!;
  sheetExt.Origin.IsPublic = pub;
  sheetExt.Sheet = sheet;
  homeworks.Sheets.set(sheet.Id, sheetExt);
}

async function createReview(sheet: Sheet) {
  const res = await controller.ReviewCreate({
    Kind: ReviewKind.KSheet,
    Id: sheet.Id
  });
  if (res == undefined) return;

  router.push({ name: "reviews", query: { id: res.Id } });
}

async function addExerciceToSheet(sheet: Sheet, exercice: VariantG) {
  const newTask = await controller.HomeworkAddExercice({
    IdExercice: exercice.Id,
    IdSheet: sheet.Id
  });
  if (newTask == undefined) {
    return;
  }

  const sh = homeworks.Sheets.get(sheet.Id)!;
  sh.Tasks = (sh.Tasks || []).concat(newTask);
}

async function addMonoquestionToSheet(sheet: Sheet, question: VariantG) {
  const newTask = await controller.HomeworkAddMonoquestion({
    IdQuestion: question.Id,
    IdSheet: sheet.Id
  });
  if (newTask == undefined) {
    return;
  }

  const sh = homeworks.Sheets.get(sheet.Id)!;
  sh.Tasks = (sh.Tasks || []).concat(newTask);
}

async function addRandomMonoquestionToSheet(
  sheet: Sheet,
  group: ResourceGroup
) {
  const newTask = await controller.HomeworkAddRandomMonoquestion({
    IdQuestiongroup: group.Id,
    IdSheet: sheet.Id
  });
  if (newTask == undefined) {
    return;
  }

  const sh = homeworks.Sheets.get(sheet.Id)!;
  sh.Tasks = (sh.Tasks || []).concat(newTask);
}

async function updateMonoquestion(sheet: Sheet, qu: Monoquestion) {
  const task = await controller.HomeworkUpdateMonoquestion(qu);
  if (task == undefined) return;

  const sh = homeworks.Sheets.get(sheet.Id)!;
  const tasks = sh.Tasks || [];
  const index = tasks.findIndex(v => v.Id == task.Id);
  tasks[index] = task;
}

async function updateRandomMonoquestion(sheet: Sheet, qu: RandomMonoquestion) {
  const task = await controller.HomeworkUpdateRandomMonoquestion(qu);
  if (task == undefined) return;

  const sh = homeworks.Sheets.get(sheet.Id)!;
  const tasks = sh.Tasks || [];
  const index = tasks.findIndex(v => v.Id == task.Id);
  tasks[index] = task;
}

async function removeTaskFromSheet(sheet: Sheet, task: TaskExt) {
  const res = await controller.HomeworkRemoveTask({ "id-task": task.Id });
  if (res === undefined) return;

  const sh = homeworks.Sheets.get(sheet.Id)!;
  sh.Tasks = (sh.Tasks || []).filter(ta => ta.Id != task.Id);
}

async function reorderSheetTasks(sheet: Sheet, tasks: TaskExt[]) {
  await controller.HomeworkReorderSheetTasks({
    IdSheet: sheet.Id,
    Tasks: tasks.map(t => t.Id)
  });

  const sh = homeworks.Sheets.get(sheet.Id)!;
  sh.Tasks = tasks;
}

async function setSheetFavorite(sheet: Sheet) {
  sheet.Anonymous = { Valid: false, ID: 0 };
  const res = await controller.HomeworkUpdateSheet(sheet);
  if (res === undefined) return;

  homeworks.Sheets.get(sheet.Id)!.Sheet = sheet;
}

async function updateTravail(travail: Travail) {
  const ok = await controller.HomeworkUpdateTravail(travail);
  if (ok == undefined) return;

  const cl = travauxByClassroom(travail.IdClassroom)!;
  const index = cl.Travaux!.findIndex(tr => tr.Id == travail.Id);
  cl.Travaux![index] = travail;
}

async function copyTravailTo(tr: Travail, idClassroom: number) {
  const res = await controller.HomeworkCopyTravail({
    IdTravail: tr.Id,
    IdClassroom: idClassroom
  });
  if (res === undefined) return;

  if (res.HasNewSheet) {
    homeworks.Sheets.set(res.NewSheet.Sheet.Id, res.NewSheet);
  }
  const cl = travauxByClassroom(idClassroom)!;
  cl.Travaux = (cl.Travaux || []).concat(res.Travail);
}

async function deleteTravail(tr: Travail) {
  const res = await controller.HomeworkDeleteTravail({ id: tr.Id });
  if (res === undefined) return;
  const cl = travauxByClassroom(tr.IdClassroom)!;
  cl.Travaux = (cl.Travaux || []).filter(t => t.Id != tr.Id);
}
</script>

<style>
:deep(.v-dialog .v-overlay__content) {
  max-width: 900px;
}
</style>
