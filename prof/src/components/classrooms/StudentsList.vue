<template>
  <v-dialog
    :model-value="studentToDelete != null"
    @update:model-value="studentToDelete = null"
    :retain-focus="false"
    max-width="800px"
  >
    <v-card title="Confirmer la suppression" v-if="studentToDelete != null">
      <v-card-text>
        Confirmez-vous la suppression du profile élève
        <i>{{ studentToDelete.Surname }} {{ studentToDelete.Name }}</i> ? <br />
        Toute progression sera supprimée, et cette action est irréversible.
      </v-card-text>
      <v-card-actions>
        <v-btn @click="studentToDelete = null">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn @click="deleteStudent" color="red">Supprimer le profile</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    :model-value="studentToUpdate != null"
    @update:model-value="studentToUpdate = null"
    :retain-focus="false"
    max-width="600"
  >
    <v-card title="Modifier le profil" v-if="studentToUpdate != null">
      <v-card-text>
        <v-text-field
          variant="outlined"
          density="compact"
          label="Nom"
          v-model="studentToUpdate.Name"
        ></v-text-field>
        <v-text-field
          variant="outlined"
          density="compact"
          label="Prénom"
          v-model="studentToUpdate.Surname"
        ></v-text-field>
        <date-field
          v-model="studentToUpdate.Birthday"
          label="Date de naissance"
        >
        </date-field>
      </v-card-text>
      <v-card-actions>
        <v-btn @click="studentToUpdate = null">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn @click="updateStudent" color="success">Modifier</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="showUploadFile" :retain-focus="false">
    <v-card
      title="Importer une liste"
      subtitle="Les élèves seront ajoutés à la liste courante."
    >
      <v-card-text>
        <v-file-input
          label="Liste"
          hint="Fichier CSV généré par Pronote."
          v-model="uploadedFile"
          show-size
          :multiple="false"
          variant="underlined"
          persistent-hint
        >
        </v-file-input>
      </v-card-text>
      <v-card-actions>
        <v-btn @click="showUploadFile = false">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn
          color="success"
          @click="importStudents"
          variant="text"
          :disabled="!uploadedFile.length"
        >
          Importer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="mx-auto pa-1">
    <v-row>
      <v-col md="9" sm="6">
        <v-card-title
          >Elèves de la classe {{ props.classroom.name }}</v-card-title
        >
      </v-col>
      <v-col style="text-align: right">
        <v-btn icon flat class="mx-2" @click="emit('closed')">
          <v-icon icon="mdi-close"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-card-text>
      <v-card width="90%" class="mx-auto">
        <v-row class="mt-1">
          <v-col>
            <v-card-subtitle>Elèves</v-card-subtitle>
          </v-col>
          <v-col align-self="center" style="text-align: right" md="8" sm="6">
            <v-btn
              class="mx-2"
              @click="generateClassroomCode"
              title="Générer un code pour rattacher des élèves à la classe"
              variant="text"
              :disabled="!students.length || classroomCode != null"
            >
              Code de connection
            </v-btn>
            <v-divider vertical></v-divider>
            <v-btn
              class="mx-2 my-1"
              @click="addStudent"
              title="Créer un profil"
            >
              <v-icon icon="mdi-plus" color="success"></v-icon>
              Créer
            </v-btn>
            <v-btn
              class="mx-2"
              @click="showUploadFile = true"
              title="Ajouter des élèves à partir d'une liste Pronote"
            >
              <v-icon icon="mdi-upload" color="success"></v-icon>
              Importer
            </v-btn>
          </v-col>
        </v-row>
        <v-card-text>
          <v-alert
            closable
            color="info"
            style="text-align: center; font-size: 18pt"
            :model-value="classroomCode != null"
            @update:model-value="classroomCode = null"
          >
            <v-chip size="18pt" class="pa-2">{{ classroomCode }}</v-chip>
          </v-alert>

          <v-list>
            <v-list-item v-for="(student, index) in students" :key="student.Id">
              <v-row
                no-gutters
                :class="{ 'bg-grey-lighten-3': index % 2 == 0, rounded: true }"
              >
                <v-col cols="6" sm="4" md="" align-self="center">
                  <v-btn
                    class="mx-2 my-1"
                    size="x-small"
                    icon
                    @click="studentToDelete = student"
                    title="Supprimer l'élève"
                  >
                    <v-icon icon="mdi-delete" color="red"></v-icon>
                  </v-btn>
                  <v-btn
                    class="mx-2"
                    size="x-small"
                    icon
                    @click="studentToUpdate = copy(student)"
                    title="Modifier le profil"
                  >
                    <v-icon icon="mdi-pencil"></v-icon>
                  </v-btn>
                </v-col>
                <v-col cols="6" sm="" align-self="center">
                  {{ student.Name }} {{ student.Surname }}
                </v-col>
                <v-col align-self="center">
                  <v-list-item-subtitle
                    >{{ formatDate(student.Birthday) }}
                  </v-list-item-subtitle>
                </v-col>
                <v-col align-self="center">
                  <v-tooltip
                    style="cursor: pointer"
                    v-if="student.IsClientAttached"
                  >
                    <template v-slot:activator="{ isActive, props }">
                      <v-badge
                        v-on="{ isActive }"
                        v-bind="props"
                        color="green"
                        inline
                      >
                      </v-badge>
                    </template>
                    L'élève a reliée son application.
                  </v-tooltip>
                </v-col>
              </v-row>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { Classroom, Student } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { copy, formatDate } from "@/controller/utils";
