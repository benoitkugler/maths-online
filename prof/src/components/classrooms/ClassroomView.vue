<template>
  <v-card
    class="mx-auto pa-1"
    :title="props.classroom.name"
    subtitle="Liste des élèves"
  >
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
      <v-divider vertical></v-divider>
      <v-btn
        class="mx-2"
        :href="
          controller.TeacherExportStudentsAdvance(
            props.classroom.id,
            controller.getToken()
          )
        "
      >
        <v-icon icon="mdi-download" color="success"></v-icon>
        Exporter
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
                <v-icon icon="mdi-cog"></v-icon>
              </v-btn>
            </v-col>
            <v-col cols="6" sm="4" md="3" align-self="center">
              {{ student.Student.Name }} {{ student.Student.Surname }}
              <v-tooltip
                style="cursor: pointer"
                v-if="student.Student.Clients?.length"
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
            <v-col cols="4" align-self="center">
              <v-row justify="center">
                <v-col>
                  <v-tooltip>
                    <template v-slot:activator="{ isActive, props }">
                      <v-chip v-on="{ isActive }" v-bind="props">
                        {{ rankLabels[student.Success.Rank] }}
                      </v-chip>
                    </template>
                    Guilde {{ student.Success.Rank }} /
                    {{ rankLabels.length - 1 }}
                  </v-tooltip>
                </v-col>
                <v-col cols="auto">
                  <v-menu location="left">
                    <template v-slot:activator="{ isActive, props }">
                      <v-chip
                        color="teal"
                        class="mx-1"
                        v-on="{ isActive }"
                        v-bind="props"
                      >
                        {{ student.Success.TotalPoints }}
                        <v-icon>mdi-crown</v-icon>
                      </v-chip>
                    </template>
                    <DetailsSuccess :stats="student.Success"></DetailsSuccess>
                  </v-menu>
                </v-col>
                <v-col cols="auto">
                  <v-chip color="orange" class="mx-1">
                    {{ student.Success.Flames }}
                    <v-icon>mdi-fire</v-icon>
                  </v-chip>
                </v-col>
              </v-row>
            </v-col>
          </v-row>
        </v-list-item>
      </v-list>
    </v-card-text>

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
          <v-list>
            <v-list-subheader>
              {{
                studentToUpdate.Clients?.length
                  ? "Appareils connectés"
                  : "Aucun appareil connecté"
              }}
            </v-list-subheader>
            <v-list-item
              v-for="(client, i) in studentToUpdate.Clients"
              :key="i"
              :title="client.Device || 'Appareil inconnu'"
              :subtitle="`le ${formatTime(client.Time, true)}`"
            ></v-list-item>
          </v-list>
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
  </v-card>
</template>

<script setup lang="ts">
import type { Classroom, Student, StudentExt } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { copy, formatDate, formatTime } from "@/controller/utils";
import { onMounted, onUnmounted, ref } from "vue";
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

onMounted(() => fetchStudents());
onUnmounted(() => clearTimeout(timerId));

const students = computed(() => {
  const out = _students.value.map((v) => v);
  out.sort((s1, s2) => {
    const c1 = s1.Student.Name.localeCompare(s2.Student.Name);
    return c1 == 0 ? s1.Student.Surname.localeCompare(s2.Student.Surname) : c1;
  });
  return out;
});

const _students = ref<StudentExt[]>([]);

async function fetchStudents() {
  const res = await controller.TeacherGetClassroomStudents({
    "id-classroom": props.classroom.id,
  });
  if (res === undefined) {
    return;
  }

  _students.value = res || [];
}

async function addStudent() {
  const res = await controller.TeacherAddStudent({
    "id-classroom": props.classroom.id,
  });
  if (res == undefined) {
    return;
  }
  controller.showMessage("Elève ajouté avec succès.");

  _students.value.push(res);

  studentToUpdate.value = copy(res.Student);
}

const studentToDelete = ref<Student | null>(null);
async function deleteStudent() {
  if (studentToDelete.value == null) {
    return;
  }
  const res = await controller.TeacherDeleteStudent({
    "id-student": studentToDelete.value.Id,
  });
  if (res === undefined) return;
  controller.showMessage("Profil supprimé avec succès.");

  studentToDelete.value = null;
  await fetchStudents();
}

const studentToUpdate = ref<Student | null>(null);
async function updateStudent() {
  if (studentToUpdate.value == null) {
    return;
  }
  const res = await controller.TeacherUpdateStudent(studentToUpdate.value);
  if (res === undefined) return;
  controller.showMessage("Profil mis à jour avec succès.");

  studentToUpdate.value = null;
  await fetchStudents();
}

const showUploadFile = ref(false);
const uploadedFile = ref<File[]>([]);
async function importStudents() {
  showUploadFile.value = false;
  if (uploadedFile.value.length == 0) {
    return;
  }
  const res = await controller.TeacherImportStudents(
    { "id-classroom": String(props.classroom.id) },
    uploadedFile.value[0]
  );
  if (res === undefined) return;
  controller.showMessage("Liste importée avec succès.");

  uploadedFile.value = [];
  await fetchStudents();
}

const classroomCode = ref<string | null>(null);
let timerId: ReturnType<typeof setTimeout>;
async function generateClassroomCode() {
  const res = await controller.TeacherGenerateClassroomCode({
    "id-classroom": props.classroom.id,
  });
  if (res == undefined) {
    return;
  }
  classroomCode.value = res.Code;

  const timeout = 5000;
  const refresh = async () => {
    if (classroomCode.value != null) {
      await fetchStudents();
      timerId = setTimeout(refresh, timeout);
    }
  };
  timerId = setTimeout(refresh, timeout);
}

const rankLabels = [
  "Novice",
  "Guilde de Pythagore",
  "Guilde de Thalès",
  "Guilde d'Al-Kashi",
  "Guilde de Newton",
  "Guilde de Gauss",
  "Guilde d'Einstein",
];
</script>

<style></style>
