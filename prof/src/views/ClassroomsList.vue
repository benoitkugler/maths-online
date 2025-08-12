<template>
  <v-dialog
    :model-value="classroomToShow != null"
    @update:model-value="classroomToShow = null"
    fullscreen
    :retain-focus="false"
  >
    <classroom-view
      v-if="classroomToShow"
      :classroom="classroomToShow"
      @closed="
        fetchClassrooms();
        classroomToShow = null;
      "
    ></classroom-view>
  </v-dialog>

  <v-dialog
    :model-value="classroomToDelete != null"
    @update:model-value="classroomToDelete = null"
    max-width="600px"
  >
    <v-card title="Confirmer" v-if="classroomToDelete">
      <v-card-text
        >Etes-vous certain de vouloir supprimer la classe
        <i>{{ classroomToDelete.Classroom.name }}</i> ? <br /><br />

        <div v-if="classroomToDelete.SharedWith?.length">
          Comme cette classe est partagée, vous serez retiré de ses enseignants,
          mais les élèves associés seront conservés.
        </div>
        <div v-else>
          Tous les élèves associés (et leur progression) seront supprimés.
          <br />Cette opération est irréversible.
        </div>
      </v-card-text>
      <v-card-actions>
        <v-btn @click="classroomToDelete = null">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn color="red" @click="deleteClassroom" variant="outlined">
          Supprimer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    :model-value="classroomToUpdate != null"
    @update:model-value="classroomToUpdate = null"
    max-width="600px"
  >
    <v-card title="Modifier la classe" v-if="classroomToUpdate != null">
      <v-card-text class="my-2">
        <v-row>
          <v-col>
            <v-text-field
              label="Nom"
              v-model="classroomToUpdate.name"
              variant="outlined"
              density="compact"
              hide-details
            ></v-text-field
          ></v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-select
              label="Seuil de la dernière guilde"
              hint="Nombre de points nécessaires pour débloquer la dernière guilde."
              persistent-hint
              v-model="classroomToUpdate.MaxRankThreshold"
              variant="outlined"
              density="compact"
              :items="maxRankItems"
            ></v-select>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-btn @click="classroomToUpdate = null" color="warning">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn color="green" @click="updateClassroom" variant="outlined">
          Enregistrer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="my-5 mx-auto" width="90%">
    <v-row>
      <v-col md="9" sm="6">
        <v-card-title>Mes classes</v-card-title>
        <v-card-subtitle>Gérer mes classes et mes élèves.</v-card-subtitle>
      </v-col>

      <v-col align-self="center" style="text-align: right" md="3" sm="6">
        <v-btn
          class="mr-4"
          @click="createClassroom"
          title="Ajouter une nouvelle classe"
        >
          <v-icon icon="mdi-plus" color="success"></v-icon>
          Créer une classe
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text>
      <v-row no-gutters>
        <v-col
          md="4"
          sm="6"
          xs="12"
          v-for="classroom in classrooms"
          :key="classroom.Classroom.id"
        >
          <classroom-card
            :classroom="classroom"
            @show-students="classroomToShow = classroom.Classroom"
            @delete="classroomToDelete = classroom"
            @update="classroomToUpdate = copy(classroom.Classroom)"
            @invite="(mail) => inviteTeacher(classroom.Classroom, mail)"
          ></classroom-card>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { Classroom, ClassroomExt } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { copy } from "@/controller/utils";
import { ref, onMounted } from "vue";
import ClassroomCard from "../components/classrooms/ClassroomCard.vue";
import ClassroomView from "@/components/classrooms/ClassroomView.vue";

const classrooms = ref<ClassroomExt[]>([]);

onMounted(() => fetchClassrooms());

async function fetchClassrooms() {
  const res = await controller.TeacherGetClassrooms();
  if (res == undefined) {
    return;
  }
  classrooms.value = res || [];
}

async function createClassroom() {
  const res = await controller.TeacherCreateClassroom();
  if (res === undefined) return;
  controller.showMessage("Classe créée avec succès.");

  await fetchClassrooms();
}

const classroomToDelete = ref<ClassroomExt | null>(null);
async function deleteClassroom() {
  if (classroomToDelete.value == null) {
    return;
  }
  const res = await controller.TeacherDeleteClassroom({
    id: classroomToDelete.value.Classroom.id,
  });
  if (res === undefined) return;
  controller.showMessage("Classe supprimée avec succès.");

  classroomToDelete.value = null;
  await fetchClassrooms();
}

const classroomToUpdate = ref<Classroom | null>(null);
async function updateClassroom() {
  if (classroomToUpdate.value == null) {
    return;
  }
  const res = await controller.TeacherUpdateClassroom(classroomToUpdate.value);
  classroomToUpdate.value = null;
  if (res == undefined) {
    return;
  }
  controller.showMessage("Classe mise à jour avec succès.");

  const index = classrooms.value.findIndex((cl) => cl.Classroom.id == res.id);
  classrooms.value[index].Classroom = res;
}

const classroomToShow = ref<Classroom | null>(null);

const maxRankItems = [
  { title: "15 000 : plus tranquille", value: 15_000 },
  { title: "25 000 : normal", value: 25_000 },
  { title: "40 000 : intense", value: 40_000 },
];

async function inviteTeacher(classroom: Classroom, mail: string) {
  const res = await controller.TeacherInviteTeacherToClassroom({
    IdClassroom: classroom.id,
    MailToInvite: mail,
  });
  if (res === undefined) return;
  controller.showMessage("Collègue ajouté avec succès.");
  await fetchClassrooms();
}
</script>

<style>
:deep(.v-dialog .v-overlay__content) {
  max-width: 900px;
}
</style>
