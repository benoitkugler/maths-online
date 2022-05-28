<template>
  <v-dialog v-model="showInscriptionValidated">
    <v-card title="Inscription validée" color="success">
      <v-card-text>
        Votre inscription a bien été validée. <br />
        Vous pouvez vous connecter avec vos nouveaux identifiants.
      </v-card-text>
    </v-card>
  </v-dialog>

  <v-dialog v-model="showSuccessInscription">
    <v-card title="Inscription enregistrée" color="success">
      <v-card-text>
        Merci pour votre inscription ! <br />
        Un mail de confirmation vous a été envoyé à l'adresse <i>{{ mail }}</i
        >. <br />
        Merci de suivre le lien présent dans le mail pour valider votre
        inscription.
      </v-card-text>
    </v-card>
  </v-dialog>

  <v-row class="my-1 mx-6 fill-height">
    <v-col cols="6" align-self="center">
      <v-card>
        <v-card-title class="bg-secondary rounded"> Se connecter </v-card-title>
        <v-progress-linear
          indeterminate
          v-show="isLoading"
          color="secondary"
        ></v-progress-linear>
        <v-form class="px-3 mt-4">
          <v-row>
            <v-col>
              <v-text-field
                label="Mail"
                v-model="mail"
                type="email"
                name="email"
                required
                :error="error.Error != '' && !error.IsPasswordError"
                :error-messages="
                  error.Error != '' && !error.IsPasswordError
                    ? [error.Error]
                    : ''
                "
              ></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field
                label="Mot de passe"
                v-model="password"
                :append-inner-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                :type="showPassword ? 'text' : 'password'"
                name="password"
                @click:append-inner="showPassword = !showPassword"
                :error="error.Error != '' && error.IsPasswordError"
                :error-messages="
                  error.Error != '' && error.IsPasswordError
                    ? [error.Error]
                    : ''
                "
              ></v-text-field>
            </v-col>
          </v-row>
        </v-form>
        <v-card-actions>
          <v-btn
            color="primary"
            variant="contained"
            :disabled="!areCredencesValid"
            @click="inscription"
          >
            S'inscrire
          </v-btn>
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            variant="contained"
            :disabled="!areCredencesValid"
            @click="connection"
          >
            Entrer
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-col>
    <v-col cols="6" align-self="center">
      <v-card>
        <v-card-title class="bg-primary rounded">
          Bienvenue sur Isyro
        </v-card-title>
        <v-card-text
          >Isyro est une plateforme pédagogique pour concevoir des exercices
          interactifs, à utiliser en classe comme à la maison.
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { AskInscriptionOut } from "@/controller/api_gen";
import { controller, ShowSuccessInscription } from "@/controller/controller";
import { computed } from "vue";
import { $ref } from "vue/macros";

const emit = defineEmits<{
  (e: "loggin"): void;
}>();

let showInscriptionValidated = $ref(ShowSuccessInscription);

let mail = $ref("");
let password = $ref("");
let showPassword = $ref(false);
let error = $ref<AskInscriptionOut>({ Error: "", IsPasswordError: false });
let showSuccessInscription = $ref(false);
let isLoading = $ref(false);

const areCredencesValid = computed(
  () => !isLoading && mail != "" && password != ""
);

async function inscription() {
  isLoading = true;
  const res = await controller.AskInscription({
    Mail: mail,
    Password: password,
  });
  isLoading = false;
  if (res == undefined) {
    return;
  }
  error = res;
  if (error.Error == "") {
    showSuccessInscription = true;
  }
}

async function connection() {
  const res = await controller.Loggin({
    Mail: mail,
    Password: password,
  });
  if (res == undefined) {
    return;
  }
  error = res;
  if (error.Error == "") {
    emit("loggin");
  }
}
</script>

<style></style>