import { computed, onMounted } from "vue";
import { $ref } from "vue/macros";
import DateField from "../DateField.vue";

interface Props {
  classroom: Classroom;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "closed"): void;
}>();

onMounted(() => fetchStudents());

let students = computed(() => {
  const out = _students.map((v) => v);
  out.sort((s1, s2) => s1.Name.localeCompare(s2.Name));
  return out;
});

let _students = $ref<Student[]>([]);

async function fetchStudents() {
  const res = await controller.TeacherGetClassroomStudents({
    "id-classroom": props.classroom.id,
  });
  if (res == undefined) {
    return;
  }
  _students = res || [];
}

async function addStudent() {
  const res = await controller.TeacherAddStudent({
    "id-classroom": props.classroom.id,
  });
  if (res == undefined) {
    return;
  }
  console.log(res);

  _students.push(res);

  studentToUpdate = copy(res);
}

let studentToDelete = $ref<Student | null>(null);
async function deleteStudent() {
  if (studentToDelete == null) {
    return;
  }
  await controller.TeacherDeleteStudent({ "id-student": studentToDelete.Id });
  studentToDelete = null;
  await fetchStudents();
}

let studentToUpdate = $ref<Student | null>(null);
async function updateStudent() {
  if (studentToUpdate == null) {
    return;
  }
  await controller.TeacherUpdateStudent(studentToUpdate);
  studentToUpdate = null;
  await fetchStudents();
}

let showUploadFile = $ref(false);
let uploadedFile = $ref<File[]>([]);
async function importStudents() {
  showUploadFile = false;
  if (uploadedFile.length == 0) {
    return;
  }
  await controller.TeacherImportStudents(
    { "id-classroom": String(props.classroom.id) },
    uploadedFile[0]
  );
  uploadedFile = [];
  await fetchStudents();
}

let classroomCode = $ref<string | null>(null);
async function generateClassroomCode() {
  const res = await controller.TeacherGenerateClassroomCode({
    "id-classroom": props.classroom.id,
  });
  console.log(res);

  if (res == undefined) {
    return;
  }
  classroomCode = res.Code;

  const timeout = 5000;
  const refresh = async () => {
    if (classroomCode != null) {
      await fetchStudents();
      setTimeout(refresh, timeout);
    }
  };
  setTimeout(refresh, timeout);
}

// async function createClassroom() {
//   await controller.TeacherCreateClassroom(null);
//   await fetchClassrooms();
// }

// let classroomToDelete = $ref<Classroom | null>(null);
// async function deleteClassroom() {
//   if (classroomToDelete == null) {
//     return;
//   }
//   await controller.TeacherDeleteClassroom({ id: classroomToDelete.id });
//   classroomToDelete = null;
//   await fetchClassrooms();
// }

// let classroomToUpdate = $ref<Classroom | null>(null);
// async function updateClassroom() {
//   if (classroomToUpdate == null) {
//     return;
//   }
//   const res = await controller.TeacherUpdateClassroom(classroomToUpdate);
//   classroomToUpdate = null;
//   if (res == undefined) {
//     return;
//   }

//   const index = classrooms.findIndex((cl) => cl.Classroom.id == res.id);
//   classrooms[index].Classroom = res;
// }
</script>

<style>
:deep(.v-dialog .v-overlay__content) {
  max-width: 900px;
}
</style>
