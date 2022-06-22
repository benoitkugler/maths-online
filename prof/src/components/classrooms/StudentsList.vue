<template>
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
          variant="contained"
          :disabled="!uploadedFile.length"
        >
          Importer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="mx-auto pa-1" width="80%">
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
      <v-card width="60%" class="mx-auto">
        <v-row>
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
            <v-list-item v-for="student in students" :key="student.Id">
              <v-row no-gutters>
                <v-col> {{ student.Name }} {{ student.Surname }} </v-col>
                <v-col>
                  <v-list-item-subtitle
                    >{{ formatDate(student.Birthday) }}
                  </v-list-item-subtitle>
                </v-col>
                <v-col>
                  <v-tooltip
                    style="cursor: pointer"
                    v-if="student.IsClientAttached"
                  >
                    <template v-slot:activator="{ isActive, props }">
                      <v-badge v-on="{ isActive }" v-bind="props" color="green"
                        >.</v-badge
                      >
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
import { formatDate } from "@/controller/utils";
import { onMounted } from "vue";
import { $ref } from "vue/macros";

interface Props {
  classroom: Classroom;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "closed"): void;
}>();

onMounted(() => fetchStudents());

let students = $ref<Student[]>([]);

async function fetchStudents() {
  const res = await controller.TeacherGetClassroomStudents({
    "id-classroom": props.classroom.id,
  });
  if (res == undefined) {
    return;
  }
  students = res || [];
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
