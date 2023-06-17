<template>
  <v-card
    class="my-5 mx-auto"
    width="80%"
    title="Paramètres"
    subtitle="Modifier les réglages de mon profil."
  >
    <v-card-text class="ma-2">
      <v-form v-if="settings != null">
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              variant="outlined"
              density="compact"
              v-model="settings.Mail"
              label="Adresse email"
              persistent-hint
              hint="Adresse utilisée comme identifiant de connection"
            ></v-text-field>
          </v-col>

          <v-col cols="12" md="6">
            <v-text-field
              variant="outlined"
              density="compact"
              v-model="settings.Password"
              label="Mot de passe"
              :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
              :rules="[(v) => v.length >= 4 || 'Au moins 4 charactères']"
              :type="showPassword ? 'text' : 'password'"
              hint="Le mot de passe doit contenir au moins 4 charactères."
              counter=""
              @click:append="showPassword = !showPassword"
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-switch
              density="compact"
              v-model="settings.HasEditorSimplified"
              label="Mode simplifié"
              color="primary"
              messages="Simplifie la présentation de l'éditeur en masquant les champs spécifiques aux mathématiques."
            ></v-switch>
          </v-col>
        </v-row>
        <v-row>
          <v-spacer></v-spacer>
          <v-col cols="auto">
            <v-btn color="success" @click="save" :disabled="!isFormValid"
              >Enregistrer</v-btn
            >
          </v-col>
        </v-row>
      </v-form>
      <v-row v-else justify="center">
        <v-col cols="auto">
          <v-progress-circular
            indeterminate
            color="primary"
          ></v-progress-circular>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { TeacherSettings } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { computed } from "vue";
import { onMounted } from "vue";
import { $ref } from "vue/macros";

let settings = $ref<TeacherSettings | null>(null);

let showPassword = $ref(false);

const isFormValid = computed(
  () =>
    settings?.Mail.includes("@") &&
    settings?.Mail.includes(".") &&
    settings.Password.length >= 4
);

onMounted(() => fetchSettings());

async function fetchSettings() {
  const res = await controller.TeacherGetSettings();
  if (res == undefined) {
    return;
  }
  settings = res;
}

async function save() {
  if (settings == null) return;
  const ok = await controller.TeacherUpdateSettings(settings);
  if (ok) {
    controller.settings = settings;
  }
}
</script>

<style></style>
