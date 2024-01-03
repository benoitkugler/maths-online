<template>
  <v-dialog
    :model-value="showInscriptionValidated"
    @update:model-value="removeInscriptionValidated"
    max-width="600px"
  >
    <v-card title="Inscription validée" color="success">
      <v-card-text>
        Votre inscription a bien été validée. <br />
        Vous pouvez vous connecter avec vos nouveaux identifiants.
      </v-card-text>
    </v-card>
  </v-dialog>

  <v-dialog v-model="showResetDone" max-width="600px">
    <v-card>
      <v-card-title class="bg-info"
        >Réinitialisation du mot de passe</v-card-title
      >
      <v-card-text>
        Un mail contenant votre nouveau mot de passe a été envoyé à l'adresse
        <div style="text-align: center">
          <i>{{ mail }}</i>
        </div>
        Vous pourrez le modifier via le pannel de réglages de votre compte.
      </v-card-text>
    </v-card>
  </v-dialog>

  <v-dialog v-model="showSuccessInscription" max-width="600px">
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

  <v-row class="my-1 mx-6 pb-3 fill-height" justify="center">
    <v-col
      cols="12"
      sm="6"
      align-self="center"
      class="d-none d-sm-flex"
      v-if="mode == 'connection'"
    >
      <v-card>
        <v-card-title class="bg-secondary rounded">
          Bienvenue sur Isyro
        </v-card-title>
        <v-card-text class="py-3"
          >Isyro est une plateforme pédagogique pour concevoir des exercices
          interactifs, à utiliser en classe comme à la maison.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="secondary"
            variant="elevated"
            @click="
              mode = 'inscription';
              showPassword = true;
            "
          >
            S'inscrire
          </v-btn>
          <v-spacer></v-spacer>
        </v-card-actions>
      </v-card>
    </v-col>

    <v-col cols="12" sm="6" align-self="center">
      <v-card>
        <v-card-title class="bg-primary rounded">
          {{ mode == "inscription" ? "S'inscrire" : "  Se connecter" }}
        </v-card-title>
        <v-progress-linear
          indeterminate
          v-show="isLoading"
          color="primary"
        ></v-progress-linear>
        <v-form
          class="px-3 mt-4"
          @keyup.enter="mode == 'inscription' ? inscription() : connection()"
        >
          <v-row>
            <v-col>
              <v-text-field
                density="comfortable"
                variant="outlined"
                label="Mail"
                v-model="mail"
                type="email"
                name="email"
                :hint="
                  mode == 'inscription'
                    ? 'Adresse utiilisée comme identifiant'
                    : ''
                "
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
                density="comfortable"
                variant="outlined"
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
                :hide-details="mode == 'inscription'"
              ></v-text-field>
            </v-col>
          </v-row>
          <v-row v-if="mode == 'inscription'">
            <v-col>
              <v-select
                density="compact"
                variant="outlined"
                v-model="favoriteMatiere"
                color="primary"
                label="Matière principale"
                :items="Object.keys(MatiereTagLabels)"
              ></v-select>
            </v-col>
          </v-row>
        </v-form>
        <v-card-actions class="mt-2">
          <v-btn v-if="mode == 'inscription'" @click="mode = 'connection'"
            >Retour</v-btn
          >
          <v-btn
            v-else
            v-show="error.Error != ''"
            :disabled="!isEmailValid"
            @click="resetPassword"
          >
            Mot de passe oublié ?
          </v-btn>
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            variant="elevated"
            :disabled="!areCredencesValid"
            @click="mode == 'inscription' ? inscription() : connection()"
          >
            {{ mode == "inscription" ? "S'inscrire" : "  Se connecter" }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import {
  MatiereTag,
  type AskInscriptionOut,
  MatiereTagLabels,
} from "@/controller/api_gen";
import { controller, isInscriptionValidated } from "@/controller/controller";
import { ref, computed } from "vue";

const emit = defineEmits<{
  (e: "loggin"): void;
}>();

const showInscriptionValidated = ref(isInscriptionValidated());
function removeInscriptionValidated() {
  window.location.search = "";
}

const mode = ref<"inscription" | "connection">("connection");

const mail = ref("");
const password = ref("");
const favoriteMatiere = ref<MatiereTag>(MatiereTag.Mathematiques);
const showPassword = ref(false);
const error = ref<AskInscriptionOut>({ Error: "", IsPasswordError: false });
const showSuccessInscription = ref(false);
const isLoading = ref(false);

const areCredencesValid = computed(
  () => !isLoading.value && isEmailValid && password.value != ""
);

const isEmailValid = computed(
  () => mail.value.includes("@") && mail.value.includes(".")
);

async function inscription() {
  if (!areCredencesValid.value) return;

  isLoading.value = true;
  const res = await controller.AskInscription({
    Mail: mail.value,
    Password: password.value,
    FavoriteMatiere: favoriteMatiere.value,
  });
  isLoading.value = false;
  if (res == undefined) {
    return;
  }
  error.value = res;
  if (error.value.Error == "") {
    showSuccessInscription.value = true;
  }
}

async function connection() {
  if (!areCredencesValid.value) return;

  const res = await controller.Loggin({
    Mail: mail.value,
    Password: password.value,
  });
  if (res == undefined) {
    return;
  }

  error.value = res;
  if (error.value.Error) return;

  const settings = await controller.TeacherGetSettings();
  if (settings) controller.settings = settings;
  emit("loggin");
}

const showResetDone = ref(false);
async function resetPassword() {
  const res = await controller.TeacherResetPassword({ mail: mail.value });
  if (res == undefined) return;
  showResetDone.value = true;
}
</script>

<style></style>
