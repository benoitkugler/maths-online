<template>
  <v-card
    class="mx-auto pa-1"
    :title="props.classroom.name"
    subtitle="Liste des élèves"
  >
    <v-dialog
      :model-value="studentToDelete != null"
      @update:model-value="studentToDelete = null"
      :retain-focus="false"
      max-width="800px"
    >
      <v-card title="Confirmer la suppression" v-if="studentToDelete != null">
        <v-card-text>
          Confirmez-vous la suppression du profile élève
          <i>{{ studentToDelete.Surname }} {{ studentToDelete.Name }}</i> ?
          <br />
          Toute progression sera supprimée, et cette action est irréversible.
        </v-card-text>
        <v-card-actions>
          <v-btn @click="studentToDelete = null" color="warning">Retour</v-btn>
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
          <DateField
            v-model="studentToUpdate.Birthday"
            label="Date de naissance"
          >
          </DateField>
          <v-checkbox
            class="mt-4"
            density="compact"
            v-model="studentToUpdate.IsClientAttached"
            label="Application élève rattachée"
            hint="Décocher pour renouveller la procédure de rattachement, en générant un code de connection."
            persistent-hint
            :disabled="!studentToUpdate.IsClientAttached"
          ></v-checkbox>
        </v-card-text>
        <v-card-actions>
          <v-btn @click="studentToUpdate = null" color="warning">Retour</v-btn>
          <v-spacer></v-spacer>
          <v-btn @click="updateStudent" color="success">Modifier</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="showUploadFile" :retain-focus="false" max-width="700px">
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
          <v-btn @click="showUploadFile = false" color="warning">Retour</v-btn>
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

    <template v-slot:append>
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
      <v-btn class="mx-2 my-1" @click="addStudent" title="Créer un profil">
        <v-icon icon="mdi-plus" color="success"></v-icon>
        Ajouter un élève
      </v-btn>
      <v-btn
        class="mx-2"
        @click="showUploadFile = true"
        title="Ajouter des élèves à partir d'une liste Pronote"
      >
        <v-icon icon="mdi-upload" color="success"></v-icon>
        Importer
      </v-btn>

      <v-btn icon flat class="mx-2" @click="emit('closed')">
        <v-icon icon="mdi-close"></v-icon>
      </v-btn>
    </template>

    <v-card-text>
      <v-alert
        closable
        color="info"
        style="text-align: center; font-size: 24pt"
        :model-value="classroomCode != null"
        @update:model-value="classroomCode = null"
      >
        <v-chip size="30pt" class="pa-2 my-2">{{ classroomCode }}</v-chip>
      </v-alert>

      <v-list>
        <v-list-item
          v-for="(student, index) in students"
          :key="student.Student.Id"
        >
          <v-row
            no-gutters
            :class="{ 'bg-grey-lighten-3': index % 2 == 0, rounded: true }"
          >
            <v-col cols="6" sm="4" md="2" align-self="center">
              <v-btn
                class="mx-2 my-1"
                size="x-small"
                icon
                @click="studentToDelete = student.Student"
                title="Supprimer l'élève"
              >
                <v-icon icon="mdi-delete" color="red"></v-icon>
              </v-btn>
              <v-btn
                class="mx-1"
                size="x-small"
                icon
                @click="studentToUpdate = copy(student.Student)"
                title="Modifier le profil"
              >
                <v-icon icon="mdi-pencil"></v-icon>
              </v-btn>
            </v-col>
            <v-col cols="6" sm="4" md="3" align-self="center">
              {{ student.Student.Name }} {{ student.Student.Surname }}
              <v-tooltip
                style="cursor: pointer"
                v-if="student.Student.IsClientAttached"
              >
                <template v-slot:activator="{ isActive, props }">
                  <v-icon
                    v-on="{ isActive }"
                    v-bind="props"
                    color="green"
                    inline
                  >
                    mdi-check
                  </v-icon>
                </template>
                L'élève a relié son application.
              </v-tooltip>
            </v-col>
            <v-col align-self="center">
              <v-list-item-subtitle
                :style="
                  classroomCode == null
                    ? ''
                    : 'color: transparent;  text-shadow: 0 0 4px rgba(0,0,0,0.5)'
                "
                >{{ formatDate(student.Student.Birthday) }}
              </v-list-item-subtitle>
            </v-col>
            <v-col cols="2" align-self="center">
              <v-menu location="left">
                <template v-slot:activator="{ isActive, props }">
                  <v-chip color="teal" v-on="{ isActive }" v-bind="props">
                    {{ student.Success.TotalPoints }}
                    <v-icon>mdi-crown</v-icon>
                  </v-chip>
                </template>
                <DetailsSuccess :stats="student.Success"></DetailsSuccess>
              </v-menu>
              <v-chip color="orange" class="mx-2">
                {{ student.Success.Flames }}
                <v-icon>mdi-fire</v-icon>
              </v-chip>
            </v-col>
          </v-row>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { ref } from "vue";
import type { Classroom, Student, StudentExt } from "@/controller/api_gen";
import { $ref } from "vue/macros";
import { controller } from "@/controller/controller";
import { copy, formatDate } from "@/controller/utils";
import { onMounted } from "vue";
import { computed } from "vue";
import DateField from "../DateField.vue";
import DetailsSuccess from "./DetailsSuccess.vue";

interface Props {
  classroom: Classroom;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "closed"): void;
}>();

const mode = ref(0);

onMounted(() => fetchStudents());

let students = computed(() => {
  const out = _students.map(v => v);
  out.sort((s1, s2) => s1.Student.Name.localeCompare(s2.Student.Name));
  return out;
});

let _students = $ref<StudentExt[]>([]);

async function fetchStudents() {
  const res = await controller.TeacherGetClassroomStudents({
    "id-classroom": props.classroom.id
  });
  if (res == undefined) {
    return;
  }

  _students = res || [];
}

async function addStudent() {
  const res = await controller.TeacherAddStudent({
    "id-classroom": props.classroom.id
  });
  if (res == undefined) {
    return;
  }

  _students.push(res);

  studentToUpdate = copy(res.Student);
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
    "id-classroom": props.classroom.id
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
</script>

<style></style>
