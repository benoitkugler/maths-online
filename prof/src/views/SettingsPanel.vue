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
              hint="Adresse utilisée comme identifiant de connection."
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
            <v-text-field
              variant="outlined"
              density="compact"
              label="Contact"
              v-model="settings.Contact.Name"
              hint="Affiché aux élèves. Laisser vide pour utiliser l'adresse mail."
              persistent-hint
            ></v-text-field>
          </v-col>
          <v-col>
            <v-text-field
              variant="outlined"
              density="compact"
              label="URL"
              v-model="settings.Contact.URL"
              placeholder="https://"
              hint="Optionnelle. La fournir pour afficher un lien."
              persistent-hint
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-select
              density="compact"
              variant="outlined"
              label="Matière principale"
              hint="Les ressources sont filtrées par défaut en utilisant cette matière."
              persistent-hint
              :items="Object.keys(MatiereTagLabels)"
              v-model="settings.FavoriteMatiere"
            ></v-select>
          </v-col>
          <v-col>
            <v-checkbox
              density="compact"
              :model-value="!settings.HasEditorSimplified"
              @update:model-value="
                (b) => {
                  if (settings) settings.HasEditorSimplified = !b;
                }
              "
              label="Mode scientifique"
              color="primary"
              messages="Décocher pour masquer les champs spécifiques aux mathématiques dans l'éditeur."
            ></v-checkbox>
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
import { MatiereTagLabels, type TeacherSettings } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { ref, computed, onMounted } from "vue";

const settings = ref<TeacherSettings | null>(null);

const showPassword = ref(false);

const isFormValid = computed(
  () =>
    settings.value?.Mail.includes("@") &&
    settings.value?.Mail.includes(".") &&
    settings.value.Password.length >= 4
);

onMounted(() => fetchSettings());

async function fetchSettings() {
  const res = await controller.TeacherGetSettings();
  if (res == undefined) {
    return;
  }
  settings.value = res;
}

async function save() {
  if (settings.value == null) return;
  const ok = await controller.TeacherUpdateSettings(settings.value);
  if (ok) {
    controller.settings = settings.value;
  }
}
</script>

<style></style>
